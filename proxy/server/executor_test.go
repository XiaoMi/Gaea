// Copyright 2019 The Gaea Authors. All Rights Reserved.
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
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/proxy/plan"
	"github.com/XiaoMi/Gaea/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
	"reflect"
	"strings"
	"sync"
	"testing"
)

var localManager *Manager

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
						t.Errorf("not equal, fromSlave: %v, actual: %v", test.expect, actual)
					}
				}
			}
		})
	}
}

func (se *SessionExecutor) forTest(sql string, ctx *util.RequestContext) error {
	_, err := se.doQuery(ctx, sql)
	return err
}

// Test for CheckSelectLock
func TestTokensSplit(t *testing.T) {
	type testCase struct {
		sql             string
		expectedTrimmed []string
	}

	testCases := []testCase{
		{"select id,name from t1", []string{"select", "id", "name", "from", "t1"}},
		{"select id,name from t1\n", []string{"select", "id", "name", "from", "t1"}},
		{"select id,name from t1\r", []string{"select", "id", "name", "from", "t1"}},
		{"select \t id,name \nfrom t1", []string{"select", "id", "name", "from", "t1"}},
		{"select / id,name from t1", []string{"select", "id", "name", "from", "t1"}},
	}

	for _, test := range testCases {
		res := strings.FieldsFunc(test.sql, func(r rune) bool {
			return r == ' ' || r == ',' ||
				r == '\t' || r == '/' ||
				r == '\n' || r == '\r'
		})
		assert.Equal(t, res, test.expectedTrimmed)
	}
}

func TestExecute(t *testing.T) {
	se, err := prepareSessionExecutor()
	if err != nil {
		t.Fatal("prepare session executer error:", err)
		return
	}

	reqCtx := util.NewRequestContext()

	reqCtx.Set(util.StmtType, parser.StmtSelect)
	reqCtx.Set(util.FromSlave, 0)

	type testCase struct {
		sql             string
		expectFromSlave int
	}

	testCases := []testCase{
		{"/*master*/ select * from t1", 0},
		{"select * from t1", 1},
		{"select /*master*/ * from t1", 0},
		{"/*master*/ select /*master*/ * from t1", 0},
		{"/*master*/ select /* master*/ * from t1", 0},
		{"select /* master*/ * from t1", 1},
		{"/*master */ select * from t1", 1},
		{"/*master*/ select * from t1 inner join t2 on t1.id = t2.id", 0},
		{"select /*master*/ * from t1 inner join t2 on t1.id = t2.id", 0},
	}

	for _, ca := range testCases {
		t.Run(ca.sql, func(t *testing.T) {
			_ = se.forTest(ca.sql, reqCtx)
			assert.Equal(t, reqCtx.Get(util.FromSlave).(int), ca.expectFromSlave)
			reqCtx.Set(util.FromSlave, 0)
		})
	}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	slice0MasterPool := backend.NewMockConnectionPool(mockCtl)
	slice1MasterPool := backend.NewMockConnectionPool(mockCtl)

	slice0Status := sync.Map{}
	slice0Status.Store(0, backend.UP)

	slice1Status := sync.Map{}
	slice1Status.Store(0, backend.UP)

	se.manager.GetNamespace("test_executor_namespace").slices["slice-0"].Master = &backend.DBInfo{ConnPool: []backend.ConnectionPool{slice0MasterPool}, StatusMap: slice0Status}
	se.manager.GetNamespace("test_executor_namespace").slices["slice-0"].Slave = &backend.DBInfo{}
	se.manager.GetNamespace("test_executor_namespace").slices["slice-1"].Master = &backend.DBInfo{ConnPool: []backend.ConnectionPool{slice1MasterPool}, StatusMap: slice1Status}
	se.manager.GetNamespace("test_executor_namespace").slices["slice-1"].Slave = &backend.DBInfo{}

	expectResult1 := &mysql.Result{}
	expectResult2 := &mysql.Result{}
	//slice-0
	slice0MasterConn := backend.NewMockPooledConnect(mockCtl)
	slice0MasterConn.EXPECT().GetConnectionID().Return(int64(1)).AnyTimes()
	slice0MasterConn.EXPECT().UseDB("db_mycat_0").Return(nil)
	slice0MasterConn.EXPECT().SetCharset("utf8", mysql.CharsetIds["utf8"]).Return(false, nil)
	slice0MasterConn.EXPECT().SetSessionVariables(mysql.NewSessionVariables()).Return(false, nil)
	slice0MasterConn.EXPECT().GetAddr().Return("127.0.0.1:3306").AnyTimes()
	slice0MasterConn.EXPECT().Execute("SELECT * FROM `tbl_mycat` WHERE `k`=0", defaultMaxSqlResultSize).Return(expectResult1, nil)
	slice0MasterConn.EXPECT().Recycle().Return()

	//slice-1
	slice1MasterConn := backend.NewMockPooledConnect(mockCtl)
	slice1MasterConn.EXPECT().GetConnectionID().Return(int64(2)).AnyTimes()
	slice1MasterConn.EXPECT().UseDB("db_mycat_2").Return(nil)
	slice1MasterConn.EXPECT().SetCharset("utf8", mysql.CharsetIds["utf8"]).Return(false, nil)
	slice1MasterConn.EXPECT().SetSessionVariables(mysql.NewSessionVariables()).Return(false, nil)
	slice1MasterConn.EXPECT().GetAddr().Return("127.0.0.1:3306").AnyTimes()
	slice1MasterConn.EXPECT().Execute("SELECT * FROM `tbl_mycat` WHERE `k`=0", defaultMaxSqlResultSize).Return(expectResult2, nil)
	slice1MasterConn.EXPECT().Recycle().Return()

	slice0MasterPool.EXPECT().Get(context.TODO()).Return(slice0MasterConn, nil)
	slice1MasterPool.EXPECT().Get(context.TODO()).Return(slice1MasterConn, nil)

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

	reqCtx.Set(util.StmtType, parser.StmtInsert)

	rs, err := se.ExecuteSQLs(reqCtx, sqls)
	assert.Equal(t, nil, err)
	assert.Equal(t, rs, ret)
}

