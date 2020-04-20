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

import "testing"

func TestMycatShardSimpleUpdate(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi'"},
					"db_mycat_1": {"UPDATE `tbl_mycat` SET `a`='hi'"},
				},
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat` SET `a`='hi'"},
					"db_mycat_3": {"UPDATE `tbl_mycat` SET `a`='hi'"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat set id = 5",
			hasErr: true, // cannot update shard column value
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat set ID = 5",
			hasErr: true, // cannot update shard column value
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where ID = 5",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `ID`=5"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat, tbl_mycat_child set id = 5",
			hasErr: true, // does not support update multiple tables in sharding
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardUpdateColumnCaseInsensitive(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat set ID = 5",
			hasErr: true, // cannot update shard column value
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where ID = 5",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `ID`=5"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardUpdateWithWhere(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 2",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 3",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=3"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 4",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=4"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 6",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=6"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set tbl_mycat.a = 'hi' where tbl_mycat.id = 6",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `tbl_mycat`.`id`=6"}, // table name in assignment is removed
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set db_mycat.tbl_mycat.a = 'hi' where db_mycat.tbl_mycat.id = 6",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `db_mycat_2`.`tbl_mycat`.`id`=6"}, // db name in assignment is removed
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id in (1,3,5)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id` IN (1,5)"},
				},
				"slice-1": {
					"db_mycat_3": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id` IN (3)"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat set id = 5 where id = 6",
			hasErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardUpdateWithOrderBy(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' ORDER BY `id`"},
					"db_mycat_1": {"UPDATE `tbl_mycat` SET `a`='hi' ORDER BY `id`"},
				},
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat` SET `a`='hi' ORDER BY `id`"},
					"db_mycat_3": {"UPDATE `tbl_mycat` SET `a`='hi' ORDER BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 0 order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=0 ORDER BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 0 order by db_mycat.tbl_mycat.id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=0 ORDER BY `db_mycat_0`.`tbl_mycat`.`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 0 order by a",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=0 ORDER BY `a`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' where id = 0 order by id desc",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' WHERE `id`=0 ORDER BY `id` DESC"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardUpdateWithLimit(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "update tbl_mycat set a = 'hi' limit 10",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat` SET `a`='hi' LIMIT 10"},
					"db_mycat_1": {"UPDATE `tbl_mycat` SET `a`='hi' LIMIT 10"},
				},
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat` SET `a`='hi' LIMIT 10"},
					"db_mycat_3": {"UPDATE `tbl_mycat` SET `a`='hi' LIMIT 10"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat set a = 'hi' limit 0, 10",
			hasErr: true, // parse sql error: line 1 column 45 near "offset 0"
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat set a = 'hi' limit 10 offset 20",
			hasErr: true, // parse sql error: line 1 column 45 near "offset 20"
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestGlobalTableUpdate(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_global_one set a = 'hi' limit 10",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
					"db_mycat_1": {"UPDATE `tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
				},
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
					"db_mycat_3": {"UPDATE `tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update db_mycat.tbl_mycat_global_one set db_mycat.tbl_mycat_global_one.a = 'hi' limit 10",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `db_mycat_0`.`tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
					"db_mycat_1": {"UPDATE `db_mycat_1`.`tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
				},
				"slice-1": {
					"db_mycat_2": {"UPDATE `db_mycat_2`.`tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
					"db_mycat_3": {"UPDATE `db_mycat_3`.`tbl_mycat_global_one` SET `a`='hi' LIMIT 10"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat_global_one set a = 'hi' limit 0, 10",
			hasErr: true, // parse sql error: line 1 column 45 near "offset 0"
		},
		{
			db:     "db_mycat",
			sql:    "update tbl_mycat_global_one set a = 'hi' limit 10 offset 20",
			hasErr: true, // parse sql error: line 1 column 45 near "offset 20"
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}
