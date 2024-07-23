package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitStatementToPieces(t *testing.T) {
	type testCase struct {
		originSql   string
		expectSqls  []string
		assertEqual bool
	}

	cases := []testCase{
		{
			"select 1;select 2;",
			[]string{"select 1", "select 2"},
			true,
		},
		{
			"/*test*/select 1;/*test*/ select 2;",
			[]string{"/*test*/select 1", "/*test*/ select 2"},
			true,
		},
		{
			"select 1;/*test;*/ select 2;",
			[]string{"select 1", "/*test;*/ select 2"},
			true,
		},
		{
			"-- test;select 1",
			[]string{"-- test", "select 1"},
			false,
		},
	}

	for _, tt := range cases {
		res, err := SplitStatementToPieces(tt.originSql)
		if err != nil {
			t.Fatalf("SplitStatementToPieces err.sql:%s\n", tt.originSql)
		}
		if tt.assertEqual && !isEqSlice(tt.expectSqls, res) {
			t.Fatalf("SplitStatementToPieces res not equal,sql:%s,res:%+q\n", tt.originSql, res)
		}
	}

}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "test simple",
			sql:  "select 1",
			want: []string{"select", "1"},
		},
		{
			name: "test simple with space",
			sql:  " select *         from t ",
			want: []string{"select", "*", "from", "t"},
		},
		{
			name: "test simple with semicolon",
			sql:  "select * from t;",
			want: []string{"select", "*", "from", "t;"},
		},
		{
			name: "test simple with semicolon and db name",
			sql:  "select a,b from A.t",
			want: []string{"select", "a", "b", "from", "A.t"},
		},
		{
			name: "test select subquery",
			sql:  "select * from (select * from t1);",
			want: []string{"select", "*", "from", "(select", "*", "from", "t1);"},
		},
		{
			name: "test select subquery 2",
			sql:  "select a.* from (select * from t1 where id>1000) as a where a.id<2000;",
			want: []string{"select", "a.*", "from", "(select", "*", "from", "t1", "where", "id>1000)", "as", "a", "where", "a.id<2000;"},
		},
		{
			name: "test simple with tab",
			sql: "select a	b from A.t",
			want: []string{"select", "a", "b", "from", "A.t"},
		},
		{
			name: "test simple with comment in sql",
			sql:  "select /*master*/ a,b from A.t",
			want: []string{"select", "*master*", "a", "b", "from", "A.t"},
		},
		{
			name: "test select master hint leading",
			sql:  "/*master*/ select a,b from A.t",
			want: []string{"select", "a", "b", "from", "A.t", "*master*"},
		},
		{
			name: "test select master hint training",
			sql:  "select a,b from t /*master*/",
			want: []string{"select", "a", "b", "from", "t", "*master*"},
		},
		{
			name: "test mycat comment in sql",
			sql:  "select a,b from A.t /* !mycat:sql=select 1 */ ",
			want: []string{"select", "a", "b", "from", "A.t", "*", "!mycat:sql=select", "1", "*"},
		},
		{
			name: "test with comment in sql 2",
			sql:  "select 1 -- test",
			want: []string{"select", "1", "--", "test"},
		},
		{
			name: "test with comment in multi line sql",
			sql: `-- aaa
			-- bbb
			ccc
			-- ddd
			eee`,
			want: []string{"ccc", "eee"},
		},
		{
			name: "test insert",
			sql:  "insert into tbl_unshard(col1) (select col1 from tbl_unshard_1);",
			want: []string{"insert", "into", "tbl_unshard(col1)", "(select", "col1", "from", "tbl_unshard_1);"},
		},
		{
			name: "test only comment",
			sql:  `/*slave*/`,
			want: []string{},
		},
		{
			name: "test special characters new line",
			sql:  "select * \nfrom t",
			want: []string{"select", "*", "from", "t"},
		},
		{
			name: "test special characters page Break",
			sql:  "select * \f from t",
			want: []string{"select", "*", "\f", "from", "t"},
		},
		{
			name: "test special characters vertical tabs",
			sql:  "select * \vfrom t",
			want: []string{"select", "*", "\vfrom", "t"},
		},
		{
			name: "test only comment",
			sql:  "select * \u00A0from t",
			want: []string{"select", "*", "\u00a0from", "t"},
		},
		{
			name: "test only comment",
			sql:  "select * from t\v",
			want: []string{"select", "*", "from", "t"},
		},
		{
			name: "test only comment",
			sql:  "select * from t\v",
			want: []string{"select", "*", "from", "t"},
		},
		{
			name: "test read_only",
			sql:  "/* hint */ select @@read_only",
			want: []string{"select", "@@read_only"},
		},
		{
			name: "test read_only",
			sql:  "select /* hint */ @@read_only",
			want: []string{"select", "*", "hint", "*", "@@read_only"},
		},
		{
			name: "test empty",
			sql:  " ",
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Tokenize(tt.sql), "Tokenize(%v)", tt.sql)
		})
	}
}

