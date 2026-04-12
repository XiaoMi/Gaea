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
	"github.com/XiaoMi/Gaea/parser/opcode"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/util"
)

// SelectPlan is the plan for select statement
type SelectPlan struct {
	basePlan
	*TableAliasStmtInfo

	stmt *ast.SelectStmt

	distinct          bool   // 是否是SELECT DISTINCT
	groupByColumn     []int  // GROUP BY 列索引
	orderByColumn     []int  // ORDER BY 列索引
	orderByDirections []bool // ORDER BY 方向, true: DESC
	originColumnCount int    // 补列前的列长度
	columnCount       int    // 补列后的列长度

	aggregateFuncs map[int]AggregateFuncMerger // key = column index

	havingAggColumnIndex   map[string]int // key = aggregate canonical key, value = merged result column index
	havingAliasColumnIndex map[string]int // key = aggregate alias lower-case, value = merged result column index

	offset int64 // LIMIT offset
	count  int64 // LIMIT count, 未设置则为-1

	havingExpr ast.ExprNode // HAVING表达式, 用于在结果合并后重新应用过滤

	sqls map[string]map[string][]string
}

// NewSelectPlan constructor of SelectPlan
// db is the session db
func NewSelectPlan(db string, sql string, r *router.Router) *SelectPlan {
	return &SelectPlan{
		TableAliasStmtInfo:     NewTableAliasStmtInfo(db, sql, r),
		aggregateFuncs:         make(map[int]AggregateFuncMerger),
		havingAggColumnIndex:   make(map[string]int),
		havingAliasColumnIndex: make(map[string]int),
		offset:                 -1,
		count:                  -1,
	}
}

// ExecuteIn implement Plan
func (s *SelectPlan) ExecuteIn(reqCtx *util.RequestContext, sess Executor) (*mysql.Result, error) {
	sqls := s.GetSQLs()
	if sqls == nil {
		return nil, fmt.Errorf("SQL has not generated")
	}

	if len(sqls) == 0 {
		r := newEmptyResultset(s, s.GetStmt())
		ret := mysql.ResultPool.Get()
		ret.Resultset = r
		return ret, nil
	}

	rs, err := sess.ExecuteSQLs(reqCtx, sqls)
	if err != nil {
		return nil, fmt.Errorf("execute in SelectPlan error: %v", err)
	}
	// fix: 修复全局表或分片表 order by/group by等情况下单分片执行时多一列的问题, 修复由于 limit offset 语句改写导致结果行数不正确问题
	if s.isExecOnSingleNode() && s.noAddColumns() && !s.HasLimit() {
		return rs[0], nil
	}

	r, err := MergeSelectResult(s, s.stmt, rs)
	if err != nil {
		return nil, fmt.Errorf("merge select result error: %v", err)
	}
	return r, nil

}

// GetStmt SelectStmt
func (s *SelectPlan) GetStmt() *ast.SelectStmt {
	return s.stmt
}

func (s *SelectPlan) setAggregateFuncMerger(idx int, merger AggregateFuncMerger) error {
	if _, ok := s.aggregateFuncs[idx]; ok {
		return fmt.Errorf("column %d already set", idx)
	}
	s.aggregateFuncs[idx] = merger
	return nil
}

func (s *SelectPlan) setHavingAggregateColumnIndex(agg *ast.AggregateFuncExpr, idx int) error {
	key, err := havingAggregateKey(agg)
	if err != nil {
		return err
	}
	if s.havingAggColumnIndex == nil {
		s.havingAggColumnIndex = make(map[string]int)
	}
	s.havingAggColumnIndex[key] = idx
	return nil
}

func (s *SelectPlan) setHavingAliasColumnIndex(alias string, idx int) {
	alias = strings.ToLower(strings.TrimSpace(alias))
	if alias == "" {
		return
	}
	if s.havingAliasColumnIndex == nil {
		s.havingAliasColumnIndex = make(map[string]int)
	}
	s.havingAliasColumnIndex[alias] = idx
}

func (s *SelectPlan) getHavingAggregateColumnIndex(agg *ast.AggregateFuncExpr) (int, bool, error) {
	key, err := havingAggregateKey(agg)
	if err != nil {
		return 0, false, err
	}
	idx, ok := s.havingAggColumnIndex[key]
	return idx, ok, nil
}

func (s *SelectPlan) getHavingAliasColumnIndex(alias string) (int, bool) {
	alias = strings.ToLower(strings.TrimSpace(alias))
	if alias == "" || s.havingAliasColumnIndex == nil {
		return 0, false
	}
	idx, ok := s.havingAliasColumnIndex[alias]
	return idx, ok
}

// HasLimit if the select statement has limit clause, return true
func (s *SelectPlan) HasLimit() bool {
	return s.count != -1
}

// GetLimitValue get offset, count in limit clause
func (s *SelectPlan) GetLimitValue() (int64, int64) {
	return s.offset, s.count
}

// HasGroupBy if the select statement has group by clause, return true
func (s *SelectPlan) HasGroupBy() bool {
	return len(s.groupByColumn) != 0
}

// HasHaving if the select statement has having clause that needs post-filtering, return true
func (s *SelectPlan) HasHaving() bool {
	return s.havingExpr != nil
}

// GetHavingExpr get the having expression for post-merge filtering
func (s *SelectPlan) GetHavingExpr() ast.ExprNode {
	return s.havingExpr
}

// GetOriginColumnCount get origin column count in statement,
// since group by and order by may add extra columns to FieldList.
func (s *SelectPlan) GetOriginColumnCount() int {
	return s.originColumnCount
}

// GetColumnCount get column count with extra columns
func (s *SelectPlan) GetColumnCount() int {
	return s.columnCount
}

