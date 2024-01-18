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
	Name      string
	Slices    []*models.Slice
	GaeaUsers []*models.User
}

func (ns *NsSlice) GetLocalSliceConn() (res map[string]*LocalSliceConn, err error) {
	sliceConns := make(map[string]*LocalSliceConn, len(ns.Slices))
	for _, slice := range ns.Slices {
		localConn := &LocalSliceConn{
			MasterConn: nil,
			SlaveConns: []*sql.DB{},
		}
		// 初始化主服务器连接
		// TODO: fix hard user and password
		localConn.MasterConn, err = InitConn(defaultMysqlAdminUser, defaultMysqlAdminPasswd, slice.Master, "")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize master connection for slice '%s': %v", slice.Name, err)
		}
		// 初始化从服务器连接
		for _, slave := range slice.Slaves {
			slaveDb, err := InitConn(defaultMysqlAdminUser, defaultMysqlAdminPasswd, slave, "")
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
func (ns *NsSlice) GetMasterAdminConn(sliceIndex int) (*sql.DB, error) {
	return InitConn(defaultMysqlAdminUser, defaultMysqlAdminPasswd, ns.Slices[sliceIndex].Master, "")
}

func (ns *NsSlice) GetMasterAdminDBConn(sliceIndex int, db string) (*sql.DB, error) {
	return InitConn(defaultMysqlAdminUser, defaultMysqlAdminPasswd, ns.Slices[sliceIndex].Master, db)
}

func (ns *NsSlice) GetMasterCommonConn(sliceIndex int) (*sql.DB, error) {
	return InitConn(defaultGaeaBackendUser, defaultGaeaBackendPass, ns.Slices[sliceIndex].Master, "")
}

func (ns *NsSlice) GetMasterCommonDBConn(sliceIndex int, db string) (*sql.DB, error) {
	return InitConn(defaultGaeaBackendUser, defaultGaeaBackendPass, ns.Slices[sliceIndex].Master, db)
}

// TODO: fix error
func (ns *NsSlice) GetSlaveAdminConn(sliceIndex int, index int) (*sql.DB, error) {
	if sliceIndex > len(ns.Slices) {
		return nil, fmt.Errorf("sliceIndex more than")
	}
	if index > len(ns.Slices[sliceIndex].Slaves) {
		return nil, fmt.Errorf("slaveIndex more than")
	}

	return InitConn(defaultMysqlAdminUser, defaultMysqlAdminPasswd, ns.Slices[sliceIndex].Slaves[index], "")
}

func (ns *NsSlice) GetSlaveCommonConn(sliceIndex int, index int) (*sql.DB, error) {
	if sliceIndex > len(ns.Slices) {
		return nil, fmt.Errorf("sliceIndex more than")
	}
	if index > len(ns.Slices[sliceIndex].Slaves) {
		return nil, fmt.Errorf("slaveIndex more than")
	}
	return InitConn(defaultGaeaBackendUser, defaultGaeaBackendPass, ns.Slices[sliceIndex].Slaves[index], "")
}

func ParseNamespaceTmpl(nsTmpl string, ns *NsSlice) (*models.Namespace, error) {
	tmpl, err := template.New(nsTmpl).Funcs(template.FuncMap{
		"lastItem": func(index int, all int) bool {
			return index == all-1
		},
	}).Parse(nsTmpl)
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

type MySqlCluster struct {
	Type   string
	Master *MysqlInstance
	Slaves []*MysqlInstance
}

type MysqlInstance struct {
	UserName string
	Password string
	Host     string
	Port     int
}

func (m *MysqlInstance) Addr() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

type GaeaCluster struct {
	Host          string
	Port          int
	ReadWriteUser *models.User
	ReadUser      *models.User
	WriteUser     *models.User
	LogDirectory  string
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
	fmt.Printf("basePath: %s\n", basePath)
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
