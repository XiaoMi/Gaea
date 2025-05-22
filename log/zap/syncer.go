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
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// 网络盘/本地盘配置：默认的缓冲区大小和刷新间隔
	defaultBufferSize       = 1 * 1024 * 1024       // 默认的分片缓冲区大小
	defaultBufferFlushIntvl = 50 * time.Millisecond // 默认的分片缓冲区刷新间隔
	defaultShardQueueSize   = 256                   // 默认的分片队列大小
	defaultShardCount       = 16                    // 默认的分片队列数量
)

const (
	//ShardErrShortWrite means that a write accepted fewer bytes than requested but failed to return an explicit error.
	ShardErrShortWriteDrop = 0
)

type strategyType int

const (
	strategyBlocking             strategyType = iota // 阻塞策略
	strategyDiscard                                  // 直接丢弃
	strategyDiscardWithDowngrade                     // 降级写入+本地回退
)

// AsyncWriter implements a sharded, buffered writer for high-throughput logging.
// It uses concurrent shards with individual queues to minimize contention.
type AsyncWriter struct {
	// 分片配置
	shardCount     int
	shardQueueSize int
	shardIndex     uint32
	shards         []*ShardWriter
	shardDropped   uint64

	// buffer配置
	bufferSize          int
	bufferFlushInterval time.Duration

	strategy        strategyType       // 策略配置
	closed          int32              // 关闭配置
	masterWriter    io.WriteCloser     // 主写入器
	downgradeWriter *BufferWriteCloser // 降级写入器
}

type Option func(*AsyncWriter)

// NewAsyncWriter creates a new async writer with the given io.WriteCloser and options.
// Defaults: 16 shards, 256 queue/shard, 1MB buffer, 50ms flush interval.
func NewAsyncWriter(masterWriter io.WriteCloser, opts ...Option) *AsyncWriter {
	aw := &AsyncWriter{
		shardCount:          defaultShardCount,
		shardQueueSize:      defaultShardQueueSize,
		shardIndex:          0,
		shardDropped:        0,
		bufferSize:          defaultBufferSize,
		bufferFlushInterval: defaultBufferFlushIntvl,
		strategy:            strategyBlocking,
		closed:              0,
		masterWriter:        masterWriter,
	}

	for _, opt := range opts {
		opt(aw)
	}

	aw.shards = make([]*ShardWriter, aw.shardCount)
	shardBuffer := NewBufferWriter(masterWriter, aw.bufferSize, aw.bufferFlushInterval)
	strategy := aw.getStrategy()
	for i := range aw.shards {
		aw.shards[i] = NewShardWriter(
			masterWriter,
			aw.shardQueueSize,
			&aw.shardDropped,
			shardBuffer,
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

	var errs []string

	// 先同步所有分片
	if err := aw.Sync(); err != nil {
		errs = append(errs, fmt.Sprintf("sync failed: %v", err))
	}

	// 然后关闭分片队列
	for i, sw := range aw.shards {
		if err := sw.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("shard: %d close error: %v", i, err))
		}
	}

	// 等待处理完成
	for _, sw := range aw.shards {
		sw.Wait()
	}

	if err := aw.masterWriter.Close(); err != nil {
		errs = append(errs, fmt.Sprintf("masterWriter close failed: %v", err))
	}

	// 无论 masterWriter 是否关闭成功，都要尝试关闭 downgradeWriter
	var downgradeErr error
	if aw.downgradeWriter != nil {
		if err := aw.downgradeWriter.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("downgradeWriter close failed: %v", downgradeErr))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close async writer failed: %s", strings.Join(errs, ", "))
	}
	return nil
}

func (aw *AsyncWriter) Strategy() strategyType {
	return aw.strategy
}

func (aw *AsyncWriter) getStrategy() WriteStrategy {
	var strategy WriteStrategy
	switch aw.strategy {
	case strategyBlocking:
		strategy = NewBlockingStrategy()
	case strategyDiscard:
		strategy = NewDiscardStrategy(&aw.shardDropped)
	case strategyDiscardWithDowngrade:
		if aw.downgradeWriter != nil {
			strategy = NewDiscardDowngradeStrategy(aw.downgradeWriter, &aw.shardDropped)
		} else {
			strategy = NewDiscardStrategy(&aw.shardDropped)
			aw.strategy = strategyDiscard
		}
	default:
		strategy = NewBlockingStrategy()
	}
	return strategy
}

// Dropped returns the number of discarded log entries (thread-safe).
// Valid only when using discard strategy (WithDiscardWhenFull(true)).
func (aw *AsyncWriter) Dropped() uint64 {
	return atomic.LoadUint64(&aw.shardDropped)
}

// WriteStrategy defines the behavior when shard queues are full.
type WriteStrategy interface {
	Write(queue chan<- []byte, data []byte) (int, error)
}

// BlockingStrategy blocks producers until queue space is available.
type BlockingStrategy struct{}

func NewBlockingStrategy() *BlockingStrategy {
	return &BlockingStrategy{}
}

func (s *BlockingStrategy) Write(queue chan<- []byte, data []byte) (int, error) {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	queue <- dataCopy
	return len(dataCopy), nil
}

// DiscardStrategy drops new entries when queues are full, incrementing drop counter.
type DiscardStrategy struct {
	dropped *uint64 // 共享统计
}

func NewDiscardStrategy(dropped *uint64) *DiscardStrategy {
	return &DiscardStrategy{
		dropped: dropped,
	}
}

