// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"context"
	"time"

	"github.com/XiaoMi/Gaea/mysql"
)

type PooledConnect interface {
	Recycle()
	Reconnect() error
	Close()
	IsClosed() bool
	UseDB(db string) error
	Execute(sql string, maxRows int) (*mysql.Result, error)
	ExecuteWithTimeout(sql string, maxRows int, timeout time.Duration) (*mysql.Result, error)
	SetAutoCommit(v uint8) error
	Begin() error
	Commit() error
	Rollback() error
	Ping() error
	PingWithTimeout(timeout time.Duration) error
	SetCharset(charset string, collation mysql.CollationID) (bool, error)
	FieldList(table string, wildcard string) ([]*mysql.Field, error)
	GetAddr() string
	SetSessionVariables(frontend *mysql.SessionVariables) (bool, error)
	SyncSessionVariables(frontend *mysql.SessionVariables) error
	WriteSetStatement() error
	GetConnectionID() int64
	GetReturnTime() time.Time
	MoreRowsExist() bool
	MoreResultsExist() bool
	FetchMoreRows(result *mysql.Result, maxRows int) error
	ReadMoreResult(maxRows int) (*mysql.Result, error)
}

type ConnectionPool interface {
	Open() error
	Addr() string
	Datacenter() string
	Close()
	Get(ctx context.Context) (PooledConnect, error)
	GetCheck(ctx context.Context) (PooledConnect, error)
	Put(pc PooledConnect)

	SetCapacity(capacity int) (err error)
	SetIdleTimeout(idleTimeout time.Duration)
	StatsJSON() string
	Capacity() int64
	Available() int64
	Active() int64
	InUse() int64
	MaxCap() int64
	WaitCount() int64
	WaitTime() time.Duration
	IdleTimeout() time.Duration
	IdleClosed() int64
	SetLastChecked()
	GetLastChecked() int64
}