// GetGroupByColumnInfo get extra column offset and length for group by
func (s *SelectPlan) GetGroupByColumnInfo() []int {
	return s.groupByColumn
}

// HasOrderBy if select statement has order by clause, return true
func (s *SelectPlan) HasOrderBy() bool {
	return len(s.orderByDirections) != 0
}

// GetOrderByColumnInfo get extra column offset and length for order by
func (s *SelectPlan) GetOrderByColumnInfo() ([]int, []bool) {
	return s.orderByColumn, s.orderByDirections
}

// GetSQLs get generated SQLs
// the first key is slice, the second key is backend database name, the value is sql list.
func (s *SelectPlan) GetSQLs() map[string]map[string][]string {
	return s.sqls
}

// 执行计划是否仅仅涉及一个分片
func (s *SelectPlan) isExecOnSingleNode() bool {
	return len(s.result.indexes) == 1
}

func (s *SelectPlan) noAddColumns() bool {
	return s.originColumnCount == s.columnCount
}

// HandleSelectStmt build a SelectPlan
// 处理SelectStmt语法树, 改写其中一些节点, 并获取路由信息和结果聚合函数
func HandleSelectStmt(p *SelectPlan, stmt *ast.SelectStmt) error {
	p.stmt = stmt // hold the reference of stmt

	p.distinct = stmt.Distinct

	if err := handleTableRefs(p, stmt); err != nil {
		return fmt.Errorf("handle From error: %v", err)
	}

	// field list的处理必须在group by之前, 因为group by, order by会补列, 而这些补充的列是已经处理过的
	if stmt.Fields != nil {
		if err := handleFieldList(p, stmt); err != nil {
			return fmt.Errorf("handle Fields error: %v", err)
		}

		// 记录补列前的Fields长度
		p.originColumnCount = len(stmt.Fields.Fields)
	}

	if err := handleWhere(p, stmt); err != nil {
		return fmt.Errorf("handle Where error: %v", err)
	}

	if err := postHandleHintDatabaseFunction(p); err != nil {
		return fmt.Errorf("handle Hint error: %v", err)
	}

	//如果是多节点执行，特殊处理groupby orderby having limit
	if !p.isExecOnSingleNode() {
		// group by的处理必须在table处理之后
		if err := handleGroupBy(p, stmt); err != nil {
			return fmt.Errorf("handle GroupBy error: %v", err)
		}

		// order by的处理必须在table处理之后
		// 与group by补列的顺序没有要求, 只要保证处理返回结果去掉这些补充列时保持相反的顺序, 这里放在group by之后
		if err := handleOrderBy(p, stmt); err != nil {
			return fmt.Errorf("handle OrderBy error: %v", err)
		}

		handleExtraFieldList(p, stmt)

		// 记录补列后的Fields长度, 后面的handler不会补列了
		if stmt.Fields != nil {
			p.columnCount = len(stmt.Fields.Fields)
		}

		if err := handleHaving(p, stmt); err != nil {
			return fmt.Errorf("handle Having error: %v", err)
		}

		if err := handleLimit(p, stmt); err != nil {
			return fmt.Errorf("handle Limit error: %v", err)
		}
	} else {
		// 即使在单个分片执行，也会产生加列的情况
		if stmt.Fields != nil {
			p.columnCount = len(stmt.Fields.Fields)
		}
	}

	if err := postHandleGlobalTableRouteResultInQuery(p.StmtInfo); err != nil {
		return fmt.Errorf("post handle global table error: %v", err)
	}

	sqls, err := generateShardingSQLs(p.stmt, p.result, p.router)
	if err != nil {
		return fmt.Errorf("generate select SQL error: %v", err)
	}

	p.sqls = sqls

	return nil
}

// 处理GroupBy, 把GroupBy的列补到FieldList中, 然后把GroupBy去掉
func handleGroupBy(p *SelectPlan, stmt *ast.SelectStmt) error {
	if stmt.GroupBy == nil {
		return nil
	}

	groupByFields, err := createSelectFieldsFromByItems(p, stmt.GroupBy.Items)
	if err != nil {
		return fmt.Errorf("get group by fields error: %v", err)
	}

	for i := 0; i < len(groupByFields); i++ {
		p.groupByColumn = append(p.groupByColumn, i+len(stmt.Fields.Fields))
	}

	// append group by fields
	stmt.Fields.Fields = append(stmt.Fields.Fields, groupByFields...)

	return nil
}

func handleOrderBy(p *SelectPlan, stmt *ast.SelectStmt) error {
	if stmt.OrderBy == nil {
		return nil
	}

	orderByFields, err := createSelectFieldsFromByItems(p, stmt.OrderBy.Items)
	if err != nil {
		return fmt.Errorf("get order by fields error: %v", err)
	}

	for i := 0; i < len(orderByFields); i++ {
		p.orderByColumn = append(p.orderByColumn, i+len(stmt.Fields.Fields))
	}

	for _, f := range stmt.OrderBy.Items {
		p.orderByDirections = append(p.orderByDirections, f.Desc)
	}

	stmt.Fields.Fields = append(stmt.Fields.Fields, orderByFields...)
	return nil
}

