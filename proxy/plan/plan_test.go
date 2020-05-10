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
	"encoding/json"
	"sync/atomic"
	"testing"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/proxy/sequence"
)

type SQLTestcase struct {
	db     string
	sql    string
	sqls   map[string]map[string][]string
	hasErr bool
}

type OrderSequence struct {
	v      int64
	db     string
	table  string
	pkName string
}

type PlanInfo struct {
	phyDBs map[string]string
	rt     *router.Router
	seqs   *sequence.SequenceManager
}

func NewOrderSequence(db, table, pkName string) *OrderSequence {
	return &OrderSequence{
		db:     db,
		table:  table,
		pkName: pkName,
	}
}

func (s *OrderSequence) GetPKName() string {
	return s.pkName
}

func (s *OrderSequence) NextSeq() (int64, error) {
	newInt := atomic.AddInt64(&s.v, 1)
	return newInt, nil
}

// 获取使用TiDB parser测试SQL改写结果的测试函数
func getTestFunc(info *PlanInfo, test SQLTestcase) func(t *testing.T) {
	return func(t *testing.T) {
		stmt, err := parser.ParseSQL(test.sql)
		if err != nil {
			if test.hasErr {
				t.Logf("parse sql error: %v", err)
				return
			}
			t.Fatalf("parse sql error: %v", err)
		}

		p, err := BuildPlan(stmt, info.phyDBs, test.db, test.sql, info.rt, info.seqs)
		if err != nil {
			if test.hasErr {
				t.Logf("BuildPlan got expect error, sql: %s, err: %v", test.sql, err)
				return
			}
			t.Fatalf("BuildPlan error, sql: %s, err: %v", test.sql, err)
		}

		var actualSQLs map[string]map[string][]string
		switch plan := p.(type) {
		case *SelectPlan:
			actualSQLs = plan.GetSQLs()
		case *InsertPlan:
			actualSQLs = plan.sqls
		case *UpdatePlan:
			actualSQLs = plan.sqls
		case *DeletePlan:
			actualSQLs = plan.sqls
		case *ExplainPlan:
			actualSQLs = plan.sqls
		case *UnshardPlan:
			actualSQLs = make(map[string]map[string][]string)
			dbSQLs := make(map[string][]string)
			if db, ok := info.phyDBs[plan.db]; ok {
				plan.db = db
			}
			dbSQLs[plan.db] = []string{plan.sql}
			actualSQLs[backend.DefaultSlice] = dbSQLs
		}

		if actualSQLs == nil {
			t.Fatalf("get sqls error: %v", err)
		}

		if !checkSQLs(test.sqls, actualSQLs) {
			t.Errorf("not equal, expect: %v, actual: %v", test.sqls, actualSQLs)
		}
	}
}

func createNamespace(nsStr string) (*models.Namespace, error) {
	ns := &models.Namespace{}
	err := json.Unmarshal([]byte(nsStr), ns)
	if err != nil {
		return nil, err
	}
	if err := ns.Verify(); err != nil {
		return nil, err
	}
	return ns, nil
}

func createRouter(nsCfg *models.Namespace) (*router.Router, error) {
	return router.NewRouter(nsCfg)
}

func createSequenceManager(nsCfg *models.Namespace) (*sequence.SequenceManager, error) {
	sequences := sequence.NewSequenceManager()
	for _, v := range nsCfg.GlobalSequences {
		seq := NewOrderSequence(v.DB, v.Table, v.PKName)
		sequences.SetSequence(v.DB, v.Table, seq)
	}
	return sequences, nil
}

func checkSQLs(expect, actual map[string]map[string][]string) bool {
	aContainsB := func(a, b map[string]map[string][]string) bool {
		for sliceName, dbSQLs := range a {
			if len(b[sliceName]) == 0 {
				return false
			}
			for dbName, sqls := range dbSQLs {
				if len(b[sliceName][dbName]) == 0 {
					return false
				}
				for i, sql := range sqls {
					if b[sliceName][dbName][i] != sql {
						return false
					}
				}
			}
		}
		return true
	}

	return aContainsB(expect, actual) && aContainsB(actual, expect)
}

