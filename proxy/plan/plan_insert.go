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

	"github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/proxy/sequence"
	"github.com/XiaoMi/Gaea/util"
)

// InsertPlan is the plan for insert statement
type InsertPlan struct {
	basePlan
	*StmtInfo

	stmt *ast.InsertStmt

	table               string
	isAssignmentMode    bool
	shardingColumnIndex int

	sequences *sequence.SequenceManager

	sqls map[string]map[string][]string
}

// NewInsertPlan constructor of InsertPlan
func NewInsertPlan(db string, sql string, r *router.Router, seq *sequence.SequenceManager) *InsertPlan {
	return &InsertPlan{
		StmtInfo:            NewStmtInfo(db, sql, r),
		shardingColumnIndex: -1,
		sequences:           seq,
	}
}

// GetStmt return InsertStmt
func (s *InsertPlan) GetStmt() *ast.InsertStmt {
	return s.stmt
}

// HandleInsertStmt build a InsertPlan
func HandleInsertStmt(p *InsertPlan, stmt *ast.InsertStmt) error {
	p.stmt = stmt

	if err := precheckInsertStmt(p); err != nil {
		return err
	}

	// 处理全局表成功时会触发fastReturn
	fastReturn, err := handleInsertTableRefs(p)
	if err != nil {
		return fmt.Errorf("handleInsertTableRefs error: %v", err)
	}
	if fastReturn {
		return nil
	}

	if err := handleInsertGlobalSequenceValue(p); err != nil {
		return fmt.Errorf("handleInsertGlobalSequenceValue error: %v", err)
	}

	if err := handleInsertColumnNames(p); err != nil {
		return fmt.Errorf("handleInsertColumnNames error: %v", err)
	}

	if err := handleInsertOnDuplicate(p); err != nil {
		return fmt.Errorf("handleInsertOnDuplicate error: %v", err)
	}

	if err := handleInsertValues(p); err != nil {
		return fmt.Errorf("handleInsertValues error: %v", err)
	}

	sqls, err := generateShardingSQLs(p.stmt, p.result, p.router)
	if err != nil {
		log.Warn("generate insert sql failed, %v", err)
		return err
	}

	p.sqls = sqls

	return nil
}

func precheckInsertStmt(p *InsertPlan) error {
	stmt := p.stmt
	// doesn't support insert into select...
	if stmt.Select != nil {
		return errors.ErrSelectInInsert
	}

	// INSERT INTO tbl SET col=val, ...
	if len(stmt.Setlist) != 0 {
		p.isAssignmentMode = true
		return nil
	}

	if len(stmt.Columns) == 0 {
		return errors.ErrIRNoColumns
	}

	values := stmt.Lists[0]
	if len(stmt.Columns) != len(values) {
		return fmt.Errorf("column count doesn't match value count")
	}

	return nil
}

func handleInsertTableRefs(p *InsertPlan) (fastReturn bool, err error) {
	if p.stmt.Table.TableRefs.Right != nil {
		return false, fmt.Errorf("have multi tables in insert")
	}
	tableSource, ok := p.stmt.Table.TableRefs.Left.(*ast.TableSource)
	if !ok {
		return false, fmt.Errorf("not a table source")
	}
	tableName := tableSource.Source.(*ast.TableName)
	p.table = tableName.Name.L

	rule, need, err := NeedCreateTableNameDecoratorWithoutAlias(p.StmtInfo, tableName)
	if err != nil {
		return false, fmt.Errorf("check table name need to decorate error: %v", err)
	}

	if !need {
		// 如果不需要装饰, 不应该走到分表逻辑, 直接报错
		return false, fmt.Errorf("not a sharding table")
	}

	decorator, err := CreateTableNameDecorator(tableName, rule, p.GetRouteResult())
	if err != nil {
		return false, fmt.Errorf("create table name decorator error: %v", err)
	}

	tableSource.Source = decorator

	// 如果是全局表, 则将记录写入所有分片
	if rule.GetType() == router.GlobalTableRuleType {
		p.result.db = rule.GetDB()
		p.result.table = rule.GetTable()
		p.result.indexes = rule.GetSubTableIndexes()
		sqls, err := generateShardingSQLs(p.stmt, p.result, p.router)
		if err != nil {
			return false, fmt.Errorf("generate global table insert sql error: %v", err)
		}
		p.sqls = sqls
		return true, nil
	}

	return false, nil
}

func handleInsertColumnNames(p *InsertPlan) error {
	if p.isAssignmentMode {
		// INSERT INTO tbl SET col = val, ...
		for i, assignment := range p.stmt.Setlist {
			col := assignment.Column
			removeSchemaAndTableInfoInColumnName(col)
			columnName := col.Name.L
			rule := p.tableRules[p.table]
			if columnName == rule.GetShardingColumn() {
				p.shardingColumnIndex = i
			}
		}
	} else {
		// INSERT INTO tbl (col, ...) VALUES (val, ...)
		for i, col := range p.stmt.Columns {
			removeSchemaAndTableInfoInColumnName(col)
			columnName := col.Name.L
			rule := p.tableRules[p.table]
			if columnName == rule.GetShardingColumn() {
				p.shardingColumnIndex = i
			}
		}
	}
	if p.shardingColumnIndex == -1 {
		return fmt.Errorf("sharding column not found")
	}
	return nil
}

