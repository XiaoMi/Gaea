package config

import (
	"bytes"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/XiaoMi/Gaea/models"
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
	LogExpression = `\[(.*?)\] \[NOTICE\] \[(\d+)\] OK - (\d+\.\d+)ms - ns=(.*?), (.*?)@(.*?)->(.*?)/(.*?), mysql_connect_id=(\d+), r=\d+\|(.*?)$`
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
	NsManager *NamespaceRegisterManager
	GCluster  *GaeaCluster
	NsSlices  map[string]*NsSlice
	BasePath  string
	StartTime time.Time
	Db        string
	Table     string
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
				MaxCapacity:     24,
				IdleTimeout:     60,
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
	return InitConn(e.GCluster.ReadWriteUser.UserName, e.GCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) GetReadWriteGaeaUserDBConn(db string) (*sql.DB, error) {
	return InitConn(e.GCluster.ReadWriteUser.UserName, e.GCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), db)
}

func (e *E2eManager) GetReadGaeaUserConn() (*sql.DB, error) {
	return InitConn(e.GCluster.ReadUser.UserName, e.GCluster.ReadUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) GetReadGaeaUserDBConn(db string) (*sql.DB, error) {
	return InitConn(e.GCluster.ReadUser.UserName, e.GCluster.ReadUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), db)
}

func (e *E2eManager) GetWriteGaeaUserConn() (*sql.DB, error) {
	return InitConn(e.GCluster.WriteUser.UserName, e.GCluster.WriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) GetWriteGaeaUserDBConn(db string) (*sql.DB, error) {
	return InitConn(e.GCluster.WriteUser.UserName, e.GCluster.WriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), db)
}

func (e *E2eManager) ModifyNamespace(m *models.Namespace) error {
	return e.NsManager.ModifyNamespace(m)
}

func (e *E2eManager) DeleteNamespace(m *models.Namespace) error {
	return e.NsManager.DeleteNamespace(m)
}

type NamespaceRegisterManager struct {
	Namespaces    map[string]*models.Namespace
	GaeaCCManager *GaeaCCManager
}

func NewNamespaceRegisterManger() *NamespaceRegisterManager {
	return &NamespaceRegisterManager{
		Namespaces: map[string]*models.Namespace{},
		GaeaCCManager: &GaeaCCManager{
			GaeaCCAdminUser:     gaeaCCAdminUsername,
			GaeaCCAdminPassword: gaeaCCAdminPassword,
			GaeaCCBaseRouter:    gaeaCCBaseRouter,
		},
	}
}

func (nr *NamespaceRegisterManager) ModifyNamespace(m *models.Namespace) error {
	// 存在就替换
	err := m.Verify()
	if err != nil {
		return err
	}
	//if _, ok := nr.Namespaces[m.Type]; ok {
	//	return fmt.Errorf("namespace has exists")
	//}
	nr.Namespaces[m.Name] = m
	err = nr.GaeaCCManager.modifyNamespace(m)
	if err != nil {
		return fmt.Errorf("modify Namespace error:%v", err)
	}
	return nil
}

func (nr *NamespaceRegisterManager) DeleteNamespace(m *models.Namespace) error {
	return nr.GaeaCCManager.deleteNamespace(m.Name)
}

func (nr *NamespaceRegisterManager) unRegisterNamespaceByKey(key string) error {
	// 在 ETCD 中删除 key
	err := nr.GaeaCCManager.deleteNamespace(key)
	delete(nr.Namespaces, key)
	if err != nil {
		return fmt.Errorf("error deleting from ETCD: %v", err)
	}
	return nil
}

// Unregister namespaces by gaeacc
func (nr *NamespaceRegisterManager) unRegisterNamespaces() error {
	for _, namespace := range nr.Namespaces {
		err := nr.unRegisterNamespaceByKey(namespace.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

type GaeaCCManager struct {
	GaeaCCAdminUser     string
	GaeaCCAdminPassword string
	GaeaCCBaseRouter    string
}

func (g *GaeaCCManager) modifyNamespace(n *models.Namespace) error {
	// 将命名空间对象转换为 JSON
	data, err := json.Marshal(n)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", g.GaeaCCBaseRouter+"modify", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (g *GaeaCCManager) deleteNamespace(key string) error {
	router := g.GaeaCCBaseRouter + "delete/" + key
	req, err := http.NewRequest("PUT", router, nil)
	if err != nil {
		return fmt.Errorf("create http request %s error: %w", router, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(g.GaeaCCAdminUser, g.GaeaCCAdminPassword)
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("create http request %s error: %w", router, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// SearchSqlLog function in the provided code is designed to search through SQL log files for specific entries.
// It takes a search string and a timestamp, and returns a slice of log entries that match the criteria.
// The function reads log files within a specified directory, looking for files named with the prefix "gaea_sql.log".
// It uses regular expressions to parse and match log entries based on the input parameters.
// If a matching entry is found, it's added to the result slice. The function handles errors such as file access issues and returns an error if any problems occur during the file reading and parsing process.
func (e *E2eManager) SearchSqlLog(searchString string, currentTime time.Time) ([]util.LogEntry, error) {
	searchString = strings.TrimSuffix(searchString, ";")
	var allEntries []util.LogEntry

	err := filepath.Walk(e.GCluster.LogDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "gaea_sql.log") {
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
