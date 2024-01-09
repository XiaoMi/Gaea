package util

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"time"
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

// GetDataFromRow retrieves data from a single SQL row and returns it as a slice of strings.
// It handles various SQL data types and converts them into a readable string format.
// If the row contains SQL null values, they are converted to a readable format.
// The function returns an error if the row scan fails.
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

	field := make([]interface{}, 0, numCols)
	for i := 0; i < numCols; i++ {
		field = append(field, new(sql.NullString))
	}

	if err := row.Scan(field...); err != nil {
		if err == sql.ErrNoRows {
			return []string{}, err
		}
		return nil, err
	}
	result := make([]string, 0, numCols)
	for i := 0; i < numCols; i++ {
		col := pr(field[i])
		result = append(result, col)
	}

	return result, nil
}

// GetDataFromResult retrieves the number of rows affected by an SQL operation.
// It is useful for understanding the impact of non-query SQL commands like INSERT, UPDATE, or DELETE.
// The function returns the number of rows affected and an error if the operation fails.
func GetDataFromResult(result sql.Result) (int64, error) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// MysqlQuery executes a given SQL query on the provided database connection and returns the results.
// The results are returned as a slice of string slices, where each inner slice represents a row.
// The function handles no-row scenarios and returns an error for other query failures.
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

// MysqlExec executes a given SQL command (like INSERT, UPDATE, DELETE) on the provided database connection.
// It returns the number of rows affected by the command and an error if the command execution fails.
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

// IsSlice function takes an interface{} type argument i and returns a boolean
// indicating whether i is of a slice type.
func IsSlice(i interface{}) bool {
	// Use reflect.TypeOf to get the reflection Type object of i.
	// If i is nil, reflect.TypeOf(i) will return nil,
	// so it's necessary to check for nil before calling the .Kind() method.
	t := reflect.TypeOf(i)
	if t == nil {
		return false
	}

	// The .Kind() method returns the specific kind of the type.
	// reflect.Slice is a constant representing the slice type.
	// If the kind of i is a slice, return true; otherwise, return false.
	return t.Kind() == reflect.Slice
}

// GetSqlFromFile reads SQL commands from a file and returns them as a slice of strings.
// It ignores lines that are comments (starting with "//", "#", or "--").
// The function returns an error if reading the file fails.
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

// CompareIgnoreSort compares two slices of string slices (representing SQL query results) without considering the order of rows.
// It is useful for comparing query results where the row order is not guaranteed.
// The function returns true if the results are equivalent, and false with an error if they are not.
func CompareIgnoreSort(gaeaRes [][]string, mysqlRes [][]string) (bool, error) {
	if len(gaeaRes) != len(mysqlRes) {
		return false, fmt.Errorf("sql.Result mismatched lengths for results. gaeaRes: %v, mysqlRes: %v", gaeaRes, mysqlRes)
	}

	elementCount := make(map[string]int)
	for _, item := range gaeaRes {
		elementCount[fmt.Sprint(item)]++
	}

	for _, item := range mysqlRes {
		elementCount[fmt.Sprint(item)]--
		if elementCount[fmt.Sprint(item)] < 0 {
			return false, fmt.Errorf("sql.Result mismatched elements for results. gaeaRes: %v, mysqlRes: %v", gaeaRes, mysqlRes)
		}
	}

	for _, count := range elementCount {
		if count != 0 {
			return false, fmt.Errorf("sql.Result mismatched elements for results. gaeaRes: %v, mysqlRes: %v", gaeaRes, mysqlRes)
		}
	}

	return true, nil
}

// SetupDatabaseAndInsertData sets up a database and a table, and then inserts data into the table.
// It first drops the database if it already exists, creates a new database, and then creates a table within that database.
// After setting up the database and table, it inserts a predefined number of rows into the table.
// The function takes a database connection, and the names of the database and table as parameters.
// It returns an error if any of the database setup or data insertion commands fail.
func SetupDatabaseAndInsertData(conn *sql.DB, db, table string) error {
	commands := []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", db),
		fmt.Sprintf("CREATE DATABASE %s", db),
		fmt.Sprintf("CREATE TABLE %s.%s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", db, table),
	}

	// 执行数据库和表的创建命令
	for _, cmd := range commands {
		if _, err := conn.Exec(cmd); err != nil {
			return fmt.Errorf("failed to execute command '%s': %v", cmd, err)
		}
	}

	// 插入数据
	for i := 0; i < 10; i++ {
		if _, err := conn.Exec(fmt.Sprintf("INSERT INTO %s.%s (name) VALUES (?)", db, table), "nameValue"); err != nil {
			return fmt.Errorf("failed to insert data into table %s: %v", table, err)
		}
	}

	return nil
}

