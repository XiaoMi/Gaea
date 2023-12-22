package config

import (
	"bytes"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

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
	// SliceMaster 表示测试的单主 MySQL 集群
	SliceMaster = "single-master"
	// SliceMultiMaster 表示测试的多主 MySQL 集群
	SliceMultiMaster = "multi-masters"
	// SliceMasterSlave 表示测试的主从 MySQL 集群
	SliceMasterSlave = "master-single-slave"
	// SliceMasterSlaves 表示测试的主从 MySQL 集群
	SliceMasterSlaves = "master-slaves"
	// LogExpression 标识 Gaea SQL Log 的格式
	LogExpression = `\[(.*?)\] \[NOTICE\] \[(\d+)\] OK - (\d+\.\d+)ms - ns=(.*?), (.*?)@(.*?)->(.*?)/(.*?), mysql_connect_id=(\d+), r=\d+\|(.*?)$`
)

var logPath = "cmd/logs/gaea_sql.log"
var E2eMgr *E2eManager

//go:embed ns/default.template
var DefaultNamespaceTmpl string

type E2eManager struct {
	// 用于管理所有的测试用例
	NsManager            *NamespaceRegisterManager
	MClusterMaster       *MySqlCluster
	MClusterMasterSlaves *MySqlCluster
	GCluster             *GaeaCluster
	NsSlices             map[string]*NsSlice
	BasePath             string
	StartTime            time.Time
	Db                   string
	Table                string
}

