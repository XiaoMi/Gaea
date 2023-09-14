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

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

const (
	MysqlDriverName = "mysql"
	BaseSqlDir      = "../../sql_cases/"
	TestFilePattern = "test_*"
	SqlStartFile    = "0-prepare.sql"
	SqlEndFile      = "1-clean.sql"
)

func Test_Main(t *testing.T) {
	main()
}

func getConnection(proxyUrl, mysqlUrl string) (*sql.DB, *sql.DB, error) {
	proxyDb, err := sql.Open(MysqlDriverName, proxyUrl)
	if err != nil {
		return nil, nil, err
	}
	proxyDb.SetConnMaxLifetime(time.Minute * 3)
	proxyDb.SetMaxOpenConns(10)
	proxyDb.SetMaxIdleConns(10)

	mysqlDb, err := sql.Open(MysqlDriverName, mysqlUrl)
	if err != nil {
		return nil, nil, err
	}
	mysqlDb.SetConnMaxLifetime(time.Minute * 3)
	mysqlDb.SetMaxOpenConns(10)
	mysqlDb.SetMaxIdleConns(10)
	return proxyDb, mysqlDb, nil
}

func rawToString(raw sql.RawBytes) string {
	if raw == nil {
		return "NULL"
	}

	return string(raw)
}

func checkResult(proxyResult, mysqlresult *sql.Rows) error {
	columns, _ := proxyResult.Columns()
	columnSize := len(columns)

	proxyValues := make([]sql.RawBytes, len(columns))
	proxyArgs := make([]interface{}, columnSize)
	for i := range proxyValues {
		proxyArgs[i] = &proxyValues[i]
	}

	mysqlValues := make([]sql.RawBytes, len(columns))
	mysqlArgs := make([]interface{}, columnSize)
	for i := range mysqlValues {
		mysqlArgs[i] = &mysqlValues[i]
	}

	for {
		b1 := proxyResult.Next()
		if !b1 {
			if mysqlresult.Next() {
				return fmt.Errorf("row number is not equal")
			}
			break
		}

		b2 := mysqlresult.Next()
		if !b2 {
			return fmt.Errorf("row number is not equal")

		}

		if err := proxyResult.Scan(proxyArgs...); err != nil {
			return err
		}

		if err := mysqlresult.Scan(mysqlArgs...); err != nil {
			return err
		}

		for idx := range proxyValues {
			if rawToString(proxyValues[idx]) != rawToString(mysqlValues[idx]) {
				return fmt.Errorf("column value is not equal. col1 = %s, col2 = %s", rawToString(proxyValues[idx]),
					rawToString(mysqlValues[idx]))
			}
		}
	}

	return nil
}

// Test the proxy
func TestIntegration(t *testing.T) {
	// the following code can be refator to a function
	// maybe we should encode the username and password
	proxyUrl := "IT_USER:IT_PASSWORD@tcp(127.0.0.1:13306)/sbtest1"
	mysqlUrl := "IT_USER:IT_PASSWORD@tcp(10.38.164.125:3308)/sbtest1"
	proxyDb, mysqlDb, err := getConnection(proxyUrl, mysqlUrl)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = proxyDb.Close()
		_ = proxyDb.Close()
	}()

	// 执行初始化文件
	if err := execFileInMySQL(BaseSqlDir+SqlStartFile, mysqlDb); err != nil {
		assert.Fail(t, err.Error())
	}

	// 遍历 SQL 测试文件
	sqlFiles, err := filepath.Glob(BaseSqlDir + TestFilePattern)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	for _, fs := range sqlFiles {
		bys, err := ioutil.ReadFile(fs)
		if err != nil {
			fmt.Printf("read file %s err:%s\n", SqlEndFile, err)
			assert.Fail(t, err.Error())
		}

		fmt.Printf("start test file: %s...\n", fs)
		sqls := getSQLFromFile(string(bys))

		for line, sqlString := range sqls {
			if ok, err := isSQLReadOnly(sqlString); err != nil {
				fmt.Printf("parse SQL err:%s\n", sqlString)
			} else if !ok {
				fmt.Printf("SQL not read_only.skip:%s\n", sqlString)
			}
			if err = retryer(proxyDb, mysqlDb, sqlString, doCheck); err != nil {
				fmt.Printf("filename: %s, line number: %d, sql = [%s],err:%s\n", fs, line+1, sqlString, err)
				assert.Fail(t, err.Error())
			}
		}

		fmt.Printf("finish test file: %s\n", fs)
	}

	// 执行清理文件
	if err := execFileInMySQL(BaseSqlDir+SqlEndFile, mysqlDb); err != nil {
		assert.Fail(t, err.Error())
	}
}

func execFileInMySQL(path string, mysqlDb *sql.DB) error {
	fmt.Printf("start exec whole file on in MySQL.file:%s\n", path)
	bys, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	cleanSQLs := getSQLFromFile(string(bys))
	for l, sql := range cleanSQLs {
		if _, err := mysqlDb.Exec(sql); err != nil {
			fmt.Printf("filename: %s, line number: %d, sql = [%s]\n", SqlEndFile, l+1, sql)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isSQLReadOnly(sqlString string) (bool, error) {
	p := parser.New()
	stmt, _, err := p.Parse(sqlString, "", "")
	if err != nil {
		return false, err
	}
	switch stmt[0].(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.ShowStmt:
		return true, err
	default:
		return false, err
	}
}

func retryer(proxyDb, mysqlDb *sql.DB, sqlString string,
	fn func(proxyDb, mysqlDb *sql.DB, sqlString string) error) error {
	var retryTimes = 3
	var err error
	for i := 0; i < retryTimes; i++ {
		if err = fn(proxyDb, mysqlDb, sqlString); err == nil {
			return nil
		}
	}

	return err
}

func doCheck(proxyDb, mysqlDb *sql.DB, sqlString string) error {
	var r1, r2 *sql.Rows
	var err error

	//To make the test more robust, we should add retry methetism
	if r1, err = proxyDb.Query(sqlString); err != nil {
		return err
	}
	if r2, err = mysqlDb.Query(sqlString); err != nil {
		return err
	}

	return checkResult(r1, r2)
}

func getSQLFromFile(sqls string) []string {
	var res []string
	sqlString := strings.Split(string(sqls), "\n")
	for _, sql := range sqlString {
		trimSql := strings.TrimSpace(sql)
		if strings.HasPrefix(trimSql, "//") || strings.HasPrefix(trimSql, "#") || trimSql == "" || strings.HasPrefix(trimSql, "-- ") {
			continue
		}
		res = append(res, sql)
	}
	return res
}
