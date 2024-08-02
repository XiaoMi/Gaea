// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zap

import (
	"os"
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkSyncLoggerWriter(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	encoder := &ZapEncoder{}
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(f), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return true
		})),
	)
	l := zap.New(core)
	g := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		g.Add(1)
		go func() {
			defer g.Done()
			l.Info("ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)")
		}()
	}
	g.Wait()
	l.Sync()
}

func BenchmarkAsyncLoggerWriter(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	encoder := &ZapEncoder{}
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, NewAsyncWriter(f), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return true
		})),
	)
	l := zap.New(core)
	g := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		g.Add(1)
		go func() {
			defer g.Done()
			l.Info("ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)")
		}()
	}
	g.Wait()
	l.Sync()
}
