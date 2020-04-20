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

func TestMycatShardSimpleDelete(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat`"},
					"db_mycat_1": {"DELETE FROM `tbl_mycat`"},
				},
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat`"},
					"db_mycat_3": {"DELETE FROM `tbl_mycat`"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardDeleteWithWhere(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` WHERE `id`=0"},
				},
			},
		},
		{ // column should be case insensitive
			db:  "db_mycat",
			sql: "delete from tbl_mycat where ID = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` WHERE `ID`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"DELETE FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 2",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat` WHERE `id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 3",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"DELETE FROM `tbl_mycat` WHERE `id`=3"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 4",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` WHERE `id`=4"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 6",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat` WHERE `id`=6"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where tbl_mycat.id = 6",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat` WHERE `tbl_mycat`.`id`=6"}, // table name in assignment is removed
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where db_mycat.tbl_mycat.id = 6",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat` WHERE `db_mycat_2`.`tbl_mycat`.`id`=6"}, // db name in assignment is removed
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id in (1,3,5)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"DELETE FROM `tbl_mycat` WHERE `id` IN (1,5)"},
				},
				"slice-1": {
					"db_mycat_3": {"DELETE FROM `tbl_mycat` WHERE `id` IN (3)"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardDeleteWithOrderBy(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` ORDER BY `id`"},
					"db_mycat_1": {"DELETE FROM `tbl_mycat` ORDER BY `id`"},
				},
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat` ORDER BY `id`"},
					"db_mycat_3": {"DELETE FROM `tbl_mycat` ORDER BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 0 order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` WHERE `id`=0 ORDER BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 0 order by db_mycat.tbl_mycat.id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` WHERE `id`=0 ORDER BY `db_mycat_0`.`tbl_mycat`.`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 0 order by a",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` WHERE `id`=0 ORDER BY `a`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat where id = 0 order by id desc",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` WHERE `id`=0 ORDER BY `id` DESC"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardDeleteWithLimit(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat limit 10",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat` LIMIT 10"},
					"db_mycat_1": {"DELETE FROM `tbl_mycat` LIMIT 10"},
				},
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat` LIMIT 10"},
					"db_mycat_3": {"DELETE FROM `tbl_mycat` LIMIT 10"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "delete from tbl_mycat limit 0, 10",
			hasErr: true, // parse sql error
		},
		{
			db:     "db_mycat",
			sql:    "delete from tbl_mycat limit 10 offset 20",
			hasErr: true, // parse sql error
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}
