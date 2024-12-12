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
				assert.Equal(t, tt.expectDatacenters[i], dbInfo.Datacenter[i])
			}
		})
	}
}

func TestGetSlaveConn(t *testing.T) {
	testCases := []struct {
		name                   string
		localSlaveReadPriority int
		proxyDc                string
		slaveAdders            []string
		slaveStatus            []StatusCode
		getCounts              int
		expectAddrs            []string
	}{
		{
			name:                   "test LocalSlaveReadClosed all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              8,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			expectAddrs: []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310",
				"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310",
				"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadClosed no local slave",
			proxyDc:                "c5",
			localSlaveReadPriority: LocalSlaveReadClosed,
			getCounts:              8,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			expectAddrs: []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310",
				"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310",
				"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadPreferred all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			expectAddrs:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadPreferred local one down",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusDown, StatusUp, StatusUp},
			expectAddrs:            []string{"c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadPreferred local all down",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusDown, StatusDown, StatusUp},
			expectAddrs:            []string{"c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"},
		},
		{
			name:                   "test LocalSlaveReadPreferred local all down 2 use last another slave",
			proxyDc:                "c4-bj",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown},
			expectAddrs:            []string{"c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadPreferred no local slave use last another slave",
			proxyDc:                "c5",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			expectAddrs:            []string{"c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"},
		},
		{
			name:                   "test LocalSlaveReadForce all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			expectAddrs:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadForce local all down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusDown},
			expectAddrs:            []string{},
		},
		{
			name:                   "test LocalSlaveReadForce all down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusDown, StatusDown, StatusDown},
			expectAddrs:            []string{},
		},
		{
			name:                   "test LocalSlaveReadForce no local slave",
			proxyDc:                "c5",
			getCounts:              10,
			localSlaveReadPriority: LocalSlaveReadForce,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{StatusUp, StatusUp, StatusUp},
			expectAddrs:            []string{},
		},
	}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			dbInfo := generateDBInfo(mockCtl, tt.slaveAdders, tt.slaveStatus)
			s := &Slice{Slave: dbInfo}
			s.ProxyDatacenter = tt.proxyDc
			for j := 0; j < tt.getCounts; j++ {
				cp, err := s.GetSlaveConn(dbInfo, tt.localSlaveReadPriority)

				if len(tt.expectAddrs) == 0 {
					assert.NotNil(t, err)
					continue
				}
				assert.Nil(t, err)
				assert.Equal(t, tt.expectAddrs[j], cp.GetAddr(), fmt.Sprintf(tt.name))
			}
		})
	}
}

func generateDBInfo(mockCtl *gomock.Controller, slaveHosts []string, slaveStatus []StatusCode) *DBInfo {
	connPool := make([]ConnectionPool, 0, len(slaveHosts))
	slaveWeights := make([]int, 0, len(slaveHosts))
	datacenter := make([]string, 0, len(slaveHosts))
	StatusMap := &sync.Map{}
	SlaveConsecutiveErrors := &sync.Map{}
	for i, host := range slaveHosts {
		dc, _ := util.GetInstanceDatacenter(host)
		pc := NewMockPooledConnect(mockCtl)
		pc.EXPECT().GetAddr().Return(host).AnyTimes()

		mcp := NewMockConnectionPool(mockCtl)
		mcp.EXPECT().Datacenter().Return(dc).AnyTimes()
		mcp.EXPECT().Get(context.TODO()).Return(pc, nil).AnyTimes()
		mcp.EXPECT().Addr().Return(host).AnyTimes()

		connPool = append(connPool, mcp)
		datacenter = append(datacenter, dc)
		slaveWeights = append(slaveWeights, 1)
		StatusMap.Store(i, slaveStatus[i])
		SlaveConsecutiveErrors.Store(i, 0)
	}
	slaveBalancer := newBalancer(slaveWeights, len(connPool))

	return &DBInfo{
		ConnPool:          connPool,
		Balancer:          slaveBalancer,
		StatusMap:         StatusMap,
		Datacenter:        datacenter,
		ConsecutiveErrors: SlaveConsecutiveErrors,
	}
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
			localSlaveReadPriority: LocalSlaveReadPreferred,
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
			localSlaveReadPriority: LocalSlaveReadPreferred,
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
			localSlaveReadPriority: LocalSlaveReadPreferred,
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
			localSlaveReadPriority: LocalSlaveReadPreferred,
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

	balancer := newBalancer(slaveWeights, len(connPool))

	return &DBInfo{
		ConnPool:          connPool,
		Balancer:          balancer,
		StatusMap:         statusMap,
		Datacenter:        datacenter,
		ConsecutiveErrors: slaveConsecutiveErrors,
	}
}
