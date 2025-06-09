// Copyright 2025 The Gaea Authors. All Rights Reserved.
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
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

// 带关闭功能的缓冲区
type bufferCloser struct {
	*bytes.Buffer
}

func (b *bufferCloser) Close() error { return nil }

func newBufferCloser() *bufferCloser {
	return &bufferCloser{bytes.NewBuffer(nil)}
}

func TestAsyncWriter_ConcurrentStress(t *testing.T) {
	bc := newBufferCloser()
	// 扩大队列容量避免丢弃
	writer := NewAsyncWriter(
		bc,
		WithShardQueueSize(1000),
		WithShardCount(5),
	)

	const (
		goroutines = 50
		perGoCount = 100
	)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < perGoCount; j++ {
				if _, err := writer.Write([]byte("concurrent data")); err != nil {
					t.Logf("Write failed: %v", err)
				}
			}
		}()
	}

	wg.Wait()
	writer.Close()

	expected := goroutines * perGoCount * len("concurrent data")
	if bc.Len() != expected {
		t.Fatalf("Expected %d bytes, got %d", expected, bc.Len())
	}
}

// 测试并发写入数据完整性和正确性(1.完整存在于缓冲区 2.仅出现一次 3.没有被污染)
func TestAsyncWriter_ConcurrentDataIntegrity(t *testing.T) {
	// 初始化写入器和缓冲区
	bc := newBufferCloser()
	writer := NewAsyncWriter(bc)
	defer writer.Close()

	// 测试参数配置
	concurrency := 100 // 并发协程数
	writeCount := 100  // 每个协程写入次数

	// 启动并发写入
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writeCount; j++ {
				// 生成唯一可识别数据
				originalData := []byte(fmt.Sprintf("data-%02d-%03d", id, j))

				// 校验数据长度
				if len(originalData) != 11 {
					t.Errorf("Invalid data length: %d", len(originalData))
					return
				}

				// 保存原始数据副本
				expected := string(originalData)

				// 执行写入
				if _, err := writer.Write(originalData); err != nil {
					t.Errorf("Write failed: %v", err)
					return
				}

				// 仅污染数据内容（删除长度修改）
				for k := range originalData {
					originalData[k] = 'X'
				}

				// 验证内容污染成功
				if strings.Contains(string(originalData), expected) {
					t.Error("Data pollution failed: original data still exists")
				}
			}
		}(i)
	}

	// 等待所有写入完成
	wg.Wait()

	// 强制同步并关闭
	writer.Close()

	// 结果校验
	result := bc.String()

	// 检查数据完整性
	for i := 0; i < concurrency; i++ {
		for j := 0; j < writeCount; j++ {
			expected := fmt.Sprintf("data-%02d-%03d", i, j)
			if !strings.Contains(result, expected) {
				t.Errorf("Missing expected data: %s", expected)
			}
		}
	}

	// 校验总量
	expectedSize := concurrency * writeCount * 11 // 11为每条数据固定长度
	if len(result) != expectedSize {
		t.Errorf("Unexpected data length, got %d want %d", len(result), expectedSize)
	}
}

type mockBuffer struct {
	bytes.Buffer
}

func (mockBuffer) Close() error {
	return nil
}

func TestNewAsyncWriter(t *testing.T) {
	var buf mockBuffer
	w := NewAsyncWriter(&buf)
	if w == nil {
		t.Error("Expected non-nil AsyncWriter")
	}
}

func TestWriteAndSync(t *testing.T) {
	var buf mockBuffer
	w := NewAsyncWriter(&buf)
	data := []byte("Hello, World!")
	n, err := w.Write(data)
	if n != len(data) || err != nil {
		t.Errorf("Write failed: %d, %v", n, err)
	}
	err = w.Sync()
	if err != nil {
		t.Errorf("Sync failed: %v", err)
	}
}

func TestIntervalFlush(t *testing.T) {
	var buf mockBuffer
	w := NewAsyncWriter(&buf) // 设置全局刷新间隔为 100ms

	data := []byte("Hello, World!")
	// Write some data to trigger flush
	_, _ = w.Write(data)
	// Wait for the interval flush to occur
	time.Sleep(200 * time.Millisecond)
	// w.Sync()
	// Check if data has been flushed
	if !bytes.Contains(buf.Bytes(), data) {
		t.Error("Data was not flushed within expected interval")
	}
}

