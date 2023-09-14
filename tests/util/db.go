package util

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
)

func GetDataFromRows(rows *sql.Rows) ([][]string, error) {
	pr := func(t interface{}) (r string) {
		r = "\\N"
		switch v := t.(type) {
		case *sql.NullBool:
			if v.Valid {
				r = strconv.FormatBool(v.Bool)
			}
		case *sql.NullString:
			if v.Valid {
				r = v.String
			}
		case *sql.NullInt64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Int64)
			}
		case *sql.NullFloat64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Float64)
			}
		case *time.Time:
			if v.Year() > 1900 {
				r = v.Format("2006-01-02 15:04:05")
			}
		default:
			r = fmt.Sprintf("%#v", t)
		}
		return
	}
	c, _ := rows.Columns()
	n := len(c)
	field := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		field = append(field, new(sql.NullString))
	}
	var converts [][]string
	for rows.Next() {
		if err := rows.Scan(field...); err != nil {
			return nil, err
		}
		row := make([]string, 0, n)
		for i := 0; i < n; i++ {
			col := pr(field[i])
			if col == "\\N" {
				col = "NULL" // 或者你期望的其他值
			}
			row = append(row, col)
		}
		converts = append(converts, row)
	}
	return converts, nil
}

func GetDataFromRow(row *sql.Row, numCols int) ([]string, error) {
	pr := func(t interface{}) (r string) {
		r = "\\N"
		switch v := t.(type) {
		case *sql.NullBool:
			if v.Valid {
				r = strconv.FormatBool(v.Bool)
			}
		case *sql.NullString:
			if v.Valid {
				r = v.String
			}
		case *sql.NullInt64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Int64)
			}
		case *sql.NullFloat64:
			if v.Valid {
				r = fmt.Sprintf("%v", v.Float64)
			}
		case *time.Time:
			if v.Year() > 1900 {
				r = v.Format("2006-01-02 15:04:05")
			}
		default:
			r = fmt.Sprintf("%#v", t)
		}
		return
	}

	// Note: The sql.Row does not have a method to get columns directly
	// You'll need to know the number of columns (numCols) beforehand
	field := make([]interface{}, 0, numCols)
	for i := 0; i < numCols; i++ {
		field = append(field, new(sql.NullString))
	}

	if err := row.Scan(field...); err != nil {
		if err == sql.ErrNoRows {
			// Return nil or []string{} as per your requirement
			return []string{}, err
		}
		return nil, err
	}
	result := make([]string, 0, numCols)
	for i := 0; i < numCols; i++ {
		col := pr(field[i])
		result = append(result, col)
	}

	// As the sql.Row doesn't have a .Columns() method like sql.Rows, you'll need to know the column names from elsewhere
	// In this example, I've just returned an empty slice for column names
	return result, nil
}

func IsSlice(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Slice
}

func CreateDatabaseAndTables(db *sql.DB, databaseName string, tables []string) error {
	// 删除数据库 (如果存在)
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", databaseName)
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("DROP DATABASE Error: %w", err)
	}

	// 创建数据库
	sql = fmt.Sprintf("CREATE DATABASE %s", databaseName)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("CREATE DATABASE Error: %w", err)
	}

	// 使用数据库
	sql = fmt.Sprintf("USE %s", databaseName)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("USE DATABASE Error: %w", err)
	}

	// 创建表
	for _, v := range tables {
		_, err = db.Exec(v)
		if err != nil {
			return fmt.Errorf("CREATE Table Error: %w", err)
		}
	}
	return nil
}

func InsertData(db *sql.DB, databaseName string, tableName string, numRows int) error {
	// 准备插入数据的 SQL 语句
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s.%s (name) VALUES (?)", databaseName, tableName))
	if err != nil {
		return fmt.Errorf("PREPARE INSERT Error: %w", err)
	}
	defer stmt.Close()

	// 循环插入指定数量的数据
	for i := 1; i <= numRows; i++ {
		_, err = stmt.Exec(fmt.Sprintf("John%d", i))
		if err != nil {
			return fmt.Errorf("EXEC INSERT Error: %w", err)
		}
	}
	return nil
}

func DeleteRow(db *sql.DB, databaseName string, tableName string, id int) error {
	// Prepare the delete SQL statement
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s.%s WHERE id=?", databaseName, tableName))
	if err != nil {
		return fmt.Errorf("PREPARE DELETE Error: %w", err)
	}
	defer stmt.Close()

	// Execute the delete
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("EXEC DELETE Error: %w", err)
	}

	return nil
}

