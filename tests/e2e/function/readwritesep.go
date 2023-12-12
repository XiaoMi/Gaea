package function

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Read-Write Splitting", func() {
	nsTemplateFile := "e2e/function/ns/default.template"
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMSName]
	table := config.DefaultE2eTable
	ginkgo.BeforeEach(func() {
		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), slice)
		gomega.Expect(err).Should(gomega.BeNil())
		masterConn, err := slice.GetMasterConn(0)
		gomega.Expect(err).Should(gomega.BeNil())
		err = util.SetupDatabaseAndInsertData(masterConn, db, table)
		gomega.Expect(err).Should(gomega.BeNil())
	})

	ginkgo.Context("When handling read and write operations", func() {
		ginkgo.It("should direct read operations to replicas", func() {
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())

			gaeaWriteConn, err := e2eMgr.GetWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())

			gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())

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
					gomega.Expect(err).Should(gomega.BeNil())
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
				gomega.Expect(err).Should(gomega.BeNil())
				gomega.Expect(res).Should(gomega.HaveLen(1))
				gomega.Expect(sqlCase.ExpectBackendAddr).Should(gomega.Equal(res[0].BackendAddr))
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})

})
