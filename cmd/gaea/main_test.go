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
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"time"
)

const MYSQL_DRIVER_NAME = "mysql"

func Test_Main(t *testing.T) {
	main()
}

func getConnection(proxyUrl, mysqlUrl string) (*sql.DB, *sql.DB, error) {
	proxyDb, err := sql.Open(MYSQL_DRIVER_NAME, proxyUrl)
	if err != nil {
		return nil, nil, err
	}
	proxyDb.SetConnMaxLifetime(time.Minute * 3)
	proxyDb.SetMaxOpenConns(10)
	proxyDb.SetMaxIdleConns(10)

	mysqlDb, err := sql.Open(MYSQL_DRIVER_NAME, mysqlUrl)
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

func checkResult(proxyResult, mysqlresult *sql.Rows, t *testing.T) error {
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
				assert.Fail(t, "row number is not equal")
			}
			break
		}

		b2 := mysqlresult.Next()
		if !b2 {
			assert.Fail(t, "row number is not equal")
			break
		}

		if err := proxyResult.Scan(proxyArgs...); err != nil {
			return err
		}

		if err := mysqlresult.Scan(mysqlArgs...); err != nil {
			return err
		}

		for idx := range proxyValues {
			assert.Equal(t, rawToString(proxyValues[idx]), rawToString(mysqlValues[idx]))
		}
	}

	return nil
}

// Test the proxy
func TestIntegration(t *testing.T) {
	// the following code can be refator to a function
	// maybe we should encode the username and password
	proxyUrl := "PROXY_URL"
	mysqlUrl := "MYSQL_URL"
	proxyDb, mysqlDb, err := getConnection(proxyUrl, mysqlUrl)
	if err != nil {
		panic(err)
	}

	defer func() {
		proxyDb.Close()
		proxyDb.Close()
	}()

	//
	sqlFiles, err := ioutil.ReadDir("../../sql_cases")
	if err != nil {
		assert.Fail(t, err.Error())
	}

	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, fs := range sqlFiles {
		bys, err := ioutil.ReadFile("../../sql_cases/" + fs.Name())
		if err != nil {
			assert.Fail(t, err.Error())
		}

		for _, sqlString := range strings.Split(string(bys), "\n") {
			trimSql := strings.TrimSpace(sqlString)
			if strings.HasPrefix(trimSql, "//") || trimSql == "" {
				continue
			}

			//This is schema sql, we should run it first
			if strings.Contains(fs.Name(), "scheam") {
				if _, err := mysqlDb.Exec(sqlString); err != nil {
					assert.Fail(t, err.Error())
				}
			}

			if err = retryer(t, proxyDb, mysqlDb, sqlString, doCheck); err != nil {
				assert.Fail(t, err.Error())
			}
		}
	}
}

func retryer(t *testing.T, proxyDb, mysqlDb *sql.DB, sqlString string,
	fn func(t *testing.T, proxyDb, mysqlDb *sql.DB, sqlString string) error) error {
	var retryTimes = 3
	var err error
	for i := 0; i < retryTimes; i++ {
		if err = fn(t, proxyDb, mysqlDb, sqlString); err == nil {
			return nil
		}
	}

	return err
}

func doCheck(t *testing.T, proxyDb, mysqlDb *sql.DB, sqlString string) error {
	var r1, r2 *sql.Rows
	var err error

	//To make the test more robust, we should add retry methetism
	if r1, err = proxyDb.Query(sqlString); err != nil {
		return err
	}
	if r2, err = mysqlDb.Query(sqlString); err != nil {
		return err
	}

	return checkResult(r1, r2, t)
}