func TestWriterClose(t *testing.T) {
	var buf mockBuffer
	w := NewAsyncWriter(&buf)
	data := []byte("Hello, World!")
	_, _ = w.Write(data)
	err := w.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
	// Ensure no data is lost
	if !bytes.Contains(buf.Bytes(), data) {
		t.Error("Data was lost during close")
	}
}

type mockWriteCloser struct {
	mu      sync.Mutex
	written []byte
}

func (mwc *mockWriteCloser) Write(p []byte) (n int, err error) {
	mwc.mu.Lock()
	defer mwc.mu.Unlock()
	mwc.written = append(mwc.written, p...)
	return len(p), nil
}

func (mwc *mockWriteCloser) Close() error {
	return nil
}

func TestIntegration(t *testing.T) {
	mwc := &mockWriteCloser{}
	w := NewAsyncWriter(mwc)
	data := []byte("Hello, World!")
	// Simulate concurrent writes
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			_, _ = w.Write(data)
		}()
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	// Check if all data has been written
	expected := strings.Repeat(string(data), 10)
	if !bytes.Equal(mwc.written, []byte(expected)) {
		t.Errorf("Written data does not match expected: got %s, want %s", string(mwc.written), expected)
	}
}

// 模拟错误写入器（实现io.WriteCloser）
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func (w *errorWriter) Close() error {
	return nil
}

func TestAsyncWriter_ErrorHandling(t *testing.T) {
	ew := &errorWriter{}
	writer := NewAsyncWriter(ew)

	// 多次写入触发错误
	for i := 0; i < 10; i++ {
		if _, err := writer.Write([]byte("test")); err != nil {
			// t.Fatal("Unexpected write error:", err)
		}
	}

	// 验证同步时是否处理错误
	if err := writer.Sync(); err != nil {
		t.Log("Sync returned expected error:", err)
	}
}

type errorCloseWriter struct{}

