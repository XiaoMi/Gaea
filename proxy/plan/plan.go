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
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/proxy/sequence"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/hack"
)

// type check
var _ Plan = &UnshardPlan{}
var _ Plan = &SelectPlan{}
var _ Plan = &DeletePlan{}
var _ Plan = &UpdatePlan{}
var _ Plan = &InsertPlan{}
var _ Plan = &SelectLastInsertIDPlan{}

// Plan is a interface for select/insert etc.
type Plan interface {
	ExecuteIn(*util.RequestContext, Executor) (*mysql.Result, error)

	// only for cache
	Size() int
}

// Executor TODO: move to package executor
type Executor interface {

	// 执行分片或非分片单条SQL
	ExecuteSQL(ctx *util.RequestContext, slice, db, sql string) (*mysql.Result, error)

	// 执行分片SQL
	ExecuteSQLs(*util.RequestContext, map[string]map[string][]string) ([]*mysql.Result, error)

	// 用于执行INSERT时设置last insert id
	SetLastInsertID(uint64)

	GetLastInsertID() uint64
}

// Checker 用于检查SelectStmt是不是分表的Visitor, 以及是否包含DB信息
type Checker struct {
	db            string
	router        *router.Router
	hasShardTable bool // 是否包含分片表
	dbInvalid     bool // SQL是否No database selected
	tableNames    []*ast.TableName
}

// NewChecker db为USE db中设置的DB名. 如果没有执行USE db, 则为空字符串
func NewChecker(db string, router *router.Router) *Checker {
	return &Checker{
		db:            db,
		router:        router,
		hasShardTable: false,
		dbInvalid:     false,
	}
}

func (s *Checker) GetUnshardTableNames() []*ast.TableName {
	return s.tableNames
}

// IsDatabaseInvalid 判断执行计划中是否包含db信息, 如果不包含, 且又含有表名, 则是一个错的执行计划, 应该返回以下错误:
// ERROR 1046 (3D000): No database selected
func (s *Checker) IsDatabaseInvalid() bool {
	return s.dbInvalid
}

// IsShard if is shard table
func (s *Checker) IsShard() bool {
	return s.hasShardTable
}

// Enter for node visit
func (s *Checker) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if s.hasShardTable {
		return n, true
	}
	switch nn := n.(type) {
	case *ast.TableName:
		if s.isTableNameDatabaseInvalid(nn) {
			s.dbInvalid = true
			return n, true
		}
		has := s.hasShardTableInTableName(nn)
		if has {
			s.hasShardTable = true
			return n, true
		}
		s.tableNames = append(s.tableNames, nn)
	}
	return n, false
}

// Leave for node visit
func (s *Checker) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, !s.dbInvalid && !s.hasShardTable
}

// 如果ast.TableName不带DB名, 且Session未设置DB, 则是不允许的SQL, 应该返回No database selected
func (s *Checker) isTableNameDatabaseInvalid(n *ast.TableName) bool {
	return s.db == "" && n.Schema.L == ""
}

func (s *Checker) hasShardTableInTableName(n *ast.TableName) bool {
	db := n.Schema.L
	if db == "" {
		db = s.db
	}
	table := n.Name.L
	_, ok := s.router.GetShardRule(db, table)
	return ok
}

func (s *Checker) hasShardTableInColumnName(n *ast.ColumnName) bool {
	db := n.Schema.L
	if db == "" {
		db = s.db
	}
	table := n.Table.L
	_, ok := s.router.GetShardRule(db, table)
	return ok
}

type basePlan struct{}

func (*basePlan) Size() int {
	return 1
}

// StmtInfo 各种Plan的一些公共属性
type StmtInfo struct {
	db               string // session db
	sql              string // origin sql
	router           *router.Router
	tableRules       map[string]router.Rule // key = table name, value = router.Rule, 记录使用到的分片表
	globalTableRules map[string]router.Rule // 记录使用到的全局表
	result           *RouteResult
}

// TableAliasStmtInfo 使用到表别名, 且依赖表别名做路由计算的StmtNode, 目前包括UPDATE, SELECT
// INSERT也可以使用表别名, 但是由于只存在一个表, 可以直接去掉, 因此不需要.
type TableAliasStmtInfo struct {
	*StmtInfo
	tableAlias map[string]string // key = table alias, value = table
	hintPhyDB  string            // 记录mycat分片时DATABASE()函数指定的物理DB名
}

