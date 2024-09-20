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

package types

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/stmtctx"
	"github.com/XiaoMi/Gaea/parser/tidb-types/json"
)

func TestDatum(t *testing.T) {
	values := []interface{}{
		int64(1),
		uint64(1),
		1.1,
		"abc",
		[]byte("abc"),
		[]int{1},
	}
	for _, val := range values {
		var d Datum
		d.SetValue(val)
		x := d.GetValue()
		require.Equal(t, val, x)
	}
}

func testDatumToBool(t *testing.T, in interface{}, res int) {
	datum := NewDatum(in)
	res64 := int64(res)
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = true
	b, err := datum.ToBool(sc)
	require.NoError(t, err)
	require.Equal(t, res64, b)
}

func TestToBool(t *testing.T) {
	testDatumToBool(t, int(0), 0)
	testDatumToBool(t, int64(0), 0)
	testDatumToBool(t, uint64(0), 0)
	testDatumToBool(t, float32(0.1), 0)
	testDatumToBool(t, float64(0.1), 0)
	testDatumToBool(t, "", 0)
	testDatumToBool(t, "0.1", 0)
	testDatumToBool(t, []byte{}, 0)
	testDatumToBool(t, []byte("0.1"), 0)
	testDatumToBool(t, NewBinaryLiteralFromUint(0, -1), 0)
	testDatumToBool(t, Enum{Name: "a", Value: 1}, 1)
	testDatumToBool(t, Set{Name: "a", Value: 1}, 1)

	t1, err := ParseTime(&stmtctx.StatementContext{TimeZone: time.UTC}, "2011-11-10 11:11:11.999999", mysql.TypeTimestamp, 6)
	require.NoError(t, err)
	testDatumToBool(t, t1, 1)

	td, err := ParseDuration(nil, "11:11:11.999999", 6)
	require.NoError(t, err)
	testDatumToBool(t, td, 1)

	ft := NewFieldType(mysql.TypeNewDecimal)
	ft.Decimal = 5
	v, err := Convert(0.1415926, ft)
	require.NoError(t, err)
	testDatumToBool(t, v, 0)
	d := NewDatum(&invalidMockType{})
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = true
	_, err = d.ToBool(sc)
	require.Error(t, err)
}

func TestEqualDatums(t *testing.T) {
	tests := []struct {
		a    []interface{}
		b    []interface{}
		same bool
	}{
		// Positive cases
		{[]interface{}{1}, []interface{}{1}, true},
		{[]interface{}{1, "aa"}, []interface{}{1, "aa"}, true},
		{[]interface{}{1, "aa", 1}, []interface{}{1, "aa", 1}, true},

		// negative cases
		{[]interface{}{1}, []interface{}{2}, false},
		{[]interface{}{1, "a"}, []interface{}{1, "aaaaaa"}, false},
		{[]interface{}{1, "aa", 3}, []interface{}{1, "aa", 2}, false},

		// Corner cases
		{[]interface{}{}, []interface{}{}, true},
		{[]interface{}{nil}, []interface{}{nil}, true},
		{[]interface{}{}, []interface{}{1}, false},
		{[]interface{}{1}, []interface{}{1, 1}, false},
		{[]interface{}{nil}, []interface{}{1}, false},
	}
	for _, tt := range tests {
		testEqualDatums(t, tt.a, tt.b, tt.same)
	}
}

func testEqualDatums(t *testing.T, a []interface{}, b []interface{}, same bool) {
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = true
	res, err := EqualDatums(sc, MakeDatums(a...), MakeDatums(b...))
	require.NoError(t, err)
	require.Equal(t, same, res)
}

func testDatumToInt64(t *testing.T, val interface{}, expect int64) {
	d := NewDatum(val)
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = true
	b, err := d.ToInt64(sc)
	require.NoError(t, err)
	require.Equal(t, expect, b)
}

