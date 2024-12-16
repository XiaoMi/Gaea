// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestLimitSelectResult(t *testing.T) {
	tests := []struct {
		offset       int64
		count        int64
		retLen       int64
		trimedRetLen int64
	}{
		{-1, -1, 5, 5},
		{0, -1, 5, 5},
		{0, 3, 5, 3},
		{0, 10, 5, 5},
		{3, 10, 5, 2},
		{3, 1, 5, 1},
		{4, 1, 5, 1},
		{5, 1, 5, 0},
		{5, 10, 5, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d:%d", test.offset, test.count), func(t *testing.T) {
			info := &SelectPlan{
				offset: test.offset,
				count:  test.count,
			}

			ret := &mysql.Result{
				Resultset: &mysql.Resultset{
					Values:   make([][]interface{}, test.retLen),
					RowDatas: make([]mysql.RowData, test.retLen),
				},
			}

			if err := limitSelectResult(info, ret); err != nil {
				t.Fatalf("limitSelectResult error: %v", err)
			}

			if int64(len(ret.Values)) != test.trimedRetLen {
				t.Errorf("len Values not equal, expect: %d, actual: %d", test.trimedRetLen, len(ret.Values))
			}
		})
	}
}

func TestResultRow_GetInt(t *testing.T) {
	tests := []struct {
		name     string
		row      ResultRow
		column   int
		expected int64
		hasError bool
	}{
		{"Valid int64", ResultRow{int64(10)}, 0, 10, false},
		{"Valid uint64", ResultRow{uint64(20)}, 0, 20, false},
		{"Valid float64", ResultRow{float64(30.5)}, 0, 30, false},
		{"Valid string", ResultRow{"40"}, 0, 40, false},
		{"Valid nil", ResultRow{nil}, 0, 0, false},
		{"Valid decimal", ResultRow{decimal.NewFromFloat(123.45)}, 0, 123, false},
		{"Valid small decimal", ResultRow{decimal.New(1, -10)}, 0, 0, false},
		{"Invalid string", ResultRow{"invalid"}, 0, 0, true},
		{"Invalid type", ResultRow{struct{}{}}, 0, 0, true},
		{"Invalid []byte", ResultRow{[]byte("invalid")}, 0, 0, true},
	}

	for _, tt := range tests {
		result, err := tt.row.GetInt(tt.column)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestResultRow_GetUint(t *testing.T) {
	tests := []struct {
		name     string
		row      ResultRow
		column   int
		expected uint64
		hasError bool
	}{
		{"Valid int64", ResultRow{uint64(10)}, 0, 10, false},
		{"Valid uint64", ResultRow{int64(20)}, 0, 20, false},
		{"Valid float64", ResultRow{float64(30.5)}, 0, 30, false},
		{"Valid string", ResultRow{"40"}, 0, 40, false},
		{"Valid nil", ResultRow{nil}, 0, 0, false},
		{"Valid decimal", ResultRow{decimal.NewFromFloat(123.45)}, 0, 123, false},
		{"Valid small decimal", ResultRow{decimal.New(1, -10)}, 0, 0, false},
		{"Invalid string", ResultRow{"invalid"}, 0, 0, true},
		{"Invalid type", ResultRow{struct{}{}}, 0, 0, true},
		{"Invalid []byte", ResultRow{[]byte("invalid")}, 0, 0, true},
	}

	for _, tt := range tests {
		result, err := tt.row.GetUint(tt.column)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestResultRow_GetFloat(t *testing.T) {
	tests := []struct {
		name     string
		row      ResultRow
		column   int
		expected float64
		hasError bool
	}{
		{"Valid float64", ResultRow{float64(10.5)}, 0, 10.5, false},
		{"Valid int64", ResultRow{int64(20)}, 0, 20, false},
		{"Valid uint64", ResultRow{uint64(30)}, 0, 30, false},
		{"Valid string", ResultRow{"40.5"}, 0, 40.5, false},
		{"Valid nil", ResultRow{nil}, 0, 0, false},
		{"Valid decimal", ResultRow{decimal.NewFromFloat(123.45)}, 0, 123.45, false},
		{"Valid small decimal", ResultRow{decimal.New(1, -10)}, 0, 1e-10, false},
		{"Invalid string", ResultRow{"invalid"}, 0, 0, true},
		{"Invalid type", ResultRow{struct{}{}}, 0, 0, true},
		{"Invalid []byte", ResultRow{[]byte("invalid")}, 0, 0, true},
	}

	for _, tt := range tests {
		result, err := tt.row.GetFloat(tt.column)
		if tt.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		}
	}
}

func TestAggregateFuncCountMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncCountMerger{aggregateFuncBaseMerger{fieldIndex: 0}}
	from := ResultRow{int64(5)}
	to := ResultRow{int64(3)}

	err := merger.MergeTo(from, to)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), to.GetValue(0))
}

