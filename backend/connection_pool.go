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
	"sync"
	"time"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
)

const (
	getConnTimeout = 2 * time.Second
)

var (
	// ErrConnectionPoolClosed means pool closed error
	ErrConnectionPoolClosed = errors.New("connection pool is closed")
	// DefaultCapacity default capacity of connection pool
	DefaultCapacity = 64
)

// ConnectionPool means connection pool with specific addr
type ConnectionPool struct {
	mu          sync.RWMutex
	connections *util.ResourcePool

	addr     string
	user     string
	password string
	db       string

	charset     string
	collationID mysql.CollationID

	capacity    int // capacity of pool
	maxCapacity int // max capacity of pool
	idleTimeout time.Duration
}

// NewConnectionPool create connection pool
func NewConnectionPool(addr, user, password, db string, capacity, maxCapacity int, idleTimeout time.Duration, charset string, collationID mysql.CollationID) *ConnectionPool {
	cp := &ConnectionPool{addr: addr, user: user, password: password, db: db, capacity: capacity, maxCapacity: maxCapacity, idleTimeout: idleTimeout, charset: charset, collationID: collationID}
	return cp
}

func (cp *ConnectionPool) pool() (p *util.ResourcePool) {
	cp.mu.Lock()
	p = cp.connections
	cp.mu.Unlock()
	return p
}

// Open open connection pool without error, should be called before use the pool
func (cp *ConnectionPool) Open() {
	if cp.capacity == 0 {
		cp.capacity = DefaultCapacity
	}

	if cp.maxCapacity == 0 {
		cp.maxCapacity = cp.capacity
	}
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.connections = util.NewResourcePool(cp.connect, cp.capacity, cp.maxCapacity, cp.idleTimeout)
	return
}

// connect is used by the resource pool to create new resource.It's factory method
func (cp *ConnectionPool) connect() (util.Resource, error) {
	c, err := NewDirectConnection(cp.addr, cp.user, cp.password, cp.db, cp.charset, cp.collationID)
	if err != nil {
		return nil, err
	}
	return &PooledConnection{directConnection: c, pool: cp}, nil
}

// Addr return addr of connection pool
func (cp *ConnectionPool) Addr() string {
	return cp.addr
}

// Close close connection pool
func (cp *ConnectionPool) Close() {
	p := cp.pool()
	if p == nil {
		return
	}
	p.Close()
	cp.mu.Lock()
	cp.connections = nil
	cp.mu.Unlock()
	return
}

// tryReuse reset params of connection before reuse
func (cp *ConnectionPool) tryReuse(pc *PooledConnection) error {
	return pc.directConnection.ResetConnection()
}

// Get return a connection, you should call PooledConnection's Recycle once done
func (cp *ConnectionPool) Get(ctx context.Context) (*PooledConnection, error) {
	p := cp.pool()
	if p == nil {
		return nil, ErrConnectionPoolClosed
	}

	getCtx, cancel := context.WithTimeout(ctx, getConnTimeout)
	defer cancel()
	r, err := p.Get(getCtx)
	if err != nil {
		return nil, err
	}
	return r.(*PooledConnection), nil
}

// Put recycle a connection into the pool
func (cp *ConnectionPool) Put(pc *PooledConnection) {
	p := cp.pool()
	if p == nil {
		panic(ErrConnectionPoolClosed)
	}

	if pc == nil {
		p.Put(nil)
	} else if err := cp.tryReuse(pc); err != nil {
		pc.Close()
		p.Put(nil)
	} else {
		p.Put(pc)
	}
}

// SetCapacity alert the size of the pool at runtime
func (cp *ConnectionPool) SetCapacity(capacity int) (err error) {
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
func (cp *ConnectionPool) SetIdleTimeout(idleTimeout time.Duration) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.connections != nil {
		cp.connections.SetIdleTimeout(idleTimeout)
	}
	cp.idleTimeout = idleTimeout
}

// StatsJSON return the pool stats as JSON object.
func (cp *ConnectionPool) StatsJSON() string {
	p := cp.pool()
	if p == nil {
		return "{}"
	}
	return p.StatsJSON()
}

// Capacity return the pool capacity
func (cp *ConnectionPool) Capacity() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Capacity()
}

// Available returns the number of available connections in the pool
func (cp *ConnectionPool) Available() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Available()
}

// Active returns the number of active connections in the pool
func (cp *ConnectionPool) Active() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.Active()
}

// InUse returns the number of in-use connections in the pool
func (cp *ConnectionPool) InUse() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.InUse()
}

// MaxCap returns the maximum size of the pool
func (cp *ConnectionPool) MaxCap() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.MaxCap()
}

// WaitCount returns how many clients are waitting for a connection
func (cp *ConnectionPool) WaitCount() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.WaitCount()
}

// WaitTime returns the time wait for a connection
func (cp *ConnectionPool) WaitTime() time.Duration {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.WaitTime()
}

// IdleTimeout returns the idle timeout for the pool
func (cp *ConnectionPool) IdleTimeout() time.Duration {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.IdleTimeout()
}

// IdleClosed return the number of closed connections for the pool
func (cp *ConnectionPool) IdleClosed() int64 {
	p := cp.pool()
	if p == nil {
		return 0
	}
	return p.IdleClosed()
}
