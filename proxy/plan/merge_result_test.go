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
	"testing"

	"github.com/XiaoMi/Gaea/mysql"
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
		row      ResultRow
		column   int
		expected int64
		hasError bool
	}{
		{ResultRow{int64(10)}, 0, 10, false},
		{ResultRow{uint64(20)}, 0, 20, false},
		{ResultRow{float64(30.5)}, 0, 30, false},
		{ResultRow{"40"}, 0, 40, false},
		{ResultRow{nil}, 0, 0, false},
		{ResultRow{"invalid"}, 0, 0, true},
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
		row      ResultRow
		column   int
		expected uint64
		hasError bool
	}{
		{ResultRow{uint64(10)}, 0, 10, false},
		{ResultRow{int64(20)}, 0, 20, false},
		{ResultRow{float64(30.5)}, 0, 30, false},
		{ResultRow{"40"}, 0, 40, false},
		{ResultRow{nil}, 0, 0, false},
		{ResultRow{"invalid"}, 0, 0, true},
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
		row      ResultRow
		column   int
		expected float64
		hasError bool
	}{
		{ResultRow{float64(10.5)}, 0, 10.5, false},
		{ResultRow{int64(20)}, 0, 20, false},
		{ResultRow{uint64(30)}, 0, 30, false},
		{ResultRow{"40.5"}, 0, 40.5, false},
		{ResultRow{nil}, 0, 0, false},
		{ResultRow{"invalid"}, 0, 0, true},
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

func TestAggregateFuncSumMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncSumMerger{aggregateFuncBaseMerger{fieldIndex: 0}}
	from := ResultRow{int64(5)}
	to := ResultRow{int64(3)}

	err := merger.MergeTo(from, to)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), to.GetValue(0))

	// Test with different data types
	fromFloat := ResultRow{float64(5.5)}
	toFloat := ResultRow{float64(3.5)}
	err = merger.MergeTo(fromFloat, toFloat)
	assert.NoError(t, err)
	assert.Equal(t, float64(9.0), toFloat.GetValue(0))
}

func TestAggregateFuncMaxMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncMaxMerger{aggregateFuncBaseMerger{fieldIndex: 0}}
	from := ResultRow{int64(10)}
	to := ResultRow{int64(5)}

	err := merger.MergeTo(from, to)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), to.GetValue(0))

	// Test with higher value in 'from'
	fromHigher := ResultRow{int64(15)}
	err = merger.MergeTo(fromHigher, to)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), to.GetValue(0))
}

func TestAggregateFuncMinMerger_MergeTo(t *testing.T) {
	merger := &AggregateFuncMinMerger{aggregateFuncBaseMerger{fieldIndex: 0}}
	from := ResultRow{int64(5)}
	to := ResultRow{int64(10)}

	err := merger.MergeTo(from, to)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), to.GetValue(0))

	// Test with lower value in 'from'
	fromLower := ResultRow{int64(3)}
	err = merger.MergeTo(fromLower, to)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), to.GetValue(0))
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
