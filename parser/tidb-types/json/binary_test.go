// Copyright 2017 PingCAP, Inc.
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

package json

import (
	"testing"

	. "github.com/pingcap/check"
	"github.com/stretchr/testify/require"
)

func TestT(t *testing.T) {
	TestingT(t)
}

func TestBinaryJSONMarshalUnmarshal(t *testing.T) {
	strs := []string{
		`{"a": [1, "2", {"aa": "bb"}, 4, null], "b": true, "c": null}`,
		`{"aaaaaaaaaaa": [1, "2", {"aa": "bb"}, 4.1], "bbbbbbbbbb": true, "ccccccccc": "d"}`,
		`[{"a": 1, "b": true}, 3, 3.5, "hello, world", null, true]`,
	}
	for _, str := range strs {
		parsedBJ := mustParseBinaryFromString(t, str)
		require.Equal(t, str, parsedBJ.String())
	}
}

func TestBinaryJSONExtract(t *testing.T) {
	bj1 := mustParseBinaryFromString(t, `{"\"hello\"": "world", "a": [1, "2", {"aa": "bb"}, 4.0, {"aa": "cc"}], "b": true, "c": ["d"]}`)
	bj2 := mustParseBinaryFromString(t, `[{"a": 1, "b": true}, 3, 3.5, "hello, world", null, true]`)

	var tests = []struct {
		bj              BinaryJSON
		pathExprStrings []string
		expected        BinaryJSON
		found           bool
		err             error
	}{
		// test extract with only one path expression.
		{bj1, []string{"$.a"}, mustParseBinaryFromString(t, `[1, "2", {"aa": "bb"}, 4.0, {"aa": "cc"}]`), true, nil},
		{bj2, []string{"$.a"}, mustParseBinaryFromString(t, "null"), false, nil},
		{bj1, []string{"$[0]"}, bj1, true, nil}, // in Extract, autowraped bj1 as an array.
		{bj2, []string{"$[0]"}, mustParseBinaryFromString(t, `{"a": 1, "b": true}`), true, nil},
		{bj1, []string{"$.a[2].aa"}, mustParseBinaryFromString(t, `"bb"`), true, nil},
		{bj1, []string{"$.a[*].aa"}, mustParseBinaryFromString(t, `["bb", "cc"]`), true, nil},
		{bj1, []string{"$.*[0]"}, mustParseBinaryFromString(t, `["world", 1, true, "d"]`), true, nil},
		{bj1, []string{`$.a[*]."aa"`}, mustParseBinaryFromString(t, `["bb", "cc"]`), true, nil},
		{bj1, []string{`$."\"hello\""`}, mustParseBinaryFromString(t, `"world"`), true, nil},
		{bj1, []string{`$**[1]`}, mustParseBinaryFromString(t, `"2"`), true, nil},

		// test extract with multi path expressions.
		{bj1, []string{"$.a", "$[5]"}, mustParseBinaryFromString(t, `[[1, "2", {"aa": "bb"}, 4.0, {"aa": "cc"}]]`), true, nil},
		{bj2, []string{"$.a", "$[0]"}, mustParseBinaryFromString(t, `[{"a": 1, "b": true}]`), true, nil},
	}

	for _, tt := range tests {
		var pathExprList = make([]PathExpression, 0)
		for _, peStr := range tt.pathExprStrings {
			pe, err := ParseJSONPathExpr(peStr)
			require.NoError(t, err)
			pathExprList = append(pathExprList, pe)
		}

		result, found := tt.bj.Extract(pathExprList)
		require.Equal(t, tt.found, found, found)
		if found {
			require.Equal(t, tt.expected.String(), result.String())
		}
	}
}

func TestBinaryJSONType(t *testing.T) {
	var tests = []struct {
		In  string
		Out string
	}{
		{`{"a": "b"}`, "OBJECT"},
		{`["a", "b"]`, "ARRAY"},
		{`3`, "INTEGER"},
		{`3.0`, "DOUBLE"},
		{`null`, "NULL"},
		{`true`, "BOOLEAN"},
	}
	for _, tt := range tests {
		bj := mustParseBinaryFromString(t, tt.In)
		require.Equal(t, tt.Out, bj.Type())
	}
	// we can't parse '9223372036854775808' to JSON::Uint64 now,
	// because go builtin JSON parser treats that as DOUBLE.
	require.Equal(t, "UNSIGNED INTEGER", CreateBinary(uint64(1<<63)).Type())
}

func TestBinaryJSONUnquote(t *testing.T) {
	var tests = []struct {
		j        string
		unquoted string
	}{
		{j: `3`, unquoted: "3"},
		{j: `"3"`, unquoted: "3"},
		{j: `"hello, \"escaped quotes\" world"`, unquoted: "hello, \"escaped quotes\" world"},
		{j: "\"\\u4f60\"", unquoted: "你"},
		{j: `true`, unquoted: "true"},
		{j: `null`, unquoted: "null"},
		{j: `{"a": [1, 2]}`, unquoted: `{"a": [1, 2]}`},
		{j: `"\""`, unquoted: `"`},
		{j: `"'"`, unquoted: `'`},
		{j: `"''"`, unquoted: `''`},
		{j: `""`, unquoted: ``},
	}
	for _, tt := range tests {
		bj := mustParseBinaryFromString(t, tt.j)
		unquoted, err := bj.Unquote()
		require.NoError(t, err)
		require.Equal(t, tt.unquoted, unquoted)
	}
}

