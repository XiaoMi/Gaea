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

package parser

import (
	"fmt"
	"testing"
	"unicode"

	"github.com/XiaoMi/Gaea/mysql"

	requires "github.com/stretchr/testify/require"
)

func TestTokenID(t *testing.T) {
	for str, tok := range tokenMap {
		l := NewScanner(str)
		var v yySymType
		tok1 := l.Lex(&v)
		requires.Equal(t, tok1, tok)
	}
}

func TestSingleChar(t *testing.T) {
	table := []byte{'|', '&', '-', '+', '*', '/', '%', '^', '~', '(', ',', ')'}
	for _, tok := range table {
		l := NewScanner(string(tok))
		var v yySymType
		tok1 := l.Lex(&v)
		requires.Equal(t, tok1, int(tok))
	}
}

type testCaseItem struct {
	str string
	tok int
}

func TestSingleCharOther(t *testing.T) {
	table := []testCaseItem{
		{"AT", identifier},
		{"?", paramMarker},
		{"PLACEHOLDER", identifier},
		{"=", eq},
		{".", int('.')},
	}
	runTest(t, table)
}

func TestAtLeadingIdentifier(t *testing.T) {
	table := []testCaseItem{
		{"@", singleAtIdentifier},
		{"@''", singleAtIdentifier},
		{"@1", singleAtIdentifier},
		{"@.1_", singleAtIdentifier},
		{"@-1.", singleAtIdentifier},
		{"@~", singleAtIdentifier},
		{"@$", singleAtIdentifier},
		{"@a_3cbbc", singleAtIdentifier},
		{"@`a_3cbbc`", singleAtIdentifier},
		{"@-3cbbc", singleAtIdentifier},
		{"@!3cbbc", singleAtIdentifier},
		{"@@global.test", doubleAtIdentifier},
		{"@@session.test", doubleAtIdentifier},
		{"@@local.test", doubleAtIdentifier},
		{"@@test", doubleAtIdentifier},
		{"@@global.`test`", doubleAtIdentifier},
		{"@@session.`test`", doubleAtIdentifier},
		{"@@local.`test`", doubleAtIdentifier},
		{"@@`test`", doubleAtIdentifier},
	}
	runTest(t, table)
}

func TestUnderscoreCS(t *testing.T) {
	var v yySymType
	scanner := NewScanner(`_utf8"string"`)
	tok := scanner.Lex(&v)
	requires.Equal(t, underscoreCS, tok)
	tok = scanner.Lex(&v)
	requires.Equal(t, stringLit, tok)

	scanner.reset("N'string'")
	tok = scanner.Lex(&v)
	requires.Equal(t, underscoreCS, tok)
	tok = scanner.Lex(&v)
	requires.Equal(t, stringLit, tok)
}

func TestLiteral(t *testing.T) {
	table := []testCaseItem{
		{`'''a'''`, stringLit},
		{`''a''`, stringLit},
		{`""a""`, stringLit},
		{`\'a\'`, int('\\')},
		{`\"a\"`, int('\\')},
		{"0.2314", decLit},
		{"1234567890123456789012345678901234567890", decLit},
		{"132.313", decLit},
		{"132.3e231", floatLit},
		{"132.3e-231", floatLit},
		{"001e-12", floatLit},
		{"23416", intLit},
		{"123test", identifier},
		{"123" + string(unicode.ReplacementChar) + "xxx", identifier},
		{"0", intLit},
		{"0x3c26", hexLit},
		{"x'13181C76734725455A'", hexLit},
		{"0b01", bitLit},
		{fmt.Sprintf("t1%c", 0), identifier},
		{"N'some text'", underscoreCS},
		{"n'some text'", underscoreCS},
		{"\\N", null},
		{".*", int('.')},       // `.`, `*`
		{".1_t_1_x", int('.')}, // `.`, `1_t_1_x`
		{"9e9e", floatLit},     // 9e9e = 9e9 + e
		// Issue #3954
		{".1e23", floatLit}, // `.1e23`
		{".123", decLit},    // `.123`
		{".1*23", decLit},   // `.1`, `*`, `23`
		{".1,23", decLit},   // `.1`, `,`, `23`
		{".1 23", decLit},   // `.1`, `23`
		// TODO: See #3963. The following test cases do not test the ambiguity.
		{".1$23", int('.')},    // `.`, `1$23`
		{".1a23", int('.')},    // `.`, `1a23`
		{".1e23$23", int('.')}, // `.`, `1e23$23`
		{".1e23a23", int('.')}, // `.`, `1e23a23`
		{".1C23", int('.')},    // `.`, `1C23`
		{".1\u0081", int('.')}, // `.`, `1\u0081`
		{".1\uff34", int('.')}, // `.`, `1\uff34`
		{`b''`, bitLit},
		{`b'0101'`, bitLit},
		{`0b0101`, bitLit},
	}
	runTest(t, table)
}