func (w *errorCloseWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func (w *errorCloseWriter) Close() error {
	return fmt.Errorf("Close error")
}

func newErrorCloseWriter() *errorCloseWriter {
	return &errorCloseWriter{}
}
func TestAsyncWriter_ErrorClose(t *testing.T) {
	ew := newErrorCloseWriter()
	writer := NewAsyncWriter(ew)

	// 多次写入触发错误
	for i := 0; i < 10; i++ {
		if _, err := writer.Write([]byte("test")); err != nil {
			// t.Fatal("Unexpected write error:", err)
		}
	}

	if err := writer.Close(); err == nil {
		t.Fatal("Expected close error:", err)
	}

}

// 测试正常写入不丢弃
func TestAsyncWriter_NormalWrite(t *testing.T) {
	bc := newBufferCloser()
	writer := NewAsyncWriter(bc)

	testData := []byte("normal write test")
	if _, err := writer.Write(testData); err != nil {
		t.Fatal("Write failed:", err)
	}

	writer.Close()

	if !strings.Contains(bc.String(), string(testData)) {
		t.Fatal("Data not written to buffer")
	}
	if writer.Dropped() > 0 {
		t.Fatal("Unexpected dropped logs in normal mode")
	}
}

// 测试队列满时丢弃日志
func TestAsyncWriter_Discard(t *testing.T) {
	bc := newBufferCloser()

	// 使用极小的队列容量快速触发丢弃
	writer := NewAsyncWriter(
		bc,
		WithStrategy(strategyDiscard),
		WithShardQueueSize(1), // 每个分片队列容量1
		WithShardCount(1),     // 单分片简化测试
		WithBufferFlushInterval(time.Millisecond),
	)

	const (
		totalWrites  = 1000            // 总写入次数
		payloadSize  = 1024            // 每次写入大小
		expectedDrop = totalWrites - 1 // 队列容量1，最多保留1条
	)

	// 快速写入数据（不等待处理）
	for i := 0; i < totalWrites; i++ {
		_, _ = writer.Write(make([]byte, payloadSize))
	}

	// 等待处理完成
	writer.Close()

	// 验证丢弃数量
	dropped := writer.Dropped()

	// 验证实际写入量
	written := bc.Len()
	// 修改后的断言逻辑
	expectedTotal := totalWrites * payloadSize
	actualTotal := written + int(dropped)*payloadSize
	if actualTotal != expectedTotal {
		t.Fatalf("Expected written(%d) + dropped(%d*%d) = %d, got %d", written, dropped, payloadSize, expectedTotal, actualTotal)
	}

}

func TestAsyncWriter_DiscardWithDowngrade(t *testing.T) {
	mbc := newBufferCloser()
	dbc := newBufferCloser()

	// 使用极小的队列容量快速触发丢弃
	writer := NewAsyncWriter(
		mbc,
		WithStrategy(strategyDiscardWithDowngrade),
		WithDefaultDowngradeWriter(dbc, defaultBufferSize, defaultBufferFlushIntvl),
		WithShardQueueSize(1), // 每个分片队列容量1
		WithShardCount(1),     // 单分片简化测试
		WithBufferFlushInterval(time.Millisecond),
	)

	const (
		totalWrites  = 1000            // 总写入次数
		payloadSize  = 1024            // 每次写入大小
		expectedDrop = totalWrites - 1 // 队列容量1，最多保留1条
	)

	// 快速写入数据（不等待处理）
	for i := 0; i < totalWrites; i++ {
		_, _ = writer.Write(make([]byte, payloadSize))
	}

	// 等待处理完成
	writer.Close()

	// 验证丢弃数量
	// dropped := writer.Dropped()

	// 验证实际写入量
	mwritten := mbc.Len()
	dwritten := dbc.Len()
	// 修改后的断言逻辑
	expectedTotal := totalWrites * payloadSize
	actualTotal := mwritten + dwritten
	if actualTotal != expectedTotal {
		t.Fatalf("Expected Total written(%d), Actual Total written(%d), master written:(%d),  downgradeWriter written:(%d)", expectedTotal, actualTotal, mwritten, dwritten)
	}
}

func TestAsyncWriter_DiscardWithDowngradeWithNil(t *testing.T) {
	mbc := newBufferCloser()
	// 使用极小的队列容量快速触发丢弃
	writer := NewAsyncWriter(
		mbc,
		WithDefaultDowngradeWriter(nil, defaultBufferSize, defaultBufferFlushIntvl),
		WithStrategy(strategyDiscardWithDowngrade),
		WithShardQueueSize(1), // 每个分片队列容量1
		WithShardCount(1),     // 单分片简化测试
		WithBufferFlushInterval(time.Millisecond),
	)

	const (
		totalWrites  = 1000            // 总写入次数
		payloadSize  = 1024            // 每次写入大小
		expectedDrop = totalWrites - 1 // 队列容量1，最多保留1条
	)

	// 快速写入数据（不等待处理）
	for i := 0; i < totalWrites; i++ {
		_, _ = writer.Write(make([]byte, payloadSize))
	}

	// 等待处理完成
	writer.Close()

	// 验证丢弃数量
	dropped := writer.Dropped()

	// 验证实际写入量
	written := mbc.Len()

	// 修改后的断言逻辑
	expectedTotal := totalWrites * payloadSize
	actualTotal := written + int(dropped)*payloadSize
	if actualTotal != expectedTotal {
		t.Fatalf("Expected written(%d) + dropped(%d*%d) = %d, got %d", written, dropped, payloadSize, expectedTotal, actualTotal)
	}
}

type slowWriter struct {
	delay time.Duration
	buf   *bytes.Buffer
}

func (w *slowWriter) Close() error {
	return nil
}
func (w *slowWriter) Write(p []byte) (n int, err error) {
	time.Sleep(w.delay) // 模拟 IO 延迟
	return w.buf.Write(p)
}
func (w *slowWriter) String() string {
	return w.buf.String()
}

func newSlowWriterCloser(delay time.Duration) *slowWriter {
	return &slowWriter{
		delay: delay,
		buf:   bytes.NewBuffer(nil),
	}
}

// 测试阻塞情况下数据完整性（1.已写入数据的正确性 2.数据无重复 3.污染防御有效性）
func TestAsyncWriterSlow_Block_ConcurrentDataIntegrity(t *testing.T) {
	// 初始化写入器和缓冲区
	bc := newSlowWriterCloser(100 * time.Microsecond)
	writer := NewAsyncWriter(bc)
	defer writer.Close()
	// 测试参数配置
	concurrency := 100 // 并发协程数
	writeCount := 100  // 每个协程写入次数
	// 启动并发写入
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writeCount; j++ {
				// 生成唯一可识别数据
				originalData := []byte(fmt.Sprintf("data-%02d-%03d", id, j))
				// 校验数据长度
				if len(originalData) != 11 {
					t.Errorf("Invalid data length: %d", len(originalData))
					return
				}
				// 保存原始数据副本
				expected := string(originalData)
				// 执行写入
				if _, err := writer.Write(originalData); err != nil {
					t.Errorf("Write failed: %v", err)
					return
				}
				// 仅污染数据内容（删除长度修改）
				for k := range originalData {
					originalData[k] = 'X'
				}
				// 验证内容污染成功
				if strings.Contains(string(originalData), expected) {
					t.Error("Data pollution failed: original data still exists")
				}
			}
		}(i)
	}
	// 等待所有写入完成
	wg.Wait()
	// 强制同步并关闭
	writer.Close()
	// 结果校验
	result := bc.String()
	// 检查数据完整性
	for i := 0; i < concurrency; i++ {
		for j := 0; j < writeCount; j++ {
			expected := fmt.Sprintf("data-%02d-%03d", i, j)
			if !strings.Contains(result, expected) {
				t.Errorf("Missing expected data: %s", expected)
			}
		}
	}
	// 校验总量
	expectedSize := concurrency * writeCount * 11 // 11为每条数据固定长度
	if len(result) != expectedSize {
		t.Errorf("Unexpected data length, got %d want %d", len(result), expectedSize)
	}
}

