package shard

import (
	"database/sql"
	"path/filepath"

	"github.com/XiaoMi/Gaea/tests/config"
	"github.com/XiaoMi/Gaea/tests/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("shard join support test in king shard hash ", func() {
	nsTemplateFile := "e2e/shard/ns/kingshard/hash.template"
	e2eMgr := config.NewE2eManager()
	sliceMulti := e2eMgr.NsSlices[config.SliceMMName]
	sliceSingle := e2eMgr.NsSlices[config.SliceSMName]
	multiMasterA, err := sliceMulti.GetMasterConn(0)
	gomega.Expect(err).Should(gomega.BeNil())

	multiMasterB, err := sliceMulti.GetMasterConn(1)
	gomega.Expect(err).Should(gomega.BeNil())

	singleMaster, err := sliceSingle.GetMasterConn(0)
	gomega.Expect(err).Should(gomega.BeNil())

	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	gomega.Expect(err).Should(gomega.BeNil())
	ginkgo.BeforeEach(func() {

		err := e2eMgr.AddNsFromFile(filepath.Join(e2eMgr.BasePath, nsTemplateFile), sliceMulti)
		gomega.Expect(err).Should(gomega.BeNil())
		_, err = gaeaConn.Exec("USE sbtest1")
		gomega.Expect(err).Should(gomega.BeNil())

		prepareCases := []struct {
			DB   *sql.DB
			file string
		}{
			{
				DB:   multiMasterA,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-slice0-kingshard.sql"),
			},
			{
				DB:   multiMasterB,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-slice1-kingshard.sql"),
			},
			{
				DB:   singleMaster,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-test-prepare.sql"),
			},
			{
				DB:   gaeaConn,
				file: filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/0-gaea-prepare.sql"),
			},
		}
		for _, v := range prepareCases {
			sqls, err := util.GetSqlFromFile(v.file)
			gomega.Expect(err).Should(gomega.BeNil())
			for _, sql := range sqls {
				err := util.VerifySqlParsable(sql)
				gomega.Expect(err).Should(gomega.BeNil())
				_, err = util.MysqlExec(v.DB, sql)
				gomega.Expect(err).Should(gomega.BeNil())
			}
		}

	})

	ginkgo.Context("shard support test", func() {
		ginkgo.It("When testing sql support", func() {
			cases := []struct {
				path       string
				resultType util.ResultType
			}{
				{
					filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/kingshard/equal.sql"),
					util.Equal,
				},
				{
					filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/kingshard/unequal.sql"),
					util.UnEqual,
				},
				{
					filepath.Join(e2eMgr.BasePath, "e2e/shard/case/join/kingshard/unsupport.sql"),
					util.UnSupport,
				},
			}
			for _, c := range cases {
				sqls, err := util.GetSqlFromFile(c.path)
				gomega.Expect(err).Should(gomega.BeNil())
				for _, sql := range sqls {
					err := util.VerifySqlParsable(sql)
					gomega.Expect(err).Should(gomega.BeNil())
					switch c.resultType {
					case util.Equal:
						res1, err := util.MysqlQuery(gaeaConn, sql)
						gomega.Expect(err).Should(gomega.BeNil())
						res2, err := util.MysqlQuery(singleMaster, sql)
						gomega.Expect(err).Should(gomega.BeNil())
						_, err = util.CompareIgnoreSort(res1, res2)
						gomega.Expect(err).Should(gomega.BeNil())
					case util.UnSupport:
						_, err := util.MysqlQuery(singleMaster, sql)
						gomega.Expect(err).Should(gomega.BeNil())
						_, err = util.MysqlQuery(gaeaConn, sql)
						gomega.Expect(err).ShouldNot(gomega.BeNil())
					case util.UnEqual:
						res1, err := util.MysqlQuery(gaeaConn, sql)
						gomega.Expect(err).Should(gomega.BeNil())
						res2, err := util.MysqlQuery(singleMaster, sql)
						gomega.Expect(err).Should(gomega.BeNil())
						_, err = util.CompareIgnoreSort(res1, res2)
						gomega.Expect(err).ShouldNot(gomega.BeNil())
					}
				}
			}
		})
	})

	ginkgo.AfterEach(func() {
		e2eMgr.Clean()
	})
})