func runTest(t *testing.T, table []testCaseItem) {
	var val yySymType
	for _, v := range table {
		l := NewScanner(v.str)
		tok := l.Lex(&val)
		requires.Equal(t, v.tok, tok, v.str)
	}
}
func TestComment(t *testing.T) {

	table := []testCaseItem{
		{"-- select --\n1", intLit},
		{"/*!40101 SET character_set_client = utf8 */;", set},
		{"/*+ BKA(t1) */", hintBegin},
		{"/* SET character_set_client = utf8 */;", int(';')},
		{"/* some comments */ SELECT ", selectKwd},
		{`-- comment continues to the end of line
SELECT`, selectKwd},
		{`# comment continues to the end of line
SELECT`, selectKwd},
		{"#comment\n123", intLit},
		{"--5", int('-')},
		{"--\nSELECT", selectKwd},
		{"--\tSELECT", 0},
		{"--\r\nSELECT", selectKwd},
		{"--", 0},
	}
	runTest(t, table)
}

func TestScanQuotedIdent(t *testing.T) {
	l := NewScanner("`fk`")
	l.r.peek()
	tok, pos, lit := scanQuotedIdent(l)
	requires.Zero(t, pos.Offset)
	requires.Equal(t, quotedIdentifier, tok)
	requires.Equal(t, "fk", lit)
}

func TestscanString(t *testing.T) {
	table := []struct {
		raw    string
		expect string
	}{
		{`' \n\tTest String'`, " \n\tTest String"},
		{`'\x\B'`, "xB"},
		{`'\0\'\"\b\n\r\t\\'`, "\000'\"\b\n\r\t\\"},
		{`'\Z'`, string(rune(26))},
		{`'\%\_'`, `\%\_`},
		{`'hello'`, "hello"},
		{`'"hello"'`, `"hello"`},
		{`'""hello""'`, `""hello""`},
		{`'hel''lo'`, "hel'lo"},
		{`'\'hello'`, "'hello"},
		{`"hello"`, "hello"},
		{`"'hello'"`, "'hello'"},
		{`"''hello''"`, "''hello''"},
		{`"hel""lo"`, `hel"lo`},
		{`"\"hello"`, `"hello`},
		{`'disappearing\ backslash'`, "disappearing backslash"},
		{"'한국의中文UTF8およびテキストトラック'", "한국의中文UTF8およびテキストトラック"},
		{"'\\a\x90'", "a\x90"},
		{`"\aèàø»"`, `aèàø»`},
	}
	for _, v := range table {
		l := NewScanner(v.raw)
		tok, pos, lit := l.scan()
		requires.Zero(t, pos.Offset)
		requires.Equal(t, stringLit, tok)
		requires.Equal(t, v.expect, lit)
	}
}

