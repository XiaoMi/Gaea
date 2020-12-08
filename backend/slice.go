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

const (
	weightSplit = "@"

	// DefaultSlice means default slice for namespace
	DefaultSlice = "slice-0"
)

// Slice means one slice of the mysql cluster
type Slice struct {
	sync.RWMutex
	Namespace string
	Cfg       models.Slice

	Master                 ConnectionPool
	Slave                  []ConnectionPool
	slaveBalancer          *balancer
	slaveWeights           []int
	StatisticSlave         []ConnectionPool
	statisticSlaveBalancer *balancer
	statisticSlaveWeights  []int

	charset             string
	collationID         mysql.CollationID
	statManager         *statManager
	healthCheckInterval int
	stopHealthCheck     chan struct{}
}

// GetSliceName return name of slice
func (s *Slice) GetSliceName() string {
	return s.Cfg.Name
}

// GetConn get backend connection from different node based on fromSlave and userType
func (s *Slice) GetConn(fromSlave bool, userType int) (pc PooledConnect, err error) {
	if fromSlave {
		if userType == models.StatisticUser {
			pc, err = s.GetStatisticSlaveConn()
			if err != nil {
				return nil, err
			}
		} else {
			pc, err = s.GetSlaveConn()
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

// GetMasterConn return a connection in master pool
func (s *Slice) GetMasterConn() (PooledConnect, error) {
	ctx := context.TODO()
	return s.Master.Get(ctx)
}

// GetSlaveConn return a connection in slave pool
func (s *Slice) GetSlaveConn() (PooledConnect, error) {
	s.Lock()
	index, err := s.slaveBalancer.next()
	s.Unlock()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	return s.Slave[index].Get(ctx)
}

// GetStatisticSlaveConn return a connection in statistic slave pool
func (s *Slice) GetStatisticSlaveConn() (PooledConnect, error) {
	s.Lock()
	index, err := s.statisticSlaveBalancer.next()
	s.Unlock()
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	return s.StatisticSlave[index].Get(ctx)
}

// Close close the pool in slice
func (s *Slice) Close() error {
	s.Lock()
	defer s.Unlock()

	//stop health check
	if s.statManager != nil {
		close(s.stopHealthCheck)
	}

	// close master
	s.Master.Close()

	// close slaves
	for i := range s.Slave {
		s.Slave[i].Close()
	}

	// close statistic slaves
	for i := range s.StatisticSlave {
		s.StatisticSlave[i].Close()
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
	s.Master = NewConnectionPool(masterStr, s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID)
	s.Master.Open()
	return nil
}

// ParseSlave create connection pool of slaves
// (127.0.0.1:3306@2,192.168.0.12:3306@3)
func (s *Slice) ParseSlave(slaves []string) error {
	if len(slaves) == 0 {
		return nil
	}

	var err error
	var weight int

	count := len(slaves)
	s.Slave = make([]ConnectionPool, 0, count)
	s.slaveWeights = make([]int, 0, count)

	//parse addr and weight
	for i := 0; i < count; i++ {
		addrAndWeight := strings.Split(slaves[i], weightSplit)
		if len(addrAndWeight) == 2 {
			weight, err = strconv.Atoi(addrAndWeight[1])
			if err != nil {
				return err
			}
		} else {
			weight = 1
		}
		s.slaveWeights = append(s.slaveWeights, weight)
		idleTimeout, err := util.Int2TimeDuration(s.Cfg.IdleTimeout)
		if err != nil {
			return err
		}
		cp := NewConnectionPool(addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID)
		cp.Open()
		s.Slave = append(s.Slave, cp)
	}
	s.slaveBalancer = newBalancer(s.slaveWeights, len(s.Slave))
	return nil
}

// ParseStatisticSlave create connection pool of statistic slaves
// slaveStr(127.0.0.1:3306@2,192.168.0.12:3306@3)
func (s *Slice) ParseStatisticSlave(statisticSlaves []string) error {
	if len(statisticSlaves) == 0 {
		return nil
	}

	var err error
	var weight int

	count := len(statisticSlaves)
	s.StatisticSlave = make([]ConnectionPool, 0, count)
	s.statisticSlaveWeights = make([]int, 0, count)

	//parse addr and weight
	for i := 0; i < count; i++ {
		addrAndWeight := strings.Split(statisticSlaves[i], weightSplit)
		if len(addrAndWeight) == 2 {
			weight, err = strconv.Atoi(addrAndWeight[1])
			if err != nil {
				return err
			}
		} else {
			weight = 1
		}
		s.statisticSlaveWeights = append(s.statisticSlaveWeights, weight)
		idleTimeout, err := util.Int2TimeDuration(s.Cfg.IdleTimeout)
		if err != nil {
			return err
		}
		cp := NewConnectionPool(addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID)
		cp.Open()
		s.StatisticSlave = append(s.StatisticSlave, cp)
	}
	s.statisticSlaveBalancer = newBalancer(s.statisticSlaveWeights, len(s.StatisticSlave))
	return nil
}

// SetCharsetInfo set charset
func (s *Slice) SetCharsetInfo(charset string, collationID mysql.CollationID) {
	s.charset = charset
	s.collationID = collationID
}

func (s *Slice) CreateSliceStatManager(cfg *models.HealthCheckConfig) {
	s.statManager = newStatManager(cfg, s)
	s.stopHealthCheck = make(chan struct{})
	s.healthCheckInterval = cfg.IntervalSeconds
}

func (s *Slice) StartRefreshSliceStat() {
	go func() {
		ticker := time.NewTicker(time.Duration(s.healthCheckInterval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				slaveWeights := make([]int, 0)
				statisticSlaveWeights := make([]int, 0)

				s.statManager.getMaster()
				for i, slave := range s.Slave {
					if !s.statManager.getSlave(i) {
						slaveWeights = append(slaveWeights, 0)
						log.Warn("namespace:%s slice:%s slave:%v is invalid", s.Namespace, s.GetSliceName(), slave.Addr())
						continue
					}
					slaveWeights = append(slaveWeights, s.slaveWeights[i])
				}
				for i, statisticSlave := range s.StatisticSlave {
					if !s.statManager.getStatisticSlave(i) {
						statisticSlaveWeights = append(statisticSlaveWeights, 0)
						log.Warn("namespace:%s slice:%s statisticSlave:%v is invalid", s.Namespace, s.GetSliceName(), statisticSlave.Addr())
						continue
					}
					statisticSlaveWeights = append(statisticSlaveWeights, s.statisticSlaveWeights[i])
				}

				s.Lock()
				s.slaveBalancer = newBalancer(slaveWeights, len(s.Slave))
				s.statisticSlaveBalancer = newBalancer(statisticSlaveWeights, len(s.StatisticSlave))
				s.Unlock()
			case <-s.stopHealthCheck:
				return
			}
		}
	}()
}
