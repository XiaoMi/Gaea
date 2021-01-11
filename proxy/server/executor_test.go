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
	"fmt"
	"testing"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/backend/mocks"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
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

func TestExecute(t *testing.T) {
	se, err := prepareSessionExecutor()
	if err != nil {
		t.Fatal("prepare session executer error:", err)
		return
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	slice0MasterPool := new(mocks.ConnectionPool)
	slice0SlavePool := new(mocks.ConnectionPool)
	slice1MasterPool := new(mocks.ConnectionPool)
	slice1SlavePool := new(mocks.ConnectionPool)
	se.manager.GetNamespace("test_executor_namespace").slices["slice-0"].Master = slice0MasterPool
	se.manager.GetNamespace("test_executor_namespace").slices["slice-0"].Slave = []backend.ConnectionPool{slice0SlavePool}
	se.manager.GetNamespace("test_executor_namespace").slices["slice-1"].Master = slice1MasterPool
	se.manager.GetNamespace("test_executor_namespace").slices["slice-1"].Slave = []backend.ConnectionPool{slice1SlavePool}

	expectResult1 := &mysql.Result{}
	expectResult2 := &mysql.Result{}
	//slice-0
	ctx := context.Background()
	slice0MasterConn := new(mocks.PooledConnect)
	slice0MasterConn.On("GetConnectionID").Return(int64(1))
	slice0MasterPool.On("Get", ctx).Return(slice0MasterConn, nil).Once()
	slice0MasterConn.On("UseDB", "db_mycat_0").Return(nil)
	slice0MasterConn.On("SetCharset", "utf8", mysql.CharsetIds["utf8"]).Return(false, nil)
	slice0MasterConn.On("SetSessionVariables", mysql.NewSessionVariables()).Return(false, nil)
	slice0MasterConn.On("GetAddr").Return("127.0.0.1:3306")
	slice0MasterConn.On("Execute", "SELECT * FROM `tbl_mycat` WHERE `k`=0", defaultMaxSqlResultSize).Return(expectResult1, nil)
	slice0MasterConn.On("Recycle").Return(nil)

	//slice-1
	slice1MasterConn := new(mocks.PooledConnect)
	slice1MasterConn.On("GetConnectionID").Return(int64(2))
	slice1MasterPool.On("Get", ctx).Return(slice1MasterConn, nil).Once()
	slice1MasterConn.On("UseDB", "db_mycat_2").Return(nil)
	slice1MasterConn.On("SetCharset", "utf8", mysql.CharsetIds["utf8"]).Return(false, nil)
	slice1MasterConn.On("SetSessionVariables", mysql.NewSessionVariables()).Return(false, nil)
	slice1MasterConn.On("GetAddr").Return("127.0.0.1:3306")
	slice1MasterConn.On("Execute", "SELECT * FROM `tbl_mycat` WHERE `k`=0", defaultMaxSqlResultSize).Return(expectResult2, nil)
	slice1MasterConn.On("Recycle").Return(nil)

	sqls := map[string]map[string][]string{
		"slice-0": {
			"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
		},
		"slice-1": {
			"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
		},
	}

	ret := make([]*mysql.Result, 0)
	ret = append(ret, expectResult1, expectResult2)

	reqCtx := util.NewRequestContext()
	reqCtx.Set(util.StmtType, parser.StmtInsert)

	rs, err := se.ExecuteSQLs(reqCtx, sqls)
	assert.Equal(t, nil, err)
	assert.Equal(t, rs, ret)
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
	proxyCfg := `
; config type, etcd/file, you can test gaea with file type, you shoud use etcd in production
config_type=etcd
;file config path, 具体配置放到file_config_path的namespace目录下，该下级目录为固定目录
file_config_path=./etc/file
;coordinator addr
coordinator_addr=http://127.0.0.1:2379
;远程配置(当前为etcd)根目录
;将会废弃该配置项，通过cluster name识别root
coordinator_root=/gaea
;etcd user config
username=root
password=root
;environ
environ=local
;service name
service_name=gaea_proxy
;gaea_proxy cluster name
cluster_name=gaea
;log config
log_path=./logs
log_level=Notice
log_filename=gaea
log_output=file

;admin addr
admin_addr=0.0.0.0:13307
; basic auth
admin_user=admin
admin_password=admin

;proxy addr
proto_type=tcp4
proxy_addr=0.0.0.0:13306
proxy_charset=utf8

;slow sql time, when execute time is higher than this, log it, unit: ms
slow_sql_time=100

;close session after session timeout, unit: seconds
session_timeout=3600

;stats conf
stats_enabled=true 

;encrypt key
encrypt_key=1234abcd5678efg*
`

	nsCfg := `
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
    "default_slice": "slice-0",
	"max_sql_execute_time": 0
}`

	//加载proxy配置
	var proxy = &models.Proxy{}
	cfg, err := ini.Load([]byte(proxyCfg))
	if err != nil {
		return nil, err
	}
	if err = cfg.MapTo(proxy); err != nil {
		return nil, err
	}

	//加载namespace配置
	namespaceName := "test_executor_namespace"
	namespaceConfig := &models.Namespace{}
	if err := json.Unmarshal([]byte(nsCfg), namespaceConfig); err != nil {
		return nil, err
	}

	m := NewManager()
	// init statistics
	statisticManager, err := CreateStatisticManager(proxy, m)
	if err != nil {
		log.Warn("init stats manager failed, %v", err)
		return nil, err
	}
	m.statistics = statisticManager

	// init namespace
	current, _, _ := m.switchIndex.Get()
	namespaceConfigs := map[string]*models.Namespace{namespaceName: namespaceConfig}
	m.namespaces[current] = CreateNamespaceManager(namespaceConfigs)
	user, err := CreateUserManager(namespaceConfigs)
	if err != nil {
		return nil, err
	}
	m.users[current] = user
	return m, nil
}
