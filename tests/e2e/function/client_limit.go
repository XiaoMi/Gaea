package function

import (
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

// This test script is designed to test the client connection limit functionality.
// It sets up a testing environment with a specified maximum number of client connections (10 in this case).
// The BeforeEach block configures a namespace and prepares a database for testing.
// The actual test, within the It block, attempts to create more connections (20) than the maximum allowed, and then checks if connections exceeding the limit are correctly rejected.
// The script counts the number of errors (expected to be equal to the number of connections exceeding the limit) and successful connections (expected to be equal to the limit).
// The AfterEach block cleans up the environment after the test.
var _ = ginkgo.Describe("test client connection limit", func() {
	e2eMgr := config.NewE2eManager()
	db := config.DefaultE2eDatabase
	slice := e2eMgr.NsSlices[config.SliceSingleMaster]
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
