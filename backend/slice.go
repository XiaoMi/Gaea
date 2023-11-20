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
	"github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

type StatusCode uint32

const (
	weightSplit     = "@"
	datacenterSplit = "#"

	// DefaultSlice means default slice for namespace
	DefaultSlice       = "slice-0"
	PingPeriod   int64 = 4

	StatusUp                StatusCode = 1
	StatusDown              StatusCode = 0
	LocalSlaveReadClosed               = 0
	LocalSlaveReadPreferred            = 1
	LocalSlaveReadForce                = 2
)

func (s *StatusCode) String() string {
	r := "StatusUp"
	if *s == StatusDown {
		r = "StatusDown"
	}

	return fmt.Sprintf(r)
}

type DBInfo struct {
	ConnPool   []ConnectionPool
	Balancer   *balancer
	StatusMap  *sync.Map
	Datacenter []string
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

// Slice means one slice of the mysql cluster
type Slice struct {
	Cfg models.Slice
	sync.RWMutex

	Master          *DBInfo
	Slave           *DBInfo
	StatisticSlave  *DBInfo
	ProxyDatacenter string
	charset         string
	collationID     mysql.CollationID
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
				log.Warn("get connection from slave failed, try to get from master, error: %s", err.Error())
				pc, err = s.GetMasterConn()
			}
		}
	} else {
		pc, err = s.GetMasterConn()
	}
	if err != nil {
		log.Warn("get connection from backend failed, error: %s", err.Error())
		return
	}
	return
}

func (s *Slice) GetDirectConn(addr string) (*DirectConnection, error) {
	return NewDirectConnection(addr, s.Cfg.UserName, s.Cfg.Password, "", s.charset, s.collationID, s.Cfg.Capability)
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

			_, err := checkInstanceStatus(name, cp)

			if time.Now().Unix()-cp.GetLastChecked() >= int64(downAfterNoAlive) {
				s.SetMasterStatus(StatusDown)
				log.Warn("[ns:%s, %s:%s] check master StatusDown for %ds. err: %s", name, s.Cfg.Name, cp.Addr(), time.Now().Unix()-cp.GetLastChecked(), err)
				continue
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

				pc, err := checkInstanceStatus(name, cp)

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
					log.Warn("[ns:%s, %s:%s] skip check slave sync, get nil conn", name, s.Cfg.Name, cp.Addr())
					continue
				}

				if alive, err := checkSlaveSyncStatus(pc, secondBehindMaster); !alive {
					db.SetStatus(idx, StatusDown)
					log.Warn("[ns:%s, %s:%s] check slave StatusDown. sync err:%s", name, s.Cfg.Name, cp.Addr(), err)
					continue
				}

				db.SetStatus(idx, StatusUp)
				if oldStatus == StatusDown {
					log.Warn("[ns:%s, %s:%s] check slave StatusUp", name, s.Cfg.Name, cp.Addr())
				}
			}
		}
	}
}

func checkInstanceStatus(name string, cp ConnectionPool) (PooledConnect, error) {
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

// GetSlaveConn get connection from salve
func (s *Slice) GetSlaveConn(slavesInfo *DBInfo, localSlaveReadPriority int) (PooledConnect, error) {
	if len(slavesInfo.ConnPool) == 0 || allSlaveIsOffline(slavesInfo.StatusMap) {
		return nil, errors.ErrNoSlaveDB
	}
	var index int
	partialFoundIndex, foundIndex := -1, -1
	// find the idx of the ConnPool that isn't mark as down
	for size := len(slavesInfo.ConnPool); size > 0; size-- {
		s.Lock()
		var err error
		index, err = slavesInfo.Balancer.next()
		s.Unlock()
		if err != nil {
			return nil, err
		}

		if status, err := slavesInfo.GetStatus(index); err != nil {
			log.Debug("get slave status addr:%s,err:%s", slavesInfo.ConnPool[index].Addr(), err)
			continue
		} else if status == StatusDown {
			log.Debug("get slave status err or down,addr:%s", slavesInfo.ConnPool[index].Addr())
			continue
		}

		// partial found slave cause slave status StatusUP
		partialFoundIndex = index

		// check localSlaveReadPriority and update foundIndex
		if localSlaveReadPriority == LocalSlaveReadClosed {
			foundIndex = partialFoundIndex
			break
		}
		// check datacenter
		if slavesInfo.ConnPool[index].Datacenter() == s.ProxyDatacenter {
			foundIndex = index
			break
		}
	}
	if foundIndex >= 0 {
		return slavesInfo.ConnPool[foundIndex].Get(context.TODO())
	}
	if partialFoundIndex >= 0 && localSlaveReadPriority != LocalSlaveReadForce {
		return slavesInfo.ConnPool[partialFoundIndex].Get(context.TODO())
	}
	return nil, fmt.Errorf("get backend conn error,no local datacenter slaves")
}

// Close close the pool in slice
func (s *Slice) Close() error {
	s.Lock()
	defer s.Unlock()
	// close master
	for i := range s.Master.ConnPool {
		s.Master.ConnPool[i].Close()
	}

	// close slaves
	for i := range s.Slave.ConnPool {
		s.Slave.ConnPool[i].Close()
	}

	// close statistic slaves
	for i := range s.StatisticSlave.ConnPool {
		s.StatisticSlave.ConnPool[i].Close()
	}

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
	connectionPool := NewConnectionPool(masterStr, s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID, s.Cfg.Capability, s.Cfg.InitConnect, dc)
	if err := connectionPool.Open(); err != nil {
		return err
	}

	status := &sync.Map{}
	status.Store(0, StatusUp)

	s.Master = &DBInfo{[]ConnectionPool{connectionPool}, nil, status, []string{dc}}
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
	slaveWeights := make([]int, 0, count)
	datacenter := make([]string, 0, count)

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
		slaveWeights = append(slaveWeights, weight)
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
		datacenter = append(datacenter, dc)

		cp := NewConnectionPool(addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID, s.Cfg.Capability, s.Cfg.InitConnect, dc)
		if err = cp.Open(); err != nil {
			return nil, err
		}
		connPool = append(connPool, cp)
	}

	if len(slaveWeights) == 0 {
		return &DBInfo{}, nil
	}
	slaveBalancer := newBalancer(slaveWeights, len(connPool))
	StatusMap := &sync.Map{}
	for idx := range connPool {
		StatusMap.Store(idx, StatusUp)
	}

	return &DBInfo{connPool, slaveBalancer, StatusMap, datacenter}, nil
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

	slaveStatus, err := GetSlaveStatus(pc)
	if err != nil {
		return false, err
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

func GetSlaveStatus(conn PooledConnect) (SlaveStatus, error) {
	var slaveStatus SlaveStatus
	res, err := conn.Execute("show slave status;", 0)
	if err != nil {
		return slaveStatus, err
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
	return slaveStatus, err
}
