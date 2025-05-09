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
	"bufio"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// 默认的缓冲区大小和刷新间隔
	defaultBufferSize       = 1 * 1024 * 1024 // 1MB
	defaultBufferFlushIntvl = 50 * time.Millisecond
	// 默认的分片队列大小和数量
	defaultShardQueueSize = 256
	defaultShardCount     = 16
	// 日志丢弃策略
	defaultDiscardWhenFull = false
)

const (
	//ShardErrShortWrite means that a write accepted fewer bytes than requested but failed to return an explicit error.
	ShardErrShortWriteDrop = 0
)

// AsyncWriter implements a sharded, buffered writer for high-throughput logging.
// It uses concurrent shards with individual queues to minimize contention.
type AsyncWriter struct {
	// 分片配置
	shardCount           int
	shardQueueSize       int
	shardIndex           uint32
	shards               []*ShardWriter
	shardDiscardWhenFull bool   // 新增配置项
	shardDropped         uint64 // 统计丢弃数量(atomic)

	// buffer配置
	bufferSize          int
	bufferFlushInterval time.Duration
	closed              int32          // atomic
	writer              io.WriteCloser // 添加统一管理底层 writer

}

type Option func(*AsyncWriter)

// NewAsyncWriter creates a new async writer with the given io.WriteCloser and options.
// Defaults: 16 shards, 256 queue/shard, 1MB buffer, 50ms flush interval.
func NewAsyncWriter(w io.WriteCloser, opts ...Option) *AsyncWriter {
	aw := &AsyncWriter{
		shardCount:           defaultShardCount,
		shardQueueSize:       defaultShardQueueSize,
		shardIndex:           0,
		shardDiscardWhenFull: defaultDiscardWhenFull,
		shardDropped:         0,
		bufferSize:           defaultBufferSize,
		bufferFlushInterval:  defaultBufferFlushIntvl,
		writer:               w,
	}

	for _, opt := range opts {
		opt(aw)
	}

	aw.shards = make([]*ShardWriter, aw.shardCount)
	bufferedWriter := NewBufferWriter(w, aw.bufferSize, aw.bufferFlushInterval)
	var strategy WriteStrategy
	if aw.shardDiscardWhenFull {
		strategy = &DiscardStrategy{}
	} else {
		strategy = &BlockingStrategy{}
	}

	for i := range aw.shards {
		aw.shards[i] = NewShardWriter(
			w,
			aw.shardQueueSize,
			&aw.shardDropped,
			bufferedWriter,
			strategy,
		)
	}
	return aw
}

// Write distributes data to shards using round-robin selection. Concurrent-safe.
// Callers must ensure no writes after Close() (see Close documentation).
func (aw *AsyncWriter) Write(data []byte) (int, error) {
	idx := atomic.AddUint32(&aw.shardIndex, 1) % uint32(aw.shardCount)
	return aw.shards[idx].Write(data)
}

// Sync flushes all buffered data from all shards. Blocks until complete.
// Returns first encountered error, if any.
func (aw *AsyncWriter) Sync() error {
	for _, sw := range aw.shards {
		if err := sw.Sync(); err != nil {
			return err
		}
	}
	return nil
}

// Close implements graceful shutdown: clears the queue, flushes the buffer, and then closes
// Underlying writer. Idempotent. Subsequent writes after Close are unsafe and may cause panic
func (aw *AsyncWriter) Close() error {
	if !atomic.CompareAndSwapInt32(&aw.closed, 0, 1) {
		return nil
	}
	atomic.StoreInt32(&aw.closed, 1)

	// 先同步所有分片
	if err := aw.Sync(); err != nil {
		return err
	}

	// 然后关闭分片队列
	for _, sw := range aw.shards {
		sw.Close()
	}

	// 等待处理完成
	for _, sw := range aw.shards {
		sw.Wait()
	}
	return aw.writer.Close()
}

// Dropped returns the number of discarded log entries (thread-safe).
// Valid only when using discard strategy (WithDiscardWhenFull(true)).
func (aw *AsyncWriter) Dropped() uint64 {
	return atomic.LoadUint64(&aw.shardDropped)
}

// WriteStrategy defines the behavior when shard queues are full.
type WriteStrategy interface {
	Write(queue chan<- []byte, data []byte, dropped *uint64) (int, error)
}

// BlockingStrategy blocks producers until queue space is available.
type BlockingStrategy struct{}

func (s *BlockingStrategy) Write(queue chan<- []byte, data []byte, _ *uint64) (int, error) {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	queue <- dataCopy
	return len(dataCopy), nil
}

// DiscardStrategy drops new entries when queues are full, incrementing drop counter.
type DiscardStrategy struct{}

