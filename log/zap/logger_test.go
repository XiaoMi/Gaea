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
	"github.com/stretchr/testify/assert"
	"io"
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

func TestCreateLogManager(t *testing.T) {
	config := map[string]string{
		"path":            "/tmp",
		"filename":        "test",
		"level":           "debug",
		"log_keep_days":   "7",
		"log_keep_counts": "30",
	}

	loggerManager, err := CreateLogManager(config)
	assert.NoError(t, err)
	assert.NotNil(t, loggerManager)
	assert.NotNil(t, loggerManager.logger)
	assert.Len(t, loggerManager.writers, 2)
}

func TestGetZapLevelFromStr(t *testing.T) {
	assert.Equal(t, zap.DebugLevel, getZapLevelFromStr("debug"))
	assert.Equal(t, zap.DebugLevel, getZapLevelFromStr("TRACE"))
	assert.Equal(t, zap.InfoLevel, getZapLevelFromStr("notice"))
	assert.Equal(t, zap.WarnLevel, getZapLevelFromStr("warn"))
	assert.Equal(t, zap.FatalLevel, getZapLevelFromStr("fatal"))
	assert.Equal(t, zapcore.Level(99), getZapLevelFromStr("none"))
	assert.Equal(t, zap.InfoLevel, getZapLevelFromStr("unknown"))
}

func TestSetLevel(t *testing.T) {
	loggerManager := &ZapLoggerManager{}
	err := loggerManager.SetLevel("test", "debug")
	assert.NoError(t, err)
}

func TestLoggerClose(t *testing.T) {
	loggerManager := &ZapLoggerManager{
		logger: zap.NewExample(),
		writers: []io.WriteCloser{
			&mockWriter{},
			&mockWriter{},
		},
	}
	loggerManager.Close()
	assert.Nil(t, loggerManager.logger)
	for _, writer := range loggerManager.writers {
		assert.Implements(t, (*io.WriteCloser)(nil), writer)
	}
}

type mockWriter struct{}

func (mw *mockWriter) Write(_ []byte) (n int, err error) { return }
func (mw *mockWriter) Close() error                      { return nil }

func TestMain(m *testing.M) {
	// 初始化一些全局资源，例如临时目录等
	code := m.Run()
	// 清理资源
	os.Exit(code)
}