func TestBinaryJSONModify(t *testing.T) {
	var tests = []struct {
		base     string
		setField string
		setValue string
		mt       ModifyType
		expected string
		success  bool
	}{
		{`null`, "$", `{}`, ModifySet, `{}`, true},
		{`{}`, "$.a", `3`, ModifySet, `{"a": 3}`, true},
		{`{"a": 3}`, "$.a", `[]`, ModifyReplace, `{"a": []}`, true},
		{`{"a": 3}`, "$.b", `"3"`, ModifySet, `{"a": 3, "b": "3"}`, true},
		{`{"a": []}`, "$.a[0]", `3`, ModifySet, `{"a": [3]}`, true},
		{`{"a": [3]}`, "$.a[1]", `4`, ModifyInsert, `{"a": [3, 4]}`, true},
		{`{"a": [3]}`, "$[0]", `4`, ModifySet, `4`, true},
		{`{"a": [3]}`, "$[1]", `4`, ModifySet, `[{"a": [3]}, 4]`, true},
		{`{"b": true}`, "$.b", `false`, ModifySet, `{"b": false}`, true},

		// nothing changed because the path is empty and we want to insert.
		{`{}`, "$", `1`, ModifyInsert, `{}`, true},
		// nothing changed because the path without last leg doesn't exist.
		{`{"a": [3, 4]}`, "$.b[1]", `3`, ModifySet, `{"a": [3, 4]}`, true},
		// nothing changed because the path without last leg doesn't exist.
		{`{"a": [3, 4]}`, "$.a[2].b", `3`, ModifySet, `{"a": [3, 4]}`, true},
		// nothing changed because we want to insert but the full path exists.
		{`{"a": [3, 4]}`, "$.a[0]", `30`, ModifyInsert, `{"a": [3, 4]}`, true},
		// nothing changed because we want to replace but the full path doesn't exist.
		{`{"a": [3, 4]}`, "$.a[2]", `30`, ModifyReplace, `{"a": [3, 4]}`, true},

		// bad path expression.
		{"null", "$.*", "{}", ModifySet, "null", false},
		{"null", "$[*]", "{}", ModifySet, "null", false},
		{"null", "$**.a", "{}", ModifySet, "null", false},
		{"null", "$**[3]", "{}", ModifySet, "null", false},
	}
	for _, tt := range tests {
		pathExpr, err := ParseJSONPathExpr(tt.setField)
		require.NoError(t, err)

		base := mustParseBinaryFromString(t, tt.base)
		value := mustParseBinaryFromString(t, tt.setValue)
		expected := mustParseBinaryFromString(t, tt.expected)
		obtain, err := base.Modify([]PathExpression{pathExpr}, []BinaryJSON{value}, tt.mt)
		if tt.success {
			require.NoError(t, err)
			require.Equal(t, expected.String(), obtain.String())
		} else {
			require.Error(t, err)
		}
	}
}

func TestBinaryJSONRemove(t *testing.T) {
	var tests = []struct {
		base     string
		path     string
		expected string
		success  bool
	}{
		{`null`, "$", `{}`, false},
		{`{"a":[3]}`, "$.a[*]", `{"a":[3]}`, false},
		{`{}`, "$.a", `{}`, true},
		{`{"a":3}`, "$.a", `{}`, true},
		{`{"a":1,"b":2,"c":3}`, "$.b", `{"a":1,"c":3}`, true},
		{`{"a":1,"b":2,"c":3}`, "$.d", `{"a":1,"b":2,"c":3}`, true},
		{`{"a":3}`, "$[0]", `{"a":3}`, true},
		{`{"a":[3,4,5]}`, "$.a[0]", `{"a":[4,5]}`, true},
		{`{"a":[3,4,5]}`, "$.a[1]", `{"a":[3,5]}`, true},
		{`{"a":[3,4,5]}`, "$.a[4]", `{"a":[3,4,5]}`, true},
		{`{"a": [1, 2, {"aa": "xx"}]}`, "$.a[2].aa", `{"a": [1, 2, {}]}`, true},
	}
	for _, tt := range tests {
		pathExpr, err := ParseJSONPathExpr(tt.path)
		require.NoError(t, err)

		base := mustParseBinaryFromString(t, tt.base)
		expected := mustParseBinaryFromString(t, tt.expected)
		obtain, err := base.Remove([]PathExpression{pathExpr})
		if tt.success {
			require.NoError(t, err)
			require.Equal(t, expected.String(), obtain.String())
		} else {
			require.Error(t, err)
		}
	}
}