func (s *DiscardStrategy) Write(queue chan<- []byte, data []byte, dropped *uint64) (int, error) {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	select {
	case queue <- dataCopy:
		return len(dataCopy), nil
	default:
		atomic.AddUint64(dropped, 1)
		return ShardErrShortWriteDrop, nil
	}
}

// ShardWriter manages a single shard's write queue and buffer.
type ShardWriter struct {
	queue          chan []byte
	bufferedWriter *BufferWriter
	wg             sync.WaitGroup

	strategy WriteStrategy // 策略接口
	dropped  *uint64       // 共享统计
}

func NewShardWriter(w io.WriteCloser, queueSize int, dropped *uint64, bufferedWriter *BufferWriter, strategy WriteStrategy) *ShardWriter {
	sw := &ShardWriter{
		queue:          make(chan []byte, queueSize),
		bufferedWriter: bufferedWriter,
		wg:             sync.WaitGroup{},
		strategy:       strategy,
		dropped:        dropped,
	}

	sw.wg.Add(1)
	go func() {
		defer sw.wg.Done()
		sw.processWrite()
	}()

	return sw
}

// processWrite continuously drains the queue to the buffered writer.
// Runs in a dedicated goroutine until queue closure.
func (sw *ShardWriter) processWrite() {
	defer sw.bufferedWriter.Sync()
	for data := range sw.queue {
		sw.bufferedWriter.Write(data)
	}
}

func (sw *ShardWriter) Write(data []byte) (int, error) {
	return sw.strategy.Write(sw.queue, data, sw.dropped)
}

// Sync flushes in-memory data to persistent storage. May block during I/O.
// Handles remaining queue items before flushing the buffer.
func (sw *ShardWriter) Sync() error {
	for {
		select {
		case data, ok := <-sw.queue:
			if !ok {
				// 队列已关闭，由 processWrite 处理
				return nil
			}
			sw.bufferedWriter.Write(data)
		default:
			return sw.bufferedWriter.Sync()
		}
	}
}

// Close initiates graceful shutdown by closing the input queue.
// Pending data will be processed before final buffer flush.
func (sw *ShardWriter) Close() error {
	close(sw.queue) // 仅关闭队列，processWrite 会处理剩余数据后退出
	return nil
}

// Wait blocks until all queued data is processed and flushed.
// Must be called after Close for safe shutdown.
func (sw *ShardWriter) Wait() error {
	sw.wg.Wait()
	return nil
}

// BufferWriter implements auto-flushing buffered I/O with periodic sync.
type BufferWriter struct {
	sync.Mutex
	writer      io.Writer
	buf         *bufio.Writer
	flushTicker *time.Ticker
	done        chan struct{}
}

func NewBufferWriter(w io.WriteCloser, size int, flushIntvl time.Duration) *BufferWriter {
	bw := &BufferWriter{
		writer:      w,
		buf:         bufio.NewWriterSize(w, size),
		flushTicker: time.NewTicker(flushIntvl),
		done:        make(chan struct{}, 1),
	}
	go bw.processWrite()
	return bw
}

func (bw *BufferWriter) processWrite() {
	for {
		select {
		case <-bw.flushTicker.C:
			bw.Sync()
		case <-bw.done:
			return
		}
	}
}

// Write appends data to the buffer. Thread-safe via mutex.
// May block if buffer needs flushing to underlying writer.
func (bw *BufferWriter) Write(data []byte) (int, error) {
	bw.Lock()
	defer bw.Unlock()
	return bw.buf.Write(data)
}

// Sync immediately flushes buffered data. Caller should handle errors.
func (bw *BufferWriter) Sync() error {
	bw.Lock()
	defer bw.Unlock()
	return bw.buf.Flush()
}

// Close stops automatic flushing and performs final sync.
// Must be called to ensure data integrity on shutdown.
func (bw *BufferWriter) Close() error {
	bw.flushTicker.Stop()
	bw.done <- struct{}{}
	return bw.Sync()
}

func WithBufferFlushInterval(d time.Duration) Option {
	return func(aw *AsyncWriter) {
		aw.bufferFlushInterval = d
	}
}

func WithBufferSize(size int) Option {
	return func(aw *AsyncWriter) {
		if size > 0 {
			aw.bufferSize = size
		}
	}
}

func WithShardQueueSize(size int) Option {
	return func(aw *AsyncWriter) {
		if size > 0 {
			aw.shardQueueSize = size
		}
	}
}

func WithShardCount(count int) Option {
	return func(aw *AsyncWriter) {
		if count > 0 {
			aw.shardCount = count
		}
	}
}

func WithDiscardWhenFull(discard bool) Option {
	return func(aw *AsyncWriter) {
		aw.shardDiscardWhenFull = discard
	}
}
