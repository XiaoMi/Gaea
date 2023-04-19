package parser

import "testing"

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
