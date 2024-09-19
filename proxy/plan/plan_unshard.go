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
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/proxy/router"
	"strings"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
	"github.com/XiaoMi/Gaea/util"
)

// UnshardPlan is the plan for unshard statement
type UnshardPlan struct {
	basePlan

	db     string
	phyDBs map[string]string
	sql    string
	stmt   ast.StmtNode
}

// SelectLastInsertIDPlan is the plan for SELECT LAST_INSERT_ID()
// TODO: fix below
// https://dev.mysql.com/doc/refman/5.6/en/information-functions.html#function_last-insert-id
// The value of LAST_INSERT_ID() is not changed if you set the AUTO_INCREMENT column of a row
// to a non-“magic” value (that is, a value that is not NULL and not 0).
type SelectLastInsertIDPlan struct {
	basePlan
	asName string
}

type SetPlan struct {
	basePlan
	sql  string
	stmt ast.StmtNode
}

type IgnorePlan struct {
	basePlan
}

// IsSelectLastInsertIDStmt check if the statement is SELECT LAST_INSERT_ID()
func IsSelectLastInsertIDStmt(stmt ast.StmtNode) bool {
	s, ok := stmt.(*ast.SelectStmt)
	if !ok {
		return false
	}

	if len(s.Fields.Fields) != 1 {
		return false
	}

	if s.From != nil || s.Where != nil || s.GroupBy != nil || s.Having != nil || s.OrderBy != nil || s.Limit != nil {
		return false
	}

	f, ok := s.Fields.Fields[0].Expr.(*ast.FuncCallExpr)
	if !ok {
		return false
	}

	return f.FnName.L == "last_insert_id"
}

// IsSetStmt check if the statement is set comment
func IsSetStmt(stmt ast.StmtNode) bool {
	_, ok := stmt.(*ast.SetStmt)
	if !ok {
		return false
	}

	return true
}

// CreateUnshardPlan constructor of UnshardPlan
func CreateUnshardPlan(stmt ast.StmtNode, phyDBs map[string]string, db string, tableNames []*ast.TableName) (*UnshardPlan, error) {
	p := &UnshardPlan{
		db:     db,
		phyDBs: phyDBs,
		stmt:   stmt,
	}
	rewriteUnshardTableName(phyDBs, tableNames)
	rsql, err := generateUnshardingSQL(stmt)
	if err != nil {
		return nil, fmt.Errorf("generate unshardPlan SQL error: %v", err)
	}
	p.sql = rsql
	return p, nil
}

// PreCreateUnshardPlan constructor of UnshardPlan
// if db not in phyDBs, return false
// if db is empty, like `select user()` or `select * from t`, will not return true and use session db cause error may return from backend mysql
func PreCreateUnshardPlan(sql string, phyDBs map[string]string, db string) (*UnshardPlan, error) {
	if phyDB, ok := phyDBs[db]; ok {
		if db != phyDB {
			return nil, fmt.Errorf("db name not match, db: %s, phyDB: %s", db, phyDB)
		}
	}

	return &UnshardPlan{
		db:     db,
		phyDBs: phyDBs,
		sql:    sql,
	}, nil
}

func rewriteUnshardTableName(phyDBs map[string]string, tableNames []*ast.TableName) {
	for _, tableName := range tableNames {
		if phyDB, ok := phyDBs[tableName.Schema.String()]; ok {
			tableName.Schema.O = phyDB
			tableName.Schema.L = strings.ToLower(phyDB)
		}
	}
}

func generateUnshardingSQL(stmt ast.StmtNode) (string, error) {
	s := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, s)
	_ = stmt.Restore(ctx)
	return s.String(), nil
}

// CreateSelectLastInsertIDPlan constructor of SelectLastInsertIDPlan
func CreateSelectLastInsertIDPlan(stmt *ast.SelectStmt) *SelectLastInsertIDPlan {
	asName := ""
	if len(stmt.Fields.Fields) > 0 && stmt.Fields.Fields[0].AsName.String() != "" {
		asName = stmt.Fields.Fields[0].AsName.String()
	}
	return &SelectLastInsertIDPlan{asName: asName}
}

// CreateSetPlan constructor of SetPlan
func CreateSetPlan(sql string, stmt ast.StmtNode) *SetPlan {
	return &SetPlan{sql: sql,
		stmt: stmt}
}

// CreateIgnoreSetPlan constructor of SetPlan
func CreateIgnorePlan() *IgnorePlan {
	return &IgnorePlan{}
}

