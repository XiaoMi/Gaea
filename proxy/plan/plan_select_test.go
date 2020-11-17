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

	"github.com/XiaoMi/Gaea/proxy/router"
)

func TestSimpleSelectShardMycatMod(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id from tbl_mycat where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id` FROM `tbl_mycat` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select tbl_mycat.id from tbl_mycat where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `tbl_mycat`.`id` FROM `tbl_mycat` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where tbl_mycat.id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `tbl_mycat`.`id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from db_mycat.tbl_mycat where db_mycat.tbl_mycat.id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `db_mycat_0`.`tbl_mycat` WHERE `db_mycat_0`.`tbl_mycat`.`id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 2",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 3",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`=3"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 4",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=4"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id in (0,1,2,3,4,6)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (0,4)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (1)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (2,6)"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (3)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where k = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `k`=0"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSimpleSelectShardMycatMurmur(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 0",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 2",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 3",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=3"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 4",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=4"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id in (0,1,2,3,4,6)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (6)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (1,2,3)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (0,4)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where k = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSimpleSelectShardMycatMurmur_ShardKeyTypeString(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id in ('0')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN ('0')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = '0'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`='0'"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSimpleSelectShardMycatString(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_string where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_string` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_string where k = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_string` WHERE `k`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_string` WHERE `k`=0"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_string` WHERE `k`=0"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat_string` WHERE `k`=0"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSimpleSelectShardMycatLong(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 2",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 3",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=3"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 256",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=256"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 512",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=512"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 768",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=768"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 769",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=769"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 1024",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=1024"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 1025",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=1025"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id in (0,1,256,257,512,513,768,769,1024,1025)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id` IN (0,1,1024,1025)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_long` WHERE `id` IN (256,257)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_long` WHERE `id` IN (512,513)"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat_long` WHERE `id` IN (768,769)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where k = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `k`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_long` WHERE `k`=0"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_long` WHERE `k`=0"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat_long` WHERE `k`=0"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectMultiTablesEQ(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child`"},
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child`"},
					"db_mycat_3": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child`"},
				},
			},
		},

		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where tbl_mycat.id = COUNT(1)", // 一边是列名, 另一边不支持, 则只替换列名
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=COUNT(1)"},
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=COUNT(1)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=COUNT(1)"},
					"db_mycat_3": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=COUNT(1)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where tbl_mycat.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where tbl_mycat.id is null",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id` IS NULL"},
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id` IS NULL"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id` IS NULL"},
					"db_mycat_3": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id` IS NULL"},
				},
			},
		},
		// TODO: 分表列是否需要支持等值比较NULL
		//{
		//	db:     "db_mycat",
		//	sql:    "select * from tbl_mycat, tbl_mycat_child where tbl_mycat.id = null",
		//	hasErr: true,
		//},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where 1 = 1 and tbl_mycat.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE 1=1 AND `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat join tbl_mycat_child on tbl_mycat.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` ON `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where tbl_mycat_child.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat_child`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat join tbl_mycat_child on tbl_mycat_child.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` ON `tbl_mycat_child`.`id`=1"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat, tbl_mycat_child where id = 1", // id is ambiguous that both sharding key in tbl_mycat and tbl_mycat_child
			hasErr: true,
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat join tbl_mycat_child on id = 1", // id is ambiguous that both sharding key in tbl_mycat and tbl_mycat_child
			hasErr: true,
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_user_child where id = 1", // id is not ambiguous
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_user_child` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat join tbl_mycat_user_child on id = 1", // id is not ambiguous
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_user_child` ON `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat join tbl_mycat_child on tbl_mycat_child.id = 1 where tbl_mycat.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` ON `tbl_mycat_child`.`id`=1 WHERE `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:   "db_mycat",
			sql:  "select * from tbl_mycat join tbl_mycat_child on tbl_mycat_child.id = 1 where tbl_mycat.id = 0",
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat join tbl_mycat_child on tbl_mycat_child.id = 1 or tbl_mycat.id = 2",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` ON `tbl_mycat_child`.`id`=1 OR `tbl_mycat`.`id`=2"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` ON `tbl_mycat_child`.`id`=1 OR `tbl_mycat`.`id`=2"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where tbl_mycat_child.id = 1 or tbl_mycat.id = 2",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat_child`.`id`=1 OR `tbl_mycat`.`id`=2"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat_child`.`id`=1 OR `tbl_mycat`.`id`=2"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat join tbl_mycat_murmur on tbl_mycat.id = 1", // tables have different route
			hasErr: true,
		},

		// expr.R
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where COUNT(1) = tbl_mycat.id", // 一边是列名, 另一边不支持, 则只替换列名
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE COUNT(1)=`tbl_mycat`.`id`"},
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE COUNT(1)=`tbl_mycat`.`id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE COUNT(1)=`tbl_mycat`.`id`"},
					"db_mycat_3": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE COUNT(1)=`tbl_mycat`.`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where 1 = tbl_mycat.id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE 1=`tbl_mycat`.`id`"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat, tbl_mycat_child where 1 = a.id",
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectJoinUsing(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat join tbl_mycat_child using(id)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` USING (`id`)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` USING (`id`)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` USING (`id`)"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` USING (`id`)"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat join tbl_mycat_child using(tbl_mycat.id)",
			hasErr: true,
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat join tbl_mycat_child using(db_mycat.tbl_mycat.id)",
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectMultiTablesAlias(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat as a, tbl_mycat_child where a.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat` AS `a`) JOIN `tbl_mycat_child` WHERE `a`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat join tbl_mycat_child as b on b.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` JOIN `tbl_mycat_child` AS `b` ON `b`.`id`=1"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat as a join tbl_mycat_child as a on a.id = 1",
			hasErr: true,
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat as a join tbl_mycat_child on b.id = 1",
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestKingshardSelectAlias(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select a.ss, a from tbl_ks as a where a.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `a`.`ss`,`a` FROM `tbl_ks_0001` AS `a` WHERE `a`.`id`=1",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestKingshardSelectBetweenAlias(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select  name from tbl_ks as a where a.id between 10 and 100",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `name` FROM `tbl_ks_0000` AS `a` WHERE `a`.`id` BETWEEN 10 AND 100",
						"SELECT `name` FROM `tbl_ks_0001` AS `a` WHERE `a`.`id` BETWEEN 10 AND 100",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT `name` FROM `tbl_ks_0002` AS `a` WHERE `a`.`id` BETWEEN 10 AND 100",
						"SELECT `name` FROM `tbl_ks_0003` AS `a` WHERE `a`.`id` BETWEEN 10 AND 100",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectColumnCaseInsensitive(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select a.ss, a from tbl_ks as a where a.ID = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `a`.`ss`,`a` FROM `tbl_ks_0001` AS `a` WHERE `a`.`ID`=1",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select a.ss, a from tbl_ks as a where 1 = a.ID",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `a`.`ss`,`a` FROM `tbl_ks_0001` AS `a` WHERE 1=`a`.`ID`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where CREATE_TIME between '2014-09-05 00:00:00' and '2014-09-07 00:00:00'", // 2014-09-01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140905` WHERE `CREATE_TIME` BETWEEN '2014-09-05 00:00:00' AND '2014-09-07 00:00:00'",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140907` WHERE `CREATE_TIME` BETWEEN '2014-09-05 00:00:00' AND '2014-09-07 00:00:00'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where ID in (0,2)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `ID` IN (0)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `ID` IN (2)"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectTableNameCaseInsensitive(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select a.ss, a from TBL_KS_UPPERCASE as a where a.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `a`.`ss`,`a` FROM `TBL_KS_UPPERCASE_0001` AS `a` WHERE `a`.`id`=1",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select ss, a from TBL_KS_UPPERCASE where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `ss`,`a` FROM `TBL_KS_UPPERCASE_0001` WHERE `id`=1",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select a.ss, a from tbl_ks_uppercase as a where a.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `a`.`ss`,`a` FROM `tbl_ks_uppercase_0001` AS `a` WHERE `a`.`id`=1",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select a.ss, a from tbl_ks_uppercase_child as a where a.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `a`.`ss`,`a` FROM `tbl_ks_uppercase_child_0001` AS `a` WHERE `a`.`id`=1",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from TBL_MYCAT where TBL_MYCAT.ID in (0,2)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `TBL_MYCAT` WHERE `TBL_MYCAT`.`ID` IN (0)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `TBL_MYCAT` WHERE `TBL_MYCAT`.`ID` IN (2)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat as A where A.ID in (0,2)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` AS `A` WHERE `A`.`ID` IN (0)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` AS `A` WHERE `A`.`ID` IN (2)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from (select id, ss from tbl_mycat) A where A.ID in (0,2)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (SELECT `id`,`ss` FROM (`tbl_mycat`)) AS `A` WHERE `A`.`ID` IN (0)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM (SELECT `id`,`ss` FROM (`tbl_mycat`)) AS `A` WHERE `A`.`ID` IN (2)"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

// TODO: range shard
func TestMycatSelectBinaryOperatorComparison(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 = id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1=`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id > 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`>1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`>1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`>1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`>1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 < id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE 1<`id`"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1<`id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE 1<`id`"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE 1<`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id >= 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`>=1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`>=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`>=1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`>=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 <= id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE 1<=`id`"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1<=`id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE 1<=`id`"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE 1<=`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id < 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`<1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`<1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`<1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`<1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 > id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE 1>`id`"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1>`id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE 1>`id`"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE 1>`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id <= 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`<=1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`<=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`<=1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`<=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 >= id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE 1>=`id`"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1>=`id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE 1>=`id`"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE 1>=`id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id <> 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id != 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`!=1"},
				},
			},
		},

		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where tbl_mycat.id = tbl_mycat_child.id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=`tbl_mycat_child`.`id`"},
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=`tbl_mycat_child`.`id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=`tbl_mycat_child`.`id`"},
					"db_mycat_3": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `tbl_mycat`.`id`=`tbl_mycat_child`.`id`"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectPatternIn(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id in (0,1,2,3,4,5,6,7)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (0,4)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (1,5)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (2,6)"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (3,7)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id in (0,0,0,1)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (0,0,0)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (1)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id not in (1)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` NOT IN (1)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` NOT IN (1)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id` NOT IN (1)"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id` NOT IN (1)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 in (1,2,3) and id=1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1 IN (1,2,3) AND `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 in (1,2,3) or id=1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE 1 IN (1,2,3) OR `id`=1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1 IN (1,2,3) OR `id`=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE 1 IN (1,2,3) OR `id`=1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE 1 IN (1,2,3) OR `id`=1"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat where id in (1 = 1)",
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectPatternInWithFuncDatabase(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where database() in ('db_mycat_0', 'db_mycat_1') and id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1') AND `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where database() not in ('db_mycat_0', 'db_mycat_1') and id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1') AND `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where database() not in ('db_mycat_0', 'db_mycat_1') or id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1') OR `id`=1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1') OR `id`=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1') OR `id`=1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1') OR `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where database() in ('db_mycat_0', 'db_mycat_1')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1')"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1')"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1')"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where database() not in ('db_mycat_0', 'db_mycat_1')",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1')"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1')"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1')"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() NOT IN ('db_mycat_0','db_mycat_1')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where database() in ('db_mycat_0', 'db_mycat_1') and id in (0,1,2,3,4,5,6,7)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1') AND `id` IN (0,4)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1') AND `id` IN (1,5)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1') AND `id` IN (2,6)"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_0','db_mycat_1') AND `id` IN (3,7)"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

