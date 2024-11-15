package plan

import (
	"testing"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/parser"
)

func TestSimpleUnionShardMycatMod(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 0 UNION select * from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select tbl_mycat.id from tbl_mycat where id = 0 UNION select tbl_mycat.id from tbl_mycat where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `tbl_mycat`.`id` FROM `tbl_mycat` WHERE `id`=0"},
					"db_mycat_1": {"SELECT `tbl_mycat`.`id` FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where tbl_mycat.id = 0 UNION select * from tbl_mycat where tbl_mycat.id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `tbl_mycat`.`id`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `tbl_mycat`.`id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id in (0,1) UNION ALL select * from tbl_mycat where id = 2",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (0)"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (1)"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`=2"},
				},
			},
		},

		{
			db:  "db_mycat",
			sql: "select id from tbl_mycat where id = 0 UNION select id from tbl_mycat where id = 1 ORDER BY id DESC",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT `id` FROM `tbl_mycat` WHERE `id`=0"},
					"db_mycat_1": {"SELECT `id` FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 0 UNION select * from tbl_mycat where id = 1 LIMIT 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`=1"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id = 0 UNION ALL select * from tbl_mycat where id = 1 UNION select * from tbl_mycat where id = 2 UNION select * from tbl_mycat where id = 3",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat` WHERE `id`=0"},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id`=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat` WHERE `id`=2"},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id`=3"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id in (0,1,2) UNION select * from tbl_mycat where id in (3,4,6)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat` WHERE `id` IN (0)",
						"SELECT * FROM `tbl_mycat` WHERE `id` IN (4)",
					},
					"db_mycat_1": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (1)"},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM `tbl_mycat` WHERE `id` IN (2)",
						"SELECT * FROM `tbl_mycat` WHERE `id` IN (6)",
					},
					"db_mycat_3": {"SELECT * FROM `tbl_mycat` WHERE `id` IN (3)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat where id in (0,4) UNION select * from tbl_mycat where id in (1) UNION select * from tbl_mycat where id in (2,6)  UNION select * from tbl_mycat where id in (3) ",
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
	}

	for _, test := range tests {
		t.Run(test.sql, getUnionTestFunc(ns, test))
	}
}

func TestSimpleUnionShardMycatMurmur(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 0 UNION select * from tbl_mycat_murmur where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=1"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=0"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 2 UNION select * from tbl_mycat_murmur where id = 3",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=2",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=3",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = 4 UNION select * from tbl_mycat_murmur where id = 4",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=4",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id`=4",
					},
				},
			},
		},

		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id in (0,1,2,3,4,6) UNION select * from tbl_mycat_murmur where id in (0,1,2,3,4,6)",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (6)",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (6)",
					},
					"db_mycat_1": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (1,2,3)",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (1,2,3)",
					},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (0,4)", "SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN (0,4)"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where k = 0 UNION select * from tbl_mycat_murmur where k = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=1",
					},
					"db_mycat_1": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=1",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=1",
					},
					"db_mycat_3": {
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_murmur` WHERE `k`=1",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getUnionTestFunc(ns, test))
	}
}

func TestSimpleUnionShardMycatMurmur_ShardKeyTypeString(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id in ('0') UNION select * from tbl_mycat_murmur where id in ('1')",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN ('0')"},
				},
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id` IN ('1')"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_murmur where id = '0' UNION select * from tbl_mycat_murmur where id = '1'",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`='0'"},
				},
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_murmur` WHERE `id`='1'"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getUnionTestFunc(ns, test))
	}
}

func TestSimpleUnionShardMycatString(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_string where id = 0 UNION select * from tbl_mycat_string where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_string` WHERE `id`=0",
						"SELECT * FROM `tbl_mycat_string` WHERE `id`=1",
					},
				},
			},
		},

		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_string where k = 0 UNION select * from tbl_mycat_string where k = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=1",
					},
					"db_mycat_1": {
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=1",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=1",
					},
					"db_mycat_3": {
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_string` WHERE `k`=1",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getUnionTestFunc(ns, test))
	}
}

func TestSimpleUnionShardMycatLong(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 0 UNION select * from tbl_mycat_long where id = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=0",
						"SELECT * FROM `tbl_mycat_long` WHERE `id`=1",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 2 UNION select * from tbl_mycat_long where id = 3",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_long` WHERE `id`=2",
						"SELECT * FROM `tbl_mycat_long` WHERE `id`=3",
					},
				},
			},
		},

		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 256 UNION select * from tbl_mycat_long where id = 512",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_1": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=256"},
				},
				"slice-1": {
					"db_mycat_2": {"SELECT * FROM `tbl_mycat_long` WHERE `id`=512"},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 768 UNION select * from tbl_mycat_long where id = 769",
			sqls: map[string]map[string][]string{
				"slice-1": {
					"db_mycat_3": {
						"SELECT * FROM `tbl_mycat_long` WHERE `id`=768",
						"SELECT * FROM `tbl_mycat_long` WHERE `id`=769",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id = 1024 UNION select * from tbl_mycat_long where id = 1025",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_long` WHERE `id`=1024",
						"SELECT * FROM `tbl_mycat_long` WHERE `id`=1025",
					},
				},
			},
		},
		{
			db:  "db_mycat",
			sql: "select * from tbl_mycat_long where id in (0,1,1024,1025) UNION select * from tbl_mycat_long where id in (256,257) UNION select * from tbl_mycat_long where id in (512,513) UNION  select * from tbl_mycat_long where id in (768,769)",
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
			sql: "select * from tbl_mycat_long where k = 0 UNION select * from tbl_mycat_long where k = 1",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_mycat_0": {
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=1",
					},
					"db_mycat_1": {
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=1",
					},
				},
				"slice-1": {
					"db_mycat_2": {
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=1",
					},
					"db_mycat_3": {
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=0",
						"SELECT * FROM `tbl_mycat_long` WHERE `k`=1",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getUnionTestFunc(ns, test))
	}
}