func prepareSessionExecutor() (*SessionExecutor, error) {
	var userName = "test_executor"
	var namespaceName = "test_executor_namespace"
	var database = "db_ks"
	var err error
	if localManager == nil {
		localManager, err = prepareNamespaceManager()
		if err != nil {
			return nil, err
		}
	}

	executor := newSessionExecutor(localManager)
	executor.user = userName

	collationID := 33 // "utf8"
	executor.SetCollationID(mysql.CollationID(collationID))
	executor.SetCharset("utf8")
	// set database
	executor.SetDatabase(database)
	executor.namespace = namespaceName
	cc := new(Session)
	cc.proxy = &Server{
		manager:       localManager,
		ServerVersion: "5.7.25-gaea",
	}
	executor.session = cc
	return executor, nil
}

type ModifyNamespaceFunc func(nsConfig *models.Namespace)

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
        },
        {
            "user_name": "test_executor_w",
            "password": "test_executor",
            "namespace": "test_executor_namespace",
            "rw_flag": 2,
            "rw_split": 0
        },
        {
            "user_name": "test_executor_r",
            "password": "test_executor",
            "namespace": "test_executor_namespace",
            "rw_flag": 1,
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

func initNamespaceConfig() *models.Namespace {
	namespaceConfig := &models.Namespace{}
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
        },
        {
            "user_name": "test_executor_w",
            "password": "test_executor",
            "namespace": "test_executor_namespace",
            "rw_flag": 2,
            "rw_split": 0
        },
        {
            "user_name": "test_executor_r",
            "password": "test_executor",
            "namespace": "test_executor_namespace",
            "rw_flag": 1,
            "rw_split": 1
        }

    ],
    "default_slice": "slice-0",
	"max_sql_execute_time": 0
}`
	if err := json.Unmarshal([]byte(nsCfg), namespaceConfig); err != nil {
		return nil
	}
	return namespaceConfig
}
func modifyDefaultNamespace(mnFunc ModifyNamespaceFunc, m *Manager) error {
	nsConfig := initNamespaceConfig()
	if mnFunc != nil {
		mnFunc(nsConfig)
	}
	if err := m.ReloadNamespacePrepare(nsConfig); err != nil {
		return err
	}
	if err := m.ReloadNamespaceCommit(nsConfig.Name); err != nil {
		return err
	}
	return nil
}

func newDefaultSessionExecutor(mnFunc ModifyNamespaceFunc) (*SessionExecutor, error) {
	var err error
	if localManager == nil {
		localManager, err = prepareNamespaceManager()
		if err != nil {
			return nil, err
		}
	}

	err = modifyDefaultNamespace(mnFunc, localManager)
	if err != nil {
		return nil, err
	}

	c := newSessionExecutor(localManager)
	c.namespace = "test_executor_namespace"
	c.user = "test_executor"
	c.db = "db_ks"
	return c, nil
}

// test checkExecuteFromSlave
func TestCanExecuteFromSlave(t *testing.T) {
	var userPriv = map[string]string{
		"read_write_split": "test_executor",   // rw_flag: 2, rw_split: 1
		"write_only":       "test_executor_w", // rw_flag: 2, rw_split: 0
		"read_only":        "test_executor_r", // rw_flag: 1, rw_split: 1
	}
	type TestCase struct {
		name             string
		sql              string
		userList         []string
		expectFromSlaves []bool
	}

	testCases := []TestCase{
		{
			name:             "test show read_only",
			sql:              `show variables like "read_only";`,
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test show read_",
			sql:              `show variables like "read_";`,
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{true, false, true},
		},
		{
			name:             "test show global status",
			sql:              `show global status like "uptime";`,
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{true, false, true},
		},
		{
			name:             "test select read_only",
			sql:              "select @@read_only;",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test select global read_only",
			sql:              "select @@global.read_only;",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test select simple",
			sql:              "select * from t",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{true, false, true},
		},
		{
			name:             "test select master hint",
			sql:              "/*master*/ select * from t",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test select master hint internal",
			sql:              "select /*master*/ * from t",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test select master hint after",
			sql:              "select * from t /*master*/",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test select for update",
			sql:              "select * from t where id=1 for update",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test select for share mode",
			sql:              "select * from t where id=1 in share mode",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, true},
		},
		{
			name:             "test update",
			sql:              "update t set col1='a' where id=1",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, false},
		},
		{
			name:             "test insert",
			sql:              "insert into t(col1) values('a')",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, false},
		},
		{
			name:             "test delete",
			sql:              "delete from t where id=1",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, false},
		},
		{
			name:             "test replace",
			sql:              "replace into t(col1) values('a')",
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, false},
		},
		{
			name:             "test only comments",
			sql:              `/*!40100 SET @@SQL_MODE='' */`,
			userList:         []string{userPriv["read_write_split"], userPriv["write_only"], userPriv["read_only"]},
			expectFromSlaves: []bool{false, false, false},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			for i, user := range tt.userList {
				se, err := newDefaultSessionExecutor(nil)
				assert.Equal(t, err, nil)
				se.user = user
				_, comments := extractPrefixCommentsAndRewrite(tt.sql, mysql.ServerVersion)
				reqCtx := util.NewRequestContext()
				reqCtx.Set(util.StmtType, parser.Preview(tt.sql))
				_, err = se.getPlan(reqCtx, se.GetNamespace(), se.db, tt.sql, nil)
				if err != nil {
					t.Fatalf("getPlan error.name:%s, sql:%s,err:%s", tt.name, tt.sql, err)
				}
				assert.Equal(t, checkExecuteFromSlave(reqCtx, se, tt.sql, comments), tt.expectFromSlaves[i], tt.name+"-"+tt.userList[i])
			}
		})
	}
}

// test extractPrefixCommentsAndRewrite
func TestCanExecuteJDBCPrefix(t *testing.T) {
	type TestCase struct {
		name         string
		mysqlVersion string
		sql          string
		trimmedSql   string
	}

	testCases := []TestCase{
		{
			name:       "test jdbc more than 8030 with mysql ",
			sql:        "/* mysql-connector-j-8.0.31 (Revision: 0c86fc148d567b62266c2302bdad0f1e7a7e4eba) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@transaction_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
			trimmedSql: "SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@tx_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
		},
		{
			name:       "test jdbc less than 8030",
			sql:        "/* mysql-connector-j-8.0.31 (Revision: 0c86fc148d567b62266c2302bdad0f1e7a7e4eba) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@transaction_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
			trimmedSql: "SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@tx_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
		},
	}

	for _, tt := range testCases {
		trimmedSql, _ := extractPrefixCommentsAndRewrite(tt.sql, mysql.ServerVersion)
		assert.Equal(t, trimmedSql, tt.trimmedSql, tt.name+"-"+tt.sql)
	}
}

// Test for unshard database and table
func TestUnshardPlan(t *testing.T) {
	type TestCase struct {
		name              string
		sql               string
		mnFunc            ModifyNamespaceFunc
		expectUnshardPlan bool
		expectDB          string
		expectSql         string
	}

	defaltDb := "db_ks"
	var (
		testCases = []TestCase{
			{
				name:              "test unshard show variable read_only without shard rules",
				sql:               `show variables like "read_only"`,
				expectUnshardPlan: true,
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
					nsConfig.ShardRules = nil
				},
				expectDB:  defaltDb,
				expectSql: `show variables like "read_only"`,
			},
			{
				name:              "test unshard show columns without shard rules",
				sql:               `show columns from db_unshard.tbl_unshard`,
				expectUnshardPlan: true,
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
					nsConfig.ShardRules = nil
				},
				expectDB:  defaltDb,
				expectSql: `show columns from db_unshard.tbl_unshard`,
			},
			{
				name:              "test unshard select with db but without shard rules",
				sql:               "select * from db_unshard.tbl_unshard",
				expectUnshardPlan: true,
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
					nsConfig.ShardRules = nil
				},
				expectDB:  defaltDb,
				expectSql: "select * from db_unshard.tbl_unshard",
			},
			{
				name:              "test unshard select with db and subquery",
				sql:               "select * from (select id from db_unshard.tbl_unshard_1) db_unshard.tbl_unshard",
				expectUnshardPlan: true,
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
					nsConfig.ShardRules = nil
				},
				expectDB:  defaltDb,
				expectSql: "select * from (select id from db_unshard.tbl_unshard_1) db_unshard.tbl_unshard",
			},
			{
				name:              "test unshard select without db but subquery",
				sql:               "select db_unshard.* from (select * from tbl_unshard where id>1000) as a where a.id<2000;",
				expectUnshardPlan: true,
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
					nsConfig.ShardRules = nil
				},
				expectDB:  defaltDb,
				expectSql: "select db_unshard.* from (select * from tbl_unshard where id>1000) as a where a.id<2000;",
			},
			{
				name: "test unshard select with db and with shard rules",
				sql:  "select * from db_unshard.tbl_unshard",
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
				},
				expectUnshardPlan: true,
				expectDB:          "db_unshard",
				expectSql:         "select * from db_unshard.tbl_unshard",
			},
			{
				name:              "test unshard select without db but with shard rules",
				sql:               "select * from tbl_unshard", // tbl_unshard is unshard in shard db db_ks
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "select * from tbl_unshard",
			},
			{
				name:              "test unshard select with db but with shard rules",
				sql:               "select * from db_ks.tbl_unshard", // tbl_unshard is unshard in shard db db_ks
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "select * from db_ks.tbl_unshard",
			},
			{
				name: "test unshard select by parsing sql with renamed origin database, but use default db",
				sql:  "select * from db_unshard_rename.tbl_unshard", // db_unshard_rename_0 is backend physical database
				mnFunc: func(nsConfig *models.Namespace) {
					// db_unshard_rename is renamed database of backend physical database db_unshard_rename_0
					nsConfig.AllowedDBS = map[string]bool{"db_unshard_rename": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard_rename": "db_unshard_rename_0"}
				},
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "SELECT * FROM `db_unshard_rename_0`.`tbl_unshard`",
			},
			{
				name: "test unshard select without db with renamed origin database, but use default db",
				sql:  "select user()", // db_unshard_rename_0 is backend physical database
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard_rename": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard_rename": "db_unshard_rename_0"}
					nsConfig.ShardRules = nil
				},
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "select user()",
			},
			{
				name:              "test unshard select without db and table",
				sql:               "select connection_id();",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "select connection_id();",
			},
			{
				name: "test unshard union select without db",
				sql:  "select * from tbl_unshard_1 union select * from tbl_unshard_2",
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.ShardRules = nil
				},
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "select * from tbl_unshard_1 union select * from tbl_unshard_2",
			},
			{
				name: "test unshard union select with db",
				sql:  "select * from db_unshard.t1 union select * from db_unshard.t2", // db_ks.tbl_ks is a shard table
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
				},
				expectUnshardPlan: true,
				expectDB:          "db_unshard",
				expectSql:         "select * from db_unshard.t1 union select * from db_unshard.t2",
			},
			{
				name:              "test shard simple select without db",
				sql:               "select * from tbl_ks",
				mnFunc:            nil,
				expectUnshardPlan: false,
			},
			{
				name:              "test shard select with db",
				sql:               "select * from db_ks.tbl_ks",
				mnFunc:            nil,
				expectUnshardPlan: false,
			},
			{
				name: "test unshard update with no shard rules",
				sql:  "UPDATE db_unshard.tbl_unshard SET col1='a' WHERE id=6",
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
					nsConfig.ShardRules = nil
				},
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "UPDATE db_unshard.tbl_unshard SET col1='a' WHERE id=6",
			},
			{
				name: "test update unshard with no db and no shard rules",
				sql:  "UPDATE tbl_unshard SET col1='a' WHERE id=6",
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = map[string]string{"db_unshard": "db_unshard"}
					nsConfig.ShardRules = nil
				},
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "UPDATE tbl_unshard SET col1='a' WHERE id=6",
			},
			{
				name:              "test update unshard without db",
				sql:               "UPDATE tbl_unshard SET col1='a' WHERE id=6",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "UPDATE tbl_unshard SET col1='a' WHERE id=6",
			},
			{
				name:              "test shard update without db",
				sql:               "UPDATE tbl_ks SET col1='a' WHERE id=6",
				expectUnshardPlan: false,
			},
			{
				name:              "test unshard insert without db",
				sql:               "insert into tbl_unshard values('a')",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "insert into tbl_unshard values('a')",
			},
			{
				name:              "test unshard insert without db and select from",
				sql:               "insert into tbl_unshard(col1) (select col1 from tbl_unshard_1);",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "insert into tbl_unshard(col1) (select col1 from tbl_unshard_1);",
			},
			{
				name:              "test unshard insert with db",
				sql:               "insert into db_unshard.tbl_unshard values('a')",
				mnFunc:            nil,
				expectUnshardPlan: true,
				expectDB:          "db_unshard",
				expectSql:         "insert into db_unshard.tbl_unshard values('a')",
			},
			{
				name:              "test shard insert without db",
				sql:               "insert into tbl_ks(id,col1) values(1,'a')",
				expectUnshardPlan: false,
			},
			{
				name:              "test unshard replace without db",
				sql:               "replace into tbl_unshard values('a')",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "replace into tbl_unshard values('a')",
			},
			{
				name:              "test shard replace without db",
				sql:               "replace into tbl_ks(id,col1) values(1,'a')",
				expectUnshardPlan: false,
			},
			{
				name:              "test unshard delete without db",
				sql:               "delete from tbl_unshard where id=1",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "delete from tbl_unshard where id=1",
			},
			{
				name:              "test shard delete without db",
				sql:               "delete from tbl_ks where id=1",
				expectUnshardPlan: false,
			},
			{
				name:              "test unshard delete with db",
				sql:               "delete from db_unshard.tbl_unshard where id=1",
				expectUnshardPlan: true,
				expectDB:          "db_unshard",
				expectSql:         "delete from db_unshard.tbl_unshard where id=1",
			},
			{
				name:              "test shard delete with db",
				sql:               "delete from db_ks.tbl_ks where id=1",
				expectUnshardPlan: false,
			},
			// TODO: fix sql execute when use renamed db
			{
				name:              "test unshard function without db",
				sql:               "SELECT HEX(WEIGHT_STRING('ab' AS CHAR(4)));",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "SELECT HEX(WEIGHT_STRING('ab' AS CHAR(4)));",
			},
			{
				name: "test unshard function with db and no phydbs",
				sql:  "SELECT HEX(WEIGHT_STRING('ab' AS CHAR(4)));",
				mnFunc: func(nsConfig *models.Namespace) {
					nsConfig.AllowedDBS = map[string]bool{"db_unshard": true}
					nsConfig.DefaultPhyDBS = nil
				},
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "SELECT HEX(WEIGHT_STRING('ab' AS CHAR(4)));",
			},
			{
				name:              "test unshard function div without db",
				sql:               "SELECT 5 DIV 2, -5 DIV 2, 5 DIV -2, -5 DIV -2;",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "SELECT 5 DIV 2, -5 DIV 2, 5 DIV -2, -5 DIV -2;",
			},
			{
				name:              "test only comment",
				sql:               "/*!40100 SET @@SQL_MODE='' */",
				expectUnshardPlan: true,
				expectDB:          defaltDb,
				expectSql:         "SET @@SESSION.`sql_mode`=''",
			},
		}
	)
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			se, err := newDefaultSessionExecutor(tt.mnFunc)
			assert.Equal(t, err, nil)
			reqCtx := util.NewRequestContext()
			p, err := se.getPlan(reqCtx, se.GetNamespace(), se.db, tt.sql, nil)
			if err != nil {
				t.Fatalf("getPlan error.name:%s,err:%s\n", tt.name, err)
			}
			up, ok := p.(*plan.UnshardPlan)
			if ok != tt.expectUnshardPlan {
				t.Fatalf("getPlan not equal, name:%s,sql:%s, plan type:%T, expect unshard plan:%v\n", tt.name, tt.sql, p, tt.expectUnshardPlan)
				return
			}
			if !ok {
				return
			}
			v := reflect.ValueOf(up)
			resDb := v.Elem().Field(1).String()
			resSql := v.Elem().Field(3).String()
			if tt.expectDB != resDb {
				t.Fatalf("getPlan db error, name:%s,sql:%s, get db:%s, expect db:%s\n", tt.name, tt.sql, resDb, tt.expectDB)
			}
			if tt.expectSql != resSql {
				t.Fatalf("getPlan sql error, name:%s,sql:%s, get sql:%s, expect sql:%s\n", tt.name, tt.sql, resSql, tt.expectSql)
			}
		})
	}
}

func TestPreRewriteSQL(t *testing.T) {
	type TestCase struct {
		name         string
		mysqlVersion string
		sql          string
		expectSQL    string
	}

	testCases := []TestCase{
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "/* mysql-connector-j-8.0.31 (Revision: 0c86fc148d567b62266c2302bdad0f1e7a7e4eba) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@transaction_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
			expectSQL:    "/* mysql-connector-j-8.0.31 (Revision: 0c86fc148d567b62266c2302bdad0f1e7a7e4eba) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@tx_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "/* mysql-connector-j-8.0.31 (Revision: 0c86fc148d567b62266c2302bdad0f1e7a7e4eba) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@transaction_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
			expectSQL:    "/* mysql-connector-j-8.0.31 (Revision: 0c86fc148d567b62266c2302bdad0f1e7a7e4eba) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@collation_server AS collation_server, @@collation_connection AS collation_connection, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_write_timeout AS net_write_timeout, @@performance_schema AS performance_schema, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@tx_isolation AS transaction_isolation, @@wait_timeout AS wait_timeout",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "/* mysql-connector-java-5.1.39 ( Revision: 3289a357af6d09ecc1a10fd3c26e95183e5790ad ) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_buffer_length AS net_buffer_length, @@net_write_timeout AS net_write_timeout, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@tx_isolation AS tx_isolation, @@wait_timeout AS wait_timeout",
			expectSQL:    "/* mysql-connector-java-5.1.39 ( Revision: 3289a357af6d09ecc1a10fd3c26e95183e5790ad ) */SELECT  @@session.auto_increment_increment AS auto_increment_increment, @@character_set_client AS character_set_client, @@character_set_connection AS character_set_connection, @@character_set_results AS character_set_results, @@character_set_server AS character_set_server, @@init_connect AS init_connect, @@interactive_timeout AS interactive_timeout, @@license AS license, @@lower_case_table_names AS lower_case_table_names, @@max_allowed_packet AS max_allowed_packet, @@net_buffer_length AS net_buffer_length, @@net_write_timeout AS net_write_timeout, @@query_cache_size AS query_cache_size, @@query_cache_type AS query_cache_type, @@sql_mode AS sql_mode, @@system_time_zone AS system_time_zone, @@time_zone AS time_zone, @@tx_isolation AS tx_isolation, @@wait_timeout AS wait_timeout",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "select @@transaction_read_only",
			expectSQL:    "select @@tx_read_only",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "select @@session.transaction_read_only",
			expectSQL:    "select @@session.tx_read_only",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "select /*master*/ @@session.transaction_read_only",
			expectSQL:    "select /*master*/ @@session.tx_read_only",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "select @@session.transaction_read_only /*test*/",
			expectSQL:    "select @@session.transaction_read_only /*test*/",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "SELECT @@session.transaction_isolation",
			expectSQL:    "SELECT @@session.tx_isolation",
		},
		{
			mysqlVersion: "5.7.25-gaea",
			sql:          "SELECT @@transaction_isolation",
			expectSQL:    "SELECT @@tx_isolation",
		},
		{
			mysqlVersion: "8.0.25-gaea",
			sql:          "select @@session.transaction_read_only",
			expectSQL:    "select @@session.transaction_read_only",
		},
		{
			mysqlVersion: "8.0.25-gaea",
			sql:          "SELECT @@transaction_isolation",
			expectSQL:    "SELECT @@transaction_isolation",
		},
	}

	for _, tt := range testCases {
		sql := preRewriteSQL(tt.sql, tt.mysqlVersion)
		assert.Equal(t, tt.expectSQL, sql, tt.sql+"-"+tt.mysqlVersion)
	}
}

func TestExtractPrefixCommentsAndRewrite(t *testing.T) {
	type args struct {
		sql     string
		version string
	}
	tests := []struct {
		name          string
		args          args
		expectTrimmed string
		expectComment parser.MarginComments
	}{
		{
			name: "test comment leading",
			args: args{
				sql:     "/*master*/ select * from t",
				version: "5.7.25-gaea",
			},
			expectTrimmed: "select * from t",
			expectComment: parser.MarginComments{
				Leading:  "/*master*/ ",
				Trailing: "",
			},
		},
		{
			name: "test comments trailing",
			args: args{
				sql:     "select * from t /*master*/",
				version: "5.7.25-gaea",
			},
			expectTrimmed: "select * from t",
			expectComment: parser.MarginComments{
				Leading:  "",
				Trailing: " /*master*/",
			},
		},
		{
			name: "test comments internal",
			args: args{
				sql:     "select /*master*/ * from t",
				version: "5.7.25-gaea",
			},
			expectTrimmed: "select /*master*/ * from t",
			expectComment: parser.MarginComments{
				Leading:  "",
				Trailing: "",
			},
		},
		{
			name: "test comments mycat trailing",
			args: args{
				sql:     "select * from t /* !mycat:sql=select 1 from xm_order where order_id = 4230329282206593 */",
				version: "5.7.25-gaea",
			},
			expectTrimmed: "select * from t",
			expectComment: parser.MarginComments{
				Leading:  "",
				Trailing: " /* !mycat:sql=select 1 from xm_order where order_id = 4230329282206593 */",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTrimmed, gotComment := extractPrefixCommentsAndRewrite(tt.args.sql, tt.args.version)
			assert.Equalf(t, tt.expectTrimmed, gotTrimmed, "extractPrefixCommentsAndRewrite(%v, %v)", tt.args.sql, tt.args.version)
			assert.Equalf(t, tt.expectComment, gotComment, "extractPrefixCommentsAndRewrite(%v, %v)", tt.args.sql, tt.args.version)
		})
	}
}
