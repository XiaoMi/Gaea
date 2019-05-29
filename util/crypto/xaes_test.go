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

package crypto

import (
	"encoding/base64"
	"testing"
)

func TestEncryptECB(t *testing.T) {

	key := "1234abcd5678efg*"
	msg := "mQSa0mS1Gi1Q8VCVLHrbU_izaHqzaPmh"
	data, err := EncryptECB(key, []byte(msg))
	if err != nil {
		t.Fatalf("encrypt failed, err:%v", err)
		return
	}
	base64Str := base64.StdEncoding.EncodeToString(data)

	t.Logf("encrypt succ, data: %v, data str: %s, len:%v", data, base64Str, len(data))

	data, _ = base64.StdEncoding.DecodeString(base64Str)
	origin, err := DecryptECB(key, data)
	if err != nil {
		t.Fatalf("decrypt failed, err:%v", err)
		return
	}

	if string(origin) != msg {
		t.Fatalf("origin not equal msg")
	}
	t.Log("decrypt succ")
}
