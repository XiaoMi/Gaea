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

func TestMycatShardSimpleInsert(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (0,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (1, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (1,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (2, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (2,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (3, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (3,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (4, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (4,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (6, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (tbl_mycat.id, tbl_mycat.a) values (6, 'hi')", // table name is removed in columns
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into db_mycat.tbl_mycat (db_mycat.tbl_mycat.id, db_mycat.tbl_mycat.a) values (6, 'hi')", // db name is removed in columns, but rewritten in table
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `db_mycat_2`.`tbl_mycat` (`id`,`a`) VALUES (6,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id,id,a) values (6,6,'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`id`,`a`) VALUES (6,6,'hi')"}, // invalid syntax, but gaea does not handle
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (a) values ('hi')",
			hasErr: true, // sharding column not found
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update a = 'hello'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi') ON DUPLICATE KEY UPDATE `a`='hello'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update a = 'hello'+'hi'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi') ON DUPLICATE KEY UPDATE `a`='hello'+'hi'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update a = left('hello',3)",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (6,'hi') ON DUPLICATE KEY UPDATE `a`=LEFT('hello', 3)"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update id = 5",
			hasErr: true, // routing key in update expression
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (id, a) values (6, 'hi') on duplicate key update id = id+3",
			hasErr: true, // routing key in update expression
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat select * from tbl_mycat_child",
			hasErr: true, // select in insert not allowed
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat () values ()",
			hasErr: true, // insert or replace must specify columns
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat values (6, 'hi')",
			hasErr: true, // insert or replace must specify columns
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (id) values (6, 'hi')",
			hasErr: true, // column count doesn't match value count
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (id, a) values (6)",
			hasErr: true, // column count doesn't match value count
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardBatchInsert(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (0, 'hi'), (4, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (0,'hi'),(4,'hi')"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (id, a) values (0, 'hi'), (1, 'hi'), (4, 'hi')",
			hasErr: true, // batch insert has cross slice values
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (id, a) values (6, 'hi'), (5, 'hello')",
			hasErr: true, // batch insert has cross slice values
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardSimpleInsertSet(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat set id = 0, a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` SET `id`=0,`a`='hi'"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSimpleWriteDMLShardMycatMurmur(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_murmur (id, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat_murmur` (`id`,`a`) VALUES (0,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_murmur set a = 'h' where id = 0",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat_murmur` SET `a`='h' WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_murmur where id = 0",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat_murmur` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_murmur (id, a) values (1, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat_murmur` (`id`,`a`) VALUES (1,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_murmur set a = 'h' where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat_murmur` SET `a`='h' WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_murmur where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"DELETE FROM `tbl_mycat_murmur` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_murmur (id, a) values (2, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat_murmur` (`id`,`a`) VALUES (2,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_murmur set a = 'h' where id = 2",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat_murmur` SET `a`='h' WHERE `id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_murmur where id = 2",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"DELETE FROM `tbl_mycat_murmur` WHERE `id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_murmur (id, a) values (3, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat_murmur` (`id`,`a`) VALUES (3,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_murmur (id, a) values (4, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat_murmur` (`id`,`a`) VALUES (4,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_murmur set a = 'h' where id = 4",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat_murmur` SET `a`='h' WHERE `id`=4"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_murmur where id = 4",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat_murmur` WHERE `id`=4"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSimpleWriteDMLShardMycatModLong(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (0,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat_long` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (1, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (1,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"DELETE FROM `tbl_mycat_long` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (256, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (256,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 256",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=256"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 256",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"DELETE FROM `tbl_mycat_long` WHERE `id`=256"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (257, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (257,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 257",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=257"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 257",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"DELETE FROM `tbl_mycat_long` WHERE `id`=257"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (512, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (512,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 512",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=512"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 512",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat_long` WHERE `id`=512"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (513, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (513,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 513",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=513"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 513",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"DELETE FROM `tbl_mycat_long` WHERE `id`=513"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (768, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (768,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 768",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=768"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 768",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"DELETE FROM `tbl_mycat_long` WHERE `id`=768"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_long (id, a) values (769, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"INSERT INTO `tbl_mycat_long` (`id`,`a`) VALUES (769,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "update tbl_mycat_long set a = 'h' where id = 769",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"UPDATE `tbl_mycat_long` SET `a`='h' WHERE `id`=769"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "delete from tbl_mycat_long where id = 769",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"DELETE FROM `tbl_mycat_long` WHERE `id`=769"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatInsertGlobalTable(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat_global_one set id = 0, a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
					"db_mycat_1": {"INSERT INTO `tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
				},
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
					"db_mycat_3": {"INSERT INTO `tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into db_mycat.tbl_mycat_global_one set id = 0, a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `db_mycat_0`.`tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
					"db_mycat_1": {"INSERT INTO `db_mycat_1`.`tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
				},
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `db_mycat_2`.`tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
					"db_mycat_3": {"INSERT INTO `db_mycat_3`.`tbl_mycat_global_one` SET `id`=0,`a`='hi'"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

// 注意这一组各个测试用例之前有关联, 因为都用到了同一个全局序列号
func TestMycatInsertSequenceShardKey(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat set id = nextval(), a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat` SET `id`=1,`a`='hi'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat set id = nextval(), a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"INSERT INTO `tbl_mycat` SET `id`=2,`a`='hi'"}, // next val
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (id, a) values (nextval(), 'hi')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"INSERT INTO `tbl_mycat` (`id`,`a`) VALUES (3,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (ID, a) values (nextval(), 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` (`ID`,`a`) VALUES (4,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat set ID = nextval(), a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"INSERT INTO `tbl_mycat` SET `ID`=5,`a`='hi'"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

// 注意这一组各个测试用例之前有关联, 因为都用到了同一个全局序列号
func TestMycatInsertSequenceUnshardKey(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "insert into tbl_ks set id = 0, user_id = nextval(), a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {"INSERT INTO `tbl_ks_0000` SET `id`=0,`user_id`=1,`a`='hi'"},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "insert into tbl_ks set id = 0, user_id = nextval(), a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {"INSERT INTO `tbl_ks_0000` SET `id`=0,`user_id`=2,`a`='hi'"},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "insert into tbl_ks (id, user_id) values (3,nextval()),(3,nextval()),(3, nextval())",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {"INSERT INTO `tbl_ks_0003` (`id`,`user_id`) VALUES (3,3),(3,4),(3,5)"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestEscapeBackslashShard(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: `insert into tbl_ks (id,name) values (1,'hello\\"world')`,
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {"INSERT INTO `tbl_ks_0001` (`id`,`name`) VALUES (1,'hello\\\\\"world')"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatShardSimpleInsertColumnCaseInsensitive(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat (ID, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` (`ID`,`a`) VALUES (0,'hi')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "insert into tbl_mycat set ID = 0, a = 'hi'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `tbl_mycat` SET `ID`=0,`a`='hi'"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "insert into tbl_mycat (ID, a) values (6, 'hi') on duplicate key update ID = 5",
			hasErr: true, // routing key in update expression
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestInserTableNameColumnCaseInsensitive(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "insert into TBL_MYCAT (ID, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"INSERT INTO `TBL_MYCAT` (`ID`,`a`) VALUES (0,'hi')"},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "insert into tbl_ks_uppercase_child (ID, a) values (0, 'hi')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {"INSERT INTO `tbl_ks_uppercase_child_0000` (`ID`,`a`) VALUES (0,'hi')"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}