func TestToInt64(t *testing.T) {
	testDatumToInt64(t, "0", int64(0))
	testDatumToInt64(t, int(0), int64(0))
	testDatumToInt64(t, int64(0), int64(0))
	testDatumToInt64(t, uint64(0), int64(0))
	testDatumToInt64(t, float32(3.1), int64(3))
	testDatumToInt64(t, float64(3.1), int64(3))
	testDatumToInt64(t, NewBinaryLiteralFromUint(100, -1), int64(100))
	testDatumToInt64(t, Enum{Name: "a", Value: 1}, int64(1))
	testDatumToInt64(t, Set{Name: "a", Value: 1}, int64(1))
	testDatumToInt64(t, json.CreateBinary(int64(3)), int64(3))

	t1, err := ParseTime(&stmtctx.StatementContext{
		TimeZone: time.UTC,
	}, "2011-11-10 11:11:11.999999", mysql.TypeTimestamp, 0)
	require.NoError(t, err)
	testDatumToInt64(t, t1, int64(20111110111112))

	td, err := ParseDuration(nil, "11:11:11.999999", 6)
	require.NoError(t, err)
	testDatumToInt64(t, td, int64(111112))

	ft := NewFieldType(mysql.TypeNewDecimal)
	ft.Decimal = 5
	v, err := Convert(3.1415926, ft)
	require.NoError(t, err)
	testDatumToInt64(t, v, int64(3))

	binLit, err := ParseHexStr("0x9999999999999999999999999999999999999999999")
	require.NoError(t, err)
	testDatumToInt64(t, binLit, -1)
}

func TestToFloat32(t *testing.T) {
	ft := NewFieldType(mysql.TypeFloat)
	var datum = NewFloat64Datum(281.37)
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = true
	converted, err := datum.ConvertTo(sc, ft)
	require.NoError(t, err)
	require.Equal(t, KindFloat32, converted.Kind())
	require.Equal(t, float32(281.37), converted.GetFloat32())

	datum.SetString("281.37")
	converted, err = datum.ConvertTo(sc, ft)
	require.NoError(t, err)
	require.Equal(t, KindFloat32, converted.Kind())
	require.Equal(t, float32(281.37), converted.GetFloat32())

	ft = NewFieldType(mysql.TypeDouble)
	datum = NewFloat32Datum(281.37)
	converted, err = datum.ConvertTo(sc, ft)
	require.NoError(t, err)
	require.Equal(t, KindFloat64, converted.Kind())
	// Convert to float32 and convert back to float64, we will get a different value.
	require.NotEqual(t, 281.37, converted.GetFloat64())
	require.Equal(t, datum.GetFloat64(), converted.GetFloat64())

}

// mustParseTimeIntoDatum is similar to ParseTime but panic if any error occurs.
func mustParseTimeIntoDatum(s string, tp byte, fsp int) (d Datum) {
	t, err := ParseTime(&stmtctx.StatementContext{TimeZone: time.UTC}, s, tp, fsp)
	if err != nil {
		panic("ParseTime fail")
	}
	d.SetMysqlTime(t)
	return
}

func TestToJSON(t *testing.T) {
	ft := NewFieldType(mysql.TypeJSON)
	sc := new(stmtctx.StatementContext)
	tests := []struct {
		datum    Datum
		expected string
		success  bool
	}{
		{NewIntDatum(1), `1.0`, true},
		{NewFloat64Datum(2), `2`, true},
		{NewStringDatum("\"hello, 世界\""), `"hello, 世界"`, true},
		{NewStringDatum("[1, 2, 3]"), `[1, 2, 3]`, true},
		{NewStringDatum("{}"), `{}`, true},
		{mustParseTimeIntoDatum("2011-11-10 11:11:11.111111", mysql.TypeTimestamp, 6), `"2011-11-10 11:11:11.111111"`, true},

		// can not parse JSON from this string, so error occurs.
		{NewStringDatum("hello, 世界"), "", false},
	}
	for _, tt := range tests {
		obtain, err := tt.datum.ConvertTo(sc, ft)
		if tt.success {
			require.NoError(t, err)

			sd := NewStringDatum(tt.expected)
			var expected Datum
			expected, err = sd.ConvertTo(sc, ft)
			require.NoError(t, err)

			var cmp int
			cmp, err = obtain.CompareDatum(sc, &expected)
			require.NoError(t, err)
			require.Equal(t, 0, cmp)
		} else {
			require.Error(t, err)
		}
	}
}

func TestIsNull(t *testing.T) {
	tests := []struct {
		data   interface{}
		isnull bool
	}{
		{nil, true},
		{0, false},
		{1, false},
		{1.1, false},
		{"string", false},
		{"", false},
	}
	for _, tt := range tests {
		testIsNull(t, tt.data, tt.isnull)
	}
}

func testIsNull(t *testing.T, data interface{}, isnull bool) {
	d := NewDatum(data)
	require.Equal(t, isnull, d.IsNull())
}

func TestToBytes(t *testing.T) {
	tests := []struct {
		a   Datum
		out []byte
	}{
		{NewIntDatum(1), []byte("1")},
		{NewDecimalDatum(NewDecFromInt(1)), []byte("1")},
		{NewFloat64Datum(1.23), []byte("1.23")},
		{NewStringDatum("abc"), []byte("abc")},
	}
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = true
	for _, tt := range tests {
		bin, err := tt.a.ToBytes()
		require.NoError(t, err)
		require.Equal(t, tt.out, bin)
	}
}

