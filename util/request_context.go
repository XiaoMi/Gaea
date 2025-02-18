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

package util

import "github.com/XiaoMi/Gaea/mysql"

// RequestContext means request scope context with values
// 旧版 thread safe，因为 context 是顺序执行的，把锁去掉，提升性能，新版本 thread unsafe
type RequestContext struct {
	tokens         []string
	cmdType        byte
	stmtType       int
	fromSlave      int
	fingerprint    string
	fingerprintMD5 string
	defaultSlice   string
}

// NewRequestContext return request scopre context
func NewRequestContext() *RequestContext {
	return &RequestContext{}
}

func (reqCtx *RequestContext) GetStmtType() int {
	return reqCtx.stmtType
}

func (reqCtx *RequestContext) SetStmtType(value int) {
	reqCtx.stmtType = value
}

// SetCmdStmtType sets the command statement type for the request context.
// This method is used to assign the `cmdType` in the RequestContext,
// which helps identify the type of command being executed (e.g., prepare or execute SQL).
func (reqCtx *RequestContext) SetCmdStmtType(value byte) {
	reqCtx.cmdType = value
}

// IsPrepareSQL checks if the current command is a prepared SQL statement (either COM_STMT_PREPARE or COM_STMT_EXECUTE).
// This method is used for logging purposes to differentiate between regular SQL and prepared SQL statements.
func (reqCtx *RequestContext) IsPrepareSQL() bool {
	return reqCtx.cmdType == mysql.ComStmtExecute || reqCtx.cmdType == mysql.ComStmtPrepare
}

func (reqCtx *RequestContext) GetTokens() []string {
	return reqCtx.tokens
}

func (reqCtx *RequestContext) SetTokens(value []string) {
	reqCtx.tokens = value
}

func (reqCtx *RequestContext) GetFromSlave() int {
	return reqCtx.fromSlave
}

func (reqCtx *RequestContext) SetFromSlave(value int) {
	reqCtx.fromSlave = value
}

func (reqCtx *RequestContext) GetFingerprint() string {
	return reqCtx.fingerprint
}

func (reqCtx *RequestContext) SetFingerprint(value string) {
	reqCtx.fingerprint = value
}

func (reqCtx *RequestContext) GetFingerprintMD5() string {
	return reqCtx.fingerprintMD5
}

func (reqCtx *RequestContext) SetFingerprintMD5(value string) {
	reqCtx.fingerprintMD5 = value
}

func (reqCtx *RequestContext) GetDefaultSlice() string {
	return reqCtx.defaultSlice
}

func (reqCtx *RequestContext) SetDefaultSlice(value string) {
	reqCtx.defaultSlice = value
}