// BuildPlan build plan for ast
func BuildPlan(stmt ast.StmtNode, phyDBs map[string]string, db, sql string, router *router.Router, seq *sequence.SequenceManager) (Plan, error) {
	if IsSelectLastInsertIDStmt(stmt) {
		return CreateSelectLastInsertIDPlan(), nil
	}

	if estmt, ok := stmt.(*ast.ExplainStmt); ok {
		return buildExplainPlan(estmt, phyDBs, db, sql, router, seq)
	}

	checker := NewChecker(db, router)
	stmt.Accept(checker)

	if checker.IsDatabaseInvalid() {
		return nil, fmt.Errorf("no database selected") // TODO: return standard MySQL error
	}

	if checker.IsShard() {
		return buildShardPlan(stmt, db, sql, router, seq)
	}
	return CreateUnshardPlan(stmt, phyDBs, db, checker.GetUnshardTableNames())
}

func buildShardPlan(stmt ast.StmtNode, db string, sql string, router *router.Router, seq *sequence.SequenceManager) (Plan, error) {
	switch s := stmt.(type) {
	case *ast.SelectStmt:
		plan := NewSelectPlan(db, sql, router)
		if err := HandleSelectStmt(plan, s); err != nil {
			return nil, err
		}
		return plan, nil
	case *ast.InsertStmt:
		// InsertStmt contains REPLACE statement
		plan := NewInsertPlan(db, sql, router, seq)
		if err := HandleInsertStmt(plan, s); err != nil {
			return nil, err
		}
		return plan, nil
	case *ast.UpdateStmt:
		plan := NewUpdatePlan(s, db, sql, router)
		if err := HandleUpdatePlan(plan); err != nil {
			return nil, err
		}
		return plan, nil
	case *ast.DeleteStmt:
		plan := NewDeletePlan(s, db, sql, router)
		if err := HandleDeletePlan(plan); err != nil {
			return nil, err
		}
		return plan, nil
	default:
		return nil, fmt.Errorf("stmt type does not support shard now")
	}
}

// NewStmtInfo constructor of StmtInfo
func NewStmtInfo(db string, sql string, r *router.Router) *StmtInfo {
	return &StmtInfo{
		db:               db,
		sql:              sql,
		router:           r,
		tableRules:       make(map[string]router.Rule),
		globalTableRules: make(map[string]router.Rule),
		result:           NewRouteResult("", "", nil), // nil route result
	}
}

// NewTableAliasStmtInfo means table alias StmtInfo
func NewTableAliasStmtInfo(db string, sql string, r *router.Router) *TableAliasStmtInfo {
	return &TableAliasStmtInfo{
		StmtInfo:   NewStmtInfo(db, sql, r),
		tableAlias: make(map[string]string),
	}
}

// GetRouteResult get route result
func (s *StmtInfo) GetRouteResult() *RouteResult {
	return s.result
}

func (s *StmtInfo) checkAndGetDB(db string) (string, error) {
	if db != "" && db != s.db {
		return "", fmt.Errorf("db not match")
	}
	return s.db, nil
}

// RecordShardTable 将表信息记录到StmtInfo中, 并返回表信息对应的路由规则
func (s *StmtInfo) RecordShardTable(db, table string) (router.Rule, error) {
	rule, err := s.getShardRule(db, table)
	if err != nil {
		return nil, fmt.Errorf("get shard rule error, db: %s, table: %s, err: %v", db, table, err)
	}

	if err := s.checkStmtRouteResult(rule); err != nil {
		return nil, fmt.Errorf("check route result error, db: %s, table: %s, err: %v", db, table, err)
	}

	return rule, nil
}

// 根据db和table获取Rule
// 如果只传table, 则使用session db.
func (s *StmtInfo) getShardRule(db, table string) (router.Rule, error) {
	validDB, err := s.checkAndGetDB(db)
	if err != nil {
		return nil, err
	}

	rule, ok := s.router.GetShardRule(validDB, table) // 这里一定是ShardingRule, 不会是DefaultRule
	if !ok {
		return nil, fmt.Errorf("rule not found")
	}

	if rule.GetType() == router.GlobalTableRuleType {
		s.globalTableRules[table] = rule
	} else {
		s.tableRules[table] = rule // 记录已经使用到的rule
	}
	return rule, nil
}

// 检查路由规则与现有RouteResult是否一致
// 一致的标准: 与RouteResult的db, table一致
func (s *StmtInfo) checkStmtRouteResult(rule router.Rule) error {
	// 如果是全局表, 不需要检查路由规则是否一致, 只记录该规则, 直接返回即可
	if rule.GetType() == router.GlobalTableRuleType {
		return nil
	}

	db := rule.GetDB()
	var table string
	if linkedRule, ok := rule.(*router.LinkedRule); ok {
		table = linkedRule.GetParentTable()
	} else {
		table = rule.GetTable()
	}

	if s.result.db == "" && s.result.table == "" {
		s.result.db = db
		s.result.table = table
		s.result.indexes = rule.GetSubTableIndexes()
	} else {
		if err := s.result.Check(db, table); err != nil {
			return fmt.Errorf("check db and table error: %v", err)
		}
	}

	return nil
}

