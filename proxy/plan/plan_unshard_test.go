// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/util"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/XiaoMi/Gaea/backend"
)

func TestUnshardPlan(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}
	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: `select * from tbl_unshard_a as a join db_mycat.tbl_unshard_b as b on a.id = b.id`,
			sqls: map[string]map[string][]string{
				backend.DefaultSlice: {
					"db_mycat_0": {"SELECT * FROM `tbl_unshard_a` AS `a` JOIN `db_mycat_0`.`tbl_unshard_b` AS `b` ON `a`.`id`=`b`.`id`"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestUnshardPlanWithoutDB(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}
	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: `select * from tbl_unshard limit 10`,
			sqls: map[string]map[string][]string{
				backend.DefaultSlice: {
					"db_mycat_0": {"SELECT * FROM `tbl_unshard` LIMIT 10"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectInsertID(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}
	tests := []SQLTestcase{
		{
			db:   "db_mycat",
			sql:  `select last_insert_id()`,
			sqls: map[string]map[string][]string{},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestPreCreateUnshardPlan(t *testing.T) {
	sql := "select 1"
	p, err := PreCreateUnshardPlan(sql, map[string]string{}, "test")
	assert.Equal(t, err, nil)
	assert.Equal(t, sql, p.sql)
	p, err = PreCreateUnshardPlan(sql, map[string]string{"test": "test_0"}, "test")
	assert.NotEqual(t, err, nil)
}

type mockExecutor struct {
}

func (*mockExecutor) ExecuteSQL(ctx *util.RequestContext, slice, db, sql string) (*mysql.Result, error) {
	return mysql.ResultPool.Get(), nil
}

func (*mockExecutor) ExecuteSQLs(*util.RequestContext, map[string]map[string][]string) ([]*mysql.Result, error) {
	return make([]*mysql.Result, 0), nil
}

func (*mockExecutor) SetLastInsertID(uint64) {
}

func (*mockExecutor) GetLastInsertID() uint64 {
	return 0
}

func (*mockExecutor) HandleSet(*util.RequestContext, string, *ast.SetStmt) (*mysql.Result, error) {
	return mysql.ResultPool.Get(), nil
}

func TestBuildPlan(t *testing.T) {
	tests := []struct {
		sql        string
		expectPlan Plan
	}{
		{
			sql:        `select last_insert_id()`,
			expectPlan: &SelectLastInsertIDPlan{},
		},
		{
			sql:        `set session read_only=1`,
			expectPlan: &SetPlan{},
		},
		{
			sql:        `select 1`,
			expectPlan: &UnshardPlan{},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, func(t *testing.T) {
			stmt, err := parser.ParseSQL(test.sql)
			assert.Equal(t, err, nil)
			p, err := BuildPlan(stmt, map[string]string{}, "test", test.sql, nil, nil, nil)
			assert.Equal(t, reflect.TypeOf(test.expectPlan), reflect.TypeOf(p))
			assert.Equal(t, err, nil)
			_, err = p.ExecuteIn(util.NewRequestContext(), &mockExecutor{})
			assert.Equal(t, err, nil)
		})
	}

}

func TestIgnorePlan(t *testing.T) {
	p := CreateIgnorePlan()
	_, err := p.ExecuteIn(util.NewRequestContext(), &mockExecutor{})
	assert.Equal(t, err, nil)
}

func newMockRouter() *router.Router {
	r, _ := router.NewRouter(&models.Namespace{
		DefaultSlice: "slice-0",
		Slices: []*models.Slice{
			{
				Name: "slice-0",
			},
		},
	})
	return r
}

var (
	testTokensSelect = []string{"SELECT", "db.test_shard_hash"}
	testTokensInsert = []string{"INSERT", "db.test_shard_hash"}
	testTokensUpdate = []string{"UPDATE", "t", "SET", "a=1"}
)

func TestCheckUnshardBase(t *testing.T) {
	type args struct {
		tokenId int
		tokens  []string
		rt      *router.Router
		db      string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "Test with valid token and shard rule",
			args: args{
				tokenId: 1,
				tokens:  testTokensSelect,
				rt:      newMockRouter(),
				db:      "db",
			},
			want:  "db",
			want1: true,
		},
		{
			name: "Test without shard rule",
			args: args{
				tokenId: 1,
				tokens:  testTokensSelect,
				rt:      newMockRouter(),
				db:      "other_db",
			},
			want:  "other_db",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CheckUnshardBase(tt.args.tokenId, tt.args.tokens, tt.args.rt, tt.args.db)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestCheckUnshardInsert(t *testing.T) {
	type args struct {
		tokens []string
		rt     *router.Router
		db     string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "Test with valid insert token and shard rule",
			args: args{
				tokens: testTokensInsert,
				rt:     newMockRouter(),
				db:     "db",
			},
			want:  "db",
			want1: true,
		},
		{
			name: "Test without shard rule",
			args: args{
				tokens: testTokensInsert,
				rt:     newMockRouter(),
				db:     "other_db",
			},
			want:  "other_db",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CheckUnshardInsert(tt.args.tokens, tt.args.rt, tt.args.db)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestCheckUnshardUpdate(t *testing.T) {
	type args struct {
		tokens []string
		rt     *router.Router
		db     string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "Test with valid update token and shard rule",
			args: args{
				tokens: testTokensUpdate,
				rt:     newMockRouter(),
				db:     "db",
			},
			want:  "db",
			want1: true,
		},
		{
			name: "Test without shard rule",
			args: args{
				tokens: testTokensUpdate,
				rt:     newMockRouter(),
				db:     "other_db",
			},
			want:  "other_db",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CheckUnshardUpdate(tt.args.tokens, tt.args.rt, tt.args.db)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
