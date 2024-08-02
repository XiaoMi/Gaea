// Copyright 2024 The Gaea Authors. All Rights Reserved.
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

// 借鉴 zap 优秀的设计

package mysql

import "sync"

type resultPool struct {
	p1 *sync.Pool
	p2 *sync.Pool
}

var ResultPool = &resultPool{
	// 有 result 的对象池
	p1: &sync.Pool{
		New: func() interface{} {
			return new(Result)
		},
	},
	// 没有 result 的对象池
	p2: &sync.Pool{
		New: func() interface{} {
			return new(Result)
		},
	},
}

func (rp *resultPool) Get() *Result {
	r := rp.p1.Get().(*Result)
	r.pool = rp
	r.Reset()
	if r.Resultset == nil {
		r.Resultset = &Resultset{}
	}
	return r
}

func (rp *resultPool) GetWithoutResultSet() *Result {
	r := rp.p2.Get().(*Result)
	r.pool = rp
	r.Reset()
	r.Resultset = nil
	return r
}

func (rp *resultPool) Put(r *Result) {
	if r.Resultset != nil {
		rp.p1.Put(r)
	} else {
		rp.p2.Put(r)
	}
}

func (r *Result) Reset() {
	r.Status = 0
	r.InsertID = 0
	r.AffectedRows = 0
	r.Warnings = 0
	r.Info = ""
	if r.Resultset != nil {
		r.Resultset.Fields = r.Resultset.Fields[:0]
		r.Resultset.Values = r.Resultset.Values[:0]
		r.Resultset.RowDatas = r.Resultset.RowDatas[:0]
		r.Resultset.FieldNames = make(map[string]int)
	}
}

func (r *Result) Free() {
	if r.pool != nil {
		r.pool.Put(r)
	}
}