// TODO: range shard
func TestMycatSelectPatternBetween(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id between 1 and 5",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 1 AND 5"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 1 AND 5"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 1 AND 5"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 1 AND 5"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id between 5 and 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 5 AND 1"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 5 AND 1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 5 AND 1"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id` BETWEEN 5 AND 1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where user between 'curry' and 'durant'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `user` BETWEEN 'curry' AND 'durant'"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `user` BETWEEN 'curry' AND 'durant'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `user` BETWEEN 'curry' AND 'durant'"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `user` BETWEEN 'curry' AND 'durant'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 between 2 and 3",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE 1 BETWEEN 2 AND 3"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1 BETWEEN 2 AND 3"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE 1 BETWEEN 2 AND 3"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE 1 BETWEEN 2 AND 3"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where 1 not between 2 and 3",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE 1 NOT BETWEEN 2 AND 3"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE 1 NOT BETWEEN 2 AND 3"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE 1 NOT BETWEEN 2 AND 3"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE 1 NOT BETWEEN 2 AND 3"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectPatternBetweenWithFuncDatabase(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where database() between 'db_mycat_0' and 'db_mycat_2'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() BETWEEN 'db_mycat_0' AND 'db_mycat_2'"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() BETWEEN 'db_mycat_0' AND 'db_mycat_2'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() BETWEEN 'db_mycat_0' AND 'db_mycat_2'"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE DATABASE() BETWEEN 'db_mycat_0' AND 'db_mycat_2'"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectPatternLike(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id like '1%'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` LIKE '1%'"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` LIKE '1%'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id` LIKE '1%'"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id` LIKE '1%'"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectPatternLogicOperator(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 0 or id in (1,2)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 OR `id` IN ()"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 OR `id` IN (1)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 OR `id` IN (2)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 0 and user = 'curry'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 AND `user`='curry'"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 0 or user = 'curry'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 OR `user`='curry'"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 OR `user`='curry'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 OR `user`='curry'"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`=0 OR `user`='curry'"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

// TODO: need more testcases
func TestMycatSelectSubqueryInTableRefs(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select id from (select user from tbl_mycat) as a", // table in subquery must be a sharding table
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id` FROM (SELECT `user` FROM (`tbl_mycat`)) AS `a`"},
					"db_mycat_1": {"SELECT `id` FROM (SELECT `user` FROM (`tbl_mycat`)) AS `a`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id` FROM (SELECT `user` FROM (`tbl_mycat`)) AS `a`"},
					"db_mycat_3": {"SELECT `id` FROM (SELECT `user` FROM (`tbl_mycat`)) AS `a`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id from (select user from tbl_mycat_unknown) as a", //unshard plan
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id` FROM (SELECT `user` FROM (`tbl_mycat_unknown`)) AS `a`"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select id from tbl_mycat as a, (select user from tbl_mycat) as a", // table alias is ambiguous
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectAggregationFunctionCount(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select count(*) from tbl_mycat",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT COUNT(1) FROM `tbl_mycat`"},
					"db_mycat_1": {"SELECT COUNT(1) FROM `tbl_mycat`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT COUNT(1) FROM `tbl_mycat`"},
					"db_mycat_3": {"SELECT COUNT(1) FROM `tbl_mycat`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select count(id) from tbl_mycat",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT COUNT(`id`) FROM `tbl_mycat`"},
					"db_mycat_1": {"SELECT COUNT(`id`) FROM `tbl_mycat`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT COUNT(`id`) FROM `tbl_mycat`"},
					"db_mycat_3": {"SELECT COUNT(`id`) FROM `tbl_mycat`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select count(user) from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT COUNT(`user`) FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select count(user) from tbl_mycat where user = 'curry'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT COUNT(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_1": {"SELECT COUNT(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT COUNT(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_3": {"SELECT COUNT(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectAggregationFunctionMax(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:     "db_mycat",
			sql:    "select max(*) from tbl_mycat",
			hasErr: true, // max(*) is invalid syntax
		},
		{
			db:  "db_mycat",
			sql: "select max(id) from tbl_mycat",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT MAX(`id`) FROM `tbl_mycat`"},
					"db_mycat_1": {"SELECT MAX(`id`) FROM `tbl_mycat`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT MAX(`id`) FROM `tbl_mycat`"},
					"db_mycat_3": {"SELECT MAX(`id`) FROM `tbl_mycat`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select max(user) from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT MAX(`user`) FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select max(user) from tbl_mycat where user = 'curry'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT MAX(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_1": {"SELECT MAX(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT MAX(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_3": {"SELECT MAX(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectAggregationFunctionMin(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:     "db_mycat",
			sql:    "select min(*) from tbl_mycat",
			hasErr: true, // min(*) is invalid syntax
		},
		{
			db:  "db_mycat",
			sql: "select min(id) from tbl_mycat",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT MIN(`id`) FROM `tbl_mycat`"},
					"db_mycat_1": {"SELECT MIN(`id`) FROM `tbl_mycat`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT MIN(`id`) FROM `tbl_mycat`"},
					"db_mycat_3": {"SELECT MIN(`id`) FROM `tbl_mycat`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select min(user) from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT MIN(`user`) FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select min(user) from tbl_mycat where user = 'curry'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT MIN(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_1": {"SELECT MIN(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT MIN(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_3": {"SELECT MIN(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectAggregationFunctionSum(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:     "db_mycat",
			sql:    "select sum(*) from tbl_mycat",
			hasErr: true, // sum(*) is invalid syntax
		},
		{
			db:  "db_mycat",
			sql: "select sum(id) from tbl_mycat",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT SUM(`id`) FROM `tbl_mycat`"},
					"db_mycat_1": {"SELECT SUM(`id`) FROM `tbl_mycat`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT SUM(`id`) FROM `tbl_mycat`"},
					"db_mycat_3": {"SELECT SUM(`id`) FROM `tbl_mycat`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select sum(user) from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT SUM(`user`) FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select sum(user) from tbl_mycat where user = 'curry'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT SUM(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_1": {"SELECT SUM(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT SUM(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
					"db_mycat_3": {"SELECT SUM(`user`) FROM `tbl_mycat` WHERE `user`='curry'"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectGroupBy(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat group by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` GROUP BY `id`"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` GROUP BY `id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` GROUP BY `id`"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` GROUP BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, count(user) from tbl_mycat group by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,COUNT(`user`) FROM `tbl_mycat` GROUP BY `id`"},
					"db_mycat_1": {"SELECT `id`,COUNT(`user`) FROM `tbl_mycat` GROUP BY `id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,COUNT(`user`) FROM `tbl_mycat` GROUP BY `id`"},
					"db_mycat_3": {"SELECT `id`,COUNT(`user`) FROM `tbl_mycat` GROUP BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, count(user) from tbl_mycat where id = 1 group by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT `id`,COUNT(`user`) FROM `tbl_mycat` WHERE `id`=1 GROUP BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select user, count(id) from tbl_mycat where user = 'curry' group by user",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `user`,COUNT(`id`) FROM `tbl_mycat` WHERE `user`='curry' GROUP BY `user`"},
					"db_mycat_1": {"SELECT `user`,COUNT(`id`) FROM `tbl_mycat` WHERE `user`='curry' GROUP BY `user`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `user`,COUNT(`id`) FROM `tbl_mycat` WHERE `user`='curry' GROUP BY `user`"},
					"db_mycat_3": {"SELECT `user`,COUNT(`id`) FROM `tbl_mycat` WHERE `user`='curry' GROUP BY `user`"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectHaving(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat having id = 1", // note: does not calculate route in having clause
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` HAVING `id`=1"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` HAVING `id`=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` HAVING `id`=1"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` HAVING `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, count(user) from tbl_mycat where id=1 group by id having count(user) > 5",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT `id`,COUNT(`user`) FROM `tbl_mycat` WHERE `id`=1 GROUP BY `id` HAVING COUNT(`user`)>5"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectOrderBy(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id`"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id`"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat where id = 1 order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` WHERE `id`=1 ORDER BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat where id in (0,1,2,3,4,5,6,7) order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` WHERE `id` IN (0,4) ORDER BY `id`"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` WHERE `id` IN (1,5) ORDER BY `id`"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` WHERE `id` IN (2,6) ORDER BY `id`"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` WHERE `id` IN (3,7) ORDER BY `id`"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat order by id desc",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id` DESC"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id` DESC"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id` DESC"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` ORDER BY `id` DESC"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestMycatSelectLimit(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat limit 10",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 10"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 10"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 10"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 10"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat limit 0, 10",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 0,10"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 0,10"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 0,10"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 0,10"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id, user from tbl_mycat limit 10, 10",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 20"},
					"db_mycat_1": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 20"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 20"},
					"db_mycat_3": {"SELECT `id`,`user` FROM `tbl_mycat` LIMIT 20"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectMycatMultiTablesDatabaseHint(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		// database function is left
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where DATABASE() = `db_mycat_0` and tbl_mycat.id = 1", // hint is column name
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE DATABASE()=`db_mycat_0` AND `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where DATABASE() = 'db_mycat_0' and tbl_mycat.id = 1", // hint is value
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE DATABASE()='db_mycat_0' AND `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where DB() = `db_mycat_0` and tbl_mycat.id = 1", // not DATABASE hint, use origin route
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE DB()=`db_mycat_0` AND `tbl_mycat`.`id`=1"},
				},
			},
		},

		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat, tbl_mycat_child where DATABASE() = count(1) and tbl_mycat.id = 1", // hint must be a value or column name
			hasErr: true,
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat, tbl_mycat_child where DATABASE() = 'db_mycat_10000' and tbl_mycat.id = 1", // phy db not found
			hasErr: true,
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks where DATABASE() = 'db_ks' and id = 1", // only mycat route support database hint
			hasErr: true,
		},
		// database function is right
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where `db_mycat_0` = DATABASE() and tbl_mycat.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `db_mycat_0`=DATABASE() AND `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where `db_mycat_0` = DB() and tbl_mycat.id = 1", // not DATABASE hint, use origin route
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE `db_mycat_0`=DB() AND `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_child where DATABASE() = 'db_mycat_0' and tbl_mycat.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_child` WHERE DATABASE()='db_mycat_0' AND `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat, tbl_mycat_child where count(1) = DATABASE() and tbl_mycat.id = 1", // hint must be a value or column name
			hasErr: true,
		},
		{
			db:     "db_mycat",
			sql:    "select * from tbl_mycat, tbl_mycat_child where 'db_mycat_10000' = DATABASE() and tbl_mycat.id = 1", // phy db not found
			hasErr: true,
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks where 'db_ks' = DATABASE() and id = 1", // only mycat route support database hint
			hasErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectMultiTablesKingshard(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select * from tbl_ks, tbl_ks_child",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM (`tbl_ks_0000`) JOIN `tbl_ks_child_0000`",
						"SELECT * FROM (`tbl_ks_0001`) JOIN `tbl_ks_child_0001`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM (`tbl_ks_0002`) JOIN `tbl_ks_child_0002`",
						"SELECT * FROM (`tbl_ks_0003`) JOIN `tbl_ks_child_0003`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks join tbl_ks_child",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_0000` JOIN `tbl_ks_child_0000`",
						"SELECT * FROM `tbl_ks_0001` JOIN `tbl_ks_child_0001`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_0002` JOIN `tbl_ks_child_0002`",
						"SELECT * FROM `tbl_ks_0003` JOIN `tbl_ks_child_0003`",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectKingshardNumRange(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_range where id < 0",
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id <= 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`<=0",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id < 50",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`<50",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id < 100",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`<100",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id <= 100",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`<=100",
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`<=100",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id < 300",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`<300",
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`<300",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`<300",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id <= 300",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`<=300",
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`<=300",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`<=300",
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`<=300",
					},
				},
			},
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where id < 400",
			hasErr: true, // shard key not in key range
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where id <= 400",
			hasErr: true, // shard key not in key range
		},

		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id > 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`>0",
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`>0",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`>0",
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`>0",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id >= 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`>=0",
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`>=0",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`>=0",
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`>=0",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id > 50",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE `id`>50",
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`>50",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`>50",
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`>50",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id > 100",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`>100",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`>100",
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`>100",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id >= 100",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`>=100",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`>=100",
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`>=100",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id > 300",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`>300",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id >= 300",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`>=300",
					},
				},
			},
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where id > 400",
			hasErr: true, // shard key not in key range
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where id >= 400",
			hasErr: true, // shard key not in key range
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id <= 300 AND id > 300",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0003` WHERE `id`<=300 AND `id`>300", // we can't handle this testcase
					},
				},
			},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_range where id < 300 AND id >= 300",
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where id > 100 AND id < 300",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0001` WHERE `id`>100 AND `id`<300",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE `id`>100 AND `id`<300",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectKingshardNumRangeReverse(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_range where 0>id",
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 0>=id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 0>=`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 50>id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 50>`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 100 > id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 100>`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 100 >= id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 100>=`id`",
						"SELECT * FROM `tbl_ks_range_0001` WHERE 100>=`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 300>id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 300>`id`",
						"SELECT * FROM `tbl_ks_range_0001` WHERE 300>`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 300>`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 300 >= id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 300>=`id`",
						"SELECT * FROM `tbl_ks_range_0001` WHERE 300>=`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 300>=`id`",
						"SELECT * FROM `tbl_ks_range_0003` WHERE 300>=`id`",
					},
				},
			},
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where 400 > id",
			hasErr: true, // shard key not in key range
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where 400 >= id",
			hasErr: true, // shard key not in key range
		},

		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 0 < id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 0<`id`",
						"SELECT * FROM `tbl_ks_range_0001` WHERE 0<`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 0<`id`",
						"SELECT * FROM `tbl_ks_range_0003` WHERE 0<`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 0 <= id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 0<=`id`",
						"SELECT * FROM `tbl_ks_range_0001` WHERE 0<=`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 0<=`id`",
						"SELECT * FROM `tbl_ks_range_0003` WHERE 0<=`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 50 < id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0000` WHERE 50<`id`",
						"SELECT * FROM `tbl_ks_range_0001` WHERE 50<`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 50<`id`",
						"SELECT * FROM `tbl_ks_range_0003` WHERE 50<`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 100 < id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0001` WHERE 100<`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 100<`id`",
						"SELECT * FROM `tbl_ks_range_0003` WHERE 100<`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 100<=id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0001` WHERE 100<=`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 100<=`id`",
						"SELECT * FROM `tbl_ks_range_0003` WHERE 100<=`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 300 < id",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0003` WHERE 300<`id`",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 300 <= id",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0003` WHERE 300<=`id`",
					},
				},
			},
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where 400 < id",
			hasErr: true, // shard key not in key range
		},
		{
			db:     "db_ks",
			sql:    "select * from tbl_ks_range where 400 <= id",
			hasErr: true, // shard key not in key range
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 300 >= id AND 300 < id",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0003` WHERE 300>=`id` AND 300<`id`", // we can't handle this testcase
					},
				},
			},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_range where 300 > id AND 300 <= id",
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_range where 100 < id AND 300 > id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0001` WHERE 100<`id` AND 300>`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_range_0002` WHERE 100<`id` AND 300>`id`",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectKingshardDateYear(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_year where create_time > 1564070400", // 2019/07/26 00:00:00
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2019` WHERE `create_time`>1564070400",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_year where create_time > 1406304000", // 2014/07/26 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2014` WHERE `create_time`>1406304000",
						"SELECT * FROM `tbl_ks_year_2015` WHERE `create_time`>1406304000",
						"SELECT * FROM `tbl_ks_year_2016` WHERE `create_time`>1406304000",
						"SELECT * FROM `tbl_ks_year_2017` WHERE `create_time`>1406304000",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2018` WHERE `create_time`>1406304000",
						"SELECT * FROM `tbl_ks_year_2019` WHERE `create_time`>1406304000",
					},
				},
			},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_year where create_time < 1388505600", // 2014/01/01 00:00:00
			sqls: map[string]map[string][]string{},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_year where create_time < 1385827200", // 2013/12/01 00:00:00
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_year where create_time <= 1388505600", // 2014/01/01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2014` WHERE `create_time`<=1388505600",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_year where create_time <= 1388505600", // 2014/01/01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2014` WHERE `create_time`<=1388505600",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_year where create_time > 1514735999", // 2017/12/31 23:59:59
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2017` WHERE `create_time`>1514735999",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2018` WHERE `create_time`>1514735999",
						"SELECT * FROM `tbl_ks_year_2019` WHERE `create_time`>1514735999",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_year where create_time >= 1514736000", // 2018/01/01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_year_2018` WHERE `create_time`>=1514736000",
						"SELECT * FROM `tbl_ks_year_2019` WHERE `create_time`>=1514736000",
					},
				},
			},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_year where create_time >= 1577808000", // 2020/01/01 00:00:00
			sqls: map[string]map[string][]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectKingshardDateMonth(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_month where create_time > 1398873600", // 2014/05/01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_month_201405` WHERE `create_time`>1398873600",
						"SELECT * FROM `tbl_ks_month_201406` WHERE `create_time`>1398873600",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_month_201408` WHERE `create_time`>1398873600",
						"SELECT * FROM `tbl_ks_month_201409` WHERE `create_time`>1398873600",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_month where create_time > 1398614400", // 2014/04/28 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_month_201405` WHERE `create_time`>1398614400",
						"SELECT * FROM `tbl_ks_month_201406` WHERE `create_time`>1398614400",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_month_201408` WHERE `create_time`>1398614400",
						"SELECT * FROM `tbl_ks_month_201409` WHERE `create_time`>1398614400",
					},
				},
			},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_month where create_time < 1398873600", // 2014/05/01 00:00:00
			sqls: map[string]map[string][]string{},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_month where create_time < 1398614400", // 2014/04/28 00:00:00
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_month where create_time < 1404144000", // 2014/07/01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_month_201405` WHERE `create_time`<1404144000",
						"SELECT * FROM `tbl_ks_month_201406` WHERE `create_time`<1404144000",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_month where create_time >= 1404144000", // 2014/07/01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_month_201408` WHERE `create_time`>=1404144000",
						"SELECT * FROM `tbl_ks_month_201409` WHERE `create_time`>=1404144000",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_month where create_time > 1412092799", // 2014/09/30 23:59:59
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_month_201409` WHERE `create_time`>1412092799",
					},
				},
			},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_month where create_time >= 1412092800", // 2014/10/01 00:00:00
			sqls: map[string]map[string][]string{},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_month where create_time >= 1412438400", // 2014/10/05 00:00:00
			sqls: map[string]map[string][]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectKingshardDateDay(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_day where create_time < 1409500800", // 2014/09/01 00:00:00
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time <= 1409500800", // 2014/09/01 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140901` WHERE `create_time`<=1409500800",
					},
				},
			},
		},
		{
			db:   "db_ks",
			sql:  "select * from tbl_ks_day where create_time >= 1410278400", // 2014/09/10 00:00:00
			sqls: map[string]map[string][]string{},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time = 1409846400", // 2014/09/05 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140905` WHERE `create_time`=1409846400",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time = 1409932800", // 2014/09/06 00:00:00
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time = 1410019200", // 2014/09/07 00:00:00
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140907` WHERE `create_time`=1410019200",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time > 1410019200", // 2014/09/07 00:00:00
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140907` WHERE `create_time`>1410019200",
						"SELECT * FROM `tbl_ks_day_20140908` WHERE `create_time`>1410019200",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time > 1410019200", // 2014/09/07 00:00:00
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140907` WHERE `create_time`>1410019200",
						"SELECT * FROM `tbl_ks_day_20140908` WHERE `create_time`>1410019200",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time < 1410278400", // 2014/09/10 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140901` WHERE `create_time`<1410278400",
						"SELECT * FROM `tbl_ks_day_20140902` WHERE `create_time`<1410278400",
						"SELECT * FROM `tbl_ks_day_20140903` WHERE `create_time`<1410278400",
						"SELECT * FROM `tbl_ks_day_20140904` WHERE `create_time`<1410278400",
						"SELECT * FROM `tbl_ks_day_20140905` WHERE `create_time`<1410278400",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140907` WHERE `create_time`<1410278400",
						"SELECT * FROM `tbl_ks_day_20140908` WHERE `create_time`<1410278400",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_day where create_time >= 1408464000", // 2014/08/20 00:00:00
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140901` WHERE `create_time`>=1408464000",
						"SELECT * FROM `tbl_ks_day_20140902` WHERE `create_time`>=1408464000",
						"SELECT * FROM `tbl_ks_day_20140903` WHERE `create_time`>=1408464000",
						"SELECT * FROM `tbl_ks_day_20140904` WHERE `create_time`>=1408464000",
						"SELECT * FROM `tbl_ks_day_20140905` WHERE `create_time`>=1408464000",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_day_20140907` WHERE `create_time`>=1408464000",
						"SELECT * FROM `tbl_ks_day_20140908` WHERE `create_time`>=1408464000",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectMultiTablesOnConditionKingshard(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select * from tbl_ks join tbl_ks_child on tbl_ks.id in (1,2,3) AND tbl_ks.id = tbl_ks_child.id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_0001` JOIN `tbl_ks_child_0001` ON `tbl_ks_0001`.`id` IN (1) AND `tbl_ks_0001`.`id`=`tbl_ks_child_0001`.`id`",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_0002` JOIN `tbl_ks_child_0002` ON `tbl_ks_0002`.`id` IN (2) AND `tbl_ks_0002`.`id`=`tbl_ks_child_0002`.`id`",
						"SELECT * FROM `tbl_ks_0003` JOIN `tbl_ks_child_0003` ON `tbl_ks_0003`.`id` IN (3) AND `tbl_ks_0003`.`id`=`tbl_ks_child_0003`.`id`",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectMultiTablesComparisonKingshard(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select * from tbl_ks where (tbl_ks.id = 3 OR 1 = 1) AND 1 = 0",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_0000` WHERE (`tbl_ks_0000`.`id`=3 OR 1=1) AND 1=0",
						"SELECT * FROM `tbl_ks_0001` WHERE (`tbl_ks_0001`.`id`=3 OR 1=1) AND 1=0",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_0002` WHERE (`tbl_ks_0002`.`id`=3 OR 1=1) AND 1=0",
						"SELECT * FROM `tbl_ks_0003` WHERE (`tbl_ks_0003`.`id`=3 OR 1=1) AND 1=0",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectShardTableWithGlobalTableKingshard(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select * from tbl_ks, tbl_ks_global_one where tbl_ks.id = 3 and tbl_ks_global_one.name='haha'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_ks": {
						"SELECT * FROM (`tbl_ks_0003`) JOIN `tbl_ks_global_one` WHERE `tbl_ks_0003`.`id`=3 AND `tbl_ks_global_one`.`name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks, tbl_ks_global_one where tbl_ks.unshard_col = 3 and tbl_ks_global_one.name='haha'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM (`tbl_ks_0000`) JOIN `tbl_ks_global_one` WHERE `tbl_ks_0000`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha'",
						"SELECT * FROM (`tbl_ks_0001`) JOIN `tbl_ks_global_one` WHERE `tbl_ks_0001`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha'",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM (`tbl_ks_0002`) JOIN `tbl_ks_global_one` WHERE `tbl_ks_0002`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha'",
						"SELECT * FROM (`tbl_ks_0003`) JOIN `tbl_ks_global_one` WHERE `tbl_ks_0003`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from tbl_ks, tbl_ks_global_one, tbl_ks_global_two where tbl_ks.unshard_col = 3 and tbl_ks_global_one.name='haha' and tbl_ks_global_two.gender='female'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM ((`tbl_ks_0000`) JOIN `tbl_ks_global_one`) JOIN `tbl_ks_global_two` WHERE `tbl_ks_0000`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha' AND `tbl_ks_global_two`.`gender`='female'",
						"SELECT * FROM ((`tbl_ks_0001`) JOIN `tbl_ks_global_one`) JOIN `tbl_ks_global_two` WHERE `tbl_ks_0001`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha' AND `tbl_ks_global_two`.`gender`='female'",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT * FROM ((`tbl_ks_0002`) JOIN `tbl_ks_global_one`) JOIN `tbl_ks_global_two` WHERE `tbl_ks_0002`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha' AND `tbl_ks_global_two`.`gender`='female'",
						"SELECT * FROM ((`tbl_ks_0003`) JOIN `tbl_ks_global_one`) JOIN `tbl_ks_global_two` WHERE `tbl_ks_0003`.`unshard_col`=3 AND `tbl_ks_global_one`.`name`='haha' AND `tbl_ks_global_two`.`gender`='female'",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectGlobalTableKingshard(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select * from tbl_ks_global_one where name='haha'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM `tbl_ks_global_one` WHERE `name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select * from db_ks.tbl_ks_global_one, tbl_ks_global_two where tbl_ks_global_one.name='haha' and tbl_ks_global_two.name='hehe'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT * FROM (`db_ks`.`tbl_ks_global_one`) JOIN `tbl_ks_global_two` WHERE `tbl_ks_global_one`.`name`='haha' AND `tbl_ks_global_two`.`name`='hehe'",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectShardTableWithGlobalTableMycat(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_global_one where tbl_mycat.id = 3 and tbl_mycat_global_one.name='haha'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {
						"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_global_one` WHERE `tbl_mycat`.`id`=3 AND `tbl_mycat_global_one`.`name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from db_mycat.tbl_mycat, db_mycat.tbl_mycat_global_one where db_mycat.tbl_mycat.id = 3 and db_mycat.tbl_mycat_global_one.name='haha'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {
						"SELECT * FROM (`db_mycat_3`.`tbl_mycat`) JOIN `db_mycat_3`.`tbl_mycat_global_one` WHERE `db_mycat_3`.`tbl_mycat`.`id`=3 AND `db_mycat_3`.`tbl_mycat_global_one`.`name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_global_one where tbl_mycat.unshard_col = 3 and tbl_mycat_global_one.name='haha'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_global_one` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha'",
					},
					"db_mycat_1": {
						"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_global_one` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha'",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_global_one` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha'",
					},
					"db_mycat_3": {
						"SELECT * FROM (`tbl_mycat`) JOIN `tbl_mycat_global_one` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat, tbl_mycat_global_one, tbl_mycat_global_two where tbl_mycat.unshard_col = 3 and tbl_mycat_global_one.name='haha' and tbl_mycat_global_two.gender='female'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM ((`tbl_mycat`) JOIN `tbl_mycat_global_one`) JOIN `tbl_mycat_global_two` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha' AND `tbl_mycat_global_two`.`gender`='female'",
					},
					"db_mycat_1": {
						"SELECT * FROM ((`tbl_mycat`) JOIN `tbl_mycat_global_one`) JOIN `tbl_mycat_global_two` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha' AND `tbl_mycat_global_two`.`gender`='female'",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM ((`tbl_mycat`) JOIN `tbl_mycat_global_one`) JOIN `tbl_mycat_global_two` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha' AND `tbl_mycat_global_two`.`gender`='female'",
					},
					"db_mycat_3": {
						"SELECT * FROM ((`tbl_mycat`) JOIN `tbl_mycat_global_one`) JOIN `tbl_mycat_global_two` WHERE `tbl_mycat`.`unshard_col`=3 AND `tbl_mycat_global_one`.`name`='haha' AND `tbl_mycat_global_two`.`gender`='female'",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectGlobalTableMycat(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_global_one where name='haha'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_global_one` WHERE `name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from db_mycat.tbl_mycat_global_one where name='haha'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `db_mycat_0`.`tbl_mycat_global_one` WHERE `name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from db_mycat.tbl_mycat_global_one where db_mycat.tbl_mycat_global_one.name='haha'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `db_mycat_0`.`tbl_mycat_global_one` WHERE `db_mycat_0`.`tbl_mycat_global_one`.`name`='haha'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from db_mycat.tbl_mycat_global_one, tbl_mycat_global_two where tbl_mycat_global_one.name='haha' and tbl_mycat_global_two.name='hehe'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM (`db_mycat_0`.`tbl_mycat_global_one`) JOIN `tbl_mycat_global_two` WHERE `tbl_mycat_global_one`.`name`='haha' AND `tbl_mycat_global_two`.`name`='hehe'",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from db_mycat.tbl_mycat_global_one, db_mycat.tbl_mycat_global_two where db_mycat.tbl_mycat_global_one.name='haha' and db_mycat.tbl_mycat_global_two.name='hehe'",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM (`db_mycat_0`.`tbl_mycat_global_one`) JOIN `db_mycat_0`.`tbl_mycat_global_two` WHERE `db_mycat_0`.`tbl_mycat_global_one`.`name`='haha' AND `db_mycat_0`.`tbl_mycat_global_two`.`name`='hehe'",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectMycatGroupByDatabase(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select database(), count(id) from tbl_mycat group by database()",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE()",
					},
					"db_mycat_1": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE()",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE()",
					},
					"db_mycat_3": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE()",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select database(), count(id) from tbl_mycat where database()='db_mycat_1' group by database()",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE() FROM `tbl_mycat` WHERE DATABASE()='db_mycat_1' GROUP BY DATABASE()",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectMycatOrderByDatabase(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select database(), count(id) from tbl_mycat group by database() order by database()",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE() ORDER BY DATABASE()",
					},
					"db_mycat_1": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE() ORDER BY DATABASE()",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE() ORDER BY DATABASE()",
					},
					"db_mycat_3": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` GROUP BY DATABASE() ORDER BY DATABASE()",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select database(), count(id) from tbl_mycat where database() in ('db_mycat_1','db_mycat_2') group by database() order by database()",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_1','db_mycat_2') GROUP BY DATABASE() ORDER BY DATABASE()",
					},
					"db_mycat_1": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_1','db_mycat_2') GROUP BY DATABASE() ORDER BY DATABASE()",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_1','db_mycat_2') GROUP BY DATABASE() ORDER BY DATABASE()",
					},
					"db_mycat_3": {
						"SELECT DATABASE(),COUNT(`id`),DATABASE(),DATABASE() FROM `tbl_mycat` WHERE DATABASE() IN ('db_mycat_1','db_mycat_2') GROUP BY DATABASE() ORDER BY DATABASE()",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectForceIndexDatabase(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat force index(id, name) where id > 100 and name = `zhangsan`",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat` FORCE INDEX (`id`, `name`) WHERE `id`>100 AND `name`=`zhangsan`",
					},
					"db_mycat_1": {
						"SELECT * FROM `tbl_mycat` FORCE INDEX (`id`, `name`) WHERE `id`>100 AND `name`=`zhangsan`",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM `tbl_mycat` FORCE INDEX (`id`, `name`) WHERE `id`>100 AND `name`=`zhangsan`",
					},
					"db_mycat_3": {
						"SELECT * FROM `tbl_mycat` FORCE INDEX (`id`, `name`) WHERE `id`>100 AND `name`=`zhangsan`",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectOrderByAliasColumn(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select count(id) as a, uid from tbl_mycat where uid = 2 order by a",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `a`",
					},
					"db_mycat_1": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `a`",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `a`",
					},
					"db_mycat_3": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `a`",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id as a, uid from tbl_mycat where uid = 2 order by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `id`",
					},
					"db_mycat_1": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `id`",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `id`",
					},
					"db_mycat_3": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 ORDER BY `id`",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func TestSelectGroupByAliasColumn(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select count(id) as a, uid from tbl_mycat where uid = 2 group by a",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `a`",
					},
					"db_mycat_1": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `a`",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `a`",
					},
					"db_mycat_3": {
						"SELECT COUNT(`id`) AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `a`",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select id as a, uid from tbl_mycat where uid = 2 group by id",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `id`",
					},
					"db_mycat_1": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `id`",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `id`",
					},
					"db_mycat_3": {
						"SELECT `id` AS `a`,`uid` FROM `tbl_mycat` WHERE `uid`=2 GROUP BY `id`",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getTestFunc(ns, test))
	}
}