// 用于WHERE条件或JOIN ON条件中, 只存在列名时, 查找对应的路由规则
func (s *StmtInfo) getSettedRuleByColumnName(column string) (router.Rule, bool, error) {
	var columnExistsInShardingTables int // 记录分片表名出现在分片表中分片列的次数
	var ret router.Rule
	for _, r := range s.tableRules {
		if r.GetShardingColumn() == column {
			columnExistsInShardingTables++
			ret = r
		}
	}

	if columnExistsInShardingTables > 1 {
		return nil, false, fmt.Errorf("column %s is ambiguous for sharding", column)
	}

	return ret, ret != nil, nil
}

// 处理SELECT只含有全局表的情况
// 这种情况只路由到默认分片
// 如果有多个全局表, 则只取第一个全局表的配置, 因此需要业务上保证这些全局表的配置是一致的.
func postHandleGlobalTableRouteResultInQuery(p *StmtInfo) error {
	if len(p.tableRules) == 0 && len(p.globalTableRules) != 0 {
		var tableName string
		var rule router.Rule
		for t, r := range p.globalTableRules {
			tableName = t
			rule = r
			break
		}
		p.result.db = rule.GetDB()
		p.result.table = tableName
		p.result.indexes = []int{0} // 全局表SELECT只取默认分片
	}
	return nil
}

// 处理UPDATE, DELETE只含有全局表的情况
// 这种情况只路由到默认分片
// 如果有多个全局表, 则只取第一个全局表的配置, 因此需要业务上保证这些全局表的配置是一致的.
func postHandleGlobalTableRouteResultInModify(p *StmtInfo) error {
	if len(p.tableRules) == 0 && len(p.globalTableRules) != 0 {
		var tableName string
		var rule router.Rule
		for t, r := range p.globalTableRules {
			tableName = t
			rule = r
			break
		}
		p.result.db = rule.GetDB()
		p.result.table = tableName
		p.result.indexes = rule.GetSubTableIndexes()
	}
	return nil
}

// RecordSubqueryTableAlias 记录表名位置的子查询的别名, 便于后续处理
// 返回已存在Rule的第一个 (任意一个即可)
// 限制: 子查询中的表对应的路由规则必须与外层查询相关联, 或者为全局表
func (t *TableAliasStmtInfo) RecordSubqueryTableAlias(alias string) (router.Rule, error) {
	if alias == "" {
		return nil, fmt.Errorf("subquery table alias is nil")
	}

	if len(t.tableRules) == 0 {
		return nil, fmt.Errorf("no explicit table exist except subquery")
	}

	table := "gaea_subquery_" + alias
	if err := t.setTableAlias(table, alias); err != nil {
		return nil, fmt.Errorf("set subquery table alias error: %v", err)
	}

	var rule router.Rule
	for _, r := range t.tableRules {
		rule = r
		break
	}

	t.tableRules[table] = rule
	return rule, nil
}

// GetSettedRuleFromColumnInfo 用于WHERE条件或JOIN ON条件中, 查找列名对应的路由规则
func (t *TableAliasStmtInfo) GetSettedRuleFromColumnInfo(db, table, column string) (router.Rule, bool, bool, error) {
	if db == "" && table == "" {
		rule, need, err := t.getSettedRuleByColumnName(column)
		return rule, need, false, err
	}

	rule, isAlias, err := t.getSettedRuleFromTable(db, table)
	return rule, rule != nil, isAlias, err
}

// 用于WHERE条件或JOIN ON条件中, 只存在列名时, 查找对应的路由规则
func (t *TableAliasStmtInfo) getSettedRuleByColumnName(column string) (router.Rule, bool, error) {
	var columnExistsInShardingTables int // 记录分片表名出现在分片表中分片列的次数
	var ret router.Rule
	for _, r := range t.tableRules {
		if r.GetShardingColumn() == column {
			columnExistsInShardingTables++
			ret = r
		}
	}

	if columnExistsInShardingTables > 1 {
		return nil, false, fmt.Errorf("column %s is ambiguous for sharding", column)
	}

	return ret, ret != nil, nil
}

// 获取FROM TABLE列表中的表数据
// 用于FieldList和Where条件中列名的判断
func (t *TableAliasStmtInfo) getSettedRuleFromTable(db, table string) (router.Rule, bool, error) {
	_, err := t.checkAndGetDB(db)
	if err != nil {
		return nil, false, err
	}
	if rule, ok := t.tableRules[table]; ok {
		return rule, false, nil
	}

	if rule, ok := t.globalTableRules[table]; ok {
		return rule, false, nil
	}

	if originTable, ok := t.getAliasTable(table); ok {
		if rule, ok := t.tableRules[originTable]; ok {
			return rule, true, nil
		}
		if rule, ok := t.globalTableRules[originTable]; ok {
			return rule, true, nil
		}
	}

	return nil, false, fmt.Errorf("rule not found")
}

