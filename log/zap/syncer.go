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
