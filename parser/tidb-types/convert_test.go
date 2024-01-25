// Copyright 2015 PingCAP, Inc.
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
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/pingcap/errors"
	"github.com/stretchr/testify/require"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/stmtctx"
	"github.com/XiaoMi/Gaea/parser/terror"
	"github.com/XiaoMi/Gaea/parser/tidb-types/json"
	"github.com/XiaoMi/Gaea/util/testleak"
)

type invalidMockType struct {
}

// Convert converts the val with type tp.
func Convert(val interface{}, target *FieldType) (v interface{}, err error) {
	d := NewDatum(val)
	sc := new(stmtctx.StatementContext)
	sc.TimeZone = time.UTC
	ret, err := d.ConvertTo(sc, target)
	if err != nil {
		return ret.GetValue(), errors.Trace(err)
	}
	return ret.GetValue(), nil
}

func TestConvertType(t *testing.T) {
	defer testleak.AfterTestT(t)()
	ft := NewFieldType(mysql.TypeBlob)
	ft.Flen = 4
	ft.Charset = "utf8"
	v, err := Convert("123456", ft)
	require.True(t, ErrDataTooLong.Equal(err))
	require.Equal(t, "1234", v)
	ft = NewFieldType(mysql.TypeString)
	ft.Flen = 4
	ft.Charset = mysql.CharsetBin
	v, err = Convert("12345", ft)
	require.True(t, ErrDataTooLong.Equal(err))
	require.Equal(t, []byte("1234"), v)

	ft = NewFieldType(mysql.TypeFloat)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(111.114, ft)
	require.NoError(t, err)
	require.Equal(t, float32(111.11), v)

	ft = NewFieldType(mysql.TypeFloat)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(999.999, ft)
	require.Error(t, err)
	require.Equal(t, float32(999.99), v)

	ft = NewFieldType(mysql.TypeFloat)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(-999.999, ft)
	require.Error(t, err)
	require.Equal(t, float32(-999.99), v)

	ft = NewFieldType(mysql.TypeFloat)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(1111.11, ft)
	require.Error(t, err)
	require.Equal(t, float32(999.99), v)

	ft = NewFieldType(mysql.TypeFloat)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(999.916, ft)
	require.NoError(t, err)
	require.Equal(t, float32(999.92), v)

	ft = NewFieldType(mysql.TypeFloat)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(999.914, ft)
	require.NoError(t, err)
	require.Equal(t, float32(999.91), v)

	ft = NewFieldType(mysql.TypeFloat)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(999.9155, ft)
	require.NoError(t, err)
	require.Equal(t, float32(999.92), v)

	// For TypeBlob
	ft = NewFieldType(mysql.TypeBlob)
	_, err = Convert(&invalidMockType{}, ft)
	require.Error(t, err)

	// Nil
	ft = NewFieldType(mysql.TypeBlob)
	v, err = Convert(nil, ft)
	require.NoError(t, err)
	require.Nil(t, v)

	// TypeDouble
	ft = NewFieldType(mysql.TypeDouble)
	ft.Flen = 5
	ft.Decimal = 2
	v, err = Convert(999.9155, ft)
	require.NoError(t, err)
	require.Equal(t, float64(999.92), v)

	// For TypeString
	ft = NewFieldType(mysql.TypeString)
	ft.Flen = 3
	v, err = Convert("12345", ft)
	require.True(t, ErrDataTooLong.Equal(err))
	require.Equal(t, "123", v)
	ft = NewFieldType(mysql.TypeString)
	ft.Flen = 3
	ft.Charset = mysql.CharsetBin
	v, err = Convert("12345", ft)
	require.True(t, ErrDataTooLong.Equal(err))
	require.Equal(t, []byte("123"), v)

	// For TypeDuration
	ft = NewFieldType(mysql.TypeDuration)
	ft.Decimal = 3
	v, err = Convert("10:11:12.123456", ft)
	require.NoError(t, err)
	require.Equal(t, "10:11:12.123", v.(Duration).String())

	ft.Decimal = 1
	vv, err := Convert(v, ft)
	require.NoError(t, err)
	require.Equal(t, "10:11:12.1", vv.(Duration).String())

	vd, err := ParseTime(nil, "2010-10-10 10:11:11.12345", mysql.TypeDatetime, 2)
	require.Equal(t, "2010-10-10 10:11:11.12", vd.String())
	require.NoError(t, err)
	v, err = Convert(vd, ft)
	require.NoError(t, err)
	require.Equal(t, "10:11:11.1", v.(Duration).String())

	vt, err := ParseTime(&stmtctx.StatementContext{TimeZone: time.UTC}, "2010-10-10 10:11:11.12345", mysql.TypeTimestamp, 2)
	require.Equal(t, "2010-10-10 10:11:11.12", vt.String())
	require.NoError(t, err)
	v, err = Convert(vt, ft)
	require.NoError(t, err)
	require.Equal(t, "10:11:11.1", v.(Duration).String())

	// For mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate
	ft = NewFieldType(mysql.TypeTimestamp)
	ft.Decimal = 3
	v, err = Convert("2010-10-10 10:11:11.12345", ft)
	require.NoError(t, err)
	require.Equal(t, "2010-10-10 10:11:11.123", v.(Time).String())

	ft.Decimal = 1
	vv, err = Convert(v, ft)
	require.NoError(t, err)
	require.Equal(t, "2010-10-10 10:11:11.1", vv.(Time).String())

	// For TypeLonglong
	ft = NewFieldType(mysql.TypeLonglong)
	v, err = Convert("100", ft)
	require.NoError(t, err)
	require.Equal(t, int64(100), v)
	// issue 4287.
	v, err = Convert(math.Pow(2, 63)-1, ft)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt64), v)
	ft = NewFieldType(mysql.TypeLonglong)
	ft.Flag |= mysql.UnsignedFlag
	v, err = Convert("100", ft)
	require.NoError(t, err)
	require.Equal(t, uint64(100), v)
	// issue 3470
	ft = NewFieldType(mysql.TypeLonglong)
	v, err = Convert(Duration{Duration: time.Duration(12*time.Hour + 59*time.Minute + 59*time.Second + 555*time.Millisecond), Fsp: 3}, ft)
	require.NoError(t, err)
	require.Equal(t, int64(130000), v)
	v, err = Convert(Time{
		Time: FromDate(2017, 1, 1, 12, 59, 59, 555000),
		Type: mysql.TypeDatetime,
		Fsp:  MaxFsp}, ft)
	require.NoError(t, err)
	require.Equal(t, int64(20170101130000), v)

	// For TypeBit
	ft = NewFieldType(mysql.TypeBit)
	ft.Flen = 24 // 3 bytes.
	v, err = Convert("100", ft)
	require.NoError(t, err)
	require.Equal(t, NewBinaryLiteralFromUint(3223600, 3), v)

	v, err = Convert(NewBinaryLiteralFromUint(100, -1), ft)
	require.NoError(t, err)
	require.Equal(t, NewBinaryLiteralFromUint(100, 3), v)

	ft.Flen = 1
	v, err = Convert(1, ft)
	require.NoError(t, err)
	require.Equal(t, NewBinaryLiteralFromUint(1, 1), v)

	_, err = Convert(2, ft)
	require.Error(t, err)

	ft.Flen = 0
	_, err = Convert(2, ft)
	require.Error(t, err)

	// For TypeNewDecimal
	ft = NewFieldType(mysql.TypeNewDecimal)
	ft.Flen = 8
	ft.Decimal = 4
	v, err = Convert(3.1416, ft)
	require.NoErrorf(t, err, errors.ErrorStack(err))
	require.Equal(t, "3.1416", v.(*MyDecimal).String())
	v, err = Convert("3.1415926", ft)
	require.True(t, terror.ErrorEqual(err, ErrTruncated))
	require.Equal(t, "3.1416", v.(*MyDecimal).String())
	v, err = Convert("99999", ft)
	require.Truef(t, terror.ErrorEqual(err, ErrOverflow), "err %v", err)
	require.Equal(t, "9999.9999", v.(*MyDecimal).String())

	v, err = Convert("-10000", ft)
	require.Truef(t, terror.ErrorEqual(err, ErrOverflow), "err %v", err)
	require.Equal(t, "-9999.9999", v.(*MyDecimal).String())

	// Test Datum.ToDecimal with bad number.
	d := NewDatum("hello")
	sc := new(stmtctx.StatementContext)
	_, err = d.ToDecimal(sc)
	require.Truef(t, terror.ErrorEqual(err, ErrBadNumber), "err %v", err)

	sc.IgnoreTruncate = true
	v, err = d.ToDecimal(sc)
	require.NoError(t, err)
	require.Equal(t, "0", v.(*MyDecimal).String())

	// For TypeYear
	ft = NewFieldType(mysql.TypeYear)
	v, err = Convert("2015", ft)
	require.NoError(t, err)
	require.Equal(t, int64(2015), v)
	v, err = Convert(2015, ft)
	require.NoError(t, err)
	require.Equal(t, int64(2015), v)
	_, err = Convert(1800, ft)
	require.Error(t, err)
	dt, err := ParseDate(nil, "2015-11-11")
	require.NoError(t, err)
	v, err = Convert(dt, ft)
	require.NoError(t, err)
	require.Equal(t, int64(2015), v)
	v, err = Convert(ZeroDuration, ft)
	require.NoError(t, err)
	require.Equal(t, int64(time.Now().Year()), v)

	// For enum
	ft = NewFieldType(mysql.TypeEnum)
	ft.Elems = []string{"a", "b", "c"}
	v, err = Convert("a", ft)
	require.NoError(t, err)
	require.Equal(t, Enum{Name: "a", Value: 1}, v)
	v, err = Convert(2, ft)
	require.NoError(t, err)
	require.Equal(t, Enum{Name: "b", Value: 2}, v)
	_, err = Convert("d", ft)
	require.Error(t, err)
	v, err = Convert(4, ft)
	require.Truef(t, terror.ErrorEqual(err, ErrTruncated), "err %v", err)
	require.Equal(t, Enum{}, v)

	ft = NewFieldType(mysql.TypeSet)
	ft.Elems = []string{"a", "b", "c"}
	v, err = Convert("a", ft)
	require.NoError(t, err)
	require.Equal(t, Set{Name: "a", Value: 1}, v)
	v, err = Convert(2, ft)
	require.NoError(t, err)
	require.Equal(t, Set{Name: "b", Value: 2}, v)
	v, err = Convert(3, ft)
	require.NoError(t, err)
	require.Equal(t, Set{Name: "a,b", Value: 3}, v)
	_, err = Convert("d", ft)
	require.Error(t, err)
	_, err = Convert(9, ft)
	require.Error(t, err)
}

