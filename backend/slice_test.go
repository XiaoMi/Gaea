// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"sync"
	"testing"
	"time"

	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
	"github.com/bytedance/mockey"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestParseDBInfo(t *testing.T) {
	testCases := []struct {
		name              string
		addrs             []string
		expectAddrs       []string
		expectDatacenters []string
		expectWeights     []int
		expectErr         bool
	}{

		{
			name:              "valid slave with weight and datacenter",
			addrs:             []string{"c3-mysql-slave01:3307@10#bj", "c4-mysql-slave02:3308@5#sgp"},
			expectAddrs:       []string{"c3-mysql-slave01:3307", "c4-mysql-slave02:3308"},
			expectDatacenters: []string{"bj", "sgp"},
			expectWeights:     []int{10, 5},
			expectErr:         false,
		},
		{
			name:              "valid slave without weight, default to 1",
			addrs:             []string{"c3-mysql-slave01:3307#bj", "c4-mysql-slave02:3308#sgp"},
			expectAddrs:       []string{"c3-mysql-slave01:3307", "c4-mysql-slave02:3308"},
			expectDatacenters: []string{"bj", "sgp"},
			expectWeights:     []int{1, 1},
			expectErr:         false,
		},

		{
			name:              "invalid weight format, should default to 1",
			addrs:             []string{"c3-mysql-slave01:3307@xyz#bj", "c4-mysql-slave02:3308@!@#sgp"},
			expectAddrs:       []string{"c3-mysql-slave01:3307", "c4-mysql-slave02:3308"},
			expectDatacenters: []string{"bj", "sgp"},
			expectWeights:     []int{1, 1}, // 非法权重默认设为 1
			expectErr:         false,
		},

		{
			name:              "empty input",
			addrs:             []string{},
			expectAddrs:       []string{},
			expectDatacenters: []string{},
			expectWeights:     []int{},
			expectErr:         false,
		},

		{
			name:              "invalid format, should return error",
			addrs:             []string{"invalid-slave"},
			expectAddrs:       []string{"invalid-slave"},
			expectDatacenters: []string{""},
			expectWeights:     []int{1},
			expectErr:         false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := new(Slice)
			err := s.ParseSlave(tt.addrs)
			assert.NoError(t, err)
			assert.NotNil(t, s.Slave)
			err = s.ParseMonitorSlave(tt.addrs)
			assert.NoError(t, err)
			assert.NotNil(t, s.MonitorSlave)

			slave := s.Slave
			monitorSlave := s.MonitorSlave
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, s.Slave)
				assert.Equal(t, len(tt.expectAddrs), len(s.Slave.Nodes))

				for i := range tt.expectAddrs {
					assert.Equal(t, tt.expectAddrs[i], slave.Nodes[i].Address)
					assert.Equal(t, tt.expectDatacenters[i], slave.Nodes[i].Datacenter)
					assert.Equal(t, tt.expectWeights[i], slave.Nodes[i].Weight)

					assert.Equal(t, tt.expectAddrs[i], monitorSlave.Nodes[i].Address)
					assert.Equal(t, tt.expectDatacenters[i], monitorSlave.Nodes[i].Datacenter)
					assert.Equal(t, tt.expectWeights[i], monitorSlave.Nodes[i].Weight)
				}
			}
		})
	}
}

func TestParseSlave(t *testing.T) {
	testCases := []struct {
		name              string
		slaveAdders       []string
		slaveStatus       []StatusCode
		expectAddrs       []string
		expectDatacenters []string
	}{
		{
			name:              "test simple",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"c3", "c3", "c4"},
		},
		{
			name:              "test simple 2",
			slaveAdders:       []string{"c3-mysql-test00:3306", "c3-mysql-t:3308", "c4-my:3310"},
			expectAddrs:       []string{"c3-mysql-test00:3306", "c3-mysql-t:3308", "c4-my:3310"},
			expectDatacenters: []string{"c3", "c3", "c4"},
		},
		{
			name:              "test with weight",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306@10", "c3-mysql-test01.bj:3308@5", "c4-mysql-test02.bj:3310"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"c3", "c3", "c4"},
		},
		{
			name:              "test with weight and datacenter",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306@10#bj", "c3-mysql-test01.bj:3308@5#sgp", "c4-mysql-test02.bj:3310@3#c4"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"bj", "sgp", "c4"},
		},
		{
			name:              "test with datacenter",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306#bj", "c3-mysql-test01.bj:3308#sgp", "c4-mysql-test02.bj:3310#c4"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"bj", "sgp", "c4"},
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := new(Slice)
			err := s.ParseSlave(tt.slaveAdders)
			assert.Nil(t, err)
			for i := range tt.expectAddrs {
				assert.Equal(t, tt.expectAddrs[i], s.Slave.Nodes[i].Address)
				assert.Equal(t, tt.expectDatacenters[i], s.Slave.Nodes[i].Datacenter)
			}
		})
	}
}

func TestParseMonitorSlave(t *testing.T) {
	testCases := []struct {
		name              string
		slaveAdders       []string
		slaveStatus       []StatusCode
		expectAddrs       []string
		expectDatacenters []string
	}{
		{
			name:              "test simple",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"c3", "c3", "c4"},
		},
		{
			name:              "test simple 2",
			slaveAdders:       []string{"c3-mysql-test00:3306", "c3-mysql-t:3308", "c4-my:3310"},
			expectAddrs:       []string{"c3-mysql-test00:3306", "c3-mysql-t:3308", "c4-my:3310"},
			expectDatacenters: []string{"c3", "c3", "c4"},
		},
		{
			name:              "test with weight",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306@10", "c3-mysql-test01.bj:3308@5", "c4-mysql-test02.bj:3310"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"c3", "c3", "c4"},
		},
		{
			name:              "test with weight and datacenter",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306@10#bj", "c3-mysql-test01.bj:3308@5#sgp", "c4-mysql-test02.bj:3310@3#c4"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"bj", "sgp", "c4"},
		},
		{
			name:              "test with datacenter",
			slaveAdders:       []string{"c3-mysql-test00.bj:3306#bj", "c3-mysql-test01.bj:3308#sgp", "c4-mysql-test02.bj:3310#c4"},
			expectAddrs:       []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectDatacenters: []string{"bj", "sgp", "c4"},
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := new(Slice)
			err := s.ParseMonitorSlave(tt.slaveAdders)
			if err != nil {
				t.Fatal(err)
			}
			for i := range tt.expectAddrs {
				assert.Equal(t, tt.expectAddrs[i], s.MonitorSlave.Nodes[i].ConnPool.Addr())
				assert.Equal(t, tt.expectDatacenters[i], s.MonitorSlave.Nodes[i].Datacenter)
			}
		})
	}
}

