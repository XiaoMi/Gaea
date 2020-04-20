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

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/util"
)

// UpdatePlan is the plan for update statement
type UpdatePlan struct {
	basePlan
	*TableAliasStmtInfo

	stmt *ast.UpdateStmt
	sqls map[string]map[string][]string
}

// NewUpdatePlan constructor of UpdatePlan
func NewUpdatePlan(stmt *ast.UpdateStmt, db, sql string, r *router.Router) *UpdatePlan {
	return &UpdatePlan{
		TableAliasStmtInfo: NewTableAliasStmtInfo(db, sql, r),
		stmt:               stmt,
	}
}

// ExecuteIn implement Plan
func (s *UpdatePlan) ExecuteIn(reqCtx *util.RequestContext, sess Executor) (*mysql.Result, error) {
	sqls := s.sqls
	if sqls == nil {
		return nil, fmt.Errorf("SQL has not generated")
	}

	if len(sqls) == 0 {
		return nil, nil
	}

	rs, err := sess.ExecuteSQLs(reqCtx, sqls)
	if err != nil {
		return nil, fmt.Errorf("execute in UpdatePlan error: %v", err)
	}

	r, err := MergeExecResult(rs)

	if err != nil {
		return nil, fmt.Errorf("merge update result error: %v", err)
	}

	return r, nil
}

// HandleUpdatePlan build a UpdatePlan
func HandleUpdatePlan(p *UpdatePlan) error {
	if err := handleUpdateTableRefs(p); err != nil {
		return fmt.Errorf("handle From error: %v", err)
	}

	if err := handleUpdateAssignmentList(p); err != nil {
		return fmt.Errorf("handle assignment list error: %v", err)
	}

	if err := handleUpdateWhere(p); err != nil {
		return fmt.Errorf("handle Where error: %v", err)
	}

	if err := handleUpdateOrderBy(p); err != nil {
		return fmt.Errorf("handle OrderBy error: %v", err)
	}

	// Limit clause does not need to handle

	// handle global table
	if err := postHandleGlobalTableRouteResultInModify(p.StmtInfo); err != nil {
		return fmt.Errorf("post handle global table error: %v", err)
	}

	sqls, err := generateShardingSQLs(p.stmt, p.GetRouteResult(), p.router)
	if err != nil {
		return fmt.Errorf("generate sqls error: %v", err)
	}

	p.sqls = sqls
	return nil
}

func handleUpdateTableRefs(p *UpdatePlan) error {
	tableRefs := p.stmt.TableRefs
	if tableRefs == nil {
		return nil
	}

	join := tableRefs.TableRefs
	if join == nil {
		return nil
	}

	if join.Right != nil {
		return fmt.Errorf("does not support update multiple tables in sharding")
	}

	return handleJoin(p.TableAliasStmtInfo, join)
}

func handleUpdateWhere(p *UpdatePlan) error {
	stmt := p.stmt
	if stmt.Where == nil {
		return nil
	}

	has, result, decorator, err := handleComparisonExpr(p.TableAliasStmtInfo, stmt.Where)
	if err != nil {
		return fmt.Errorf("rewrite Where error: %v", err)
	}
	if has {
		p.GetRouteResult().Inter(result)
	}
	stmt.Where = decorator
	return nil
}

func handleUpdateOrderBy(p *UpdatePlan) error {
	order := p.stmt.Order
	if order == nil {
		return nil
	}

	for _, item := range order.Items {
		columnExpr, ok := item.Expr.(*ast.ColumnNameExpr)
		if !ok {
			return fmt.Errorf("ByItem.Expr is not a ColumnNameExpr")
		}

		rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInField(p.TableAliasStmtInfo, columnExpr)
		if err != nil {
			return err
		}

		if need {
			decorator := CreateColumnNameExprDecorator(columnExpr, rule, isAlias, p.GetRouteResult())
			item.Expr = decorator
		}
	}

	return nil
}

// TODO: Assignment直接引用ColumnName, 不能做表名的装饰器. 采用的解决办法是UPDATE只支持一个表, 然后把DB名和表名去掉.
func handleUpdateAssignmentList(p *UpdatePlan) error {
	l := p.stmt.List
	for _, assignment := range l {
		r, need, _, err := needCreateColumnNameDecorator(p.TableAliasStmtInfo, assignment.Column)
		if err != nil {
			return err
		}

		if need && r.GetShardingColumn() == assignment.Column.Name.L {
			return fmt.Errorf("cannot update shard column value")
		}
		removeSchemaAndTableInfoInColumnName(assignment.Column)
	}
	return nil
}
