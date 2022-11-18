package server

import (
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/backend/mocks"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlaveIsLagBehand(t *testing.T) {
	//first we should
	type CaseAndResult struct {
		c         SlaveStatus
		namespace *Namespace
		r         bool
	}

	cases := []CaseAndResult{
		{
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Yes"},
			&Namespace{secondsBehindMaster: 20},
			true,
		},
		{
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Yes"},
			&Namespace{secondsBehindMaster: 2000},
			false,
		},
		{
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Connecting", SlaveSQLRunning: "Yes"},
			&Namespace{secondsBehindMaster: 2000},
			true,
		},
		{
			SlaveStatus{SecondsBehindMaster: 1000, SlaveIORunning: "Yes", SlaveSQLRunning: "Connecting"},
			&Namespace{secondsBehindMaster: 2000},
			true,
		},
	}

	for _, ca := range cases {
		convey.Convey("123", t, func() {
			p := gomonkey.NewPatches()
			defer p.Reset()
			p.ApplyFunc(GetSlaveStatus, func(conn backend.PooledConnect) (SlaveStatus, error) {
				return ca.c, nil
			})

			//use mock to test interface
			slice0MasterConn := new(mocks.PooledConnect)
			slice0MasterConn.On("GetAddr").Return(string("127.0.0.1:13307"))

			//use gomonkey to test struct
			isBehand, _ := slaveIsLagBehand(slice0MasterConn, ca.namespace)
			assert.Equal(t, isBehand, ca.r)
		})
	}
}