func handleExtraFieldList(p *SelectPlan, stmt *ast.SelectStmt) {
	selectFields := make(map[string]int)
	for i := 0; i < p.originColumnCount; i++ {
		field := stmt.Fields.Fields[i]
		if field.AsName.L != "" {
			selectFields[field.AsName.L] = i
		}
		if field, isColumnExpr := stmt.Fields.Fields[i].Expr.(*ast.ColumnNameExpr); isColumnExpr {
			selectFields[field.Name.Name.L] = i
		}
	}

	deleteNum := 0
	for i := 0; i < len(p.groupByColumn); i++ {
		p.groupByColumn[i] -= deleteNum
		currColumnIndex := p.originColumnCount + i - deleteNum
		field, isColumnExpr := stmt.Fields.Fields[currColumnIndex].Expr.(*ast.ColumnNameExpr)
		if !isColumnExpr {
			continue
		}
		if index, ok := selectFields[field.Name.Name.L]; !ok {
			continue
		} else {
			stmt.Fields.Fields = append(stmt.Fields.Fields[:currColumnIndex], stmt.Fields.Fields[currColumnIndex+1:]...)
			p.groupByColumn[i] = index
			deleteNum++
		}
	}

	for i := 0; i < len(p.orderByColumn); i++ {
		p.orderByColumn[i] -= deleteNum
		currColumnIndex := p.originColumnCount + len(p.groupByColumn) + i - deleteNum
		field, isColumnExpr := stmt.Fields.Fields[currColumnIndex].Expr.(*ast.ColumnNameExpr)
		if !isColumnExpr {
			continue
		}
		if index, ok := selectFields[field.Name.Name.L]; !ok {
			continue
		} else {
			stmt.Fields.Fields = append(stmt.Fields.Fields[:currColumnIndex], stmt.Fields.Fields[currColumnIndex+1:]...)
			p.orderByColumn[i] = index
			deleteNum++
		}
	}
}

func createSelectFieldsFromByItems(p *SelectPlan, items []*ast.ByItem) ([]*ast.SelectField, error) {
	var ret []*ast.SelectField
	for _, item := range items {
		selectField, err := createSelectFieldFromByItem(p, item)
		if err != nil {
			return nil, err
		}
		ret = append(ret, selectField)
	}
	return ret, nil
}

func createSelectFieldFromByItem(p *SelectPlan, item *ast.ByItem) (*ast.SelectField, error) {
	var columnExpr *ast.ColumnNameExpr

	// support order by null/order by 1/order by max/order by min
	switch item.Expr.(type) {
	case *ast.FuncCallExpr:
		// 特殊处理DATABASE()这种情况
		funcExpr := item.Expr.(*ast.FuncCallExpr)
		if funcExpr.FnName.L == "database" {
			ret := &ast.SelectField{
				Expr: item.Expr,
			}
			return ret, nil
		}
		return nil, fmt.Errorf("ByItem.Expr is a FuncCallExpr but not DATABASE()")
	case *driver.ValueExpr:
		return &ast.SelectField{Expr: item.Expr}, nil
	case *ast.PositionExpr:
		return &ast.SelectField{Expr: item.Expr}, nil
	case *ast.AggregateFuncExpr:
		// 处理 order by max()/min() 等情况
		return &ast.SelectField{Expr: item.Expr}, nil
	case *ast.ColumnNameExpr:
		columnExpr = item.Expr.(*ast.ColumnNameExpr)
	default:
		return nil, fmt.Errorf("ByItem.Expr is not a ColumnNameExpr")
	}

	rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInField(p.TableAliasStmtInfo, columnExpr)
	if err != nil {
		return nil, err
	}

	if need {
		decorator := CreateColumnNameExprDecorator(columnExpr, rule, isAlias, p.GetRouteResult())
		item.Expr = decorator
	}

	ret := &ast.SelectField{
		Expr: columnExpr,
	}
	return ret, nil
}

// 处理from table和join on部分
// 主要是改写table ExprNode, 并找到路由条件
func handleTableRefs(p *SelectPlan, stmt *ast.SelectStmt) error {
	tableRefs := stmt.From
	if tableRefs == nil {
		return nil
	}

	join := tableRefs.TableRefs
	if join == nil {
		return nil
	}

	return handleJoin(p.TableAliasStmtInfo, join)
}

