// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
)

type StatusCode uint32

const (
	weightSplit     = "@"
	datacenterSplit = "#"

	// DefaultSlice means default slice for namespace
	DefaultSlice       = "slice-0"
	PingPeriod   int64 = 4

	StatusUp             StatusCode = 1
	StatusDown           StatusCode = 0
	LocalSlaveReadClosed            = 0 // 不启用本地优先
	LocalSlaveReadPrefer            = 1 // 优先本地，但无本地时降级使用全局
	LocalSlaveReadForce             = 2 // 强制只用本地
	CheckRepeat                     = 3
)

const (
	NormalUser    = "normal user"
	MonitorUser   = "monitor user"
	StatisticUser = "statistic user"
)

const (
	DirectMaster        int = 0
	DirectSlave         int = 1
	SlaveFallbackMaster int = 2
)

func (s *StatusCode) String() string {
	r := "StatusUp"
	if *s == StatusDown {
		r = "StatusDown"
	}

	return r
}

type DBInfo struct {
	sync.Mutex
	Nodes          []*NodeInfo // 节点信息
	LocalBalancer  *balancer   // 仅包含数据中心与代理匹配的节点
	RemoteBalancer *balancer   // 包含除本机外的所有其他数据节点
	GlobalBalancer *balancer   // 包含所有节点
}

func (dbInfo *DBInfo) InitFuseRecoveryPolicy(fuseWindowSize int64, fuseMinErrorCount int64, fuseCooldownPeriod int64) error {
	// 参数校验
	if fuseWindowSize < 0 {
		return fmt.Errorf("invalid fuse window size: %d (must > 0)", fuseWindowSize)
	}

	// 为每个节点初始化独立的滑动窗口
	for _, node := range dbInfo.Nodes {
		if node == nil {
			continue // 跳过空节点
		}

		// 创建默认的熔断策略（滑动窗口熔断）（参数不合法时会返回 disabled 状态）
		node.FuseStrategy = NewSlidingWindow(int64(fuseWindowSize), fuseMinErrorCount)

		// 不禁用熔断必须要使用恢复策略
		if fuseCooldownPeriod > 0 {
			// 使用硬熔断冷却期恢复
			node.RecoveryStrategy = NewHardCoolDown(fuseCooldownPeriod)
		} else {
			// 采用默认的渐进式恢复策略
			node.RecoveryStrategy = NewGradualRecovery()
		}
	}
	return nil
}

// GetNode 根据索引返回 `NodeInfo` 结构
func (d *DBInfo) GetNode(index int) (*NodeInfo, error) {
	if index < 0 || index >= len(d.Nodes) {
		return nil, fmt.Errorf("index %d out of range", index)
	}
	return d.Nodes[index], nil
}

type IndexWeightList struct {
	Indices []int
	Weights []int
}

func newIndexWeightList() *IndexWeightList {
	return &IndexWeightList{}
}

// getIndicesAndWeights 返回 (本地数据中心的索引+权重, 远程数据中心的索引+权重, 全部索引+权重) 权重为0的节点将不会被记录
func (dbi *DBInfo) getIndicesAndWeights(proxyDatacenter string) (*IndexWeightList, *IndexWeightList, *IndexWeightList) {
	local := newIndexWeightList()
	remote := newIndexWeightList()
	global := newIndexWeightList()

	for idx, node := range dbi.Nodes {
		// 过滤掉权重为 0 以及权重为负数的节点
		if node.Weight <= 0 {
			continue
		}
		// 全局列表存储权重不为 0 的节点
		global.Indices = append(global.Indices, idx)
		global.Weights = append(global.Weights, node.Weight)

		// 按照数据中心划分
		if node.Datacenter == proxyDatacenter {
			local.Indices = append(local.Indices, idx)
			local.Weights = append(local.Weights, node.Weight)
		} else {
			remote.Indices = append(remote.Indices, idx)
			remote.Weights = append(remote.Weights, node.Weight)
		}
	}

	return local, remote, global
}

// Slice means one slice of the mysql cluster
type Slice struct {
	Cfg models.Slice
	sync.RWMutex

	Namespace      string
	Master         *DBInfo // 业务主库
	Slave          *DBInfo // 业务从库
	StatisticSlave *DBInfo // 业务统计用从库(StatisticUser)

	MonitorMaster *DBInfo // 监控用户的主库
	MonitorSlave  *DBInfo // 监控用户的从库

	ProxyDatacenter string
	charset         string
	collationID     mysql.CollationID
	HealthCheckSql  string

	HandshakeTimeout            time.Duration
	FallbackToMasterOnSlaveFail string // 控制从库获取失败时是否回退到主库

	// 熔断基础参数
	FuseEnabled        string // "ON"|"OFF" 熔断总开关
	FuseWindowSize     int64  // 滑动窗口时长（秒），建议 10-60
	FuseMinErrorCount  int64
	FuseCooldownPeriod int64 // 硬冷却的持续时间
}

type getConnError struct {
	Namespace string
	Addr      string
	Err       error
}

func NewGetConnError(err error, opts ...string) error {
	var namespace, addr string
	if len(opts) >= 1 {
		namespace = opts[0]
	}
	if len(opts) >= 2 {
		addr = opts[1]
	}
	return &getConnError{
		Namespace: namespace,
		Addr:      addr,
		Err:       err,
	}
}

func (e *getConnError) Error() string {
	if len(e.Addr) == 0 {
		return fmt.Sprintf("[ns=%s] %v", e.Namespace, e.Err)
	}
	return fmt.Sprintf("[ns=%s][addr=%s] %v", e.Namespace, e.Addr, e.Err)
}

// GetSliceName return name of slice
func (s *Slice) GetSliceName() string {
	return s.Cfg.Name
}

// InitFuseRecoveryPolicy init fuse slide window
func (s *Slice) InitFuseRecoveryPolicy(dbInfo *DBInfo) error {
	// 初始化 `FuseRecoveryPolicy`
	if err := dbInfo.InitFuseRecoveryPolicy(s.FuseWindowSize, s.FuseMinErrorCount, s.FuseCooldownPeriod); err != nil {
		return fmt.Errorf("failed to initialize fuse and recovery strategy: %w", err)
	}
	return nil
}

