package plan

import (
	"fmt"
	"strings"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
	driver "github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/proxy/sequence"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/math"
)

// UnionPlan is the plan for union statement
type UnionPlan struct {
	basePlan

	db                          string
	router                      *router.Router
	seq                         *sequence.SequenceManager
	hintPlan                    Plan
	stmt                        *ast.UnionStmt
	subPlans                    []Plan
	subPlanIsAfterUnionDistinct []bool

	unionOrderBy *ast.OrderByClause
	unionLimit   *ast.Limit
}

// NewUnionPlan constructor of UnionPlan
func NewUnionPlan(db string, sql string, r *router.Router, seq *sequence.SequenceManager, hintPlan Plan) *UnionPlan {
	return &UnionPlan{
		db:                          db,
		router:                      r,
		seq:                         seq,
		hintPlan:                    hintPlan,
		subPlans:                    []Plan{},
		subPlanIsAfterUnionDistinct: []bool{},
	}
}

// HandleUnionStmt processes the union statement
func (p *UnionPlan) HandleUnionStmt(stmt *ast.UnionStmt) error {
	if stmt == nil {
		return fmt.Errorf("handle union stmt error: ast.UnionStmt is nil")
	}

	p.stmt = stmt
	for _, selectStmt := range stmt.SelectList.Selects {
		selectSQL, err := restoreSelectStmt(selectStmt)
		if err != nil {
			return fmt.Errorf("failed to restore select statement: %v", err)
		}
		plan, err := buildShardPlan(selectStmt, p.db, selectSQL, p.router, p.seq, p.hintPlan)
		if err != nil {
			return fmt.Errorf("failed to build shard plan: %v", err)
		}
		p.subPlans = append(p.subPlans, plan)
		p.subPlanIsAfterUnionDistinct = append(p.subPlanIsAfterUnionDistinct, selectStmt.IsAfterUnionDistinct)
	}
	// Handle ORDER BY and LIMIT at the UNION level
	if stmt.OrderBy != nil {
		p.unionOrderBy = stmt.OrderBy
	}
	if stmt.Limit != nil {
		p.unionLimit = stmt.Limit
	}
	return nil
}

// ExecuteIn implements the Plan interface
func (p *UnionPlan) ExecuteIn(reqCtx *util.RequestContext, sess Executor) (*mysql.Result, error) {
	var results []*mysql.Result
	for _, plan := range p.subPlans {
		result, err := plan.ExecuteIn(reqCtx, sess)
		if err != nil {
			return nil, fmt.Errorf("failed to execute sub-plan: %v", err)
		}
		results = append(results, result)
	}

	// Merge results according to UNION semantics
	return p.MergeUnionResult(results)
}

// MergeUnionResult merges all sub-plan results into a single result according to UNION semantics
func (p *UnionPlan) MergeUnionResult(rs []*mysql.Result) (mergedResult *mysql.Result, err error) {
	mergedResult, err = p.mergeMultiResultSet(rs)
	if err != nil {
		return mergedResult, err
	}

	// Apply ORDER BY if required
	if p.unionOrderBy != nil {
		if err := p.handleUnionOrderBy(mergedResult); err != nil {
			return nil, fmt.Errorf("failed to sort result: %v", err)
		}
	}

	// Apply LIMIT if required
	if p.unionLimit != nil {
		if err := p.handleUnionLimit(mergedResult); err != nil {
			return nil, fmt.Errorf("failed to apply limit/offset: %v", err)
		}
	}
	// Construct RowData to return
	if err := GenerateSelectResultRowData(mergedResult); err != nil {
		return nil, fmt.Errorf("failed to generate select result row data: %v", err)
	}

	return mergedResult, nil
}

