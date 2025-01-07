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

func (e *E2eManager) Clean() error {
	namespaces, err := e.ListNamespace()
	if err != nil {
		return fmt.Errorf("error listing namespaces: %v", err)
	}
	for _, name := range namespaces {
		if err = e.DeleteNamespace(name); err != nil {
			fmt.Printf("Failed to delete namespace '%s': %v\n", name, err)
			return err
		} else {
			fmt.Printf("Successfully deleted namespace '%s'\n", name)
		}
	}

	if e.GCluster != nil {
		if e.GCluster.readWriteConn != nil {
			fmt.Println("Closing readWriteConn connection...")
			if err := e.GCluster.readWriteConn.Close(); err != nil {
				fmt.Println("Failed to close readWriteConn connection:", err)
			} else {
				fmt.Println("Successfully closed readWriteConn connection")
			}
		}

		if e.GCluster.writeConn != nil {
			fmt.Println("Closing writeConn connection...")
			if err := e.GCluster.writeConn.Close(); err != nil {
				fmt.Println("Failed to close writeConn connection:", err)
			} else {
				fmt.Println("Successfully closed writeConn connection")
			}
		}

		if e.GCluster.readConn != nil {
			fmt.Println("Closing readConn connection...")
			if err := e.GCluster.readConn.Close(); err != nil {
				fmt.Println("Failed to close readConn connection:", err)
			} else {
				fmt.Println("Successfully closed readConn connection")
			}
		}
	}

	// 关闭 openConnections 连接
	for _, conn := range e.openConnections {
		if conn != nil {
			fmt.Println("Closing open connection...")
			if err := conn.Close(); err != nil {
				fmt.Println("Failed to close open connection:", err)
			}
		}
	}

	// 清空连接池，避免重复关闭
	e.openConnections = []*sql.DB{}
	fmt.Println("All open connections have been closed and cleared.")

	return nil
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