// GetConn get backend connection from different node based on fromSlave and userType
func (s *Slice) GetConn(reqCtx *util.RequestContext, userType int, localSlaveReadPriority int) (pc PooledConnect, err error) {
	switch userType {
	case models.MonitorUser:
		// 监控用户，走监控专用连接池
		return s.getMonitorConnection(reqCtx, localSlaveReadPriority)
	case models.StatisticUser:
		// 统计用户，走业务后端连接池
		return s.getStatisticConnection(localSlaveReadPriority)
	default:
		// 普通用户: 走业务后端连接池
		return s.getNormalConnection(reqCtx, localSlaveReadPriority)
	}
}

// getStatisticConnection get connection from statistic slave, not to get from master, Most of the SQL queries for user statistics are slow queries
func (s *Slice) getStatisticConnection(localSlaveReadPriority int) (pc PooledConnect, err error) {
	return s.handleSlaveStrict(StatisticUser, s.StatisticSlave, localSlaveReadPriority)
}

func (s *Slice) getNormalConnection(reqCtx *util.RequestContext, localSlaveReadPriority int) (pc PooledConnect, err error) {
	connectionMode := getconnectionMode(reqCtx.GetFromSlave(), s.ShouldFallbackToMasterOnSlaveFail())
	switch connectionMode {
	case DirectMaster:
		return s.handleDirectMaster(NormalUser, s.Master)
	case SlaveFallbackMaster:
		return s.handleSlaveWithFallback(NormalUser, s.Slave, s.Master, localSlaveReadPriority)
	case DirectSlave:
		return s.handleSlaveStrict(NormalUser, s.Slave, localSlaveReadPriority)
	default:
		return nil, fmt.Errorf("user:%s, invalid connection mode: %d", NormalUser, connectionMode)
	}
}

func (s *Slice) getMonitorConnection(reqCtx *util.RequestContext, localSlaveReadPriority int) (pc PooledConnect, err error) {
	connectionMode := getconnectionMode(reqCtx.GetFromSlave(), s.ShouldFallbackToMasterOnSlaveFail())
	switch connectionMode {
	case DirectMaster:
		return s.handleDirectMaster(MonitorUser, s.MonitorMaster)
	case SlaveFallbackMaster:
		return s.handleSlaveWithFallback(MonitorUser, s.MonitorSlave, s.MonitorMaster, localSlaveReadPriority)
	case DirectSlave:
		return s.handleSlaveStrict(MonitorUser, s.MonitorSlave, localSlaveReadPriority)
	default:
		return nil, fmt.Errorf("user:%s, invalid connection mode: %d", MonitorUser, connectionMode)
	}
}

// 处理直接主库连接路径
func (s *Slice) handleDirectMaster(userType string, master *DBInfo) (pc PooledConnect, err error) {
	pc, err = s.GetMasterConn(master)
	if err != nil {
		log.Warn("[direct→master] try master conn failed, user: %s, master error: %v", userType, err)
		return nil, err
	}
	return pc, nil
}

// 处理严格从库路径
func (s *Slice) handleSlaveStrict(userType string, slave *DBInfo, localSlaveReadPriority int) (pc PooledConnect, err error) {
	pc, err = s.GetSlaveConn(slave, localSlaveReadPriority)
	if err != nil {
		log.Warn("[direct→slave] try slave conn failed, no fallback, user: %s, localSlaveReadPriority: %d, slave error: %v", userType, localSlaveReadPriority, err)
		return nil, err
	}
	return pc, nil
}

// 处理从库可回退路径
func (s *Slice) handleSlaveWithFallback(userType string, slave *DBInfo, master *DBInfo, localSlaveReadPriority int) (PooledConnect, error) {
	// 尝试从库连接
	pc, slaveErr := s.GetSlaveConn(slave, localSlaveReadPriority)
	if slaveErr == nil {
		return pc, nil
	}
	// 获取主库连接
	pc, masterErr := s.GetMasterConn(master)
	if masterErr == nil {
		log.Warn("[slave→master] try slave conn failed, fallback to master success, user: %s, localSlaveReadPriority: %d, slave error: %v", userType, localSlaveReadPriority, slaveErr)
	} else {
		log.Warn("[slave→master] try slave conn failed, fallback to master failed, user: %s, localSlaveReadPriority: %d, slave error: %v, master error:%v", userType, localSlaveReadPriority, slaveErr, masterErr)
	}
	return pc, masterErr
}

func getconnectionMode(fromSlave bool, fallbackToMaster bool) (connectionMode int) {
	switch {
	case !fromSlave:
		// 第1判断条件（fromSlave=false时执行）
		return DirectMaster
	case fallbackToMaster:
		// 第2判断条件（fromSlave=true且fallbackToMaster=true）
		return SlaveFallbackMaster
	default:
		// 第3条件（fromSlave=true且fallbackToMaster=false）
		return DirectSlave
	}
}

func (s *Slice) ShouldFallbackToMasterOnSlaveFail() bool {
	val := s.FallbackToMasterOnSlaveFail

	// 用户配置了该字段，根据配置值来判断
	switch strings.ToLower(val) {
	case "on":
		return true
	case "off":
		return false
	default:
		// 如果用户填写了其他值，这里默认回退
		return true
	}
}

func (s *Slice) IsFuseEnabled() bool {
	val := s.FuseEnabled

	// 用户配置了该字段，根据配置值来判断
	switch strings.ToLower(val) {
	case "on":
		return true
	case "off":
		return false
	default:
		// 如果用户填写了其他值，这里默认开启
		return true
	}
}

func (s *Slice) GetDirectConn(addr string) (*DirectConnection, error) {
	return NewDirectConnection(addr, s.Cfg.UserName, s.Cfg.Password, "", s.charset, s.collationID, s.Cfg.Capability, s.HandshakeTimeout)
}

