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

	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/util"
)

// BetweenExprDecorator decorate BetweenExpr
// Between只需要改写表名并计算路由, 不需要改写边界值.
type BetweenExprDecorator struct {
	*ast.BetweenExpr // origin
	column           *ColumnNameExprDecorator

	tableIndexes []int

	rule   router.Rule
	result *RouteResult
}

// NeedCreateBetweenExprDecorator check if BetweenExpr needs decoration
func NeedCreateBetweenExprDecorator(p *TableAliasStmtInfo, n *ast.BetweenExpr) (router.Rule, bool, bool, error) {
	// 如果不是ColumnNameExpr, 则不做任何路由计算和装饰, 直接返回
	columnNameExpr, ok := n.Expr.(*ast.ColumnNameExpr)
	if !ok {
		return nil, false, false, nil
	}

	rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInCondition(p, columnNameExpr)
	if err != nil {
		return nil, false, false, fmt.Errorf("check ColumnName error: %v", err)
	}

	return rule, need, isAlias, nil
}

// CreateBetweenExprDecorator create BetweenExprDecorator
func CreateBetweenExprDecorator(n *ast.BetweenExpr, rule router.Rule, isAlias bool, result *RouteResult) (*BetweenExprDecorator, error) {
	columnNameExpr := n.Expr.(*ast.ColumnNameExpr)
	columnNameExprDecorator := CreateColumnNameExprDecorator(columnNameExpr, rule, isAlias, result)

	tableIndexes, err := getBetweenExprRouteResult(rule, n)
	if err != nil {
		return nil, fmt.Errorf("getBetweenExprRouteResult error: %v", err)
	}

	ret := &BetweenExprDecorator{
		BetweenExpr:  n,
		column:       columnNameExprDecorator,
		tableIndexes: tableIndexes,
		rule:         rule,
		result:       result,
	}
	return ret, nil
}

// Accept do nothing and return current decorator
func (b *BetweenExprDecorator) Accept(v ast.Visitor) (ast.Node, bool) {
	return b, true
}

// Restore column name restore is different from BetweenExpr
func (b *BetweenExprDecorator) Restore(ctx *format.RestoreCtx) error {
	if err := b.column.Restore(ctx); err != nil {
		return fmt.Errorf("an error occurred while restore BetweenExpr.Expr: %v", err)
	}
	if b.Not {
		ctx.WriteKeyWord(" NOT BETWEEN ")
	} else {
		ctx.WriteKeyWord(" BETWEEN ")
	}
	if err := b.Left.Restore(ctx); err != nil {
		return fmt.Errorf("an error occurred while restore BetweenExpr.Left: %v", err)
	}
	ctx.WriteKeyWord(" AND ")
	if err := b.Right.Restore(ctx); err != nil {
		return fmt.Errorf("an error occurred while restore BetweenExpr.Right: %v", err)
	}
	return nil
}

// GetCurrentRouteResult return route result
func (b *BetweenExprDecorator) GetCurrentRouteResult() []int {
	return b.tableIndexes
}

func getBetweenExprRouteResult(rule router.Rule, n *ast.BetweenExpr) ([]int, error) {
	//如果是全局表, 则返回广播路由
	if rule.GetType() == router.GlobalTableRuleType {
		indexes := rule.GetSubTableIndexes()
		return indexes, nil
	}

	columnNameExpr := n.Expr.(*ast.ColumnNameExpr)
	_, _, column := getColumnInfoFromColumnName(columnNameExpr.Name)

	if rule.GetShardingColumn() != column {
		indexes := rule.GetSubTableIndexes()
		return indexes, nil
	}

	if _, ok := rule.GetShard().(router.RangeShard); ok {
		return getShardBetweenExprRouteResult(rule, n)
	}

	indexes := rule.GetSubTableIndexes()
	return indexes, nil
}

// copy from origin PlanBuilder.getRangeShardTableIndex
func getShardBetweenExprRouteResult(rule router.Rule, n *ast.BetweenExpr) ([]int, error) {
	rangeShard := rule.GetShard().(router.RangeShard)

	leftValueExpr, ok := n.Left.(*driver.ValueExpr)
	if !ok {
		return nil, fmt.Errorf("n.Left is not a ValueExpr, type: %T", n.Left)
	}
	leftValue, err := util.GetValueExprResult(leftValueExpr)
	if err != nil {
		return nil, fmt.Errorf("get value from n.Left error: %v", err)
	}

	rightValueExpr, ok := n.Right.(*driver.ValueExpr)
	if !ok {
		return nil, fmt.Errorf("n.Left is not a ValueExpr, type: %T", n.Right)
	}
	rightValue, err := util.GetValueExprResult(rightValueExpr)
	if err != nil {
		return nil, fmt.Errorf("get value from n.Right error: %v", err)
	}

	start, err := rule.FindTableIndex(leftValue)
	if err != nil {
		return nil, fmt.Errorf("FindTableIndex for n.Left error: %v", err)
	}
	last, err := rule.FindTableIndex(rightValue)
	if err != nil {
		return nil, fmt.Errorf("FindTableIndex for n.Right error: %v", err)
	}

	if n.Not {
		if start > last {
			start, last = last, start
			start = adjustShardIndex(rangeShard, rightValue, start)
		} else {
			start = adjustShardIndex(rangeShard, leftValue, start)
		}

		l1 := makeList(rule.GetFirstTableIndex(), start+1)
		l2 := makeList(last, rule.GetLastTableIndex()+1)
		return unionList(l1, l2), nil

	}
	if start > last {
		start, last = last, start
	}
	return makeList(start, last+1), nil
}
