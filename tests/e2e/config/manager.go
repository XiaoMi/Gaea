// Copyright 2024 The Gaea Authors. All Rights Reserved.
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

package config

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/cc"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/XiaoMi/Gaea/util/requests"
)

const (
	// 注意这里的端口要和gaea_cc.ini 中的addr=0.0.0.0:23306
	gaeaCCBaseRouter        = "http://localhost:23306/api/cc/namespace/"
	gaeaCCAdminUsername     = "admin"
	gaeaCCAdminPassword     = "admin"
	defaultMysqlAdminUser   = "superroot"
	defaultMysqlAdminPasswd = "superroot"
	defaultGaeaBackendUser  = "gaea_backend_user"
	defaultGaeaBackendPass  = "gaea_backend_pass"
	defaultGaeaUser         = "gaea_user"
	defaultGaeaPass         = "gaea_pass"
	defaultHost             = "127.0.0.1"

	DefaultE2eDatabase = "db_e2e_test"
	DefaultE2eTable    = "tbl_e2e_test"
	// SliceSingleTestMaster 表示测试（对照）的单主 MySQL 集群 3379
	SliceSingleTestMaster = "slice-single-test-master"
	// SliceSingleMaster 表示gaea使用的单主 MySQL 集群 3349
	SliceSingleMaster = "slice-single-master"
	// SliceDualMaster 表示测试的多主 MySQL 集群 3319 3349
	SliceDualMaster = "slice-dual-master"
	// SliceSingleSlave 表示测试的主从 MySQL 集群 3319 3329
	SliceSingleSlave = "slice-single-slave"
	// SliceDualSlave 表示测试的主从 MySQL 集群 3319 3329 3339
	SliceDualSlave = "slice-dual-slave"
	// LogExpression 标识 Gaea SQL Log 的格式
	LogExpression = `\[(.*?)\] \[INFO\] \[(\d+)\] OK - (\d+\.\d+)ms - ns=(.*?), (.*?)@(.*?)->(.*?)/(.*?), connect_id=(\d+), mysql_connect_id=(\d+), transaction=(\w+)\|(.*?)$`
)

var logDirectory = "cmd/logs"
var E2eMgr *E2eManager

// Embed namespace templates directly into the Go binary.
var (
	//go:embed ns/default.template
	DefaultNamespaceTmpl string

	//go:embed ns/unshard.template
	UnShardNamespaceTmpl string

	//go:embed ns/mycat_long.template
	MycatLongNamespaceTmpl string

	//go:embed ns/mycat_mod.template
	MycatModNamespaceTmpl string

	//go:embed ns/mycat_murmur.template
	MycatMurmurNamespaceTmpl string

	//go:embed ns/mycat_string.template
	MycatStringNamespaceTmpl string

	//go:embed ns/shard.template
	ShardNamespaceTmpl string

	//go:embed ns/kingshard_hash.template
	KingShardHashNamespaceTmpl string

	//go:embed ns/kingshard_mod.template
	KingShardModNamespaceTmpl string

	//go:embed ns/unshard_dml.template
	UnShardDMLNamespaceTmpl string
)

type E2eManager struct {
	// 用于管理所有的测试用例
	NsManager       *NamespaceRegisterManager
	GCluster        *GaeaCluster
	NsSlices        map[string]*NsSlice
	BasePath        string
	StartTime       time.Time
	Db              string
	Table           string
	openConnections []*sql.DB // 存储所有打开的数据库连接
}

