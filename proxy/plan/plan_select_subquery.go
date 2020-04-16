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
	"github.com/XiaoMi/Gaea/proxy/router"
)

// SubqueryColumnNameRewriteVisitor visit ColumnNameExpr in subquery, check if need decorate, and then decorate it.
type SubqueryColumnNameRewriteVisitor struct {
	info *TableAliasStmtInfo
}

// NewSubqueryColumnNameRewriteVisitor consturctor of SubqueryColumnNameRewriteVisitor
func NewSubqueryColumnNameRewriteVisitor(p *TableAliasStmtInfo) *SubqueryColumnNameRewriteVisitor {
	return &SubqueryColumnNameRewriteVisitor{
		info: p,
	}
}

// Enter implement ast.Visitor
func (s *SubqueryColumnNameRewriteVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	return n, false
}

// Leave implement ast.Visitor
func (s *SubqueryColumnNameRewriteVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	field, ok := n.(*ast.ColumnNameExpr)
	if !ok {
		return n, true
	}

	db, table, _ := getColumnInfoFromColumnName(field.Name)

	rule, _ := s.info.getShardRule(db, table)
	if rule == nil || rule.GetType() == router.GlobalTableRuleType {
		return n, true
	}
 
	decorator := CreateColumnNameExprDecorator(field, rule, false, s.info.GetRouteResult())
	return decorator, true
}

func handleSubquerySelectStmt(p *TableAliasStmtInfo, subquery *ast.SelectStmt) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("handleSubqueryExpr panic: %v", v)
		}
	}()

	if err = handleSubqueryTableRefs(p, subquery); err != nil {
		return fmt.Errorf("handle From error: %v", err)
	}

	// 对所有可能含有ColumnName的Node做装饰.
	columnRewritter := NewSubqueryColumnNameRewriteVisitor(p)
	if subquery.Where != nil {
		subquery.Where.Accept(columnRewritter)
	}
	if subquery.Fields != nil {
		subquery.Fields.Accept(columnRewritter)
	}
	if subquery.GroupBy != nil {
		subquery.GroupBy.Accept(columnRewritter)
	}
	if subquery.Having != nil {
		subquery.Having.Accept(columnRewritter)
	}
	if subquery.OrderBy != nil {
		subquery.OrderBy.Accept(columnRewritter)
	}

	return nil
}

// 处理from table和join on部分
// 主要是改写table ExprNode, 并找到路由条件
func handleSubqueryTableRefs(p *TableAliasStmtInfo, stmt *ast.SelectStmt) error {
	tableRefs := stmt.From
	if tableRefs == nil {
		return nil
	}

	join := tableRefs.TableRefs
	if join == nil {
		return nil
	}

	if err := handleSubqueryJoin(p, join); err != nil {
		return fmt.Errorf("handleSubqueryTableRefs error: %v", err)
	}
	return nil
}

func handleSubqueryJoin(p *TableAliasStmtInfo, join *ast.Join) error {
	if err := precheckJoinClause(join); err != nil {
		return fmt.Errorf("precheck Join error: %v", err)
	}

	// 只允许最多两个表的JOIN
	if join.Left != nil {
		switch left := join.Left.(type) {
		case *ast.TableSource:
			// 改写两个表的node
			if err := rewriteSubqueryTableSource(p, left); err != nil {
				return fmt.Errorf("rewrite left TableSource error: %v", err)
			}
		case *ast.Join:
			if err := handleSubqueryJoin(p, left); err != nil {
				return fmt.Errorf("handle nested left Join error: %v", err)
			}
		default:
			return fmt.Errorf("invalid left Join type: %T", join.Left)
		}
	}
	if join.Right != nil {
		right, ok := join.Right.(*ast.TableSource)
		if !ok {
			return fmt.Errorf("right is not TableSource, type: %T", join.Right)
		}

		if err := rewriteSubqueryTableSource(p, right); err != nil {
			return fmt.Errorf("rewrite right TableSource error: %v", err)
		}
	}

	// 只改写表名, 不计算路由
	if join.On != nil {
		rewritter := NewSubqueryColumnNameRewriteVisitor(p)
		join.On.Accept(rewritter)
	}

	return nil
}

// gaea规定在子查询的FROM表名中不能再出现子查询
func rewriteSubqueryTableSource(p *TableAliasStmtInfo, tableSource *ast.TableSource) error {
	switch tableSource.Source.(type) {
	case *ast.TableName:
		return rewriteSubqueryTableNameInTableSource(p, tableSource)
	case *ast.SelectStmt:
		return fmt.Errorf("cannot handle subquery in subquery")
	default:
		return fmt.Errorf("field Source cannot handle, type: %T", tableSource.Source)
	}
}

func rewriteSubqueryTableNameInTableSource(p *TableAliasStmtInfo, tableSource *ast.TableSource) error {
	tableName, ok := tableSource.Source.(*ast.TableName)
	if !ok {
		return fmt.Errorf("field Source is not type of TableName, type: %T", tableSource.Source)
	}

	// 不记录子查询的表名alias
	rule, need, err := NeedCreateTableNameDecorator(p, tableName, "")
	if err != nil {
		return fmt.Errorf("check NeedCreateTableNameDecorator error: %v", err)
	}

	if !need {
		return nil
	}

	// 这是一个分片表或关联表, 创建一个TableName的装饰器, 并替换原有节点
	d, err := CreateTableNameDecorator(tableName, rule, p.GetRouteResult())
	if err != nil {
		return fmt.Errorf("create TableNameDecorator error: %v", err)
	}
	tableSource.Source = d
	return nil
}
