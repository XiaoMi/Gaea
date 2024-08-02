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

package unshard

import (
	"path/filepath"

	config "github.com/XiaoMi/Gaea/tests/e2e/config"
	"github.com/XiaoMi/Gaea/tests/e2e/util"
	"github.com/onsi/ginkgo/v2"
)

var _ = ginkgo.Describe("sql_support_test", func() {
	planManagers := []*config.PlanManager{}
	e2eMgr := config.NewE2eManager()
	slice := e2eMgr.NsSlices[config.SliceSingleSlave]
	gaeaConn, err := e2eMgr.GetReadWriteGaeaUserConn()
	util.ExpectNoError(err)

	ginkgo.BeforeEach(func() {
		// 注册
		ns, err := config.ParseNamespaceTmpl(config.UnShardDMLNamespaceTmpl, slice)
		util.ExpectNoError(err)
		err = e2eMgr.ModifyNamespace(ns)
		util.ExpectNoError(err)

		casesPath, err := config.GetJSONFilesFromDir(filepath.Join(e2eMgr.BasePath, "unshard/case/dml"))
		util.ExpectNoError(err)

		conns, err := slice.GetLocalSliceConn()
		util.ExpectNoError(err)

		ginkgo.By("get sql plan")
		planManagers = []*config.PlanManager{}
		for _, v := range casesPath {
			p := &config.PlanManager{
				PlanPath:         v,
				MysqlClusterConn: conns,
				GaeaDB:           gaeaConn,
			}
			planManagers = append(planManagers, p)
		}
	})

	ginkgo.Context("sql support test", func() {
		ginkgo.It("sql support", func() {
			for _, p := range planManagers {
				err := p.Init()
				util.ExpectNoError(err)
				err = p.Run()
				util.ExpectNoError(err)
			}
		})
	})

	ginkgo.AfterEach(func() {
		for _, p := range planManagers {
			p.MysqlClusterConnClose()
		}
		e2eMgr.Clean()
	})
})