func DeleteVerify(db *sql.DB, databaseName string, tableName string, id int) error {
	// Prepare the select SQL statement
	stmt, err := db.Prepare(fmt.Sprintf("SELECT name FROM %s.%s WHERE id=?", databaseName, tableName))
	if err != nil {
		return fmt.Errorf("PREPARE SELECT Error: %w", err)
	}

	// Try to execute the select to fetch the deleted value
	var deletedValue string
	err = stmt.QueryRow(id).Scan(&deletedValue)

	// Since the row was deleted, we expect an error indicating no rows in result set
	if err == nil || !strings.Contains(err.Error(), "no rows in result set") {
		return fmt.Errorf("QUERY SELECT Error: expected no rows in result set, got %w", err)
	}

	return nil
}

func UpdateRow(db *sql.DB, databaseName string, tableName string, id int, newValue string) error {
	// Prepare the update SQL statement
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s.%s SET name=? WHERE id=?", databaseName, tableName))
	if err != nil {
		return fmt.Errorf("PREPARE UPDATE Error: %w", err)
	}
	defer stmt.Close()

	// Execute the update
	_, err = stmt.Exec(newValue, id)
	if err != nil {
		return fmt.Errorf("EXEC UPDATE Error: %w", err)
	}

	return nil
}

func UpdateVerify(db *sql.DB, databaseName string, tableName string, id int, expectedValue string) error {
	// Prepare the select SQL statement
	stmt, err := db.Prepare(fmt.Sprintf("SELECT name FROM %s.%s WHERE id=?", databaseName, tableName))
	if err != nil {
		return fmt.Errorf("PREPARE SELECT Error: %w", err)
	}
	// Execute the select to fetch the updated value
	var updatedValue string
	err = stmt.QueryRow(id).Scan(&updatedValue)
	if err != nil {
		return fmt.Errorf("QUERY SELECT Error: %w", err)
	}
	// Verify the updated value
	if updatedValue != expectedValue {
		return fmt.Errorf("VERIFY UPDATE Error: expected %s, got %s", expectedValue, updatedValue)
	}

	return nil
}

func InsertRow(db *sql.DB, databaseName string, tableName string, newValue string) (int64, error) {
	// Prepare the insert SQL statement
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s.%s (name) VALUES (?)", databaseName, tableName))
	if err != nil {
		return 0, fmt.Errorf("PREPARE INSERT Error: %w", err)
	}
	defer stmt.Close()

	// Execute the insert
	res, err := stmt.Exec(newValue)
	if err != nil {
		return 0, fmt.Errorf("EXEC INSERT Error: %w", err)
	}

	// Get the ID of the inserted row
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("GET LAST INSERT ID Error: %w", err)
	}

	return id, nil
}

func InsertVerify(db *sql.DB, databaseName string, tableName string, id int64, expectedValue string) error {
	// Prepare the select SQL statement
	stmt, err := db.Prepare(fmt.Sprintf("SELECT name FROM %s.%s WHERE id=?", databaseName, tableName))
	if err != nil {
		return fmt.Errorf("PREPARE SELECT Error: %w", err)
	}

	// Execute the select to fetch the inserted value
	var insertedValue string
	err = stmt.QueryRow(id).Scan(&insertedValue)
	if err != nil {
		return fmt.Errorf("QUERY SELECT Error: %w", err)
	}

	// Verify the inserted value
	if insertedValue != expectedValue {
		return fmt.Errorf("VERIFY INSERT Error: expected %s, got %s", expectedValue, insertedValue)
	}

	return nil
}

func ExecuteTransactionAndReturnId(db *sql.DB, sql string, args ...interface{}) (int64, error) {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("BEGIN TRANSACTION Error: %w", err)
	}

	// Prepare the SQL statement
	stmt, err := tx.Prepare(sql)
	if err != nil {
		// If an error occurred while preparing the statement, rollback the transaction
		_ = tx.Rollback()
		return 0, fmt.Errorf("PREPARE TRANSACTION Error: %w", err)
	}
	defer stmt.Close()

	// Execute the SQL statement with provided arguments
	res, err := stmt.Exec(args...)
	if err != nil {
		// If an error occurred while executing the statement, rollback the transaction
		_ = tx.Rollback()
		return 0, fmt.Errorf("EXEC TRANSACTION Error: %w", err)
	}

	// Get the ID of the last inserted row
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("GET LAST INSERT ID Error: %w", err)
	}

	// If no errors occurred, commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("COMMIT TRANSACTION Error: %w", err)
	}

	return id, nil
}