func prepareShardKingshardRouter() (*router.Router, error) {
	nsStr := `
{
    "name": "gaea_namespace_1",
    "online": true,
    "read_only": true,
    "allowed_dbs": {
        "test": true
    },
    "default_phy_dbs": {
        "test": "db_mycat_0"
    },
    "slices": [
        {
            "name": "slice-0",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:3306",
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        },
        {
            "name": "slice-1",
            "user_name": "root",
            "password": "root",
            "master": "127.0.0.1:3307",
            "capacity": 64,
            "max_capacity": 128,
            "idle_timeout": 3600
        }
    ],
    "shard_rules": [
        {
            "db": "test",
            "table": "tbl_ks",
            "type": "hash",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        },
        {
            "db": "test",
            "table": "test_hash_1",
            "type": "hash",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ]
        },
		{
			"db": "test",
            "table": "tbl_ks_child",
            "type": "linked",
			"parent_table": "tbl_ks",
			"key": "id"
		}
    ],
    "users": [
        {
            "user_name": "test_shard_hash",
            "password": "test_shard_hash",
            "namespace": "gaea_namespace_1",
            "rw_flag": 2,
            "rw_split": 1
        }
    ],
    "default_slice": "slice-0"
}
`

	nsModel, err := createNamespace(nsStr)
	if err != nil {
		return nil, err
	}

	return createRouter(nsModel)
}