func TestParseSlaveWithWeights(t *testing.T) {
	testCases := []struct {
		name           string
		slaveAdders    []string
		expectAddrs    []string
		expectWeights  []int
		expectIndices  []int
		expectingError bool // 是否预期解析失败
	}{

		// ✅ 测试正常权重解析
		{
			name:           "valid weights",
			slaveAdders:    []string{"c3-mysql-test00.bj:3306@10", "c3-mysql-test01.bj:3308@5", "c4-mysql-test02.bj:3310@3"},
			expectAddrs:    []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectWeights:  []int{10, 5, 3},
			expectIndices:  []int{0, 1, 2},
			expectingError: false,
		},

		// ✅ 测试默认权重解析（未指定权重默认为 1）
		{
			name:           "default weights",
			slaveAdders:    []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectAddrs:    []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectWeights:  []int{1, 1, 1}, // 默认权重
			expectIndices:  []int{0, 1, 2},
			expectingError: false,
		},

		// ✅ 测试负数权重（负数可以被检测，在创建balancer之前被过滤掉了）
		{
			name:           "negative weight",
			slaveAdders:    []string{"c3-mysql-test00.bj:3306@-5", "c3-mysql-test01.bj:3308@3"},
			expectAddrs:    []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
			expectIndices:  []int{0, 1},
			expectWeights:  []int{1, 3}, // 默认权重
			expectingError: false,
		},

		// ✅ 测试 0 权重（应当过滤掉）
		{
			name:          "zero weight filtered",
			slaveAdders:   []string{"c3-mysql-test00.bj:3306@0", "c3-mysql-test01.bj:3308@4", "c4-mysql-test02.bj:3310@0"},
			expectAddrs:   []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			expectIndices: []int{0, 1, 2},
			expectWeights: []int{0, 4, 0},

			expectingError: false,
		},

		// ✅ 测试混合非法和合法的情况
		{
			name:           "mixed valid and invalid weights",
			slaveAdders:    []string{"c3-mysql-test00.bj:3306@2", "c3-mysql-test01.bj:3308@-1", "c4-mysql-test02.bj:3310@0"},
			expectAddrs:    []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"}, // 过滤掉权重为 0 和负数的
			expectIndices:  []int{0, 1, 2},
			expectWeights:  []int{2, 1, 0},
			expectingError: false,
		},

		// ✅ 测试混合非法字符权重测试
		{
			name:           "invalid weights",
			slaveAdders:    []string{"c3-mysql-test00.bj:3306@&", "c3-mysql-test01.bj:3308@+", "c4-mysql-test02.bj:3310@-"},
			expectAddrs:    []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"}, // 过滤掉权重为 0 和负数的
			expectIndices:  []int{0, 1, 2},
			expectWeights:  []int{1, 1, 1},
			expectingError: false,
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := new(Slice)
			err := s.ParseSlave(tt.slaveAdders)
			assert.Nil(t, err)

			if tt.expectingError {
				assert.Error(t, err, "Expected an error but got none")
				return
			} else {
				assert.NoError(t, err, "Unexpected error during parsing")
			}

			// 检查解析出的地址、权重和索引
			var actualAddrs []string
			var actualWeights []int
			var actualIndices []int
			for i, node := range s.Slave.Nodes {
				actualAddrs = append(actualAddrs, node.Address)
				actualWeights = append(actualWeights, node.Weight)
				actualIndices = append(actualIndices, i) // 直接使用 i 作为索引
			}

			assert.Equal(t, tt.expectAddrs, actualAddrs, "Address mismatch")
			assert.Equal(t, tt.expectWeights, actualWeights, "Weight mismatch")
			assert.Equal(t, tt.expectIndices, actualIndices, "Index mismatch")
		})
	}
}

func TestGetSlaveConnWhenLocalSlaveReadClosed(t *testing.T) {
	testCases := []struct {
		name                   string
		localSlaveReadPriority int
		proxyDc                string
		slaveAdders            []string
		slaveStatus            []StatusCode
		getCounts              int
		weights                []int
		expectAddrs            map[string]int
		expectErr              int
	}{
		{
			name:                   "local slave all up, remote slave all up, equal weigh",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              6,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 1},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 2,
				"c3-mysql-test01.bj:3308": 2,
				"c4-mysql-test02.bj:3310": 2,
			},
			expectErr: 0,
		},
		{
			name:                   "local slave all up, remote slave all up, not equal weight",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			weights:                []int{2, 1, 1},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 6,
				"c3-mysql-test01.bj:3308": 3,
				"c4-mysql-test02.bj:3310": 3,
			},
			expectErr: 0,
		},
		{
			name:                   "local slave all up, remote slave all up, not equal weight",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 2, 2},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 2,
				"c3-mysql-test01.bj:3308": 2,
				"c4-mysql-test01.bj:3310": 4,
				"c4-mysql-test02.bj:3310": 4,
			},
			expectErr: 0,
		},
		{
			name:                   "local slave all down, remote slave all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusDown, StatusDown, StatusUp, StatusUp},
			weights:                []int{1, 1, 2, 2},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 6,
				"c4-mysql-test02.bj:3310": 6,
			},
			expectErr: 0,
		},
		{
			name:                   "local slave part down, remote slave all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              10,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusDown, StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 2, 2},
			expectAddrs: map[string]int{
				"c3-mysql-test01.bj:3308": 2,
				"c4-mysql-test01.bj:3310": 4,
				"c4-mysql-test02.bj:3310": 4,
			},
			expectErr: 0,
		},
		{
			name:                   "remote slave all down, local slave all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusDown},
			weights:                []int{1, 1, 2, 2},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 6,
				"c3-mysql-test01.bj:3308": 6,
			},
			expectErr: 0,
		},
		{
			name:                   "remote slave part down, local slave all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusUp},
			weights:                []int{1, 1, 2, 2},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 3,
				"c3-mysql-test01.bj:3308": 3,
				"c4-mysql-test02.bj:3310": 6,
			},
			expectErr: 0,
		},
		{
			name:                   "unknown dc c5, all slave up",
			proxyDc:                "c5",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 2, 2},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 2,
				"c3-mysql-test01.bj:3308": 2,
				"c4-mysql-test01.bj:3310": 4,
				"c4-mysql-test02.bj:3310": 4,
			},
			expectErr: 0,
		},

		{
			name:                   "unknown dc c5, part slave up",
			proxyDc:                "c5",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp, StatusDown},
			weights:                []int{1, 1, 2, 2},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 3,
				"c3-mysql-test01.bj:3308": 3,
				"c4-mysql-test01.bj:3310": 6,
			},
			expectErr: 0,
		},
		{
			name:                   "unknown dc c5, all slave down",
			proxyDc:                "c5",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusDown, StatusDown, StatusDown, StatusDown},
			weights:                []int{1, 1, 2, 2},
			expectAddrs:            nil,
			expectErr:              12,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			dbInfo, err := generateDBInfoWithWeights(mockCtl, tt.slaveAdders, tt.slaveStatus, tt.weights)
			assert.Nil(t, err)
			s := &Slice{Slave: dbInfo}
			s.ProxyDatacenter = tt.proxyDc
			err = dbInfo.InitBalancers(s.ProxyDatacenter)
			assert.Nil(t, err)
			actualError := 0
			actualSuccess := map[string]int{}
			for j := 0; j < tt.getCounts; j++ {
				cp, err := s.GetSlaveConn(dbInfo, tt.localSlaveReadPriority)
				if err != nil {
					actualError++
				} else {
					actualSuccess[cp.GetAddr()]++
				}
			}
			for actAdd, actSucc := range actualSuccess {
				assert.Equal(t, actSucc, tt.expectAddrs[actAdd], fmt.Errorf("case: %s, actual address %s acualCount %d  expect count %d", tt.name, actAdd, actSucc, tt.expectAddrs[actAdd]))
			}
			if actualError != tt.expectErr {
				assert.Equal(t, actualError, tt.expectErr, fmt.Errorf("case: %s, actual error %d  expect error %d", tt.name, actualError, tt.expectErr))
			}
		})
	}

}

