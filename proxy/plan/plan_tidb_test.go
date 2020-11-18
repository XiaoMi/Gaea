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
	"strings"
	"testing"

	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
)

func BenchmarkSelectStmtCheckShard(b *testing.B) {
	r, err := prepareShardKingshardRouter()
	if err != nil {
		b.Fatal(err)
	}
	benchmarks := []struct {
		sql     string
		isShard bool
	}{
		{"select * from test_hash_0, test_hash_1 where test_hash_0.id = test_hash_1.id", true},
		{"select * from test_a, test_b left join test_c on test_a.id=test_c.id where test_a.id in (1,2,3) or test_b.k = 0 order by test_a.id desc limit 10", false},
	}
	for _, bm := range benchmarks {
		b.Run(bm.sql, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				//for i := 0; i < 1; i++ {
				stmt, err := parser.ParseSQL(bm.sql)
				if err != nil {
					b.Fatal(err)
				}
				selectStmt, ok := stmt.(*ast.SelectStmt)
				if !ok {
					b.Fatal("not a select stmt")
				}
				visitor := NewChecker("test", r)
				selectStmt.Accept(visitor)
				if visitor.IsShard() != bm.isShard {
					b.Errorf("isShard not equal, expect: %v, actual: %v", bm.isShard, visitor.IsShard())
				}
			}
		})
	}
}

func TestSelectStmtCheckShard(t *testing.T) {
	r, err := prepareShardKingshardRouter()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		sql     string
		isShard bool
	}{
		{"select * from test_hash_0, test_hash_1 where test_hash_0.id = test_hash_1.id", true},
		{"select * from test_a, test_b left join test_c on test_a.id=test_c.id where test_a.id in (1,2,3) or test_b.k = 0 order by test_a.id desc limit 10", false},
	}
	for _, test := range tests {
		t.Run(test.sql, func(t *testing.T) {
			//for i := 0; i < 1; i++ {
			stmt, err := parser.ParseSQL(test.sql)
			if err != nil {
				t.Fatal(err)
			}
			selectStmt, ok := stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatal("not a select stmt")
			}
			visitor := NewChecker("test", r)
			selectStmt.Accept(visitor)
			if visitor.IsShard() != test.isShard {
				t.Errorf("isShard not equal, expect: %v, actual: %v", test.isShard, visitor.IsShard())
			}
		})
	}
}

// TODO: no router, panic, change to table test function
func _TestGroupByRewriting(t *testing.T) {
	tests := []struct {
		sql        string
		rewrite    string
		groupByCol []int
		count      int
	}{
		{"select * from tbl1 group by a, b", "SELECT *,a,b FROM tbl1",
			[]int{1, 2}, 2},
	}
	for _, test := range tests {
		t.Run(test.sql, func(t *testing.T) {
			stmt, err := parser.ParseSQL(test.sql)
			if err != nil {
				t.Fatal(err)
			}
			selectStmt, ok := stmt.(*ast.SelectStmt)
			if !ok {
				t.Fatal("not a select stmt")
			}

			info := NewSelectPlan("test", test.sql, nil)
			if err := HandleSelectStmt(info, selectStmt); err != nil {
				t.Fatal(err)
			}
			s := &strings.Builder{}
			selectStmt.Restore(format.NewRestoreCtx(0, s))
			rewriteSQL := s.String()
			if rewriteSQL != test.rewrite {
				t.Errorf("rewrite sql not equal, expect: %v, actual: %v", test.rewrite, rewriteSQL)
			}
			if len(info.GetGroupByColumnInfo()) != test.count {
				t.Errorf("rewrite sql not equal, expect: %v, actual: %v", test.rewrite, rewriteSQL)
			}
			for i, columnsIndex := range info.GetGroupByColumnInfo() {
				if test.groupByCol[i] != columnsIndex {
					t.Errorf("groupByColumnStart not equal, expect: %v, actual: %v", test.groupByCol[i], columnsIndex)
				}
			}
		})
	}
}
