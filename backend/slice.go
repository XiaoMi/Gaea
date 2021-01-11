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
	Cfg models.Slice
	sync.RWMutex

	Master ConnectionPool

	Slave         []ConnectionPool
	slaveBalancer *balancer

	StatisticSlave         []ConnectionPool
	statisticSlaveBalancer *balancer

	charset     string
	collationID mysql.CollationID
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

func (s *Slice) GetDirectConn(addr string) (*DirectConnection, error) {
	return NewDirectConnection(addr, s.Cfg.UserName, s.Cfg.Password, "", s.charset, s.collationID)
}

// GetMasterConn return a connection in master pool
func (s *Slice) GetMasterConn() (PooledConnect, error) {
	ctx := context.TODO()
	return s.Master.Get(ctx)
}

// GetSlaveConn return a connection in slave pool
func (s *Slice) GetSlaveConn() (PooledConnect, error) {
	if len(s.Slave) == 0 {
		return nil, errors.ErrNoDatabase
	}

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
	if len(s.StatisticSlave) == 0 {
		return nil, errors.ErrNoDatabase
	}

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
	slaveWeights := make([]int, 0, count)

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
		slaveWeights = append(slaveWeights, weight)
		idleTimeout, err := util.Int2TimeDuration(s.Cfg.IdleTimeout)
		if err != nil {
			return err
		}
		cp := NewConnectionPool(addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID)
		cp.Open()
		s.Slave = append(s.Slave, cp)
	}
	s.slaveBalancer = newBalancer(slaveWeights, len(s.Slave))
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
	statisticSlaveWeights := make([]int, 0, count)

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
		statisticSlaveWeights = append(statisticSlaveWeights, weight)
		idleTimeout, err := util.Int2TimeDuration(s.Cfg.IdleTimeout)
		if err != nil {
			return err
		}
		cp := NewConnectionPool(addrAndWeight[0], s.Cfg.UserName, s.Cfg.Password, "", s.Cfg.Capacity, s.Cfg.MaxCapacity, idleTimeout, s.charset, s.collationID)
		cp.Open()
		s.StatisticSlave = append(s.StatisticSlave, cp)
	}
	s.statisticSlaveBalancer = newBalancer(statisticSlaveWeights, len(s.StatisticSlave))
	return nil
}

// SetCharsetInfo set charset
func (s *Slice) SetCharsetInfo(charset string, collationID mysql.CollationID) {
	s.charset = charset
	s.collationID = collationID
}
