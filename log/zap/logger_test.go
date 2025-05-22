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
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCreateLogManager(t *testing.T) {
	config := map[string]string{
		"path":            "/tmp",
		"filename":        "test",
		"level":           "debug",
		"log_keep_days":   "7",
		"log_keep_counts": "30",
		"log_strategy":    "",
		"log_local_path":  "",
	}

	loggerManager, err := CreateLogManager(config)
	assert.NoError(t, err)
	assert.NotNil(t, loggerManager)
	assert.NotNil(t, loggerManager.logger)
	assert.Len(t, loggerManager.writers, 2)
}

// 辅助函数：带重试的文件存在检查
func assertFileExists(t *testing.T, path string) {
	t.Helper()
	const maxRetries = 5
	const delay = 100 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(delay)
	}
	t.Errorf("File not created: %s", path)
}

func TestCreateLogManagerAsyncWriter(t *testing.T) {
	// 正常创建带降级日志
	t.Run("Normal creation with downgrade log", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpLocalDir := t.TempDir()

		config := map[string]string{
			"path":           tmpDir,
			"log_local_path": tmpLocalDir,
			"filename":       "test",
			"level":          "debug",
			"log_strategy":   "non-blocking",
		}

		mgr, err := CreateLogManager(config)
		if err != nil {
			t.Fatalf("Creation failed: %v", err)
		}

		// 1. 立即写入测试日志
		mgr.logger.Debug("test debug message")
		mgr.logger.Error("test error message")

		// 2. 显式同步日志
		if syncErr := mgr.logger.Sync(); syncErr != nil {
			t.Fatalf("Log synchronization failed: %v", syncErr)
		}

		// 3. 先关闭再检查（确保缓冲区刷新）
		mgr.Close()

		// 4.验证文件存在
		assertFileExists(t, filepath.Join(tmpDir, "test.log"))
	})

	t.Run("No downgrade log path", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := map[string]string{
			"path":     tmpDir,
			"filename": "test",
			"level":    "info",
		}

		mgr, err := CreateLogManager(config)
		if err != nil {
			t.Fatalf("Creation failed: %v", err)
		}
		defer mgr.Close()

		// 验证没有本地日志目录创建
		localLog := filepath.Join("", "test.log")
		if _, err := os.Stat(localLog); err == nil {
			t.Errorf("Local log file created unexpectedly: %s", localLog)
		}
	})

	t.Run("Missing required parameters", func(t *testing.T) {
		testCases := []struct {
			name   string
			config map[string]string
			expect string
		}{
			{"missing path", map[string]string{"filename": "test"}, "not found path"},
			{"missing filename", map[string]string{"path": "/tmp"}, "not found filename"},
			{"missing level", map[string]string{"path": "/tmp", "filename": "test"}, "not found level"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := CreateLogManager(tc.config)
				if err == nil || !strings.Contains(err.Error(), tc.expect) {
					t.Errorf("The expected error contains %q, actually get %v", tc.expect, err)
				}
			})
		}
	})

	t.Run("Log policy configuration", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := map[string]string{
			"path":         tmpDir,
			"filename":     "test",
			"level":        "warn",
			"log_strategy": "discard",
		}

		mgr, err := CreateLogManager(config)
		if err != nil {
			t.Fatalf("Creation failed: %v", err)
		}
		defer mgr.Close()
		aw, ok := mgr.writers[0].(*AsyncWriter)
		if ok && aw.Strategy() != strategyBlocking {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

}

// 测试启用降级日志的场景
func TestCreateWithDowngrade(t *testing.T) {
	tmpDir := t.TempDir()
	tmpLocalDir := t.TempDir()

	config := map[string]string{
		"path":            tmpDir,
		"log_local_path":  tmpLocalDir,
		"filename":        "downgrade_test",
		"level":           "info",
		"log_strategy":    "multi",
		"log_keep_days":   "7",
		"log_keep_counts": "30",
	}

	mgr, err := CreateLogManager(config)
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	defer mgr.Close()

	// 验证异步写入器配置
	t.Run("Verify infoWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[0].(*AsyncWriter)
		if ok && aw.Strategy() != strategyDiscardWithDowngrade {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

	t.Run("Verify warnWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[1].(*AsyncWriter)
		if ok && aw.Strategy() != strategyDiscardWithDowngrade {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

}

func TestCreateWithDowngradeWithLocalPath(t *testing.T) {
	tmpDir := t.TempDir()
	config := map[string]string{
		"path":            tmpDir,
		"log_local_path":  "",
		"filename":        "downgrade_test",
		"level":           "info",
		"log_strategy":    "multi",
		"log_keep_days":   "7",
		"log_keep_counts": "30",
	}

	mgr, err := CreateLogManager(config)
	if err != nil {
		t.Fatalf("Creation failed: %v", err)
	}
	defer mgr.Close()

	// 验证异步写入器配置
	t.Run("Verify infoWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[0].(*AsyncWriter)
		if ok && aw.Strategy() != strategyDiscard {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

	t.Run("Verify warnWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[1].(*AsyncWriter)
		if ok && aw.Strategy() != strategyDiscard {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

}

func TestCreateWithDiscard(t *testing.T) {
	tmpDir := t.TempDir()
	config := map[string]string{
		"path":            tmpDir,
		"filename":        "downgrade_test",
		"level":           "info",
		"log_local_path":  "",
		"log_strategy":    "async",
		"log_keep_days":   "7",
		"log_keep_counts": "30",
	}

	mgr, err := CreateLogManager(config)
	if err != nil {
		t.Fatalf("Creation failed: %v", err)
	}
	defer mgr.Close()

	// 验证异步写入器配置
	t.Run("Verify infoWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[0].(*AsyncWriter)
		if ok && aw.Strategy() != strategyDiscard {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

	t.Run("Verify warnWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[1].(*AsyncWriter)
		if ok && aw.Strategy() != strategyDiscard {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

}

func TestCreateWithBlock(t *testing.T) {
	tmpDir := t.TempDir()

	config := map[string]string{
		"path":            tmpDir,
		"filename":        "downgrade_test",
		"level":           "info",
		"log_local_path":  "",
		"log_strategy":    "sync",
		"log_keep_days":   "7",
		"log_keep_counts": "30",
	}

	mgr, err := CreateLogManager(config)
	if err != nil {
		t.Fatalf("Creation failed: %v", err)
	}
	defer mgr.Close()

	// 验证异步写入器配置
	t.Run("Verify infoWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[0].(*AsyncWriter)
		if ok && aw.Strategy() != strategyBlocking {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

	t.Run("Verify warnWriter downgrade configuration", func(t *testing.T) {
		aw, ok := mgr.writers[1].(*AsyncWriter)
		if ok && aw.Strategy() != strategyBlocking {
			t.Errorf("Expected policy is discard, actual policy is %v", aw.Strategy())
		}
	})

}

// 测试不启用降级日志的场景
func TestCreateWithoutDowngrade(t *testing.T) {
	tmpDir := t.TempDir()

	config := map[string]string{
		"path":         tmpDir,
		"filename":     "no_downgrade_test",
		"level":        "debug",
		"log_strategy": "blocking",
	}

	mgr, err := CreateLogManager(config)
	if err != nil {
		t.Fatalf("Creation failed: %v", err)
	}
	defer mgr.Close()

	// 验证异步写入器配置
	t.Run("Verify that infoWriter is not downgraded", func(t *testing.T) {
		verifyNoDowngrade(t, mgr.writers[0])
	})

	t.Run("验证warnWriter无降级", func(t *testing.T) {
		verifyNoDowngrade(t, mgr.writers[1])
	})
}

func verifyNoDowngrade(t *testing.T, writer AsyncWriterCloser) {
	aw, ok := writer.(*AsyncWriter)
	if !ok {
		t.Fatal("写入器类型断言失败")
	}

	if aw.downgradeWriter != nil {
		t.Error("不启用降级时应无降级写入器")
	}
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
		writers: []AsyncWriterCloser{
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
func (mw *mockWriter) Dropped() uint64                   { return 0 }

func TestMain(m *testing.M) {
	// 初始化一些全局资源，例如临时目录等
	code := m.Run()
	// 清理资源
	os.Exit(code)
}
