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
	"strings"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/opcode"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
	"github.com/shopspring/decimal"
)

// HavingFilter evaluates HAVING expressions on merged result sets
type HavingFilter struct {
	// columnIndexByName maps column name (lowercase) to column index
	columnIndexByName map[string]int
	// columnIndexByAlias maps column alias (lowercase) to column index
	columnIndexByAlias map[string]int
	// fields is the field list from the result set
	fields []*mysql.Field
	// plan contains HAVING aggregate/alias column metadata from planning
	plan *SelectPlan
}

// NewHavingFilter creates a new HavingFilter from result fields
func NewHavingFilter(fields []*mysql.Field, p *SelectPlan) *HavingFilter {
	f := &HavingFilter{
		columnIndexByName:  make(map[string]int),
		columnIndexByAlias: make(map[string]int),
		fields:             fields,
		plan:               p,
	}

	for i, field := range fields {
		name := strings.ToLower(string(field.Name))
		if name != "" {
			f.columnIndexByName[name] = i
		}
		// OrgName is the original column name (before AS)
		orgName := strings.ToLower(string(field.OrgName))
		if orgName != "" && orgName != name {
			f.columnIndexByName[orgName] = i
		}
	}

	return f
}

// FilterResult filters the result set based on HAVING expression
// Returns the filtered result
func (f *HavingFilter) FilterResult(havingExpr ast.ExprNode, r *mysql.Result) error {
	if havingExpr == nil || r == nil || len(r.Values) == 0 {
		return nil
	}

	var filteredValues [][]interface{}

	for _, row := range r.Values {
		match, err := f.evaluateExpr(havingExpr, ResultRow(row))
		if err != nil {
			return fmt.Errorf("evaluate having expression error: %v", err)
		}
		if match {
			filteredValues = append(filteredValues, row)
		}
	}

	r.Values = filteredValues
	r.RowDatas = nil // will be regenerated later

	return nil
}

// evaluateExpr evaluates an expression against a result row and returns a boolean
func (f *HavingFilter) evaluateExpr(expr ast.ExprNode, row ResultRow) (bool, error) {
	switch e := expr.(type) {
	case *ast.BinaryOperationExpr:
		return f.evaluateBinaryOp(e, row)
	case *ast.ParenthesesExpr:
		return f.evaluateExpr(e.Expr, row)
	case *driver.ValueExpr:
		// A constant value - treat non-zero/non-empty as true
		return f.valueAsBool(e.Datum.GetValue())
	default:
		return false, fmt.Errorf("unsupported expression type in HAVING: %T", expr)
	}
}

// evaluateBinaryOp evaluates a binary operation expression
func (f *HavingFilter) evaluateBinaryOp(expr *ast.BinaryOperationExpr, row ResultRow) (bool, error) {
	switch expr.Op {
	case opcode.LogicAnd:
		left, err := f.evaluateExpr(expr.L, row)
		if err != nil {
			return false, err
		}
		if !left {
			return false, nil
		}
		return f.evaluateExpr(expr.R, row)

	case opcode.LogicOr:
		left, err := f.evaluateExpr(expr.L, row)
		if err != nil {
			return false, err
		}
		if left {
			return true, nil
		}
		return f.evaluateExpr(expr.R, row)

	case opcode.EQ, opcode.NE, opcode.GT, opcode.GE, opcode.LT, opcode.LE:
		leftVal, err := f.getExprValue(expr.L, row)
		if err != nil {
			return false, err
		}
		rightVal, err := f.getExprValue(expr.R, row)
		if err != nil {
			return false, err
		}
		return f.compareValues(expr.Op, leftVal, rightVal)

	default:
		return false, fmt.Errorf("unsupported operator in HAVING: %v", expr.Op)
	}
}

// getExprValue gets the value of an expression from the result row
func (f *HavingFilter) getExprValue(expr ast.ExprNode, row ResultRow) (interface{}, error) {
	switch e := expr.(type) {
	case *ast.ColumnNameExpr:
		return f.getColumnValue(e.Name, row)
	case *ast.AggregateFuncExpr:
		// For aggregate functions, we need to find the corresponding column
		// The aggregate function result should already be in the result row
		return f.getAggregateFuncValue(e, row)
	case *driver.ValueExpr:
		return e.Datum.GetValue(), nil
	case *ColumnNameExprDecorator:
		// Handle decorated column name expressions
		// ColumnNameExprDecorator embeds *ast.ColumnNameExpr
		return f.getColumnValue(e.ColumnNameExpr.Name, row)
	default:
		return nil, fmt.Errorf("unsupported expression type for value: %T", expr)
	}
}

