package config

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/XiaoMi/Gaea/tests/util"

	"github.com/XiaoMi/Gaea/models"
)

const (
	// 注意这里的端口要和gaea_cc.ini 中的addr=0.0.0.0:23306
	gaeaCCBaseRouter    = "http://localhost:23306/api/cc/namespace/"
	gaeaCCAdminUsername = "admin"
	gaeaCCAdminPassword = "admin"
	defaultMysqlUser    = "superroot"
	defaultMysqlPasswd  = "superroot"
	defaultHost         = "127.0.0.1"

	DefaultE2eDatabase = "db_e2e_test"
	DefaultE2eTable    = "tbl_e2e_test"
	SliceSMName        = "single-master"
	SliceMMName        = "multi-masters"
	SliceMSSName       = "master-single-slave"
	SliceMSName        = "master-slaves"
	LogExpression      = `\[(.*?)\] \[NOTICE\] \[(\d+)\] OK - (\d+\.\d+)ms - ns=(.*?), (.*?)@(.*?)->(.*?)/(.*?), mysql_connect_id=(\d+), r=\d+\|(.*?)$`
)

var logPath = "cmd/logs/gaea_sql.log"
var E2eMgr *E2eManager

type E2eManager struct {
	// 用于管理所有的测试用例
	NsManager            *NamespaceRegisterManager
	MClusterMaster       *MCluster
	MClusterMasterSlaves *MCluster
	GCluster             *GCluster
	NsSlices             map[string]*NsSlice
	BasePath             string
	StartTime            string
}

func NewE2eManager() *E2eManager {
	mClusterMaster := &MCluster{
		Name:   SliceSMName,
		Master: &MInstance{defaultMysqlUser, defaultMysqlPasswd, defaultHost, 3379},
	}
	mClusterMasterSlaves := &MCluster{
		Name:   SliceMSName,
		Master: &MInstance{defaultMysqlUser, defaultMysqlPasswd, defaultHost, 3319},
		Slaves: []*MInstance{{defaultMysqlUser, defaultMysqlPasswd, defaultHost, 3329},
			{defaultMysqlUser, defaultMysqlPasswd, defaultHost, 3339}},
	}
	GaeaUsers := []*models.User{
		{
			UserName: defaultMysqlUser,
			Password: defaultMysqlPasswd,
			RWFlag:   2,
			RWSplit:  1,
		},
		{
			UserName: defaultMysqlUser + "_r",
			Password: defaultMysqlPasswd + "_r",
			RWFlag:   1,
			RWSplit:  0,
		},
		{
			UserName: defaultMysqlUser + "_w",
			Password: defaultMysqlPasswd + "_w",
			RWFlag:   2,
			RWSplit:  0,
		},
	}

	// 3379
	sliceSingleMaster := &NsSlice{
		Name: SliceSMName,
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
		Name: SliceMSName,
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
		Name: SliceMSSName,
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
		Name: SliceMMName,
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
		GCluster: &GCluster{
			Host:          defaultHost,
			Port:          13306,
			ReadWriteUser: GaeaUsers[0],
			ReadUser:      GaeaUsers[1],
			WriteUser:     GaeaUsers[2],
			LogPath:       filepath.Join(basePath, logPath),
		},
		NsSlices:  nsSlices,
		BasePath:  basePath,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	return E2eMgr
}

func (e *E2eManager) GetReadWriteGaeaUserConn() (*sql.DB, error) {
	if e.GCluster.readWriteConn != nil {
		return e.GCluster.readWriteConn, nil
	}
	return InitConn(e.GCluster.ReadWriteUser.UserName, e.GCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) NewReadWriteGaeaUserConn() (*sql.DB, error) {
	return InitConn(e.GCluster.ReadWriteUser.UserName, e.GCluster.ReadWriteUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) GetReadGaeaUserConn() (*sql.DB, error) {
	if e.GCluster.readConn != nil {
		return e.GCluster.readConn, nil
	}
	return InitConn(e.GCluster.ReadUser.UserName, e.GCluster.ReadUser.Password, fmt.Sprintf("%s:%d", e.GCluster.Host, e.GCluster.Port), "")
}

func (e *E2eManager) GetWriteGaeaUserConn() (*sql.DB, error) {
	if e.GCluster.writeConn != nil {
		return e.GCluster.writeConn, nil
	}
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

func (nr *NamespaceRegisterManager) AddNamespace(m *models.Namespace) error {
	// 存在就替换
	err := m.Verify()
	if err != nil {
		return err
	}
	if _, ok := nr.Namespaces[m.Name]; ok {
		return fmt.Errorf("namespace has exists")
	}
	nr.Namespaces[m.Name] = m
	err = nr.GaeaCCManager.modifyNamespace(m)
	if err != nil {
		return fmt.Errorf("modify Namespace error:%v", err)
	}
	return nil
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
		Timeout: time.Second * 5,
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

func (mgr *E2eManager) SearchLog(searchString string, currentTime time.Time) ([]util.LogEntry, error) {
	// 打开文件
	file, err := os.Open(mgr.GCluster.LogPath)
	if err != nil {
		return []util.LogEntry{}, fmt.Errorf("open file:%s error %v ", mgr.GCluster.LogPath, err)
	}
	defer file.Close()
	// 正则表达式
	re := regexp.MustCompile(LogExpression)
	return util.ParseLogEntries(file, re, currentTime, searchString)
}
