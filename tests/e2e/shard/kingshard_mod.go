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

package shard

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("shard join support test in king shard mod", func() {
	e2eMgr := config.NewE2eManager()
	sliceTest := e2eMgr.NsSlices[config.SliceSingleTestMaster]
	sliceMulti := e2eMgr.NsSlices[config.SliceDualMaster]
	ginkgo.BeforeEach(func() {
		// 1. MySQL 创建数据库并清理
		multiMasterA, err := config.FetchAndCleanSliceMasterConn(sliceMulti, 0)
		util.ExpectNoError(err)

		multiMasterB, err := config.FetchAndCleanSliceMasterConn(sliceMulti, 1)
		util.ExpectNoError(err)

		singleMaster, err := config.FetchAndCleanSliceMasterConn(sliceTest, 0)
		util.ExpectNoError(err)

		// 2. 执行 MySQL 准备数据 SQL
		mysqlPrepareCases := []struct {
			DB   *sql.DB
			file string
		}{
			{
				DB:   multiMasterA,
				file: filepath.Join(e2eMgr.BasePath, "shard/case/join/0-slice0-kingshard.sql"),
			},
			{
				DB:   multiMasterB,
				file: filepath.Join(e2eMgr.BasePath, "shard/case/join/0-slice1-kingshard.sql"),
			},
			{
				DB:   singleMaster,
				file: filepath.Join(e2eMgr.BasePath, "shard/case/join/0-test-prepare.sql"),
			},
		}
		// 2. 执行 MySQL 准备数据 SQL
		for _, v := range mysqlPrepareCases {
			sqls, err := util.GetSqlFromFile(v.file)
			util.ExpectNoError(err)
			for _, sql := range sqls {
				_, err = util.MysqlExec(v.DB, sql)
				util.ExpectNoError(err)
			}
		}

		// 3. 注册命名空间
		ns, err := config.ParseNamespaceTmpl(config.KingShardModNamespaceTmpl, sliceMulti)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(ns)
		util.ExpectNoError(err)

		// 4. 获取 Gaea 连接
		gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
		util.ExpectNoError(err)

		gaeaConn.Exec("USE sbtest")
		util.ExpectNoError(err)

		// 5. 执行 Gaea 准备数据 SQL
		sqls, err := util.GetSqlFromFile(filepath.Join(e2eMgr.BasePath, "shard/case/join/0-gaea-prepare.sql"))
		util.ExpectNoError(err)
		for _, sql := range sqls {
			_, err = util.MysqlExec(gaeaConn, sql)
			util.ExpectNoError(err)
		}

	})

	ginkgo.Context("shard support test", func() {
		ginkgo.It("When testing sql support", func() {
			singleMaster, err := sliceTest.GetMasterCommonDBConn(0, "sbtest")
			util.ExpectNoError(err)
			defaultShardConn, err := sliceMulti.GetMasterCommonDBConn(0, "sbtest")
			util.ExpectNoError(err)
			gaeaConn, err := e2eMgr.GetReadWriteGaeaUserDBConn("sbtest")
			util.ExpectNoError(err)

			cases := []struct {
				path           string
				resultType     util.ResultType
				gaeaConn       *sql.DB
				comparisonConn *sql.DB
			}{
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/kingshard/equal.sql"),
					util.Equal,
					gaeaConn,
					singleMaster,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/kingshard/show.sql"),
					util.Equal,
					gaeaConn,
					defaultShardConn,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/kingshard/unequal.sql"),
					util.UnEqual,
					gaeaConn,
					singleMaster,
				},
				{
					filepath.Join(e2eMgr.BasePath, "shard/case/join/kingshard/unsupport.sql"),
					util.UnSupport,
					gaeaConn,
					singleMaster,
				},
			}
			for _, c := range cases {
				sqls, err := util.GetSqlFromFile(c.path)
				util.ExpectNoError(err)
				for _, sql := range sqls {
					switch c.resultType {
					case util.Equal:
						gaeaRes, err := util.MysqlQuery(c.gaeaConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("gaea exec equal sql err:%v, err:%v\n", sql, err))
						mysqlRes, err := util.MysqlQuery(c.comparisonConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec equal sql err:%v, err:%v\n", sql, err))
						_, err = util.CompareIgnoreSort(gaeaRes, mysqlRes)
						util.ExpectNoError(err, fmt.Sprintf("compare equal sql get uneuqal res.sql:%v, err:%v\n", sql, err))
					case util.UnSupport:
						_, err := util.MysqlQuery(c.comparisonConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec unsupport sql err:%v, err:%v\n", sql, err))
						gaeaRes, err := util.MysqlQuery(c.gaeaConn, sql)
						util.ExpectError(err, fmt.Sprintf("gaea exec unsupport sql get no error, sql:%v,res:%v\n", sql, gaeaRes))
					case util.UnEqual:
						gaeaRes, err := util.MysqlQuery(c.gaeaConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("gaea exec unequal sql err:%v, err:%v\n", sql, err))
						mysqlRes, err := util.MysqlQuery(c.comparisonConn, sql)
						util.ExpectNoError(err, fmt.Sprintf("mysql exec unequal sql err.sql:%v, err:%v\n", sql, err))
						_, err = util.CompareIgnoreSort(gaeaRes, mysqlRes)
						util.ExpectError(err, fmt.Sprintf("compare unequal sql get equal res.sql:%v, err:%v\n", sql, err), gaeaRes, mysqlRes)
					}
				}
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
