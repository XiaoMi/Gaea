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
package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/golang/mock/gomock"

	"github.com/XiaoMi/Gaea/models"

	"fmt"

	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	mockmysql "github.com/siddontang/go-mysql/mysql"
)

func TestGetVariableExprResult(t *testing.T) {
	tests := []struct {
		variable []string
		expect   string
	}{
		{[]string{"ON", "on", "'on'", "`on`"}, "on"},
		{[]string{"OFF", "off", "'off'", "`off`"}, "off"},
		{[]string{"1", "'1'", "`1`"}, "1"},
		{[]string{"0", "'0'", "`0`"}, "0"},
	}
	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			for _, v := range test.variable {
				sql := fmt.Sprintf("set autocommit = %s", v)
				s, err := parser.ParseSQL(sql)
				if err != nil {
					t.Fatal(err)
				}
				stmt := s.(*ast.SetStmt)
				for _, v := range stmt.Variables {
					actual := getVariableExprResult(v.Value)
					if actual != test.expect {
						t.Errorf("not equal, expect: %v, actual: %v", test.expect, actual)
					}
				}
			}
		})
	}
}

/*

type PooledConnect interface {
	Recycle()
	Reconnect() error
	Close()
	IsClosed() bool
	UseDB(db string) error
	Execute(sql string) (*mysql.Result, error)
	SetAutoCommit(v uint8) error
	Begin() error
	Commit() error
	Rollback() error
	SetCharset(charset string, collation mysql.CollationID) (bool, error)
	FieldList(table string, wildcard string) ([]*mysql.Field, error)
	GetAddr() string
	SetSessionVariables(frontend *mysql.SessionVariables) (bool, error)
	WriteSetStatement() error
}
*/

func TestExecute(t *testing.T) {
	se, err := prepareSessionExecutor()
	if err != nil {
		t.Fatal("prepare session executer error:", err)
		return
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	slice0MasterPool := backend.NewMockConnectionPool(mockCtrl)
	slice0SlavePool := backend.NewMockConnectionPool(mockCtrl)
	slice1MasterPool := backend.NewMockConnectionPool(mockCtrl)
	slice1SlavePool := backend.NewMockConnectionPool(mockCtrl)
	se.manager.GetNamespace("test_executor_namespace").slices["slice-0"].Master = slice0MasterPool
	se.manager.GetNamespace("test_executor_namespace").slices["slice-0"].Slave = []backend.ConnectionPool{slice0SlavePool}
	se.manager.GetNamespace("test_executor_namespace").slices["slice-1"].Master = slice1MasterPool
	se.manager.GetNamespace("test_executor_namespace").slices["slice-1"].Slave = []backend.ConnectionPool{slice1SlavePool}

	slice0slaveConn := backend.NewMockPooledConnect(mockCtrl)
	//	slice1slaveConn := backend.NewMockPooledConnect(mockCtrl)
	//slice0slaveConn.EXPECT().Execute("SELECT * FROM `tbl_mycat` WHERE `k`=0").Return(nil, nil).Times(2)
	//slice1slaveConn.EXPECT().Execute("SELECT * FROM `tbl_mycat` WHERE `k`=0").Return(nil, nil).Times(2)

	slice0SlavePool.EXPECT().Get(context.TODO()).Return(slice0slaveConn, nil)
	//slice1SlavePool.EXPECT().Get(context.TODO()).Return(slice1slaveConn, nil)
	//slice0SlavePool.EXPECT().Put(slice0slaveConn).Return()
	//slice1SlavePool.EXPECT().Put(slice1slaveConn).Return()

	sqls := map[string]map[string][]string{
		"slice-0": {
			"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
			"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
		},
		"slice-1": {
			"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
			"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
		},
	}
	expectResult := &mockmysql.Result{}
	//ns := m.GetNamespace(namespaceName)
	//assert.Equal()
	//se.ExecuteCommand(mysql.ComQuery, []byte(sql))

	reqCtx := util.NewRequestContext()
	reqCtx.Set(util.StmtType, parser.StmtSelect)

	//se.isAutoCommit()
	rs, err := se.ExecuteSQLs(reqCtx, sqls)
	assert.Equal(t, nil, err)
	assert.Equal(t, rs, expectResult)
}

func prepareSessionExecutor() (*SessionExecutor, error) {
	var userName = "test_executor"
	var namespaceName = "test_executor_namespace"
	var database = "db_ks"

	m, err := prepareNamespaceManager()
	if err != nil {
		return nil, err
	}
	executor := newSessionExecutor(m)
	executor.user = userName

	collationID := 33 // "utf8"
	executor.SetCollationID(mysql.CollationID(collationID))
	executor.SetCharset("utf8")
	// set database
	executor.SetDatabase(database)
	executor.namespace = namespaceName
	return executor, nil
}

func prepareNamespaceManager() (*Manager, error) {
	nsStr := `
{
    "name": "test_executor_namespace",
    "online": true,
    "read_only": true,
    "allowed_dbs": {
        "db_ks": true,
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_ks": "db_ks",
        "db_mycat": "db_mycat_0"
    },
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:3306",
			"slave":[
				"127.0.0.1:3307"	
			],
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:13306",
			"slave":[
				"127.0.0.1:13307"
			],
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        }
    ],
    "shard_rules": [
        {
            "db": "db_ks",
            "table": "tbl_ks",
            "type": "mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        } 
	],
    "users": [
        {
            "user_name": "test_executor",
            "password": "test_executor",
            "namespace": "test_executor_namespace",
            "rw_flag": 2,
            "rw_split": 1
        }
    ],
    "default_slice": "slice-0"
}`

	namespaceName := "test_executor_namespace"
	namespaceConfig := &models.Namespace{}
	if err := json.Unmarshal([]byte(nsStr), namespaceConfig); err != nil {
		fmt.Println(err)
		return nil, err
	}

	m := NewManager()
	current, _, _ := m.switchIndex.Get()
	// init namespace
	namespaceConfigs := map[string]*models.Namespace{namespaceName: namespaceConfig}
	m.namespaces[current] = CreateNamespaceManager(namespaceConfigs)

	user, err := CreateUserManager(namespaceConfigs)
	if err != nil {
		fmt.Println(err, "===============")
		return nil, err
	}
	m.users[current] = user

	return m, nil
}
