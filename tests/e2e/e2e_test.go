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

package e2e

import (
	"testing"
	"time"

	"github.com/XiaoMi/Gaea/tests/e2e/util"

	_ "github.com/XiaoMi/Gaea/tests/e2e/dml"
	_ "github.com/XiaoMi/Gaea/tests/e2e/function"
	_ "github.com/XiaoMi/Gaea/tests/e2e/shard"
	_ "github.com/XiaoMi/Gaea/tests/e2e/unshard"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode.")
	}
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "gaea E2E Testing")
}

var _ = ginkgo.BeforeSuite(func() {
	ginkgo.By("start gaea default.")
	err := util.StartGaeaDefault()
	util.ExpectNoError(err)
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea default success.")

	ginkgo.By("start gaea-cc default.")
	err = util.StartGaeaCCDefault()
	util.ExpectNoError(err)
	time.Sleep(5 * time.Second)
	ginkgo.By("start gaea-cc default success.")
})

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("stop gaea-cc default.")
	err := util.StopGaeaCCDefault()
	util.ExpectNoError(err)
	ginkgo.By("stop gaea-cc default success.")

	ginkgo.By("stop gaea default.")
	err = util.StopGaeaDefault()
	util.ExpectNoError(err)
	ginkgo.By("stop gaea default success.")

})
