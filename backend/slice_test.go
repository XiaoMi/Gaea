package backend

import (
	"context"
	"fmt"
	"github.com/XiaoMi/Gaea/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
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
			for i, _ := range tt.expectAddrs {
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
			slaveStatus:            []StatusCode{UP, UP, UP},
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
			slaveStatus:            []StatusCode{UP, UP, UP},
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
			slaveStatus:            []StatusCode{UP, UP, UP},
			expectAddrs:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadPreferred local one down",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{DOWN, UP, UP},
			expectAddrs:            []string{"c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadPreferred local all down",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{DOWN, DOWN, UP},
			expectAddrs:            []string{"c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"},
		},
		{
			name:                   "test LocalSlaveReadPreferred local all down 2 use last another slave",
			proxyDc:                "c4-bj",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{UP, UP, DOWN},
			expectAddrs:            []string{"c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadPreferred no local slave use last another slave",
			proxyDc:                "c5",
			localSlaveReadPriority: LocalSlaveReadPreferred,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{UP, UP, UP},
			expectAddrs:            []string{"c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"},
		},
		{
			name:                   "test LocalSlaveReadForce all up",
			proxyDc:                "c3",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{UP, UP, UP},
			expectAddrs:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308"},
		},
		{
			name:                   "test LocalSlaveReadForce local all down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{UP, UP, DOWN},
			expectAddrs:            []string{},
		},
		{
			name:                   "test LocalSlaveReadForce all down",
			proxyDc:                "c4",
			localSlaveReadPriority: LocalSlaveReadForce,
			getCounts:              4,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{DOWN, DOWN, DOWN},
			expectAddrs:            []string{},
		},
		{
			name:                   "test LocalSlaveReadForce no local slave",
			proxyDc:                "c5",
			getCounts:              10,
			localSlaveReadPriority: LocalSlaveReadForce,
			slaveAdders:            []string{"c3-mysql-test00.bj:3306", "c3-mysql-test01.bj:3308", "c4-mysql-test02.bj:3310"},
			slaveStatus:            []StatusCode{UP, UP, UP},
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
	StatusMap := sync.Map{}
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
	}
	slaveBalancer := newBalancer(slaveWeights, len(connPool))

	return &DBInfo{connPool, slaveBalancer, StatusMap, datacenter}
}