func NewE2eManager() *E2eManager {
	GaeaUsers := []*models.User{
		{
			UserName: defaultGaeaUser,
			Password: defaultGaeaPass,
			RWFlag:   2,
			RWSplit:  1,
		},
		{
			UserName: defaultGaeaUser + "_r",
			Password: defaultGaeaPass + "_r",
			RWFlag:   1,
			RWSplit:  0,
		},
		{
			UserName: defaultGaeaUser + "_w",
			Password: defaultGaeaPass + "_w",
			RWFlag:   2,
			RWSplit:  0,
		},
	}

	// 3379
	sliceSingleTestMaster := &NsSlice{
		Name: SliceSingleTestMaster,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        defaultGaeaBackendUser,
				Password:        defaultGaeaBackendPass,
				Master:          fmt.Sprintf("%s:%d", defaultHost, 3379),
				Slaves:          nil,
				StatisticSlaves: nil,
				Capacity:        12,
				HealthCheckSql:  "",
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		GaeaUsers: GaeaUsers,
	}
	// 3349
	sliceSingleMaster := &NsSlice{
		Name: SliceSingleMaster,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        defaultGaeaBackendUser,
				Password:        defaultGaeaBackendPass,
				Master:          fmt.Sprintf("%s:%d", defaultHost, 3349),
				Slaves:          nil,
				StatisticSlaves: nil,
				Capacity:        12,
				HealthCheckSql:  "",
				MaxCapacity:     24,
				IdleTimeout:     60,
				Capability:      500357,
			},
		},
		GaeaUsers: GaeaUsers,
	}
	// 3319 3329 3339
	sliceMasterSlaves := &NsSlice{
		Name: SliceDualSlave,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        defaultGaeaBackendUser,
				Password:        defaultGaeaBackendPass,
				Master:          fmt.Sprintf("%s:%d", defaultHost, 3319),
				Slaves:          []string{fmt.Sprintf("%s:%d", defaultHost, 3329), fmt.Sprintf("%s:%d", defaultHost, 3339)},
				StatisticSlaves: nil,
				Capacity:        12,
				HealthCheckSql:  "",
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		GaeaUsers: GaeaUsers,
	}

	// 3319 3329
	sliceMasterSingleSlave := &NsSlice{
		Name: SliceSingleSlave,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        defaultGaeaBackendUser,
				Password:        defaultGaeaBackendPass,
				Master:          fmt.Sprintf("%s:%d", defaultHost, 3319),
				Slaves:          []string{fmt.Sprintf("%s:%d", defaultHost, 3329)},
				StatisticSlaves: nil,
				Capacity:        12,
				HealthCheckSql:  "",
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		GaeaUsers: GaeaUsers,
	}
	// 3319 3349
	sliceMultiMasters := &NsSlice{
		Name: SliceDualMaster,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        defaultGaeaBackendUser,
				Password:        defaultGaeaBackendPass,
				Master:          fmt.Sprintf("%s:%d", defaultHost, 3319),
				Slaves:          nil,
				StatisticSlaves: nil,
				Capacity:        12,
				HealthCheckSql:  "",
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
			{
				Name:            "slice-1",
				UserName:        defaultGaeaBackendUser,
				Password:        defaultGaeaBackendPass,
				Master:          fmt.Sprintf("%s:%d", defaultHost, 3349),
				Slaves:          nil,
				StatisticSlaves: nil,
				Capacity:        12,
				HealthCheckSql:  "",
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		GaeaUsers: GaeaUsers,
	}

	nsSlices := map[string]*NsSlice{
		//3349
		sliceSingleMaster.Name: sliceSingleMaster,
		//3379
		sliceSingleTestMaster.Name: sliceSingleTestMaster,
		//3319 3329 3339
		sliceMasterSlaves.Name: sliceMasterSlaves,
		//3329 3339
		sliceMultiMasters.Name: sliceMultiMasters,
		//3319 3329
		sliceMasterSingleSlave.Name: sliceMasterSingleSlave,
	}
	E2eMgr = &E2eManager{
		NsManager: NewNamespaceRegisterManger(),
		GCluster: &GaeaCluster{
			Host:          defaultHost,
			Port:          13306,
			ReadWriteUser: GaeaUsers[0],
			ReadUser:      GaeaUsers[1],
			WriteUser:     GaeaUsers[2],
			LogDirectory:  filepath.Join(basePath, logDirectory),
		},
		NsSlices:  nsSlices,
		BasePath:  basePath,
		StartTime: time.Now(),
		Db:        DefaultE2eDatabase,
		Table:     DefaultE2eTable,
	}
	return E2eMgr
}

func (e *E2eManager) GetReadWriteGaeaUserConn() (*sql.DB, error) {
	conn, err := InitConn(e.GCluster.ReadWriteUser.UserName, e.GCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
	if err != nil {
		return nil, err
	}
	e.openConnections = append(e.openConnections, conn) // 添加到追踪列表
	return conn, nil
}

func (e *E2eManager) GetReadWriteGaeaUserDBConn(db string) (*sql.DB, error) {
	conn, err := InitConn(e.GCluster.ReadWriteUser.UserName, e.GCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), db)
	if err != nil {
		return nil, err
	}
	e.openConnections = append(e.openConnections, conn) // 添加到追踪列表
	return conn, nil
}

func (e *E2eManager) GetReadGaeaUserConn() (*sql.DB, error) {
	conn, err := InitConn(e.GCluster.ReadUser.UserName, e.GCluster.ReadUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
	if err != nil {
		return nil, err
	}
	e.openConnections = append(e.openConnections, conn) // 添加到追踪列表
	return conn, err
}

func (e *E2eManager) GetReadGaeaUserDBConn(db string) (*sql.DB, error) {
	conn, err := InitConn(e.GCluster.ReadUser.UserName, e.GCluster.ReadUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), db)
	if err != nil {
		return nil, err
	}
	e.openConnections = append(e.openConnections, conn) // 添加到追踪列表
	return conn, err
}

func (e *E2eManager) GetWriteGaeaUserConn() (*sql.DB, error) {
	conn, err := InitConn(e.GCluster.WriteUser.UserName, e.GCluster.WriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
	if err != nil {
		return nil, err
	}
	e.openConnections = append(e.openConnections, conn) // 添加到追踪列表
	return conn, err
}

func (e *E2eManager) GetWriteGaeaUserDBConn(db string) (*sql.DB, error) {
	conn, err := InitConn(e.GCluster.WriteUser.UserName, e.GCluster.WriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), db)
	if err != nil {
		return nil, err
	}
	e.openConnections = append(e.openConnections, conn) // 添加到追踪列表
	return conn, err
}

