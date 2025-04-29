// Copyright 2016 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package ast_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	. "github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/auth"
	"github.com/XiaoMi/Gaea/parser/format"
)

type visitor struct{}

func (v visitor) Enter(in Node) (Node, bool) {
	return in, false
}

func (v visitor) Leave(in Node) (Node, bool) {
	return in, true
}

type visitor1 struct {
	visitor
}

func (visitor1) Enter(in Node) (Node, bool) {
	return in, true
}

func TestMiscVisitorCover(t *testing.T) {
	valueExpr := NewValueExpr(42)
	stmts := []Node{
		&AdminStmt{},
		&AlterUserStmt{},
		&BeginStmt{},
		&BinlogStmt{},
		&CommitStmt{},
		&CreateUserStmt{},
		&DeallocateStmt{},
		&DoStmt{},
		&ExecuteStmt{UsingVars: []ExprNode{valueExpr}},
		&ExplainStmt{Stmt: &ShowStmt{}},
		&GrantStmt{},
		&PrepareStmt{SQLVar: &VariableExpr{Value: valueExpr}},
		&RollbackStmt{},
		&SetPwdStmt{},
		&SetStmt{Variables: []*VariableAssignment{
			{
				Value: valueExpr,
			},
		}},
		&UseStmt{},
		&AnalyzeTableStmt{
			TableNames: []*TableName{
				{},
			},
		},
		&FlushStmt{},
		&PrivElem{},
		&VariableAssignment{Value: valueExpr},
		&KillStmt{},
		&DropStatsStmt{Table: &TableName{}},
	}

	for _, v := range stmts {
		v.Accept(visitor{})
		v.Accept(visitor1{})
	}
}

func TestDDLVisitorCoverMisc(t *testing.T) {
	sql := `
create table t (c1 smallint unsigned, c2 int unsigned);
alter table t add column a smallint unsigned after b;
create index t_i on t (id);
create database test character set utf8;
drop database test;
drop index t_i on t;
drop table t;
truncate t;
create table t (
jobAbbr char(4) not null,
constraint foreign key (jobabbr) references ffxi_jobtype (jobabbr) on delete cascade on update cascade
);
`
	parse := parser.New()
	stmts, _, err := parse.Parse(sql, "", "")
	require.NoError(t, err)
	for _, stmt := range stmts {
		stmt.Accept(visitor{})
		stmt.Accept(visitor1{})
	}
}

func TestDMLVistorCover(t *testing.T) {
	sql := `delete from somelog where user = 'jcole' order by timestamp_column limit 1;
delete t1, t2 from t1 inner join t2 inner join t3 where t1.id=t2.id and t2.id=t3.id;
select * from t where exists(select * from t k where t.c = k.c having sum(c) = 1);
insert into t_copy select * from t where t.x > 5;
(select /*+ TIDB_INLJ(t1) */ a from t1 where a=10 and b=1) union (select /*+ TIDB_SMJ(t2) */ a from t2 where a=11 and b=2) order by a limit 10;
update t1 set col1 = col1 + 1, col2 = col1;
show create table t;
load data infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b';`

	p := parser.New()
	stmts, _, err := p.Parse(sql, "", "")
	require.NoError(t, err)
	for _, stmt := range stmts {
		stmt.Accept(visitor{})
		stmt.Accept(visitor1{})
	}
}

func TestSensitiveStatement(t *testing.T) {
	positive := []StmtNode{
		&SetPwdStmt{},
		&CreateUserStmt{},
		&AlterUserStmt{},
		&GrantStmt{},
	}
	for i, stmt := range positive {
		_, ok := stmt.(SensitiveStmtNode)
		require.Truef(t, ok, "%d, %#v fail", i, stmt)
	}

	negative := []StmtNode{
		&DropUserStmt{},
		&RevokeStmt{},
		&AlterTableStmt{},
		&CreateDatabaseStmt{},
		&CreateIndexStmt{},
		&CreateTableStmt{},
		&DropDatabaseStmt{},
		&DropIndexStmt{},
		&DropTableStmt{},
		&RenameTableStmt{},
		&TruncateTableStmt{},
	}
	for _, stmt := range negative {
		_, ok := stmt.(SensitiveStmtNode)
		require.False(t, ok)
	}
}