func handleJoin(p *TableAliasStmtInfo, join *ast.Join) error {
	if err := precheckJoinClause(join); err != nil {
		return fmt.Errorf("precheck Join error: %v", err)
	}

	// 只允许最多两个表的JOIN
	if join.Left != nil {
		switch left := join.Left.(type) {
		case *ast.TableSource:
			// 改写两个表的node
			if err := rewriteTableSource(p, left); err != nil {
				return fmt.Errorf("rewrite left TableSource error: %v", err)
			}
		case *ast.Join:
			if err := handleJoin(p, left); err != nil {
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

		if err := rewriteTableSource(p, right); err != nil {
			return fmt.Errorf("rewrite right TableSource error: %v", err)
		}
	}

	// 改写ON条件
	if join.On != nil {
		err := rewriteOnCondition(p, join.On)
		if err != nil {
			return fmt.Errorf("rewrite on condition error: %v", err)
		}
	}

	return nil
}

func handleWhere(p *SelectPlan, stmt *ast.SelectStmt) (err error) {
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

// 检查TableRefs中存在的不允许在分表中执行的语法
func precheckJoinClause(join *ast.Join) error {
	// 不允许USING的列名中出现DB名和表名, 因为目前Join子句的TableName不方便加装饰器
	for _, c := range join.Using {
		if c.Schema.String() != "" {
			return fmt.Errorf("JOIN does not support USING column with schema")
		}
		if c.Table.String() != "" {
			return fmt.Errorf("JOIN does not support USING column with table")
		}
	}
	return nil
}

// 改写TableSource节点, 得到一个装饰器
// Source必须为TableName节点或子查询
func rewriteTableSource(p *TableAliasStmtInfo, tableSource *ast.TableSource) error {
	switch ss := tableSource.Source.(type) {
	case *ast.TableName:
		return rewriteTableNameInTableSource(p, tableSource)
	case *ast.SelectStmt:
		// 递归处理子查询
		if err := handleSubquerySelectStmt(p, ss); err != nil {
			return fmt.Errorf("handleSubquerySelectStmt error: %v", err)
		}
		alias := tableSource.AsName.L
		if alias != "" {
			// 记录子查询的别名
			if _, err := p.RecordSubqueryTableAlias(alias); err != nil {
				return fmt.Errorf("record subquery alias error: %v", err)
			}
		}
		return nil
	default:
		return fmt.Errorf("field Source cannot handle, type: %T", tableSource.Source)
	}
}

func rewriteTableNameInTableSource(p *TableAliasStmtInfo, tableSource *ast.TableSource) error {
	tableName, ok := tableSource.Source.(*ast.TableName)
	if !ok {
		return fmt.Errorf("field Source is not type of TableName, type: %T", tableSource.Source)
	}
	alias := tableSource.AsName.L

	rule, need, err := NeedCreateTableNameDecorator(p, tableName, alias)
	if err != nil {
		return fmt.Errorf("check NeedCreateTableNameDecorator error: %v", err)
	}

	if !need {
		return nil
	}

	// 这是一个分片表或关联表, 创建一个TableName的装饰器, 并替换原有节点
	d, err := CreateTableNameDecorator(tableName, rule, p.GetRouteResult())
	d.Alias = tableSource.AsName.String()
	if err != nil {
		return fmt.Errorf("create TableNameDecorator error: %v", err)
	}
	tableSource.Source = d
	tableSource.AsName.L = ""
	tableSource.AsName.O = ""
	return nil
}

func rewriteOnCondition(p *TableAliasStmtInfo, on *ast.OnCondition) error {
	has, result, decorator, err := handleComparisonExpr(p, on.Expr)
	if err != nil {
		return fmt.Errorf("rewrite Expr in OnCondition error: %v", err)
	}
	if has {
		p.GetRouteResult().Inter(result)
	}
	on.Expr = decorator
	return nil
}

// 处理info中的hint
// 目前只有mycat路由方式支持
// hint路由会覆盖遍历语法树时计算出的路由
func postHandleHintDatabaseFunction(p *SelectPlan) error {
	if p.hintPhyDB == "" {
		return nil
	}

	rule, ok := p.router.GetShardRule(p.result.db, p.result.table)
	if !ok {
		return fmt.Errorf("sharding rule of route result not found, result: %v", p.result)
	}
	mr, ok := rule.(router.MycatRule)
	if !ok {
		return fmt.Errorf("sharding rule is not mycat mode, result: %v", p.result)
	}

	if !router.IsMycatShardingRule(mr.GetType()) { // TODO: need refactor, why is MycatRule's type not mycat rule?
		return fmt.Errorf("only mycat rule supports database function hint")
	}

	idx, ok := mr.GetTableIndexByDatabaseName(p.hintPhyDB)
	if !ok {
		return fmt.Errorf("hint db not found: %s", p.hintPhyDB)
	}

	p.result.indexes = []int{idx}
	return nil
}

// ColumnNameRewriteVisitor visit ColumnNameExpr, check if need decorate, and then decorate it.
type ColumnNameRewriteVisitor struct {
	info *TableAliasStmtInfo
}

// NewColumnNameRewriteVisitor constructor of ColumnNameRewriteVisitor
func NewColumnNameRewriteVisitor(p *TableAliasStmtInfo) *ColumnNameRewriteVisitor {
	return &ColumnNameRewriteVisitor{
		info: p,
	}
}

// Enter implement ast.Visitor
func (s *ColumnNameRewriteVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	return n, false
}

// Leave implement ast.Visitor
func (s *ColumnNameRewriteVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	field, ok := n.(*ast.ColumnNameExpr)
	if !ok {
		return n, true
	}
	rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInField(s.info, field)
	if err != nil {
		panic(fmt.Errorf("check NeedCreateColumnNameExprDecoratorInField in ColumnNameExpr error: %v", err))
	}
	if need {
		decorator := CreateColumnNameExprDecorator(field, rule, isAlias, s.info.GetRouteResult())
		return decorator, true
	}

	return n, true
}

func handleFieldList(p *SelectPlan, stmt *ast.SelectStmt) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("handleFieldList panic: %v", e)
		}
	}()

	fields := stmt.Fields
	if fields == nil {
		return nil
	}

	// 先用一个Visitor生成一个替换表名的装饰器
	// 这里如果出错, 只能通过panic返回err
	columnNameRewriter := NewColumnNameRewriteVisitor(p.TableAliasStmtInfo)
	fields.Accept(columnNameRewriter)

	// 如果最外层是聚合函数, 则生成一个聚合函数装饰器, 并记录对应的列位置
	// 只处理最外层的聚合函数.
	for i, f := range fields.Fields {
		switch field := f.Expr.(type) {
		case *ast.AggregateFuncExpr:
			merger, err := CreateAggregateFunctionMerger(field, i)
			if err != nil {
				return fmt.Errorf("create aggregate function merger error, column index: %d, err: %v", i, err)
			}
			if err := p.setAggregateFuncMerger(i, merger); err != nil {
				return fmt.Errorf("set aggregate function merger error, column index: %d, err: %v", i, err)
			}
			if err := p.setHavingAggregateColumnIndex(field, i); err != nil {
				return fmt.Errorf("set having aggregate column index error, column index: %d, err: %v", i, err)
			}
			p.setHavingAliasColumnIndex(f.AsName.L, i)
		default:
			// do nothing
		}
	}
	return nil
}