// restoreSelectStmt is helper function to restore SelectStmt to SQL string
func restoreSelectStmt(stmt *ast.SelectStmt) (string, error) {
	var sb strings.Builder
	ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, &sb)
	if err := stmt.Restore(ctx); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (p *UnionPlan) mergeMultiResultSet(rs []*mysql.Result) (mergedResult *mysql.Result, err error) {
	// An empty result set (which has no rows but may have column definitions) is different from a nil result set (which represents no result set at all).
	if len(rs) == 0 {
		// Return an empty result set with empty Fields and Values
		mergedResult = &mysql.Result{
			Resultset: &mysql.Resultset{
				Fields:     []*mysql.Field{},  // Empty slice of fields
				FieldNames: map[string]int{},  // Empty map of field names
				Values:     [][]interface{}{}, // Empty slice of values
				RowDatas:   []mysql.RowData{}, // Empty slice of row data
			},
			Status:       0,
			InsertID:     0,
			AffectedRows: 0,
			Warnings:     0,
			Info:         "",
		}
		return mergedResult, nil
	}

	if len(rs) == 1 {
		return rs[0], nil
	}
	// Initialize mergedResult with the first result
	mergedResult = rs[0]
	mergedResult.Fields = rs[0].Fields

	for i := 1; i < len(rs); i++ {
		// Check if the number of columns matches
		if len(rs[i].Fields) != len(mergedResult.Fields) {
			return nil, fmt.Errorf("unmatched result set column count: %d not equal to: %d", len(rs[i].Fields), len(mergedResult.Fields))
		}
		// Ensure that the status flags of the result set can reflect the execution status of all sub-queries.
		mergedResult.Status |= rs[i].Status
		// Check whether the corresponding column types are consistent
		for j, field := range rs[i].Fields {
			mergedField := mergedResult.Fields[j]
			if mergedField.Type != field.Type {
				return nil, fmt.Errorf(
					"unmatched result set (column `%s`) type: %s not equal to (column `%s`) type: %s",
					string(field.Name), mysql.MysqlTypeName(field.Type),
					string(mergedField.Name), mysql.MysqlTypeName(mergedField.Type),
				)
			}
		}

		if p.subPlanIsAfterUnionDistinct[i] {
			// UNION DISTINCT: Merge and eliminate duplicates
			combinedValues := append(mergedResult.Values, rs[i].Values...)
			uniqueValues, err := p.removeDuplicateValues(combinedValues, len(mergedResult.Fields))
			if err != nil {
				return nil, fmt.Errorf("failed to remove duplicates: %v", err)
			}
			mergedResult.Values = uniqueValues

		} else {
			// UNION ALL: Append without removing duplicates
			mergedResult.Values = append(mergedResult.Values, rs[i].Values...)
		}
	}
	// reset RowDatas
	mergedResult.RowDatas = nil
	return mergedResult, nil
}

// handleUnionOrderBy handle order by in the union result
func (p *UnionPlan) handleUnionOrderBy(res *mysql.Result) error {
	if res == nil || res.Resultset == nil {
		return nil
	}

	if err := p.sortResult(res); err != nil {
		return fmt.Errorf("union plan failed to sort result: %v", err)
	}
	return nil
}

// handleUnionLimit limit to the result in the union result
func (p *UnionPlan) handleUnionLimit(res *mysql.Result) error {
	// If the result set is empty, return directly
	if res == nil || res.Resultset == nil {
		return nil
	}

	//Call limitSelectResult to perform the limit operation
	if err := p.limitResult(res); err != nil {
		return fmt.Errorf("union plan failed to limit result: %v", err)
	}
	return nil
}

// extractLimitInfo extracts limit information from a UnionStmt
// offset = 0 means start from the first row of the result set.
// count = -1 means return all rows (no limit)
func (p *UnionPlan) extractLimitInfo(stmt *ast.UnionStmt) (offset int64, count int64, err error) {
	offset, count = 0, -1
	if stmt == nil || stmt.Limit == nil {
		return offset, count, nil
	}

	limit := stmt.Limit

	if limit.Count != nil {
		if countExpr, ok := limit.Count.(*driver.ValueExpr); ok {
			count = countExpr.GetInt64()
		} else {
			return -1, -1, fmt.Errorf("unexpected type for limit count: %T", limit.Count)
		}
	} else {
		count = -1
	}

	if limit.Offset != nil {
		if offsetExpr, ok := limit.Offset.(*driver.ValueExpr); ok {
			offset = offsetExpr.GetInt64()
		} else {
			return -1, -1, fmt.Errorf("unexpected type for limit offset: %T", limit.Offset)
		}
	} else {
		offset = 0
	}
	return offset, count, nil
}

func (p *UnionPlan) removeDuplicateValues(values [][]interface{}, colCnt int) ([][]interface{}, error) {
	distinctKeySet := make(map[string]bool)
	// Change to 0 capacity to correctly initialize newValues
	newValues := make([][]interface{}, 0, len(values))

	for _, row := range values {
		if colCnt > len(row) {
			return nil, fmt.Errorf("columns: %d is greater than the number of columns in the returned row data columns: %d", colCnt, len(row))
		}
		keySlice := row[0:colCnt]
		mk, err := generateMapKey(keySlice)
		if err != nil {
			return newValues, err
		}
		// If the key does not exist yet, record the key and add the current row to the new result set
		if _, ok := distinctKeySet[mk]; !ok {
			distinctKeySet[mk] = true
			newValues = append(newValues, row)
		}
	}

	return newValues, nil
}

