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

	UP                      StatusCode = 1
	DOWN                    StatusCode = 0
	LocalSlaveReadClosed               = 0
	LocalSlaveReadPreferred            = 1
	LocalSlaveReadForce                = 2
)

func (s *StatusCode) String() string {
	r := "UP"
	if *s == DOWN {
		r = "DOWN"
	}

	return fmt.Sprintf(r)
}

type DBInfo struct {
	ConnPool   []ConnectionPool
	Balancer   *balancer
	StatusMap  sync.Map
	Datacenter []string
}

func (dbi *DBInfo) GetStatus(index int) (StatusCode, error) {
	if value, ok := dbi.StatusMap.Load(index); ok {
		return value.(StatusCode), nil
	}
	return DOWN, fmt.Errorf("can't get status of index:%d", index)
}

func (dbi *DBInfo) GetMasterStatus(index int) (StatusCode, error) {
	return dbi.GetStatus(index)
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
	if v, _ := s.Master.StatusMap.Load(0); v != UP {
		return nil, fmt.Errorf("master:%s is Down", s.Cfg.Master)
	}

	ctx := context.TODO()
	return s.Master.ConnPool[0].Get(ctx)
}

func allSlaveIsOffline(SlaveStatusMap *sync.Map) bool {
	var result = true
	SlaveStatusMap.Range(func(k, v interface{}) bool {
		if v == UP {
			result = false
			return false
		}
		return true
	})

	return result
}

// GetSlaveConn get connection from salve
func (s *Slice) GetSlaveConn(slavesInfo *DBInfo, localSlaveReadPriority int) (PooledConnect, error) {
	if len(slavesInfo.ConnPool) == 0 || allSlaveIsOffline(&slavesInfo.StatusMap) {
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
		} else if status == DOWN {
			log.Debug("get slave status err or down,addr:%s", slavesInfo.ConnPool[index].Addr())
			continue
		}

		// partial found slave cause slave status UP
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

	status := sync.Map{}
	status.Store(0, UP)

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
	StatusMap := sync.Map{}
	for idx := range connPool {
		StatusMap.Store(idx, UP)
	}

	return &DBInfo{connPool, slaveBalancer, StatusMap, datacenter}, nil
}

// SetCharsetInfo set charset
func (s *Slice) SetCharsetInfo(charset string, collationID mysql.CollationID) {
	s.charset = charset
	s.collationID = collationID
}
