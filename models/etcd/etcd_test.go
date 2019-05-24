// Copyright 2019 The Gaea Authors. All Rights Reserved.
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

package etcdclient

import (
	"testing"

	"github.com/coreos/etcd/client"
)

func Test_isErrNoNode(t *testing.T) {
	err := client.Error{}
	err.Code = client.ErrorCodeKeyNotFound
	if !isErrNoNode(err) {
		t.Fatalf("test isErrNoNode failed, %v", err)
	}
	err.Code = client.ErrorCodeNotFile
	if isErrNoNode(err) {
		t.Fatalf("test isErrNoNode failed, %v", err)
	}
}

func Test_isErrNodeExists(t *testing.T) {
	err := client.Error{}
	err.Code = client.ErrorCodeNodeExist
	if !isErrNodeExists(err) {
		t.Fatalf("test isErrNodeExists failed, %v", err)
	}
	err.Code = client.ErrorCodeNotFile
	if isErrNodeExists(err) {
		t.Fatalf("test isErrNodeExists failed, %v", err)
	}
}