func preparePlanInfo() (*PlanInfo, error) {
	nsStr := `
{
    "name": "gaea_namespace_1",
    "online": true,
    "read_only": true,
    "allowed_dbs": {
        "db_ks": true,
        "db_mycat": true
    },
    "default_phy_dbs": {
        "db_ks": "db_ks",
        "db_mycat": "db_mycat_0"
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
            "db": "db_ks",
            "table": "tbl_ks",
            "type": "mod",
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
            "db": "db_ks",
            "table": "tbl_ks_child",
            "type": "linked",
            "key": "id",
            "parent_table": "tbl_ks"
        },
        {
            "db": "db_ks",
            "table": "tbl_ks_user_child",
            "type": "linked",
            "key": "user_id",
            "parent_table": "tbl_ks"
        },
		{
            "db": "db_ks",
            "table": "tbl_ks_global_one",
            "type": "global",
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
            "db": "db_ks",
            "table": "tbl_ks_global_two",
            "type": "global",
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
			"db": "db_ks",
            "table": "tbl_ks_range",
            "type": "range",
			"key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"table_row_limit": 100
		},
		{
			"db": "db_ks",
            "table": "tbl_ks_year",
            "type": "date_year",
			"key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"date_range": [
				"2014-2017",
				"2018-2019"
			]
		},
		{
			"db": "db_ks",
            "table": "tbl_ks_month",
            "type": "date_month",
			"key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"date_range": [
				"201405-201406",
				"201408-201409"
			]
		},
		{
			"db": "db_ks",
            "table": "tbl_ks_day",
            "type": "date_day",
			"key": "create_time",
            "slices": [
                "slice-0",
                "slice-1"
            ],
			"date_range": [
				"20140901-20140905",
				"20140907-20140908"
			]
        },
        {
            "db": "db_ks",
            "table": "TBL_KS_UPPERCASE",
            "type": "mod",
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
            "db": "db_ks",
            "table": "TBL_KS_UPPERCASE_CHILD",
            "type": "linked",
            "key": "ID",
            "parent_table": "TBL_KS_UPPERCASE"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat",
            "type": "mycat_mod",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_child",
            "type": "linked",
            "parent_table": "tbl_mycat",
            "key": "id"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_user_child",
            "type": "linked",
            "parent_table": "tbl_mycat",
            "key": "user_id"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_murmur",
            "type": "mycat_murmur",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_0","db_mycat_1","db_mycat_2","db_mycat_3"
            ],
			"seed": "0",
			"virtual_bucket_times": "160"
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_long",
            "type": "mycat_long",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
			"partition_count": "4",
			"partition_length": "256"
        },
		{
            "db": "db_mycat",
            "table": "tbl_mycat_global_one",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
		{
            "db": "db_mycat",
            "table": "tbl_mycat_global_two",
            "type": "global",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ]
        },
        {
            "db": "db_mycat",
            "table": "tbl_mycat_string",
            "type": "mycat_string",
            "key": "id",
            "locations": [
                2,
                2
            ],
            "slices": [
                "slice-0",
                "slice-1"
            ],
            "databases": [
                "db_mycat_[0-3]"
            ],
			"partition_count": "4",
			"partition_length": "256",
			"hash_slice": "20"
        }
    ],
	"global_sequences": [
		{
			"db": "db_mycat",
			"table": "tbl_mycat",
			"type": "test",
			"pk_name": "id"
		},
		{
			"db": "db_ks",
			"table": "tbl_ks",
			"type": "test",
			"pk_name": "user_id"
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
}`
	nsModel, err := createNamespace(nsStr)
	if err != nil {
		return nil, err
	}

	rt, err := createRouter(nsModel)
	if err != nil {
		return nil, err
	}

	seqs, err := createSequenceManager(nsModel)
	if err != nil {
		return nil, err
	}

	planInfo := &PlanInfo{
		phyDBs: nsModel.DefaultPhyDBS,
		rt:     rt,
		seqs:   seqs,
	}
	return planInfo, nil
}