func TestGetSlaveConnWhenLocalSlaveReadPrefer(t *testing.T) {
	testCases := []struct {
		name                   string
		localSlaveReadPriority int
		proxyDc                string
		slaveAdders            []string
		slaveStatus            []StatusCode
		weights                []int
		getCounts              int
		expectAddrs            map[string]int
		expectErr              int
	}{
		// 本地和非本地上都只有一个slave，且都是up
		{
			name:                   "local slave all up, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test01.bj:3308"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp},
			weights:                []int{1, 3},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3308": 4,
			},
			expectErr: 0,
		},

		// 本地机房全部UP
		{
			name:                   "local slave all up, remote slave all up, equal weight",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              6,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 1},
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 6,
			},
			expectErr: 0,
		},
		{
			name:                   "local slave all up, remote slave all up, equal weight, equal number",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 1, 1},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 2,
				"c4-mysql-test02.bj:3310": 2,
			},
			expectErr: 0,
		},
		{
			name:                   "local slave all up, remote slave all up, not equal weight",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              10,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 4, 4},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 5,
				"c4-mysql-test02.bj:3310": 5,
			},
			expectErr: 0,
		},
		// 本地机房全部宕机
		{
			name:                   "local slave all down, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown},
			weights:                []int{1, 1, 1},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 6,
				"c3-mysql-test01.bj:3308": 6,
			},
			expectErr: 0,
		},
		// 本地机房部分宕机
		{
			name:                   "local slave part down, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              10,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310", "c4-mysql-test01.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusUp},
			weights:                []int{1, 1, 4, 4},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 10,
			},
			expectErr: 0,
		},
		{
			name:                   "local slave part down, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              6,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310", "c4-mysql-test01.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusUp},
			weights:                []int{1, 1, 3, 1},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 6,
			},
			expectErr: 0,
		},

		// 本地优先，但本地没有可用的slave
		{
			name:                   "no local slave, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              2,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp},
			weights:                []int{1, 1},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 1,
				"c3-mysql-test01.bj:3308": 1,
			},
			expectErr: 0,
		},

		// 本地全部宕机down 非本地 up
		{
			name:                   "local slave down, remote slave up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test01.bj:3308"},
			slaveStatus:            []StatusCode{StatusUp, StatusDown},
			weights:                []int{1, 3},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 4,
			},
			expectErr: 0,
		},

		// 本地和非本地上都只有一个slave，且都是down
		{
			name:                   "local slave down, remote slave down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test01.bj:3308"},
			slaveStatus:            []StatusCode{StatusDown, StatusDown},
			weights:                []int{1, 3},
			expectAddrs:            nil,
			expectErr:              4,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			dbInfo, err := generateDBInfoWithWeights(mockCtl, tt.slaveAdders, tt.slaveStatus, tt.weights)
			assert.Nil(t, err)
			s := &Slice{Slave: dbInfo}
			s.ProxyDatacenter = tt.proxyDc
			err = dbInfo.InitBalancers(s.ProxyDatacenter)
			assert.Nil(t, err)
			actualError := 0
			actualSuccess := map[string]int{}

			for j := 0; j < tt.getCounts; j++ {
				cp, err := s.GetSlaveConn(dbInfo, tt.localSlaveReadPriority)
				if err != nil {
					actualError++
				} else {
					actualSuccess[cp.GetAddr()]++
				}
			}
			for actAdd, actSucc := range actualSuccess {
				assert.Equal(t, actSucc, tt.expectAddrs[actAdd], fmt.Errorf("case: %s, actual address %s acualCount %d  expect count %d", tt.name, actAdd, actSucc, tt.expectAddrs[actAdd]))
			}
			if actualError != tt.expectErr {
				assert.Equal(t, actualError, tt.expectErr, fmt.Errorf("case: %s, actual error %d  expect error %d", tt.name, actualError, tt.expectErr))
			}

		})
	}

}

func TestGetSlaveConnWhenLocalSlaveReadForce(t *testing.T) {
	testCases := []struct {
		name                   string
		localSlaveReadPriority int
		proxyDc                string
		slaveAdders            []string
		slaveStatus            []StatusCode
		weights                []int
		getCounts              int
		expectGetErrorCounts   int
		expectAddrs            map[string]int
	}{
		// 本地机房全部UP
		{
			name:                   "local slave up, remote slave up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              6,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 1},
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 6,
			},
			expectGetErrorCounts: 0,
		},
		{
			name:                   "local slave all up, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 1, 1},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 2,
				"c4-mysql-test02.bj:3310": 2,
			},
			expectGetErrorCounts: 0,
		},
		{
			name:                   "local slave all up, remote slave all up, not equal weight",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              10,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test01.bj:3310", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp, StatusUp},
			weights:                []int{1, 1, 4, 4},
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 5,
				"c4-mysql-test02.bj:3310": 5,
			},
			expectGetErrorCounts: 0,
		},

		// 本地机房全部宕机
		{
			name:                   "local slave all down, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown},
			weights:                []int{1, 1, 1},
			expectGetErrorCounts:   12,
			expectAddrs:            nil,
		},
		// 本地机房部分宕机
		{
			name:                   "local slave part down, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              10,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310", "c4-mysql-test01.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusUp},
			weights:                []int{1, 1, 4, 4},
			expectGetErrorCounts:   0,
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 10,
			},
		},
		{
			name:                   "local slave part down, remote slave all up",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310", "c4-mysql-test01.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusUp},
			weights:                []int{1, 1, 3, 1},
			expectGetErrorCounts:   0,
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 12,
			},
		},

		// 本地和异地都有一台机器down掉
		{
			name:                   "local slave one down, remote slave one down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310", "c4-mysql-test01.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusDown, StatusDown, StatusUp},
			weights:                []int{1, 1, 3, 1},
			expectGetErrorCounts:   0,
			expectAddrs: map[string]int{
				"c4-mysql-test01.bj:3310": 12,
			},
		},
		// 本地没有slave
		{
			name:                   "no local slave, remote slave part down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
			slaveStatus:            []StatusCode{StatusUp, StatusDown},
			weights:                []int{1, 1},
			expectGetErrorCounts:   12,
			expectAddrs:            nil,
		},

		// 本地没有slave
		{
			name:                   "no local slave, remote slave all down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              12,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
			slaveStatus:            []StatusCode{StatusDown, StatusDown},
			weights:                []int{1, 1},
			expectGetErrorCounts:   12,
			expectAddrs:            nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()
			dbInfo, err := generateDBInfoWithWeights(mockCtl, tt.slaveAdders, tt.slaveStatus, tt.weights)
			assert.Nil(t, err)
			s := &Slice{Slave: dbInfo}
			s.ProxyDatacenter = tt.proxyDc
			err = dbInfo.InitBalancers(s.ProxyDatacenter)
			assert.Nil(t, err)
			actualError := 0
			actualSuccess := map[string]int{}
			for j := 0; j < tt.getCounts; j++ {
				cp, err := s.GetSlaveConn(dbInfo, tt.localSlaveReadPriority)
				if err != nil {
					actualError++
				} else {
					actualSuccess[cp.GetAddr()]++
				}
			}
			for actAdd, actSucc := range actualSuccess {
				assert.Equal(t, actSucc, tt.expectAddrs[actAdd], fmt.Errorf("case: %s, actual address %s acualCount %d  expect count %d", tt.name, actAdd, actSucc, tt.expectAddrs[actAdd]))
			}
			if actualError != tt.expectGetErrorCounts {
				assert.Equal(t, actualError, tt.expectGetErrorCounts, fmt.Errorf("case: %s, actual error %d  expect error %d", tt.name, actualError, tt.expectGetErrorCounts))
			}
		})
	}
}

