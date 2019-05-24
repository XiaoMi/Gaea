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
	"testing"
	"time"
)

type A struct {
	a int
	b string
}

func callback() {
	fmt.Println("timeout")
}

func newTimeWheel() *TimeWheel {
	tw, err := NewTimeWheel(time.Second, 3600)
	if err != nil {
		panic(err)
	}
	tw.Start()
	return tw
}

func TestAdd(t *testing.T) {
	tw := newTimeWheel()
	err := tw.Add(time.Second*1, "test", callback)
	if err != nil {
		t.Fatalf("test add failed, %v", err)
	}
	time.Sleep(time.Second * 5)
	tw.Stop()
}

func TestRemove(t *testing.T) {
	a := &A{a: 10, b: "test"}
	tw := newTimeWheel()
	err := tw.Add(time.Second*1, a, callback)
	if err != nil {
		t.Fatalf("test add failed, %v", err)
	}
	tw.Remove(a)
	time.Sleep(time.Second * 5)
	tw.Stop()
}

func BenchmarkAdd(b *testing.B) {
	tw := newTimeWheel()
	for i := 0; i < b.N; i++ {
		key := "test" + strconv.Itoa(i)
		err := tw.Add(time.Second, key, callback)
		if err != nil {
			b.Fatalf("benchmark Add failed, %v", err)
		}
	}
}