// GetMasterStatus return master status
func (s *Slice) GetMasterStatus() (StatusCode, error) {
	// 取 Master 节点
	node, err := s.Master.GetNode(0)
	if err != nil {
		return StatusDown, err
	}
	return node.GetStatus(), nil
}

func (s *Slice) GetSlaveConn(slavesInfo *DBInfo, localSlaveReadPriority int) (PooledConnect, error) {
	// 如果整个 `ConnPool` 为空，或者所有节点都宕机，则直接返回错误
	if len(slavesInfo.Nodes) == 0 || allSlaveIsOffline(slavesInfo) {
		return nil, errors.ErrNoSlaveDB
	}

	switch localSlaveReadPriority {
	case LocalSlaveReadForce:
		// 强制使用本地，从库不可用直接返回错误
		if slavesInfo.LocalBalancer == nil {
			return nil, fmt.Errorf("no local balancer available")
		}
		return s.getConnFromBalancer(slavesInfo, slavesInfo.LocalBalancer)

	case LocalSlaveReadClosed:
		// 不启用本地优先，直接使用 `GlobalBalancer`
		if slavesInfo.GlobalBalancer == nil {
			return nil, fmt.Errorf("no global balancer available")
		}
		return s.getConnFromBalancer(slavesInfo, slavesInfo.GlobalBalancer)

	case LocalSlaveReadPrefer:
		// 优先尝试本地 `LocalBalancer`
		if slavesInfo.LocalBalancer != nil {
			if conn, err := s.getConnFromBalancer(slavesInfo, slavesInfo.LocalBalancer); err == nil {
				return conn, nil
			}
		}
		// 降级尝试 `RemoteBalancer`
		if slavesInfo.RemoteBalancer != nil {
			if conn, err := s.getConnFromBalancer(slavesInfo, slavesInfo.RemoteBalancer); err == nil {
				return conn, nil
			}
		}
		// 如果本地和远程都不可用，直接返回错误
		return nil, fmt.Errorf("no available slave DB in local or remote data centers")
	default:
		// 默认使用 `GlobalBalancer`
		if slavesInfo.GlobalBalancer == nil {
			return nil, fmt.Errorf("no global balancer available,invalid localSlaveReadPriority: %d	", localSlaveReadPriority)
		}
		return s.getConnFromBalancer(slavesInfo, slavesInfo.GlobalBalancer)
	}
}

func (s *Slice) GetMasterConn(masterInfo *DBInfo) (PooledConnect, error) {
	// 如果整个 `ConnPool` 为空，或者所有节点都宕机，则直接返回错误
	if masterIsOffline(masterInfo) {
		return nil, errors.ErrNoMasterDB
	}
	return s.getConnFromMaster(masterInfo)

}

func (s *Slice) getConnFromMaster(master *DBInfo) (pc PooledConnect, err error) {
	node, err := master.GetNode(0)
	if err != nil {
		return nil, fmt.Errorf("master node is not available, err: %v", err)
	}

	if !node.IsStatusUp() {
		return nil, fmt.Errorf("master: %s is Down", node.Address)
	}

	// 获取连接
	ctx := context.TODO()
	pc, err = node.ConnPool.Get(ctx)
	if err != nil {
		return nil, &getConnError{
			Namespace: s.Namespace,
			Addr:      node.Address,
			Err:       err,
		}
	}
	return pc, nil
}

// CheckStatus check slice instance status
func (s *Slice) CheckStatus(ctx context.Context, downAfterNoAlive int, secondsBehindMaster int) {
	go s.checkBackendMasterStatus(ctx, downAfterNoAlive)
	go s.checkBackendSlaveStatus(ctx, s.Slave, downAfterNoAlive, secondsBehindMaster)
	go s.checkBackendSlaveStatus(ctx, s.StatisticSlave, downAfterNoAlive, secondsBehindMaster)
}

func (s *Slice) checkBackendMasterStatus(ctx context.Context, downAfterNoAlive int) {
	defer func() {
		if err := recover(); err != nil {
			log.Warn("[ns:%s, %s] check master status, panic: %v\n%s", s.Namespace, s.Cfg.Name, err, debug.Stack())
		}
	}()

	// 确保 Master 节点存在
	if s.Master == nil || len(s.Master.Nodes) == 0 {
		log.Warn("[ns:%s, %s] check master status, master node is empty", s.Namespace, s.Cfg.Name)
		return
	}

	ticker := time.NewTicker(time.Duration(PingPeriod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Warn("[ns:%s, %s] check master status, context canceled", s.Namespace, s.Cfg.Name)
			return
		case <-ticker.C:
			node := s.Master.Nodes[0]
			log.Debug("[ns:%s, %s:%s] check master status, start", s.Namespace, s.Cfg.Name, node.Address)

			// 1. Master 连接池健康检查
			conn, err := node.GetPooledConnectWithHealthCheck(s.Namespace, s.HealthCheckSql)
			if err != nil {
				log.Warn("[ns:%s, %s:%s] check master status, Get master conn error: %v, last check time: %s", s.Namespace, s.Cfg.Name, node.Address, err, time.Unix(node.ConnPool.GetLastChecked(), 0).Format(mysql.TimeFormat))
				// 连接池可能因为网络原因导致Ping失败，这个时候继续向下检查
			}

			// 2. 判断 Master 是否要下线
			shouldSetDown, elapsed := node.ShouldDownAfterNoAlive(downAfterNoAlive)
			if shouldSetDown {
				node.SetStatusDown()
				log.Warn("[ns:%s, %s:%s] check master status, Marked as StatusDown for %ds", s.Namespace, s.Cfg.Name, node.Address, elapsed)
				continue
			}

			// 3. 更新 Master 状态
			if conn != nil && node.IsStatusDown() {
				node.SetStatusUp()
				log.Warn("[ns:%s, %s:%s] check master status, Master recovered from down, now StatusUp", s.Namespace, s.Cfg.Name, node.Address)
			}
		}
	}
}

