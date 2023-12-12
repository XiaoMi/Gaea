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

var _ = ginkgo.Describe("Load Balancing", func() {
	nsTemplateFile := "e2e/function/ns/default.template"
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMSName]
	table := config.DefaultE2eTable
	currentTime := time.Now()
	ginkgo.BeforeEach(func() {

		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), slice)
		gomega.Expect(err).Should(gomega.BeNil())
		masterConn, err := slice.GetMasterConn(0)
		gomega.Expect(err).Should(gomega.BeNil())
		err = util.SetupDatabaseAndInsertData(masterConn, db, table)
		gomega.Expect(err).Should(gomega.BeNil())
	})

	ginkgo.Context("When distributing queries among replicas", func() {
		ginkgo.It("should evenly distribute read queries", func() {

			gaeaReadConn, err := e2eMgr.GetReadGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())

			gaeaWriteConn, err := e2eMgr.GetWriteGaeaUserConn()
			gomega.Expect(err).Should(gomega.BeNil())

			// 定义 SQL 测试用例
			sqlCases := []struct {
				GaeaConn          *sql.DB
				GaeaSQL           string
				ExpectBackendAddr string
			}{
				//读用户-读
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case1 */SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case2 */ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[1].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case3 */ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[0].Addr(),
				},
				{
					GaeaConn:          gaeaReadConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case4 */ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Slaves[1].Addr(),
				},
				// 写用户-读
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case5 */SELECT * FROM `%s`.`%s` WHERE `id`= 1", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case6 */ SELECT * FROM `%s`.`%s` WHERE `id`= 2", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case7 */ SELECT * FROM `%s`.`%s` WHERE `id`= 3", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case8 */ SELECT * FROM `%s`.`%s` WHERE `id`= 4", db, table),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				// 写用户-写
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case9 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case10 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case11 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
				{
					GaeaConn:          gaeaWriteConn,
					GaeaSQL:           fmt.Sprintf("/*Load Balancing Case12 */ INSERT INTO `%s`.`%s` (name) VALUES ('%s')", db, table, "nameValue"),
					ExpectBackendAddr: e2eMgr.MClusterMasterSlaves.Master.Addr(),
				},
			}
			// 执行 SQL 测试用例
			for _, sqlCase := range sqlCases {
				_, err := sqlCase.GaeaConn.Exec(sqlCase.GaeaSQL)
				gomega.Expect(err).Should(gomega.BeNil())
				res, err := e2eMgr.SearchLog(sqlCase.GaeaSQL, currentTime)
				gomega.Expect(err).Should(gomega.BeNil())
				// 避免扫到以前的数据
				gomega.Expect(res).Should(gomega.HaveLen(1))
				gomega.Expect(sqlCase.ExpectBackendAddr).Should(gomega.Equal(res[0].BackendAddr))
			}
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})

})
