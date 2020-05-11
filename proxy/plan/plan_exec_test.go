package plan

import (
	"testing"

	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/util"
)

func TestExecuteIn(t *testing.T) {
	planInfo, _ := preparePlanInfo()
	sql := "SELECT * FROM tbl_mycat_murmur WHERE tbl_mycat_murmur.id=5 AND tbl_mycat_murmur.id=4"
	stmt, _ := parser.ParseSQL(sql)
	plan, err := BuildPlan(stmt, nil, "db_mycat", sql, planInfo.rt, planInfo.seqs)
	if err != nil {
		t.Fatalf("build plan error: %v", err)
	}
	ret, err := plan.ExecuteIn(util.NewRequestContext(), nil)
	if err != nil {
		t.Fatalf("execute error: %v", err)
	}
	t.Logf("result: %v", ret)
}