func isEqSlice(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestGetInsertDBTable(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantDB    string
		wantTable string
	}{
		{
			name:      "test simple",
			token:     "t",
			wantDB:    "",
			wantTable: "t",
		},
		{
			name:      "test simple with db",
			token:     "db.t",
			wantDB:    "db",
			wantTable: "t",
		},
		{
			name:      "test simple with db and quote",
			token:     "`db`.`t`",
			wantDB:    "db",
			wantTable: "t",
		},
		{
			name:      "test simple with db and bracket",
			token:     "`db`.`t`(col1,col2)",
			wantDB:    "db",
			wantTable: "t",
		},
		{
			name:      "test simple without db and quote",
			token:     "t(col1,col2)",
			wantDB:    "",
			wantTable: "t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDB, gotTable := GetInsertDBTable(tt.token)
			assert.Equalf(t, tt.wantDB, gotDB, "GetInsertDBTable(%v)", tt.token)
			assert.Equalf(t, tt.wantTable, gotTable, "GetInsertDBTable(%v)", tt.token)
		})
	}
}

func TestPreview(t *testing.T) {
	testcases := []struct {
		sql  string
		want int
	}{
		{"select ...", StmtSelect},
		{"    select ...", StmtSelect},
		{"(select ...", StmtSelect},
		{"( select ...", StmtSelect},
		{"insert ...", StmtInsert},
		{"replace ....", StmtReplace},
		{"   update ...", StmtUpdate},
		{"Update", StmtUpdate},
		{"UPDATE ...", StmtUpdate},
		{"\n\t    delete ...", StmtDelete},
		{"", StmtUnknown},
		{" ", StmtUnknown},
		{"begin", StmtBegin},
		{" begin", StmtBegin},
		{" begin ", StmtBegin},
		{"\n\t begin ", StmtBegin},
		{"... begin ", StmtUnknown},
		{"begin ...", StmtUnknown},
		{"begin /* ... */", StmtBegin},
		{"begin /* ... *//*test*/", StmtBegin},
		{"begin;", StmtBegin},
		{"begin ;", StmtBegin},
		{"begin; /*...*/", StmtBegin},
		{"start transaction", StmtBegin},
		{"commit", StmtCommit},
		{"commit /*...*/", StmtCommit},
		{"rollback", StmtRollback},
		{"rollback /*...*/", StmtRollback},
		{"rollback to point", StmeSRollback},
		{"rollback to point /*...*/", StmeSRollback},
		{"savepoint point", StmtSavepoint},
		{"release savepoint point", StmtRelease},
		{"release savepoint point /*...*/", StmtRelease},
		{"lock tables", StmtLockTables},
		{"lock table t1 read, t2 write", StmtLockTables},
		{"unlock tables", StmtUnlockTables},
		{"flush", StmtFlush},
		{"create", StmtDDL},
		{"alter", StmtDDL},
		{"rename", StmtDDL},
		{"drop", StmtDDL},
		{"set", StmtSet},
		{"show", StmtShow},
		{"use", StmtUse},
		{"repair", StmtOther},
		{"optimize", StmtOther},
		{"truncate", StmtDDL},
		{"explain ", StmtExplain},
		{"unknown", StmtUnknown},

		{"/* leading comment */ select ...", StmtSelect},
		{"/* leading comment */ (select ...", StmtSelect},
		{"/* leading comment */ /* leading comment 2 */ select ...", StmtSelect},
		{"/*! MySQL-specific comment */", StmtComment},
		{"/*!50708 MySQL-version comment */", StmtComment},
		{"-- leading single line comment \n select ...", StmtSelect},
		{"-- leading single line comment \n -- leading single line comment 2\n select ...", StmtSelect},

		{"/* leading comment no end select ...", StmtUnknown},
		{"-- leading single line comment no end select ...", StmtUnknown},
		{"/*!40000 ALTER TABLE `t1` DISABLE KEYS */", StmtComment}}

	for _, tt := range testcases {
		t.Run(tt.sql, func(t *testing.T) {
			if got := Preview(tt.sql); got != tt.want {
				t.Errorf("Preview(%s): %v, want %v", tt.sql, got, tt.want)
			}
		})

	}
}
