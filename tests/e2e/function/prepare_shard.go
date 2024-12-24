// Copyright 2024 The Gaea Authors. All Rights Reserved.
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

package function

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

type tbl_mycat struct {
	ID          int
	Name        string
	Age         int
	Salary      decimal.Decimal // DECIMAL 类型
	Rating      float64         // FLOAT 类型（使用 float64 来处理）
	IsActive    bool
	CreatedAt   []uint8 // TIME 类型 MySQL存储的是[]uint8 类型
	Description string
	UpdatedAt   []uint8
}

// The test purpose of "Sharded CRUD Test with Multiple MySQL Types including Decimal" is to verify that CRUD operations of different data types (including Decimal) can be executed correctly in a sharded database environment to ensure the following key behaviors:
// Consistency of sharded data:
// Verify that databases and tables are created correctly across multiple shards and ensure that data read and write operations are routed to the correct shards.
// CRUD operations across shards:
// Verify that insert, query, update, delete operations on multiple data types such as Decimal, Float, Datetime, etc. are as expected.
var _ = ginkgo.Describe("Sharded CRUD Test with Multiple MySQL Types including Decimal", func() {
	var (
		masterAdminConnShard1 *sql.DB
		masterAdminConnShard2 *sql.DB
		err                   error
	)
	e2eMgr := config.NewE2eManager()
	nsTmpl := config.ShardNamespaceTmpl
	shard1DBs := []string{"db_mycat_mod_0", "db_mycat_mod_1"}
	shard2DBs := []string{"db_mycat_mod_2", "db_mycat_mod_3"}

	// 表结构 SQL 模板
	createTableSQLTemplate := `
		CREATE TABLE IF NOT EXISTS %s.tbl_mycat (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255),
			age INT ,
			salary DECIMAL(10,2) ,
			rating FLOAT ,
			is_active BOOLEAN,
			created_at DATETIME,
			description TEXT,
			updated_at DATETIME
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	sliceConfig := e2eMgr.NsSlices[config.SliceDualMaster]

	// BeforeEach: 设置多个分片上的数据库和表
	ginkgo.BeforeEach(func() {
		// 获取连接
		masterAdminConnShard1, err = sliceConfig.GetMasterAdminConn(0)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = util.CleanUpDatabases(masterAdminConnShard1)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		masterAdminConnShard2, err = sliceConfig.GetMasterAdminConn(1)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = util.CleanUpDatabases(masterAdminConnShard2)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		time.Sleep(1 * time.Second)

		// 初始化分片 1 的数据库和表
		for _, dbName := range shard1DBs {
			dbSQL := fmt.Sprintf("CREATE DATABASE %s ", dbName)
			tableSQL := fmt.Sprintf(createTableSQLTemplate, dbName)
			_, err := masterAdminConnShard1.Exec(dbSQL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			time.Sleep(100 * time.Millisecond)
			_, err = masterAdminConnShard1.Exec(tableSQL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("Failed to create table in database %s", dbName))
		}

		// 初始化分片 2 的数据库和表
		for _, dbName := range shard2DBs {
			dbSQL := fmt.Sprintf("CREATE DATABASE %s ", dbName)
			tableSQL := fmt.Sprintf(createTableSQLTemplate, dbName)
			_, err := masterAdminConnShard2.Exec(dbSQL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			time.Sleep(100 * time.Millisecond)
			_, err = masterAdminConnShard2.Exec(tableSQL)
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("Failed to create table in database %s", dbName))
		}

		// 解析并修改命名空间配置
		initNs, err := config.ParseNamespaceTmpl(nsTmpl, sliceConfig)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = e2eMgr.ModifyNamespace(initNs)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// 等待数据同步和命名空间加载
		time.Sleep(2 * time.Second)
	})

	// CRUD 操作测试
	ginkgo.Context("CRUD Operations with Prepared Statements for Multiple Types", func() {
		ginkgo.It("Decimal Precision Cases", func() {
			// 初始化 Gaea 连接
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			_, err = gaeaConn.Exec("USE db_mycat_mod")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// 插入数据
			insertStmt, err := gaeaConn.Prepare(`
				INSERT INTO tbl_mycat (id, name, age, salary, rating, is_active, created_at, description, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			createdAt := []uint8(time.Now().Format("2006-01-02 15:04:05"))
			updatedAt := []uint8(time.Now().Format("2006-01-02 15:04:05"))

			// 插入一条数据
			_, err = insertStmt.Exec(
				1,                             // id
				"John Doe",                    // name
				30,                            // age
				decimal.NewFromFloat(1234.56), // salary
				4.5,                           // rating
				true,                          // is_active
				createdAt,                     // created_at
				"Sample description",          // description
				updatedAt,                     // updated_at
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// 查询数据
			query := `
				SELECT id, name, age, salary, rating, 
					is_active, created_at, description, updated_at 
				FROM tbl_mycat
				WHERE id = ?
			`

			selectStmt, err := gaeaConn.Prepare(query)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			var sample tbl_mycat
			err = selectStmt.QueryRow(1).Scan(
				&sample.ID,
				&sample.Name,
				&sample.Age,
				&sample.Salary,
				&sample.Rating,
				&sample.IsActive,
				&sample.CreatedAt,
				&sample.Description,
				&sample.UpdatedAt,
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// 验证数据
			gomega.Expect(sample.ID).To(gomega.Equal(1))
			gomega.Expect(sample.Name).To(gomega.Equal("John Doe"))
			gomega.Expect(sample.Age).To(gomega.Equal(30))
			gomega.Expect(sample.Salary.String()).To(gomega.Equal(decimal.NewFromFloat(1234.56).String()))
			gomega.Expect(sample.Rating).To(gomega.Equal(4.5))
			gomega.Expect(sample.IsActive).To(gomega.BeTrue())
			gomega.Expect(sample.Description).To(gomega.Equal("Sample description"))
			gomega.Expect(sample.CreatedAt).To(gomega.Equal(createdAt))
			gomega.Expect(sample.UpdatedAt).To(gomega.Equal(updatedAt))

			// Update data
			updateStmt, err := gaeaConn.Prepare("UPDATE tbl_mycat SET age = ? WHERE id = ?")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			_, err = updateStmt.Exec(26, 1)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Delete data
			deleteStmt, err := gaeaConn.Prepare("DELETE FROM tbl_mycat WHERE id = ?")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			_, err = deleteStmt.Exec(1)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("Edge Cases", func() {
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			_, err = gaeaConn.Exec("USE db_mycat_mod")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			insertStmt, err := gaeaConn.Prepare("INSERT INTO tbl_mycat (id, salary) VALUES (?, ?)")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			testCases := []struct {
				id     int
				salary float64
			}{
				{id: 1, salary: 123.00},
				{id: 2, salary: 456.70},
				{id: 3, salary: 0.01},
			}

			for _, tc := range testCases {
				_, err = insertStmt.Exec(tc.id, tc.salary)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			}

			selectStmt, err := gaeaConn.Prepare("SELECT id, salary FROM tbl_mycat WHERE id = ?")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			for _, tc := range testCases {
				var (
					id     int
					salary float64
				)
				err = selectStmt.QueryRow(tc.id).Scan(&id, &salary)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(id).To(gomega.Equal(tc.id))
				gomega.Expect(salary).To(gomega.BeNumerically("~", tc.salary, 0.001))
			}
		})
	})

	// AfterEach: 清理测试环境
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