func TestConcurrentGetSlaveConnWhenLocalSlaveReadForce(t *testing.T) {

	testCases := []struct {
		name                   string
		localSlaveReadPriority int
		proxyDc                string
		slaveAdders            []string
		slaveStatus            []StatusCode
		weights                []int
		getCounts              int
		expectGetErrorCounts   int
		expectAddrs            map[string]int
		concurrency            int // 并发请求的 Goroutine 数
	}{
		// 🚀 新增测试用例：本地只有 1 台从库，并发请求，请求都成功
		{
			name:                   "concurrent access with single local slave",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              100000,                                                         // 总请求数
			concurrency:            100,                                                            // 10 个 Goroutine 并发请求
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310"}, // 只有 1 个本地从库
			slaveStatus:            []StatusCode{StatusUp, StatusUp},                               // 本地从库在线
			weights:                []int{1, 1},                                                    // 权重均等
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 100000, // 预计所有请求都打到这台服务器
			},
			expectGetErrorCounts: 0, // 理论上不应该报错
		},
		// 🚀 新增测试用例：本地 1 台从库，远程 2 台从库，强制本地，但本地挂了，所有请求失败
		{
			name:                   "concurrent access with single local slave down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              100000,
			concurrency:            100,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown}, // 本地从库宕机
			weights:                []int{1, 1, 1},
			expectAddrs:            nil,    // 无法获取可用连接
			expectGetErrorCounts:   100000, // 预计所有请求都会失败
		},

		// 🚀 新增测试用例：本地有 2 台从库，一台正常，一台宕机，强制本地，并发请求，所有请求成功
		{
			name:                   "concurrent access with single local slave",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              100000,                                                                                    // 总请求数
			concurrency:            100,                                                                                       // 10 个 Goroutine 并发请求
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"}, // 2 个本地从库,
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown},                                              // 1 个 up 一个 down
			weights:                []int{1, 1, 1},                                                                            // 权重均等
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 100000, // 预计所有请求都打到这台服务器
			},
			expectGetErrorCounts: 0, // 理论上不应该报错
		},

		// 🚀 新增测试用例：本地有 3 台从库，一台正常，两台宕机，强制本地，并发请求，所有请求成功
		{
			name:                   "concurrent access with single local slave",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              100000,                                                                                                               // 总请求数
			concurrency:            100,                                                                                                                  // 10 个 Goroutine 并发请求
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test03.bj:3310"}, // 3 个本地从库,
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusUp},                                                               // 2 个 up 一个 down
			weights:                []int{1, 1, 1, 1},                                                                                                    // 权重均等
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 50000, // 预计一半请求都打到这台服务器
				"c4-mysql-test03.bj:3310": 50000, // 预计一半请求都打到这台服务器
			},
			expectGetErrorCounts: 0, // 理论上不应该报错
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			// 生成数据库连接池
			dbInfo, err := generateDBInfoWithWeights(mockCtl, tt.slaveAdders, tt.slaveStatus, tt.weights)
			assert.Nil(t, err)
			s := &Slice{Slave: dbInfo}
			s.ProxyDatacenter = tt.proxyDc

			err = dbInfo.InitBalancers(s.ProxyDatacenter)
			assert.Nil(t, err)

			actualSuccess := make(map[string]int)
			actualError := 0
			var mu sync.Mutex

			// 并发请求
			var wg sync.WaitGroup
			concurrency := tt.concurrency
			requestsPerGoroutine := tt.getCounts / concurrency

			for i := 0; i < concurrency; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < requestsPerGoroutine; j++ {
						cp, err := s.GetSlaveConn(dbInfo, tt.localSlaveReadPriority)

						mu.Lock() // 保证统计数据线程安全
						if err != nil {
							actualError++
						} else {
							actualSuccess[cp.GetAddr()]++
						}
						mu.Unlock()
					}
				}()
			}
			wg.Wait() // 等待所有 Goroutine 完成
			// 校验成功的连接请求数
			for actAdd, actSucc := range actualSuccess {
				assert.Equal(t, tt.expectAddrs[actAdd], actSucc, fmt.Errorf("case: %s, actual address %s actualCount %d expect count %d",
					tt.name, actAdd, actSucc, tt.expectAddrs[actAdd]))
			}

			// 校验失败的请求数
			assert.Equal(t, actualError, tt.expectGetErrorCounts, fmt.Errorf("case: %s, actual error %d expect error %d",
				tt.name, actualError, tt.expectGetErrorCounts))
		})
	}
}

func TestConcurrentGetSlaveConnWhenLocalSlaveReadPrefer(t *testing.T) {

	testCases := []struct {
		name                   string
		localSlaveReadPriority int
		proxyDc                string
		slaveAdders            []string
		slaveStatus            []StatusCode
		weights                []int
		getCounts              int
		expectGetErrorCounts   int
		expectAddrs            map[string]int
		concurrency            int // 并发请求的 Goroutine 数
	}{
		// 🚀 新增测试用例：本地只有 1 台从库，并发请求，请求都成功
		{
			name:                   "concurrent access with single local slave",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              100000,                                                         // 总请求数
			concurrency:            100,                                                            // 10 个 Goroutine 并发请求
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310"}, // 只有 1 个本地从库
			slaveStatus:            []StatusCode{StatusUp, StatusUp},                               // 本地从库在线
			weights:                []int{1, 1},                                                    // 权重均等
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 100000, // 预计所有请求都打到这台服务器
			},
			expectGetErrorCounts: 0, // 理论上不应该报错
		},
		// 🚀 新增测试用例：本地 1 台从库，远程 2 台从库，强制本地，但本地挂了，所有请求失败
		{
			name:                   "concurrent access with single local slave down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              100000,
			concurrency:            100,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown}, // 本地从库宕机
			weights:                []int{1, 1, 1},
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3306": 50000, // 预计一半
				"c3-mysql-test01.bj:3308": 50000, // 预计另一半
			},
			expectGetErrorCounts: 0, // 预计所有请求都会失败
		},

		// 🚀 新增测试用例：本地有 2 台从库，一台正常，一台宕机，强制本地，并发请求，所有请求成功
		{
			name:                   "concurrent access with single local slave",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              100000,                                                                                    // 总请求数
			concurrency:            100,                                                                                       // 10 个 Goroutine 并发请求
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"}, // 2 个本地从库,
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown},                                              // 1 个 up 一个 down
			weights:                []int{1, 1, 1},                                                                            // 权重均等
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 100000, // 预计所有请求都打到这台服务器
			},
			expectGetErrorCounts: 0, // 理论上不应该报错
		},

		// 🚀 新增测试用例：本地有 3 台从库，一台正常，两台宕机，强制本地，并发请求，所有请求成功
		{
			name:                   "concurrent access with single local slave",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			getCounts:              100000,                                                                                                               // 总请求数
			concurrency:            100,                                                                                                                  // 10 个 Goroutine 并发请求
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test03.bj:3310"}, // 3 个本地从库,
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown, StatusUp},                                                               // 2 个 up 一个 down
			weights:                []int{1, 1, 1, 1},                                                                                                    // 权重均等
			expectAddrs: map[string]int{
				"c4-mysql-test02.bj:3310": 50000, // 预计一半请求都打到这台服务器
				"c4-mysql-test03.bj:3310": 50000, // 预计一半请求都打到这台服务器
			},
			expectGetErrorCounts: 0, // 理论上不应该报错
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			// 生成数据库连接池
			dbInfo, err := generateDBInfoWithWeights(mockCtl, tt.slaveAdders, tt.slaveStatus, tt.weights)
			assert.Nil(t, err)
			s := &Slice{Slave: dbInfo}
			s.ProxyDatacenter = tt.proxyDc

			err = dbInfo.InitBalancers(s.ProxyDatacenter)
			assert.Nil(t, err)

			actualSuccess := make(map[string]int)
			actualError := 0
			var mu sync.Mutex

			// 并发请求
			var wg sync.WaitGroup
			concurrency := tt.concurrency
			requestsPerGoroutine := tt.getCounts / concurrency

			for i := 0; i < concurrency; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < requestsPerGoroutine; j++ {
						cp, err := s.GetSlaveConn(dbInfo, tt.localSlaveReadPriority)

						mu.Lock() // 保证统计数据线程安全
						if err != nil {
							actualError++
						} else {
							actualSuccess[cp.GetAddr()]++
						}
						mu.Unlock()
					}
				}()
			}
			wg.Wait() // 等待所有 Goroutine 完成
			// 校验成功的连接请求数
			for actAdd, actSucc := range actualSuccess {
				assert.Equal(t, tt.expectAddrs[actAdd], actSucc, fmt.Errorf("case: %s, actual address %s actualCount %d expect count %d",
					tt.name, actAdd, actSucc, tt.expectAddrs[actAdd]))
			}

			// 校验失败的请求数
			assert.Equal(t, actualError, tt.expectGetErrorCounts, fmt.Errorf("case: %s, actual error %d expect error %d",
				tt.name, actualError, tt.expectGetErrorCounts))
		})
	}
}

