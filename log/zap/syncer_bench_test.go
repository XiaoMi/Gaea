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
	"io"
	"testing"
)

// 实现无操作 Close 的 WriteCloser
type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil // 空操作关闭
}

// 创建兼容 WriteCloser 的 Discard
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

func BenchmarkAsyncWriter(b *testing.B) {
	// 使用 Discard 避免实际 I/O 影响测试
	writer := NewAsyncWriter(NopWriteCloser(io.Discard))
	defer writer.Close()

	// 测试不同写入粒度
	testCases := []struct {
		name string
		data []byte
	}{
		{"16B", make([]byte, 16)},
		{"1KB", make([]byte, 1024)},
		{"16KB", make([]byte, 1024*16)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.SetParallelism(16) // 模拟高并发
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					writer.Write(tc.data)
				}
			})
		})
	}
}
func BenchmarkSyncWrite(b *testing.B) {
	// 对比组：直接写无缓冲
	writer := io.Discard
	data := make([]byte, 256)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer.Write(data)
		}
	})
}
