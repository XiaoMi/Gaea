package util

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/XiaoMi/Gaea/parser"
)

// This function does not call rows.Close()
// so the release of resources should be managed by external code that calls this function.
// When using this function, be sure to call rows.Close() externally to avoid leaking database resources.
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
				col = "NULL" // 或者期望的其他值
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

func GetDataFromResult(result sql.Result) (int64, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func MysqlQuery(db *sql.DB, execSql string) ([][]string, error) {
	rows, err := db.Query(execSql)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, but it's not an actual error, so we return nil
			return nil, nil
		}
		// An error occurred while querying
		return nil, err
	}
	defer rows.Close()
	res, err := GetDataFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("get data from rows Error: %v", err)
	}
	return res, nil
}

func MysqlExec(db *sql.DB, execSql string) (int64, error) {
	result, err := db.Exec(execSql)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows affected but not an error, so return 0
			return 0, nil
		}
		return -1, fmt.Errorf("sql %s Exec Error: %v", execSql, err)
	}
	rowsAffected, err := GetDataFromResult(result)
	if err != nil {
		return 0, fmt.Errorf("get data from to result Error: %v", err)
	}
	return rowsAffected, nil
}

func IsSlice(i interface{}) bool {
	return reflect.TypeOf(i).Kind() == reflect.Slice
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

func GetSqlFromFile(path string) ([]string, error) {
	bys, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file error")
	}
	var res []string
	sqlString := strings.Split(string(bys), "\n")
	for _, sql := range sqlString {
		trimSql := strings.TrimSpace(sql)
		if strings.HasPrefix(trimSql, "//") || strings.HasPrefix(trimSql, "#") || trimSql == "" || strings.HasPrefix(trimSql, "-- ") || strings.HasPrefix(trimSql, "--") {
			continue
		}
		res = append(res, sql)
	}
	return res, nil
}

func VerifySqlParsable(sqlString string) error {
	p := parser.New()
	_, _, err := p.Parse(sqlString, "", "")
	if err != nil {
		return fmt.Errorf("sql parsing error: %v, SQL: %s", err, sqlString)
	}
	return nil
}

func OutPutResult(sql string, outPutFile string) error {
	// 尝试打开或创建文件
	file, err := os.OpenFile(outPutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file %s: %v", outPutFile, err)
	}
	defer file.Close()

	// 将SQL语句写入文件
	_, err = file.WriteString(sql + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func Compare(res1 [][]string, res2 [][]string) (bool, error) {
	if reflect.DeepEqual(res1, res2) {
		return true, nil
	} else {
		return false, fmt.Errorf("sql.Result mismatched types for results. res1: %v, res2: %v", res1, res2)
	}
}
