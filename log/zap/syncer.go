package zap

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"
)

const writerBuffSize = 1024 * 1024

type LogAsyncWriter struct {
	sync.Mutex
	writer io.WriteCloser
	buf    *bufio.Writer
	timer  *time.Ticker
	quit   context.CancelFunc
	ctx    context.Context
}

func NewAsyncWriter(w io.WriteCloser) *LogAsyncWriter {
	buf := bufio.NewWriterSize(w, writerBuffSize)
	syncer := &LogAsyncWriter{
		writer: w,
		buf:    buf,
	}
	go syncer.intervalFlush()
	return syncer
}

func (l *LogAsyncWriter) Write(data []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.buf.Write(data)
}

func (l *LogAsyncWriter) Sync() error {
	l.Lock()
	defer l.Unlock()
	return l.buf.Flush()
}

func (l *LogAsyncWriter) intervalFlush() {
	l.ctx, l.quit = context.WithCancel(context.Background())
	l.timer = time.NewTicker(time.Millisecond * 100)
	for {
		select {
		case <-l.timer.C:
			l.Sync()
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *LogAsyncWriter) Close() error {
	defer l.writer.Close()
	l.timer.Stop()
	l.quit()
	return l.Sync()
}
