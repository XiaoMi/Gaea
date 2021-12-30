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
	SetAutoCommit(v uint8) error
	Begin() error
	Commit() error
	Rollback() error
	Ping() error
	SetCharset(charset string, collation mysql.CollationID) (bool, error)
	FieldList(table string, wildcard string) ([]*mysql.Field, error)
	GetAddr() string
	SetSessionVariables(frontend *mysql.SessionVariables) (bool, error)
	WriteSetStatement() error
	GetConnectionID() int64
	GetReturnTime() time.Time
}

type ConnectionPool interface {
	Open()
	Addr() string
	Close()
	Get(ctx context.Context) (PooledConnect, error)
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
}
