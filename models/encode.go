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

package models

import (
	"encoding/json"

	"github.com/XiaoMi/Gaea/log"
)

// JSONEncode return json encoding of v
func JSONEncode(v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		//TODO panic
		log.Fatal("encode to json failed, %v", err)
		return nil
	}
	return b
}

// JSONDecode parses the JSON-encoded data and stores the result in the value pointed to by v
func JSONDecode(v interface{}, data []byte) error {
	return json.Unmarshal(data, v)
}