func NewE2eManager() *E2eManager {
	mClusterMaster := &MySqlCluster{
		Type:   SliceMaster,
		Master: &MysqlInstance{defaultGaeaBackendUser, defaultGaeaBackendPass, defaultHost, 3379},
	}
	mClusterMasterSlaves := &MySqlCluster{
		Type:   SliceMasterSlaves,
		Master: &MysqlInstance{defaultGaeaBackendUser, defaultGaeaBackendPass, defaultHost, 3319},
		Slaves: []*MysqlInstance{{defaultGaeaBackendUser, defaultGaeaBackendPass, defaultHost, 3329},
			{defaultGaeaBackendUser, defaultGaeaBackendPass, defaultHost, 3339}},
	}
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
	sliceSingleMaster := &NsSlice{
		Name: SliceMaster,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        mClusterMaster.Master.UserName,
				Password:        mClusterMaster.Master.Password,
				Master:          mClusterMaster.Master.Addr(),
				Slaves:          nil,
				StatisticSlaves: nil,
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		SliceConns: []*LocalSliceConn{
			{
				MasterConn: nil,
				SlaveConns: nil,
			},
		},
		GaeaUsers: GaeaUsers,
	}
	// 3319 3329 3339
	sliceMasterSlaves := &NsSlice{
		Name: SliceMasterSlaves,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        mClusterMasterSlaves.Master.UserName,
				Password:        mClusterMasterSlaves.Master.Password,
				Master:          mClusterMasterSlaves.Master.Addr(),
				Slaves:          []string{mClusterMasterSlaves.Slaves[0].Addr(), mClusterMasterSlaves.Slaves[1].Addr()},
				StatisticSlaves: nil,
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		SliceConns: []*LocalSliceConn{
			{
				MasterConn: nil,
				SlaveConns: []*sql.DB{nil, nil},
			},
		},
		GaeaUsers: GaeaUsers,
	}

	// 3319 3329
	sliceMasterSingleSlave := &NsSlice{
		Name: SliceMasterSlave,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        mClusterMasterSlaves.Master.UserName,
				Password:        mClusterMasterSlaves.Master.Password,
				Master:          mClusterMasterSlaves.Master.Addr(),
				Slaves:          []string{mClusterMasterSlaves.Slaves[0].Addr()},
				StatisticSlaves: nil,
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		SliceConns: []*LocalSliceConn{
			{
				MasterConn: nil,
				SlaveConns: []*sql.DB{nil},
			},
		},
		GaeaUsers: GaeaUsers,
	}
	// 3329 3339
	sliceMultiMasters := &NsSlice{
		Name: SliceMultiMaster,
		Slices: []*models.Slice{
			{
				Name:            "slice-0",
				UserName:        mClusterMaster.Master.UserName,
				Password:        mClusterMaster.Master.Password,
				Master:          mClusterMasterSlaves.Slaves[0].Addr(),
				Slaves:          nil,
				StatisticSlaves: nil,
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
			{
				Name:            "slice-1",
				UserName:        mClusterMaster.Master.UserName,
				Password:        mClusterMaster.Master.Password,
				Master:          mClusterMasterSlaves.Slaves[1].Addr(),
				Slaves:          nil,
				StatisticSlaves: nil,
				Capacity:        12,
				MaxCapacity:     24,
				IdleTimeout:     60,
			},
		},
		SliceConns: []*LocalSliceConn{
			{
				MasterConn: nil,
				SlaveConns: nil,
			},
			{
				MasterConn: nil,
				SlaveConns: nil,
			},
		},
		GaeaUsers: GaeaUsers,
	}
	nsSlices := map[string]*NsSlice{
		//3379
		sliceSingleMaster.Name: sliceSingleMaster,
		//3319 3329 3339
		sliceMasterSlaves.Name: sliceMasterSlaves,
		//3329 3339
		sliceMultiMasters.Name: sliceMultiMasters,
		//3319 3329
		sliceMasterSingleSlave.Name: sliceMasterSingleSlave,
	}
	E2eMgr = &E2eManager{
		NsManager:            NewNamespaceRegisterManger(),
		MClusterMaster:       mClusterMaster,
		MClusterMasterSlaves: mClusterMasterSlaves,
		GCluster: &GaeaCluster{
			Host:          defaultHost,
			Port:          13306,
			ReadWriteUser: GaeaUsers[0],
			ReadUser:      GaeaUsers[1],
			WriteUser:     GaeaUsers[2],
			LogPath:       filepath.Join(basePath, logPath),
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
	//if e.GaeaCluster.readWriteConn != nil {
	//	return e.GaeaCluster.readWriteConn, nil
	//}
	//var err error
	//e.GaeaCluster.readWriteConn, err = InitConn(e.GaeaCluster.ReadWriteUser.UserName, e.GaeaCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GaeaCluster.Host, e.GaeaCluster.Port), "")
	//return e.GaeaCluster.readWriteConn, err
	return InitConn(e.GCluster.ReadWriteUser.UserName, e.GCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) GetReadGaeaUserConn() (*sql.DB, error) {
	//if e.GaeaCluster.readConn != nil {
	//	return e.GaeaCluster.readConn, nil
	//}
	//var err error
	//e.GaeaCluster.readConn, err = InitConn(e.GaeaCluster.ReadUser.UserName, e.GaeaCluster.ReadUser.Password, fmt.Sprintf("%s:%d", e.GaeaCluster.Host, e.GaeaCluster.Port), "")
	//return e.GaeaCluster.readConn, err
	return InitConn(e.GCluster.ReadUser.UserName, e.GCluster.ReadUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) GetWriteGaeaUserConn() (*sql.DB, error) {
	//if e.GaeaCluster.writeConn != nil {
	//	return e.GaeaCluster.writeConn, nil
	//}
	//var err error
	//e.GaeaCluster.writeConn, err = InitConn(e.GaeaCluster.WriteUser.UserName, e.GaeaCluster.WriteUser.Password, fmt.Sprintf("%s:%d", e.GaeaCluster.Host, e.GaeaCluster.Port), "")
	//return e.GaeaCluster.writeConn, err
	return InitConn(e.GCluster.WriteUser.UserName, e.GCluster.WriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
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

func (e *E2eManager) SearchLog(searchString string, currentTime time.Time) ([]util.LogEntry, error) {
	// 打开文件
	file, err := os.Open(e.GCluster.LogPath)
	if err != nil {
		return []util.LogEntry{}, fmt.Errorf("open file:%s error %v ", e.GCluster.LogPath, err)
	}
	defer file.Close()
	// 正则表达式
	re := regexp.MustCompile(LogExpression)
	return util.ParseLogEntries(file, re, currentTime, searchString)
}
