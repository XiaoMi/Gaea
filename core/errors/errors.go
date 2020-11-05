// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package errors

import (
	"errors"
)

var (
	// ErrNoMasterConn no master connection
	ErrNoMasterConn = errors.New("no master connection")
	// ErrNoSlaveConn no slave connection
	ErrNoSlaveConn = errors.New("no slave connection")
	// ErrNoDefaultSlice no default slice
	ErrNoDefaultSlice = errors.New("no default slice")
	// ErrNoMasterDB no master database
	ErrNoMasterDB = errors.New("no master database")
	// ErrNoSlaveDB no slve database
	ErrNoSlaveDB = errors.New("no slave database")
	// ErrNoDatabase no database
	ErrNoDatabase = errors.New("no database")

	// ErrMasterDown master is down
	ErrMasterDown = errors.New("master is down")
	// ErrSlaveDown slave is down
	ErrSlaveDown = errors.New("slave is down")

	// ErrInvalidArgument invalid arguments
	ErrInvalidArgument = errors.New("argument is invalid")
	// ErrInvalidCharset invalid charset
	ErrInvalidCharset = errors.New("charset is invalid")
	// ErrCmdUnsupport unsupport command
	ErrCmdUnsupport = errors.New("command unsupport")

	// ErrLocationsCount locations count is wrong
	ErrLocationsCount = errors.New("locations count is not equal")
	// ErrInvalidMycatLocations mycat locations is wrong
	ErrInvalidMycatLocations = errors.New("mycat locations must be all 1")
	// ErrNoCriteria no shard condition
	ErrNoCriteria = errors.New("plan have no criteria")
	// ErrSelectInInsert insert sql has sub select sql, not support
	ErrSelectInInsert = errors.New("select in insert not allowed")
	// ErrInsertInMulti insert into multi slice
	ErrInsertInMulti = errors.New("insert in multi slice")
	// ErrUpdateInMulti update multi slice
	ErrUpdateInMulti = errors.New("update in multi slice")
	// ErrDeleteInMulti delete multi slice
	ErrDeleteInMulti = errors.New("delete in multi slice")
	// ErrReplaceInMulti replace multi slice
	ErrReplaceInMulti = errors.New("replace in multi slice")
	// ErrExecInMulti execute in multi slice
	ErrExecInMulti = errors.New("exec in multi slice")
	// ErrTransInMulti transaction cross slices, not support
	ErrTransInMulti = errors.New("transaction in multi slice")
	// ErrUnsupportedShard unsupport shard type
	ErrUnsupportedShard = errors.New("sql is unsupported in shard mode")

	// ErrNoPlan get no plan
	ErrNoPlan = errors.New("statement have no plan")
	// ErrUpdateKey update shard key
	ErrUpdateKey = errors.New("routing key in update expression")
	// ErrConnNotEqual count of sqls doesn't match lenth of conns
	ErrConnNotEqual = errors.New("the length of conns not equal sqls")
	// ErrKeyOutOfRange shard key out of range
	ErrKeyOutOfRange = errors.New("shard key not in key range")
	// ErrIRNoColumns insert or replace doesn't specify columns
	ErrIRNoColumns = errors.New("insert or replace must specify columns")
	// ErrIRNoShardingKey insert or replace not contain sharding key
	ErrIRNoShardingKey = errors.New("insert or replace not contain sharding key")
	// ErrDateRangeIllegal date range invalided
	ErrDateRangeIllegal = errors.New("date range format illegal")
	// ErrDateRangeCount date range count invalided
	ErrDateRangeCount = errors.New("date range count is not equal")
	// ErrDateRangeOverlap date range overlap
	ErrDateRangeOverlap = errors.New("date range is overlapped")
	// ErrUnknownRuleType unknown rule type
	ErrUnknownRuleType = errors.New("unknown rule type")
	// ErrInvalidSQLMode invalid sql mode
	ErrInvalidSQLMode = errors.New("invalid sql mode")
	// ErrInternalServer internal error
	ErrInternalServer = errors.New("internal server error")
	// ErrUserIsReadOnly user is readonly
	ErrUserIsReadOnly = errors.New("user is readonly")

	// ErrNamespaceNotPrepared commit namespace config without prepare
	ErrNamespaceNotPrepared = errors.New("namespace is not prepared")
)
