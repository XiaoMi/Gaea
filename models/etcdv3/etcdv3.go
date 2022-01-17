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

package etcdclientv3

import (
	"context"
	"errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"sync"
	"time"
)

// ErrClosedEtcdClient means etcd client closed
var ErrClosedEtcdClient = errors.New("use of closed etcd client")

const (
	defaultEtcdPrefix = "/gaea"
)

// EtcdClientV3 为新版的 etcd client
type EtcdClientV3 struct {
	sync.Mutex
	kapi clientv3.Client // 就这里改成新版的 API

	closed  bool
	timeout time.Duration
	Prefix  string
}

// New 建立新的 etcd v3 的客户端
func New(addr string, timeout time.Duration, username, passwd, root string) (*EtcdClientV3, error) {
	endpoints := strings.Split(addr, ",")
	for i, s := range endpoints {
		if s != "" && !strings.HasPrefix(s, "http://") {
			endpoints[i] = "http://" + s
		}
	}
	config := clientv3.Config{
		Endpoints:            endpoints,
		Username:             username,
		Password:             passwd,
		DialTimeout:          timeout, // 只设定第一次连线时间的逾时，之后不用太担心连线，连线失败后，会自动重连
		DialKeepAliveTimeout: timeout, // 之后维持 etcd 连线的逾时
	}
	c, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(root) == "" {
		root = defaultEtcdPrefix
	}
	return &EtcdClientV3{
		kapi:    *c,
		timeout: timeout, // 每个指令的逾时时间也是用同一个设定值
		Prefix:  root,
	}, nil
}

// Close 关闭新的 etcd v3 的客户端
func (c *EtcdClientV3) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil
	}
	err := c.kapi.Close() // 如果能够成功关闭 etcd 连线
	if err != nil {
		return err
	}
	c.closed = true // 如果没有错误，就特别标记连线已正确关闭
	return nil
}

func (c *EtcdClientV3) contextWithTimeout() (context.Context, context.CancelFunc) {
	if c.timeout == 0 {
		return context.Background(), func() {}
	}
	return context.WithTimeout(context.Background(), c.timeout)
}

// isErrNodeExists (目前此函式并不支援 v3)
/*func isErrNoNode(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			return e.Code == client.ErrorCodeKeyNotFound
		}
	}
	return false
}*/

// isErrNodeExists (v3 版没有在使用此函式)
/*func isErrNodeExists(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			return e.Code == client.ErrorCodeNodeExist
		}
	}
	return false
}*/

// Mkdir create directory (v2 的版本也没在用这函式)
/*func (c *EtcdClientV3) Mkdir(dir string) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	return c.mkdir(dir)
}*/

// mkdir create directory (v2 的版本也没在用这函式)
/*func (c *EtcdClientV3) mkdir(dir string) error {
	if dir == "" || dir == "/" {
		return nil
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_, err := c.kapi.Put(cntx, dir, "", nil) // 找不太到参数
	if err != nil {
		if isErrNodeExists(err) {
			return nil
		}
		return err
	}
	return nil
}*/

// Create create path with data (v2 的版本也没在用这函式)
/*func (c *EtcdClientV3) Create(path string, data []byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	log.Debug("etcd create node %s", path)
	_, err := c.kapi.Put(cntx, path, string(data), nil) // 找不太到参数
	if err != nil {
		log.Debug("etcd create node %s failed: %s", path, err)
		return err
	}
	log.Debug("etcd create node OK")
	return nil
}*/

// 参考文件在 https://etcd.io/docs/v3.5/tutorials/how-to-get-key-by-prefix/
// 1 WithLease 是用来定时删除 key
// 2 WithLimit 是用来限制 etcd 的回传数量
// 3 WithRev，Revision 为 etcd key 的唯一值，为在 Revision 找值
// 4 WithMaxCreateRev 为在小于某一个 Revision 中找最大 Revision 的值
// 5 WithSort 为使用 get 时，对回传结果进行排序
// 6 WithPrefix 为取出前缀为 key 的值，比如 前缀为 foo，会回传 foo1 foo2
// 7 WithRange 为取出 key 值的范围，end key 值语义上要大于 key 值
// 8 WithFromKey 为取出 key 值的范围，但回传的 end key 会等于参数 end key
// 9 WithSerializable，Linearizability 为用来增加资料的正确性，Serializable 为用来减少延迟
// 10 WithKeysOnly 为当使用 get 时，只回传 key
// 11 WithCountOnly 为当使用 get 时，只回传 key
// 12 WithMinModRev 为会过滤小于 Revision 的 修改 Key
// 13 WithMaxModRev 为会过滤大于 Revision 的 修改 Key
// 14 WithMinCreateRev 为会过滤小于 Revision 的 新建 Key
// 15 WithMaxCreateRev 为会过滤大于 Revision 的 新建 Key