func TestCompareBinary(t *testing.T) {
	jNull := mustParseBinaryFromString(t, `null`)
	jBoolTrue := mustParseBinaryFromString(t, `true`)
	jBoolFalse := mustParseBinaryFromString(t, `false`)
	jIntegerLarge := CreateBinary(uint64(1 << 63))
	jIntegerSmall := mustParseBinaryFromString(t, `3`)
	jStringLarge := mustParseBinaryFromString(t, `"hello, world"`)
	jStringSmall := mustParseBinaryFromString(t, `"hello"`)
	jArrayLarge := mustParseBinaryFromString(t, `["a", "c"]`)
	jArraySmall := mustParseBinaryFromString(t, `["a", "b"]`)
	jObject := mustParseBinaryFromString(t, `{"a": "b"}`)

	var tests = []struct {
		left  BinaryJSON
		right BinaryJSON
	}{
		{jNull, jIntegerSmall},
		{jIntegerSmall, jIntegerLarge},
		{jIntegerLarge, jStringSmall},
		{jStringSmall, jStringLarge},
		{jStringLarge, jObject},
		{jObject, jArraySmall},
		{jArraySmall, jArrayLarge},
		{jArrayLarge, jBoolFalse},
		{jBoolFalse, jBoolTrue},
	}
	for _, tt := range tests {
		cmp := CompareBinary(tt.left, tt.right)
		require.True(t, cmp < 0)
	}
}

func TestBinaryJSONMerge(t *testing.T) {
	var tests = []struct {
		suffixes []string
		expected string
	}{
		{[]string{`{"a": 1}`, `{"b": 2}`}, `{"a": 1, "b": 2}`},
		{[]string{`{"a": 1}`, `{"a": 2}`}, `{"a": [1, 2]}`},
		{[]string{`[1]`, `[2]`}, `[1, 2]`},
		{[]string{`{"a": 1}`, `[1]`}, `[{"a": 1}, 1]`},
		{[]string{`[1]`, `{"a": 1}`}, `[1, {"a": 1}]`},
		{[]string{`{"a": 1}`, `4`}, `[{"a": 1}, 4]`},
		{[]string{`[1]`, `4`}, `[1, 4]`},
		{[]string{`4`, `{"a": 1}`}, `[4, {"a": 1}]`},
		{[]string{`4`, `1`}, `[4, 1]`},
		{[]string{`{}`, `[]`}, `[{}]`},
	}

	for _, tt := range tests {
		suffixes := make([]BinaryJSON, 0, len(tt.suffixes)+1)
		for _, s := range tt.suffixes {
			suffixes = append(suffixes, mustParseBinaryFromString(t, s))
		}
		result := MergeBinary(suffixes)
		cmp := CompareBinary(result, mustParseBinaryFromString(t, tt.expected))
		require.Equal(t, 0, cmp)
	}
}

func mustParseBinaryFromString(t *testing.T, s string) BinaryJSON {
	bj, err := ParseBinaryFromString(s)
	require.NoError(t, err)
	return bj
}

const benchStr = `{"a":[1,"2",{"aa":"bb"},4,null],"b":true,"c":null}`

func BenchmarkBinaryMarshal(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(int64(len(benchStr)))
	bj, _ := ParseBinaryFromString(benchStr)
	for i := 0; i < b.N; i++ {
		bj.MarshalJSON()
	}
}

func TestBinaryJSONContains(t *testing.T) {
	var tests = []struct {
		input    string
		target   string
		expected bool
	}{
		{`{}`, `{}`, true},
		{`{"a":1}`, `{}`, true},
		{`{"a":1}`, `1`, false},
		{`{"a":[1]}`, `[1]`, false},
		{`{"b":2, "c":3}`, `{"c":3}`, true},
		{`1`, `1`, true},
		{`[1]`, `1`, true},
		{`[1,2]`, `[1]`, true},
		{`[1,2]`, `[1,3]`, false},
		{`[1,2]`, `["1"]`, false},
		{`[1,2,[1,3]]`, `[1,3]`, true},
		{`[1,2,[1,[5,[3]]]]`, `[1,3]`, true},
		{`[1,2,[1,[5,{"a":[2,3]}]]]`, `[1,{"a":[3]}]`, true},
		{`[{"a":1}]`, `{"a":1}`, true},
		{`[{"a":1,"b":2}]`, `{"a":1}`, true},
		{`[{"a":{"a":1},"b":2}]`, `{"a":1}`, false},
	}

	for _, tt := range tests {
		obj := mustParseBinaryFromString(t, tt.input)
		target := mustParseBinaryFromString(t, tt.target)
		require.Equal(t, tt.expected, ContainsBinary(obj, target))
	}
}

func TestBinaryJSONDepth(t *testing.T) {
	var tests = []struct {
		input    string
		expected int
	}{
		{`{}`, 1},
		{`[]`, 1},
		{`true`, 1},
		{`[10, 20]`, 2},
		{`[[], {}]`, 2},
		{`[10, {"a": 20}]`, 3},
		{`{"Person": {"Name": "Homer", "Age": 39, "Hobbies": ["Eating", "Sleeping"]} }`, 4},
	}

	for _, tt := range tests {
		obj := mustParseBinaryFromString(t, tt.input)
		require.Equal(t, tt.expected, obj.GetElemDepth())
	}
}
