package zap

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func TestZapEncoder_EncodeEntry(t *testing.T) {
	zapBufferPool = buffer.NewPool()
	e := &ZapEncoder{}

	// 创建一个模拟的时间和日志条目
	now := time.Now()
	entry := zapcore.Entry{
		Time:    now,
		Level:   zapcore.InfoLevel,
		Message: "test message",
	}

	// 调用 EncodeEntry 方法并获取结果缓冲区
	buf, err := e.EncodeEntry(entry, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// 检查缓冲区的内容是否符合预期
	expectedOutput := fmt.Sprintf("[%s] [%s] %s\n", now.Format("2006-01-02 15:04:05.000"), strings.ToUpper("info"), "test message")
	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\ngot:\n%s", expectedOutput, buf.String())
	}
}

func TestZapEncoder_Clone(t *testing.T) {
	e := &ZapEncoder{}
	clone := e.Clone()
	if clone == nil || clone != e {
		t.Error("Clone method did not return a same instance")
	}
}

func BenchmarkZapEncoder_EncodeEntry(b *testing.B) {
	zapBufferPool = buffer.NewPool()
	e := &ZapEncoder{}

	// 创建一个模拟的时间和日志条目
	now := time.Now()
	entry := zapcore.Entry{
		Time:    now,
		Level:   zapcore.InfoLevel,
		Message: "benchmark test message",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = e.EncodeEntry(entry, nil)
	}
}