func havingAggregateKey(agg *ast.AggregateFuncExpr) (string, error) {
	if agg == nil {
		return "", fmt.Errorf("nil aggregate function")
	}

	args := make([]string, 0, len(agg.Args))
	for _, arg := range agg.Args {
		key, err := havingAggregateArgKey(arg)
		if err != nil {
			return "", err
		}
		args = append(args, key)
	}

	distinct := ""
	if agg.Distinct {
		distinct = "distinct "
	}
	funcName := strings.ToLower(strings.TrimSpace(agg.F))
	return fmt.Sprintf("%s(%s%s)", funcName, distinct, strings.Join(args, ",")), nil
}

func havingAggregateArgKey(expr ast.ExprNode) (string, error) {
	switch e := expr.(type) {
	case *ast.ColumnNameExpr:
		return havingColumnNameKey(e.Name), nil
	case *ColumnNameExprDecorator:
		return havingColumnNameKey(e.ColumnNameExpr.Name), nil
	default:
		key, err := parser.NodeToStringWithoutQuote(expr)
		if err != nil {
			return "", fmt.Errorf("restore aggregate argument error: %v", err)
		}
		return strings.ToLower(strings.TrimSpace(key)), nil
	}
}

func havingColumnNameKey(name *ast.ColumnName) string {
	if name == nil {
		return ""
	}

	parts := make([]string, 0, 3)
	if schema := havingCIStrKey(name.Schema.O, name.Schema.L); schema != "" {
		parts = append(parts, schema)
	}
	if table := havingCIStrKey(name.Table.O, name.Table.L); table != "" {
		parts = append(parts, table)
	}
	if column := havingCIStrKey(name.Name.O, name.Name.L); column != "" {
		parts = append(parts, column)
	}
	return strings.Join(parts, ".")
}

func havingCIStrKey(original, lower string) string {
	if original != "" {
		return strings.ToLower(strings.TrimSpace(original))
	}
	return strings.ToLower(strings.TrimSpace(lower))
}

func buildAggregateAliasMap(fields []*ast.SelectField) map[string]*ast.AggregateFuncExpr {
	ret := make(map[string]*ast.AggregateFuncExpr)
	for _, f := range fields {
		if f.AsName.L == "" {
			continue
		}
		if agg, ok := f.Expr.(*ast.AggregateFuncExpr); ok {
			ret[strings.ToLower(f.AsName.L)] = agg
		}
	}
	return ret
}

type havingAggregateCollector struct {
	aliasAgg map[string]*ast.AggregateFuncExpr
	found    bool
	aggs     []*ast.AggregateFuncExpr
}

// Enter implement ast.Visitor
func (v *havingAggregateCollector) Enter(n ast.Node) (ast.Node, bool) {
	switch expr := n.(type) {
	case *ast.AggregateFuncExpr:
		v.found = true
		v.aggs = append(v.aggs, expr)
		return n, true
	case *ast.ColumnNameExpr:
		name := strings.ToLower(expr.Name.Name.L)
		if _, ok := v.aliasAgg[name]; ok {
			v.found = true
		}
		return n, true
	case *ColumnNameExprDecorator:
		name := strings.ToLower(expr.ColumnNameExpr.Name.Name.L)
		if _, ok := v.aliasAgg[name]; ok {
			v.found = true
		}
		return n, true
	default:
		return n, false
	}
}

// Leave implement ast.Visitor
func (v *havingAggregateCollector) Leave(n ast.Node) (ast.Node, bool) {
	return n, true
}

func handleHaving(p *SelectPlan, stmt *ast.SelectStmt) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("handleHaving panic: %v", e)
		}
	}()

	having := stmt.Having
	if having == nil {
		return nil
	}

	// 先用一个Visitor生成一个替换表名的装饰器
	// 这里如果出错, 只能通过panic返回err
	columnNameRewriter := NewColumnNameRewriteVisitor(p.TableAliasStmtInfo)
	having.Accept(columnNameRewriter)

	aliasAgg := buildAggregateAliasMap(stmt.Fields.Fields)
	collector := &havingAggregateCollector{aliasAgg: aliasAgg}
	having.Expr.Accept(collector)

	// 检查HAVING表达式是否包含聚合函数或引用了聚合别名
	// 只有包含聚合函数或聚合别名时才需要在全局层重新应用过滤
	// 简单的列比较条件（如 HAVING id=1）可以安全地下推到各分片
	if collector.found {
		// 保存HAVING表达式用于在合并结果后重新应用过滤
		// 同时从stmt中移除HAVING, 使下推到各分片的SQL不包含HAVING子句
		// 这是因为聚合函数的结果需要跨分片汇总后才能正确评估HAVING条件
		p.havingExpr = having.Expr
		stmt.Having = nil

		// Build a set of existing aggregate expressions in SELECT fields
		existingAgg := make(map[string]struct{})
		for i, f := range stmt.Fields.Fields {
			if agg, ok := f.Expr.(*ast.AggregateFuncExpr); ok {
				key, err := havingAggregateKey(agg)
				if err != nil {
					return fmt.Errorf("build having aggregate key error: %v", err)
				}
				if err := p.setHavingAggregateColumnIndex(agg, i); err != nil {
					return fmt.Errorf("set having aggregate column index error, column index: %d, err: %v", i, err)
				}
				p.setHavingAliasColumnIndex(f.AsName.L, i)
				existingAgg[key] = struct{}{}
			}
		}

		// Ensure HAVING aggregate expressions are present in SELECT fields
		for _, agg := range collector.aggs {
			key, err := havingAggregateKey(agg)
			if err != nil {
				return fmt.Errorf("build having aggregate key error: %v", err)
			}
			if _, ok := existingAgg[key]; ok {
				continue
			}

			field := &ast.SelectField{Expr: agg, Auxiliary: true}
			stmt.Fields.Fields = append(stmt.Fields.Fields, field)
			idx := len(stmt.Fields.Fields) - 1
			merger, err := CreateAggregateFunctionMerger(agg, idx)
			if err != nil {
				return fmt.Errorf("create aggregate function merger error, column index: %d, err: %v", idx, err)
			}
			if err := p.setAggregateFuncMerger(idx, merger); err != nil {
				return fmt.Errorf("set aggregate function merger error, column index: %d, err: %v", idx, err)
			}
			if err := p.setHavingAggregateColumnIndex(agg, idx); err != nil {
				return fmt.Errorf("set having aggregate column index error, column index: %d, err: %v", idx, err)
			}
			existingAgg[key] = struct{}{}
		}

		p.columnCount = len(stmt.Fields.Fields)
	}

	return nil
}

