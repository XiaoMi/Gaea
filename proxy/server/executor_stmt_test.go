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

package server

import (
	"testing"
)

func Test_calcParams(t *testing.T) {
	sql := "update micf_order_0 set order_status=4, update_time=1541831505" +
		"where\n" +
		"order_id in ('1321989216361392') and\n" +
		"project_id = 371 and\n" +
		"order_status = 2"
	paramCount, offsets, err := calcParams(sql)
	t.Log(paramCount)
	t.Log(offsets)
	t.Log(err)
	if err != nil {
		t.Logf("test calcParams failed, %v\n", err)
	}

	sql = "select * from t1 where id = ? and col = ?"
	paramCount, offsets, err = calcParams(sql)
	t.Log(paramCount)
	t.Log(offsets)
	t.Log(err)
	if err != nil || paramCount != 2 {
		t.Logf("test calcParams failed, %v\n", err)
	}
}
