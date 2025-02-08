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
)

func (s *StatusCode) String() string {
	r := "StatusUp"
	if *s == StatusDown {
		r = "StatusDown"
	}

	return fmt.Sprintf(r)
}

type DBInfo struct {
	ConnPool       []ConnectionPool
	Nodes          []NodeInfo // 存储所有节点的详细信息，索引对应 ConnPool
	LocalBalancer  *balancer  // 仅包含数据中心与代理匹配的节点
	RemoteBalancer *balancer  // 包含除本机外的所有其他数据节点
	GlobalBalancer *balancer  // 包含所有节点

	StatusMap *sync.Map
	// 用于记录连续错误次数，键为索引，值为错误计数
	ConsecutiveErrors *sync.Map
}

type NodeInfo struct {
	Address    string // 节点地址，如 "c3-mysql-test00.bj:3306"
	Datacenter string // 节点所属的数据中心，例如 "c3"、"c4"
	Weight     int    // 该节点的负载均衡权重
}

func (dbi *DBInfo) GetStatus(index int) (StatusCode, error) {
	if index > len(dbi.ConnPool) {
		return StatusDown, fmt.Errorf("index:%d out of range", index)
	}
	if value, ok := dbi.StatusMap.Load(index); ok {
		return value.(StatusCode), nil
	}
	return StatusDown, fmt.Errorf("can't get status of index:%d", index)
}

func (dbi *DBInfo) SetStatus(index int, status StatusCode) {
	dbi.StatusMap.Store(index, status)
}

func (dbi *DBInfo) IncrementErrorCount(index int) {
	count, _ := dbi.ConsecutiveErrors.LoadOrStore(index, 0)
	dbi.ConsecutiveErrors.Store(index, count.(int)+1)
}

func (dbi *DBInfo) GetErrorCount(index int) int {
	count, _ := dbi.ConsecutiveErrors.LoadOrStore(index, 0)
	return count.(int)
}

func (dbi *DBInfo) ResetErrorCount(index int) {
	dbi.ConsecutiveErrors.Store(index, 0)
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

	Master                      *DBInfo
	Slave                       *DBInfo
	StatisticSlave              *DBInfo
	ProxyDatacenter             string
	charset                     string
	collationID                 mysql.CollationID
	HealthCheckSql              string
	MaxSlaveFuseErrorCount      int
	HandshakeTimeout            time.Duration
	FallbackToMasterOnSlaveFail string // 控制从库获取失败时是否回退到主库
}

// GetSliceName return name of slice
func (s *Slice) GetSliceName() string {
	return s.Cfg.Name
}

// GetConn get backend connection from different node based on fromSlave and userType
func (s *Slice) GetConn(fromSlave bool, userType int, localSlaveReadPriority int) (pc PooledConnect, err error) {
	if fromSlave {
		if userType == models.StatisticUser {
			pc, err = s.GetSlaveConn(s.StatisticSlave, localSlaveReadPriority)
			if err != nil {
				return nil, err
			}
		} else {
			pc, err = s.GetSlaveConn(s.Slave, localSlaveReadPriority)
			if err != nil {
				// 如果从库连接失败，根据配置决定是否回退
				if s.ShouldFallbackToMasterOnSlaveFail() {
					log.Warn("get connection from slave failed, try to get from master, error: %s", err.Error())
					pc, err = s.GetMasterConn()
				} else {
					return nil, err
				}
			}
		}
	} else {
		pc, err = s.GetMasterConn()
	}
	if err != nil {
		log.Warn("get connection from backend failed, error: %s", err.Error())
	}
	return
}

