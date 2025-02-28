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

package backend

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
)

const (
	ExecTimeOut             = 2 * time.Second
	GetConnTimeout          = 2 * time.Second
	pingPeriod              = 4 * time.Second
	handshakeTimeoutDefault = 500 * time.Millisecond
)

const (
	MonitorDefaultCapacity    = 3
	MonitorDefaultMaxCapacity = 5
)

var (
	// ErrConnectionPoolClosed means pool closed error
	ErrConnectionPoolClosed = errors.New("connection pool is closed")
	// DefaultCapacity default capacity of connection pool
	DefaultCapacity = 64
)

// connectionPoolImpl means connection pool with specific addr
type connectionPoolImpl struct {
	mu          sync.RWMutex
	connections *util.ResourcePool
	checkConn   *pooledConnectImpl

	addr       string
	datacenter string
	user       string
	password   string
	db         string

	charset     string
	collationID mysql.CollationID

	capacity         int // capacity of pool
	maxCapacity      int // max capacity of pool
	idleTimeout      time.Duration
	clientCapability uint32
	initConnect      string
	lastChecked      int64
	handshakeTimeout time.Duration
}

// NewConnectionPool create connection pool
func NewConnectionPool(addr, user, password, db string, capacity, maxCapacity int, idleTimeout time.Duration, charset string, collationID mysql.CollationID, clientCapability uint32, initConnect string, dc string, handshakeTimeout time.Duration) ConnectionPool {
	return &connectionPoolImpl{
		addr:             addr,
		datacenter:       dc,
		user:             user,
		password:         password,
		db:               db,
		capacity:         capacity,
		maxCapacity:      maxCapacity,
		idleTimeout:      idleTimeout,
		charset:          charset,
		collationID:      collationID,
		clientCapability: clientCapability,
		initConnect:      strings.Trim(strings.TrimSpace(initConnect), ";"),
		lastChecked:      time.Now().Unix(),
		handshakeTimeout: handshakeTimeout,
	}
}

func (cp *connectionPoolImpl) pool() (p *util.ResourcePool) {
	cp.mu.Lock()
	p = cp.connections
	cp.mu.Unlock()
	return p
}