// getColumnValue gets the value of a column by name from the result row
func (f *HavingFilter) getColumnValue(name *ast.ColumnName, row ResultRow) (interface{}, error) {
	colName := strings.ToLower(name.Name.L)

	if f.plan != nil {
		if idx, ok := f.plan.getHavingAliasColumnIndex(colName); ok {
			return f.getRowValueByPlanIndex(row, idx, fmt.Sprintf("having alias %s", colName))
		}
	}

	// Try to find by name
	if idx, ok := f.columnIndexByName[colName]; ok {
		if idx < len(row) {
			return row[idx], nil
		}
	}

	// Try to find by alias
	if idx, ok := f.columnIndexByAlias[colName]; ok {
		if idx < len(row) {
			return row[idx], nil
		}
	}

	return nil, fmt.Errorf("column %s not found in result", colName)
}

// getAggregateFuncValue gets the value of an aggregate function from the result row
// It uses planner metadata to resolve the merged aggregate column.
func (f *HavingFilter) getAggregateFuncValue(agg *ast.AggregateFuncExpr, row ResultRow) (interface{}, error) {
	if f.plan != nil {
		if idx, ok, err := f.plan.getHavingAggregateColumnIndex(agg); err != nil {
			return nil, err
		} else if ok {
			key, err := havingAggregateKey(agg)
			if err != nil {
				return nil, err
			}
			return f.getRowValueByPlanIndex(row, idx, fmt.Sprintf("aggregate function %s", key))
		}
	}

	return f.getAggregateFuncValueByFieldName(agg, row)
}

func (f *HavingFilter) getRowValueByPlanIndex(row ResultRow, planIdx int, label string) (interface{}, error) {
	idx := planIdx + f.resultColumnIndexDelta()
	if idx < 0 || idx >= len(row) {
		return nil, fmt.Errorf("%s column index %d out of bounds for row with %d columns", label, idx, len(row))
	}
	return row[idx], nil
}

func (f *HavingFilter) resultColumnIndexDelta() int {
	if f.plan == nil {
		return 0
	}
	return len(f.fields) - f.plan.GetColumnCount()
}

func (f *HavingFilter) getAggregateFuncValueByFieldName(agg *ast.AggregateFuncExpr, row ResultRow) (interface{}, error) {
	key, err := havingAggregateKey(agg)
	if err != nil {
		return nil, err
	}

	for i, field := range f.fields {
		if normalizeHavingFieldName(field.Name) == key || normalizeHavingFieldName(field.OrgName) == key {
			if i >= len(row) {
				return nil, fmt.Errorf("aggregate function %s column index %d out of bounds for row with %d columns", key, i, len(row))
			}
			return row[i], nil
		}
	}

	return nil, fmt.Errorf("aggregate function %s not found in result", agg.F)
}

func normalizeHavingFieldName(name []byte) string {
	return strings.ToLower(strings.TrimSpace(string(name)))
}

// compareValues compares two values using the given operator
func (f *HavingFilter) compareValues(op opcode.Op, left, right interface{}) (bool, error) {
	// Handle nil values
	if left == nil || right == nil {
		// NULL comparisons return false in SQL (except for IS NULL)
		return false, nil
	}

	// Try to compare as numbers
	leftDec, leftErr := f.toDecimal(left)
	rightDec, rightErr := f.toDecimal(right)

	if leftErr == nil && rightErr == nil {
		// Both values can be compared as decimals
		cmp := leftDec.Cmp(rightDec)
		switch op {
		case opcode.EQ:
			return cmp == 0, nil
		case opcode.NE:
			return cmp != 0, nil
		case opcode.GT:
			return cmp > 0, nil
		case opcode.GE:
			return cmp >= 0, nil
		case opcode.LT:
			return cmp < 0, nil
		case opcode.LE:
			return cmp <= 0, nil
		}
	}

	// Fall back to string comparison
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)
	cmp := strings.Compare(leftStr, rightStr)

	switch op {
	case opcode.EQ:
		return cmp == 0, nil
	case opcode.NE:
		return cmp != 0, nil
	case opcode.GT:
		return cmp > 0, nil
	case opcode.GE:
		return cmp >= 0, nil
	case opcode.LT:
		return cmp < 0, nil
	case opcode.LE:
		return cmp <= 0, nil
	}

	return false, fmt.Errorf("unsupported comparison operator: %v", op)
}

// toDecimal converts a value to decimal for comparison
func (f *HavingFilter) toDecimal(val interface{}) (decimal.Decimal, error) {
	switch v := val.(type) {
	case int64:
		return decimal.NewFromInt(v), nil
	case uint64:
		return decimal.NewFromString(fmt.Sprintf("%d", v))
	case float64:
		return decimal.NewFromFloat(v), nil
	case string:
		return decimal.NewFromString(v)
	case []byte:
		return decimal.NewFromString(string(v))
	case decimal.Decimal:
		return v, nil
	default:
		return decimal.Zero, fmt.Errorf("cannot convert %T to decimal", val)
	}
}

// valueAsBool converts a value to boolean
func (f *HavingFilter) valueAsBool(val interface{}) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case int64:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case float64:
		return v != 0, nil
	case string:
		return v != "" && v != "0", nil
	case []byte:
		return len(v) > 0 && string(v) != "0", nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", val)
	}
}
