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

package plan

import (
	"fmt"
)

// RouteResult is the route result of a statement
// 遍历AST之后得到的路由结果
// db, table唯一确定了一个路由, 这里只记录分片表的db和table, 如果是关联表, 必须关联到同一个父表
type RouteResult struct {
	db    string
	table string

	currentIndex int   // 当前遍历indexes位置下标
	indexes      []int // 分片索引列表, 是有序的
}

// NewRouteResult constructor of RouteResult
func NewRouteResult(db, table string, initIndexes []int) *RouteResult {
	return &RouteResult{
		db:      db,
		table:   table,
		indexes: initIndexes,
	}
}

// Check check if the db and table is valid
func (r *RouteResult) Check(db, table string) error {
	if r.db != db {
		return fmt.Errorf("db not equal, origin: %v, current: %v", r.db, db)
	}
	if r.table != table {
		return fmt.Errorf("table not equal, origin: %v, current: %v", r.table, table)
	}
	return nil
}

// Inter inter indexes with origin indexes in RouteResult
// 如果是关联表, db, table需要用父表的db和table
func (r *RouteResult) Inter(indexes []int) {
	r.indexes = interList(r.indexes, indexes)
}

// Union union indexes with origin indexes in RouteResult
// 如果是关联表, db, table需要用父表的db和table
func (r *RouteResult) Union(indexes []int) {
	r.indexes = unionList(r.indexes, indexes)
}

// GetShardIndexes get shard indexes
func (r *RouteResult) GetShardIndexes() []int {
	return r.indexes
}

// GetCurrentTableIndex get current table index
func (r *RouteResult) GetCurrentTableIndex() (int, error) {
	if r.currentIndex >= len(r.indexes) {
		return -1, fmt.Errorf("table index out of range")
	}
	return r.indexes[r.currentIndex], nil
}

// Next get next table index
func (r *RouteResult) Next() int {
	idx := r.currentIndex
	r.currentIndex++
	return r.indexes[idx]
}

// HasNext check if has next index
func (r *RouteResult) HasNext() bool {
	return r.currentIndex < len(r.indexes)
}

// Reset reset the cursor of index
func (r *RouteResult) Reset() {
	r.currentIndex = 0
}