func testToString(t *testing.T, val interface{}, expect string) {
	b, err := ToString(val)
	require.NoError(t, err)
	require.Equal(t, expect, b)
}

func TestConvertToString(t *testing.T) {
	defer testleak.AfterTestT(t)()
	testToString(t, "0", "0")
	testToString(t, true, "1")
	testToString(t, "false", "false")
	testToString(t, int(0), "0")
	testToString(t, int64(0), "0")
	testToString(t, uint64(0), "0")
	testToString(t, float32(1.6), "1.6")
	testToString(t, float64(-0.6), "-0.6")
	testToString(t, []byte{1}, "\x01")
	testToString(t, NewBinaryLiteralFromUint(0x4D7953514C, -1), "MySQL")
	testToString(t, NewBinaryLiteralFromUint(0x41, -1), "A")
	testToString(t, Enum{Name: "a", Value: 1}, "a")
	testToString(t, Set{Name: "a", Value: 1}, "a")

	t1, err := ParseTime(&stmtctx.StatementContext{TimeZone: time.UTC},
		"2011-11-10 11:11:11.999999", mysql.TypeTimestamp, 6)
	require.NoError(t, err)
	testToString(t, t1, "2011-11-10 11:11:11.999999")

	td, err := ParseDuration(nil, "11:11:11.999999", 6)
	require.NoError(t, err)
	testToString(t, td, "11:11:11.999999")

	ft := NewFieldType(mysql.TypeNewDecimal)
	ft.Flen = 10
	ft.Decimal = 5
	v, err := Convert(3.1415926, ft)
	require.True(t, terror.ErrorEqual(err, ErrTruncated))
	testToString(t, v, "3.14159")

	_, err = ToString(&invalidMockType{})
	require.Error(t, err)

	// test truncate
	tests := []struct {
		flen    int
		charset string
		input   string
		output  string
	}{
		{5, "utf8", "你好，世界", "你好，世界"},
		{5, "utf8mb4", "你好，世界", "你好，世界"},
		{4, "utf8", "你好，世界", "你好，世"},
		{4, "utf8mb4", "你好，世界", "你好，世"},
		{15, "binary", "你好，世界", "你好，世界"},
		{12, "binary", "你好，世界", "你好，世"},
		{0, "binary", "你好，世界", ""},
	}
	for _, tt := range tests {
		ft = NewFieldType(mysql.TypeVarchar)
		ft.Flen = tt.flen
		ft.Charset = tt.charset
		inputDatum := NewStringDatum(tt.input)
		sc := new(stmtctx.StatementContext)
		outputDatum, err := inputDatum.ConvertTo(sc, ft)
		if tt.input != tt.output {
			require.True(t, ErrDataTooLong.Equal(err))
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tt.output, outputDatum.GetString())
	}
}