// containsAggregateFunc checks if an expression contains aggregate functions
func containsAggregateFunc(expr ast.ExprNode) bool {
	if expr == nil {
		return false
	}

	switch e := expr.(type) {
	case *ast.AggregateFuncExpr:
		return true
	case *ast.BinaryOperationExpr:
		return containsAggregateFunc(e.L) || containsAggregateFunc(e.R)
	case *ast.ParenthesesExpr:
		return containsAggregateFunc(e.Expr)
	case *ast.UnaryOperationExpr:
		return containsAggregateFunc(e.V)
	case *ast.FuncCallExpr:
		for _, arg := range e.Args {
			if containsAggregateFunc(arg) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func handleComparisonExpr(p *TableAliasStmtInfo, comp ast.ExprNode) (bool, []int, ast.ExprNode, error) {
	switch expr := comp.(type) {
	case *ast.BinaryOperationExpr:
		return handleBinaryOperationExpr(p, expr)
	case *ast.PatternInExpr:
		return handlePatternInExpr(p, expr)
	case *ast.BetweenExpr:
		return handleBetweenExpr(p, expr)
	case *ast.ParenthesesExpr:
		has, routeResult, newExpr, err := handleComparisonExpr(p, expr.Expr)
		expr.Expr = newExpr
		return has, routeResult, expr, err
	default:
		// 其他情况只替换表名 (但是不处理根节点是ColumnNameExpr的情况, 理论上也不会出现这种情况)
		columnNameRewriter := NewColumnNameRewriteVisitor(p)
		expr.Accept(columnNameRewriter)
		return false, p.GetRouteResult().GetShardIndexes(), comp, nil
	}
}

func handlePatternInExpr(p *TableAliasStmtInfo, expr *ast.PatternInExpr) (bool, []int, ast.ExprNode, error) {
	rule, need, isAlias, err := NeedCreatePatternInExprDecorator(p, expr)
	if err != nil {
		return false, nil, nil, fmt.Errorf("check PatternInExpr error: %v", err)
	}
	if !need {
		return false, nil, expr, nil
	}
	decorator, err := CreatePatternInExprDecorator(expr, rule, isAlias, p.GetRouteResult())
	if err != nil {
		return false, nil, nil, fmt.Errorf("create PatternInExprDecorator error: %v", err)
	}
	return true, decorator.GetCurrentRouteResult(), decorator, nil
}

func handleBetweenExpr(p *TableAliasStmtInfo, expr *ast.BetweenExpr) (bool, []int, ast.ExprNode, error) {
	rule, need, isAlias, err := NeedCreateBetweenExprDecorator(p, expr)
	if err != nil {
		return false, nil, nil, fmt.Errorf("check BetweenExpr error: %v", err)
	}
	if !need {
		return false, nil, expr, nil
	}

	decorator, err := CreateBetweenExprDecorator(expr, rule, isAlias, p.GetRouteResult())
	if err != nil {
		return false, nil, nil, fmt.Errorf("create CreateBetweenExprDecorator error: %v", err)
	}

	return true, decorator.GetCurrentRouteResult(), decorator, nil
}

// return value: hasRoutingResult, RouteResult, Decorator, error
// the Decorator must not be nil. If no modification to the input expr, just return it.
func handleBinaryOperationExpr(p *TableAliasStmtInfo, expr *ast.BinaryOperationExpr) (bool, []int, ast.ExprNode, error) {
	_, ok := opcode.Ops[expr.Op]
	if !ok {
		return false, nil, nil, fmt.Errorf("unknown BinaryOperationExpr.Op: %v", expr.Op)
	}

	switch expr.Op {
	case opcode.LogicAnd, opcode.LogicOr:
		return handleBinaryOperationExprLogic(p, expr)
	case opcode.EQ, opcode.NE, opcode.GT, opcode.GE, opcode.LT, opcode.LE:
		return handleBinaryOperationExprMathCompare(p, expr)
	default:
		return handleBinaryOperationExprOther(p, expr)
	}
}

// 处理逻辑比较运算
func handleBinaryOperationExprLogic(p *TableAliasStmtInfo, expr *ast.BinaryOperationExpr) (bool, []int, ast.ExprNode, error) {
	lHas, lResult, lDecorator, lErr := handleComparisonExpr(p, expr.L)
	if lErr != nil {
		return false, nil, nil, fmt.Errorf("handle BinaryOperationExpr.L error: %v", lErr)
	}
	rHas, rResult, rDecorator, rErr := handleComparisonExpr(p, expr.R)
	if rErr != nil {
		return false, nil, nil, fmt.Errorf("handle BinaryOperationExpr.R error: %v", rErr)
	}

	if lDecorator != nil {
		expr.L = lDecorator
	}
	if rDecorator != nil {
		expr.R = rDecorator
	}

	has, result := mergeBinaryOperationRouteResult(expr.Op, lHas, lResult, rHas, rResult)
	return has, result, expr, nil
}

// 处理算术比较运算
// 如果出现列名, 则必须为列名与列名比较, 列名与值比较, 否则会报错 (比如 id + 2 = 3 就会报错, 因为 id + 2 处理不了)
// 如果是其他情况, 则直接返回 (如 1 = 1 这种)
func handleBinaryOperationExprMathCompare(p *TableAliasStmtInfo, expr *ast.BinaryOperationExpr) (bool, []int, ast.ExprNode, error) {
	lType := getExprNodeTypeInBinaryOperation(expr.L)
	rType := getExprNodeTypeInBinaryOperation(expr.R)

	// handle hint database function: SELECT * from tbl where DATABASE() = db_0 / 'db_0' / `db_0`
	if expr.Op == opcode.EQ {
		if lType == FuncCallExpr {
			hintDB, err := getDatabaseFuncHint(expr.L.(*ast.FuncCallExpr), expr.R)
			if err != nil {
				return false, nil, nil, fmt.Errorf("get database function hint error: %v", err)
			}
			if hintDB != "" {
				p.hintPhyDB = hintDB
				return false, nil, expr, nil
			}
		} else if rType == FuncCallExpr {
			hintDB, err := getDatabaseFuncHint(expr.R.(*ast.FuncCallExpr), expr.L)
			if err != nil {
				return false, nil, nil, fmt.Errorf("get database function hint error: %v", err)
			}
			if hintDB != "" {
				p.hintPhyDB = hintDB
				return false, nil, expr, nil
			}
		}
	}

	if lType == ColumnNameExpr && rType == ColumnNameExpr {
		return handleBinaryOperationExprCompareLeftColumnRightColumn(p, expr)
	}

	if lType == ColumnNameExpr {
		if rType == ValueExpr {
			return handleBinaryOperationExprCompareLeftColumnRightValue(p, expr, getFindTableIndexesFunc(expr.Op))
		}
		column := expr.L.(*ast.ColumnNameExpr)
		rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInCondition(p, column)
		if err != nil {
			return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.L: %v", err)
		}
		if !need {
			return false, nil, expr, nil
		}

		decorator := CreateColumnNameExprDecorator(column, rule, isAlias, p.GetRouteResult())
		expr.L = decorator
		return false, nil, expr, nil
	}

	if rType == ColumnNameExpr {
		if lType == ValueExpr {
			return handleBinaryOperationExprCompareLeftValueRightColumn(p, expr, getFindTableIndexesFunc(inverseOperator(expr.Op)))
		}
		column := expr.R.(*ast.ColumnNameExpr)
		rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInCondition(p, column)
		if err != nil {
			return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.R: %v", err)
		}
		if !need {
			return false, nil, expr, nil
		}

		decorator := CreateColumnNameExprDecorator(column, rule, isAlias, p.GetRouteResult())
		expr.R = decorator
		return false, nil, expr, nil
	}

	return false, nil, expr, nil
}

// 处理其他情况的运算
// 如果出现分表列, 只创建一个替换表名的装饰器, 不计算路由. 因此返回结果前两个一定是false, nil
func handleBinaryOperationExprOther(p *TableAliasStmtInfo, expr *ast.BinaryOperationExpr) (bool, []int, ast.ExprNode, error) {
	if lColumn, ok := expr.L.(*ast.ColumnNameExpr); ok {
		lRule, lNeed, lIsAlias, lErr := NeedCreateColumnNameExprDecoratorInCondition(p, lColumn)
		if lErr != nil {
			return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.L: %v", lErr)
		}

		if lNeed {
			lDecorator := CreateColumnNameExprDecorator(lColumn, lRule, lIsAlias, p.GetRouteResult())
			expr.L = lDecorator
		}
	}
	if rColumn, ok := expr.R.(*ast.ColumnNameExpr); ok {
		rRule, rNeed, rIsAlias, rErr := NeedCreateColumnNameExprDecoratorInCondition(p, rColumn)
		if rErr != nil {
			return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.R: %v", rErr)
		}
		if rNeed {
			rDecorator := CreateColumnNameExprDecorator(rColumn, rRule, rIsAlias, p.GetRouteResult())
			expr.R = rDecorator
		}
	}
	return false, nil, expr, nil
}

// 获取mycat路由模式下的hint物理DB名
func getDatabaseFuncHint(f *ast.FuncCallExpr, v ast.ExprNode) (string, error) {
	if f.FnName.L != "database" {
		return "", nil
	}
	switch vv := v.(type) {
	case *ast.ColumnNameExpr:
		return vv.Name.Name.String(), nil
	case *driver.ValueExpr:
		return vv.GetString(), nil
	default:
		return "", fmt.Errorf("invalid value type of database function hint: %T", v)
	}
}

// 返回一个根据路由信息和路由值获取路由结果的函数
// 左边为列名, 右边为参数
func getFindTableIndexesFunc(op opcode.Op) func(rule router.Rule, columnName string, v interface{}) ([]int, error) {
	findTableIndexesFunc := func(rule router.Rule, columnName string, v interface{}) ([]int, error) {
		// 如果不是分表列, 则需要返回所有分片
		if rule.GetShardingColumn() != columnName {
			return rule.GetSubTableIndexes(), nil
		}

		// 如果是分表列, 还需要根据运算符判断
		switch op {
		case opcode.EQ:
			index, err := rule.FindTableIndex(v)
			if err != nil {
				return nil, err
			}
			return []int{index}, nil
		case opcode.NE:
			return rule.GetSubTableIndexes(), nil
		case opcode.GT, opcode.GE, opcode.LT, opcode.LE:
			// 如果是range路由, 需要做一些特殊处理
			if rangeShard, ok := rule.GetShard().(router.RangeShard); ok {
				index, err := rule.FindTableIndex(v)
				if err != nil {
					return nil, err
				}
				if op == opcode.LT || op == opcode.LE {
					if op == opcode.LT {
						index = adjustShardIndex(rangeShard, v, index)
					}
					return makeList(rule.GetFirstTableIndex(), index+1), nil
				} else {
					return makeList(index, rule.GetLastTableIndex()+1), nil
				}
			}

			// 如果不是 (即hash路由), 则返回所有分片
			return rule.GetSubTableIndexes(), nil
		default: // should not going here
			return rule.GetSubTableIndexes(), nil
		}
	}

	return findTableIndexesFunc
}

// copy from PlanBuilder.adjustShardIndex()
func adjustShardIndex(s router.RangeShard, value interface{}, index int) int {
	if s.EqualStart(value, index) {
		return index - 1
	}
	return index
}

func inverseOperator(op opcode.Op) opcode.Op {
	switch op {
	case opcode.GT:
		return opcode.LT
	case opcode.GE:
		return opcode.LE
	case opcode.LT:
		return opcode.GT
	case opcode.LE:
		return opcode.GE
	default:
		return op
	}
}

func handleBinaryOperationExprCompareLeftColumnRightColumn(p *TableAliasStmtInfo, expr *ast.BinaryOperationExpr) (bool, []int, ast.ExprNode, error) {
	lColumn := expr.L.(*ast.ColumnNameExpr)
	lRule, lNeed, lIsAlias, lErr := NeedCreateColumnNameExprDecoratorInCondition(p, lColumn)
	if lErr != nil {
		return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.L: %v", lErr)
	}
	rColumn := expr.R.(*ast.ColumnNameExpr)
	rRule, rNeed, rIsAlias, rErr := NeedCreateColumnNameExprDecoratorInCondition(p, rColumn)
	if rErr != nil {
		return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.R: %v", rErr)
	}

	if lNeed {
		lDecorator := CreateColumnNameExprDecorator(lColumn, lRule, lIsAlias, p.GetRouteResult())
		expr.L = lDecorator
	}
	if rNeed {
		rDecorator := CreateColumnNameExprDecorator(rColumn, rRule, rIsAlias, p.GetRouteResult())
		expr.R = rDecorator
	}
	return false, nil, expr, nil
}

func handleBinaryOperationExprCompareLeftColumnRightValue(p *TableAliasStmtInfo, expr *ast.BinaryOperationExpr, findTableIndexes func(router.Rule, string, interface{}) ([]int, error)) (bool, []int, ast.ExprNode, error) {
	column := expr.L.(*ast.ColumnNameExpr)
	rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInCondition(p, column)
	if err != nil {
		return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.L: %v", err)
	}
	if !need {
		return false, nil, expr, nil
	}

	decorator := CreateColumnNameExprDecorator(column, rule, isAlias, p.GetRouteResult())
	expr.L = decorator

	if rule.GetType() == router.GlobalTableRuleType {
		return false, nil, expr, nil
	}

	valueExpr := expr.R.(*driver.ValueExpr)
	v, err := util.GetValueExprResult(valueExpr)
	if err != nil {
		return false, nil, nil, fmt.Errorf("get ValueExpr value error: %v", err)
	}

	tableIndexes, err := findTableIndexes(rule, column.Name.Name.L, v)
	if err != nil {
		return false, nil, nil, fmt.Errorf("find table index error: %v", err)
	}

	return true, tableIndexes, expr, nil
}

func handleBinaryOperationExprCompareLeftValueRightColumn(p *TableAliasStmtInfo, expr *ast.BinaryOperationExpr, findTableIndexes func(router.Rule, string, interface{}) ([]int, error)) (bool, []int, ast.ExprNode, error) {
	column := expr.R.(*ast.ColumnNameExpr)
	rule, need, isAlias, err := NeedCreateColumnNameExprDecoratorInCondition(p, column)
	if err != nil {
		return false, nil, nil, fmt.Errorf("check ColumnNameExpr error in BinaryOperationExpr.R: %v", err)
	}
	if !need {
		return false, nil, expr, nil
	}

	decorator := CreateColumnNameExprDecorator(column, rule, isAlias, p.GetRouteResult())
	expr.R = decorator

	if rule.GetType() == router.GlobalTableRuleType {
		return false, nil, expr, nil
	}

	valueExpr := expr.L.(*driver.ValueExpr)
	v, err := util.GetValueExprResult(valueExpr)
	if err != nil {
		return false, nil, nil, fmt.Errorf("get ValueExpr value error: %v", err)
	}

	tableIndexes, err := findTableIndexes(rule, column.Name.Name.L, v)
	if err != nil {
		return false, nil, nil, fmt.Errorf("find table index error: %v", err)
	}

	return true, tableIndexes, expr, nil
}

func mergeBinaryOperationRouteResult(op opcode.Op, lHas bool, lResult []int, rHas bool, rResult []int) (bool, []int) {
	switch op {
	case opcode.LogicAnd:
		if lHas == false && rHas == false {
			return false, nil
		}
		if lHas && rHas {
			return true, interList(lResult, rResult)
		}
		if lHas {
			return true, lResult
		}
		if rHas {
			return true, rResult
		}
	case opcode.LogicOr:
		if lHas && rHas {
			return true, unionList(lResult, rResult)
		}
		return false, nil
	}
	return false, nil
}

func handleLimit(p *SelectPlan, stmt *ast.SelectStmt) error {
	need, originOffset, originCount, newLimit := NeedRewriteLimitOrCreateRewrite(stmt)
	p.offset = originOffset
	p.count = originCount
	if need {
		stmt.Limit = newLimit
	}
	return nil
}

func getTableInfoFromTableName(t *ast.TableName) (string, string) {
	return t.Schema.O, t.Name.L
}

func getColumnInfoFromColumnName(t *ast.ColumnName) (string, string, string) {
	return t.Schema.O, t.Table.L, t.Name.L
}