// Update update path with data
func (c *EtcdClientV3) Update(path string, data []byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_ = log.Debug("etcd update node %s", path)

	_, err := c.kapi.Put(cntx, path, string(data))
	if err != nil {
		_ = log.Debug("etcd update node %s failed: %s", path, err)
		return err
	}
	_ = log.Debug("etcd update node OK")
	return nil
}

// UpdateWithTTL update path with data and ttl
func (c *EtcdClientV3) UpdateWithTTL(path string, data []byte, ttl time.Duration) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_ = log.Debug("etcd update node %s with ttl %f", path, ttl.Seconds())

	lse, err := c.kapi.Grant(cntx, int64(ttl.Seconds()))
	if err != nil {
		_ = log.Debug("etcd lease node with ttl %d failed: %s ", ttl.Seconds(), err)
		return err
	}

	_, err = c.kapi.Put(cntx, path, string(data), clientv3.WithLease(lse.ID))
	if err != nil {
		log.Debug("etcd update node %s failed: %s", path, err)
		return err
	}
	log.Debug("etcd update node OK")
	return nil
}

// Lease create lease in etcd
func (c *EtcdClientV3) Lease(ttl time.Duration) (clientv3.LeaseID, error) {
	/*c.Lock()
	defer c.Unlock()
	if c.closed {
		return -1, ErrClosedEtcdClient
	}*/
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_ = log.Debug("etcd lease node with ttl %d", ttl)

	lse, err := c.kapi.Grant(cntx, int64(ttl.Seconds()))
	if err != nil {
		_ = log.Debug("etcd lease node with ttl %d failed: %s ", ttl, err)
		return -1, err
	}
	_ = log.Debug("etcd create lease OK with ID %d", lse.ID)
	return lse.ID, nil
}

// UpdateWithLease update path with data and ttl by using lease
func (c *EtcdClientV3) UpdateWithLease(path string, data []byte, leaseID clientv3.LeaseID) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_ = log.Debug("etcd update node %s with lease %d", path, leaseID)

	_, err := c.kapi.Put(cntx, path, string(data), clientv3.WithLease(leaseID))
	if err != nil {
		_ = log.Debug("etcd update node %s failed: %s", path, err)
		return err
	}
	_ = log.Debug("etcd update node OK")
	return nil
}

// Delete delete path
func (c *EtcdClientV3) Delete(path string) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_ = log.Debug("etcd delete node %s", path)
	_, err := c.kapi.Delete(cntx, path, clientv3.WithPrevKV())
	if err != nil {
		_ = log.Debug("etcd delete node %s failed: %s", path, err)
		return err
	}
	_ = log.Debug("etcd delete node OK")
	return nil
}

// Read read path data
func (c *EtcdClientV3) Read(path string) ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil, ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_ = log.Debug("etcd read node %s", path)
	r, err := c.kapi.Get(cntx, path, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	} else {
		if len(r.Kvs) > 0 {
			return r.Kvs[0].Value, nil
		}
	}
	return []byte{}, nil
}

// List list path, return slice of all paths
func (c *EtcdClientV3) List(path string) ([]string, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil, ErrClosedEtcdClient
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_ = log.Debug("etcd list node %s", path)
	r, err := c.kapi.Get(cntx, path, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	} else if r == nil {
		return nil, nil
	} else {
		var files []string
		for _, node := range r.Kvs {
			files = append(files, string(node.Key))
		}
		return files, nil
	}
}

// Watch watch path
func (c *EtcdClientV3) Watch(path string, ch chan string) error {
	c.Lock() // 在这里上锁
	// defer c.Unlock() // 移除此行，避免死结发生
	if c.closed {
		c.Unlock() // 上锁后记得解锁，去防止死结问题发生
		panic(ErrClosedEtcdClient)
	}

	rch := c.kapi.Watch(context.Background(), path, clientv3.WithPrefix())

	c.Unlock() // 上锁后在适当时机解锁，去防止死结问题发生
	// 在这里解锁是最好的，因为解锁后立刻可以进行监听

	for wresp := range rch {
		for _, ev := range wresp.Events {
			ch <- string(ev.Kv.Key)
		}
	}

	return nil
}

// BasePrefix return base prefix
func (c *EtcdClientV3) BasePrefix() string {
	return c.Prefix
}
