package util

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 处理后端DB的一些问题

func Open(host string, port int, User string, Password string, connectTimeout int, readTimeout int) (*sql.DB, error) {
	mysqlUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=%ds&readTimeout=%ds&interpolateParams=true",
		User,
		Password,
		host, port,
		connectTimeout,
		readTimeout,
	)
	db, err := sql.Open("mysql", mysqlUri)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db, err
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
