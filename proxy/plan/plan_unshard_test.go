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
