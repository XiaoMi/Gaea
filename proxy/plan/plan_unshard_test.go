package plan

import (
	"strings"
	"testing"

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
	tests := []UnshardSQLTestcase{
		{
			db:        "db_mycat",
			sql:       `select * from tbl_unshard_a as a join db_mycat.tbl_unshard_b as b on a.id = b.id`,
			expectSQL: "SELECT * FROM `tbl_unshard_a` AS `a` JOIN `db_mycat_0`.`tbl_unshard_b` AS `b` ON `a`.`id`=`b`.`id`",
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getUnshardTestFunc(ns, test))
	}
}


func TestUnshardPlanWithoutDB(t *testing.T) {
	ns, err := preparePlanInfo()
	if err != nil {
		t.Fatalf("prepare namespace error: %v", err)
	}
	tests := []UnshardSQLTestcase{
		{
			db:        "db_mycat",
			sql:       `select * from tbl_unshard_a`,
			expectSQL: "SELECT * FROM `tbl_unshard_a`",
		},
	}
	for _, test := range tests {
		t.Run(test.sql, getUnshardTestFunc(ns, test))
	}
}


// 获取使用TiDB parser测试SQL改写结果的测试函数
func getUnshardTestFunc(info *PlanInfo, test UnshardSQLTestcase) func(t *testing.T) {
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

		var actualSQLs string
		switch plan := p.(type) {
		case *UnshardPlan:
			actualSQLs = plan.sql
		default:
			t.Fatal("invalid plan type, only support unShard plan")
		}

		if actualSQLs == "" {
			t.Fatalf("get sqls error: %v", err)
		}

		if test.expectSQL != actualSQLs {
			t.Errorf("not equal, expect: %v, actual: %v", test.expectSQL, actualSQLs)
		}
	}
}
