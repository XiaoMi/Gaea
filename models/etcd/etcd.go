// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

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

package etcdclient

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/client"

	"github.com/XiaoMi/Gaea/log"
)

// ErrClosedEtcdClient means etcd client closed
var ErrClosedEtcdClient = errors.New("use of closed etcd client")

const (
	defaultEtcdPrefix = "/gaea"
)

// EtcdClient etcd client
type EtcdClient struct {
	sync.Mutex
	kapi client.KeysAPI

	closed  bool
	timeout time.Duration
	Prefix  string
}

// New constructor of EtcdClient
func New(addr string, timeout time.Duration, username, passwd, root string) (*EtcdClient, error) {
	endpoints := strings.Split(addr, ",")
	for i, s := range endpoints {
		if s != "" && !strings.HasPrefix(s, "http://") {
			endpoints[i] = "http://" + s
		}
	}
	config := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		Username:                username,
		Password:                passwd,
		HeaderTimeoutPerRequest: time.Second * 10,
	}
	c, err := client.New(config)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(root) == "" {
		root = defaultEtcdPrefix
	}
	return &EtcdClient{
		kapi:    client.NewKeysAPI(c),
		timeout: timeout,
		Prefix:  root,
	}, nil
}

// Close close etcd client
func (c *EtcdClient) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	return nil
}

func (c *EtcdClient) contextWithTimeout() (context.Context, context.CancelFunc) {
	if c.timeout == 0 {
		return context.Background(), func() {}
	}
	return context.WithTimeout(context.Background(), c.timeout)
}

func isErrNoNode(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			return e.Code == client.ErrorCodeKeyNotFound
		}
	}
	return false
}

func isErrNodeExists(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			return e.Code == client.ErrorCodeNodeExist
		}
	}
	return false
}

// Mkdir create directory
func (c *EtcdClient) Mkdir(dir string) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	return c.mkdir(dir)
}

func (c *EtcdClient) mkdir(dir string) error {
	if dir == "" || dir == "/" {
		return nil
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_, err := c.kapi.Set(cntx, dir, "", &client.SetOptions{Dir: true, PrevExist: client.PrevNoExist})
	if err != nil {
		if isErrNodeExists(err) {
			return nil
		}
		return err
	}
	return nil
}

// Create create path with data
func (c *EtcdClient) Create(path string, data []byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	log.Debug("etcd create node %s", path)
	_, err := c.kapi.Set(cntx, path, string(data), &client.SetOptions{PrevExist: client.PrevNoExist})
	if err != nil {
		log.Debug("etcd create node %s failed: %s", path, err)
		return err
	}
	log.Debug("etcd create node OK")
	return nil
}

// Update update path with data
func (c *EtcdClient) Update(path string, data []byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	log.Debug("etcd update node %s", path)
	_, err := c.kapi.Set(cntx, path, string(data), &client.SetOptions{PrevExist: client.PrevIgnore})
	if err != nil {
		log.Debug("etcd update node %s failed: %s", path, err)
		return err
	}
	log.Debug("etcd update node OK")
	return nil
}

// UpdateWithTTL update path with data and ttl
func (c *EtcdClient) UpdateWithTTL(path string, data []byte, ttl time.Duration) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	log.Debug("etcd update node %s with ttl %d", path, ttl)
	_, err := c.kapi.Set(cntx, path, string(data), &client.SetOptions{PrevExist: client.PrevIgnore, TTL: ttl})
	if err != nil {
		log.Debug("etcd update node %s failed: %s", path, err)
		return err
	}
	log.Debug("etcd update node OK")
	return nil
}

// Delete delete path
func (c *EtcdClient) Delete(path string) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	log.Debug("etcd delete node %s", path)
	_, err := c.kapi.Delete(cntx, path, nil)
	if err != nil && !isErrNoNode(err) {
		log.Debug("etcd delete node %s failed: %s", path, err)
		return err
	}
	log.Debug("etcd delete node OK")
	return nil
}

// Read read path data
func (c *EtcdClient) Read(path string) ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil, ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	log.Debug("etcd read node %s", path)
	r, err := c.kapi.Get(cntx, path, nil)
	if err != nil && !isErrNoNode(err) {
		return nil, err
	} else if r == nil || r.Node.Dir {
		return nil, nil
	} else {
		return []byte(r.Node.Value), nil
	}
}

// List list path, return slice of all paths
func (c *EtcdClient) List(path string) ([]string, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil, ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	log.Debug("etcd list node %s", path)
	r, err := c.kapi.Get(cntx, path, nil)
	if err != nil && !isErrNoNode(err) {
		return nil, err
	} else if r == nil || !r.Node.Dir {
		return nil, nil
	} else {
		var files []string
		for _, node := range r.Node.Nodes {
			files = append(files, node.Key)
		}
		return files, nil
	}
}

// Watch watch path
func (c *EtcdClient) Watch(path string, ch chan string) error {
    c.Lock() // 在这里上锁
	// defer c.Unlock() // 移除此行，避免死结发生
	if c.closed {
		c.Unlock() // 上锁后记得解锁，去防止死结问题发生
		panic(ErrClosedEtcdClient)
	}
	
	watcher := c.kapi.Watcher(path, &client.WatcherOptions{Recursive: true})
    
    c.Unlock() // 上锁后在适当时机解锁，去防止死结问题发生
	// 在这里解锁是最好的，因为解锁后立刻可以进行监听
    
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			panic(err)
		}
		ch <- res.Action
	}
}

// BasePrefix return base prefix
func (c *EtcdClient) BasePrefix() string {
	return c.Prefix
}