// Open open connection pool without error, should be called before use the pool
func (cp *connectionPoolImpl) Open() error {
	if cp.capacity == 0 {
		cp.capacity = DefaultCapacity
	}

	if cp.maxCapacity == 0 {
		cp.maxCapacity = cp.capacity
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()
	var err error = nil
	cp.connections, err = util.NewResourcePool(
		cp.connect, cp.capacity, cp.maxCapacity, cp.idleTimeout,
	)
	return err
}

// connect is used by the resource pool to create new resource.It's factory method
func (cp *connectionPoolImpl) connect() (util.Resource, error) {
	c, err := NewDirectConnection(cp.addr, cp.user, cp.password, cp.db, cp.charset, cp.collationID, cp.clientCapability, cp.handshakeTimeout)
	if err != nil {
		return nil, err
	}
	if cp.initConnect != "" {
		for _, sql := range strings.Split(cp.initConnect, ";") {
			_, err := c.Execute(sql, 0)
			if err != nil {
				return nil, err
			}
		}
	}
	return &pooledConnectImpl{directConnection: c, pool: cp}, nil
}

// Addr return addr of connection pool
func (cp *connectionPoolImpl) Addr() string {
	return cp.addr
}

// Datacenter return datacenter of connection pool
func (cp *connectionPoolImpl) Datacenter() string {
	return cp.datacenter
}

// Close close connection pool
func (cp *connectionPoolImpl) Close() {
	p := cp.pool()
	if p == nil {
		return
	}
	p.Close()
	cp.mu.Lock()
	// close check conn
	if cp.checkConn != nil {
		cp.checkConn.Close()
		cp.checkConn = nil
	}
	cp.connections = nil
	cp.mu.Unlock()
	return
}

// tryReuse reset params of connection before reuse
func (cp *connectionPoolImpl) tryReuse(pc *pooledConnectImpl) error {
	return pc.directConnection.ResetConnection()
}

// Get return a connection, you should call PooledConnect's Recycle once done
func (cp *connectionPoolImpl) Get(ctx context.Context) (pc PooledConnect, err error) {
	p := cp.pool()
	if p == nil {
		return nil, ErrConnectionPoolClosed
	}

	getCtx, cancel := context.WithTimeout(ctx, GetConnTimeout)
	defer cancel()
	r, err := p.Get(getCtx)
	if err != nil {
		return nil, err
	}

	pc = r.(*pooledConnectImpl)

	//do ping when over the ping time. if error happen, create new one
	if !pc.GetReturnTime().IsZero() && time.Until(pc.GetReturnTime().Add(pingPeriod)) < 0 {
		if err = pc.PingWithTimeout(GetConnTimeout); err != nil {
			err = pc.Reconnect()
		}
		// 这边防止连接泄漏
		if err != nil {
			pc.Recycle()
			return nil, err
		}
	}

	return pc, err
}

// GetCheck return a check backend db connection, which independent with connection pool
func (cp *connectionPoolImpl) GetCheck(ctx context.Context) (PooledConnect, error) {
	if cp.checkConn != nil && !cp.checkConn.IsClosed() {
		return cp.checkConn, nil
	}

	getCtx, cancel := context.WithTimeout(ctx, GetConnTimeout)
	defer cancel()

	getConnChan := make(chan error)
	go func() {
		// connect timeout will be in 2s
		checkConn, err := cp.connect()
		if err != nil {
			return
		}
		cp.checkConn = checkConn.(*pooledConnectImpl)

		if cp.checkConn.IsClosed() {
			if err := cp.checkConn.Reconnect(); err != nil {
				return
			}
		}
		getConnChan <- err
	}()

	select {
	case <-getCtx.Done():
		return nil, fmt.Errorf("get conn timeout")
	case err1 := <-getConnChan:
		if err1 != nil {
			return nil, err1
		}
		return cp.checkConn, nil
	}

}

// Put recycle a connection into the pool
func (cp *connectionPoolImpl) Put(pc PooledConnect) {
	p := cp.pool()
	if p == nil {
		panic(ErrConnectionPoolClosed)
	}

	if pc == nil {
		p.Put(nil)
	} else if err := cp.tryReuse(pc.(*pooledConnectImpl)); err != nil {
		pc.Close()
		p.Put(nil)
	} else {
		p.Put(pc)
	}
}

// SetCapacity alert the size of the pool at runtime
func (cp *connectionPoolImpl) SetCapacity(capacity int) (err error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.connections != nil {
		err = cp.connections.SetCapacity(capacity)
		if err != nil {
			return err
		}
	}
	cp.capacity = capacity
	return nil
}

// SetIdleTimeout set the idleTimeout of the pool
func (cp *connectionPoolImpl) SetIdleTimeout(idleTimeout time.Duration) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.connections != nil {
		cp.connections.SetIdleTimeout(idleTimeout)
	}
	cp.idleTimeout = idleTimeout
}

// StatsJSON return the pool stats as JSON object.
func (cp *connectionPoolImpl) StatsJSON() string {
	p := cp.pool()
	if p == nil {
		return "{}"
	}
	return p.StatsJSON()
}

// Capacity return the pool capacity
func (cp *connectionPoolImpl) Capacity() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Capacity()
}

// Available returns the number of available connections in the pool
func (cp *connectionPoolImpl) Available() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Available()
}

// Active returns the number of active connections in the pool
func (cp *connectionPoolImpl) Active() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Active()
}

// InUse returns the number of in-use connections in the pool
func (cp *connectionPoolImpl) InUse() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.InUse()
}

// MaxCap returns the maximum size of the pool
func (cp *connectionPoolImpl) MaxCap() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.MaxCap()
}

// WaitCount returns how many clients are waitting for a connection
func (cp *connectionPoolImpl) WaitCount() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.WaitCount()
}

// WaitTime returns the time wait for a connection
func (cp *connectionPoolImpl) WaitTime() time.Duration {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.WaitTime()
}

// IdleTimeout returns the idle timeout for the pool
func (cp *connectionPoolImpl) IdleTimeout() time.Duration {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.IdleTimeout()
}

// IdleClosed return the number of closed connections for the pool
func (cp *connectionPoolImpl) IdleClosed() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.IdleClosed()
}

// SetLastChecked set last checked time
func (cp *connectionPoolImpl) SetLastChecked() {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp == nil {
		return
	}
	cp.lastChecked = time.Now().Unix()
}

// GetLastChecked get last checked time
func (cp *connectionPoolImpl) GetLastChecked() int64 {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	if cp == nil {
		return 0
	}
	return cp.lastChecked
}
