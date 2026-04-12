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
)

func buildHavingSelectPlan(t *testing.T, sql string) (*SelectPlan, int) {
	t.Helper()

	info, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	stmt, err := parser.ParseSQL(sql)
	if err != nil {
		t.Fatalf("parse sql error: %v", err)
	}

	sel, ok := stmt.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected parsed SelectStmt, got %T", stmt)
	}
	originalFieldCount := len(sel.Fields.Fields)

	plan, err := BuildPlan(stmt, info.phyDBs, "db_mycat", sql, info.rt, info.seqs, nil)
	if err != nil {
		t.Fatalf("BuildPlan error: %v", err)
	}

	sp, ok := plan.(*SelectPlan)
	if !ok {
		t.Fatalf("expected SelectPlan, got %T", plan)
	}

	return sp, originalFieldCount
}

func assertHavingPulledUp(t *testing.T, sp *SelectPlan) {
	t.Helper()

	if sp.havingExpr == nil {
		t.Fatalf("expected havingExpr to be set")
	}
	if sp.stmt.Having != nil {
		t.Fatalf("expected stmt.Having to be removed for global re-filter")
	}
}

func assertHiddenHavingSumField(t *testing.T, fields []*ast.SelectField) {
	t.Helper()

	for _, field := range fields {
		agg, ok := field.Expr.(*ast.AggregateFuncExpr)
		if !ok || strings.ToLower(agg.F) != "sum" || len(agg.Args) != 1 {
			continue
		}

		var tableName, colName string
		switch arg := agg.Args[0].(type) {
		case *ast.ColumnNameExpr:
			tableName = arg.Name.Table.L
			colName = arg.Name.Name.L
		case *ColumnNameExprDecorator:
			tableName = arg.ColumnNameExpr.Name.Table.L
			colName = arg.ColumnNameExpr.Name.Name.L
		default:
			t.Fatalf("expected SUM argument to be column name, got %T", agg.Args[0])
		}
		if tableName != "t" || colName != "b" {
			continue
		}
		if !field.Auxiliary {
			t.Fatalf("expected HAVING aggregate field to be auxiliary")
		}
		return
	}

	t.Fatalf("expected appended HAVING aggregate field SUM(t.b)")
}

func TestHandleHaving_QualifiedAggregateTriggersGlobal(t *testing.T) {
	sql := "select t.a, sum(t.b) as total_b from tbl_mycat as t group by t.a having sum(t.b) > 50"

	sp, _ := buildHavingSelectPlan(t, sql)

	assertHavingPulledUp(t, sp)
}

func TestHandleHaving_AddsQualifiedAggregateFieldWhenMissing(t *testing.T) {
	sql := "select t.a from tbl_mycat t group by t.a having sum(t.b) < 20"

	sp, originalFieldCount := buildHavingSelectPlan(t, sql)

	assertHavingPulledUp(t, sp)
	if len(sp.stmt.Fields.Fields) <= originalFieldCount {
		t.Fatalf("expected extra aggregate field appended for HAVING, fields: %d, original: %d", len(sp.stmt.Fields.Fields), originalFieldCount)
	}
	assertHiddenHavingSumField(t, sp.stmt.Fields.Fields[originalFieldCount:])
}

func TestHandleHaving_AggregateAliasReferenceTriggersGlobal(t *testing.T) {
	sql := "select t.a as group_name, sum(t.b) as total_val from tbl_mycat t group by t.a having total_val = 80"

	sp, _ := buildHavingSelectPlan(t, sql)

	assertHavingPulledUp(t, sp)
}
