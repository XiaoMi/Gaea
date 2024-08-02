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
	"github.com/XiaoMi/Gaea/util"
)

func TestExecuteIn(t *testing.T) {
	planInfo, _ := preparePlanInfo()
	sql := "SELECT * FROM tbl_mycat_murmur WHERE tbl_mycat_murmur.id=5 AND tbl_mycat_murmur.id=4"
	stmt, _ := parser.ParseSQL(sql)
	plan, err := BuildPlan(stmt, nil, "db_mycat", sql, planInfo.rt, planInfo.seqs, nil)
	if err != nil {
		t.Fatalf("build plan error: %v", err)
	}
	ret, err := plan.ExecuteIn(util.NewRequestContext(), nil)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}
	t.Logf("result: %v", ret)
}
