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

import "errors"

// Slice means config model of slice
type Slice struct {
	Name            string   `json:"name"`
	UserName        string   `json:"user_name"`
	Password        string   `json:"password"`
	Master          string   `json:"master"`
	Slaves          []string `json:"slaves"`
	StatisticSlaves []string `json:"statistic_slaves"`

	Capacity    int `json:"capacity"`     // connection pool capacity
	MaxCapacity int `json:"max_capacity"` // max connection pool capacity
	IdleTimeout int `json:"idle_timeout"` // close backend direct connection after idle_timeout,unit: seconds
}

func (s *Slice) verify() error {
	if s.Name == "" {
		return errors.New("must specify slice name")
	}

	if s.UserName == "" {
		return errors.New("missing user")
	}

	if s.Master == "" && len(s.Slaves) == 0 {
		return errors.New("both master and slaves empty")
	}

	for _, slave := range s.Slaves {
		if slave == "" {
			return errors.New("illegal slave addr")
		}
	}

	if s.Capacity <= 0 {
		return errors.New("connection pool capacity should be > 0")
	}

	if s.MaxCapacity <= 0 {
		return errors.New("max connection pool capactiy should be > 0")
	}

	return nil
}