func TestGetConn(t *testing.T) {
	testCases := []struct {
		name                   string
		proxyDc                string
		userType               int
		fromSlave              bool
		localSlaveReadPriority int
		master                 string
		masterStatus           StatusCode
		slaveAddrs             []string
		slaveStatus            []StatusCode
		weights                []int
		getCounts              int
		expectAddrs            map[string]int // 预期每个连接地址被命中的次数
		expectErr              int            // 预期的错误情况，0表示无错误
	}{
		{
			name:                   "MonitorUser, fromSlave true, local slave up, should get slave connection",
			proxyDc:                "c3",
			userType:               models.MonitorUser,
			fromSlave:              true,
			localSlaveReadPriority: LocalSlaveReadPrefer,
			master:                 "c3-mysql-test00.bj:3319",
			masterStatus:           StatusUp,
			slaveAddrs:             []string{"c3-mysql-test00.bj:3329", "c3-mysql-test00.bj:3339"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp},
			weights:                []int{1, 1},
			getCounts:              6,
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3329": 3,
				"c3-mysql-test00.bj:3339": 3,
			},
			expectErr: 0,
		},
		{
			name:                   "MonitorUser, fromSlave true, slave down, fallback to master",
			proxyDc:                "c3",
			userType:               models.MonitorUser,
			fromSlave:              true,
			localSlaveReadPriority: LocalSlaveReadForce,
			master:                 "c3-mysql-test00.bj:3319",
			masterStatus:           StatusUp,
			slaveAddrs:             []string{"c3-mysql-test00.bj:3329", "c3-mysql-test00.bj:3339"},
			weights:                []int{1, 1},
			slaveStatus:            []StatusCode{StatusDown, StatusUp},
			getCounts:              6,
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3339": 6, // All connections should go to the available slave
			},
			expectErr: 0,
		},
		{
			name:                   "MonitorUser, fromSlave true, all slaves down, master up, should return error",
			proxyDc:                "c3",
			userType:               models.MonitorUser,
			fromSlave:              true,
			localSlaveReadPriority: LocalSlaveReadForce,
			master:                 "c3-mysql-test00.bj:3319",
			masterStatus:           StatusUp,
			slaveAddrs:             []string{"c3-mysql-test00.bj:3329", "c3-mysql-test00.bj:3339"},
			weights:                []int{1, 1},
			slaveStatus:            []StatusCode{StatusDown, StatusDown},
			getCounts:              6,
			expectAddrs: map[string]int{
				"c3-mysql-test00.bj:3319": 6,
			}, // No connection should be made
			expectErr: 0, // Expecting 6 errors
		},
		{
			name:                   "MonitorUser, fromSlave true, all slaves down, master down, should return error",
			proxyDc:                "c3",
			userType:               models.MonitorUser,
			fromSlave:              true,
			localSlaveReadPriority: LocalSlaveReadForce,
			master:                 "c3-mysql-test00.bj:3319",
			masterStatus:           StatusDown,
			slaveAddrs:             []string{"c3-mysql-test00.bj:3329", "c3-mysql-test00.bj:3339"},
			weights:                []int{1, 1},
			slaveStatus:            []StatusCode{StatusDown, StatusDown},
			getCounts:              6,
			expectAddrs:            nil, // No connection should be made
			expectErr:              6,   // Expecting 6 errors
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			// 生成数据库连接池
			master, err := generateDBInfoWithWeights(mockCtl, []string{tc.master}, []StatusCode{tc.masterStatus}, []int{1})
			assert.Nil(t, err)
			master.InitBalancers(tc.proxyDc)

			// 生成数据库连接池
			slave, err := generateDBInfoWithWeights(mockCtl, tc.slaveAddrs, tc.slaveStatus, tc.weights)
			slave.InitBalancers(tc.proxyDc)
			assert.Nil(t, err)

			s := &Slice{
				MonitorMaster: master,
				MonitorSlave:  slave,
			}
			s.ProxyDatacenter = tc.proxyDc
			var actualError int
			var actualAddrs = make(map[string]int)
			// Call GetConn the expected number of times and track the connection counts
			for j := 0; j < tc.getCounts; j++ {
				reqCtx := util.NewRequestContext()
				reqCtx.SetFromSlave(tc.fromSlave)
				pc, err := s.GetConn(reqCtx, tc.userType, tc.localSlaveReadPriority)
				if err != nil {
					actualError++
				} else {
					// Track the address of the connection
					actualAddrs[pc.GetAddr()]++
				}
			}

			// Check if the addresses match the expected values
			for expectedAddr, expectedCount := range tc.expectAddrs {
				if actualAddrs[expectedAddr] != expectedCount {
					t.Errorf("for case %s, expected %d calls to %s, but got %d", tc.name, expectedCount, expectedAddr, actualAddrs[expectedAddr])
				}
			}

			// Verify the number of errors
			if actualError != tc.expectErr {
				t.Errorf("for case %s, expected %d errors, but got %d", tc.name, tc.expectErr, actualError)
			}
		})
	}
}

func generateDBInfoWithWeights(mockCtl *gomock.Controller, slaveHosts []string, slaveStatus []StatusCode, slaveWeights []int) (*DBInfo, error) {
	if len(slaveHosts) != len(slaveStatus) || len(slaveHosts) != len(slaveWeights) {
		return nil, fmt.Errorf("mismatched lengths: hosts=%d, status=%d, weights=%d", len(slaveHosts), len(slaveStatus), len(slaveWeights))
	}

	nodes := make([]*NodeInfo, len(slaveHosts)) // 直接存储 `NodeInfo`

	for i, host := range slaveHosts {
		dc, _ := util.GetInstanceDatacenter(host)
		pc := NewMockPooledConnect(mockCtl)
		pc.EXPECT().GetAddr().Return(host).AnyTimes()

		mcp := NewMockConnectionPool(mockCtl)
		mcp.EXPECT().Datacenter().Return(dc).AnyTimes()
		mcp.EXPECT().Get(context.TODO()).Return(pc, nil).AnyTimes()
		mcp.EXPECT().Addr().Return(host).AnyTimes()

		// 创建 `NodeInfo`，将 `ConnectionPool` 存入其中
		nodes[i] = &NodeInfo{
			Address:    host,
			Datacenter: dc,
			Weight:     slaveWeights[i],
			ConnPool:   mcp,
			Status:     slaveStatus[i],
		}
	}

	// 直接返回 `DBInfo`，删除 `StatusMap` 和 `ConsecutiveErrors`
	return &DBInfo{Nodes: nodes}, nil
}

