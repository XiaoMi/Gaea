package function

import (
	"database/sql"
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Read-Write Splitting", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMasterSlaves]
	table := config.DefaultE2eTable
	ginkgo.BeforeEach(func() {
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		err = e2eMgr.NsManager.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)
	})

	ginkgo.Context("When handling read and write operations", func() {
		ginkgo.It("should direct read operations to replicas", func() {
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err)
			defer gaeaReadWriteConn.Close()

			gaeaWriteConn, err := e2eMgr.GetWriteGaeaUserConn()
			util.ExpectNoError(err)
			defer gaeaWriteConn.Close()

			gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
			util.ExpectNoError(err)
			defer gaeaReadConn.Close()

			sqlCases := []struct {
				GaeaConn          *sql.DB
				GaeaSQL           string
				ExpectBackendAddr string
				IsSuccess         bool
			}{
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1 FOR UPDATE", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("DELETE FROM %s.%s WHERE id=2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("UPDATE %s.%s SET name= '%s' WHERE id=3", db, table, "newName"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaReadWriteConn,
					GaeaSQL:           fmt.Sprintf("INSERT INTO %s.%s (name) VALUES ('tempValue')", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				// 写用户
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1 FOR UPDATE", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("DELETE FROM %s.%s WHERE id=2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("UPDATE %s.%s SET name= '%s' WHERE id=3", db, table, "newName"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("INSERT INTO %s.%s (name) VALUES ('tempValue')", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
					IsSuccess:         true,
				},
				//读用户
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1 FOR UPDATE", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[1].Addr(),
					IsSuccess:         true,
				},

				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
					IsSuccess:         true,
				},
				{
					GaeaConn:  gaeaReadConn,
					GaeaSQL:   fmt.Sprintf("DELETE FROM %s.%s  WHERE `id`= 2", db, table),
					IsSuccess: false,
				},
				{
					GaeaConn:  gaeaReadConn,
					GaeaSQL:   fmt.Sprintf("UPDATE %s.%s SET name= '%s' WHERE id=3", db, table, "newName"),
					IsSuccess: false,
				},
				{
					GaeaConn:  gaeaReadConn,
					GaeaSQL:   fmt.Sprintf("INSERT INTO %s.%s (name) VALUES ('tempValue')", db, table),
					IsSuccess: false,
				},
			}
			for _, sqlCase := range sqlCases {
				currentTime := time.Now()
				_, err := sqlCase.GaeaConn.Exec(sqlCase.GaeaSQL)
				if !sqlCase.IsSuccess {
					gomega.Expect(err).ShouldNot(gomega.BeNil())
					continue
				} else {
					util.ExpectNoError(err)
				}

				var res []util.LogEntry
				retryCount := 3 // 设置重试次数
				for i := 0; i < retryCount; i++ {
					time.Sleep(1 * time.Millisecond) // 等待一段时间再重试
					res, err = e2eMgr.SearchLog(sqlCase.GaeaSQL, currentTime)
					if err == nil && len(res) == 1 {
						break
					}
				}
				// 检查结果
				util.ExpectNoError(err)
				gomega.Expect(res).Should(gomega.HaveLen(1))
				gomega.Expect(sqlCase.ExpectBackendAddr).Should(gomega.Equal(res[0].BackendAddr))
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})

})