func testStrToInt(t *testing.T, str string, expect int64, truncateAsErr bool, expectErr error) {
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = !truncateAsErr
	val, err := StrToInt(sc, str)
	if expectErr != nil {
		require.Truef(t, terror.ErrorEqual(err, expectErr), "err %v", err)
	} else {
		require.NoError(t, err)
		require.Equal(t, expect, val)
	}
}

func testStrToUint(t *testing.T, str string, expect uint64, truncateAsErr bool, expectErr error) {
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = !truncateAsErr
	val, err := StrToUint(sc, str)
	if expectErr != nil {
		require.Truef(t, terror.ErrorEqual(err, expectErr), "err %v", err)
	} else {
		require.NoError(t, err)
		require.Equal(t, expect, val)
	}
}

func testStrToFloat(t *testing.T, str string, expect float64, truncateAsErr bool, expectErr error) {
	sc := new(stmtctx.StatementContext)
	sc.IgnoreTruncate = !truncateAsErr
	val, err := StrToFloat(sc, str)
	if expectErr != nil {
		require.Truef(t, terror.ErrorEqual(err, expectErr), "err %v", err)
	} else {
		require.NoError(t, err)
		require.Equal(t, expect, val)
	}
}

func TestStrToNum(t *testing.T) {
	defer testleak.AfterTestT(t)()
	testStrToInt(t, "0", 0, true, nil)
	testStrToInt(t, "-1", -1, true, nil)
	testStrToInt(t, "100", 100, true, nil)
	testStrToInt(t, "65.0", 65, false, nil)
	testStrToInt(t, "65.0", 65, true, nil)
	testStrToInt(t, "", 0, false, nil)
	testStrToInt(t, "", 0, true, ErrTruncated)
	testStrToInt(t, "xx", 0, true, ErrTruncated)
	testStrToInt(t, "xx", 0, false, nil)
	testStrToInt(t, "11xx", 11, true, ErrTruncated)
	testStrToInt(t, "11xx", 11, false, nil)
	testStrToInt(t, "xx11", 0, false, nil)

	testStrToUint(t, "0", 0, true, nil)
	testStrToUint(t, "", 0, false, nil)
	testStrToUint(t, "", 0, false, nil)
	testStrToUint(t, "-1", 0xffffffffffffffff, false, ErrOverflow)
	testStrToUint(t, "100", 100, true, nil)
	testStrToUint(t, "+100", 100, true, nil)
	testStrToUint(t, "65.0", 65, true, nil)
	testStrToUint(t, "xx", 0, true, ErrTruncated)
	testStrToUint(t, "11xx", 11, true, ErrTruncated)
	testStrToUint(t, "xx11", 0, true, ErrTruncated)

	// TODO: makes StrToFloat return truncated value instead of zero to make it pass.
	testStrToFloat(t, "", 0, true, ErrTruncated)
	testStrToFloat(t, "-1", -1.0, true, nil)
	testStrToFloat(t, "1.11", 1.11, true, nil)
	testStrToFloat(t, "1.11.00", 1.11, false, nil)
	testStrToFloat(t, "1.11.00", 1.11, true, ErrTruncated)
	testStrToFloat(t, "xx", 0.0, false, nil)
	testStrToFloat(t, "0x00", 0.0, false, nil)
	testStrToFloat(t, "11.xx", 11.0, false, nil)
	testStrToFloat(t, "11.xx", 11.0, true, ErrTruncated)
	testStrToFloat(t, "xx.11", 0.0, false, nil)

	// for issue #5111
	testStrToFloat(t, "1e649", math.MaxFloat64, true, ErrTruncatedWrongVal)
	testStrToFloat(t, "1e649", math.MaxFloat64, false, nil)
	testStrToFloat(t, "-1e649", -math.MaxFloat64, true, ErrTruncatedWrongVal)
	testStrToFloat(t, "-1e649", -math.MaxFloat64, false, nil)
}