func (s *Slice) ShouldFallbackToMasterOnSlaveFail() bool {
	val := s.Cfg.FallbackToMasterOnSlaveFail

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

func (s *Slice) GetDirectConn(addr string) (*DirectConnection, error) {
	return NewDirectConnection(addr, s.Cfg.UserName, s.Cfg.Password, "", s.charset, s.collationID, s.Cfg.Capability, s.HandshakeTimeout)
}

// GetMasterConn return a connection in master pool
func (s *Slice) GetMasterConn() (PooledConnect, error) {
	if v, _ := s.Master.StatusMap.Load(0); v != StatusUp {
		return nil, fmt.Errorf("master:%s is Down", s.Cfg.Master)
	}

	ctx := context.TODO()
	return s.Master.ConnPool[0].Get(ctx)
}

// GetMasterStatus return master status
func (s *Slice) GetMasterStatus() (StatusCode, error) {
	return s.Master.GetStatus(0)
}

// SetMasterStatus set master status
func (s *Slice) SetMasterStatus(code StatusCode) {
	s.Master.SetStatus(0, code)
}

// CheckStatus check slice instance status
func (s *Slice) CheckStatus(ctx context.Context, name string, downAfterNoAlive int, secondsBehindMaster int) {
	go s.checkBackendMasterStatus(ctx, name, downAfterNoAlive)
	go s.checkBackendSlaveStatus(ctx, s.Slave, name, downAfterNoAlive, secondsBehindMaster)
	go s.checkBackendSlaveStatus(ctx, s.StatisticSlave, name, downAfterNoAlive, secondsBehindMaster)
}

func (s *Slice) checkBackendMasterStatus(ctx context.Context, name string, downAfterNoAlive int) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("[ns:%s, %s] check master status panic:%s", name, s.Cfg.Name, err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Warn("[ns:%s, %s] check master status canceled", name, s.Cfg.Name)
			return
		case <-time.After(time.Duration(PingPeriod) * time.Second):
			if len(s.Master.ConnPool) == 0 {
				log.Warn("[ns:%s, %s] master is empty", name, s.Cfg.Name)
				continue
			}
			cp := s.Master.ConnPool[0]
			log.Debug("[ns:%s, %s:%s] start check master", name, s.Cfg.Name, cp.Addr())
			_, err := checkInstanceStatus(name, cp, s.HealthCheckSql)

			if time.Now().Unix()-cp.GetLastChecked() >= int64(downAfterNoAlive) {
				s.SetMasterStatus(StatusDown)
				log.Warn("[ns:%s, %s:%s] check master StatusDown for %ds. err: %s", name, s.Cfg.Name, cp.Addr(), time.Now().Unix()-cp.GetLastChecked(), err)
				continue
			}
			if err != nil {
				log.Warn("[ns:%s, %s:%s] check master error:%s", name, s.Cfg.Name, cp.Addr(), err)
			}
			oldStatus, err := s.GetMasterStatus()
			if err != nil {
				log.Warn("[ns:%s, %s:%s] get master master status error:%s", name, s.Cfg.Name, cp.Addr(), err)
				continue
			}

			s.SetMasterStatus(StatusUp)
			if oldStatus == StatusDown {
				log.Warn("[ns:%s, %s:%s] check master StatusUp", name, s.Cfg.Name, cp.Addr())
			}
		}
	}
}

func (s *Slice) checkBackendSlaveStatus(ctx context.Context, db *DBInfo, name string, downAfterNoAlive int, secondBehindMaster int) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("[ns:%s, %s] check slave status panic:%s", name, s.Cfg.Name, err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Warn("[ns:%s, %s] check slave status canceled", name, s.Cfg.Name)
			return
		case <-time.After(time.Duration(PingPeriod) * time.Second):
			for idx, cp := range db.ConnPool {
				log.Debug("[ns:%s, %s:%s] start check slave", name, s.Cfg.Name, cp.Addr())

				oldStatus, err := db.GetStatus(idx)
				if err != nil {
					log.Warn("[ns:%s, %s:%s] get slave status error:%s", name, s.Cfg.Name, cp.Addr(), err)
					continue
				}
				pc, err := checkInstanceStatus(name, cp, s.HealthCheckSql)
				// check slave status
				if time.Now().Unix()-cp.GetLastChecked() >= int64(downAfterNoAlive) {
					db.SetStatus(idx, StatusDown)
					log.Warn("[ns:%s, %s:%s] check slave StatusDown for %ds. err:%s", name, s.Cfg.Name, cp.Addr(), time.Now().Unix()-cp.GetLastChecked(), err)
					continue
				}

				// check master status, if master is down, we should not check slave sync status,cause slave io thread is close
				if masterStatus, err := s.GetMasterStatus(); err != nil {
					log.Warn("[ns:%s, %s:%s] get master status error:%s", name, s.Cfg.Name, cp.Addr(), err)
					continue
				} else if masterStatus == StatusDown {
					// set slave status to up to avoid slave down when master is down on startup
					db.SetStatus(idx, StatusUp)
					if oldStatus == StatusDown {
						log.Warn("[ns:%s, %s:%s] check slave StatusUp", name, s.Cfg.Name, cp.Addr())
					}
					continue
				}

				if pc == nil {
					errInfo := "get nil conn"
					if err != nil {
						errInfo += ", " + err.Error()
					}
					log.Warn("[ns:%s, %s:%s] skip check slave sync, %s", name, s.Cfg.Name, cp.Addr(), errInfo)
					continue
				}

				if alive, err := checkSlaveSyncStatus(pc, secondBehindMaster); !alive {
					db.SetStatus(idx, StatusDown)
					log.Warn("[ns:%s, %s:%s] check slave sync error:%s", name, s.Cfg.Name, cp.Addr(), err)
					continue
				}

				db.SetStatus(idx, StatusUp)
				db.ResetErrorCount(idx) // 探活成功，错误计数置0
				if oldStatus == StatusDown {
					log.Warn("[ns:%s, %s:%s] check slave StatusUp", name, s.Cfg.Name, cp.Addr())
				}
			}
		}
	}
}

