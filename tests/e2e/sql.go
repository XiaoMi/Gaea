package e2e

import (
	"database/sql"
	"fmt"
	"github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("sql", func() {

	var cfg *config.Config
	var gDB *sql.DB
	var mDB *sql.DB
	var sDB *sql.DB
	var err error

	ginkgo.BeforeEach(func() {
		cfg = config.TestConfig
		gDB, err = util.Open(cfg.GaeaHost, cfg.GaeaPort, cfg.GaeaUser, cfg.GaeaPassword, 10, 3)
		gomega.Expect(err).Should(gomega.BeNil())
		mDB, err = util.Open(cfg.MasterHost, cfg.MasterPort, cfg.MasterUser, cfg.MasterPassword, 10, 3)
		gomega.Expect(err).Should(gomega.BeNil())
		sDB, err = util.Open(cfg.SlaveHost, cfg.SlavePort, cfg.SlaveUser, cfg.SlavePassword, 10, 3)
		gomega.Expect(err).Should(gomega.BeNil())
		// init DB
		ginkgo.By("init mysql data.")

		for _, s := range cfg.GenDropTableSqls() {
			_, err = mDB.Exec(s)
			gomega.Expect(err).Should(gomega.BeNil())
		}

		for _, s := range cfg.GenCreateTableSqls() {
			_, err = mDB.Exec(s)
			gomega.Expect(err).Should(gomega.BeNil())
		}

		_, err = mDB.Exec(fmt.Sprintf(`insert into %s.%s values(0,'a'),(1,'b'),(2,'c')`, cfg.Database, cfg.Tables[0]))
		gomega.Expect(err).Should(gomega.BeNil())

		sTx, err := sDB.Begin()
		gomega.Expect(err).Should(gomega.BeNil())
		_, err = sTx.Exec(fmt.Sprintf(`insert into %s.%s values(3,'d')`, cfg.Database, cfg.Tables[0]))
		gomega.Expect(err).Should(gomega.BeNil())
		err = sTx.Commit()
		gomega.Expect(err).Should(gomega.BeNil())
	})

	ginkgo.AfterEach(func() {
		// fix slave
		ginkgo.By("restore mysql data.")
		sTx, err := sDB.Begin()
		gomega.Expect(err).Should(gomega.BeNil())

		_, err = sTx.Exec(fmt.Sprintf(`delete from %s.%s where id=3;`, cfg.Database, cfg.Tables[0]))

		err = sTx.Commit()
		gomega.Expect(err).Should(gomega.BeNil())

		// delete cluster data
		_, err = mDB.Exec(fmt.Sprintf(`delete from %s.%s`, cfg.Database, cfg.Tables[0]))
		gomega.Expect(err).Should(gomega.BeNil())

	})

	ginkgo.Context("hint", func() {
		ginkgo.It("master hint", func() {
			ssql := fmt.Sprintf(`select /*master*/ count(*) from %s.%s`, cfg.Database, cfg.Tables[0])
			checkFunc := func(rows1 *sql.Rows, rows2 *sql.Rows) bool {
				defer rows1.Close()
				defer rows2.Close()
				for rows1.Next() {
					var count1 int
					var count2 int
					err1 := rows1.Scan(&count1)
					if ok := rows2.Next(); !ok {
						return false
					}
					err2 := rows2.Scan(&count2)
					if err1 != err2 {
						return false
					}
					if count1 != count2 {
						return false
					}
				}
				if rows2.Next() {
					return false
				}
				return true
			}
			same := util.CheckDBRes(gDB, mDB, ssql, checkFunc)
			gomega.Expect(same).Should(gomega.BeTrue())
			same = util.CheckDBRes(sDB, gDB, ssql, checkFunc)
			gomega.Expect(same).Should(gomega.BeFalse())
		})
	})
})
