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
	"io"
	"testing"

	. "github.com/pingcap/check"
	"github.com/stretchr/testify/require"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/terror"
	"github.com/XiaoMi/Gaea/util/testleak"
)

func TestT(t *testing.T) {
	TestingT(t)
}

func testIsTypeBlob(t *testing.T, tp byte, expect bool) {
	v := IsTypeBlob(tp)
	require.Equal(t, expect, v)
}

func testIsTypeChar(t *testing.T, tp byte, expect bool) {
	v := IsTypeChar(tp)
	require.Equal(t, expect, v)
}

func TestIsType(t *testing.T) {
	defer testleak.AfterTestT(t)()
	testIsTypeBlob(t, mysql.TypeTinyBlob, true)
	testIsTypeBlob(t, mysql.TypeMediumBlob, true)
	testIsTypeBlob(t, mysql.TypeBlob, true)
	testIsTypeBlob(t, mysql.TypeLongBlob, true)
	testIsTypeBlob(t, mysql.TypeInt24, false)

	testIsTypeChar(t, mysql.TypeString, true)
	testIsTypeChar(t, mysql.TypeVarchar, true)
	testIsTypeChar(t, mysql.TypeLong, false)
}

func testTypeStr(t *testing.T, tp byte, expect string) {
	v := TypeStr(tp)
	require.Equal(t, expect, v)
}

func testTypeToStr(t *testing.T, tp byte, charset string, expect string) {
	v := TypeToStr(tp, charset)
	require.Equal(t, expect, v)
}

func TestTypeToStr(t *testing.T) {
	defer testleak.AfterTestT(t)()
	testTypeStr(t, mysql.TypeYear, "year")
	testTypeStr(t, 0xdd, "")

	testTypeToStr(t, mysql.TypeBlob, "utf8", "text")
	testTypeToStr(t, mysql.TypeLongBlob, "utf8", "longtext")
	testTypeToStr(t, mysql.TypeTinyBlob, "utf8", "tinytext")
	testTypeToStr(t, mysql.TypeMediumBlob, "utf8", "mediumtext")
	testTypeToStr(t, mysql.TypeVarchar, "binary", "varbinary")
	testTypeToStr(t, mysql.TypeString, "binary", "binary")
	testTypeToStr(t, mysql.TypeTiny, "binary", "tinyint")
	testTypeToStr(t, mysql.TypeBlob, "binary", "blob")
	testTypeToStr(t, mysql.TypeLongBlob, "binary", "longblob")
	testTypeToStr(t, mysql.TypeTinyBlob, "binary", "tinyblob")
	testTypeToStr(t, mysql.TypeMediumBlob, "binary", "mediumblob")
	testTypeToStr(t, mysql.TypeVarchar, "utf8", "varchar")
	testTypeToStr(t, mysql.TypeString, "utf8", "char")
	testTypeToStr(t, mysql.TypeShort, "binary", "smallint")
	testTypeToStr(t, mysql.TypeInt24, "binary", "mediumint")
	testTypeToStr(t, mysql.TypeLong, "binary", "int")
	testTypeToStr(t, mysql.TypeLonglong, "binary", "bigint")
	testTypeToStr(t, mysql.TypeFloat, "binary", "float")
	testTypeToStr(t, mysql.TypeDouble, "binary", "double")
	testTypeToStr(t, mysql.TypeYear, "binary", "year")
	testTypeToStr(t, mysql.TypeDuration, "binary", "time")
	testTypeToStr(t, mysql.TypeDatetime, "binary", "datetime")
	testTypeToStr(t, mysql.TypeDate, "binary", "date")
	testTypeToStr(t, mysql.TypeTimestamp, "binary", "timestamp")
	testTypeToStr(t, mysql.TypeNewDecimal, "binary", "decimal")
	testTypeToStr(t, mysql.TypeUnspecified, "binary", "unspecified")
	testTypeToStr(t, 0xdd, "binary", "")
	testTypeToStr(t, mysql.TypeBit, "binary", "bit")
	testTypeToStr(t, mysql.TypeEnum, "binary", "enum")
	testTypeToStr(t, mysql.TypeSet, "binary", "set")
}

func TestEOFAsNil(t *testing.T) {
	defer testleak.AfterTestT(t)()
	err := EOFAsNil(io.EOF)
	require.NoError(t, err)
}

func TestMaxFloat(t *testing.T) {
	defer testleak.AfterTestT(t)()
	tbl := []struct {
		Flen    int
		Decimal int
		Expect  float64
	}{
		{3, 2, 9.99},
		{5, 2, 999.99},
		{10, 1, 999999999.9},
		{5, 5, 0.99999},
	}

	for _, test := range tbl {
		f := GetMaxFloat(test.Flen, test.Decimal)
		require.Equal(t, test.Expect, f)
	}
}

func TestRoundFloat(t *testing.T) {
	defer testleak.AfterTestT(t)()
	tbl := []struct {
		Input  float64
		Expect float64
	}{
		{2.5, 3},
		{1.5, 2},
		{0.5, 1},
		{0.49999999999999997, 0},
		{0, 0},
		{-0.49999999999999997, 0},
		{-0.5, -1},
		{-2.5, -3},
		{-1.5, -2},
	}

	for _, test := range tbl {
		f := RoundFloat(test.Input)
		require.Equal(t, test.Expect, f)
	}
}

func TestRound(t *testing.T) {
	defer testleak.AfterTestT(t)()
	tbl := []struct {
		Input  float64
		Dec    int
		Expect float64
	}{
		{-1.23, 0, -1},
		{-1.58, 0, -2},
		{1.58, 0, 2},
		{1.298, 1, 1.3},
		{1.298, 0, 1},
		{23.298, -1, 20},
	}

	for _, test := range tbl {
		f := Round(test.Input, test.Dec)
		require.Equal(t, test.Expect, f)
	}
}

func TestTruncate(t *testing.T) {
	defer testleak.AfterTestT(t)()
	tbl := []struct {
		Input   float64
		Flen    int
		Decimal int
		Expect  float64
		Err     error
	}{
		{100.114, 10, 2, 100.11, nil},
		{100.115, 10, 2, 100.12, nil},
		{100.1156, 10, 3, 100.116, nil},
		{100.1156, 3, 1, 99.9, ErrOverflow},
		{1.36, 10, 2, 1.36, nil},
	}

	for _, test := range tbl {
		f, err := TruncateFloat(test.Input, test.Flen, test.Decimal)
		require.Equal(t, test.Expect, f)
		require.True(t, terror.ErrorEqual(err, test.Err))
	}
}
