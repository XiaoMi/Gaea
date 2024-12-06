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
		time.Sleep(50 * time.Millisecond)
	})

	ginkgo.Context("test client qps limit", func() {
		ginkgo.It("should limit client qps correctly", func() {
			sqlCases := []struct {
				limitQPS  uint32
				actualQPS int
				limitTx   bool
			}{
				{
					limitQPS:  5,
					actualQPS: 5,
					limitTx:   false,
				},
				{
					limitQPS:  0,
					actualQPS: 20,
					limitTx:   false,
				},
				{
					limitQPS:  5,
					actualQPS: 5,
					limitTx:   true,
				},
				{
					limitQPS:  0,
					actualQPS: 60,
					limitTx:   true,
				},
			}

			testCount := 20
			for index, test := range sqlCases {
				initNs.ClientQPSLimit = test.limitQPS
				initNs.SupportLimitTransaction = test.limitTx
				err = e2eMgr.ModifyNamespace(initNs)
				util.ExpectNoError(err)
				time.Sleep(100 * time.Millisecond)
				qpsCount := 0
				now := time.Now()
				if test.limitTx {
					for i := 0; i < testCount; i++ {
						gaeaTestConn, _ := e2eMgr.GetReadWriteGaeaUserConn()
						tx, err := gaeaTestConn.Begin()
						if err == nil {
							qpsCount++
						} else {
							continue
						}
						_, err = tx.Exec("select 1")
						if err == nil {
							qpsCount++
						} else {
							continue
						}
						err = tx.Commit()
						if err == nil {
							qpsCount++
						}
					}
				} else {
					for i := 0; i < testCount; i++ {
						gaeaTestConn, err := e2eMgr.GetReadWriteGaeaUserConn()
						util.ExpectNoError(err)
						_, err = gaeaTestConn.Exec("select 1")
						gaeaTestConn.Close()
						if err == nil {
							qpsCount++
						} else {
							fmt.Printf("********qps limit: %v", err)
						}
					}
				}

				duration := time.Since(now).Milliseconds()
				if test.limitQPS > 0 {
					util.ExpectEqual(qpsCount < 2*test.actualQPS, true, fmt.Sprintf("test %d: actual qps not expected, duration: %d", index, duration))
				} else {
					util.ExpectEqual(qpsCount == test.actualQPS, true, fmt.Sprintf("test %d: actual qps not expected, duration: %d", index, duration))
				}
			}
		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