func checkInstanceStatus(name string, cp ConnectionPool, healthCheckSql string) (PooledConnect, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("[ns:%s, %s] check instance status panic:%s", name, cp.Addr(), err)
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

	cp.SetLastChecked()
	return pc, nil
}

func allSlaveIsOffline(SlaveStatusMap *sync.Map) bool {
	var result = true
	SlaveStatusMap.Range(func(k, v interface{}) bool {
		if v == StatusUp {
			result = false
			return false
		}
		return true
	})

	return result
}

func (dbInfo *DBInfo) InitBalancers(proxyDatacenter string) error {
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

func (s *Slice) GetSlaveConn(slavesInfo *DBInfo, localSlaveReadPriority int) (PooledConnect, error) {
	// 如果整个 `ConnPool` 为空，或者所有节点都宕机，则直接返回错误
	if len(slavesInfo.ConnPool) == 0 || allSlaveIsOffline(slavesInfo.StatusMap) {
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

// getConnFromBalancer 封装从给定 balancer 中依次尝试选取健康连接的逻辑
func (s *Slice) getConnFromBalancer(slavesInfo *DBInfo, bal *balancer) (PooledConnect, error) {
	// 加锁保证同一时刻只有一个 Session 在使用该 balancer
	s.Lock()
	defer s.Unlock()
	for i := 0; i < len(bal.roundRobinQ); i++ {
		index, err := bal.next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next index from balancer: %v", err)
		}
		// 检查该节点是否健康
		if status, err := slavesInfo.GetStatus(index); err != nil || status == StatusDown {
			continue
		}
		// 找到一个健康节点，返回连接（并处理熔断逻辑）
		return s.getConnWithFuse(slavesInfo, index)
	}
	return nil, fmt.Errorf("no healthy connection available from selected balancer")
}

// getConnWithFuse 封装从连接池获取连接并处理熔断的逻辑
func (s *Slice) getConnWithFuse(slavesInfo *DBInfo, idx int) (PooledConnect, error) {
	pc, err := slavesInfo.ConnPool[idx].Get(context.TODO())
	if err != nil {
		// 如果没有开启熔断 或者类型不是连接超时类型的错误直接返回错误
		if s.MaxSlaveFuseErrorCount == 0 || !mysql.IsConnectionTimeoutError(err) {
			log.Warn("Failed to get connection from pool (index: %d): %v", idx, err)
			return pc, err
		}
		// 增加错误计数
		slavesInfo.IncrementErrorCount(idx)
		// 检查是否需要熔断
		if slavesInfo.GetErrorCount(idx) >= s.MaxSlaveFuseErrorCount {
			slavesInfo.SetStatus(idx, StatusDown)
			log.Warn("Fuse triggered for pool (index: %d). Set status to StatusDown.", idx)
		}
		return nil, err
	}

	// 如果获取连接成功，且有熔断机制，则重置错误计数
	if s.MaxSlaveFuseErrorCount > 0 {
		slavesInfo.ResetErrorCount(idx)
	}
	return pc, nil
}

// Close close the pool in slice
// Concurrently close the connection to prevent the main database transaction from not being committed
// and causing the slave database connection to not be released.
func (s *Slice) Close() error {
	s.Lock()
	defer s.Unlock()

	var wg sync.WaitGroup
	closePool := func(connPools []ConnectionPool) {
		defer wg.Done()
		for i := range connPools {
			connPools[i].Close()
		}
	}
	// close master
	wg.Add(1)
	go closePool(s.Master.ConnPool)

	// close slaves
	wg.Add(1)
	go closePool(s.Slave.ConnPool)

	// close statistic slaves
	wg.Add(1)
	go closePool(s.StatisticSlave.ConnPool)

	wg.Wait()
	return nil
}

// ParseMaster create master connection pool
func (s *Slice) ParseMaster(masterStr string) error {
	if len(masterStr) == 0 {
		return errors.ErrNoMasterDB
	}
	idleTimeout, err := util.Int2TimeDuration(s.Cfg.IdleTimeout)
	if err != nil {
		return err
	}
	dc, err := util.GetInstanceDatacenter(masterStr)
	if err != nil {
		log.Warn("get master(%s) datacenter err:%s,will use default proxy datacenter.", masterStr, err)
		dc = s.ProxyDatacenter
	}
	connectionPool := NewConnectionPool(masterStr, s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID, s.Cfg.Capability, s.Cfg.InitConnect, dc, s.HandshakeTimeout)
	if err := connectionPool.Open(); err != nil {
		return err
	}

	status := &sync.Map{}
	status.Store(0, StatusUp)

	slaveConsecutiveErrors := &sync.Map{}
	slaveConsecutiveErrors.Store(0, 0)

	s.Master = &DBInfo{
		ConnPool:          []ConnectionPool{connectionPool},
		StatusMap:         status,
		ConsecutiveErrors: slaveConsecutiveErrors,
	}
	return nil
}

// ParseSlave create connection pool of slaves
// (127.0.0.1:3306@2,192.168.0.12:3306@3)
func (s *Slice) ParseSlave(slaves []string) (*DBInfo, error) {
	if len(slaves) == 0 {
		return &DBInfo{}, nil
	}

	var err error
	var weight int

	count := len(slaves)
	connPool := make([]ConnectionPool, 0, count)
	nodes := make([]NodeInfo, 0, count)

	//parse addr and weight
	for i := 0; i < count; i++ {
		if slaves[i] == "" {
			continue
		}
		// slave[i] 格式: c3-mysql-test00.bj:3306@10#bj
		dc := ""
		addrAndWeightDatacenter := strings.Split(slaves[i], datacenterSplit)
		if len(addrAndWeightDatacenter) == 2 {
			dc = addrAndWeightDatacenter[1]
			slaves[i] = addrAndWeightDatacenter[0]
		}

		addrAndWeight := strings.Split(slaves[i], weightSplit)
		if len(addrAndWeight) == 2 {
			weight, err = strconv.Atoi(addrAndWeight[1])
			if err != nil {
				return nil, err
			}
		} else {
			weight = 1
		}

		idleTimeout, err := util.Int2TimeDuration(s.Cfg.IdleTimeout)
		if err != nil {
			return nil, err
		}
		// if dc not config, get hostname prefix and suffix
		if dc == "" {
			dc, err = util.GetInstanceDatacenter(addrAndWeight[0])
			if err != nil {
				log.Warn("get master(%s) datacenter err:%s,will use default proxy datacenter.", addrAndWeight[0], err)
				dc = s.ProxyDatacenter
			}
		}

		cp := NewConnectionPool(addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID, s.Cfg.Capability, s.Cfg.InitConnect, dc, s.HandshakeTimeout)
		if err = cp.Open(); err != nil {
			return nil, err
		}
		connPool = append(connPool, cp)
		nodes = append(nodes, NodeInfo{
			Address:    addrAndWeight[0],
			Datacenter: dc,
			Weight:     weight,
		})
	}

	if len(nodes) == 0 {
		return &DBInfo{}, nil
	}

	StatusMap := &sync.Map{}
	slaveConsecutiveErrors := &sync.Map{}
	for idx := range connPool {
		StatusMap.Store(idx, StatusUp)
		slaveConsecutiveErrors.Store(idx, 0)
	}

	return &DBInfo{
		ConnPool:          connPool,
		Nodes:             nodes,
		StatusMap:         StatusMap,
		ConsecutiveErrors: slaveConsecutiveErrors,
	}, nil
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
			log.Fatal("check slave sync status panic:%s", err)
		}
	}()

	// if secondsBehindMaster is 0, we won't check slave sync status
	if secondsBehindMaster == 0 {
		return true, nil
	}

	skipCheck, slaveStatus, err := GetSlaveStatus(pc)
	if err != nil {
		return false, fmt.Errorf("get slave status error:%s", err)
	}
	// if suspectedMaster is true, we think this is a master
	if skipCheck {
		return true, nil
	}

	if slaveStatus.SecondsBehindMaster > uint64(secondsBehindMaster) {
		return false, fmt.Errorf("SecondsBehindMaster(%d) larger than %d", slaveStatus.SecondsBehindMaster, secondsBehindMaster)
	}

	if slaveStatus.SlaveIORunning != "Yes" {
		return false, fmt.Errorf("io thread not running")
	}
	if slaveStatus.SlaveSQLRunning != "Yes" {
		return false, fmt.Errorf("sql thread not running")
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