func TestFieldTypeToStr(t *testing.T) {
	defer testleak.AfterTestT(t)()
	v := TypeToStr(mysql.TypeUnspecified, "not binary")
	require.Equal(t, TypeStr(mysql.TypeUnspecified), v)
	v = TypeToStr(mysql.TypeBlob, mysql.CharsetBin)
	require.Equal(t, "blob", v)
	v = TypeToStr(mysql.TypeString, mysql.CharsetBin)
	require.Equal(t, "binary", v)
}

func accept(t *testing.T, tp byte, value interface{}, unsigned bool, expected string) {
	ft := NewFieldType(tp)
	if unsigned {
		ft.Flag |= mysql.UnsignedFlag
	}
	d := NewDatum(value)
	sc := new(stmtctx.StatementContext)
	sc.TimeZone = time.UTC
	sc.IgnoreTruncate = true
	casted, err := d.ConvertTo(sc, ft)
	require.NoErrorf(t, err, "%v", ft)
	if casted.IsNull() {
		require.Equal(t, "<nil>", expected)
	} else {
		str, err := casted.ToString()
		require.NoError(t, err)
		require.Equal(t, expected, str)
	}
}

func unsignedAccept(t *testing.T, tp byte, value interface{}, expected string) {
	accept(t, tp, value, true, expected)
}

