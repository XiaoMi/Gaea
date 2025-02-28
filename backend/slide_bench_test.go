package backend

import (
	"sync"
	"testing"
)

func BenchmarkSlidingWindow_All(b *testing.B) {
	b.Run("HighConcurrency", BenchmarkSlidingWindow_HighConcurrency)
	b.Run("HighTraffic", BenchmarkSlidingWindow_HighTraffic)
	b.Run("HighErrorRate", BenchmarkSlidingWindow_HighErrorRate)
	b.Run("ExpiredRequests", BenchmarkSlidingWindow_ExpiredRequests)
}
func BenchmarkSlidingWindow_HighErrorRate(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := int64(1000)
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, minRequestCount)

	// 每个请求都为错误请求，模拟错误率超过阈值的情况
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Trigger(int64(i))
	}
}

func BenchmarkSlidingWindow_HighConcurrency(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := int64(1000)
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, minRequestCount)

	// 使用并发模拟大量请求
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sw.Trigger(int64(i))
		}(i)
	}
	wg.Wait() // 等待所有 goroutine 完成
}

func BenchmarkSlidingWindow_HighTraffic(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := int64(1000)
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, minRequestCount)

	// 每秒模拟 1000 个请求
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Trigger(int64(i))
	}
}

func BenchmarkSlidingWindow_ExpiredRequests(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := int64(1000)
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, minRequestCount)

	// 每秒请求一个正常的请求
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Trigger(int64(i))
	}
}
