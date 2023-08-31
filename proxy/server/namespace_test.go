package server

import (
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlaveIsLagBehand(t *testing.T) {
	testCases := []struct {
		name      string
		c         SlaveStatus
		namespace *Namespace
		r         bool
	}{
		{
			"test SecondsBehindMaster more than ns config",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Yes"},
			&Namespace{secondsBehindMaster: 20},
			true,
		},
		{
			"test SecondsBehindMaster less than ns config",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Yes"},
			&Namespace{secondsBehindMaster: 2000},
			false,
		},
		{
			"test SlaveIORunning not running",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Connecting", SlaveSQLRunning: "Yes"},
			&Namespace{secondsBehindMaster: 2000},
			true,
		},
		{
			"test SlaveSQLRunning not running",
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Connecting"},
			&Namespace{secondsBehindMaster: 2000},
			true,
		},
	}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	for _, ca := range testCases {
		t.Run(ca.name, func(t *testing.T) {
			slice0MasterConn := backend.NewMockPooledConnect(mockCtl)
			slice0MasterConn.EXPECT().GetAddr().Return("127.0.0.1:13307").AnyTimes()
			slice0MasterConn.EXPECT().Execute("show slave status;", 0).Return(&mysql.Result{
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

			isBehind, _ := slaveIsLagBehind(slice0MasterConn, ca.namespace)
			assert.Equal(t, isBehind, ca.r)
		})
	}
}
