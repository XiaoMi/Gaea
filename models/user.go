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
	"errors"
	"fmt"
	"strings"
)

// 用户只读标识
// read write privileges
const (
	// 只读
	ReadOnly = 1
	// 可读可写
	ReadWrite = 2
)

// 用户读写分离标识
const (
	// NoReadWriteSplit 非读写分离
	NoReadWriteSplit = 0
	// ReadWriteSplit 读写分离
	ReadWriteSplit = 1
	// StatisticUser 统计用户
	StatisticUser = 1
)

// User meand user struct
type User struct {
	UserName      string `json:"user_name"`
	Password      string `json:"password"`
	Namespace     string `json:"namespace"`
	RWFlag        int    `json:"rw_flag"`        //1: 只读 2:读写
	RWSplit       int    `json:"rw_split"`       //0: 不采用读写分离 1:读写分离
	OtherProperty int    `json:"other_property"` // 1:统计用户
}

func (p *User) verify() error {
	if p.UserName == "" {
		return errors.New("missing user name")
	}
	p.UserName = strings.TrimSpace(p.UserName)

	if p.Namespace == "" {
		return fmt.Errorf("missing namespace name: %s", p.UserName)
	}
	p.Namespace = strings.TrimSpace(p.Namespace)

	if p.Password == "" {
		return fmt.Errorf("missing password: [%s]%s", p.Namespace, p.UserName)
	}
	p.Password = strings.TrimSpace(p.Password)

	if p.RWFlag != ReadOnly && p.RWFlag != ReadWrite {
		return fmt.Errorf("invalid RWFlag, user: %s, rwflag: %d", p.UserName, p.RWFlag)
	}

	if p.RWSplit != NoReadWriteSplit && p.RWSplit != ReadWriteSplit {
		return fmt.Errorf("invalid RWSplit, user: %s, rwsplit: %d", p.UserName, p.RWSplit)
	}

	if p.OtherProperty != StatisticUser && p.OtherProperty != 0 {
		return fmt.Errorf("invalid other property, user: %s, %d", p.UserName, p.OtherProperty)
	}

	return nil
}