func signedAccept(t *testing.T, tp byte, value interface{}, expected string) {
	accept(t, tp, value, false, expected)
}

func deny(t *testing.T, tp byte, value interface{}, unsigned bool, expected string) {
	ft := NewFieldType(tp)
	if unsigned {
		ft.Flag |= mysql.UnsignedFlag
	}
	d := NewDatum(value)
	sc := new(stmtctx.StatementContext)
	casted, err := d.ConvertTo(sc, ft)
	require.Error(t, err)
	if casted.IsNull() {
		require.Equal(t, "<nil>", expected)
	} else {
		str, err := casted.ToString()
		require.NoError(t, err)
		require.Equal(t, expected, str)
	}
}

func unsignedDeny(t *testing.T, tp byte, value interface{}, expected string) {
	deny(t, tp, value, true, expected)
}

func signedDeny(t *testing.T, tp byte, value interface{}, expected string) {
	deny(t, tp, value, false, expected)
}

func strvalue(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func TestConvert(t *testing.T) {
	defer testleak.AfterTestT(t)()
	// integer ranges
	signedDeny(t, mysql.TypeTiny, -129, "-128")
	signedAccept(t, mysql.TypeTiny, -128, "-128")
	signedAccept(t, mysql.TypeTiny, 127, "127")
	signedDeny(t, mysql.TypeTiny, 128, "127")
	unsignedDeny(t, mysql.TypeTiny, -1, "255")
	unsignedAccept(t, mysql.TypeTiny, 0, "0")
	unsignedAccept(t, mysql.TypeTiny, 255, "255")
	unsignedDeny(t, mysql.TypeTiny, 256, "255")

	signedDeny(t, mysql.TypeShort, int64(math.MinInt16)-1, strvalue(int64(math.MinInt16)))
	signedAccept(t, mysql.TypeShort, int64(math.MinInt16), strvalue(int64(math.MinInt16)))
	signedAccept(t, mysql.TypeShort, int64(math.MaxInt16), strvalue(int64(math.MaxInt16)))
	signedDeny(t, mysql.TypeShort, int64(math.MaxInt16)+1, strvalue(int64(math.MaxInt16)))
	unsignedDeny(t, mysql.TypeShort, -1, "65535")
	unsignedAccept(t, mysql.TypeShort, 0, "0")
	unsignedAccept(t, mysql.TypeShort, uint64(math.MaxUint16), strvalue(uint64(math.MaxUint16)))
	unsignedDeny(t, mysql.TypeShort, uint64(math.MaxUint16)+1, strvalue(uint64(math.MaxUint16)))

	signedDeny(t, mysql.TypeInt24, -1<<23-1, strvalue(-1<<23))
	signedAccept(t, mysql.TypeInt24, -1<<23, strvalue(-1<<23))
	signedAccept(t, mysql.TypeInt24, 1<<23-1, strvalue(1<<23-1))
	signedDeny(t, mysql.TypeInt24, 1<<23, strvalue(1<<23-1))
	unsignedDeny(t, mysql.TypeInt24, -1, "16777215")
	unsignedAccept(t, mysql.TypeInt24, 0, "0")
	unsignedAccept(t, mysql.TypeInt24, 1<<24-1, strvalue(1<<24-1))
	unsignedDeny(t, mysql.TypeInt24, 1<<24, strvalue(1<<24-1))

	signedDeny(t, mysql.TypeLong, int64(math.MinInt32)-1, strvalue(int64(math.MinInt32)))
	signedAccept(t, mysql.TypeLong, int64(math.MinInt32), strvalue(int64(math.MinInt32)))
	signedAccept(t, mysql.TypeLong, int64(math.MaxInt32), strvalue(int64(math.MaxInt32)))
	signedDeny(t, mysql.TypeLong, uint64(math.MaxUint64), strvalue(uint64(math.MaxInt32)))
	signedDeny(t, mysql.TypeLong, int64(math.MaxInt32)+1, strvalue(int64(math.MaxInt32)))
	signedDeny(t, mysql.TypeLong, "1343545435346432587475", strvalue(int64(math.MaxInt32)))
	unsignedDeny(t, mysql.TypeLong, -1, "4294967295")
	unsignedAccept(t, mysql.TypeLong, 0, "0")
	unsignedAccept(t, mysql.TypeLong, uint64(math.MaxUint32), strvalue(uint64(math.MaxUint32)))
	unsignedDeny(t, mysql.TypeLong, uint64(math.MaxUint32)+1, strvalue(uint64(math.MaxUint32)))

	signedDeny(t, mysql.TypeLonglong, math.MinInt64*1.1, strvalue(int64(math.MinInt64)))
	signedAccept(t, mysql.TypeLonglong, int64(math.MinInt64), strvalue(int64(math.MinInt64)))
	signedAccept(t, mysql.TypeLonglong, int64(math.MaxInt64), strvalue(int64(math.MaxInt64)))
	signedDeny(t, mysql.TypeLonglong, math.MaxInt64*1.1, strvalue(int64(math.MaxInt64)))
	unsignedAccept(t, mysql.TypeLonglong, -1, "18446744073709551615")
	unsignedAccept(t, mysql.TypeLonglong, 0, "0")
	unsignedAccept(t, mysql.TypeLonglong, uint64(math.MaxUint64), strvalue(uint64(math.MaxUint64)))
	unsignedDeny(t, mysql.TypeLonglong, math.MaxUint64*1.1, strvalue(uint64(math.MaxUint64)))

	// integer from string
	signedAccept(t, mysql.TypeLong, "	  234  ", "234")
	signedAccept(t, mysql.TypeLong, " 2.35e3  ", "2350")
	signedAccept(t, mysql.TypeLong, " 2.e3  ", "2000")
	signedAccept(t, mysql.TypeLong, " -2.e3  ", "-2000")
	signedAccept(t, mysql.TypeLong, " 2e2  ", "200")
	signedAccept(t, mysql.TypeLong, " 0.002e3  ", "2")
	signedAccept(t, mysql.TypeLong, " .002e3  ", "2")
	signedAccept(t, mysql.TypeLong, " 20e-2  ", "0")
	signedAccept(t, mysql.TypeLong, " -20e-2  ", "0")
	signedAccept(t, mysql.TypeLong, " +2.51 ", "3")
	signedAccept(t, mysql.TypeLong, " -9999.5 ", "-10000")
	signedAccept(t, mysql.TypeLong, " 999.4", "999")
	signedAccept(t, mysql.TypeLong, " -3.58", "-4")
	signedDeny(t, mysql.TypeLong, " 1a ", "1")
	signedDeny(t, mysql.TypeLong, " +1+ ", "1")

	// integer from float
	signedAccept(t, mysql.TypeLong, 234.5456, "235")
	signedAccept(t, mysql.TypeLong, -23.45, "-23")
	unsignedAccept(t, mysql.TypeLonglong, 234.5456, "235")
	unsignedDeny(t, mysql.TypeLonglong, -23.45, "18446744073709551593")

	// float from string
	signedAccept(t, mysql.TypeFloat, "23.523", "23.523")
	signedAccept(t, mysql.TypeFloat, int64(123), "123")
	signedAccept(t, mysql.TypeFloat, uint64(123), "123")
	signedAccept(t, mysql.TypeFloat, int(123), "123")
	signedAccept(t, mysql.TypeFloat, float32(123), "123")
	signedAccept(t, mysql.TypeFloat, float64(123), "123")
	signedAccept(t, mysql.TypeDouble, " -23.54", "-23.54")
	signedDeny(t, mysql.TypeDouble, "-23.54a", "-23.54")
	signedDeny(t, mysql.TypeDouble, "-23.54e2e", "-2354")
	signedDeny(t, mysql.TypeDouble, "+.e", "0")
	signedAccept(t, mysql.TypeDouble, "1e+1", "10")

	// year
	signedDeny(t, mysql.TypeYear, 123, "<nil>")
	signedDeny(t, mysql.TypeYear, 3000, "<nil>")
	signedAccept(t, mysql.TypeYear, "2000", "2000")

	// time from string
	signedAccept(t, mysql.TypeDate, "2012-08-23", "2012-08-23")
	signedAccept(t, mysql.TypeDatetime, "2012-08-23 12:34:03.123456", "2012-08-23 12:34:03")
	signedAccept(t, mysql.TypeDatetime, ZeroDatetime, "0000-00-00 00:00:00")
	signedAccept(t, mysql.TypeDatetime, int64(0), "0000-00-00 00:00:00")
	signedAccept(t, mysql.TypeTimestamp, "2012-08-23 12:34:03.123456", "2012-08-23 12:34:03")
	signedAccept(t, mysql.TypeDuration, "10:11:12", "10:11:12")
	signedAccept(t, mysql.TypeDuration, ZeroDatetime, "00:00:00")
	signedAccept(t, mysql.TypeDuration, ZeroDuration, "00:00:00")
	signedAccept(t, mysql.TypeDuration, 0, "00:00:00")

	signedDeny(t, mysql.TypeDate, "2012-08-x", "0000-00-00")
	signedDeny(t, mysql.TypeDatetime, "2012-08-x", "0000-00-00 00:00:00")
	signedDeny(t, mysql.TypeTimestamp, "2012-08-x", "0000-00-00 00:00:00")
	signedDeny(t, mysql.TypeDuration, "2012-08-x", "00:00:00")

	// string from string
	signedAccept(t, mysql.TypeString, "abc", "abc")

	// string from integer
	signedAccept(t, mysql.TypeString, 5678, "5678")
	signedAccept(t, mysql.TypeString, ZeroDuration, "00:00:00")
	signedAccept(t, mysql.TypeString, ZeroDatetime, "0000-00-00 00:00:00")
	signedAccept(t, mysql.TypeString, []byte("123"), "123")

	//TODO add more tests
	signedAccept(t, mysql.TypeNewDecimal, 123, "123")
	signedAccept(t, mysql.TypeNewDecimal, int64(123), "123")
	signedAccept(t, mysql.TypeNewDecimal, uint64(123), "123")
	signedAccept(t, mysql.TypeNewDecimal, float32(123), "123")
	signedAccept(t, mysql.TypeNewDecimal, 123.456, "123.456")
	signedAccept(t, mysql.TypeNewDecimal, "-123.456", "-123.456")
	signedAccept(t, mysql.TypeNewDecimal, NewDecFromInt(12300000), "12300000")
	dec := NewDecFromInt(-123)
	dec.Shift(-5)
	dec.Round(dec, 5, ModeHalfEven)
	signedAccept(t, mysql.TypeNewDecimal, dec, "-0.00123")
}

func TestGetValidFloat(t *testing.T) {
	tests := []struct {
		origin string
		valid  string
	}{
		{"-100", "-100"},
		{"1abc", "1"},
		{"-1-1", "-1"},
		{"+1+1", "+1"},
		{"123..34", "123."},
		{"123.23E-10", "123.23E-10"},
		{"1.1e1.3", "1.1e1"},
		{"11e1.3", "11e1"},
		{"1.1e-13a", "1.1e-13"},
		{"1.", "1."},
		{".1", ".1"},
		{"", "0"},
		{"123e+", "123"},
		{"123.e", "123."},
	}
	sc := new(stmtctx.StatementContext)
	for _, tt := range tests {
		prefix, _ := getValidFloatPrefix(sc, tt.origin)
		require.Equal(t, tt.valid, prefix)
		_, err := strconv.ParseFloat(prefix, 64)
		require.NoError(t, err)
	}
	floatStr, err := floatStrToIntStr(sc, "1e9223372036854775807", "1e9223372036854775807")
	require.NoError(t, err)
	require.Equal(t, "1", floatStr)
	floatStr, err = floatStrToIntStr(sc, "125e342", "125e342.83")
	require.NoError(t, err)
	require.Equal(t, "125", floatStr)
	floatStr, err = floatStrToIntStr(sc, "1e21", "1e21")
	require.NoError(t, err)
	require.Equal(t, "1", floatStr)
}

// TestConvertTime tests time related conversion.
// time conversion is complicated including Date/Datetime/Time/Timestamp etc,
// Timestamp may involving timezone.
func TestConvertTime(t *testing.T) {
	timezones := []*time.Location{
		time.UTC,
		time.FixedZone("", 3*3600),
		time.Local,
	}

	for _, timezone := range timezones {
		sc := &stmtctx.StatementContext{
			TimeZone: timezone,
		}
		testConvertTimeTimeZone(t, sc)
	}
}

func testConvertTimeTimeZone(t *testing.T, sc *stmtctx.StatementContext) {
	raw := FromDate(2002, 3, 4, 4, 6, 7, 8)
	tests := []struct {
		input  Time
		target *FieldType
		expect Time
	}{
		{
			input:  Time{Type: mysql.TypeDatetime, Time: raw},
			target: NewFieldType(mysql.TypeTimestamp),
			expect: Time{Type: mysql.TypeTimestamp, Time: raw},
		},
		{
			input:  Time{Type: mysql.TypeDatetime, Time: raw},
			target: NewFieldType(mysql.TypeTimestamp),
			expect: Time{Type: mysql.TypeTimestamp, Time: raw},
		},
		{
			input:  Time{Type: mysql.TypeDatetime, Time: raw},
			target: NewFieldType(mysql.TypeTimestamp),
			expect: Time{Type: mysql.TypeTimestamp, Time: raw},
		},
		{
			input:  Time{Type: mysql.TypeTimestamp, Time: raw},
			target: NewFieldType(mysql.TypeDatetime),
			expect: Time{Type: mysql.TypeDatetime, Time: raw},
		},
	}

	for _, test := range tests {
		var d Datum
		d.SetMysqlTime(test.input)
		nd, err := d.ConvertTo(sc, test.target)
		require.NoError(t, err)
		v := nd.GetMysqlTime()
		require.Equal(t, test.expect.Type, v.Type)
		require.Equal(t, test.expect.Time, v.Time)
	}
}

func TestConvertJSONToInt(t *testing.T) {
	var tests = []struct {
		In  string
		Out int64
	}{
		{`{}`, 0},
		{`[]`, 0},
		{`3`, 3},
		{`-3`, -3},
		{`4.5`, 5},
		{`true`, 1},
		{`false`, 0},
		{`null`, 0},
		{`"hello"`, 0},
		{`"123hello"`, 123},
		{`"1234"`, 1234},
	}
	for _, tt := range tests {
		j, err := json.ParseBinaryFromString(tt.In)
		require.NoError(t, err)
		casted, _ := ConvertJSONToInt(new(stmtctx.StatementContext), j, false)
		require.Equal(t, tt.Out, casted)
	}
}

func TestConvertJSONToFloat(t *testing.T) {
	var tests = []struct {
		In  string
		Out float64
	}{
		{`{}`, 0},
		{`[]`, 0},
		{`3`, 3},
		{`-3`, -3},
		{`4.5`, 4.5},
		{`true`, 1},
		{`false`, 0},
		{`null`, 0},
		{`"hello"`, 0},
		{`"123.456hello"`, 123.456},
		{`"1234"`, 1234},
	}
	for _, tt := range tests {
		j, err := json.ParseBinaryFromString(tt.In)
		require.NoError(t, err)
		casted, _ := ConvertJSONToFloat(new(stmtctx.StatementContext), j)
		require.Equal(t, tt.Out, casted)
	}
}

func TestConvertJSONToDecimal(t *testing.T) {
	var tests = []struct {
		In  string
		Out *MyDecimal
	}{
		{`{}`, NewDecFromStringForTest("0")},
		{`[]`, NewDecFromStringForTest("0")},
		{`3`, NewDecFromStringForTest("3")},
		{`-3`, NewDecFromStringForTest("-3")},
		{`4.5`, NewDecFromStringForTest("4.5")},
		{`"1234"`, NewDecFromStringForTest("1234")},
		{`"1234567890123456789012345678901234567890123456789012345"`, NewDecFromStringForTest("1234567890123456789012345678901234567890123456789012345")},
	}
	for _, tt := range tests {
		j, err := json.ParseBinaryFromString(tt.In)
		require.NoError(t, err)
		casted, _ := ConvertJSONToDecimal(new(stmtctx.StatementContext), j)
		require.Equal(t, 0, casted.Compare(tt.Out))
	}
}

func TestNumberToDuration(t *testing.T) {
	var testCases = []struct {
		number int64
		fsp    int
		hasErr bool
		year   int
		month  int
		day    int
		hour   int
		minute int
		second int
	}{
		{20171222, 0, true, 0, 0, 0, 0, 0, 0},
		{171222, 0, false, 0, 0, 0, 17, 12, 22},
		{20171222020005, 0, false, 2017, 12, 22, 02, 00, 05},
		{10000000000, 0, true, 0, 0, 0, 0, 0, 0},
		{171222, 1, false, 0, 0, 0, 17, 12, 22},
		{176022, 1, true, 0, 0, 0, 0, 0, 0},
		{8391222, 1, true, 0, 0, 0, 0, 0, 0},
		{8381222, 0, false, 0, 0, 0, 838, 12, 22},
		{1001222, 0, false, 0, 0, 0, 100, 12, 22},
		{171260, 1, true, 0, 0, 0, 0, 0, 0},
	}

	for _, tc := range testCases {
		dur, err := NumberToDuration(tc.number, tc.fsp)
		if tc.hasErr {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, tc.hour, dur.Hour())
		require.Equal(t, tc.minute, dur.Minute())
		require.Equal(t, tc.second, dur.Second())
	}

	var testCases1 = []struct {
		number int64
		dur    time.Duration
	}{
		{171222, 17*time.Hour + 12*time.Minute + 22*time.Second},
		{-171222, -(17*time.Hour + 12*time.Minute + 22*time.Second)},
	}

	for _, tc := range testCases1 {
		dur, err := NumberToDuration(tc.number, 0)
		require.NoError(t, err)
		require.Equal(t, tc.dur, dur.Duration)
	}
}
