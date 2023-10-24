package config

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sync"
	"text/template"

	"github.com/XiaoMi/Gaea/models"
)

var (
	singleMasterClusterSlice = []*models.Slice{
		{

			Name:            "slice-0",
			UserName:        "superroot",
			Password:        "superroot",
			Master:          "127.0.0.1:3379",
			Slaves:          nil,
			StatisticSlaves: nil,
			Capacity:        12,
			MaxCapacity:     24,
			IdleTimeout:     60,
		},
	}
	masterSlaveClusterSlice = []*models.Slice{
		{

			Name:            "slice-0",
			UserName:        "superroot",
			Password:        "superroot",
			Master:          "127.0.0.1:3319",
			Slaves:          []string{"127.0.0.1:3329", "127.0.0.1:3339"},
			StatisticSlaves: nil,
			Capacity:        12,
			MaxCapacity:     24,
			IdleTimeout:     60,
		},
	}
	multiMasterClusterSlice = []*models.Slice{
		{

			Name:            "slice-0",
			UserName:        "superroot",
			Password:        "superroot",
			Master:          "127.0.0.1:3329",
			Slaves:          nil,
			StatisticSlaves: nil,
			Capacity:        12,
			MaxCapacity:     24,
			IdleTimeout:     60,
		},
		{

			Name:            "slice-1",
			UserName:        "superroot",
			Password:        "superroot",
			Master:          "127.0.0.1:3339",
			Slaves:          nil,
			StatisticSlaves: nil,
			Capacity:        12,
			MaxCapacity:     24,
			IdleTimeout:     60,
		},
	}
	gaeaUser = &models.User{
		UserName:      "superroot",
		Password:      "superroot",
		RWFlag:        2,
		RWSplit:       1,
		OtherProperty: 0,
	}
)

type E2eConfig struct {
	SingleMasterCluster *MysqlClusterConfig
	MultiMasterCluster  *MysqlClusterConfig
	MasterSlaveCluster  *MysqlClusterConfig
	NamespaceManager    *NamespaceRegisterManager
	GaeaUser            *models.User
	GaeaHost            string
	TestsBasePath       string
}

var (
	once      sync.Once
	e2eConfig *E2eConfig
)

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(currentFile)
	// 获取上级目录
	basePath = filepath.Dir(configDir)
}

func GetDefaultE2eConfig() *E2eConfig {
	once.Do(func() {
		e2eConfig = &E2eConfig{
			SingleMasterCluster: &MysqlClusterConfig{
				Slices:     singleMasterClusterSlice,
				SliceConns: map[string]*SliceConn{},
			},
			MasterSlaveCluster: &MysqlClusterConfig{
				Slices:     masterSlaveClusterSlice,
				SliceConns: map[string]*SliceConn{},
			},
			MultiMasterCluster: &MysqlClusterConfig{
				Slices:     multiMasterClusterSlice,
				SliceConns: map[string]*SliceConn{},
			},
			NamespaceManager: NewNamespaceRegisterManger(),
			GaeaUser:         gaeaUser,
			GaeaHost:         "127.0.0.1:13306",
			TestsBasePath:    basePath,
		}
	})
	return e2eConfig
}

func (e *E2eConfig) FilepathJoin(path string) string {
	// 拼接路径
	return filepath.Join(e2eConfig.TestsBasePath, path)
}

func (e *E2eConfig) RegisterNamespaces(ns ...*models.Namespace) error {
	for _, v := range ns {
		err := e.NamespaceManager.AddNamespace(v)
		if err != nil {
			return err
		}
	}
	err := e.NamespaceManager.RegisterAllNamespaces()
	if err != nil {
		return err
	}
	return nil
}

func (e *E2eConfig) UnRegisterNamespaces() error {
	err := e.NamespaceManager.UnRegisterNamespaces()
	if err != nil {
		return err
	}
	return nil
}

type MysqlClusterConfig struct {
	Slices     []*models.Slice
	SliceConns map[string]*SliceConn
}

func (m *MysqlClusterConfig) InitMysqlClusterConn() error {
	m.SliceConns = make(map[string]*SliceConn)

	for _, slice := range m.Slices {
		master, err := InitConn(slice.UserName, slice.Password, slice.Master, "")
		if err != nil {
			return err
		}
		sliceConn := &SliceConn{
			Master: master,
			Slaves: []*sql.DB{},
		}

		for _, slave := range slice.Slaves {
			slaveDB, err := InitConn(slice.UserName, slice.Password, slave, "")
			if err != nil {
				return err
			}
			sliceConn.Slaves = append(sliceConn.Slaves, slaveDB)
		}
		m.SliceConns[slice.Name] = sliceConn
	}
	return nil
}

func getLastItem() template.FuncMap {
	return template.FuncMap{
		"lastItem": func(index int, all int) bool {
			return index == all-1
		},
	}
}

func (m *MysqlClusterConfig) TemplateParse(filenames string) (*models.Namespace, error) {

	temp := filepath.Base(filenames)
	tmpl, err := template.New(temp).Funcs(getLastItem()).ParseFiles(filenames)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, m)
	if err != nil {
		return nil, err
	}

	var nameSpace = &models.Namespace{}
	err = json.NewDecoder(&buf).Decode(&nameSpace)
	if err != nil {
		return nil, err
	}
	return nameSpace, nil
}

func GetJSONFilesFromDir(dir string) ([]string, error) {
	var casesPath []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			absPath := filepath.Join(dir, f.Name())
			casesPath = append(casesPath, absPath)
		}
	}

	return casesPath, nil
}
