package function

import (
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"time"
)

var _ = ginkgo.Describe("test client connection limit", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceMaster]
	table := config.DefaultE2eTable
	maxConnections := 10
	ginkgo.BeforeEach(func() {
		initNs, err := config.ParseNamespaceTmpl(config.DefaultNamespaceTmpl, slice)
		util.ExpectNoError(err, "parse namespace template")
		initNs.MaxClientConnections = maxConnections
		e2eMgr.NsManager.DeleteNamespace(initNs)
		initNs.Name = "test_client_limit"
		for i := 0; i < len(initNs.Users); i++ {
			initNs.Users[i].Namespace = initNs.Name
		}

		err = e2eMgr.NsManager.ModifyNamespace(initNs)
		util.ExpectNoError(err, "create namespace")
		masterAdminConn, err := slice.GetMasterAdminConn(0)
		util.ExpectNoError(err, "get master admin conn")
		err = util.SetupDatabaseAndInsertData(masterAdminConn, db, table)
		util.ExpectNoError(err, "setup database and insert data")
	})
	ginkgo.Context("When handling client limit operations", func() {
		ginkgo.It("should limit exceeded maximum connection ", func() {
			testCounts := 20
			gaeaReadWriteConn, err := e2eMgr.GetReadWriteGaeaUserConn()
			util.ExpectNoError(err, "get gaea read write conn")
			errChan := make(chan error, testCounts)
			for i := 0; i < testCounts; i++ {
				index := i
				go func() {
					time.Sleep(time.Duration(index) * 10 * time.Millisecond)
					_, err := gaeaReadWriteConn.Query("select sleep(0.1)")
					errChan <- err
				}()
			}
			// check err count
			errCount, noErrCount := 0, 0
			for i := 0; i < testCounts; i++ {
				err := <-errChan
				if err != nil {
					errCount++
				} else {
					noErrCount++
				}
			}
			util.ExpectEqual(errCount, testCounts-maxConnections, "err count")
			util.ExpectEqual(noErrCount, maxConnections, "no err count")
		})
	})
	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
