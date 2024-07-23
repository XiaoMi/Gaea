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
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/proxy/sequence"
	"github.com/XiaoMi/Gaea/util"
)

// constants of ShardType
const (
	ShardTypeUnshard = "unshard"
	ShardTypeShard   = "shard"
	ExplainKey       = "explain"
)

// ExplainPlan is the plan for explain statement
type ExplainPlan struct {
	shardType string
	sqls      map[string]map[string][]string
}

func buildExplainPlan(stmt *ast.ExplainStmt, phyDBs map[string]string, db, sql string, r *router.Router, seq *sequence.SequenceManager, hintPlan Plan) (*ExplainPlan, error) {
	stmtToExplain := stmt.Stmt
	if _, ok := stmtToExplain.(*ast.ExplainStmt); ok {
		return nil, fmt.Errorf("nested explain")
	}

	p, err := BuildPlan(stmtToExplain, phyDBs, db, sql, r, seq, hintPlan)
	if err != nil {
		return nil, fmt.Errorf("build plan to explain error: %v", err)
	}

	ep := &ExplainPlan{}

	switch pl := p.(type) {
	case *SelectPlan:
		ep.shardType = ShardTypeShard
		ep.sqls = pl.sqls
		return ep, nil
	case *DeletePlan:
		ep.shardType = ShardTypeShard
		ep.sqls = pl.sqls
		return ep, nil
	case *UpdatePlan:
		ep.shardType = ShardTypeShard
		ep.sqls = pl.sqls
		return ep, nil
	case *InsertPlan:
		ep.shardType = ShardTypeShard
		ep.sqls = pl.sqls
		return ep, nil
	case *UnshardPlan:
		ep.shardType = ShardTypeUnshard
		ep.sqls = make(map[string]map[string][]string)
		dbSQLs := make(map[string][]string)
		dbSQLs[pl.db] = []string{pl.sql}
		ep.sqls[r.GetDefaultRule().GetSlice(0)] = dbSQLs
		return ep, nil
	default:
		return nil, fmt.Errorf("unsupport plan to explain, type: %T", p)
	}
}

// ExecuteIn implement Plan
func (p *ExplainPlan) ExecuteIn(reqCtx *util.RequestContext, se Executor) (*mysql.Result, error) {
	var rows [][]interface{}
	var names = []string{"shard_type", "slice", "db", "sql", "select_type", "table", "partitions", "type",
		"possible_keys", "key", "key_len", "ref", "rows", "filtered", "Extra"}

	for slice, dbSQLs := range p.sqls {
		for db, tableSQLs := range dbSQLs {
			for _, sql := range tableSQLs {
				s, err := se.ExecuteSQL(reqCtx, slice, db, ExplainKey+" "+sql)
				if err != nil {
					return nil, err
				}

				for i := 0; i < len(s.Resultset.Values); i++ {
					row := []interface{}{p.shardType, slice, db, sql}
					for j := 1; j < len(s.Resultset.Values[i]); j++ {
						if s.Resultset.Values[i][j] == nil {
							s.Resultset.Values[i][j] = "NULL"
						}
						row = append(row, s.Resultset.Values[i][j])
					}
					rows = append(rows, row)
				}
			}
		}
	}

	r, _ := mysql.BuildResultset(nil, names, rows)
	ret := mysql.ResultPool.Get()
	ret.Resultset = r

	return ret, nil
}

// Size implement Plan
func (p *ExplainPlan) Size() int {
	return 1
}