// sortSelectResult method sorts the result set according to the ORDER BY clause in the SQL statement.
// It extracts the sort column index and sort direction to construct a list of sort keys, sortKeys.
// Finally, it calls the res.SortWithoutColumnName method to sort the result set according to the constructed sort keys.
func (p *UnionPlan) sortResult(res *mysql.Result) error {
	// Check if the result set is empty
	if res == nil || res.Resultset == nil {
		return nil
	}

	// Check if OrderBy information exists
	if p.stmt == nil || p.stmt.OrderBy == nil {
		return nil
	}

	// Extract OrderBy information
	orderByColumns, orderByDirections, err := p.extractOrderByInfo(res.Fields, p.stmt.OrderBy)
	if err != nil {
		return fmt.Errorf("failed to extract order by info: %v", err)
	}

	// Check if the extracted sort column and direction match
	if len(orderByColumns) != len(orderByDirections) {
		return fmt.Errorf("order by columns and directions length mismatch: columns=%d, directions=%d",
			len(orderByColumns), len(orderByDirections))
	}

	// Pre-allocate the SortKey list to reduce the number of memory allocations
	sortKeys := make([]mysql.SortKey, len(orderByColumns))

	for i := 0; i < len(orderByDirections); i++ {
		sortKey := mysql.SortKey{
			Column:    orderByColumns[i],
			Direction: mysql.SortAsc,
		}
		if orderByDirections[i] {
			sortKey.Direction = mysql.SortDesc
		}
		sortKeys[i] = sortKey
	}

	if err := res.SortWithoutColumnName(sortKeys); err != nil {
		return fmt.Errorf("failed to sort result without column names: %v", err)
	}

	return nil
}

// limitSelectResult intercepts the result set according to the LIMIT clause in the SQL statement.
// It extracts the offset offset and the limit count, verifies their validity, and then calculates the actual interception range of the result set.
// If the offset exceeds the result set length, an empty result is returned;
// otherwise, the rows of the result set are intercepted according to the calculated range.
func (p *UnionPlan) limitResult(ret *mysql.Result) error {
	offset, count, err := p.extractLimitInfo(p.stmt)
	if err != nil {
		return fmt.Errorf("fail to extract limit info: %v", err)
	}

	// Check if offset is valid
	if offset < 0 {
		return fmt.Errorf("invalid offset value: %d", offset)
	}

	// Ensure count is valid, if count is negative and not -1, return an error
	if count < -1 {
		return fmt.Errorf("invalid count value: %d", count)
	}

	// Get the length of the result set
	rowLen := int64(len(ret.Values))

	// If count is -1, it means there is no limit and all rows are taken
	if count == -1 {
		count = rowLen - offset
	}

	// Determine the end position, which cannot exceed the length of the result set
	end := math.MinInt64(offset+count, rowLen)

	// If offset exceeds the result set length, return an empty result
	if offset >= rowLen {
		ret.RowDatas = ret.RowDatas[:0]
		ret.Values = ret.Values[:0]
		return nil
	}

	// Extract the required part
	ret.Values = ret.Values[offset:end]
	return nil
}

// extractOrderByInfo extracts order by information from fields and OrderByClause
func (p *UnionPlan) extractOrderByInfo(fields []*mysql.Field, orderBy *ast.OrderByClause) ([]int, []bool, error) {
	var orderByColumns []int
	var orderByDirections []bool

	for _, item := range orderBy.Items {
		index, err := p.getOrderByColumnIndex(fields, item.Expr)
		if err != nil {
			return nil, nil, err
		}
		orderByColumns = append(orderByColumns, index)
		orderByDirections = append(orderByDirections, item.Desc)
	}

	return orderByColumns, orderByDirections, nil
}

// getOrderByColumnIndex gets the index of the corresponding column in the result field list according to the sort expression.
// For a column name expression (*ast.ColumnNameExpr), it extracts the column name and finds the matching field in the result field list and returns its index.
// If no matching column name is found or an unsupported expression type is encountered, an error is returned.
func (p *UnionPlan) getOrderByColumnIndex(fields []*mysql.Field, expr ast.ExprNode) (int, error) {
	switch x := expr.(type) {
	case *ast.ColumnNameExpr:
		colName := x.Name.Name.L
		for i, field := range fields {
			if string(field.Name) == colName {
				return i, nil
			}
		}
		return -1, fmt.Errorf("column %s not found in result fields", colName)
	case *ast.PositionExpr:
		// PositionExpr.N 是从 1 开始的，需要转换为 Go 中的索引，从 0 开始
		// example: (SELECT id AS identifier FROM employees) UNION (SELECT id FROM contractors) ORDER BY 1;
		idx := x.N - 1
		if idx < 0 || idx >= len(fields) {
			return -1, fmt.Errorf("ORDER BY column index %d out of range", x.N)
		}
		return idx, nil
	default:
		return -1, fmt.Errorf("unsupported ORDER BY expression: %v", expr)
	}
}