func TestSimpleUnionShardKingShardMod(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}

	tests := []SQLTestcase{
		{
			db:  "db_ks",
			sql: "select  name from tbl_ks as a where a.id =4 UNION select  name from tbl_ks as a where a.id =8 ",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `name` FROM `tbl_ks_0000` AS `a` WHERE `a`.`id`=4",
						"SELECT `name` FROM `tbl_ks_0000` AS `a` WHERE `a`.`id`=8",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select  name from tbl_ks as a where a.id =1 UNION select  name from tbl_ks as a where a.id =5 UNION select  name from tbl_ks as a where a.id =9 ",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `name` FROM `tbl_ks_0001` AS `a` WHERE `a`.`id`=1",
						"SELECT `name` FROM `tbl_ks_0001` AS `a` WHERE `a`.`id`=5",
						"SELECT `name` FROM `tbl_ks_0001` AS `a` WHERE `a`.`id`=9",
					},
				},
			},
		},
		{
			db:  "db_ks",
			sql: "select  name from tbl_ks as a where a.id =1 UNION select  name from tbl_ks as a where a.id =2 UNION select  name from tbl_ks as a where a.id =3 UNION  select  name from tbl_ks as a where a.id =4 ",
			sqls: map[string]map[string][]string{
				"slice-0": {
					"db_ks": {
						"SELECT `name` FROM `tbl_ks_0001` AS `a` WHERE `a`.`id`=1",
						"SELECT `name` FROM `tbl_ks_0000` AS `a` WHERE `a`.`id`=4",
					},
				},
				"slice-1": {
					"db_ks": {
						"SELECT `name` FROM `tbl_ks_0002` AS `a` WHERE `a`.`id`=2",
						"SELECT `name` FROM `tbl_ks_0003` AS `a` WHERE `a`.`id`=3",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.sql, getUnionTestFunc(ns, test))
	}
}

func getUnionTestFunc(info *PlanInfo, test SQLTestcase) func(t *testing.T) {
	return func(t *testing.T) {
		// Parse the SQL
		stmt, err := parser.ParseSQL(test.sql)
		if err != nil {
			if test.hasErr {
				t.Logf("parse sql error: %v", err)
				return
			}
			t.Fatalf("parse sql error: %v", err)
		}

		// Build the plan from the parsed statement
		p, err := BuildPlan(stmt, info.phyDBs, test.db, test.sql, info.rt, info.seqs, nil)
		if err != nil {
			if test.hasErr {
				t.Logf("BuildPlan got expected error, sql: %s, err: %v", test.sql, err)
				return
			}
			t.Fatalf("BuildPlan error, sql: %s, err: %v", test.sql, err)
		}

		// Define a helper function to check subPlan SQLs recursively
		var collectSQLs func(p Plan) (map[string]map[string][]string, error)
		collectSQLs = func(plan Plan) (map[string]map[string][]string, error) {
			var actualSQLs map[string]map[string][]string

			switch p := plan.(type) {
			case *SelectPlan:
				actualSQLs = p.GetSQLs()
			case *InsertPlan:
				actualSQLs = p.sqls
			case *UpdatePlan:
				actualSQLs = p.sqls
			case *DeletePlan:
				actualSQLs = p.sqls
			case *ExplainPlan:
				actualSQLs = p.sqls
			case *SelectLastInsertIDPlan:
				actualSQLs = make(map[string]map[string][]string)
			case *UnshardPlan:
				actualSQLs = make(map[string]map[string][]string)
				dbSQLs := make(map[string][]string)
				if db, ok := info.phyDBs[p.db]; ok {
					p.db = db
				}
				dbSQLs[p.db] = []string{p.sql}
				actualSQLs[backend.DefaultSlice] = dbSQLs
			case *UnionPlan:
				// If it's a UnionPlan, collect the SQLs for each subplan
				actualSQLs = make(map[string]map[string][]string)
				for _, subPlan := range p.subPlans {
					subPlanSQLs, err := collectSQLs(subPlan)
					if err != nil {
						return nil, err
					}
					// Merge subPlanSQLs into the actualSQLs map
					for key, val := range subPlanSQLs {
						if actualSQLs[key] == nil {
							actualSQLs[key] = make(map[string][]string)
						}
						for dbKey, dbVal := range val {
							actualSQLs[key][dbKey] = append(actualSQLs[key][dbKey], dbVal...)
						}
					}
				}
			}

			return actualSQLs, nil
		}

		// Collect SQLs for the plan
		actualSQLs, err := collectSQLs(p)
		if err != nil {
			t.Fatalf("failed to collect SQLs: %v", err)
		}

		// Assert that the collected SQLs match the expected ones
		if len(test.randomSqls) > 0 {
			if !checkRandomSQLs(test.randomSqls, actualSQLs) {
				t.Errorf("not equal, expect: %v, actual: %v", test.randomSqls, actualSQLs)
				return
			}
			return
		}

		if !checkSQLs(test.sqls, actualSQLs) {
			t.Errorf("not equal, expect: %v, actual: %v", test.sqls, actualSQLs)
		}
	}
}
