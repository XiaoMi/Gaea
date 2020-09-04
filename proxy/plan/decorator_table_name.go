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
	"github.com/pingcap/errors"
)

// TableNameDecorator decorate TableName
type TableNameDecorator struct {
	origin *ast.TableName
	rule   router.Rule
	result *RouteResult
}

// NeedCreateTableNameDecoratorWithoutAlias check if TableName without alias needs decorate
// 不带表别名时, 只用StmtInfo就可以判断
func NeedCreateTableNameDecoratorWithoutAlias(p *StmtInfo, n *ast.TableName) (router.Rule, bool, error) {
	db, table := getTableInfoFromTableName(n)

	rule, err := p.RecordShardTable(db, table)
	if err != nil {
		return nil, false, fmt.Errorf("record shard table in SelectPlan error, db: %s, table, %s, err: %v", db, table, err)
	}

	return rule, true, nil
}

// NeedCreateTableNameDecorator check if TableName with alias needs decorate
// SELECT语句可能带有表别名, 需要记录表别名
func NeedCreateTableNameDecorator(p *TableAliasStmtInfo, n *ast.TableName, alias string) (router.Rule, bool, error) {
	db, table := getTableInfoFromTableName(n)

	rule, err := p.RecordShardTable(db, table, alias)
	if err != nil {
		return nil, false, fmt.Errorf("record shard table in SelectPlan error, db: %s, table, %s, alias: %s, err: %v", db, table, alias, err)
	}

	return rule, true, nil
}

// CreateTableNameDecorator create TableNameDecorator
// the table has been checked before
func CreateTableNameDecorator(n *ast.TableName, rule router.Rule, result *RouteResult) (*TableNameDecorator, error) {
	if len(n.PartitionNames) != 0 {
		return nil, fmt.Errorf("TableName does not support PartitionNames in sharding")
	}

	ret := &TableNameDecorator{
		origin: n,
		rule:   rule,
		result: result,
	}

	return ret, nil
}

// Restore implement ast.Node
func (t *TableNameDecorator) Restore(ctx *format.RestoreCtx) error {
	tableIndex, err := t.result.GetCurrentTableIndex()
	if err != nil {
		return err
	}

	ruleType := t.rule.GetType()

	// kingshard不需改写库名, mycat需要改写, 全局表需要改写
	if t.origin.Schema.String() != "" {
		if ruleType == router.GlobalTableRuleType {
			dbName, err := t.rule.GetDatabaseNameByTableIndex(tableIndex)
			if err != nil {
				return fmt.Errorf("get mycat database name error: %v", err)
			}
			ctx.WriteName(dbName)
			ctx.WritePlain(".")
		} else if router.IsMycatShardingRule(ruleType) {
			dbName, err := t.rule.GetDatabaseNameByTableIndex(tableIndex)
			if err != nil {
				return fmt.Errorf("get mycat database name error: %v", err)
			}
			ctx.WriteName(dbName)
			ctx.WritePlain(".")
		} else {
			ctx.WriteName(t.origin.Schema.String())
			ctx.WritePlain(".")
		}
	}

	// kingshard需要改写表名, mycat不需要改写, 全局表需要改写
	if ruleType == router.GlobalTableRuleType {
		ctx.WriteName(t.origin.Name.String())
	} else if router.IsMycatShardingRule(ruleType) {
		ctx.WriteName(t.origin.Name.String())
	} else {
		ctx.WriteName(fmt.Sprintf("%s_%04d", t.origin.Name.String(), tableIndex))
	}

	for _, value := range t.origin.IndexHints {
		ctx.WritePlain(" ")
		if err := value.Restore(ctx); err != nil {
			return errors.Annotate(err, "An error occurred while splicing IndexHints")
		}
	}
	return nil
}

// Accept implement ast.Node
// do nothing and return current decorator
func (t *TableNameDecorator) Accept(v ast.Visitor) (ast.Node, bool) {
	return t, true
}

// Text implement ast.Node
func (t *TableNameDecorator) Text() string {
	return t.origin.Text()
}

// SetText implement ast.Node
func (t *TableNameDecorator) SetText(text string) {
	t.origin.SetText(text)
}