func (s *Slice) checkBackendSlaveStatus(ctx context.Context, slave *DBInfo, downAfterNoAlive int, secondBehindMaster int) {
	defer func() {
		if err := recover(); err != nil {
			log.Warn("[ns:%s, %s] check slave status, panic: %v\n%s", s.Namespace, s.Cfg.Name, err, debug.Stack())
		}
	}()
	// 先检查 `slave` 是否 `nil` 或者 `slave.Nodes` 是否为空
	if slave == nil || len(slave.Nodes) == 0 {
		log.Warn("[ns:%s, %s] check slave status, slave DBInfo is nil or no slave nodes", s.Namespace, s.Cfg.Name)
		return
	}

	ticker := time.NewTicker(time.Duration(PingPeriod) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Warn("[ns:%s, %s] check slave status, context canceled", s.Namespace, s.Cfg.Name)
			return
		case <-ticker.C:
			for _, node := range slave.Nodes {
				if err := s.TryRecover(node, downAfterNoAlive, secondBehindMaster); err != nil {
					log.Warn("[ns:%s, %s:%s] check slave status, error: %v", s.Namespace, s.Cfg.Name, node.Address, err)
				}
			}
		}
	}
}

func (s *Slice) TryRecover(node *NodeInfo, downAfterNoAlive int, secondBehindMaster int) error {
	if node.FuseStrategy == nil || node.RecoveryStrategy == nil {
		s.checkWithNoRecovery(node, downAfterNoAlive, secondBehindMaster)
		return nil
	}

	// 根据策略类型动态分支
	switch strategy := node.RecoveryStrategy.(type) {
	case *HardCoolDownStrategy: // 硬恢复策略
		s.checkWithHardRecovery(node, downAfterNoAlive, secondBehindMaster, strategy)
		return nil
	case *GradualRecoveryStrategy: // 渐进式恢复策略
		s.checkWithGradualRecovery(node, downAfterNoAlive, secondBehindMaster, strategy)
		return nil
	default:
		return fmt.Errorf("unsupported recovery strategy: %T", strategy)
	}
}

