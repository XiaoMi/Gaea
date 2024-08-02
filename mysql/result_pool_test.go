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

package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test result set
func TestGetResultSet(t *testing.T) {
	rs := ResultPool.Get()
	assert.Equal(t, rs.Resultset == nil, false)
	rs.Free()
	rs = ResultPool.GetWithoutResultSet()
	assert.Equal(t, rs.Resultset == nil, true)
	rs.Free()
	rss := make([]*Result, 0)
	for i := 0; i < 10; i++ {
		rs = ResultPool.Get()
		assert.Equal(t, rs.Resultset == nil, false)
		rss = append(rss, rs)
	}
	for i := 0; i < 10; i++ {
		rs = ResultPool.GetWithoutResultSet()
		assert.Equal(t, rs.Resultset == nil, true)
		rss = append(rss, rs)
	}
	for _, rs := range rss {
		rs.Free()
	}
	for i := 0; i < 10; i++ {
		rs = ResultPool.Get()
		assert.Equal(t, rs.Resultset == nil, false)
		rs.Free()
	}
	for i := 0; i < 10; i++ {
		rs = ResultPool.GetWithoutResultSet()
		assert.Equal(t, rs.Resultset == nil, true)
		rs.Free()
	}
}
