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
	fileclient "github.com/XiaoMi/Gaea/models/file"
)

// config type
const (
	ConfigFile = "file"
	ConfigEtcd = "etcd"
)

// Client client interface
type Client interface {
	Create(path string, data []byte) error
	Update(path string, data []byte) error
	UpdateWithTTL(path string, data []byte, ttl time.Duration) error
	Delete(path string) error
	Read(path string) ([]byte, error)
	List(path string) ([]string, error)
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
		c, err := fileclient.New(root)
		if err != nil {
			log.Warn("create fileclient failed, %s", addr)
			return nil
		}
		return c
	case ConfigEtcd:
		// etcd
		c, err := etcdclient.New(addr, time.Minute, username, password, root)
		if err != nil {
			log.Fatal("create etcdclient to %s failed, %v", addr, err)
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
func (s *Store) NamespaceBase() string {
	return filepath.Join(s.prefix, "namespace")
}

// NamespacePath concat namespace path
func (s *Store) NamespacePath(name string) string {
	return filepath.Join(s.prefix, "namespace", name)
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
	return s.client.Update(s.ProxyPath(p.Token), p.Encode())
}

// DeleteProxy delete proxy path
func (s *Store) DeleteProxy(token string) error {
	return s.client.Delete(s.ProxyPath(token))
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
		return nil, fmt.Errorf("node %s not exists", s.NamespacePath(name))
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