func (s *Slice) checkWithNoRecovery(node *NodeInfo, downAfterNoAlive int, secondBehindMaster int) {
	start := time.Now()
	log.Debug("[ns:%s, %s:%s] check slave status with no strategy, start", s.Namespace, s.Cfg.Name, node.Address)

	// 1. 从库连接池健康检查
	conn, err := node.GetPooledConnectWithHealthCheck(s.Namespace, s.HealthCheckSql)
	if err != nil {
		log.Warn("[ns:%s, %s:%s] check slave status with no strategy, Get slave conn error: %v, last check time: %s, duration: %v", s.Namespace, s.Cfg.Name, node.Address, err, time.Unix(node.ConnPool.GetLastChecked(), 0).Format(mysql.TimeFormat), time.Since(start))
		// 连接池可能因为网络原因导致Ping失败，这个时候继续向下检查
	}

	// 2. 判断是否需要下线
	shouldSetDown, elapsed := node.ShouldDownAfterNoAlive(downAfterNoAlive)
	if shouldSetDown {
		node.SetStatusDown()
		log.Warn("[ns:%s, %s:%s] check slave status with no strategy, Marked as StatusDown for %ds, duration: %v", s.Namespace, s.Cfg.Name, node.Address, elapsed, time.Since(start))
		return
	}

	// 3. 获取主库状态
	masterStatus, err := s.GetMasterStatus()
	if err != nil || masterStatus == StatusDown {
		log.Warn("[ns:%s, %s:%s] check slave status with no strategy, Get master status: %s, get master err: %v, duration: %v", s.Namespace, s.Cfg.Name, node.Address, masterStatus.String(), err, time.Since(start))
		if node.IsStatusDown() {
			node.SetStatusUp()
			log.Warn("[ns:%s, %s:%s] check slave status with no strategy, Marked as StatusUp success, Slave recovered from down, (case master down), duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
		}
		return
	}

	// 4. 检查从库同步状态
	alive, err := checkSlaveSyncStatus(conn, secondBehindMaster)
	if !alive {
		node.SetStatusDown()
		if err != nil {
			log.Warn("[ns:%s, %s:%s] check slave status with no strategy, Marked as StatusDown due to sync thread failure: %v, duration: %v", s.Namespace, s.Cfg.Name, node.Address, err, time.Since(start))
		} else {
			log.Warn("[ns:%s, %s:%s] check slave status with no strategy, Marked as StatusDown due to excessive replication delay, duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
		}
		return
	}

	// 5. 所有检查通过，最终确认 `StatusUp`,如果 pc == nil,说明连接池可能因为网络原因断开，这个时候保持原状态, 直到通过 ShouldDownAfterNoAlive 超过阈值下线
	if conn != nil && node.IsStatusDown() {
		node.SetStatusUp()
		log.Warn("[ns:%s, %s:%s] check slave status with no strategy, Marked as StatusUp success, Slave recovered from down, (case master up), duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
	}
}

func (s *Slice) checkWithHardRecovery(node *NodeInfo, downAfterNoAlive int, secondBehindMaster int, strategy *HardCoolDownStrategy) {
	log.Debug("[ns:%s, %s:%s] check slave status with hard strategy, start", s.Namespace, s.Cfg.Name, node.Address)
	start := time.Now()
	// 1. 从库连接池健康检查
	conn, err := node.GetPooledConnectWithHealthCheck(s.Namespace, s.HealthCheckSql)
	if err != nil {
		log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, Get slave conn error: %v, last check time: %s, duration: %v", s.Namespace, s.Cfg.Name, node.Address, err, time.Unix(node.ConnPool.GetLastChecked(), 0).Format(mysql.TimeFormat), time.Since(start))
		// 连接池可能因为网络原因导致Ping失败，这个时候继续向下检查
	}

	// 2. 判断是否需要下线
	shouldSetDown, elapsed := node.ShouldDownAfterNoAlive(downAfterNoAlive)
	if shouldSetDown {
		node.SetStatusDown()
		log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, Marked as StatusDown for %ds, duration: %v", s.Namespace, s.Cfg.Name, node.Address, elapsed, time.Since(start))
		return
	}

	// 3. 获取主库状态
	masterStatus, err := s.GetMasterStatus()
	if err != nil || masterStatus == StatusDown {
		log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, Get master status: %s, get master err: %v, duration: %v", s.Namespace, s.Cfg.Name, node.Address, masterStatus.String(), err, time.Since(start))
		if node.IsStatusDown() {
			node.SetStatusUp()
			log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, Marked as StatusUp success, Slave recovered from down, (case master down), duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
		}
		return
	}

	// 4. 检查从库同步状态
	alive, err := checkSlaveSyncStatus(conn, secondBehindMaster)
	if !alive {
		node.SetStatusDown()
		if err != nil {
			log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, Marked as StatusDown due to sync thread failure: %v, duration: %v", s.Namespace, s.Cfg.Name, node.Address, err, time.Since(start))
		} else {
			log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, Marked as StatusDown due to excessive replication delay, duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
		}
		return
	}

	// 5. 所有检查通过，设置状态为 UP
	if conn != nil && node.IsStatusDown() {
		// 5.1. 检查硬恢复策略是否允许恢复
		if !strategy.AllowRecovery() {
			log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, still StatusDown in cooldown period, duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
			return
		} else {
			node.SetStatusUp()
			log.Warn("[ns:%s, %s:%s] check slave status with hard strategy, Marked as StatusUp success, duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
		}
	}
}

func (s *Slice) checkWithGradualRecovery(node *NodeInfo, downAfterNoAlive int, secondBehindMaster int, strategy *GradualRecoveryStrategy) {
	log.Debug("[ns:%s, %s:%s] check slave status with gradual strategy, start", s.Namespace, s.Cfg.Name, node.Address)
	start := time.Now()

	// 1. 从库连接池健康检查
	conn, err := node.GetPooledConnectWithHealthCheck(s.Namespace, s.HealthCheckSql)
	if err != nil {
		log.Warn("[ns:%s, %s:%s] check slave status with gradual strategy, Get slave conn error: %v, duration: %v", s.Namespace, s.Cfg.Name, node.Address, err, time.Since(start))

	}

	// 探活失败，则刷新计数
	if conn == nil && node.IsStatusDown() {
		strategy.RefreshCoolDownCount()
	}

	// 2. 判断是否需要下线
	shouldSetDown, elapsed := node.ShouldDownAfterNoAlive(downAfterNoAlive)
	if shouldSetDown {
		node.SetStatusDown()
		log.Warn("[ns:%s, %s:%s] check slave status with gradual strategy, Marked as StatusDown for %ds, duration: %v", s.Namespace, s.Cfg.Name, node.Address, elapsed, time.Since(start))
		return
	}

	// 3. 获取主库状态
	masterStatus, err := s.GetMasterStatus()
	if err != nil || masterStatus == StatusDown {
		log.Warn("[ns:%s, %s:%s] check slave status with gradual strategy, Skipping slave sync check, master status: %s, duration: %v", s.Namespace, s.Cfg.Name, node.Address, masterStatus.String(), time.Since(start))
		return
	}

	// 4. 检查从库同步状态
	alive, err := checkSlaveSyncStatus(conn, secondBehindMaster)
	if !alive {
		node.SetStatusDown()
		if err != nil {
			log.Warn("[ns:%s, %s:%s] check slave status with gradual strategy, Marked as StatusDown due to sync thread failure: %v, duration: %v", s.Namespace, s.Cfg.Name, node.Address, err, time.Since(start))
		} else {
			log.Warn("[ns:%s, %s:%s] check slave status with gradual strategy, Marked as StatusDown due to excessive replication delay, duration: %v", s.Namespace, s.Cfg.Name, node.Address, time.Since(start))
		}
		return
	}

	// 5. 所有检查通过，处理恢复逻辑
	if conn != nil && node.IsStatusDown() {
		// 5.1. 检查惩罚恢复策略是否允许恢复
		// UP -> DOWN 记录一次误恢复
		if strategy.AllowRecovery() {
			strategy.UpdateLastRecoveryTime()
			node.SetStatusUp()
			log.Warn("[ns:%s, %s:%s] check slave status with gradual strategy, Marked as StatusUp success, from bad recovery: %d, duration: %v", s.Namespace, s.Cfg.Name, node.Address, strategy.errorRecoveryCount.Get(), time.Since(start))
		} else {
			log.Warn("[ns:%s, %s:%s] check slave status with gradual strategy, still StatusDown in cooldown period, remain skip: %d, duration: %v", s.Namespace, s.Cfg.Name, node.Address, strategy.consecutiveSuccessCheckCount.Get(), time.Since(start))
			return
		}
	}
}

func checkInstanceStatus(name string, cp ConnectionPool, healthCheckSql string) (PooledConnect, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("[ns:%s, %s] check instance status panic:%s", name, cp.Addr(), err, debug.Stack())
		}
	}()

	pc, err := cp.GetCheck(context.Background())
	if err != nil {
		if pc != nil {
			pc.Close()
		}
		return nil, fmt.Errorf("get check conn err:%s", err)
	}

	if pc == nil {
		return nil, fmt.Errorf("get nil check conn, ins:%s", cp.Addr())
	}

	// 重复多次检测通过，才算通过，网络不稳定时，可以使得检测更准确
	for i := 0; i <= CheckRepeat; i++ {
		if len(healthCheckSql) > 0 {
			_, err := pc.ExecuteWithTimeout(healthCheckSql, 0, ExecTimeOut)
			if err == nil {
				cp.SetLastChecked()
				return pc, nil
			}
			log.Warn("[ns:%s instance:%s] exec health check sql:%s sqlError:%v", name, cp.Addr(), healthCheckSql, err)
			if mysql.IsServerShutdownErr(err) || mysql.IsTableSpaceMissingErr(err) || mysql.IsTableSpaceDiscardeErr(err) || err == ErrExecuteTimeout {
				pc.Close()
				return nil, fmt.Errorf("exec health check query error:%s", err)
			}
		}
		if err = pc.PingWithTimeout(GetConnTimeout); err != nil {
			pc.Close()
			return nil, fmt.Errorf("ping conn error:%s", err)
		}
		// 增加 select 1 探活
		if _, err = pc.ExecuteWithTimeout("select 1", 0, SelectSimpleTimeOut); err != nil {
			pc.Close()
			return nil, fmt.Errorf("ping conn error:%s", err)
		}
	}

	cp.SetLastChecked()
	return pc, nil
}

func allSlaveIsOffline(slavesInfo *DBInfo) bool {
	for _, node := range slavesInfo.Nodes {
		if node.IsStatusUp() {
			return false // 只要有一个 Slave 处于 UP 状态，则返回 false
		}
	}
	return true // 所有 Slave 都 Down，则返回 true
}

func masterIsOffline(master *DBInfo) bool {
	if master == nil {
		// Master 是 Down，则返回 true
		return true
	}

	if len(master.Nodes) == 0 {
		return true
	}

	if master.Nodes[0].IsStatusUp() {
		return false // Master 处于 UP 状态，则返回 false
	}
	return true // Master 是 Down，则返回 true
}

func (dbInfo *DBInfo) InitBalancers(proxyDatacenter string) error {
	if len(dbInfo.Nodes) == 0 {
		return nil
	}
	// 1️ 获取本地、远程和全局的索引+权重
	local, remote, global := dbInfo.getIndicesAndWeights(proxyDatacenter)

	// 2️ 初始化 `GlobalBalancer`（所有节点）
	var globalBalancerErr, localBalancerErr, remoteBalancerErr error
	if len(global.Indices) > 0 {
		dbInfo.GlobalBalancer, globalBalancerErr = newBalancer(global.Indices, global.Weights)
	}

	// 3️ 初始化 `LocalBalancer`（仅本地数据中心）
	if len(local.Indices) > 0 {
		dbInfo.LocalBalancer, localBalancerErr = newBalancer(local.Indices, local.Weights)
	}

	// 4️ 初始化 `RemoteBalancer`（仅远程数据中心）
	if len(remote.Indices) > 0 {
		dbInfo.RemoteBalancer, remoteBalancerErr = newBalancer(remote.Indices, remote.Weights)
	}

	// 5️ 检查错误
	if globalBalancerErr != nil || localBalancerErr != nil || remoteBalancerErr != nil {
		return fmt.Errorf("failed to initialize balancers: global=%v, local=%v, remote=%v",
			globalBalancerErr, localBalancerErr, remoteBalancerErr)
	}
	return nil
}

// getConnFromBalancer 封装从给定 balancer 中依次尝试选取健康连接的逻辑
func (s *Slice) getConnFromBalancer(slavesInfo *DBInfo, bal *balancer) (PooledConnect, error) {
	node, err := s.getNodeFromBalancer(slavesInfo, bal)
	if err != nil || node == nil {
		return nil, err
	}
	return s.getConnWithFuse(node)
}

func (s *Slice) getNodeFromBalancer(slavesInfo *DBInfo, bal *balancer) (*NodeInfo, error) {
	// 加锁保证同一时刻只有一个 Session 在使用该 balancer
	slavesInfo.Lock()
	defer slavesInfo.Unlock()
	for i := 0; i < len(bal.roundRobinQ); i++ {
		index, err := bal.next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next index from balancer: %v", err)
		}
		if node, _ := slavesInfo.GetNode(index); node != nil && node.IsStatusUp() {
			return node, nil
		}
	}
	return nil, fmt.Errorf("no healthy connection available from selected balancer")
}

// getConnWithFuse 封装从连接池获取连接并处理熔断的逻辑
func (s *Slice) getConnWithFuse(node *NodeInfo) (PooledConnect, error) {
	pc, err := node.ConnPool.Get(context.TODO())
	// 熔断埋点
	s.TryFuse(node, err)
	if err != nil {
		return nil, &getConnError{
			Namespace: s.Namespace,
			Addr:      node.Address,
			Err:       err,
		}
	}
	return pc, nil
}

// tryFuse 处理熔断触发及策略相关副作用
func (s *Slice) TryFuse(node *NodeInfo, err error) {
	if node.FuseStrategy == nil || node.RecoveryStrategy == nil {
		return
	}

	if !mysql.AsConnError(err) {
		return
	}

	now := time.Now()
	if !node.FuseStrategy.Trigger(now.Unix()) {
		return
	}
	// 熔断后下线节点
	changed := node.SetStatusDown()
	log.Warn("[ns=%s][addr=%s] Triggered fuse, node marked as DOWN at %s", s.Namespace, node.Address, now.Format(mysql.TimeFormat))

	// 根据恢复策略类型执行不同处理
	switch strategy := node.RecoveryStrategy.(type) {
	case *GradualRecoveryStrategy:
		if !changed {
			return
		}
		// 渐进式恢复策略需要更新熔断时间、误恢复次数
		strategy.UpdateFuseTime(now.Unix())
		isBad := strategy.IsBadRecovery(now.Unix())
		if isBad {
			strategy.UpdateCoolDownCount()
		} else {
			strategy.ResetBadRecovery(now.Unix())
		}
	case *HardCoolDownStrategy:
		// 硬熔断策略需要更新熔断时间
		strategy.UpdateFuseTime(now.Unix())
	case nil:
		// 无恢复策略的默认处理
	default:
		log.Warn("Unhandled recovery strategy type: %T", strategy)
	}
}

// Close close the pool in slice
// Concurrently close the connection to prevent the main database transaction from not being committed
// and causing the slave database connection to not be released.
func (s *Slice) Close() error {
	s.Lock()
	defer s.Unlock()

	var wg sync.WaitGroup
	closePool := func(info *DBInfo) {
		defer wg.Done()
		if info != nil {
			for _, node := range info.Nodes {
				node.ConnPool.Close()
			}
		}
	}
	// close master
	wg.Add(1)
	go closePool(s.Master)

	// close slaves
	wg.Add(1)
	go closePool(s.Slave)

	// close statistic slaves
	wg.Add(1)
	go closePool(s.StatisticSlave)

	// close monitor master
	wg.Add(1)
	go closePool(s.MonitorMaster)

	// close monitor slaves
	wg.Add(1)
	go closePool(s.MonitorSlave)

	wg.Wait()
	return nil
}

// ParseMaster 解析主库信息，并初始化负载均衡器
func (s *Slice) ParseMaster(masterAddr string) error {
	s.Master = &DBInfo{Nodes: []*NodeInfo{}} // 确保 Master 不为 nil
	if len(masterAddr) == 0 {
		return nil
	}

	dbInfo, err := s.parseDBInfo([]string{masterAddr}, false)
	if err != nil {
		return err
	}

	// 初始化 `Balancer`
	if err := dbInfo.InitBalancers(s.ProxyDatacenter); err != nil {
		return fmt.Errorf("failed to initialize master balancer: %w", err)
	}

	s.Master = dbInfo
	return nil
}

func (s *Slice) ParseMonitorMaster(masterAddr string) error {
	s.MonitorMaster = &DBInfo{Nodes: []*NodeInfo{}} // 确保 Master 不为 nil
	if len(masterAddr) == 0 {
		return nil
	}

	dbInfo, err := s.parseDBInfo([]string{masterAddr}, true)
	if err != nil {
		return fmt.Errorf("failed to parse monitor master db info: %w", err)
	}

	// 初始化 `Balancer`
	if err := dbInfo.InitBalancers(s.ProxyDatacenter); err != nil {
		return fmt.Errorf("failed to initialize master balancer: %w", err)
	}

	s.MonitorMaster = dbInfo
	return nil
}

// ParseSlave 解析从库信息，并初始化负载均衡器
func (s *Slice) ParseSlave(slaveAddrs []string) error {
	// 确保 Slave 不为 nil
	s.Slave = &DBInfo{Nodes: []*NodeInfo{}}
	if len(slaveAddrs) == 0 {
		return nil
	}

	dbInfo, err := s.parseDBInfo(slaveAddrs, false)
	if err != nil {
		return err
	}

	// 初始化 `Balancer`
	if err := dbInfo.InitBalancers(s.ProxyDatacenter); err != nil {
		return fmt.Errorf("failed to initialize slave balancer: %w", err)
	}

	s.Slave = dbInfo
	return nil
}

func (s *Slice) ParseMonitorSlave(slaveAddrs []string) error {
	// 确保 Slave 不为 nil
	s.MonitorSlave = &DBInfo{Nodes: []*NodeInfo{}}
	if len(slaveAddrs) == 0 {
		return nil
	}

	dbInfo, err := s.parseDBInfo(slaveAddrs, true)

	if err != nil {
		return fmt.Errorf("failed to parse monitor slave db info: %w", err)
	}

	// 初始化 `Balancer`
	if err := dbInfo.InitBalancers(s.ProxyDatacenter); err != nil {
		return fmt.Errorf("failed to initialize slave balancer: %w", err)
	}

	s.MonitorSlave = dbInfo
	return nil
}

// ParseStatisticSlave 解析统计从库信息（权重一般为 1）
func (s *Slice) ParseStatisticSlave(statisticSlaveAddrs []string) error {
	// 确保 StatisticSlave 不为 nil
	s.StatisticSlave = &DBInfo{Nodes: []*NodeInfo{}}
	if len(statisticSlaveAddrs) == 0 {
		return nil
	}

	dbInfo, err := s.parseDBInfo(statisticSlaveAddrs, false)
	if err != nil {
		return err
	}

	// 初始化 `Balancer`
	if err := dbInfo.InitBalancers(s.ProxyDatacenter); err != nil {
		return fmt.Errorf("failed to initialize statistic slave balancer: %w", err)
	}

	s.StatisticSlave = dbInfo
	return nil
}

// parseDBInfo 解析数据库节点信息，并返回 `DBInfo`,解析权重时：
// 如果解析失败（非法字符，如 @&, @+, @-），默认权重 1。
// 如果解析成功但 权重 < 0，默认权重 1。
// 如果解析成功但 权重 = 0，默认 权重 0。
// 其他情况下，使用解析出的权重。注意，`0` 权重的节点并没有被过滤。
func (s *Slice) parseDBInfo(dbAddrs []string, isMonitorUser bool) (*DBInfo, error) {
	if len(dbAddrs) == 0 {
		return &DBInfo{Nodes: []*NodeInfo{}}, nil // 确保 `Nodes` 为空 slice，而不是 `nil`
	}

	var err error
	nodes := make([]*NodeInfo, 0, len(dbAddrs))

	for _, addr := range dbAddrs {
		if addr == "" {
			continue
		}

		// 解析数据中心和权重
		dc := ""
		addrAndWeightDatacenter := strings.Split(addr, datacenterSplit) // 格式: "host:port@权重#数据中心"
		if len(addrAndWeightDatacenter) == 2 {
			dc = addrAndWeightDatacenter[1] // 提取数据中心
			addr = addrAndWeightDatacenter[0]
		}

		// 解析权重（默认值为 `1`）
		weight := 1
		addrAndWeight := strings.Split(addr, weightSplit) // 格式: "host:port@权重"
		if len(addrAndWeight) == 2 {
			parsedWeight, err := strconv.Atoi(addrAndWeight[1])
			if err != nil {
				log.Warn("parseDBInfo Error, Invalid weight format for %s: %v. defaulting to weight=1", addr, err)
				weight = 1
			} else if parsedWeight < 0 {
				log.Warn("parseDBInfo Error, Negative or zero weight for %s. defaulting to weight=1", addr)
				weight = 1
			} else {
				weight = parsedWeight
			}
		}

		// 获取数据中心（如果 `#` 解析不到，则调用 `util.GetInstanceDatacenter`）
		if dc == "" {
			dc, err = util.GetInstanceDatacenter(addrAndWeight[0])
			if err != nil {
				log.Warn("get datacenter failed for %s: %s, using default gaea proxy datacenter: %s", addrAndWeight[0], err, s.ProxyDatacenter)
				dc = s.ProxyDatacenter
			}
		}

		// 创建连接池
		idleTimeout, err := util.Int2TimeDuration(s.Cfg.IdleTimeout)
		if err != nil {
			return nil, err
		}
		var cp ConnectionPool
		if isMonitorUser {
			cp = NewConnectionPool(
				addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", MonitorDefaultCapacity,
				MonitorDefaultMaxCapacity, idleTimeout, s.charset, s.collationID,
				s.Cfg.Capability, s.Cfg.InitConnect, dc, s.HandshakeTimeout,
			)
		} else {
			cp = NewConnectionPool(
				addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity,
				s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID,
				s.Cfg.Capability, s.Cfg.InitConnect, dc, s.HandshakeTimeout,
			)
		}

		if err = cp.Open(); err != nil {
			return nil, err
		}

		// 创建 `NodeInfo`
		nodes = append(nodes, &NodeInfo{
			Address:    addrAndWeight[0],
			Datacenter: dc,
			Weight:     weight,
			ConnPool:   cp,
			Status:     StatusUp, // 默认 `Up`
		})
	}

	return &DBInfo{Nodes: nodes}, nil
}

// SetCharsetInfo set charset
func (s *Slice) SetCharsetInfo(charset string, collationID mysql.CollationID) {
	s.charset = charset
	s.collationID = collationID
}

type SlaveStatus struct {
	SecondsBehindMaster uint64
	SlaveIORunning      string
	SlaveSQLRunning     string
	MasterLogFile       string
	ReadMasterLogPos    uint64
	RelayMasterLogFile  string
	ExecMasterLogPos    uint64
}

// checkSlaveSyncStatus check slave sync status, if slave is not sync, return false
func checkSlaveSyncStatus(pc PooledConnect, secondsBehindMaster int) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Warn("check slave sync status panic: %v", err)
		}
	}()

	// if secondsBehindMaster is 0, we won't check slave sync status
	if secondsBehindMaster == 0 {
		return true, nil
	}

	if pc == nil {
		return true, fmt.Errorf("check slave sync status error: pc is nil")
	}

	skipCheck, slaveStatus, err := GetSlaveStatus(pc)
	if err != nil {
		return false, fmt.Errorf("check slave sync status error: %v", err)
	}
	// if suspectedMaster is true, we think this is a master
	if skipCheck {
		return true, nil
	}

	// 如果 SecondsBehindMaster 超出阈值，标记 `StatusDown`
	if slaveStatus.SecondsBehindMaster > uint64(secondsBehindMaster) {
		log.Warn("Slave has sync delay: SecondsBehindMaster=%d, Threshold=%d - marking as StatusDown", slaveStatus.SecondsBehindMaster, secondsBehindMaster)
		return false, nil
	}

	// 处理同步线程状态异常（I/O 或 SQL 线程未运行），直接标记为 `StatusDown`
	if slaveStatus.SlaveIORunning != "Yes" {
		return false, fmt.Errorf("slave sync threads not running, IO=%s", slaveStatus.SlaveIORunning)
	}

	if slaveStatus.SlaveSQLRunning != "Yes" {
		return false, fmt.Errorf("slave sync threads not running, SQL=%s", slaveStatus.SlaveSQLRunning)
	}

	return true, nil
}

