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
	"fmt"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

func BenchmarkSyncLoggerWriterWithWaitGroup(b *testing.B) {
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

func BenchmarkAsyncLoggerWriterWithWaitGroup(b *testing.B) {
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

func BenchmarkSyncLoggerWriter(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	encoder := &ZapEncoder{}
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(f),
		zap.LevelEnablerFunc(func(zapcore.Level) bool { return true }),
	)
	logger := zap.New(core)
	defer logger.Sync()

	// 预先生成测试数据避免内存分配干扰
	logMsg := "ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(logMsg)
		}
	})

	b.ReportAllocs() // 必须放在测试逻辑之后
}

func BenchmarkAsyncLoggerWriter(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	// 创建带缓冲的异步写入器
	asyncWriter := NewAsyncWriter(f)
	defer asyncWriter.Close()

	encoder := &ZapEncoder{}
	core := zapcore.NewCore(
		encoder,
		asyncWriter,
		zap.LevelEnablerFunc(func(zapcore.Level) bool { return true }),
	)
	logger := zap.New(core)
	defer logger.Sync()

	// 设置并行度参数
	b.SetParallelism(16) // 设置每个CPU核心的并发因子
	logMsg := "ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(logMsg)
		}
	})

	// 报告关键指标
	b.ReportMetric(float64(asyncWriter.Dropped()), "drops/op")
	b.ReportAllocs()
}

func BenchmarkAsyncLoggerWriterDiscard(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	// 创建带缓冲的异步写入器
	asyncWriter := NewAsyncWriter(
		f,
		WithStrategy(strategyDiscard),
	)
	defer asyncWriter.Close()

	encoder := &ZapEncoder{}
	core := zapcore.NewCore(
		encoder,
		asyncWriter,
		zap.LevelEnablerFunc(func(zapcore.Level) bool { return true }),
	)
	logger := zap.New(core)
	defer logger.Sync()

	// 设置并行度参数
	b.SetParallelism(16) // 设置每个CPU核心的并发因子
	logMsg := "ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(logMsg)
		}
	})

	// 报告关键指标
	drops := asyncWriter.Dropped()
	b.ReportMetric(float64(drops)/float64(b.N), "drops/op")
	b.ReportAllocs()
}

func BenchmarkAsyncLoggerWriterDiscardDowngrade(b *testing.B) {
	f1, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	defer f1.Close()

	f2, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	defer f2.Close()

	// 创建带缓冲的异步写入器
	asyncWriter := NewAsyncWriter(
		f1,
		WithDefaultDowngradeWriter(f2, defaultBufferSize, defaultBufferFlushIntvl),
		WithStrategy(strategyDiscardWithDowngrade),
	)
	defer asyncWriter.Close()

	encoder := &ZapEncoder{}
	core := zapcore.NewCore(
		encoder,
		asyncWriter,
		zap.LevelEnablerFunc(func(zapcore.Level) bool { return true }),
	)
	logger := zap.New(core)
	defer logger.Sync()

	// 设置并行度参数
	b.SetParallelism(16) // 设置每个CPU核心的并发因子
	logMsg := "ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info(logMsg)
		}
	})

	// 报告关键指标
	drops := asyncWriter.Dropped()
	b.ReportMetric(float64(drops)/float64(b.N), "drops/op")
	b.ReportAllocs()
}

func BenchmarkDiscardAsyncLoggerWriter_QPS(b *testing.B) {
	// 100 parallelism
	b.Run("QPS_1k Parallelism_0.1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 1_000, 100)
	})

	b.Run("QPS_5k Parallelism_0.1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 5_000, 100)
	})

	b.Run("QPS_10k Parallelism_0.1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 10_000, 100)
	})

	b.Run("QPS_50k Parallelism_0.1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 50_000, 100)
	})

	// 1000 parallelism
	b.Run("QPS_1k Parallelism_1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 1_000, 1000)
	})

	b.Run("QPS_5k Parallelism_1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 5_000, 1000)
	})

	b.Run("QPS_10k Parallelism_1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 10_000, 1000)
	})

	b.Run("QPS_50k Parallelism_1k", func(b *testing.B) {
		benchmarkQPSLevel(b, 50_000, 1000)
	})

	// 10000 parallelism
	b.Run("QPS_1k Parallelism_10k", func(b *testing.B) {
		benchmarkQPSLevel(b, 1_000, 10_000)
	})

	b.Run("QPS_5k Parallelism_10k", func(b *testing.B) {
		benchmarkQPSLevel(b, 5_000, 10_000)
	})

	b.Run("QPS_10k Parallelism_10k", func(b *testing.B) {
		benchmarkQPSLevel(b, 10_000, 10_000)
	})

	b.Run("QPS_50k Parallelism_10k", func(b *testing.B) {
		benchmarkQPSLevel(b, 50_000, 10_000)
	})

	// 100000 parallelism
	b.Run("QPS_5k Parallelism_100k", func(b *testing.B) {
		benchmarkQPSLevel(b, 5_000, 100_000)
	})

	b.Run("QPS_10k Parallelism_100k", func(b *testing.B) {
		benchmarkQPSLevel(b, 10_000, 100_000)
	})

	b.Run("QPS_50k Parallelism_100k", func(b *testing.B) {
		benchmarkQPSLevel(b, 50_000, 100_000)
	})

	b.Run("QPS_100k Parallelism_100k", func(b *testing.B) {
		benchmarkQPSLevel(b, 100_000, 100_000)
	})
}

// 基准测试核心逻辑（参数显式传递）
func benchmarkQPSLevel(b *testing.B, qps int, parallelism int) {
	// 1. 资源初始化
	f, err := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		b.Fatalf("打开文件失败: %v", err)
	}

	// 2. 限流器配置（带QPS验证）
	limiter := rate.NewLimiter(rate.Limit(qps), 1)

	// 3. 异步写入器（带容量监控）
	asyncWriter := NewAsyncWriter(
		f,
		WithStrategy(strategyDiscard),
	)
	defer asyncWriter.Close()

	// 4. 日志核心（带压力检测）
	encoder := &ZapEncoder{}
	core := zapcore.NewCore(
		encoder,
		asyncWriter,
		zap.LevelEnablerFunc(func(zapcore.Level) bool { return true }),
	)
	logger := zap.New(core)
	defer logger.Sync()

	// 5. 动态日志内容（防止编译器优化）
	logMsg := fmt.Sprintf("load_test_qps=%d_%d", qps, time.Now().UnixNano())

	// 6. 执行参数配置
	b.SetParallelism(parallelism)
	b.ResetTimer()

	// 7. 压力测试执行
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if limiter.Allow() {
				logger.Info(logMsg)
			}
		}
	})

	// 8. 结果验证
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	b.ReportMetric(float64(memStats.Mallocs)/float64(b.N), "mallocs/op")

	drops := asyncWriter.Dropped()
	b.ReportMetric(float64(drops)/float64(b.N), "drops/op")
	b.ReportAllocs()
}