func (e *E2eManager) ModifyNamespace(m *models.Namespace) error {
	if e.NsManager == nil {
		return fmt.Errorf("namespace manager has not been initialized")
	}
	return e.NsManager.modifyNamespace(m)
}

func (e *E2eManager) DeleteNamespace(name string) error {
	return e.NsManager.DeleteNamespace(name)
}

func (e *E2eManager) ListNamespace() ([]string, error) {
	return e.NsManager.ListNamespaces()
}

type NamespaceRegisterManager struct {
	GaeaCCManager *GaeaCCManager
}

func NewNamespaceRegisterManger() *NamespaceRegisterManager {
	return &NamespaceRegisterManager{
		GaeaCCManager: &GaeaCCManager{
			GaeaCCAdminUser:     gaeaCCAdminUsername,
			GaeaCCAdminPassword: gaeaCCAdminPassword,
			GaeaCCBaseRouter:    gaeaCCBaseRouter,
		},
	}
}

func (nr *NamespaceRegisterManager) modifyNamespace(m *models.Namespace) error {
	// 存在就替换
	err := m.Verify()
	if err != nil {
		return err
	}
	err = nr.GaeaCCManager.modifyNamespace(m)
	if err != nil {
		return fmt.Errorf("modify Namespace error:%v", err)
	}
	return nil
}

func (nr *NamespaceRegisterManager) DeleteNamespace(name string) error {
	return nr.GaeaCCManager.deleteNamespace(name)
}

func (nr *NamespaceRegisterManager) ListNamespaces() ([]string, error) {
	return nr.GaeaCCManager.listNamespaces()
}

type GaeaCCManager struct {
	GaeaCCAdminUser     string
	GaeaCCAdminPassword string
	GaeaCCBaseRouter    string
}

func (g *GaeaCCManager) modifyNamespace(n *models.Namespace) error {
	// Serialize the namespace object into JSON
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	req := requests.NewRequest(g.GaeaCCBaseRouter+"modify", requests.Put, map[string]string{"Content-Type": "application/json"}, nil, data)
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	// Send the request
	resp, err := requests.Send(req)
	if err != nil {
		return err
	}
	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(resp.Body))
	}
	return nil
}

// deleteNamespace delete Namespace by name
func (g *GaeaCCManager) deleteNamespace(key string) error {
	// Build the full URL for the DELETE request
	url := g.GaeaCCBaseRouter + "delete/" + key
	// Send the request
	return requests.SendPut(url, g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
}

// listNamespace queries the GaeaCC for a list of namespaces.
func (g *GaeaCCManager) listNamespaces() ([]string, error) {
	// Build the full URL for the DELETE request
	url := g.GaeaCCBaseRouter + "list"
	// Send the request
	resp, err := requests.SendGet(url, g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	if err != nil {
		return []string{}, err
	}
	var listResp cc.ListNamespaceResp
	err = json.Unmarshal(resp.Body, &listResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	// Check if the list namespaces was successful
	return listResp.Data, nil
}

// SearchSqlLog function in the provided code is designed to search through SQL log files for specific entries.
// It takes a search string and a timestamp, and returns a slice of log entries that match the criteria.
// The incoming searchString will be unblocked.
// The function reads log files within a specified directory, looking for files named with the prefix "gaea_sql.log".
// It uses regular expressions to parse and match log entries based on the input parameters.
// If a matching entry is found, it's added to the result slice. The function handles errors such as file access issues and returns an error if any problems occur during the file reading and parsing process.
func (e *E2eManager) SearchSqlLog(searchString string, currentTime time.Time) ([]util.LogEntry, error) {
	// 等待日志落盘
	time.Sleep(100 * time.Millisecond)
	searchString = strings.TrimSuffix(searchString, ";")
	var allEntries []util.LogEntry

	err := filepath.Walk(e.GCluster.LogDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "gaea_sql.log") && info.Mode()&os.ModeSymlink == 0 {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open file:%s error %v", path, err)
			}
			defer file.Close()

			re := regexp.MustCompile(LogExpression)

			entries, err := util.ParseLogEntries(file, re, currentTime, searchString)
			if err != nil {
				return err
			}
			allEntries = append(allEntries, entries...)
		}
		return nil
	})

	if err != nil {
		return []util.LogEntry{}, err
	}
	return allEntries, nil
}
func (e *E2eManager) ClearSqlLog() error {
	err := filepath.Walk(e.GCluster.LogDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "gaea_sql.log") {
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
			if err != nil {
				return fmt.Errorf("open file:%s error %v", path, err)
			}
			defer file.Close()
			file.Truncate(0)
		}
		return nil
	})
	return err
}