func mustParseDurationDatum(str string, fsp int) Datum {
	dur, err := ParseDuration(nil, str, fsp)
	if err != nil {
		panic(err)
	}
	return NewDurationDatum(dur)
}

func TestComputePlusAndMinus(t *testing.T) {
	sc := &stmtctx.StatementContext{TimeZone: time.UTC}
	tests := []struct {
		a      Datum
		b      Datum
		plus   Datum
		minus  Datum
		hasErr bool
	}{
		{NewIntDatum(72), NewIntDatum(28), NewIntDatum(100), NewIntDatum(44), false},
		{NewIntDatum(72), NewUintDatum(28), NewIntDatum(100), NewIntDatum(44), false},
		{NewUintDatum(72), NewUintDatum(28), NewUintDatum(100), NewUintDatum(44), false},
		{NewUintDatum(72), NewIntDatum(28), NewUintDatum(100), NewUintDatum(44), false},
		{NewFloat64Datum(72.0), NewFloat64Datum(28.0), NewFloat64Datum(100.0), NewFloat64Datum(44.0), false},
		{NewDecimalDatum(NewDecFromStringForTest("72.5")), NewDecimalDatum(NewDecFromInt(3)), NewDecimalDatum(NewDecFromStringForTest("75.5")), NewDecimalDatum(NewDecFromStringForTest("69.5")), false},
		{NewIntDatum(72), NewFloat64Datum(42), Datum{}, Datum{}, true},
		{NewStringDatum("abcd"), NewIntDatum(42), Datum{}, Datum{}, true},
	}

	for ith, tt := range tests {
		got, err := ComputePlus(tt.a, tt.b)
		require.Equal(t, tt.hasErr, err != nil)
		v, err := got.CompareDatum(sc, &tt.plus)
		require.NoError(t, err)
		require.Equalf(t, 0, v, "%dth got:%#v, %#v, expect:%#v, %#v", ith, got, got.x, tt.plus, tt.plus.x)
	}
}

func TestCopyDatum(t *testing.T) {
	var raw Datum
	raw.b = []byte("raw")
	raw.k = KindRaw
	tests := []Datum{
		NewIntDatum(72),
		NewUintDatum(72),
		NewStringDatum("abcd"),
		NewBytesDatum([]byte("abcd")),
		raw,
	}

	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = true
	for _, tt := range tests {
		tt1 := CopyDatum(tt)
		res, err := tt.CompareDatum(sc, &tt1)
		require.NoError(t, err)
		require.Equal(t, 0, res)
		if tt.b != nil {
			require.NotSame(t, &tt1.b[0], &tt.b[0])
		}
	}
}

func prepareCompareDatums() ([]Datum, []Datum) {
	vals := make([]Datum, 0, 5)
	vals = append(vals, NewIntDatum(1))
	vals = append(vals, NewFloat64Datum(1.23))
	vals = append(vals, NewStringDatum("abcde"))
	vals = append(vals, NewDecimalDatum(NewDecFromStringForTest("1.2345")))
	vals = append(vals, NewTimeDatum(Time{Time: FromGoTime(time.Date(2018, 3, 8, 16, 1, 0, 315313000, time.UTC)), Fsp: 6, Type: mysql.TypeTimestamp}))

	vals1 := make([]Datum, 0, 5)
	vals1 = append(vals1, NewIntDatum(1))
	vals1 = append(vals1, NewFloat64Datum(1.23))
	vals1 = append(vals1, NewStringDatum("abcde"))
	vals1 = append(vals1, NewDecimalDatum(NewDecFromStringForTest("1.2345")))
	vals1 = append(vals1, NewTimeDatum(Time{Time: FromGoTime(time.Date(2018, 3, 8, 16, 1, 0, 315313000, time.UTC)), Fsp: 6, Type: mysql.TypeTimestamp}))
	return vals, vals1
}

func BenchmarkCompareDatum(b *testing.B) {
	vals, vals1 := prepareCompareDatums()
	sc := new(stmtctx.StatementContext)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j, v := range vals {
			v.CompareDatum(sc, &vals1[j])
		}
	}
}

func BenchmarkCompareDatumByReflect(b *testing.B) {
	vals, vals1 := prepareCompareDatums()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reflect.DeepEqual(vals, vals1)
	}
}
