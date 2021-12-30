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
	"time"

	"github.com/XiaoMi/Gaea/mysql"
)

// PooledConnect app use this object to exec sql
type pooledConnectImpl struct {
	directConnection *DirectConnection
	pool             *connectionPoolImpl
	returnTime       time.Time
}

// Recycle return PooledConnect to the pool
func (pc *pooledConnectImpl) Recycle() {
	//if has error,the connection can’t be recycled
	if pc.directConnection.pkgErr != nil {
		pc.Close()
	}

	if pc.IsClosed() {
		pc.pool.Put(nil)
	} else {
		pc.pool.Put(pc)
		pc.returnTime = time.Now()
	}
}

// Reconnect replaces the existing underlying connection with a new one.
// If we get "MySQL server has gone away (errno 2006)", then call Reconnect
func (pc *pooledConnectImpl) Reconnect() error {
	pc.directConnection.Close()
	newConn, err := NewDirectConnection(pc.pool.addr, pc.pool.user, pc.pool.password, pc.pool.db, pc.pool.charset, pc.pool.collationID)
	if err != nil {
		return err
	}
	pc.directConnection = newConn
	return nil
}

// Close implement util.Resource interface
func (pc *pooledConnectImpl) Close() {
	pc.directConnection.Close()
}

// IsClosed check if pooled connection closed
func (pc *pooledConnectImpl) IsClosed() bool {
	if pc.directConnection == nil {
		return true
	}
	return pc.directConnection.IsClosed()
}

// UseDB  wrapper of direct connection, init database
func (pc *pooledConnectImpl) UseDB(db string) error {
	return pc.directConnection.UseDB(db)
}

func (pc *pooledConnectImpl) Ping() error {
	return pc.directConnection.Ping()
}
// Execute wrapper of direct connection, execute sql
func (pc *pooledConnectImpl) Execute(sql string, maxRows int) (*mysql.Result, error) {
	return pc.directConnection.Execute(sql, maxRows)
}

// SetAutoCommit wrapper of direct connection, set autocommit
func (pc *pooledConnectImpl) SetAutoCommit(v uint8) error {
	return pc.directConnection.SetAutoCommit(v)
}

// Begin wrapper of direct connection, begin transaction
func (pc *pooledConnectImpl) Begin() error {
	return pc.directConnection.Begin()
}

// Commit wrapper of direct connection, commit transaction
func (pc *pooledConnectImpl) Commit() error {
	return pc.directConnection.Commit()
}

// Rollback wrapper of direct connection, rollback transaction
func (pc *pooledConnectImpl) Rollback() error {
	return pc.directConnection.Rollback()
}

// SetCharset wrapper of direct connection, set charset of connection
func (pc *pooledConnectImpl) SetCharset(charset string, collation mysql.CollationID) (bool, error) {
	return pc.directConnection.SetCharset(charset, collation)
}

// FieldList wrapper of direct connection, send field list to mysql
func (pc *pooledConnectImpl) FieldList(table string, wildcard string) ([]*mysql.Field, error) {
	return pc.directConnection.FieldList(table, wildcard)
}

// GetAddr wrapper of return addr of direct connection
func (pc *pooledConnectImpl) GetAddr() string {
	return pc.directConnection.GetAddr()
}

// SetSessionVariables set pc variables according to session
func (pc *pooledConnectImpl) SetSessionVariables(frontend *mysql.SessionVariables) (bool, error) {
	return pc.directConnection.SetSessionVariables(frontend)
}

// WriteSetStatement exec sql
func (pc *pooledConnectImpl) WriteSetStatement() error {
	return pc.directConnection.WriteSetStatement()
}

func (pc *pooledConnectImpl) GetConnectionID() int64 {
	return int64(pc.directConnection.conn.ConnectionID)
}

func (pc *pooledConnectImpl) GetReturnTime() time.Time {
	return pc.returnTime
}