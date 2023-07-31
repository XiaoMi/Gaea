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
	"strings"
	"sync"
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

	reqCtx.Set("stmtType", 1)
	reqCtx.Set(util.FromSlave, 0)

	type sqlAndExpectations struct {
		sql      string
		expected int
	}

	testCase := []sqlAndExpectations{
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

	for _, ca := range testCase {
		_ = se.forTest(ca.sql, reqCtx)
		assert.Equal(t, reqCtx.Get(util.FromSlave).(int), ca.expected)
		reqCtx.Set(util.FromSlave, 0)
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	slice0MasterPool := new(mocks.ConnectionPool)
	slice1MasterPool := new(mocks.ConnectionPool)

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

func TestCanExecuteFromSlave(t *testing.T) {
	type userPriv struct {
		name string
		user string
	}
	type TestCase struct {
		name          string
		sql           string
		userList      []userPriv
		fromSlaveList []bool
	}

	userList := []userPriv{
		{
			name: "read_write",
			user: "test_executor",
		},
		{
			name: "write_only",
			user: "test_executor_w",
		},
		{
			name: "read_only",
			user: "test_executor_r",
		},
	}
	testCases := []TestCase{
		{
			name:          "test select simple",
			sql:           "select * from t",
			userList:      userList,
			fromSlaveList: []bool{true, false, true},
		},
		{
			name:          "test select master hint",
			sql:           "/*master*/ select * from t",
			userList:      userList,
			fromSlaveList: []bool{false, false, true},
		},
		{
			name:          "test select master hint 2",
			sql:           "/*master*/ select * from t",
			userList:      userList,
			fromSlaveList: []bool{false, false, true},
		},
		// test select ... for update/select ... lock in share mode
		{
			name:          "test select for update",
			sql:           "select * from t where id=1 for update",
			userList:      userList,
			fromSlaveList: []bool{false, false, true},
		},
		// test update
		{
			name:          "test update",
			sql:           "update t set col1='a' where id=1",
			userList:      userList,
			fromSlaveList: []bool{false, false, false},
		},
	}

	for _, tt := range testCases {
		for i, user := range tt.userList {
			c, err := newDefaultSessionExecutor()
			assert.Equal(t, err, nil)
			c.user = user.user
			trimmedSql, comments := extractPrefixCommentsAndRewrite(tt.sql, mysql.ServerVersion)
			_, stmt, err := c.getPlan(c.GetNamespace(), c.db, trimmedSql, nil)
			if err != nil {
				t.Fatal("getPlan error:", tt.name)
			}
			assert.Equal(t, canExecuteFromSlave(c, tt.sql, stmt, comments), tt.fromSlaveList[i], tt.name+"-"+tt.userList[i].name)
		}

	}
}

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
		c, err := newDefaultSessionExecutor()
		assert.Equal(t, err, nil)
		trimmedSql, _ := extractPrefixCommentsAndRewrite(tt.sql, mysql.ServerVersion)
		fmt.Println("trimmedSql:", trimmedSql)
		_, _, err = c.getPlan(c.GetNamespace(), c.db, trimmedSql, nil)
		if err != nil {
			t.Fatal("getPlan error:", tt.name)
		}
		assert.Equal(t, trimmedSql, tt.trimmedSql, tt.name+"-"+tt.sql)
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

func newDefaultSessionExecutor() (*SessionExecutor, error) {
	var err error
	if localManager == nil {
		localManager, err = prepareNamespaceManager()
		if err != nil {
			return nil, err
		}
	}

	c := newSessionExecutor(localManager)
	c.namespace = "test_executor_namespace"
	c.user = "test_executor"
	c.db = "db_ks"
	return c, nil
}
