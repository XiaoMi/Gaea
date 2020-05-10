package plan

import (
	"strings"
	"testing"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/format"
)

type UnshardSQLTestcase struct {
	db        string
	sql       string
	expectSQL string
	hasErr    bool
}

func TestRewriteUnshardTableName(t *testing.T) {
	sql := `select * from tbl_a as a join db_unshard.tbl_b as b on a.id = b.id`
	expectSQL := "SELECT * FROM `tbl_a` AS `a` JOIN `db_unshard_0`.`tbl_b` AS `b` ON `a`.`id`=`b`.`id`"
	n, err := parser.ParseSQL(sql)
	if err != nil {
		t.Fatalf("parse sql error: %v", err)
	}
	v := &UnshardTableRewriteVisitor{PhyDBs: map[string]string{
		"db_unshard": "db_unshard_0",
	}}

	n.Accept(v)
	s := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, s)
	_ = n.Restore(ctx)
	if s.String() != expectSQL {
		t.Fatal("error rewrite SQL", s.String())
	}
}

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
