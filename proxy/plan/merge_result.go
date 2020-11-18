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
	"strconv"
	"strings"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/util/hack"
	"github.com/XiaoMi/Gaea/util/math"
)

// ResultRow is one Row in Result
type ResultRow []interface{}

// GetInt get int value from column
// copy from Resultset.GetInt()
func (r ResultRow) GetInt(column int) (int64, error) {
	d := r[column]
	switch v := d.(type) {
	case uint64:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

// GetUint get uint64 value from column
func (r ResultRow) GetUint(column int) (uint64, error) {
	d := r[column]
	switch v := d.(type) {
	case uint64:
		return v, nil
	case int64:
		return uint64(v), nil
	case float64:
		return uint64(v), nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	case []byte:
		return strconv.ParseUint(string(v), 10, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

// GetFloat get float64 value from column
func (r ResultRow) GetFloat(column int) (float64, error) {
	d := r[column]
	switch v := d.(type) {
	case float64:
		return v, nil
	case uint64:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("data type is %T", v)
	}
}

// SetValue set value to column
func (r ResultRow) SetValue(column int, value interface{}) {
	r[column] = value
}

// GetValue get value from column
func (r ResultRow) GetValue(column int) interface{} {
	return r[column]
}

// AggregateFuncMerger is the merger of aggregate function
type AggregateFuncMerger interface {
	// MergeTo 合并结果集, from为待合并行, to为结果聚合行
	MergeTo(from, to ResultRow) error
}

type aggregateFuncBaseMerger struct {
	fieldIndex int // 所在列位置
}

// CreateAggregateFunctionMerger create AggregateFunctionMerger by function type
// currently support: "count", "sum", "max", "min"
func CreateAggregateFunctionMerger(funcType string, fieldIndex int) (AggregateFuncMerger, error) {
	switch strings.ToLower(funcType) {
	case "count":
		ret := new(AggregateFuncCountMerger)
		ret.fieldIndex = fieldIndex
		return ret, nil
	case "sum":
		ret := new(AggregateFuncSumMerger)
		ret.fieldIndex = fieldIndex
		return ret, nil
	case "max":
		ret := new(AggregateFuncMaxMerger)
		ret.fieldIndex = fieldIndex
		return ret, nil
	case "min":
		ret := new(AggregateFuncMinMerger)
		ret.fieldIndex = fieldIndex
		return ret, nil
	default:
		return nil, fmt.Errorf("aggregate function type is not support: %s", funcType)
	}
}

// AggregateFuncCountMerger merge COUNT() column in result
type AggregateFuncCountMerger struct {
	aggregateFuncBaseMerger
}

// MergeTo implement AggregateFuncMerger
func (a *AggregateFuncCountMerger) MergeTo(from, to ResultRow) error {
	idx := a.fieldIndex
	if idx >= len(from) || idx >= len(to) {
		return fmt.Errorf("field index out of bound: %d", a.fieldIndex)
	}

	valueToMerge, err := from.GetInt(idx)
	if err != nil {
		return fmt.Errorf("get from int value error: %v", err)
	}
	originValue, err := to.GetInt(idx)
	if err != nil {
		return fmt.Errorf("get to int value error: %v", err)
	}
	to.SetValue(idx, originValue+valueToMerge)
	return nil
}

// AggregateFuncSumMerger merge SUM() column in result
type AggregateFuncSumMerger struct {
	aggregateFuncBaseMerger
}

// MergeTo implement AggregateFuncMerger
func (a *AggregateFuncSumMerger) MergeTo(from, to ResultRow) error {
	idx := a.fieldIndex
	if idx >= len(from) || idx >= len(to) {
		return fmt.Errorf("field index out of bound: %d", a.fieldIndex)
	}

	fromValueI := from.GetValue(idx)

	// nil对应NULL, NULL不参与比较
	if fromValueI == nil {
		return nil
	}

	switch to.GetValue(idx).(type) {
	case int64:
		return a.sumToInt64(from, to)
	case uint64:
		return a.sumToUint64(from, to)
	case float64, string, []byte, nil:
		return a.sumToFloat64(from, to)
	default:
		fromValue := from.GetValue(idx)
		toValue := to.GetValue(idx)
		return fmt.Errorf("cannot sum value %v (%T) to %v (%T)", fromValue, fromValue, toValue, toValue)
	}
}

func (a *AggregateFuncSumMerger) sumToInt64(from, to ResultRow) error {
	idx := a.fieldIndex // does not need to check
	valueToMerge, err := from.GetInt(idx)
	if err != nil {
		return fmt.Errorf("get from int value error: %v", err)
	}
	originValue, err := to.GetInt(idx)
	if err != nil {
		return fmt.Errorf("get to int value error: %v", err)
	}
	to.SetValue(idx, originValue+valueToMerge)
	return nil
}

func (a *AggregateFuncSumMerger) sumToUint64(from, to ResultRow) error {
	idx := a.fieldIndex // does not need to check
	valueToMerge, err := from.GetUint(idx)
	if err != nil {
		return fmt.Errorf("get from int value error: %v", err)
	}
	originValue, err := to.GetUint(idx)
	if err != nil {
		return fmt.Errorf("get to int value error: %v", err)
	}
	to.SetValue(idx, originValue+valueToMerge)
	return nil
}

func (a *AggregateFuncSumMerger) sumToFloat64(from, to ResultRow) error {
	idx := a.fieldIndex // does not need to check
	valueToMerge, err := from.GetFloat(idx)
	if err != nil {
		return fmt.Errorf("get from int value error: %v", err)
	}
	originValue, err := to.GetFloat(idx)
	if err != nil {
		return fmt.Errorf("get to int value error: %v", err)
	}
	to.SetValue(idx, originValue+valueToMerge)
	return nil
}

// AggregateFuncMaxMerger merge MAX() column in result
type AggregateFuncMaxMerger struct {
	aggregateFuncBaseMerger
}

// MergeTo implement AggregateFuncMerger
func (a *AggregateFuncMaxMerger) MergeTo(from, to ResultRow) error {
	idx := a.fieldIndex
	if idx >= len(from) || idx >= len(to) {
		return fmt.Errorf("field index out of bound: %d", a.fieldIndex)
	}

	fromValueI := from.GetValue(idx)
	toValueI := to.GetValue(idx)

	// nil对应NULL, NULL不参与比较
	if fromValueI == nil {
		return nil
	}

	switch toValue := toValueI.(type) {
	case nil:
		to.SetValue(idx, fromValueI)
		return nil
	case int64:
		if fromValueI.(int64) > toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
	case uint64:
		if fromValueI.(uint64) > toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
	case float64:
		if fromValueI.(float64) > toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
	case string:
		if fromValueI.(string) > toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
	// does not handle []byte
	default:
		return fmt.Errorf("cannot compare value %v (%T) to %v (%T)", fromValueI, fromValueI, toValueI, toValueI)
	}
}

// AggregateFuncMinMerger merge MIN() column in result
type AggregateFuncMinMerger struct {
	aggregateFuncBaseMerger
}

// MergeTo implement AggregateFuncMerger
func (a *AggregateFuncMinMerger) MergeTo(from, to ResultRow) error {
	idx := a.fieldIndex
	if idx >= len(from) || idx >= len(to) {
		return fmt.Errorf("field index out of bound: %d", a.fieldIndex)
	}

	fromValueI := from.GetValue(idx)
	toValueI := to.GetValue(idx)

	// nil对应NULL, NULL不参与比较
	if fromValueI == nil {
		return nil
	}

	switch toValue := toValueI.(type) {
	case nil:
		to.SetValue(idx, fromValueI)
		return nil
	case int64:
		if fromValueI.(int64) < toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
	case uint64:
		if fromValueI.(uint64) < toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
	case float64:
		if fromValueI.(float64) < toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
	case string:
		if fromValueI.(string) < toValue {
			to.SetValue(idx, fromValueI)
		}
		return nil
		// does not handle []byte
	default:
		return fmt.Errorf("cannot compare value %v (%T) to %v (%T)", fromValueI, fromValueI, toValueI, toValueI)
	}
}

// MergeExecResult merge execution results, like UPDATE, INSERT, DELETE, ...
func MergeExecResult(rs []*mysql.Result) (*mysql.Result, error) {
	r := new(mysql.Result)
	for _, v := range rs {
		r.Status |= v.Status
		r.AffectedRows += v.AffectedRows
		if r.InsertID == 0 {
			r.InsertID = v.InsertID
		} else if v.InsertID != 0 && r.InsertID > v.InsertID {
			//last insert id is first gen id for multi row inserted
			//see http://dev.mysql.com/doc/refman/5.6/en/information-functions.html#function_last-insert-id
			r.InsertID = v.InsertID
		}
	}

	return r, nil
}

// MergeSelectResult merge select results
func MergeSelectResult(p *SelectPlan, stmt *ast.SelectStmt, rs []*mysql.Result) (*mysql.Result, error) {
	ret := mergeMultiResultSet(rs)

	if p.distinct {
		if err := removeDistinctRowInResult(p, ret); err != nil {
			return nil, err
		}
	}

	if stmt.GroupBy != nil {
		if err := buildSelectGroupByResult(p, ret); err != nil {
			return nil, err
		}
	} else {
		if err := buildSelectOnlyResult(p, ret); err != nil {
			return nil, err
		}
	}

	if err := sortSelectResult(p, stmt, ret); err != nil {
		return nil, err
	}

	if err := limitSelectResult(p, ret); err != nil {
		return nil, err
	}

	if err := trimExtraFields(p, ret); err != nil {
		return nil, fmt.Errorf("trimExtraFields error: %v", err)
	}

	if err := GenerateSelectResultRowData(ret); err != nil {
		return nil, fmt.Errorf("generate RowData error: %v", err)
	}

	return ret, nil
}

// 合并结果集, 返回一个Result
func mergeMultiResultSet(rs []*mysql.Result) *mysql.Result {
	if len(rs) == 1 {
		return rs[0]
	}

	// 列信息认为相同, 因此只合并结果
	for i := 1; i < len(rs); i++ {
		rs[0].Status |= rs[i].Status
		rs[0].Values = append(rs[0].Values, rs[i].Values...)
		rs[0].RowDatas = append(rs[0].RowDatas, rs[i].RowDatas...)
	}

	return rs[0]
}

func removeDistinctRowInResult(p *SelectPlan, r *mysql.Result) error {
	distinctKeySet := make(map[string]bool)
	var rowToRemove []int
	// 计算除补列之外的原始列数
	resultFieldLength := len(r.Fields)
	originColumnCount := p.GetColumnCount()
	deltaColumnCount := resultFieldLength - originColumnCount
	colCnt := p.originColumnCount + deltaColumnCount

	// 根据原始列的值进行去重
	rowCount := len(r.Values)
	for i := 0; i < rowCount; i++ {
		keySlice := r.Values[i][0:colCnt]
		mk, err := generateMapKey(keySlice)
		if err != nil {
			return err
		}

		_, ok := distinctKeySet[mk]
		if !ok {
			distinctKeySet[mk] = true
		} else {
			rowToRemove = append(rowToRemove, i)
		}
	}

	rowToRemoveCnt := len(rowToRemove)
	if rowToRemoveCnt == 0 {
		return nil
	}

	originRows := r.Values
	r.RowDatas = nil
	r.Values = originRows[:0]

	var j int
	for i := 0; i < rowCount; i++ {
		if j == rowToRemoveCnt {
			r.Values = append(r.Values, originRows[i:]...)
			break
		}
		if i == rowToRemove[j] {
			j++
		} else {
			r.Values = append(r.Values, originRows[i])
		}
	}

	return nil
}

// contains mergeGroupByWithoutFunc() and mergeGroupByWithFunc()
func buildSelectGroupByResult(p *SelectPlan, r *mysql.Result) error {
	resultMap := make(map[string]ResultRow)

	resultFieldLength := len(r.Fields)
	originColumnCount := p.GetColumnCount()
	deltaColumnCount := resultFieldLength - originColumnCount

	// 根据group by的列进行结果聚合
	for i, v := range r.Values {
		keySlice := make([]interface{}, 0)
		for _, index := range p.GetGroupByColumnInfo() {
			keySlice = append(keySlice, v[index+deltaColumnCount])
		}
		mk, err := generateMapKey(keySlice)
		if err != nil {
			return err
		}

		// 用找到的第一个结果行作为聚合结果
		_, ok := resultMap[mk]
		if !ok {
			resultMap[mk] = ResultRow(r.Values[i])
			continue
		}

		if len(p.aggregateFuncs) == 0 {
			continue
		}

		// 如果存在聚合函数, 则对聚合列进行结果聚合, 非聚合列不处理
		retToMerge := ResultRow(r.Values[i])
		for _, mfunc := range p.aggregateFuncs {
			if err := mfunc.MergeTo(retToMerge, resultMap[mk]); err != nil {
				return fmt.Errorf("MergeTo error, func: %v, value: %v, err: %v", mfunc, retToMerge, err)
			}
		}
	}

	err := buildResultFromResultMap(r, resultMap)
	if err != nil {
		return fmt.Errorf("buildResultFromResultMap error: %v", err)
	}

	return nil
}

func buildSelectOnlyResult(p *SelectPlan, rs *mysql.Result) error {
	r := rs.Resultset
	// 没有聚合函数, 直接把所有分片结果添加到同一个ResultSet下面
	if len(p.aggregateFuncs) == 0 {
		return nil
	}

	// 存在聚合函数, 需要改写聚合列的值, 然后返回 (应该只有一行记录)
	isSet := false
	var currRet ResultRow
	for i, v := range r.Values {
		if !isSet {
			isSet = true
			currRet = ResultRow(v)
			continue
		}

		retToMerge := ResultRow(r.Values[i])
		for _, mfunc := range p.aggregateFuncs {
			if err := mfunc.MergeTo(retToMerge, currRet); err != nil {
				return fmt.Errorf("MergeTo error, func: %v, value: %v, err: %v", mfunc, retToMerge, err)
			}
		}
	}

	r.Values = r.Values[:0]
	r.Values = append(r.Values, currRet)
	r.RowDatas = nil

	return nil
}

// this function modifies the first value of origin results
func buildResultFromResultMap(r *mysql.Result, resultMap map[string]ResultRow) error {
	// no group by result means the result row count is 0, so return the first result
	if len(resultMap) == 0 {
		return nil
	}

	r.Values = nil
	r.RowDatas = nil
	for _, v := range resultMap {
		r.Values = append(r.Values, v)
	}

	return nil
}

// 去掉补充的列
// 与补充列的顺序相反, 先去掉ORDER BY补充的列, 再去掉GROUP BY补充的列
func trimExtraFields(p *SelectPlan, r *mysql.Result) error {
	resultFieldLength := len(r.Fields)
	originColumnCount := p.GetColumnCount()
	deltaColumnCount := resultFieldLength - originColumnCount
	extraFieldStartIndex := deltaColumnCount + p.GetOriginColumnCount()

	if extraFieldStartIndex != -1 {
		r.Fields = r.Fields[0:extraFieldStartIndex]
		for i := 0; i < len(r.Values); i++ {
			r.Values[i] = r.Values[i][0:extraFieldStartIndex]
		}
	}

	return nil
}

func sortSelectResult(p *SelectPlan, stmt *ast.SelectStmt, ret *mysql.Result) error {
	if !p.HasOrderBy() {
		return nil
	}

	resultFieldLength := len(ret.Fields)
	originColumnCount := p.GetColumnCount()
	deltaColumnCount := resultFieldLength - originColumnCount

	orderByColumns, orderByDirections := p.GetOrderByColumnInfo()
	var sortKeys []mysql.SortKey
	for i := 0; i < len(orderByDirections); i++ {
		sortKey := mysql.SortKey{}
		sortKey.Column = orderByColumns[i] + deltaColumnCount
		if orderByDirections[i] {
			sortKey.Direction = mysql.SortDesc
		} else {
			sortKey.Direction = mysql.SortAsc
		}
		sortKeys = append(sortKeys, sortKey)
	}

	return ret.SortWithoutColumnName(sortKeys)
}

// the result from backend is aggregated and offset = 0, count = (originOffset + originCount)
func limitSelectResult(p *SelectPlan, ret *mysql.Result) error {
	if !p.HasLimit() {
		return nil
	}

	start, count := p.GetLimitValue()

	rowLen := int64(len(ret.Values))
	end := math.MinInt64(start+count, rowLen)

	if start >= rowLen {
		ret.RowDatas = ret.RowDatas[:0]
		ret.Values = ret.Values[:0]
		return nil
	}

	ret.Values = ret.Values[start:end]
	return nil
}

// GenerateSelectResultRowData generate raw RowData from values
// 根据value反向构造RowData
// copy from server.buildResultset()
func GenerateSelectResultRowData(r *mysql.Result) error {
	r.RowDatas = nil
	for i, vs := range r.Values {
		if len(vs) != len(r.Fields) {
			return fmt.Errorf("row %d has %d column not equal %d", i, len(vs), len(r.Fields))
		}

		var row []byte
		for _, value := range vs {
			// build row values
			if value == nil {
				row = append(row, 0xfb)
			} else {
				b, err := formatValue(value)
				if err != nil {
					return err
				}
				row = mysql.AppendLenEncStringBytes(row, b)
			}
		}

		r.RowDatas = append(r.RowDatas, row)
	}

	return nil
}

// copy from server.generateMapKey()
func generateMapKey(groupColumns []interface{}) (string, error) {
	bk := make([]byte, 0, 8)
	separatorBuf, err := formatValue("+")
	if err != nil {
		return "", err
	}

	for _, v := range groupColumns {
		b, err := formatValue(v)
		if err != nil {
			return "", err
		}
		bk = append(bk, b...)
		bk = append(bk, separatorBuf...)
	}

	return string(bk), nil
}

// copy from server.formatValue()
// formatValue encode value into a string format
func formatValue(value interface{}) ([]byte, error) {
	if value == nil {
		return hack.Slice("NULL"), nil
	}
	switch v := value.(type) {
	case int8:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int16:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int32:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int64:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case int:
		return strconv.AppendInt(nil, int64(v), 10), nil
	case uint8:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint16:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint32:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint64:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case uint:
		return strconv.AppendUint(nil, uint64(v), 10), nil
	case float32:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case float64:
		return strconv.AppendFloat(nil, float64(v), 'f', -1, 64), nil
	case []byte:
		return v, nil
	case string:
		return hack.Slice(v), nil
	default:
		return nil, fmt.Errorf("invalid type %T", value)
	}
}
