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

// GlobalSequence means config of global sequences with different types
type GlobalSequence struct {
	DB        string `json:"db"`
	Table     string `json:"table"`
	Type      string `json:"type"`       // 全局序列号类型,目前只兼容mycat的数据库方式
	SliceName string `json:"slice_name"` // 对应sequence表所在的分片，默认都在0号片
	PKName    string `json:"pk_name"`    // 全局序列号字段名称
}

// Encode means encode for easy use
func (p *GlobalSequence) Encode() []byte {
	return JSONEncode(p)
}