// ExecuteIn implement Plan
func (p *UnshardPlan) ExecuteIn(reqCtx *util.RequestContext, se Executor) (*mysql.Result, error) {
	r, err := se.ExecuteSQL(reqCtx, reqCtx.GetDefaultSlice(), p.db, p.sql)
	if err != nil {
		return nil, err
	}

	// set last insert id to session
	if r.InsertID != 0 {
		se.SetLastInsertID(r.InsertID)
	}

	return r, nil
}

// ExecuteIn implement Plan
func (p *SelectLastInsertIDPlan) ExecuteIn(reqCtx *util.RequestContext, se Executor) (*mysql.Result, error) {
	r := createLastInsertIDResult(se.GetLastInsertID(), p.asName)
	return r, nil
}

// ExecuteIn implement Plan
func (p *SetPlan) ExecuteIn(reqCtx *util.RequestContext, se Executor) (*mysql.Result, error) {
	if stmt, ok := p.stmt.(*ast.SetStmt); ok {
		se.HandleSet(reqCtx, p.sql, stmt)
	}

	return nil, nil
}

// ExecuteIn implement Plan
func (p *IgnorePlan) ExecuteIn(reqCtx *util.RequestContext, se Executor) (*mysql.Result, error) {
	return mysql.ResultPool.GetWithoutResultSet(), nil
}

func createLastInsertIDResult(lastInsertID uint64, asName string) *mysql.Result {
	name := "last_insert_id()"
	if asName != "" {
		name = asName
	}
	var column = 1
	var rows [][]uint64
	var names = []string{
		name,
	}

	rows = append(rows, []uint64{lastInsertID})

	var values = make([][]interface{}, len(rows))
	for i := range rows {
		values[i] = make([]interface{}, column)
		for j := range rows[i] {
			values[i][j] = rows[i][j]
		}
	}

	r, _ := mysql.BuildResultset(nil, names, values)
	ret := mysql.ResultPool.Get()
	ret.Resultset = r

	return ret
}

func CheckUnshardBase(tokenId int, tokens []string, rt *router.Router, db string) (string, bool) {
	ruleDB := db
	tokensLen := len(tokens)
	for i := 0; i < tokensLen; i++ {
		if strings.ToLower(tokens[i]) != mysql.ParseTokenIdStrMap[tokenId] {
			continue
		}
		if i+1 >= tokensLen {
			continue
		}
		// select: select col1 from db.t where...
		// delete: delete from db.t where...
		dbName, tableName := parser.GetDBTable(tokens[i+1])
		//if the token[i+1] like this: db.test_shard_hash
		if dbName != "" {
			ruleDB = dbName
		}
		// if table in shard rule, is shard plan
		if rt.GetRule(ruleDB, tableName) != rt.GetDefaultRule() {
			return ruleDB, false
		}
	}
	return ruleDB, true
}

func CheckUnshardInsert(tokens []string, rt *router.Router, db string) (string, bool) {
	ruleDB := db
	tokensLen := len(tokens)
	for i := 0; i < tokensLen; i++ {
		if strings.ToLower(tokens[i]) != mysql.ParseTokenIdStrMap[mysql.TkIdInsert] {
			continue
		}
		if i+1 >= tokensLen {
			continue
		}
		// insert: insert into db.t(col1) values...
		// replace: replace into db.t(col1) values...
		dbName, tableName := parser.GetInsertDBTable(tokens[i+1])
		//if the token[i+1] like this: db.test_shard_hash
		if dbName != "" {
			ruleDB = dbName
		}
		// if table in shard rule, is shard plan
		if rt.GetRule(ruleDB, tableName) != rt.GetDefaultRule() {
			return ruleDB, false
		}
	}
	return ruleDB, true
}

func CheckUnshardUpdate(tokens []string, rt *router.Router, db string) (string, bool) {
	ruleDB := db
	tokensLen := len(tokens)
	for i := 1; i < tokensLen; i++ {
		if strings.ToLower(tokens[i]) != mysql.ParseTokenIdStrMap[mysql.TkIdUpdate] {
			continue
		}
		if i+1 >= tokensLen {
			continue
		}
		// update: update t set a=1 where ...
		dbName, tableName := parser.GetDBTable(tokens[i-1])
		if dbName != "" {
			ruleDB = dbName
		}
		// if table in shard rule, is shard plan
		if rt.GetRule(ruleDB, tableName) != rt.GetDefaultRule() {
			return ruleDB, false
		}
	}
	return ruleDB, true
}