func generateDBInfoDefault(mockCtl *gomock.Controller, slaveHosts []string, slaveStatus []StatusCode) (*DBInfo, error) {
	if len(slaveHosts) != len(slaveStatus) {
		return nil, fmt.Errorf("mismatched lengths: hosts=%d, status=%d", len(slaveHosts), len(slaveStatus))
	}

	nodes := make([]*NodeInfo, len(slaveHosts)) // 直接存储 `NodeInfo`

	for i, host := range slaveHosts {
		dc, _ := util.GetInstanceDatacenter(host)

		pc := NewMockPooledConnect(mockCtl)
		pc.EXPECT().GetAddr().Return(host).AnyTimes()
		pc.EXPECT().PingWithTimeout(gomock.Any()).Return(nil).AnyTimes()
		pc.EXPECT().Begin().Return(nil).AnyTimes()
		pc.EXPECT().Commit().Return(nil).AnyTimes()
		pc.EXPECT().Rollback().Return(nil).AnyTimes()
		pc.EXPECT().IsClosed().Return(false).AnyTimes()
		pc.EXPECT().Reconnect().Return(nil).AnyTimes()
		pc.EXPECT().UseDB(gomock.Any()).Return(nil).AnyTimes()

		mcp := NewMockConnectionPool(mockCtl)
		mcp.EXPECT().Datacenter().Return(dc).AnyTimes()
		mcp.EXPECT().Get(context.TODO()).Return(pc, nil).AnyTimes()
		mcp.EXPECT().Addr().Return(host).AnyTimes()
		mcp.EXPECT().GetCheck(context.Background()).Return(pc, nil).AnyTimes()
		mcp.EXPECT().SetLastChecked().Return().AnyTimes()

		// 创建 `NodeInfo`，默认 `Weight=1`
		nodes[i] = &NodeInfo{
			Address:    host,
			Datacenter: dc,
			Weight:     1,
			ConnPool:   mcp,
			Status:     slaveStatus[i],
		}
	}

	// 直接返回 `DBInfo`
	return &DBInfo{Nodes: nodes}, nil
}

func TestCheckSlaveSyncStatus(t *testing.T) {
	testCases := []struct {
		name                string
		c                   SlaveStatus
		secondsBehindMaster int
		statusAlive         bool
	}{
		{
			"test SecondsBehindMaster more than ns config",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Yes"},
			20,
			false,
		},
		{
			"test SecondsBehindMaster less than ns config",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Yes"},
			2000,
			true,
		},
		{
			"test SlaveIORunning not running",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Connecting", SlaveSQLRunning: "Yes"},
			2000,
			false,
		},
		{
			"test SlaveSQLRunning not running",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Connecting"},
			2000,
			false,
		},
	}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	for _, ca := range testCases {
		t.Run(ca.name, func(t *testing.T) {
			slice0SlaveConn := NewMockPooledConnect(mockCtl)
			slice0SlaveConn.EXPECT().GetAddr().Return("127.0.0.1:13307").AnyTimes()
			slice0SlaveConn.EXPECT().Execute("show slave status;", 0).Return(&mysql.Result{
				Status: 2,
				Resultset: &mysql.Resultset{
					Fields: []*mysql.Field{
						{Name: []byte("Seconds_Behind_Master")},
						{Name: []byte("Slave_IO_Running")},
						{Name: []byte("Slave_SQL_Running")},
						{Name: []byte("Master_Log_File")},
						{Name: []byte("Read_Master_Log_Pos")},
						{Name: []byte("Relay_Master_Log_File")},
						{Name: []byte("Exec_Master_Log_Pos")},
					},
					FieldNames: map[string]int{
						"Seconds_Behind_Master": 0,
						"Slave_IO_Running":      1,
						"Slave_SQL_Running":     2,
						"Master_Log_File":       3,
						"Read_Master_Log_Pos":   4,
						"Relay_Master_Log_File": 5,
						"Exec_Master_Log_Pos":   6,
					},
					Values: [][]interface{}{
						{ca.c.SecondsBehindMaster, ca.c.SlaveIORunning, ca.c.SlaveSQLRunning, ca.c.MasterLogFile, ca.c.ReadMasterLogPos, ca.c.RelayMasterLogFile, ca.c.ExecMasterLogPos},
					},
				},
			}, nil)

			isBehind, _ := checkSlaveSyncStatus(slice0SlaveConn, ca.secondsBehindMaster)
			assert.Equal(t, isBehind, ca.statusAlive)
		})
	}
}

func TestSlice_Close(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// 初始化 MockConnectionPool
	mockMasterPool := NewMockConnectionPool(mockCtl)
	mockSlavePool := NewMockConnectionPool(mockCtl)
	mockStatisticSlavePool := NewMockConnectionPool(mockCtl)

	// 设置每个连接池的 Close 方法应当被调用一次
	mockMasterPool.EXPECT().Close().Times(1)
	mockSlavePool.EXPECT().Close().Times(1)
	mockStatisticSlavePool.EXPECT().Close().Times(1)

	// 创建 Slice 实例，并使用 mock 连接池
	slice := &Slice{
		Master: &DBInfo{
			Nodes: []*NodeInfo{
				{
					Address:  "master_host",
					ConnPool: mockMasterPool,
				},
			},
		},
		Slave: &DBInfo{
			Nodes: []*NodeInfo{
				{
					Address:  "slave_host",
					ConnPool: mockSlavePool,
				},
			},
		},
		StatisticSlave: &DBInfo{
			Nodes: []*NodeInfo{
				{
					Address:  "stat_slave_host",
					ConnPool: mockStatisticSlavePool,
				},
			},
		},
	}

	// 调用 Close 方法
	slice.Close()
}

