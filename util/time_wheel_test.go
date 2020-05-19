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
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type A struct {
	a            int
	b            string
	isCallbacked int32
}

func (a *A) callback() {
	atomic.StoreInt32(&a.isCallbacked, 1)
}

func (a *A) getCallbackValue() int32 {
	return atomic.LoadInt32(&a.isCallbacked)
}

func newTimeWheel() *TimeWheel {
	tw, err := NewTimeWheel(time.Second, 3600)
	if err != nil {
		panic(err)
	}
	tw.Start()
	return tw
}

func TestNewTimeWheel(t *testing.T) {
	tests := []struct {
		name      string
		tick      time.Duration
		bucketNum int
		hasErr    bool
	}{
		{tick: time.Second, bucketNum: 0, hasErr: true},
		{tick: time.Millisecond, bucketNum: 1, hasErr: true},
		{tick: time.Second, bucketNum: 1, hasErr: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewTimeWheel(test.tick, test.bucketNum)
			assert.Equal(t, test.hasErr, err != nil)
		})
	}
}

func TestAdd(t *testing.T) {
	tw := newTimeWheel()
	a := &A{}
	err := tw.Add(time.Second*1, "test", a.callback)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, int32(0), a.getCallbackValue())
	time.Sleep(time.Second * 2)
	assert.Equal(t, int32(1), a.getCallbackValue())
	tw.Stop()
}

func TestAddMultipleTimes(t *testing.T) {
	a := &A{}
	tw := newTimeWheel()
	for i := 0; i < 4; i++ {
		err := tw.Add(time.Second, "test", a.callback)
		assert.NoError(t, err)
		time.Sleep(time.Millisecond * 500)
		t.Logf("current: %d", i)
		assert.Equal(t, int32(0), a.getCallbackValue())
	}

	time.Sleep(time.Second * 2)
	assert.Equal(t, int32(1), a.getCallbackValue())
	tw.Stop()
}

func TestRemove(t *testing.T) {
	a := &A{a: 10, b: "test"}
	tw := newTimeWheel()
	err := tw.Add(time.Second*1, a, a.callback)
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, int32(0), a.getCallbackValue())
	err = tw.Remove(a)
	assert.NoError(t, err)
	time.Sleep(time.Second * 2)
	assert.Equal(t, int32(0), a.getCallbackValue())
	tw.Stop()
}

func BenchmarkAdd(b *testing.B) {
	a := &A{}
	tw := newTimeWheel()
	for i := 0; i < b.N; i++ {
		key := "test" + strconv.Itoa(i)
		err := tw.Add(time.Second, key, a.callback)
		if err != nil {
			b.Fatalf("benchmark Add failed, %v", err)
		}
	}
}
