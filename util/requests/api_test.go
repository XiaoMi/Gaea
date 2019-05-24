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

package requests

import (
	"encoding/json"
	"net/http/httputil"
	"testing"
)

func TestAddParameters(t *testing.T) {
	url := "http://127.0.0.1:8080"
	expect := "http://127.0.0.1:8080?a=test1&b=test2"
	params := make(map[string]string, 2)
	params["a"] = "test1"
	params["b"] = "test2"
	r := AddParameters(url, params)
	if r != expect {
		t.Errorf("test AddParameters failed, expect: %s, result: %s", expect, r)
	}
}

func TestBuildHTTPRequest(t *testing.T) {
	url := "http://127.0.0.1:8080"
	method := Put

	header := make(map[string]string, 2)
	header["Content-Type"] = "application/json"

	params := make(map[string]string, 2)
	params["a"] = "test1"
	params["b"] = "test2"

	type ReqBody struct {
		Field1 string `json:"field_1"`
		Field2 int    `json:"field_2"`
	}
	rb := ReqBody{Field1: "contest1", Field2: 99}
	body, _ := json.Marshal(rb)

	req := NewRequest(url, method, header, params, body)
	req.SetBasicAuth("test", "test")
	httpReq, err := BuildHTTPRequest(req)
	if err != nil {
		t.Errorf("test BuildHTTPRequest failed, %v", err)
	}

	httpReqDump, err := httputil.DumpRequest(httpReq, true)
	if err != nil {
		t.Errorf("test BuildHTTPRequest failed, %v", err)
	}
	t.Logf("dumped http request: %s\n", string(httpReqDump))
}

func TestEncodeURL(t *testing.T) {
	url := EncodeURL("127.0.0.1:8080", "/test/test/:%d", 1)
	if url != "http://127.0.0.1:8080/test/test/:1" {
		t.Fatalf("test EncodeURL failed, %s", url)
	}
}