// GetSlaveStatus get slave status, will check bellow cases:
// 1. if we have no privileges to get slave status, will return skipCheck true.
// 2. if slave status result is nil,maybe it's master but configured as slave, will return skipCheck true.
// 3. return slave status result with skipCheck false.
func GetSlaveStatus(conn PooledConnect) (bool, SlaveStatus, error) {
	var slaveStatus SlaveStatus
	res, err := conn.Execute("show slave status;", 0)

	// if exec error is syntax error or no privilege, will return skipCheck true.
	if err != nil {
		if mysql.IsSQLNoPrivilegeErr(err) {
			log.Warn("addr:%s, get slave status error,maybe configured error.err:%s.", conn.GetAddr(), err)
			return true, slaveStatus, nil
		}
		return false, slaveStatus, fmt.Errorf("execute show slave status error:%s", err)
	}

	// if we have no privileges to get slave status, will return skipCheck true.
	if res.RowNumber() == 0 {
		log.Debug("addr:%s, slave status is empty,maybe is master\n", conn.GetAddr())
		return true, slaveStatus, nil
	}

	for _, f := range res.Fields {
		fieldName := string(f.Name)
		var col interface{}
		col, err = res.GetValueByName(0, fieldName)
		if err != nil {
			_ = log.Warn("get field name Get '%s' failed in SlaveStatus, err: %v", fieldName, err)
			break
		}

		switch strings.ToLower(fieldName) {
		case "seconds_behind_master":
			switch col.(type) {
			case uint64:
				slaveStatus.SecondsBehindMaster = col.(uint64)
			default:
				slaveStatus.SecondsBehindMaster = 0
			}
		case "slave_io_running":
			switch col.(type) {
			case string:
				slaveStatus.SlaveIORunning = col.(string)
			default:
				slaveStatus.SlaveIORunning = "No"
			}
		case "slave_sql_running":
			switch col.(type) {
			case string:
				slaveStatus.SlaveSQLRunning = col.(string)
			default:
				slaveStatus.SlaveSQLRunning = "No"
			}
		case "master_log_file":
			switch col.(type) {
			case string:
				slaveStatus.MasterLogFile = col.(string)
			default:
				slaveStatus.MasterLogFile = ""
			}
		case "read_master_log_pos":
			switch col.(type) {
			case uint64:
				slaveStatus.ReadMasterLogPos = col.(uint64)
			default:
				slaveStatus.ReadMasterLogPos = 0
			}
		case "relay_master_log_file":
			switch col.(type) {
			case string:
				slaveStatus.RelayMasterLogFile = col.(string)
			default:
				slaveStatus.RelayMasterLogFile = ""
			}
		case "exec_master_log_pos":
			switch col.(type) {
			case uint64:
				slaveStatus.ExecMasterLogPos = col.(uint64)
			default:
				slaveStatus.ExecMasterLogPos = 0
			}
		default:
			continue
		}
	}
	return false, slaveStatus, err
}
