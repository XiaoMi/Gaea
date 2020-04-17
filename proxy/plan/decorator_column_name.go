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
	"github.com/XiaoMi/Gaea/proxy/router"
)

// ColumnNameExprDecorator decorate ColumnNameExpr to rewrite table name
type ColumnNameExprDecorator struct {
	*ast.ColumnNameExpr
	Name *ColumnNameDecorator
}

// ColumnNameDecorator decorate ColumnName to rewrite table name
type ColumnNameDecorator struct {
	origin  *ast.ColumnName
	rule    router.Rule
	result  *RouteResult
	isAlias bool
}

// NeedCreateColumnNameExprDecoratorInField check if ColumnNameExpr in field needs decoration
// 用于Field列表中判断列名是否需要装饰
// 如果db名和表名都不存在, 则不需要装饰
func NeedCreateColumnNameExprDecoratorInField(p *TableAliasStmtInfo, n *ast.ColumnNameExpr) (router.Rule, bool, bool, error) {
	db, table, _ := getColumnInfoFromColumnName(n.Name)
	if db == "" && table == "" {
		return nil, false, false, nil
	}
	return needCreateColumnNameDecorator(p, n.Name)
}

// NeedCreateColumnNameExprDecoratorInCondition check if ColumnNameExpr in condition needs decoration
// 用于JOIN ON条件或WHERE条件中判断列名是否需要装饰
// 与上面的区别在于, 当只存在列名, 不存在db名和表名时, 还会根据列名去查找对应的条件 (因为装饰之后需要在比较条件中计算路由)
func NeedCreateColumnNameExprDecoratorInCondition(p *TableAliasStmtInfo, n *ast.ColumnNameExpr) (router.Rule, bool, bool, error) {
	return needCreateColumnNameDecorator(p, n.Name)
}

// 是否需要装饰ColumnName, 需要则返回ture
// 在CreateColumnNameDecorator之前调用, 用来检查
// 返回结果bool表示是否需要创建装饰器
func needCreateColumnNameDecorator(p *TableAliasStmtInfo, n *ast.ColumnName) (router.Rule, bool, bool, error) {
	db, table, column := getColumnInfoFromColumnName(n)

	rule, ok, isAlias, err := p.GetSettedRuleFromColumnInfo(db, table, column)
	if err != nil {
		return nil, false, false, err
	}
	if !ok {
		return nil, false, false, nil
	}

	return rule, true, isAlias, nil
}

// CreateColumnNameExprDecorator create ColumnNameExprDecorator
func CreateColumnNameExprDecorator(n *ast.ColumnNameExpr, rule router.Rule, isAlias bool, result *RouteResult) *ColumnNameExprDecorator {
	columnName := createColumnNameDecorator(n.Name, rule, isAlias, result)
	return &ColumnNameExprDecorator{
		ColumnNameExpr: n,
		Name:           columnName,
	}
}

func createColumnNameDecorator(n *ast.ColumnName, rule router.Rule, isAlias bool, result *RouteResult) *ColumnNameDecorator {
	ret := &ColumnNameDecorator{
		origin:  n,
		rule:    rule,
		result:  result,
		isAlias: isAlias,
	}
	return ret
}

// GetColumnInfo get column info, return db, table, column
func (c *ColumnNameDecorator) GetColumnInfo() (string, string, string) {
	return getColumnInfoFromColumnName(c.origin)
}

// Restore implement ast.Node
func (c *ColumnNameDecorator) Restore(ctx *format.RestoreCtx) error {
	tableIndex, err := c.result.GetCurrentTableIndex()
	if err != nil {
		return err
	}

	ruleType := c.rule.GetType()

	// kingshard不需改写库名, mycat需要改写
	if c.origin.Schema.O != "" {
		if ruleType == router.GlobalTableRuleType {
			dbName, err := c.rule.GetDatabaseNameByTableIndex(tableIndex)
			if err != nil {
				return fmt.Errorf("get mycat database name error: %v", err)
			}
			ctx.WriteName(dbName)
			ctx.WritePlain(".")
		} else if router.IsMycatShardingRule(ruleType) {
			dbName, err := c.rule.GetDatabaseNameByTableIndex(tableIndex)
			if err != nil {
				return fmt.Errorf("get mycat database name error: %v", err)
			}
			ctx.WriteName(dbName)
			ctx.WritePlain(".")
		} else {
			ctx.WriteName(c.origin.Schema.String())
			ctx.WritePlain(".")
		}
	}

	// kingshard需要改写表名, mycat不需要改写, 全局表不需要改写
	if c.origin.Table.O != "" {
		if ruleType == router.GlobalTableRuleType {
			ctx.WriteName(c.origin.Table.String())
			ctx.WritePlain(".")
		} else if router.IsMycatShardingRule(ruleType) {
			ctx.WriteName(c.origin.Table.String())
			ctx.WritePlain(".")
		} else {
			if c.isAlias {
				ctx.WriteName(c.origin.Table.String())
				ctx.WritePlain(".")
			} else {
				ctx.WriteName(fmt.Sprintf("%s_%04d", c.origin.Table.String(), tableIndex))
				ctx.WritePlain(".")
			}
		}
	}

	// 列名不需要改写
	ctx.WriteName(c.origin.Name.O)

	return nil
}

// Accept implement ast.Node
// do nothing and return current decorator
func (c *ColumnNameDecorator) Accept(v ast.Visitor) (ast.Node, bool) {
	return c, true
}

// Text implement ast.Node
func (c *ColumnNameDecorator) Text() string {
	return c.origin.Text()
}

// SetText implement ast.Node
func (c *ColumnNameDecorator) SetText(text string) {
	c.origin.SetText(text)
}

// Restore implement ast.Node
func (cc *ColumnNameExprDecorator) Restore(ctx *format.RestoreCtx) error {
	if err := cc.Name.Restore(ctx); err != nil {
		return fmt.Errorf("restore ColumnNameExprDecorator error: %v", err)
	}
	return nil
}