func TestUserSpec(t *testing.T) {
	hashString := "*3D56A309CD04FA2EEF181462E59011F075C89548"
	u := UserSpec{
		User: &auth.UserIdentity{
			Username: "test",
		},
		AuthOpt: &AuthOption{
			ByAuthString: false,
			AuthString:   "xxx",
			HashString:   hashString,
		},
	}
	pwd, ok := u.EncodedPassword()
	require.True(t, ok)
	require.Equal(t, u.AuthOpt.HashString, pwd)

	u.AuthOpt.HashString = "not-good-password-format"
	pwd, ok = u.EncodedPassword()
	require.False(t, ok)

	u.AuthOpt.ByAuthString = true
	pwd, ok = u.EncodedPassword()
	require.True(t, ok)
	require.Equal(t, hashString, pwd)

	u.AuthOpt.AuthString = ""
	pwd, ok = u.EncodedPassword()
	require.True(t, ok)
	require.Equal(t, "", pwd)
}

func TestTableOptimizerHintRestore(t *testing.T) {
	testCases := []NodeRestoreTestCase{
		{"TIDB_SMJ(`t1`)", "TIDB_SMJ(`t1`)"},
		{"TIDB_SMJ(t1)", "TIDB_SMJ(`t1`)"},
		{"TIDB_SMJ(t1,t2)", "TIDB_SMJ(`t1`, `t2`)"},
		{"TIDB_INLJ(t1,t2)", "TIDB_INLJ(`t1`, `t2`)"},
		{"TIDB_HJ(t1,t2)", "TIDB_HJ(`t1`, `t2`)"},
		{"MAX_EXECUTION_TIME(3000)", "MAX_EXECUTION_TIME(3000)"},
	}
	extractNodeFunc := func(node Node) Node {
		return node.(*SelectStmt).TableHints[0]
	}
	runNodeRestoreTest(t, testCases, "select /*+ %s */ * from t1 join t2", extractNodeFunc)
}

// 测试用户变量 @my_var
func TestRestoreUserVariable(t *testing.T) {
	// 构造用户变量赋值
	va := &ast.VariableAssignment{
		Name:     "my_var",
		Value:    ast.NewValueExpr("hello"),
		IsSystem: false, // 明确标记为用户变量
	}

	// 执行 Restore 方法
	sb := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, sb)
	err := va.Restore(ctx)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	// 验证结果
	expected := "@`my_var`='hello'"
	actual := sb.String()
	if actual != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, actual)
	}
}

// 测试全局系统变量 @@GLOBAL.autocommit
func TestRestoreGlobalSystemVariable(t *testing.T) {
	va := &ast.VariableAssignment{
		Name:     "autocommit",
		Value:    ast.NewValueExpr(1),
		IsSystem: true,
		IsGlobal: true, // 标记为全局系统变量
	}

	sb := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, sb)
	err := va.Restore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := "@@GLOBAL.`autocommit`=1"
	actual := sb.String()
	if actual != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, actual)
	}
}

// 测试会话系统变量 @@SESSION.autocommit
func TestRestoreSessionSystemVariable(t *testing.T) {
	va := &ast.VariableAssignment{
		Name:     "autocommit",
		Value:    ast.NewValueExpr(0),
		IsSystem: true,
		IsGlobal: false, // 默认会话作用域
	}

	sb := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, sb)
	err := va.Restore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := "@@SESSION.`autocommit`=0"
	actual := sb.String()
	if actual != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, actual)
	}
}

// 测试 SET NAMES utf8mb4
func TestRestoreSetNames(t *testing.T) {
	va := &ast.VariableAssignment{
		Name:  ast.SetNames, // 特殊标记
		Value: ast.NewValueExpr("utf8mb4"),
	}

	sb := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, sb)
	err := va.Restore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := "NAMES 'utf8mb4'"
	actual := sb.String()
	if actual != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, actual)
	}
}

// 测试带 COLLATE 的 SET NAMES
func TestRestoreSetNamesWithCollate(t *testing.T) {
	va := &ast.VariableAssignment{
		Name:        ast.SetNames,
		Value:       ast.NewValueExpr("utf8mb4"),
		ExtendValue: ast.NewValueExpr("utf8mb4_bin"),
	}

	sb := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, sb)
	err := va.Restore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := "NAMES 'utf8mb4' COLLATE 'utf8mb4_bin'"
	actual := sb.String()
	if actual != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, actual)
	}
}

// 测试隐式作用域系统变量
func TestRestoreImplicitScopeSystemVariable(t *testing.T) {
	va := &ast.VariableAssignment{
		Name:     "autocommit",
		Value:    ast.NewValueExpr(1),
		IsSystem: true,
		IsGlobal: false, // 默认会话作用域但不显式写 SESSION
	}

	sb := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, sb)
	err := va.Restore(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := "@@SESSION.`autocommit`=1"
	actual := sb.String()
	if actual != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, actual)
	}
}
