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

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
			dbInfo, _ := s.ParseSlave(tt.slaveAdders)
			for i := range tt.expectAddrs {
				assert.Equal(t, tt.expectAddrs[i], dbInfo.ConnPool[i].Addr())
				assert.Equal(t, tt.expectDatacenters[i], dbInfo.Nodes[i].Datacenter)
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
			expectWeights:  []int{-5, 3}, // 默认权重
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
			expectWeights:  []int{2, -1, 0},
			expectingError: false,
		},

		// ✅ 测试混合非法字符权重测试
		{
			name:           "invalid weights",
			slaveAdders:    []string{"c3-mysql-test00.bj:3306@&", "c3-mysql-test01.bj:3308@+", "c4-mysql-test02.bj:3310@-"},
			expectAddrs:    []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"}, // 过滤掉权重为 0 和负数的
			expectIndices:  nil,
			expectWeights:  nil,
			expectingError: true,
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := new(Slice)
			dbInfo, err := s.ParseSlave(tt.slaveAdders)

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
			for i, node := range dbInfo.Nodes {
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

func generateDBInfoWithWeights(mockCtl *gomock.Controller, slaveHosts []string, slaveStatus []StatusCode, slaveWeights []int) (*DBInfo, error) {
	if len(slaveHosts) != len(slaveStatus) || len(slaveHosts) != len(slaveWeights) {
		return nil, fmt.Errorf("mismatched lengths: hosts=%d, status=%d, weights=%d", len(slaveHosts), len(slaveStatus), len(slaveWeights))
	}

	connPool := make([]ConnectionPool, len(slaveHosts))
	nodes := make([]NodeInfo, len(slaveHosts))
	statusMap := &sync.Map{}
	slaveConsecutiveErrors := &sync.Map{}

	for i, host := range slaveHosts {
		dc, _ := util.GetInstanceDatacenter(host)
		pc := NewMockPooledConnect(mockCtl)
		pc.EXPECT().GetAddr().Return(host).AnyTimes()

		mcp := NewMockConnectionPool(mockCtl)
		mcp.EXPECT().Datacenter().Return(dc).AnyTimes()
		mcp.EXPECT().Get(context.TODO()).Return(pc, nil).AnyTimes()
		mcp.EXPECT().Addr().Return(host).AnyTimes()

		connPool[i] = mcp
		statusMap.Store(i, slaveStatus[i])
		slaveConsecutiveErrors.Store(i, 0)

		// **新增 NodeInfo**
		nodes[i] = NodeInfo{
			Address:    host,
			Datacenter: dc,
			Weight:     slaveWeights[i],
		}
	}

	dbinfo := &DBInfo{
		ConnPool:          connPool,
		Nodes:             nodes, // **直接存储所有节点信息**
		StatusMap:         statusMap,
		ConsecutiveErrors: slaveConsecutiveErrors,
	}
	return dbinfo, nil
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
	//requirement := require.New(t)
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

	// 创建 Slice 实例并使用 mock 连接池
	slice := &Slice{
		Master: &DBInfo{
			ConnPool: []ConnectionPool{mockMasterPool},
		},
		Slave: &DBInfo{
			ConnPool: []ConnectionPool{mockSlavePool},
		},
		StatisticSlave: &DBInfo{
			ConnPool: []ConnectionPool{mockStatisticSlavePool},
		},
	}

	// 调用 Close 方法
	slice.Close()
}

func TestSlaveConsecutiveErrorCircuitBreaker(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	// 定义测试用例
	testCases := []struct {
		name                   string
		proxyDc                string
		localSlaveReadPriority int
		slaveAddrs             []string
		slaveStatus            []StatusCode
		errorThreshold         int
		operations             []string // "success" 或 "fail"
		expectedStatuses       []StatusCode
	}{
		// 连续错误未达到阈值，不熔断
		{
			name:                   "Continuous errors do not reach the threshold and the fuse does not blow",
			proxyDc:                "dc1",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			slaveAddrs:             []string{"slave1.dc1:3306", "slave2.dc1:3306"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp},
			errorThreshold:         3,
			operations:             []string{"fail", "fail", "success"},
			expectedStatuses:       []StatusCode{StatusUp, StatusUp},
		},
		// 连续错误达到阈值，熔断从库
		{
			name:                   "Continuous errors reach the threshold, fuse slave",
			proxyDc:                "dc1",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			slaveAddrs:             []string{"slave1.dc1:3306", "slave2.dc1:3306"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp},
			errorThreshold:         3,
			operations:             []string{"fail", "fail", "fail"},
			expectedStatuses:       []StatusCode{StatusDown, StatusUp},
		},
		// 错误后成功，错误计数重置，不熔断
		{
			name:                   "Success after error, error count reset, no fuse",
			proxyDc:                "dc1",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			slaveAddrs:             []string{"slave1.dc1:3306", "slave2.dc1:3306"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp},
			errorThreshold:         3,
			operations:             []string{"fail", "success", "fail", "fail", "fail"},
			expectedStatuses:       []StatusCode{StatusDown, StatusUp},
		},
		// 多个从库分别统计错误计数
		{
			name:                   "Multiple slaves count error counts separately",
			proxyDc:                "dc1",
			localSlaveReadPriority: LocalSlaveReadPrefer,
			slaveAddrs:             []string{"slave1.dc1:3306", "slave2.dc1:3306"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp},
			errorThreshold:         3,
			operations:             []string{"fail", "switch", "fail", "fail", "fail"},
			expectedStatuses:       []StatusCode{StatusUp, StatusDown},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初始化 Slice 和 DBInfo
			s := &Slice{
				ProxyDatacenter: tc.proxyDc,
				Slave:           generateDBInfoWithMockPools(mockCtl, tc.slaveAddrs, tc.slaveStatus),
			}

			s.Slave.ConsecutiveErrors = &sync.Map{}
			// 设置错误阈值
			s.MaxSlaveFuseErrorCount = tc.errorThreshold

			// 模拟操作
			currentSlaveIndex := 0
			for _, op := range tc.operations {
				var err error
				var pc PooledConnect

				if op == "fail" {
					// 模拟获取连接失败
					pc, err = s.getSlaveConnWithMockError(s.Slave, tc.localSlaveReadPriority, currentSlaveIndex)
					assert.NotNil(t, err)
				} else if op == "success" {
					// 模拟成功获取连接
					pc, err = s.getSlaveConnWithMockSuccess(s.Slave, tc.localSlaveReadPriority, currentSlaveIndex)
					assert.Nil(t, err)
					assert.NotNil(t, pc)
				} else if op == "switch" {
					// 切换到下一个从库
					currentSlaveIndex = (currentSlaveIndex + 1) % len(tc.slaveAddrs)
					continue
				}
			}

			// 检查最终的从库状态
			for idx, expectedStatus := range tc.expectedStatuses {
				status, _ := s.Slave.GetStatus(idx)
				assert.Equal(t, expectedStatus, status, fmt.Sprintf("Slave %d status should be %v", idx, expectedStatus))
			}
		})
	}
}

// 模拟获取从库连接失败，增加错误计数
func (s *Slice) getSlaveConnWithMockError(slavesInfo *DBInfo, localSlaveReadPriority int, index int) (PooledConnect, error) {
	// 增加错误计数
	slavesInfo.IncrementErrorCount(index)
	// 检查是否需要熔断
	if slavesInfo.GetErrorCount(index) >= s.MaxSlaveFuseErrorCount {
		slavesInfo.SetStatus(index, StatusDown)
	}
	return nil, fmt.Errorf("mock connection error")
}

// 模拟成功获取从库连接，重置错误计数
func (s *Slice) getSlaveConnWithMockSuccess(slavesInfo *DBInfo, localSlaveReadPriority int, index int) (PooledConnect, error) {
	// 重置错误计数
	slavesInfo.ResetErrorCount(index)
	// 返回模拟的连接
	pc := NewMockPooledConnect(gomock.NewController(nil))
	pc.EXPECT().GetAddr().Return(slavesInfo.ConnPool[index].Addr()).AnyTimes()
	return pc, nil
}

func generateDBInfoWithMockPools(mockCtl *gomock.Controller, slaveAddrs []string, slaveStatus []StatusCode) *DBInfo {
	connPool := make([]ConnectionPool, len(slaveAddrs))
	slaveWeights := make([]int, len(slaveAddrs))
	datacenter := make([]string, len(slaveAddrs))
	statusMap := &sync.Map{}
	slaveConsecutiveErrors := &sync.Map{}

	for i, addr := range slaveAddrs {
		dc, _ := util.GetInstanceDatacenter(addr)

		mcp := NewMockConnectionPool(mockCtl)
		mcp.EXPECT().Datacenter().Return(dc).AnyTimes()
		mcp.EXPECT().Addr().Return(addr).AnyTimes()

		connPool[i] = mcp
		datacenter[i] = dc
		slaveWeights[i] = 1
		statusMap.Store(i, slaveStatus[i])
		slaveConsecutiveErrors.Store(i, 0)
	}
	return &DBInfo{
		ConnPool:          connPool,
		StatusMap:         statusMap,
		ConsecutiveErrors: slaveConsecutiveErrors,
	}
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
				Nodes: []NodeInfo{
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
				Nodes: []NodeInfo{
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
				Nodes: []NodeInfo{
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
				Nodes: []NodeInfo{
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
				Nodes: []NodeInfo{
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