func TestAggregateFuncSumMerger_MergeTo_WithNil(t *testing.T) {
	merger := &AggregateFuncSumMerger{aggregateFuncBaseMerger{fieldIndex: 0}}

	// Define test cases
	tests := []struct {
		from    interface{}
		to      interface{}
		expect  interface{}
		message string
	}{
		// int64
		{nil, int64(5), int64(5), "from is nil, to is int64"},
		{int64(5), nil, int64(5), "from is int64, to is nil"},

		// float64
		{nil, float64(5.5), float64(5.5), "from is nil, to is float64"},
		{float64(5.5), nil, float64(5.5), "from is float64, to is nil"},

		// decimal.Decimal
		{nil, decimal.NewFromFloat(5.5), decimal.NewFromFloat(5.5), "from is nil, to is decimal.Decimal"},
		{decimal.NewFromFloat(5.5), nil, decimal.NewFromFloat(5.5), "from is decimal.Decimal, to is nil"},
	}

	for _, tt := range tests {
		from := ResultRow{tt.from}
		to := ResultRow{tt.to}

		err := merger.MergeTo(from, to)
		assert.NoError(t, err, tt.message)
		assert.Equal(t, tt.expect, to.GetValue(0), tt.message)
	}
}

func TestAggregateFuncSumMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncSumMerger{aggregateFuncBaseMerger{fieldIndex: 0}}

	// Test with int64
	fromInt := ResultRow{int64(5)}
	toInt := ResultRow{int64(3)}
	err := merger.MergeTo(fromInt, toInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), toInt.GetValue(0))

	// Test with uint64
	fromUint := ResultRow{uint64(5)}
	toUint := ResultRow{uint64(3)}
	err = merger.MergeTo(fromUint, toUint)
	assert.NoError(t, err)
	assert.Equal(t, uint64(8), toUint.GetValue(0))

	// Test with float64
	fromFloat := ResultRow{float64(5.5)}
	toFloat := ResultRow{float64(3.5)}
	err = merger.MergeTo(fromFloat, toFloat)
	assert.NoError(t, err)
	assert.Equal(t, float64(9.0), toFloat.GetValue(0))

	// Test with decimal.Decimal
	fromDecimal := ResultRow{decimal.NewFromFloat(5.5)}
	toDecimal := ResultRow{decimal.NewFromFloat(3.5)}
	err = merger.MergeTo(fromDecimal, toDecimal)
	assert.NoError(t, err)
	expectedDecimal := decimal.NewFromFloat(9.0)
	actualDecimal := toDecimal.GetValue(0).(decimal.Decimal)
	assert.True(t, expectedDecimal.Equal(actualDecimal))

	// Test with unsupported type
	fromUnsupported := ResultRow{struct{}{}}
	toUnsupported := ResultRow{int64(3)}
	err = merger.MergeTo(fromUnsupported, toUnsupported)
	assert.Error(t, err)

	// Test with negative int64
	fromNegativeInt := ResultRow{int64(-10)}
	toNegativeInt := ResultRow{int64(5)}
	err = merger.MergeTo(fromNegativeInt, toNegativeInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(-5), toNegativeInt.GetValue(0))

	// Test with negative decimal.Decimal
	fromNegativeDecimal := ResultRow{decimal.NewFromFloat(-10.5)}
	toNegativeDecimal := ResultRow{decimal.NewFromFloat(5.5)}
	err = merger.MergeTo(fromNegativeDecimal, toNegativeDecimal)
	assert.NoError(t, err)
	expectedNegativeDecimal := decimal.NewFromFloat(-5.0)
	actualNegativeDecimal := toNegativeDecimal.GetValue(0).(decimal.Decimal)
	assert.True(t, expectedNegativeDecimal.Equal(actualNegativeDecimal))

	// Test with max int64
	fromMaxInt := ResultRow{int64(1<<63 - 1)}
	toMaxInt := ResultRow{int64(0)}
	err = merger.MergeTo(fromMaxInt, toMaxInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(1<<63-1), toMaxInt.GetValue(0))

	// Test with max uint64
	fromMaxUint := ResultRow{uint64(1<<64 - 1)}
	toMaxUint := ResultRow{uint64(0)}
	err = merger.MergeTo(fromMaxUint, toMaxUint)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1<<64-1), toMaxUint.GetValue(0))

	/*
		// TODO: Test mixing int64 and decimal.Decimal
		fromMixed := ResultRow{int64(5)}
		toMixed := ResultRow{decimal.NewFromFloat(3.5)}
		err = merger.MergeTo(fromMixed, toMixed)
		assert.Error(t, err)

		// TODO: Test mixing float64 and uint64
		fromMixedFloat := ResultRow{float64(5.5)}
		toMixedUint := ResultRow{uint64(3)}
		err = merger.MergeTo(fromMixedFloat, toMixedUint)
		assert.Error(t, err)
	*/

	// Test with multiple columns
	fromMulti := ResultRow{int64(5), decimal.NewFromFloat(10.5)}
	toMulti := ResultRow{int64(3), decimal.NewFromFloat(2.5)}

	for i := 0; i < len(fromMulti); i++ {
		merger.fieldIndex = i
		err = merger.MergeTo(fromMulti, toMulti)
		assert.NoError(t, err)
	}

	// Verify each column's result
	assert.Equal(t, int64(8), toMulti.GetValue(0))                                          // int64 加和
	assert.True(t, decimal.NewFromFloat(13.0).Equal(toMulti.GetValue(1).(decimal.Decimal))) // decimal 加和

}

func TestAggregateFuncMaxMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncMaxMerger{aggregateFuncBaseMerger{fieldIndex: 0}}

	// Test with int64
	fromInt := ResultRow{int64(10)}
	toInt := ResultRow{int64(5)}
	err := merger.MergeTo(fromInt, toInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), toInt.GetValue(0))

	// Test with higher value in 'from'
	fromHigher := ResultRow{int64(15)}
	err = merger.MergeTo(fromHigher, toInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), toInt.GetValue(0))
	// Test with uint64
	fromUint := ResultRow{uint64(20)}
	toUint := ResultRow{uint64(15)}
	err = merger.MergeTo(fromUint, toUint)
	assert.NoError(t, err)
	assert.Equal(t, uint64(20), toUint.GetValue(0))

	// Test with float64
	fromFloat := ResultRow{float64(7.5)}
	toFloat := ResultRow{float64(3.5)}
	err = merger.MergeTo(fromFloat, toFloat)
	assert.NoError(t, err)
	assert.Equal(t, float64(7.5), toFloat.GetValue(0))

	// Test with decimal.Decimal
	fromDecimal := ResultRow{decimal.NewFromFloat(12.5)}
	toDecimal := ResultRow{decimal.NewFromFloat(10.5)}
	err = merger.MergeTo(fromDecimal, toDecimal)
	assert.NoError(t, err)
	expectedDecimal := decimal.NewFromFloat(12.5)
	actualDecimal := toDecimal.GetValue(0).(decimal.Decimal)
	assert.True(t, expectedDecimal.Equal(actualDecimal))

	// Test with nil value
	fromNil := ResultRow{nil}
	toNil := ResultRow{int64(10)}
	err = merger.MergeTo(fromNil, toNil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), toNil.GetValue(0)) // No change for nil

	// Test with unsupported type
	fromUnsupported := ResultRow{struct{}{}}
	toUnsupported := ResultRow{int64(10)}
	err = merger.MergeTo(fromUnsupported, toUnsupported)
	assert.Error(t, err)

	// Test with lower value in 'from'
	fromLower := ResultRow{int64(5)}
	err = merger.MergeTo(fromLower, toInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), toInt.GetValue(0)) // 'to' remains unchanged

}

func TestAggregateFuncMinMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncMinMerger{aggregateFuncBaseMerger{fieldIndex: 0}}
	// Test with int64
	fromInt := ResultRow{int64(5)}
	toInt := ResultRow{int64(10)}
	err := merger.MergeTo(fromInt, toInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), toInt.GetValue(0))

	// Test with lower value in 'from'
	fromLowerInt := ResultRow{int64(3)}
	err = merger.MergeTo(fromLowerInt, toInt)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), toInt.GetValue(0))

	// Test with uint64
	fromUint := ResultRow{uint64(5)}
	toUint := ResultRow{uint64(10)}
	err = merger.MergeTo(fromUint, toUint)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5), toUint.GetValue(0))

	// Test with lower value in 'from'
	fromLowerUint := ResultRow{uint64(3)}
	err = merger.MergeTo(fromLowerUint, toUint)
	assert.NoError(t, err)
	assert.Equal(t, uint64(3), toUint.GetValue(0))

	// Test with float64
	fromFloat := ResultRow{float64(3.5)}
	toFloat := ResultRow{float64(7.5)}
	err = merger.MergeTo(fromFloat, toFloat)
	assert.NoError(t, err)
	assert.Equal(t, float64(3.5), toFloat.GetValue(0))

	// Test with higher value in 'from'
	fromHigherFloat := ResultRow{float64(9.5)}
	err = merger.MergeTo(fromHigherFloat, toFloat)
	assert.NoError(t, err)
	assert.Equal(t, float64(3.5), toFloat.GetValue(0))

	// Test with decimal.Decimal
	fromDecimal := ResultRow{decimal.NewFromFloat(2.5)}
	toDecimal := ResultRow{decimal.NewFromFloat(5.5)}
	err = merger.MergeTo(fromDecimal, toDecimal)
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat(2.5).Equal(toDecimal.GetValue(0).(decimal.Decimal)))

	// Test with higher value in 'from'
	fromHigherDecimal := ResultRow{decimal.NewFromFloat(6.5)}
	err = merger.MergeTo(fromHigherDecimal, toDecimal)
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat(2.5).Equal(toDecimal.GetValue(0).(decimal.Decimal)))

	// Test with nil value in 'from'
	fromNil := ResultRow{nil}
	toNotNil := ResultRow{int64(10)}
	err = merger.MergeTo(fromNil, toNotNil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), toNotNil.GetValue(0)) // No change for nil

	// Test with nil value in 'to'
	fromNotNil := ResultRow{int64(10)}
	toNil := ResultRow{nil}
	err = merger.MergeTo(fromNotNil, toNil)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), toNil.GetValue(0))

	// Test with unsupported type
	fromUnsupported := ResultRow{struct{}{}}
	toUnsupported := ResultRow{int64(10)}
	err = merger.MergeTo(fromUnsupported, toUnsupported)
	assert.Error(t, err)
}

func TestAggregateFuncGroupConcatMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncGroupConcatMerger{aggregateFuncBaseMerger{fieldIndex: 0, distinct: false}}
	from := ResultRow{"apple"}
	to := ResultRow{"banana"}

	err := merger.MergeTo(from, to)
	assert.NoError(t, err)
	assert.Equal(t, "banana,apple", to.GetValue(0))

	// Test with distinct
	merger.distinct = true
	fromDistinct := ResultRow{"banana"}
	toDistinct := ResultRow{"banana,apple"}

	err = merger.MergeTo(fromDistinct, toDistinct)
	assert.NoError(t, err)
	assert.Equal(t, "banana,apple", toDistinct.GetValue(0))
}

func TestMergeExecResult(t *testing.T) {
	rs := []*mysql.Result{
		{AffectedRows: 1, InsertID: 2},
		{AffectedRows: 3, InsertID: 0},
	}

	result, err := MergeExecResult(rs)
	assert.NoError(t, err)
	assert.Equal(t, uint64(4), result.AffectedRows)
	assert.Equal(t, uint64(2), result.InsertID) // 最后插入 ID
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		hasError bool
	}{
		// Integer types
		{"int8", int8(127), "127", false},
		{"int16", int16(32767), "32767", false},
		{"int32", int32(2147483647), "2147483647", false},
		{"int64", int64(9223372036854775807), "9223372036854775807", false},
		{"int", int(42), "42", false},

		// Unsigned integer types
		{"uint8", uint8(255), "255", false},
		{"uint16", uint16(65535), "65535", false},
		{"uint32", uint32(4294967295), "4294967295", false},
		{"uint64", uint64(18446744073709551615), "18446744073709551615", false},
		{"uint", uint(99), "99", false},

		// Float types
		// Floating point numbers (float32 and float64) have limited precision and cannot represent some decimal numbers exactly.
		{"float32", float32(123.456), "123.45600128173828", false},
		{"float32", float32(123.45600), "123.45600128173828", false},
		{"float64", float64(9876.54321), "9876.54321", false},
		{"float64", float64(9876.54321000), "9876.54321", false},

		// String types
		{"string", "hello", "hello", false},

		// Byte slice
		{"[]byte", []byte("byte slice"), "byte slice", false},

		// Decimal type (NewFromFloat)
		{"decimal.Decimal", decimal.NewFromFloat(12345.6789), "12345.6789", false},
		{"decimal.Decimal", decimal.NewFromFloat(12345.678900), "12345.6789", false},

		// Decimal type (NewFromString)
		{"decimal.NewFromString valid", func() interface{} {
			d, _ := decimal.NewFromString("12345.6789")
			return d
		}(), "12345.6789", false},
		{"decimal.NewFromString trailing zeros", func() interface{} {
			d, _ := decimal.NewFromString("12345.678900")
			return d
		}(), "12345.678900", false},
		{"decimal.NewFromString invalid", func() interface{} {
			d, err := decimal.NewFromString("invalid")
			if err != nil {
				return err
			}
			return d
		}(), "", true},

		// Decimal type (NewFromInt)
		{"decimal.NewFromInt", decimal.NewFromInt(12345), "12345", false},
		{"decimal.NewFromInt negative", decimal.NewFromInt(-12345), "-12345", false},

		// Decimal type (NewFromBigInt)
		{"decimal.NewFromBigInt", func() interface{} {
			return decimal.NewFromBigInt(big.NewInt(1234567890), 0)
		}(), "1234567890", false},
		{"decimal.NewFromBigInt with scale", func() interface{} {
			return decimal.NewFromBigInt(big.NewInt(1234567890), -2)
		}(), "12345678.90", false},

		// Nil value
		{"nil", nil, "NULL", false},

		// Invalid type
		{"invalid type", struct{}{}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := formatValue(tt.input)
			if (err != nil) != tt.hasError {
				t.Fatalf("expected error: %v, got: %v", tt.hasError, err)
			}
			if string(output) != tt.expected && !tt.hasError {
				t.Fatalf("expected: %s, got: %s", tt.expected, string(output))
			}
		})
	}
}
