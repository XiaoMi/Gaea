// Copyright 2019 The Gaea Authors. All Rights Reserved.
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

package models

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/log"
	etcdclient "github.com/XiaoMi/Gaea/models/etcd"
	etcdclientv3 "github.com/XiaoMi/Gaea/models/etcdv3"
	fileclient "github.com/XiaoMi/Gaea/models/file"
	"github.com/coreos/etcd/client"
)

// config type
const (
	ConfigFile   = "file"
	ConfigEtcd   = "etcd"
	ConfigEtcdV3 = "etcdv3"
)

// Client client interface
type Client interface {
	Create(path string, data []byte) error
	Update(path string, data []byte) error
	UpdateWithTTL(path string, data []byte, ttl time.Duration) error
	Delete(path string) error
	Read(path string) ([]byte, error)
	List(path string) ([]string, error)
	ListWithValues(path string) (map[string]string, error)
	Close() error
	BasePrefix() string
}

// Store means exported client to use
type Store struct {
	client Client
	prefix string
}

// NewClient constructor to create client by case etcd/file/zk etc.
func NewClient(configType, addr, username, password, root string) Client {
	switch configType {
	case ConfigFile:
		// 使用文档 File 去读取设定值
		c, err := fileclient.New(root)
		if err != nil {
			log.Warn("create fileclient failed, %s", addr)
			return nil
		}
		return c
	case ConfigEtcd:
		// 使用 Etcd V2 API 去读取设定值
		c, err := etcdclient.New(addr, time.Minute, username, password, root)
		if err != nil {
			log.Fatal("create etcdclient v2 to %s failed, %v", addr, err)
			return nil
		}
		return c
	case ConfigEtcdV3:
		// 使用 Etcd V3 API 去读取设定值
		c, err := etcdclientv3.New(addr, time.Minute, username, password, root)
		if err != nil {
			log.Fatal("create etcdclient v3 to %s failed, %v", addr, err)
			return nil
		}
		return c
	}
	log.Fatal("unknown config type")
	return nil
}

// NewStore constructor of Store
func NewStore(client Client) *Store {
	return &Store{
		client: client,
		prefix: client.BasePrefix(),
	}
}

// Close close store
func (s *Store) Close() error {
	return s.client.Close()
}

// NamespaceBase return namespace path base
func (s *Store) GaeaProxyBase() string {
	return filepath.Join(s.prefix, "proxy")
}

// NamespaceBase return namespace path base
func (s *Store) NamespaceBase() string {
	return filepath.Join(s.prefix, "namespace")
}

// NamespacePath concat namespace path
func (s *Store) NamespacePath(name string) string {
	return filepath.Join(s.prefix, "namespace", name)
}

// GaeaNodePath concat gaea node path by ip and port
func (s *Store) GaeaNodePath(ip string, port string) string {
	return filepath.Join(s.prefix, "proxy", "proxy-"+ip+":"+port)
}

// ProxyBase return proxy path base
func (s *Store) ProxyBase() string {
	return filepath.Join(s.prefix, "proxy")
}

// ProxyPath concat proxy path
func (s *Store) ProxyPath(token string) string {
	return filepath.Join(s.prefix, "proxy", fmt.Sprintf("proxy-%s", token))
}

// CreateProxy create proxy model
func (s *Store) CreateProxy(p *ProxyInfo) error {
	// 就在这里，会在 etcd 服务器上新增一个 key /gaea_cluster/proxy/proxy-127.0.0.1:13306
	return s.client.Update(s.ProxyPath(p.Token), p.Encode())
}

// DeleteProxy delete proxy path
func (s *Store) DeleteProxy(token string) error {
	return s.client.Delete(s.ProxyPath(token))
}

// ListNamespace list namespace
func (s *Store) ListGaeaNode() ([]string, error) {
	files, err := s.client.List(s.GaeaProxyBase())
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(files); i++ {
		tmp := strings.Split(files[i], "/")
		files[i] = tmp[len(tmp)-1]
	}
	return files, nil
}

// ListNamespace list namespace
func (s *Store) ListNamespace() ([]string, error) {
	files, err := s.client.List(s.NamespaceBase())
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(files); i++ {
		tmp := strings.Split(files[i], "/")
		files[i] = tmp[len(tmp)-1]
	}
	return files, nil
}

// LoadNamespace load namespace value
func (s *Store) LoadNamespace(key, name string) (*Namespace, error) {
	b, err := s.client.Read(s.NamespacePath(name))
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, client.Error{
			Code:    client.ErrorCodeKeyNotFound,
			Message: fmt.Sprintf("node %s not exists", s.NamespacePath(name)),
		}
	}
	p := &Namespace{}
	if err = json.Unmarshal(b, p); err != nil {
		return nil, err
	}

	if err = p.Verify(); err != nil {
		return nil, err
	}

	if err = p.Decrypt(key); err != nil {
		return nil, err
	}

	return p, nil
}

// UpdateNamespace update namespace path with data
func (s *Store) UpdateNamespace(p *Namespace) error {
	return s.client.Update(s.NamespacePath(p.Name), p.Encode())
}

// DelNamespace delete namespace
func (s *Store) DelNamespace(name string) error {
	return s.client.Delete(s.NamespacePath(name))
}

// DelGaeaPort delete gaea node
func (s *Store) DelGaeaPort(ip string, port string) error {
	if _, err := s.client.List(s.GaeaNodePath(ip, port)); err != nil {
		return err
	}
	return s.client.Delete(s.GaeaNodePath(ip, port))
}

// ListProxyMonitorMetrics list proxies in proxy register path
func (s *Store) ListProxyMonitorMetrics() (map[string]*ProxyMonitorMetric, error) {
	files, err := s.client.List(s.ProxyBase())
	if err != nil {
		return nil, err
	}
	proxy := make(map[string]*ProxyMonitorMetric)
	for _, path := range files {
		b, err := s.client.Read(path)
		if err != nil {
			return nil, err
		}
		p := &ProxyMonitorMetric{}
		if err := JSONDecode(p, b); err != nil {
			return nil, err
		}
		proxy[p.Token] = p
	}
	return proxy, nil
}

func (s *Store) ListNamespaces() (map[string]*Namespace, error) {
	values, err := s.client.ListWithValues(s.NamespaceBase())
	if err != nil {
		return nil, err
	}
	// 初始化结果map
	res := make(map[string]*Namespace, len(values))
	for key, value := range values {
		ns := &Namespace{}
		// 反序列化每个值到Namespace结构体
		if err := json.Unmarshal([]byte(value), ns); err != nil {
			// 如果出现错误，处理错误，例如记录或返回错误
			return nil, err
		}
		// 将反序列化的Namespace添加到结果map
		res[key] = ns
	}
	return res, nil
}
