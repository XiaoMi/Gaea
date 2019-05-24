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

// constants of shard type
const (
	ShardDefault         = "default"
	ShardGlobal          = "global"
	ShardLinked          = "linked"
	ShardMod             = "mod"
	ShardHash            = "hash"
	ShardRange           = "range"
	ShardYear            = "date_year"
	ShardMonth           = "date_month"
	ShardDay             = "date_day"
	ShardMycatMod        = "mycat_mod"
	ShardMycatLong       = "mycat_long"
	ShardMycatString     = "mycat_string"
	ShardMycatMURMUR     = "mycat_murmur"
	ShardMycatPaddingMod = "mycat_padding_mod"
)

// Shard means shard model in etcd
type Shard struct {
	DB            string   `json:"db"`
	Table         string   `json:"table"`
	ParentTable   string   `json:"parent_table"`
	Type          string   `json:"type"` // 表类型: 包括分表如hash/range/data,关联表如: linked 全局表如: global等
	Key           string   `json:"key"`
	Locations     []int    `json:"locations"`
	Slices        []string `json:"slices"`
	DateRange     []string `json:"date_range"`
	TableRowLimit int      `json:"table_row_limit"`

	// only used in mycat logic database (schema)
	Databases []string `json:"databases"`

	// used in mycat partition long shard and partition string shard
	PartitionCount  string `json:"partition_count"`
	PartitionLength string `json:"partition_length"`

	// used in mycat partition string shard
	HashSlice string `json:"hash_slice"`

	// used in mycat murmur hash shard
	Seed               string `json:"seed"`
	VirtualBucketTimes string `json:"virtual_bucket_times"`

	// used in mycat padding mod shard
	PadFrom   string `json:"pad_from"`
	PadLength string `json:"pad_length"`
	ModBegin  string `json:"mod_begin"`
	ModEnd    string `json:"mod_end"`
}

func (p *Shard) verify() error {
	return nil
}

// Encode encode json
func (p *Shard) Encode() []byte {
	return JSONEncode(p)
}