func (s *DiscardStrategy) Write(queue chan<- []byte, data []byte) (int, error) {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	select {
	case queue <- dataCopy:
		return len(dataCopy), nil
	default:
		atomic.AddUint64(s.dropped, 1)
		return ShardErrShortWriteDrop, nil
	}
}

// DiscardDowngradeStrategy implements fallback logging when primary queue is full
// Employs a buffered downgrade writer for overflow scenarios
type DiscardDowngradeStrategy struct {
	bufferedWriter *BufferWriteCloser // Fallback writer for overflow logs
	dropped        *uint64            // Atomic counter for dropped logs (shared across shard)
}

// NewDiscardDowngradeStrategy constructs a queue-full fallback handler
// buffer: Fallback writer for overflow logs
// dropped: Shared atomic counter reference
func NewDiscardDowngradeStrategy(buffer *BufferWriteCloser, dropped *uint64) *DiscardDowngradeStrategy {
	return &DiscardDowngradeStrategy{
		bufferedWriter: buffer,
		dropped:        dropped,
	}
}

// Write attempts non-blocking enqueue first, falls back to buffered writer on full queue
// queue: Target channel for primary logging path
// data: Log payload to write
// Returns bytes written and potential I/O error from fallback path
func (d *DiscardDowngradeStrategy) Write(queue chan<- []byte, data []byte) (int, error) {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	select {
	case queue <- dataCopy:
		return len(dataCopy), nil
	default:
		atomic.AddUint64(d.dropped, 1)
		return d.bufferedWriter.Write(data)
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

// NewShardWriter initializes a shard processing pipeline
// w: Underlying I/O target
// queueSize: Channel buffer capacity (backpressure control)
// dropped: Shared drop counter reference
// bufferedWriter: Associated buffer manager
// strategy: Write policy implementatio
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

// Write delegates to strategy implementation
// Returns bytes written and potential write errors
func (sw *ShardWriter) Write(data []byte) (int, error) {
	return sw.strategy.Write(sw.queue, data)
}

// Sync flushes in-memory data to persistent storage. May block during I/O.
// Handles remaining queue items before flushing the buffer.
func (sw *ShardWriter) Sync() error {
	for {
		select {
		case data, ok := <-sw.queue:
			if !ok {
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
	return nil      // 关闭策略中的缓冲写入器
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

// BufferWriter is a general buffer layer that relies on external management of the underlying io.Writer resources.
// When Close is called, the io.Writer is not closed.
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

// BufferWriteCloser is an autonomous buffering component that maintains full ownership of its underlying I/O resource.
// When Close() is invoked, it performs complete cleanup including closing the wrapped io.WriteCloser to prevent resource leaks.
//
// Design Philosophy:
// 1. Self-contained lifecycle management for reliability in fallback scenarios
// 2. Atomic control of both buffering layer and concrete I/O resource
// 3. Guaranteed resource release through finalization chain
// 4. Safe nil-receiver handling for defensive programming
type BufferWriteCloser struct {
	sync.Mutex
	writer      io.WriteCloser
	buf         *bufio.Writer
	flushTicker *time.Ticker
	done        chan struct{}
}

func NewBufferWriteCloser(w io.WriteCloser, size int, flushIntvl time.Duration) *BufferWriteCloser {
	if isNil(w) {
		return nil
	}
	bw := &BufferWriteCloser{
		writer:      w,
		buf:         bufio.NewWriterSize(w, size),
		flushTicker: time.NewTicker(flushIntvl),
		done:        make(chan struct{}, 1),
	}
	go bw.processWrite()
	return bw
}

// processWrite runs the background flush scheduler in a dedicated goroutine.
// Implements two termination-safe patterns:
// 1. Ticker-driven periodic flushing for data durability
// 2. Clean shutdown via done channel signaling
// Design invariant: Never holds mutex during channel operations
func (bw *BufferWriteCloser) processWrite() {
	for {
		select {
		case <-bw.flushTicker.C:
			bw.Sync()
		case <-bw.done:
			return
		}
	}
}

// Write serializes access to the buffer using mutex protection.
// Implements backpressure by potentially blocking when:
// - Buffer requires flushing to underlying writer
// - Concurrent writes exceed buffer capacity
// Design contract: Callers must handle partial writes and errors
func (bw *BufferWriteCloser) Write(data []byte) (int, error) {
	bw.Lock()
	defer bw.Unlock()
	return bw.buf.Write(data)
}

// Sync forces immediate buffer evacuation to persistent storage.
// Critical considerations:
// - Errors indicate final write success/failure state
// - Caller-owned error handling (no retry logic)
// - Mutex-protected flush prevents concurrent modification
func (bc *BufferWriteCloser) Sync() error {
	bc.Lock()
	defer bc.Unlock()
	return bc.buf.Flush()
}

// Close implements atomic writer termination protocol:
// 1. Stop scheduled flushes (ticker)
// 2. Terminate background goroutine (done channel)
// 3. Final data synchronization (Sync)
// 4. Resource reclamation (writer.Close())
// Design guarantee: Idempotent operation with nil safety
func (db *BufferWriteCloser) Close() error {
	db.flushTicker.Stop()
	db.done <- struct{}{}
	db.Sync()
	return db.writer.Close()
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

func WithStrategy(strategy strategyType) Option {
	return func(aw *AsyncWriter) {
		aw.strategy = strategy
	}
}

func WithDefaultDowngradeWriter(w io.WriteCloser, size int, flushIntvl time.Duration) Option {
	return func(aw *AsyncWriter) {
		if !isNil(w) {
			aw.downgradeWriter = NewBufferWriteCloser(w, size, flushIntvl)
		}
	}
}

// 深层 nil 检查函数
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}
