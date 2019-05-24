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

package util

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

// BoolIndex rolled array switch mark
type BoolIndex struct {
	index int32
}

// Set set index value
func (b *BoolIndex) Set(index bool) {
	if index {
		atomic.StoreInt32(&b.index, 1)
	} else {
		atomic.StoreInt32(&b.index, 0)
	}
}

// Get return current, next, current bool value
func (b *BoolIndex) Get() (int32, int32, bool) {
	index := atomic.LoadInt32(&b.index)
	if index == 1 {
		return 1, 0, true
	}
	return 0, 1, false
}

// ItoString interface to string
func ItoString(a interface{}) (bool, string) {
	switch a.(type) {
	case nil:
		return false, "NULL"
	case []byte:
		return true, string(a.([]byte))
	default:
		return false, fmt.Sprintf("%v", a)
	}
}

// Int2TimeDuration convert int to Time.Duration
func Int2TimeDuration(t int) (time.Duration, error) {
	tmp := strconv.Itoa(t)
	tmp = tmp + "s"
	idleTimeout, err := time.ParseDuration(tmp)
	if err != nil {
		return 0, err

	}
	return idleTimeout, nil
}
