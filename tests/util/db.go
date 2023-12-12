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

// setupDatabaseAndInsertData 创建数据库和表，然后插入数据
func SetupDatabaseAndInsertData(conn *sql.DB, db, table string) error {
	commands := []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", db),
		fmt.Sprintf("CREATE DATABASE %s", db),
		fmt.Sprintf("USE %s", db),
		fmt.Sprintf("CREATE TABLE %s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", table),
	}

	// 执行数据库和表的创建命令
	for _, cmd := range commands {
		if _, err := conn.Exec(cmd); err != nil {
			return err
		}
	}

	// 插入数据
	for i := 0; i < 10; i++ {
		if _, err := conn.Exec(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", table), "nameValue"); err != nil {
			return err
		}
	}

	return nil
}

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
			return fmt.Errorf("column mismatch")
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

// sendReadRequest 发送读请求
func GetRealSeverID(conn *sql.DB) (serverID int, err error) {
	query := "SELECT @@server_id"
	err = conn.QueryRow(query).Scan(&serverID)
	if err != nil {
		return serverID, fmt.Errorf("query server ID Error:%v", err)
	}

	return serverID, nil
}

func LockWithReadLock(db *sql.DB) error {
	_, err := db.Exec("FLUSH TABLES WITH READ LOCK")
	if err != nil {
		return fmt.Errorf("failed to lock slave for read: %v", err)
	}
	return nil
}

func UnLockReadLock(db *sql.DB) error {
	_, err := db.Exec("UNLOCK TABLES")
	if err != nil {
		return fmt.Errorf("failed to unlock slave: %v", err)
	}
	return nil
}
