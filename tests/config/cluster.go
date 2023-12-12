package config

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/XiaoMi/Gaea/models"
)

type NsSlice struct {
	Name       string
	Slices     []*models.Slice
	SliceConns []*LocalSliceConn
	GaeaUsers  []*models.User
}

func (ns *NsSlice) GetLocalSliceConn() (res map[string]*LocalSliceConn, err error) {
	sliceConns := make(map[string]*LocalSliceConn, len(ns.Slices))
	for _, slice := range ns.Slices {
		localConn := &LocalSliceConn{
			MasterConn: nil,
			SlaveConns: []*sql.DB{},
		}
		// 初始化主服务器连接
		localConn.MasterConn, err = InitConn(slice.UserName, slice.Password, slice.Master, "")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize master connection for slice '%s': %v", slice.Name, err)
		}
		// 初始化从服务器连接
		for _, slave := range slice.Slaves {
			slaveDb, err := InitConn(slice.UserName, slice.Password, slave, "")
			if err != nil {
				return nil, fmt.Errorf("failed to initialize slave connection for slice '%s': %v", slice.Name, err)
			}
			localConn.SlaveConns = append(localConn.SlaveConns, slaveDb)
		}

		sliceConns[slice.Name] = localConn
	}

	return sliceConns, nil
}

// TODO: fix to rename SliceConn
type LocalSliceConn struct {
	MasterConn *sql.DB
	SlaveConns []*sql.DB
}

// TODO: fix error
func (ns *NsSlice) GetMasterConn(sliceIndex int) (*sql.DB, error) {
	var err error
	if sliceIndex > len(ns.Slices) {
		return nil, fmt.Errorf("sliceIndex more than")
	}
	if ns.SliceConns == nil {
		ns.SliceConns = make([]*LocalSliceConn, len(ns.Slices))
	}
	if ns.SliceConns[sliceIndex].MasterConn != nil {
		return ns.SliceConns[sliceIndex].MasterConn, nil
	}
	ns.SliceConns[sliceIndex].MasterConn, err = InitConn(ns.Slices[sliceIndex].UserName, ns.Slices[sliceIndex].Password, ns.Slices[sliceIndex].Master, "")
	return ns.SliceConns[sliceIndex].MasterConn, err
}

// TODO: fix error
func (ns *NsSlice) GetSlaveConn(sliceIndex int, index int) (*sql.DB, error) {
	if sliceIndex > len(ns.Slices) {
		return nil, fmt.Errorf("sliceIndex more than")
	}
	if index > len(ns.Slices[sliceIndex].Slaves) {
		return nil, fmt.Errorf("slaveIndex more than")
	}
	if ns.SliceConns[sliceIndex].SlaveConns[index] != nil {
		return ns.SliceConns[sliceIndex].SlaveConns[index], nil
	}
	ns.SliceConns[sliceIndex].SlaveConns[index], _ = InitConn(ns.Slices[sliceIndex].UserName, ns.Slices[sliceIndex].Password, ns.Slices[sliceIndex].Slaves[index], "")
	return ns.SliceConns[sliceIndex].SlaveConns[index], nil
}

func ParseTemplate(filenames string, ns *NsSlice) (*models.Namespace, error) {
	temp := filepath.Base(filenames)
	tmpl, err := template.New(temp).Funcs(template.FuncMap{
		"lastItem": func(index int, all int) bool {
			return index == all-1
		},
	}).ParseFiles(filenames)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, ns)
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

type MCluster struct {
	Name   string
	Master *MInstance
	Slaves []*MInstance
}

type MInstance struct {
	UserName string
	Password string
	Host     string
	Port     int
}

func (m *MInstance) Addr() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

type GCluster struct {
	Host          string
	Port          int
	ReadWriteUser *models.User
	ReadUser      *models.User
	WriteUser     *models.User
	LogPath       string
	readWriteConn *sql.DB
	readConn      *sql.DB
	writeConn     *sql.DB
}

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(currentFile)
	// 获取上级目录
	basePath = filepath.Dir(configDir)
}

func (e *E2eManager) AddNsFromFile(filenames string, nss *NsSlice) error {
	ns, err := ParseTemplate(filenames, nss)
	if err != nil {
		return err
	}
	return e.AddNamespace(ns)
}

func (e *E2eManager) AddNamespace(ns *models.Namespace) error {
	err := e.NsManager.AddNamespace(ns)
	if err != nil {
		return err
	}
	return nil
}

func (e *E2eManager) Clean() {
	// 删除所有的namespace
	_ = e.NsManager.unRegisterNamespaces()
	if e.GCluster != nil {
		if e.GCluster.readWriteConn != nil {
			_ = e.GCluster.readWriteConn.Close()
		}
		if e.GCluster.writeConn != nil {
			_ = e.GCluster.writeConn.Close()
		}
		if e.GCluster.readConn != nil {
			_ = e.GCluster.readConn.Close()
		}
	}

	for _, v := range e.NsSlices {
		for _, sliceConn := range v.SliceConns {
			if sliceConn.MasterConn != nil {
				_ = sliceConn.MasterConn.Close()
			}

			for _, slaveConn := range sliceConn.SlaveConns {
				if slaveConn != nil {
					_ = slaveConn.Close()
				}
			}
		}
	}
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
