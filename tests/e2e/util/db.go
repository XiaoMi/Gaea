package util

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	_ "github.com/go-sql-driver/mysql"
)

// ConnectAndCreateDB 连接到 MySQL 数据库并创建一个新的数据库
func ConnectMysql(config *config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&parseTime=%v&loc=%s&charset=%s",
		config.User, config.Password, config.Url, config.Timeout, config.ReadTimeout, config.WriteTimeout, config.ParseTime, config.Loc, config.Charset)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db, nil
}

func ConnectGaea(config *config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/?timeout=%ds&readTimeout=%ds&interpolateParams=true",
		config.User, config.Password, config.Url, config.Timeout, config.ReadTimeout)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db, nil
}

func CreateDatabases(db *sql.DB, databases []string) error {
	if db == nil {
		return fmt.Errorf("sql.DB is nil")
	}
	for _, dbName := range databases {
		_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database %s: %v", dbName, err)
		}
	}
	return nil
}

func DropDatabase(db *sql.DB, databases []string) error {
	if db == nil {
		return fmt.Errorf("sql.DB is nil")
	}
	for _, dbName := range databases {
		_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to drop database %s: %v", dbName, err)
		}
	}
	return nil
}

func GenCreateTableSqls(database string, tableSQl []string) []string {
	sqls := []string{}
	for _, t := range tableSQl {
		s := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (
			id int(11) NOT NULL AUTO_INCREMENT,
			name varchar(20) DEFAULT NULL,
			PRIMARY KEY (id)
		  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`, database, t)
		sqls = append(sqls, s)
	}
	return sqls
}

func GenInsertDataInTableSqls(dbName string, tableName string, num int) []string {
	sqls := []string{}
	for i := 0; i < num; i++ {
		sqls = append(sqls, fmt.Sprintf("INSERT INTO `%s`.`%s` (id ,name) VALUES (?,?)", dbName, tableName))
	}
	return sqls
}

func GenDropTablesSqls(database string, tableSQl []string) []string {
	var sqls []string
	for _, t := range tableSQl {
		s := fmt.Sprintf(`DROP TABLE IF EXISTS %s.%s`, database, t)
		sqls = append(sqls, s)
	}
	return sqls
}

func Exec(db *sql.DB, database string, sql string) (sql.Result, error) {
	if db == nil {
		return nil, fmt.Errorf("sql.DB is nil")
	}
	_, err := db.Exec(fmt.Sprintf("USE %s", database))
	if err != nil {
		return nil, fmt.Errorf("failed to switch to database %s: %v", database, err)
	}
	result, err := db.Exec(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to exec sql %s: %v", sql, err)
	}
	return result, nil

}

func Query(db *sql.DB, database string, query string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, fmt.Errorf("sql.DB is nil")
	}
	_, err := db.Exec(fmt.Sprintf("USE %s", database))
	if err != nil {
		return nil, fmt.Errorf("failed to switch to database %s: %v", database, err)
	}
	result, err := db.Query(query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to exec sql %s: %v", query, err)
	}
	return result, nil
}

func CheckDBRes(db1 *sql.DB, db2 *sql.DB, sql string, checkFunc func(rows1 *sql.Rows, rows2 *sql.Rows) bool) bool {
	if db1 == nil || db2 == nil {
		return false
	}
	rows1, err1 := db1.Query(sql)
	rows2, err2 := db2.Query(sql)
	if err1 != err2 {
		fmt.Printf("inequality of error.err1:%v, err2:%v\n", err1, err2)
		return false
	}

	if err1 != nil {
		fmt.Printf("error.err1:%s", err1)
		return false
	}

	return checkFunc(rows1, rows2)
}

func CreateDatabaseAndInsertData(db *sql.DB, databaseName string, tableName string, numRows int) error {
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
	sql = fmt.Sprintf("CREATE TABLE %s (id INT AUTO_INCREMENT, name VARCHAR(20), PRIMARY KEY (id))", tableName)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("CREATE Table Error: %w", err)
	}

	// 插入数据到新创建的表
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", tableName))
	if err != nil {
		return fmt.Errorf("PREPARE INSERT Error: %w", err)
	}
	defer stmt.Close()

	for i := 1; i <= numRows; i++ {
		_, err = stmt.Exec(fmt.Sprintf("John%d", i))
		if err != nil {
			return fmt.Errorf("EXEC INSERT Error: %w", err)
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

func QueryGeneralLogArgument(rows *sql.Rows) ([]string, error) {
	var result []string
	var eventTime, userHost, commandType, argument string
	var threadID, serverID int

	for rows.Next() {
		err := rows.Scan(&eventTime, &userHost, &threadID, &serverID, &commandType, &argument)
		if err != nil {
			return nil, err
		}
		result = append(result, argument)
	}
	return result, nil
}

func ContainsTargetValue(slice []string, target string) bool {
	for _, str := range slice {
		if strings.Contains(str, target) {
			return true
		}
	}
	return false
}
