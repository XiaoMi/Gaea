package function

import (
	"fmt"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"

	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("client qps limit test", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceDualSlave]
	table := config.DefaultE2eTable
	initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
	util.ExpectNoError(err, "parse namespace template")

	ginkgo.BeforeEach(func() {
		// mysql prepare
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err)
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(initNs)
		util.ExpectNoError(err)
		// wait mysql data  sync and namespace load
		time.Sleep(500 * time.Millisecond)
	})

	ginkgo.Context("test client qps limit", func() {
		ginkgo.It("should limit client qps correctly", func() {
			sqlCases := []struct {
				limitQPS  uint32
				actualQPS int
			}{
				{
					limitQPS:  10,
					actualQPS: 10,
				},
				{
					limitQPS:  0,
					actualQPS: 100,
				},
			}

			testCount := 100
			for index, test := range sqlCases {
				initNs.ClientQPSLimit = test.limitQPS
				err = e2eMgr.ModifyNamespace(initNs)
				util.ExpectNoError(err)
				time.Sleep(10 * time.Millisecond)
				gaeaTestConn, _ := e2eMgr.GetReadWriteGaeaUserConn()
				qpsCount := 0
				now := time.Now()
				for i := 0; i < testCount; i++ {
					_, err := gaeaTestConn.Exec("select 1")
					if err == nil {
						qpsCount++
					}
				}
				duration := time.Since(now).Milliseconds()
				util.ExpectEqual(qpsCount, test.actualQPS, fmt.Sprintf("test %d: actual qps not expected, duration: %d", index, duration))
			}
		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
