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
	"io"
	"sort"

	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
	"github.com/XiaoMi/Gaea/parser/types"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/util"
)

// type check
var _ ast.ExprNode = &PatternInExprDecorator{}

// PatternInExprDecorator decorate PatternInExpr
// 这里记录tableIndexes和indexValueMap是没有问题的, 因为如果是OR条件, 导致路由索引[]int变宽,
// 改写的SQL只是IN这一项没有值, 并不会影响SQL正确性和执行结果.
type PatternInExprDecorator struct {
	Expr ast.ExprNode
	List []ast.ExprNode
	Not  bool

	tableIndexes  []int
	indexValueMap map[int][]ast.ExprNode // tableIndex - valueList

	rule   router.Rule
	result *RouteResult
}

// NeedCreatePatternInExprDecorator check if PatternInExpr needs decoration
func NeedCreatePatternInExprDecorator(p *TableAliasStmtInfo, n *ast.PatternInExpr) (router.Rule, bool, bool, error) {
	if n.Sel != nil {
		return nil, false, false, fmt.Errorf("TableName does not support Sel in sharding")
	}

	// 如果不是ColumnNameExpr, 则不做任何路由计算和装饰, 直接返回
	columnNameExpr, ok := n.Expr.(*ast.ColumnNameExpr)
	if !ok {
		return nil, false, false, nil
	}

	rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInCondition(p, columnNameExpr)
	if err != nil {
		return nil, false, false, fmt.Errorf("check ColumnName error: %v", err)
	}

	if !need && rule == nil {
		return nil, false, false, nil
	}

	// ColumnName不需要装饰, 不代表PatternInExpr不需要装饰, 对全局表来说, PatternInExpr也需要装饰
	if rule.GetType() == router.GlobalTableRuleType {
		return rule, true, isAlias, nil
	}

	return rule, need, isAlias, nil
}

// CreatePatternInExprDecorator create PatternInExprDecorator
// 必须先检查是否需要装饰
func CreatePatternInExprDecorator(n *ast.PatternInExpr, rule router.Rule, isAlias bool, result *RouteResult) (*PatternInExprDecorator, error) {
	columnNameExpr := n.Expr.(*ast.ColumnNameExpr)
	columnNameExprDecorator := CreateColumnNameExprDecorator(columnNameExpr, rule, isAlias, result)

	tableIndexes, indexValueMap, err := getPatternInRouteResult(columnNameExpr.Name, n.Not, rule, n.List)
	if err != nil {
		return nil, fmt.Errorf("getPatternInRouteResult error: %v", err)
	}

	ret := &PatternInExprDecorator{
		Expr:          columnNameExprDecorator,
		List:          n.List,
		Not:           n.Not,
		rule:          rule,
		result:        result,
		tableIndexes:  tableIndexes,
		indexValueMap: indexValueMap,
	}

	return ret, nil
}

// 返回路由, 并构建路由索引到值的映射.
// 如果是分片条件, 则构建值到索引的映射.
// 例如, 1,2,3,4分别映射到索引0,2则[]int = [0,2], map=[0:[1,2], 2:[3,4]]
// 如果是全路由, 则每个分片都要返回所有的值.
func getPatternInRouteResult(n *ast.ColumnName, isNotIn bool, rule router.Rule, values []ast.ExprNode) ([]int, map[int][]ast.ExprNode, error) {
	// 如果是全局表, 则返回广播路由
	if rule.GetType() == router.GlobalTableRuleType {
		indexes := rule.GetSubTableIndexes()
		valueMap := getBroadcastValueMap(indexes, values)
		return indexes, valueMap, nil
	}

	if err := checkValueType(values); err != nil {
		return nil, nil, fmt.Errorf("check value error: %v", err)
	}

	_, _, column := getColumnInfoFromColumnName(n)

	if isNotIn {
		indexes := rule.GetSubTableIndexes()
		valueMap := getBroadcastValueMap(indexes, values)
		return indexes, valueMap, nil
	}
	if rule.GetShardingColumn() != column {
		indexes := rule.GetSubTableIndexes()
		valueMap := getBroadcastValueMap(indexes, values)
		return indexes, valueMap, nil
	}

	var indexes []int
	valueMap := make(map[int][]ast.ExprNode)
	for _, vi := range values {
		v, _ := vi.(*driver.ValueExpr)
		value, err := util.GetValueExprResult(v)
		if err != nil {
			return nil, nil, err
		}
		idx, err := rule.FindTableIndex(value)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := valueMap[idx]; !ok {
			indexes = append(indexes, idx)
		}
		valueMap[idx] = append(valueMap[idx], vi)
	}
	sort.Ints(indexes)
	return indexes, valueMap, nil
}

// 所有的值类型必须为*driver.ValueExpr
func checkValueType(values []ast.ExprNode) error {
	for i, v := range values {
		if _, ok := v.(*driver.ValueExpr); !ok {
			return fmt.Errorf("value is not ValueExpr, index: %d, type: %T", i, v)
		}
	}
	return nil
}

func getBroadcastValueMap(subTableIndexes []int, nodes []ast.ExprNode) map[int][]ast.ExprNode {
	ret := make(map[int][]ast.ExprNode)
	for _, idx := range subTableIndexes {
		ret[idx] = append(ret[idx], nodes...)
	}
	return ret
}

// GetCurrentRouteResult get route result of current decorator
func (p *PatternInExprDecorator) GetCurrentRouteResult() []int {
	return p.tableIndexes
}

// Restore implement ast.Node
func (p *PatternInExprDecorator) Restore(ctx *format.RestoreCtx) error {
	tableIndex, err := p.result.GetCurrentTableIndex()
	if err != nil {
		return err
	}

	if err := p.Expr.Restore(ctx); err != nil {
		return fmt.Errorf("an error occurred while restore PatternInExpr.Expr: %v", err)
	}
	if p.Not {
		ctx.WriteKeyWord(" NOT IN ")
	} else {
		ctx.WriteKeyWord(" IN ")
	}

	ctx.WritePlain("(")
	for i, expr := range p.indexValueMap[tableIndex] {
		if i != 0 {
			ctx.WritePlain(",")
		}
		if err := expr.Restore(ctx); err != nil {
			return fmt.Errorf("an error occurred while restore PatternInExpr.List[%d], err: %v", i, err)
		}
	}
	ctx.WritePlain(")")

	return nil
}

// Accept implement ast.Node
func (p *PatternInExprDecorator) Accept(v ast.Visitor) (node ast.Node, ok bool) {
	return p, ok
}

// Text implement ast.Node
func (p *PatternInExprDecorator) Text() string {
	return ""
}

// SetText implement ast.Node
func (p *PatternInExprDecorator) SetText(text string) {
	return
}

// SetType implement ast.ExprNode
func (p *PatternInExprDecorator) SetType(tp *types.FieldType) {
	return
}

// GetType implement ast.ExprNode
func (p *PatternInExprDecorator) GetType() *types.FieldType {
	return nil
}

// SetFlag implement ast.ExprNode
func (p *PatternInExprDecorator) SetFlag(flag uint64) {
	return
}

// GetFlag implement ast.ExprNode
func (p *PatternInExprDecorator) GetFlag() uint64 {
	return 0
}

// Format implement ast.ExprNode
func (p *PatternInExprDecorator) Format(w io.Writer) {
	return
}