// CompareQueryRows executes a given query on two different database connections (db1 and db2)
// and compares the results to check if they are identical. It compares both the structure (columns)
// and the content (rows) of the query result sets.
// The function returns nil if the results are identical, or an error if there are any differences
// or issues encountered during the execution or comparison process.
func CompareQueryRows(db1 *sql.DB, db2 *sql.DB, query string) error {
	rows1, err := db1.Query(query)
	if err != nil {
		return err
	}
	defer rows1.Close()

	rows2, err := db2.Query(query)
	if err != nil {
		return err
	}
	defer rows2.Close()
	// 检查第一个结果集是否为空
	hasRow1 := rows1.Next()
	hasRow2 := rows2.Next()

	// 如果两个结果集都为空，那么它们是相同的
	if !hasRow1 && !hasRow2 {
		return nil
	}

	// 如果一个结果集为空而另一个不为空，则它们不同
	if hasRow1 != hasRow2 {
		return fmt.Errorf("one of the result sets is empty while the other is not")
	}

	for {
		// 两个结果集都至少有一行数据
		cols1, err := rows1.Columns()
		if err != nil {
			return err
		}
		cols2, err := rows2.Columns()
		if err != nil {
			return err
		}

		// 比较列的数量和名称
		if len(cols1) != len(cols2) || !reflect.DeepEqual(cols1, cols2) {
			return fmt.Errorf("column mismatch,res1:%+v,res2:%+v", cols1, cols2)
		}

		// 创建接收数据的切片
		vals1 := make([]interface{}, len(cols1))
		vals2 := make([]interface{}, len(cols2))
		for i := range cols1 {
			vals1[i] = new(sql.RawBytes)
			vals2[i] = new(sql.RawBytes)
		}

		// 扫描数据
		if err := rows1.Scan(vals1...); err != nil {
			return err
		}
		if err := rows2.Scan(vals2...); err != nil {
			return err
		}

		// 比较每一列的值
		if !reflect.DeepEqual(vals1, vals2) {
			return fmt.Errorf("row values mismatch,sql:%s,vals1:%v ,vals2:%v ", query, vals1, vals2)
		}

		// 检查是否还有更多的行
		hasRow1 = rows1.Next()
		hasRow2 = rows2.Next()
		if hasRow1 != hasRow2 {
			return fmt.Errorf("number of rows mismatch")
		}
		if !hasRow1 {
			// 如果没有更多的行，则完成比较
			break
		}
	}

	return nil
}

// GetSecondsBehindMaster retrieves the replication delay (in seconds) of a MySQL slave server from its master.
// This function executes the 'SHOW SLAVE STATUS' SQL command and parses the output to find the 'Seconds_Behind_Master' value.
// It returns the number of seconds the slave is behind the master, or an error if the query fails, the value is null, or the field is not found.
func GetSecondsBehindMaster(db *sql.DB) (secondsBehindMaster int, err error) {
	rows, err := db.Query(`SHOW SLAVE STATUS`)
	if err != nil {
		return -1, fmt.Errorf("query slave status error:%v", err)
	}
	defer rows.Close()
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
	for rows.Next() {
		if err := rows.Scan(field...); err != nil {
			return -1, err
		}
		for i := 0; i < n; i++ {
			col := pr(field[i])
			if c[i] == "Seconds_Behind_Master" {
				if col == "\\N" {
					col = "NULL" // 期望的其他值
					return -1, fmt.Errorf("Seconds_Behind_Master is null")
				}
				delay, err := strconv.Atoi(col)
				return delay, err
			}
			if col == "\\N" {
				col = "NULL" // 或者期望的其他值
			}
		}
	}
	return -1, fmt.Errorf("can not find Seconds Behind Master")
}

// GetRealServerID retrieves the server ID of the MySQL instance connected through the given database connection.
// It executes a query to fetch the server ID and returns it.
// If the query fails, it returns an error with a descriptive message.
func GetRealSeverID(conn *sql.DB) (serverID int, err error) {
	query := "SELECT @@server_id"
	err = conn.QueryRow(query).Scan(&serverID)
	if err != nil {
		return serverID, fmt.Errorf("query server ID Error:%v", err)
	}

	return serverID, nil
}

// LockWithReadLock executes a 'FLUSH TABLES WITH READ LOCK' SQL command on the provided database connection.
// This command is used to lock all tables for read operations, which is useful for creating consistent backups
// or synchronizing replication slaves with the master.
// It returns an error if the locking operation fails.
func LockWithReadLock(db *sql.DB) error {
	_, err := db.Exec("FLUSH TABLES WITH READ LOCK")
	if err != nil {
		return fmt.Errorf("failed to lock slave for read: %v", err)
	}
	return nil
}

// UnLockReadLock releases any table locks held by the current session on the provided database connection.
// This function is typically used to unlock the tables that were locked using the 'FLUSH TABLES WITH READ LOCK' command.
// It returns an error if the unlocking operation fails.
func UnLockReadLock(db *sql.DB) error {
	_, err := db.Exec("UNLOCK TABLES")
	if err != nil {
		return fmt.Errorf("failed to unlock slave: %v", err)
	}
	return nil
}

// CleanUpDatabases removes all non-system databases from the MySQL server.
// This function is intended to be used for cleaning up a database server by
// dropping all user-created databases while preserving the essential system databases.
// System databases like 'information_schema', 'mysql', 'performance_schema', and 'sys'
// are not affected by this operation.
func CleanUpDatabases(db *sql.DB) error {
	// 获取非系统数据库列表
	rows, err := db.Query("SHOW DATABASES WHERE `Database` NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')")
	if err != nil {
		return fmt.Errorf("error querying non-system databases: %w", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			return fmt.Errorf("error scanning database result: %w", err)
		}
		databases = append(databases, database)
	}

	// 检查是否在获取数据库列表时出错
	if err = rows.Err(); err != nil {
		return fmt.Errorf("error after iterating over rows: %w", err)
	}

	// 清理非系统数据库
	for _, database := range databases {
		_, err := db.Exec(fmt.Sprintf("DROP DATABASE `%s`", database))
		if err != nil {
			return fmt.Errorf("error dropping database %s: %w", database, err)
		}
	}
	return nil
}
