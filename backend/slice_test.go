package backend

import (
	"context"
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
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
