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

	"github.com/XiaoMi/Gaea/mysql"
)

// PooledConnection app use this object to exec sql
type PooledConnection struct {
	directConnection *DirectConnection
	pool             *ConnectionPool
}

// Recycle return PooledConnection to the pool
func (pc *PooledConnection) Recycle() {
	if pc.IsClosed() {
		pc.pool.Put(nil)
	} else {
		pc.pool.Put(pc)
	}
}

// Reconnect replaces the existing underlying connection with a new one.
// If we get "MySQL server has gone away (errno 2006)", then call Reconnect
func (pc *PooledConnection) Reconnect() error {
	pc.directConnection.Close()
	newConn, err := NewDirectConnection(pc.pool.addr, pc.pool.user, pc.pool.password, pc.pool.db, pc.pool.charset, pc.pool.collationID)
	if err != nil {
		return err
	}
	pc.directConnection = newConn
	return nil
}

// Close implement util.Resource interface
func (pc *PooledConnection) Close() {
	pc.directConnection.Close()
}

// IsClosed check if pooled connection closed
func (pc *PooledConnection) IsClosed() bool {
	if pc.directConnection == nil {
		return true
	}
	return pc.directConnection.IsClosed()
}

// UseDB  wrapper of direct connection, init database
func (pc *PooledConnection) UseDB(db string) error {
	return pc.directConnection.UseDB(db)
}

// Execute wrapper of direct connection, execute sql
func (pc *PooledConnection) Execute(sql string) (*mysql.Result, error) {
	return pc.directConnection.Execute(sql)
}

func (pc *PooledConnection) ExecuteWithCtx(sql string, ctx context.Context, maxSelectResultSet int64) (*mysql.Result, error) {
	return pc.directConnection.ExecuteWithCtx(sql, ctx, maxSelectResultSet)
}

// SetAutoCommit wrapper of direct connection, set autocommit
func (pc *PooledConnection) SetAutoCommit(v uint8) error {
	return pc.directConnection.SetAutoCommit(v)
}

// Begin wrapper of direct connection, begin transaction
func (pc *PooledConnection) Begin() error {
	return pc.directConnection.Begin()
}

// Commit wrapper of direct connection, commit transaction
func (pc *PooledConnection) Commit() error {
	return pc.directConnection.Commit()
}

// Rollback wrapper of direct connection, rollback transaction
func (pc *PooledConnection) Rollback() error {
	return pc.directConnection.Rollback()
}

// SetCharset wrapper of direct connection, set charset of connection
func (pc *PooledConnection) SetCharset(charset string, collation mysql.CollationID) (bool, error) {
	return pc.directConnection.SetCharset(charset, collation)
}

// FieldList wrapper of direct connection, send field list to mysql
func (pc *PooledConnection) FieldList(table string, wildcard string) ([]*mysql.Field, error) {
	return pc.directConnection.FieldList(table, wildcard)
}

// GetAddr wrapper of return addr of direct connection
func (pc *PooledConnection) GetAddr() string {
	return pc.directConnection.GetAddr()
}

// SetSessionVariables set pc variables according to session
func (pc *PooledConnection) SetSessionVariables(frontend *mysql.SessionVariables) (bool, error) {
	return pc.directConnection.SetSessionVariables(frontend)
}

// WriteSetStatement exec sql
func (pc *PooledConnection) WriteSetStatement() error {
	return pc.directConnection.WriteSetStatement()
}
