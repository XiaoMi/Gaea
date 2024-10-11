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

package server

import (
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleUseDB(t *testing.T) {
	var se SessionExecutor
	mockey.PatchConvey("TestHandleUseDB", t, func() {
		mockey.Mock((*SessionExecutor).GetNamespace).Return(&Namespace{allowedDBs: map[string]bool{"test": true}}).Build()
		err := se.handleUseDB("test")
		assert.Nil(t, err)

		mockey.UnPatchAll()
		mockey.Mock((*SessionExecutor).GetNamespace).Return(&Namespace{allowedDBs: map[string]bool{"information_schema": true}}).Build()
		err = se.handleUseDB("information_schema")
		assert.Nil(t, err)

		mockey.UnPatchAll()
		mockey.Mock((*SessionExecutor).GetNamespace).Return(&Namespace{allowedDBs: map[string]bool{"information_schema": true}}).Build()
		err = se.handleUseDB("information_SCHEMA")
		assert.Nil(t, err)

		mockey.UnPatchAll()
		mockey.Mock((*SessionExecutor).GetNamespace).Return(&Namespace{allowedDBs: map[string]bool{"test": true}}).Build()
		err = se.handleUseDB("test1")
		assert.NotNil(t, err)

		mockey.UnPatchAll()
		mockey.Mock((*SessionExecutor).GetNamespace).Return(&Namespace{}).Build()
		err = se.handleUseDB("")
		assert.NotNil(t, err)
	})
}
