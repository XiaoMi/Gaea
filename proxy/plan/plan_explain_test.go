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

func TestExplainMycatShardSimpleInsert(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (0,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (1, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (1,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (2, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (2,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (3, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (3,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (4, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (4,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (6, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (tbl_mycat.id, tbl_mycat.a) values (6, 'hi')", // table name is removed in columns
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into db_mycat.tbl_mycat (db_mycat.tbl_mycat.id, db_mycat.tbl_mycat.a) values (6, 'hi')", // db name is removed in columns, but rewritten in table
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `db_mycat_2`.`tbl_mycat` (`id`,`a`) VALUES (6,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id,id,a) values (6,6,'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`id`,`a`) VALUES (6,6,'hi')"}, // invalid syntax, but gaea does not handle
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat (a) values ('hi')",
			hasErr: true, // sharding column not found
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update a = 'hello'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi') ON DUPLICATE KEY UPDATE `a`='hello'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update a = 'hello'+'hi'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi') ON DUPLICATE KEY UPDATE `a`='hello'+'hi'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update a = left('hello',3)",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi') ON DUPLICATE KEY UPDATE `a`=LEFT('hello', 3)"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update id = 5",
			hasErr: true, // routing key in update expression
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update id = id+3",
			hasErr: true, // routing key in update expression
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat select * from tbl_mycat_child",
			hasErr: true, // select in insert not allowed
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat () values ()",
			hasErr: true, // insert or replace must specify columns
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat values (6, 'hi')",
			hasErr: true, // insert or replace must specify columns
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat (id) values (6, 'hi')",
			hasErr: true, // column count doesn't match value count
		},
		{
			db:     "db_mycat",
			sql:    "explain insert into tbl_mycat (id, a) values (6)",
			hasErr: true, // column count doesn't match value count
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestExplainUnshardInsert(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "explain insert into tbl_unshard (id, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_unshard` (`id`,`a`) VALUES (0,'hi')"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestExplainUnshardInsertWithDb(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "explain insert into db_mycat.tbl_unshard (id, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `db_mycat_0`.`tbl_unshard` (`id`,`a`) VALUES (0,'hi')"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}