func TestIdentifier(t *testing.T) {
	replacementString := string(unicode.ReplacementChar) + "xxx"
	table := [][2]string{
		{`哈哈`, "哈哈"},
		{"`numeric`", "numeric"},
		{"\r\n \r \n \tthere\t \n", "there"},
		{`5number`, `5number`},
		{"1_x", "1_x"},
		{"0_x", "0_x"},
		{replacementString, replacementString},
		{"9e", "9e"},
		{"0b", "0b"},
		{"0b123", "0b123"},
		{"0b1ab", "0b1ab"},
		{"0B01", "0B01"},
		{"0x", "0x"},
		{"0x7fz3", "0x7fz3"},
		{"023a4", "023a4"},
		{"9eTSs", "9eTSs"},
		{fmt.Sprintf("t1%cxxx", 0), "t1"},
	}
	l := &Scanner{}
	for _, item := range table {
		l.reset(item[0])
		var v yySymType
		tok := l.Lex(&v)
		requires.Equal(t, identifier, tok, item)
		requires.Equal(t, item[1], v.ident, item)
	}
}

func TestSpecialComment(t *testing.T) {
	l := NewScanner("/*!40101 select\n5*/")
	tok, pos, lit := l.scan()
	requires.Equal(t, identifier, tok)
	requires.Equal(t, "select", lit)
	requires.Equal(t, Pos{0, 0, 9}, pos)

	tok, pos, lit = l.scan()
	requires.Equal(t, intLit, tok)
	requires.Equal(t, "5", lit)
	requires.Equal(t, Pos{1, 1, 16}, pos)
}

func TestOptimizerHint(t *testing.T) {
	l := NewScanner("  /*+ BKA(t1) */")
	tokens := []struct {
		tok int
		lit string
		pos int
	}{
		{hintBegin, "", 2},
		{identifier, "BKA", 6},
		{int('('), "(", 9},
		{identifier, "t1", 10},
		{int(')'), ")", 12},
		{hintEnd, "", 14},
	}
	for i := 0; ; i++ {
		tok, pos, lit := l.scan()
		if tok == 0 {
			return
		}
		requires.Equal(t, tokens[i].tok, tok, i)
		requires.Equal(t, tokens[i].lit, lit, i)
		requires.Equal(t, tokens[i].pos, pos.Offset, i)
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		input  string
		expect uint64
	}{
		{"01000001783", 1000001783},
		{"00001783", 1783},
		{"0", 0},
		{"0000", 0},
		{"01", 1},
		{"10", 10},
	}
	scanner := NewScanner("")
	for _, test := range tests {
		var v yySymType
		scanner.reset(test.input)
		tok := scanner.Lex(&v)
		requires.Equal(t, intLit, tok)
		switch i := v.item.(type) {
		case int64:
			requires.Equal(t, test.expect, uint64(i))
		case uint64:
			requires.Equal(t, test.expect, i)
		default:
			t.Fail()
		}
	}
}

func TestSQLModeANSIQuotes(t *testing.T) {
	tests := []struct {
		input string
		tok   int
		ident string
	}{
		{`"identifier"`, identifier, "identifier"},
		{"`identifier`", identifier, "identifier"},
		{`"identifier""and"`, identifier, `identifier"and`},
		{`'string''string'`, stringLit, "string'string"},
		{`"identifier"'and'`, identifier, "identifier"},
		{`'string'"identifier"`, stringLit, "string"},
	}
	scanner := NewScanner("")
	scanner.SetSQLMode(mysql.ModeANSIQuotes)
	for _, test := range tests {
		var v yySymType
		scanner.reset(test.input)
		tok := scanner.Lex(&v)
		requires.Equal(t, test.tok, tok)
		requires.Equal(t, test.ident, v.ident)
	}
	scanner.reset(`'string' 'string'`)
	var v yySymType
	tok := scanner.Lex(&v)
	requires.Equal(t, stringLit, tok)
	requires.Equal(t, "string", v.ident)
	tok = scanner.Lex(&v)
	requires.Equal(t, stringLit, tok)
	requires.Equal(t, "string", v.ident)

}

func TestIllegal(t *testing.T) {
	table := []testCaseItem{
		{"'", 0},
		{"'fu", 0},
		{"'\\n", 0},
		{"'\\", 0},
		{fmt.Sprintf("%c", 0), invalid},
		{"`", 0},
		{`"`, 0},
		{"@`", 0},
		{"@'", 0},
		{`@"`, 0},
		{"@@`", 0},
		{"@@global.`", 0},
	}
	runTest(t, table)
}