func DBExec(db *sql.DB, execSql string, expectType string) (interface{}, error) {
	switch expectType {
	case "Query":
		rows, err := db.Query(execSql)
		if err == sql.ErrNoRows {
			return rows, nil
		}
		if err != nil {
			return nil, fmt.Errorf("SQL: %s, ExpectType: %s, Error: %v", execSql, expectType, err)
		}
		return rows, nil
	case "Exec":
		result, err := db.Exec(execSql)
		if err == sql.ErrNoRows {
			return result, nil
		}
		if err != nil {
			return nil, fmt.Errorf("SQL: %s, ExpectType: %s, Error: %v", execSql, expectType, err)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("SQL: %s, ExpectType: %s, Error: unsupported expect type", execSql, expectType)
	}
}

// Define an enum-like type for SQL operation methods
type SQLOperation string

var ErrSkippedSQLOperation = errors.New("skipped SQL operation")
var ErrUnsupportedSQLOperation = errors.New("unsupported SQL operation type")

const (
	Query   SQLOperation = "Query"
	Exec    SQLOperation = "Exec"
	UnKnown SQLOperation = "UnKnown"
	Comment SQLOperation = "Comment"
)

func GetSqlFromFile(path string) ([]string, error) {
	bys, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file error")
	}
	var res []string
	sqlString := strings.Split(string(bys), "\n")
	for _, sql := range sqlString {
		trimSql := strings.TrimSpace(sql)
		if strings.HasPrefix(trimSql, "//") || strings.HasPrefix(trimSql, "#") || trimSql == "" || strings.HasPrefix(trimSql, "-- ") {
			continue
		}
		res = append(res, sql)
	}
	return res, nil
}

func GetSqlType(sqlString string) (SQLOperation, error) {
	p := parser.New()
	stmt, _, err := p.Parse(sqlString, "", "")
	if err != nil {
		return UnKnown, err
	}
	switch stmt[0].(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.ShowStmt, *ast.ExplainStmt:
		return Query, nil
	case *ast.InsertStmt:
		return Exec, nil
	case *ast.UpdateStmt:
		return Exec, nil
	case *ast.DeleteStmt:
		return Exec, nil
	case *ast.CreateDatabaseStmt, *ast.CreateTableStmt:
		return Exec, nil
	case *ast.AlterTableStmt:
		return Exec, nil
	case *ast.DropDatabaseStmt, *ast.DropTableStmt:
		return Exec, nil
	case *ast.UseStmt:
		return Exec, nil
	default:
		return UnKnown, nil
	}
}

func OutPutResult(sql string, outPutFile string) error {
	file, err := os.OpenFile(outPutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(sql + "\n")
	return err
}

func ExecuteSQL(db *sql.DB, operationType SQLOperation, sqlStatement string) (interface{}, error) {
	switch operationType {
	case Query:
		return DBExec(db, sqlStatement, "Query")
	case Exec:
		return DBExec(db, sqlStatement, "Exec")
	case UnKnown:
		return nil, ErrSkippedSQLOperation
	default:
		return nil, nil
	}
}

func CompareDBExecResults(result1 interface{}, result2 interface{}) (bool, error) {
	switch r1 := result1.(type) {
	case *sql.Rows:
		r2, ok := result2.(*sql.Rows)
		if !ok {
			return false, errors.New("mismatched types for results")
		}
		defer r1.Close()
		defer r2.Close()
		t1, err := GetDataFromRows(r1)
		if err != nil {
			return false, err
		}
		t2, err := GetDataFromRows(r2)
		if err != nil {
			return false, err
		}
		return reflect.DeepEqual(t1, t2), nil
	case sql.Result:
		r2, ok := result2.(sql.Result)
		if !ok {
			return false, errors.New("mismatched types for results")
		}
		rowsAffected1, err1 := r1.RowsAffected()
		if err1 != nil {
			return false, err1
		}
		rowsAffected2, err2 := r2.RowsAffected()
		if err2 != nil {
			return false, err2
		}
		// Compares affected rows
		return rowsAffected1 == rowsAffected2, nil
	default:
		return false, errors.New("unsupported result type")
	}
}
