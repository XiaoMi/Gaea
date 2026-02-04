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
	"testing"

	"github.com/XiaoMi/Gaea/parser"
)

func TestHandleHaving_AggregateAliasTriggersGlobal(t *testing.T) {
	info, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	sql := "select a, sum(b) as total_b from tbl_mycat group by a having total_b > 40"
	stmt, err := parser.ParseSQL(sql)
	if err != nil {
		t.Fatalf("parse sql error: %v", err)
	}

	plan, err := BuildPlan(stmt, info.phyDBs, "db_mycat", sql, info.rt, info.seqs, nil)
	if err != nil {
		t.Fatalf("BuildPlan error: %v", err)
	}

	sp, ok := plan.(*SelectPlan)
	if !ok {
		t.Fatalf("expected SelectPlan, got %T", plan)
	}

	if sp.havingExpr == nil {
		t.Fatalf("expected havingExpr to be set")
	}
	if sp.stmt.Having != nil {
		t.Fatalf("expected stmt.Having to be removed for global re-filter")
	}
}

func TestHandleHaving_AddsAggregateFieldWhenMissing(t *testing.T) {
	info, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	sql := "select a from tbl_mycat group by a having sum(b) < 20"
	stmt, err := parser.ParseSQL(sql)
	if err != nil {
		t.Fatalf("parse sql error: %v", err)
	}

	plan, err := BuildPlan(stmt, info.phyDBs, "db_mycat", sql, info.rt, info.seqs, nil)
	if err != nil {
		t.Fatalf("BuildPlan error: %v", err)
	}

	sp, ok := plan.(*SelectPlan)
	if !ok {
		t.Fatalf("expected SelectPlan, got %T", plan)
	}

	// original SELECT has 1 column; we expect an extra SUM(b) field appended
	if len(sp.stmt.Fields.Fields) < 2 {
		t.Fatalf("expected extra aggregate field appended for HAVING")
	}
}
