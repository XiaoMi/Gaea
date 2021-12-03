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

import (
	"sync"
)

const (
	// StmtType stmt type
	StmtType = "stmtType" // SQL类型, 值类型为int (对应parser.Preview()得到的值)
	// FromSlave if read from slave
	FromSlave    = "fromSlave"    // 读写分离标识, 值类型为int, false = 0, true = 1
	DefaultSlice = "defaultSlice" // 默认分片标识 string 类型
)

// RequestContext means request scope context with values
// thread safe
type RequestContext struct {
	lock *sync.RWMutex
	ctx  map[string]interface{}
}

// NewRequestContext return request scopre context
func NewRequestContext() *RequestContext {
	return &RequestContext{ctx: make(map[string]interface{}, 2), lock: new(sync.RWMutex)}
}

// Get return context in RequestContext
func (reqCtx *RequestContext) Get(key string) interface{} {
	reqCtx.lock.RLock()
	v, ok := reqCtx.ctx[key]
	reqCtx.lock.RUnlock()
	if ok {
		return v
	}
	return nil
}

// Set set value with specific key
func (reqCtx *RequestContext) Set(key string, value interface{}) {
	reqCtx.lock.Lock()
	reqCtx.ctx[key] = value
	reqCtx.lock.Unlock()
}
