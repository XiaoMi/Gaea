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
