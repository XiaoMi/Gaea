package zap

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

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
	w := NewAsyncWriter(&buf)
	data := []byte("Hello, World!")
	// Write some data to trigger flush
	_, _ = w.Write(data)
	// Wait for the interval flush to occur
	time.Sleep(200 * time.Millisecond)
	// Check if data has been flushed
	if !bytes.Contains(buf.Bytes(), data) {
		t.Error("Data was not flushed within expected interval")
	}
}

func TestWriterClose(t *testing.T) {
	var buf mockBuffer
	w := NewAsyncWriter(&buf)
	w.timer = time.NewTicker(time.Millisecond * 100)
	w.ctx, w.quit = context.WithCancel(context.Background())
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
	w.timer = time.NewTicker(time.Millisecond * 100)
	w.ctx, w.quit = context.WithCancel(context.Background())
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
	// Close the writer to ensure all data is flushed
	_ = w.Close()
	// Check if all data has been written
	expected := strings.Repeat(string(data), 10)
	if !bytes.Equal(mwc.written, []byte(expected)) {
		t.Errorf("Written data does not match expected: got %s, want %s", string(mwc.written), expected)
	}
}