func TestGetIndicesAndWeights(t *testing.T) {
	tests := []struct {
		name            string
		dbInfo          *DBInfo
		proxyDatacenter string
		expectedLocal   IndexWeightList
		expectedRemote  IndexWeightList
		expectedGlobal  IndexWeightList
	}{
		{
			name: "Valid indices",
			dbInfo: &DBInfo{
				Nodes: []*NodeInfo{
					{Weight: 4, Datacenter: "c1"},
					{Weight: 3, Datacenter: "c1"},
					{Weight: 2, Datacenter: "c2"},
					{Weight: 1, Datacenter: "c2"},
				},
			},
			proxyDatacenter: "c1",
			expectedLocal: IndexWeightList{
				Indices: []int{0, 1},
				Weights: []int{4, 3},
			},
			expectedRemote: IndexWeightList{
				Indices: []int{2, 3},
				Weights: []int{2, 1},
			},
			expectedGlobal: IndexWeightList{
				Indices: []int{0, 1, 2, 3},
				Weights: []int{4, 3, 2, 1},
			},
		},
		{
			name: "All nodes have weight 0",
			dbInfo: &DBInfo{
				Nodes: []*NodeInfo{
					{Weight: 0, Datacenter: "c1"},
					{Weight: 0, Datacenter: "c1"},
					{Weight: 0, Datacenter: "c2"},
					{Weight: 0, Datacenter: "c2"},
				},
			},
			proxyDatacenter: "c1",
			expectedLocal:   IndexWeightList{Indices: []int{}, Weights: []int{}},
			expectedRemote:  IndexWeightList{Indices: []int{}, Weights: []int{}},
			expectedGlobal:  IndexWeightList{Indices: []int{}, Weights: []int{}},
		},
		{
			name: "Some nodes have weight 0",
			dbInfo: &DBInfo{
				Nodes: []*NodeInfo{
					{Weight: 4, Datacenter: "c1"},
					{Weight: 0, Datacenter: "c1"},
					{Weight: 2, Datacenter: "c2"},
					{Weight: 0, Datacenter: "c2"},
				},
			},
			proxyDatacenter: "c1",
			expectedLocal: IndexWeightList{
				Indices: []int{0},
				Weights: []int{4},
			},
			expectedRemote: IndexWeightList{
				Indices: []int{2},
				Weights: []int{2},
			},
			expectedGlobal: IndexWeightList{
				Indices: []int{0, 2},
				Weights: []int{4, 2},
			},
		},
		{
			name: "Some nodes have weight -1",
			dbInfo: &DBInfo{
				Nodes: []*NodeInfo{
					{Weight: 4, Datacenter: "c1"},
					{Weight: -1, Datacenter: "c1"},
					{Weight: 2, Datacenter: "c2"},
					{Weight: 0, Datacenter: "c2"},
				},
			},
			proxyDatacenter: "c1",
			expectedLocal: IndexWeightList{
				Indices: []int{0},
				Weights: []int{4},
			},
			expectedRemote: IndexWeightList{
				Indices: []int{2},
				Weights: []int{2},
			},
			expectedGlobal: IndexWeightList{
				Indices: []int{0, 2},
				Weights: []int{4, 2},
			},
		},
		{
			name: "Some nodes have weight -1",
			dbInfo: &DBInfo{
				Nodes: []*NodeInfo{
					{Weight: 5, Datacenter: "c1"},
					{Weight: -1, Datacenter: "c1"},
					{Weight: 0, Datacenter: "c2"},
					{Weight: 1, Datacenter: "c2"},
				},
			},
			proxyDatacenter: "c1",
			expectedLocal: IndexWeightList{
				Indices: []int{0},
				Weights: []int{5},
			},
			expectedRemote: IndexWeightList{
				Indices: []int{3},
				Weights: []int{1},
			},
			expectedGlobal: IndexWeightList{
				Indices: []int{0, 3},
				Weights: []int{5, 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			local, remote, global := tt.dbInfo.getIndicesAndWeights(tt.proxyDatacenter)

			// 确保只比较有效节点（去除权重为0的情况）
			filteredExpectedLocal := filterZeroWeight(tt.expectedLocal)
			filteredExpectedRemote := filterZeroWeight(tt.expectedRemote)
			filteredExpectedGlobal := filterZeroWeight(tt.expectedGlobal)

			assert.Equal(t, filteredExpectedLocal, *local, "Local indices/weights mismatch")
			assert.Equal(t, filteredExpectedRemote, *remote, "Remote indices/weights mismatch")
			assert.Equal(t, filteredExpectedGlobal, *global, "Global indices/weights mismatch")
		})
	}
}

// 过滤掉权重为 0 的节点
func filterZeroWeight(iwl IndexWeightList) IndexWeightList {
	var filtered IndexWeightList
	for i, weight := range iwl.Weights {
		if weight > 0 {
			filtered.Indices = append(filtered.Indices, iwl.Indices[i])
			filtered.Weights = append(filtered.Weights, weight)
		}
	}
	return filtered
}

// 测试所有检查通过，从库应该保持 StatusUp。
func TestCheckBackendSlaveStatus_All_Passed(t *testing.T) {
	mockey.PatchConvey("Test checkBackendSlaveStatus - All Checks Passed", t, func() {
		// 创建一个可控的 ticker
		mockTicker := time.NewTicker(100 * time.Millisecond)
		defer mockTicker.Stop()

		// 创建一个不为nil的连接
		mockConn := &pooledConnectImpl{}
		// Mock time.NewTicker，使其返回 mockTicker
		mockey.Mock(time.NewTicker).
			Return(mockTicker).
			Build()

		mockCtl := gomock.NewController(t)
		defer mockCtl.Finish()

		slave, err := generateDBInfoDefault(mockCtl, []string{"127.0.0.1:3329"}, []StatusCode{StatusDown})
		assert.Nil(t, err)

		mockSlice := &Slice{
			Namespace:      "test_ns",
			Slave:          slave,
			StatisticSlave: slave,
		}

		// Mock `CheckConnectionPool` 从库连接池检查成功
		mockey.Mock((*NodeInfo).GetPooledConnectWithHealthCheck).
			Return(mockConn, nil).
			Build()

		// Mock `CheckDownAfterNoAlive` 从库没有超过下线阈值
		mockey.Mock((*NodeInfo).ShouldDownAfterNoAlive).
			Return(false, int64(10)).
			Build()

		// Mock `GetMasterStatus` 主库没有宕机
		mockey.Mock((*Slice).GetMasterStatus).
			Return(StatusUp, nil).
			Build()

		// Mock `checkSlaveSyncStatus` 从库同步正常
		mockey.Mock(checkSlaveSyncStatus).
			Return(true, nil).
			Build()

		// 执行 checkBackendSlaveStatus
		ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
		defer cancel()

		mockSlice.checkBackendSlaveStatus(ctx, slave, 10, 30)

		// Expect slave status to remain up even when master is down
		for _, node := range slave.Nodes {
			assert.Equal(t, StatusUp, node.GetStatus())
		}
	})
}

func TestCheckBackendSlaveStatus_CoonPool_Check_Error(t *testing.T) {
	mockey.PatchConvey("Test checkBackendSlaveStatus - Slave ConnPool connect is nil, Start is Down, Not Exceeding the offline threshold", t, func() {

		// 创建一个可控的 ticker
		mockTicker := time.NewTicker(100 * time.Millisecond)
		defer mockTicker.Stop()

		// 创建一个为nil的连接池连接
		var mockNilConn PooledConnect
		// Mock time.NewTicker，使其返回 mockTicker
		mockey.Mock(time.NewTicker).
			Return(mockTicker).
			Build()

		mockCtl := gomock.NewController(t)
		defer mockCtl.Finish()

		slave, err := generateDBInfoDefault(mockCtl, []string{"127.0.0.1:3329"}, []StatusCode{StatusDown})
		assert.Nil(t, err)

		mockSlice := &Slice{
			Namespace:      "test_ns",
			Slave:          slave,
			StatisticSlave: slave,
		}

		// Mock `CheckConnectionPool` 从库连接池检查失败
		mockey.Mock((*NodeInfo).GetPooledConnectWithHealthCheck).
			Return(mockNilConn, nil).
			Build()

		// Mock `CheckDownAfterNoAlive` 从库没有超过下线阈值
		mockey.Mock((*NodeInfo).ShouldDownAfterNoAlive).
			Return(false, int64(10)).
			Build()

		// Mock `GetMasterStatus` 主库没有宕机
		mockey.Mock((*Slice).GetMasterStatus).
			Return(StatusUp, nil).
			Build()

		// Mock `checkSlaveSyncStatus` 从库同步正常
		mockey.Mock(checkSlaveSyncStatus).
			Return(true, fmt.Errorf("check slave sync status error: pc is nil")).
			Build()

		// 执行 checkBackendSlaveStatus
		ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
		defer cancel()

		mockSlice.checkBackendSlaveStatus(ctx, slave, 10, 30)

		// Expect slave status to remain up even when master is down
		for _, node := range slave.Nodes {
			assert.Equal(t, StatusDown, node.GetStatus())
		}
	})
}

func TestCheckBackendSlaveStatus_CoonPool_Check_Nil(t *testing.T) {
	mockey.PatchConvey("Test checkBackendSlaveStatus - Slave ConnPool connect is nil, Start is UP, Exceeding the offline threshold", t, func() {

		// 创建一个可控的 ticker
		mockTicker := time.NewTicker(100 * time.Millisecond)
		defer mockTicker.Stop()

		// 创建一个为nil的连接池连接
		var mockNilConn PooledConnect
		// Mock time.NewTicker，使其返回 mockTicker
		mockey.Mock(time.NewTicker).
			Return(mockTicker).
			Build()

		mockCtl := gomock.NewController(t)
		defer mockCtl.Finish()

		slave, err := generateDBInfoDefault(mockCtl, []string{"127.0.0.1:3329"}, []StatusCode{StatusUp})
		assert.Nil(t, err)

		mockSlice := &Slice{
			Namespace:      "test_ns",
			Slave:          slave,
			StatisticSlave: slave,
		}

		// Mock `CheckConnectionPool` 从库连接池检查失败
		mockey.Mock((*NodeInfo).GetPooledConnectWithHealthCheck).
			Return(mockNilConn, nil).
			Build()

		// Mock `CheckDownAfterNoAlive` 从库没有超过下线阈值
		mockey.Mock((*NodeInfo).ShouldDownAfterNoAlive).
			Return(true, int64(10)).
			Build()

		// Mock `GetMasterStatus` 主库没有宕机
		mockey.Mock((*Slice).GetMasterStatus).
			Return(StatusUp, nil).
			Build()

		// Mock `checkSlaveSyncStatus` 从库同步正常
		mockey.Mock(checkSlaveSyncStatus).
			Return(true, fmt.Errorf("check slave sync status error: pc is nil")).
			Build()

		// 执行 checkBackendSlaveStatus
		ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
		defer cancel()

		mockSlice.checkBackendSlaveStatus(ctx, slave, 1, 30)

		// Expect slave status to remain up even when master is down
		for _, node := range slave.Nodes {
			assert.Equal(t, StatusDown, node.GetStatus())
		}
	})
}

// 主库下线，从库跳过同步
func TestCheckBackendSlaveStatus_MasterDown(t *testing.T) {
	mockey.PatchConvey("Test checkBackendSlaveStatus - Master down", t, func() {
		// 创建一个可控的 ticker
		mockTicker := time.NewTicker(100 * time.Millisecond)
		defer mockTicker.Stop()

		// 创建一个不为nil的连接
		mockConn := &pooledConnectImpl{}

		// Mock time.NewTicker，使其返回 mockTicker
		mockey.Mock(time.NewTicker).
			Return(mockTicker).
			Build()

		mockCtl := gomock.NewController(t)
		defer mockCtl.Finish()

		slave, err := generateDBInfoDefault(mockCtl, []string{"127.0.0.1:3329"}, []StatusCode{StatusDown})
		assert.Nil(t, err)

		mockSlice := &Slice{
			Namespace:      "test_ns",
			Slave:          slave,
			StatisticSlave: slave,
		}

		// Mock `CheckConnectionPool` 从库连接池检查成功
		mockey.Mock((*NodeInfo).GetPooledConnectWithHealthCheck).
			Return(mockConn, nil).
			Build()

		// Mock `CheckDownAfterNoAlive` 从库没有超过下线阈值
		mockey.Mock((*NodeInfo).ShouldDownAfterNoAlive).
			Return(false, int64(10)).
			Build()

		// Mock `GetMasterStatus` 主库下线，但没有错误
		mockey.Mock((*Slice).GetMasterStatus).
			Return(StatusDown, nil).
			Build()

		// Mock `checkSlaveSyncStatus` 从库同步没有超过阈值
		mockey.Mock(checkSlaveSyncStatus).
			Return(true, nil).
			Build()

		// 执行 checkBackendSlaveStatus
		ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
		defer cancel()

		mockSlice.checkBackendSlaveStatus(ctx, slave, 10, 30)

		// Expect slave status to remain up even when master is down
		for _, node := range slave.Nodes {
			assert.Equal(t, StatusUp, node.GetStatus())
		}
	})
}

// 主库宕机，从库跳过同步
func TestCheckBackendSlaveStatus_MasterError(t *testing.T) {
	mockey.PatchConvey("Test checkBackendSlaveStatus - Master error", t, func() {
		// 创建一个可控的 ticker
		mockTicker := time.NewTicker(100 * time.Millisecond)
		defer mockTicker.Stop()

		// 创建一个不为nil的连接
		mockConn := &pooledConnectImpl{}

		// Mock time.NewTicker，使其返回 mockTicker
		mockey.Mock(time.NewTicker).
			Return(mockTicker).
			Build()

		mockCtl := gomock.NewController(t)
		defer mockCtl.Finish()

		slave, err := generateDBInfoDefault(mockCtl, []string{"127.0.0.1:3329"}, []StatusCode{StatusDown})
		assert.Nil(t, err)

		mockSlice := &Slice{
			Namespace:      "test_ns",
			Slave:          slave,
			StatisticSlave: slave,
		}

		// Mock `CheckConnectionPool` 从库连接池检查成功
		mockey.Mock((*NodeInfo).GetPooledConnectWithHealthCheck).
			Return(mockConn, nil).
			Build()

		// Mock `CheckDownAfterNoAlive` 从库没有超过下线阈值
		mockey.Mock((*NodeInfo).ShouldDownAfterNoAlive).
			Return(false, int64(10)).
			Build()

		// Mock `GetMasterStatus` 主库宕机，返回错误
		mockey.Mock((*Slice).GetMasterStatus).
			Return(StatusDown, fmt.Errorf("get master status error")).
			Build()

		// 执行 checkBackendSlaveStatus
		ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
		defer cancel()

		mockSlice.checkBackendSlaveStatus(ctx, slave, 10, 30)

		// Expect slave status to remain up even when master is down
		for _, node := range slave.Nodes {
			assert.Equal(t, StatusUp, node.GetStatus())
		}
	})
}

// 主库正常，从库同步延迟
func TestCheckBackendSlaveStatus_SlaveSyncDelay(t *testing.T) {
	mockey.PatchConvey("Test checkBackendSlaveStatus - Slave Sync Delay", t, func() {
		// 创建一个可控的 ticker
		mockTicker := time.NewTicker(100 * time.Millisecond)
		defer mockTicker.Stop()

		// 创建一个不为nil的连接
		mockConn := &pooledConnectImpl{}

		// Mock time.NewTicker，使其返回 mockTicker
		mockey.Mock(time.NewTicker).
			Return(mockTicker).
			Build()

		mockCtl := gomock.NewController(t)
		defer mockCtl.Finish()

		slave, err := generateDBInfoDefault(mockCtl, []string{"127.0.0.1:3329"}, []StatusCode{StatusUp})
		assert.Nil(t, err)

		mockSlice := &Slice{
			Namespace:      "test_ns",
			Slave:          slave,
			StatisticSlave: slave,
		}

		// Mock `CheckConnectionPool` 从库连接池检查成功
		mockey.Mock((*NodeInfo).GetPooledConnectWithHealthCheck).
			Return(mockConn, nil).
			Build()

		// Mock `CheckDownAfterNoAlive` 从库没有超过下线阈值
		mockey.Mock((*NodeInfo).ShouldDownAfterNoAlive).
			Return(false, int64(10)).
			Build()

		// Mock `GetMasterStatus` 主库没有宕机
		mockey.Mock((*Slice).GetMasterStatus).
			Return(StatusUp, nil).
			Build()

		// Mock `checkSlaveSyncStatus` 从库同步延迟超过阈值
		mockey.Mock(checkSlaveSyncStatus).
			Return(false, fmt.Errorf("sync delay exceeded")).
			Build()

		// 执行 checkBackendSlaveStatus
		ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
		defer cancel()

		mockSlice.checkBackendSlaveStatus(ctx, slave, 10, 30)

		// Expect slave status to be down due to sync delay
		for _, node := range slave.Nodes {
			assert.Equal(t, StatusDown, node.GetStatus())
		}
	})
}