func TestAsyncWriterNotSlow_Block_ConcurrentDataIntegrity(t *testing.T) {
	// 初始化写入器和缓冲区
	bc := newBufferCloser()
	writer := NewAsyncWriter(bc)
	defer writer.Close()
	// 测试参数配置
	concurrency := 100 // 并发协程数
	writeCount := 100  // 每个协程写入次数
	// 启动并发写入
	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writeCount; j++ {
				// 生成唯一可识别数据
				originalData := []byte(fmt.Sprintf("data-%02d-%03d", id, j))
				// 校验数据长度
				if len(originalData) != 11 {
					t.Errorf("Invalid data length: %d", len(originalData))
					return
				}
				// 保存原始数据副本
				expected := string(originalData)
				// 执行写入
				if _, err := writer.Write(originalData); err != nil {
					t.Errorf("Write failed: %v", err)
					return
				}
				// 仅污染数据内容（删除长度修改）
				for k := range originalData {
					originalData[k] = 'X'
				}
				// 验证内容污染成功
				if strings.Contains(string(originalData), expected) {
					t.Error("Data pollution failed: original data still exists")
				}
			}
		}(i)
	}
	// 等待所有写入完成
	wg.Wait()
	// 强制同步并关闭
	writer.Close()
	// 结果校验
	result := bc.String()
	// 检查数据完整性
	for i := 0; i < concurrency; i++ {
		for j := 0; j < writeCount; j++ {
			expected := fmt.Sprintf("data-%02d-%03d", i, j)
			if !strings.Contains(result, expected) {
				t.Errorf("Missing expected data: %s", expected)
			}
		}
	}
	// 校验总量
	expectedSize := concurrency * writeCount * 11 // 11为每条数据固定长度
	if len(result) != expectedSize {
		t.Errorf("Unexpected data length, got %d want %d", len(result), expectedSize)
	}
}

func TestBufferCloserSpeed(t *testing.T) {
	bc := newBufferCloser()
	start := time.Now()
	for i := 0; i < 1_000_000; i++ {
		bc.Write([]byte("test"))
	}
	t.Logf("Write speed: %.1f ns/op", time.Since(start).Seconds()*1e9/1e6)
}

// 情况1：测试阻塞模式下的行为
func TestAsyncWriter_SlowWriterBackPressure(t *testing.T) {
	bc := newSlowWriterCloser(100 * time.Microsecond)
	writer := NewAsyncWriter(
		bc,
		WithShardQueueSize(10),
		WithShardCount(1),
		WithBufferSize(0), // 禁用缓冲
	)

	// 验证写入阻塞不丢数据
	for i := 0; i < 100; i++ {
		writer.Write([]byte("data"))
	}

	if writer.Dropped() > 0 {
		t.Fatal("Dropped logs in blocking mode")
	}
}

// 情况2：测试丢弃模式下的背压
func TestAsyncWriter_DiscardUnderPressure(t *testing.T) {
	slowWriter := &slowWriter{delay: 100 * time.Millisecond}
	writer := NewAsyncWriter(
		slowWriter,
		WithStrategy(strategyDiscard),
		WithShardQueueSize(10),
		WithShardCount(1),
	)

	// 快速写入触发丢弃
	for i := 0; i < 100; i++ {
		writer.Write([]byte("data"))
	}

	if writer.Dropped() == 0 {
		t.Fatal("Expected dropped logs under pressure")
	}
}