// RecordShardTable 将表信息记录到StmtInfo中, 并返回表信息对应的路由规则
func (t *TableAliasStmtInfo) RecordShardTable(db, table, alias string) (router.Rule, error) {
	rule, err := t.StmtInfo.RecordShardTable(db, table)
	if err != nil {
		return nil, fmt.Errorf("record shard table error, db: %s, table: %s, alias: %s, err: %v", db, table, alias, err)
	}

	if alias != "" {
		if err := t.setTableAlias(table, alias); err != nil {
			return nil, fmt.Errorf("set table alias error: %v", err)
		}
	}

	return rule, nil
}

func (t *TableAliasStmtInfo) setTableAlias(table, alias string) error {
	// if not set, set without check
	originTable, ok := t.tableAlias[alias]
	if !ok {
		t.tableAlias[alias] = table
		return nil
	}

	if originTable != table {
		return fmt.Errorf("table alias is set but not match, table: %s, originTable: %s", table, originTable)
	}

	// already set, return
	return nil
}

func (t *TableAliasStmtInfo) getAliasTable(alias string) (string, bool) {
	table, ok := t.tableAlias[alias]
	return table, ok
}

// 根据StmtNode和路由信息生成分片SQL
func generateShardingSQLs(stmt ast.StmtNode, result *RouteResult, router *router.Router) (map[string]map[string][]string, error) {
	ret := make(map[string]map[string][]string)

	for result.HasNext() {
		sb := &strings.Builder{}
		ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, sb)
		if err := stmt.Restore(ctx); err != nil {
			return nil, err
		}

		index := result.Next()
		rule, ok := router.GetShardRule(result.db, result.table)
		if !ok {
			return nil, fmt.Errorf("cannot find shard rule, db: %s, table: %s", result.db, result.table)
		}
		sliceIndex := rule.GetSliceIndexFromTableIndex(index)
		sliceName := rule.GetSlice(sliceIndex)
		dbName, _ := rule.GetDatabaseNameByTableIndex(index)
		sliceSQLs, ok := ret[sliceName]
		if !ok {
			sliceSQLs = make(map[string][]string)
			ret[sliceName] = sliceSQLs
		}

		ret[sliceName][dbName] = append(ret[sliceName][dbName], sb.String())
	}

	result.Reset() // must reset the cursor for next call

	return ret, nil
}

// 根据原始SQL生成后端对应slice和db的SQL
func generateSQLResultFromOriginSQL(sql string, result *RouteResult, router *router.Router) (map[string]map[string][]string, error) {
	rule := router.GetRule(result.db, result.table)
	indexes := rule.GetSubTableIndexes()
	ret := make(map[string]map[string][]string)
	for _, index := range indexes {
		sliceIndex := rule.GetSliceIndexFromTableIndex(index)
		sliceName := rule.GetSlice(sliceIndex)
		dbName, _ := rule.GetDatabaseNameByTableIndex(index)
		sliceSQLs, ok := ret[sliceName]
		if !ok {
			sliceSQLs = make(map[string][]string)
			ret[sliceName] = sliceSQLs
		}

		ret[sliceName][dbName] = append(ret[sliceName][dbName], sql)
	}

	return ret, nil
}

// copy from newEmptyResultset
// 注意去掉补充的列
func newEmptyResultset(info *SelectPlan, stmt *ast.SelectStmt) *mysql.Resultset {
	r := new(mysql.Resultset)

	fieldLen := len(stmt.Fields.Fields)
	fieldLen -= info.columnCount - info.originColumnCount

	r.Fields = make([]*mysql.Field, fieldLen)
	for i, expr := range stmt.Fields.Fields {
		r.Fields[i] = &mysql.Field{}
		if expr.WildCard != nil {
			r.Fields[i].Name = []byte("*")
		} else {
			if expr.AsName.String() != "" {
				r.Fields[i].Name = hack.Slice(expr.AsName.String())
				name, _ := parser.NodeToStringWithoutQuote(expr.Expr)
				r.Fields[i].OrgName = hack.Slice(name)
			} else {
				name, _ := parser.NodeToStringWithoutQuote(expr.Expr)
				r.Fields[i].Name = hack.Slice(name)
			}
		}
	}

	r.Values = make([][]interface{}, 0)
	r.RowDatas = make([]mysql.RowData, 0)

	return r
}