// 只有一个表, 直接去掉DB名和表名, 就不需要加装饰器了
func removeSchemaAndTableInfoInColumnName(column *ast.ColumnName) {
	column.Schema.O = ""
	column.Schema.L = ""
	column.Table.O = ""
	column.Table.L = ""
}

// TODO: refactor
func handleInsertValues(p *InsertPlan) error {
	// assignment mode
	if p.isAssignmentMode {
		valueItem := p.stmt.Setlist[p.shardingColumnIndex].Expr
		switch x := valueItem.(type) {
		case *driver.ValueExpr:
			v, err := util.GetValueExprResult(x)
			if err != nil {
				return fmt.Errorf("get value expr result failed, %v", err)
			}
			if v == nil {
				return fmt.Errorf("sharding value cannot be null")
			}
			routeIdx, err := p.tableRules[p.table].FindTableIndex(v)
			if err != nil {
				return fmt.Errorf("find table index error: %v", err)
			}
			p.result.Inter([]int{routeIdx})
		}
		return nil
	}

	// not assignment mode
	for _, valueList := range p.stmt.Lists {
		valueItem := valueList[p.shardingColumnIndex]
		switch x := valueItem.(type) {
		case *driver.ValueExpr:
			v, err := util.GetValueExprResult(x)
			if err != nil {
				return fmt.Errorf("get value expr result failed, %v", err)
			}
			if v == nil {
				return fmt.Errorf("sharding value cannot be null")
			}
			routeIdx, err := p.tableRules[p.table].FindTableIndex(v)
			if err != nil {
				return fmt.Errorf("find table index error: %v", err)
			}
			p.result.Inter([]int{routeIdx})
		}
	}
	if len(p.result.GetShardIndexes()) == 0 {
		return fmt.Errorf("batch insert has cross slice values or no route found")
	}
	return nil
}

// check on duplicate key
// 不管分片表的配置信息, 只要在OnDuplicate出现分片列, 就返回错误
// 去掉ColumnName中的DB名和表名
func handleInsertOnDuplicate(p *InsertPlan) error {
	if p.stmt.OnDuplicate == nil {
		return nil
	}

	shardingColumnName := p.tableRules[p.table].GetShardingColumn()
	for _, a := range p.stmt.OnDuplicate {
		if a.Column.Name.L == shardingColumnName {
			return errors.ErrUpdateKey
		}
		removeSchemaAndTableInfoInColumnName(a.Column)
	}

	return nil
}

// 处理全局序列号, 目前一条SQL中只允许一个列使用全局序列号
func handleInsertGlobalSequenceValue(p *InsertPlan) error {
	seq, ok := p.sequences.GetSequence(p.db, p.table)
	if !ok {
		return nil
	}
	pkName := seq.GetPKName()

	// not assignment mode
	if p.isAssignmentMode {
		for _, assignment := range p.stmt.Setlist {
			columnName := assignment.Column.Name.L
			if columnName == pkName {
				if x, ok := assignment.Expr.(*ast.FuncCallExpr); ok {
					if x.FnName.L == "nextval" {
						id, err := seq.NextSeq()
						if err != nil {
							return fmt.Errorf("get next seq error: %v", err)
						}
						assignment.Expr = ast.NewValueExpr(id)
						break
					}
				}
			}
		}
		return nil
	}

	// not assignment mode
	var seqIndex = -1
	for i, column := range p.stmt.Columns {
		columnName := column.Name.L
		if columnName == pkName {
			seqIndex = i
			break
		}
	}

	// global sequence column not found
	if seqIndex == -1 {
		return nil
	}

	for _, valueList := range p.stmt.Lists {
		if x, ok := valueList[seqIndex].(*ast.FuncCallExpr); ok {
			if x.FnName.L == "nextval" {
				id, err := seq.NextSeq()
				if err != nil {
					return fmt.Errorf("get next seq error: %v", err)
				}
				valueList[seqIndex] = ast.NewValueExpr(id)
			}
		}
	}

	return nil
}

// ExecuteIn implement Plan
func (s *InsertPlan) ExecuteIn(reqCtx *util.RequestContext, sess Executor) (*mysql.Result, error) {
	rs, err := sess.ExecuteSQLs(reqCtx, s.sqls)
	if err != nil {
		return nil, fmt.Errorf("execute in InsertPlan error: %v", err)
	}

	r, err := MergeExecResult(rs)
	if err != nil {
		return nil, err
	}

	if r.InsertID != 0 {
		sess.SetLastInsertID(r.InsertID)
	}

	return r, nil
}
