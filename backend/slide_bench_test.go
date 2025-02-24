package backend

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkSlidingWindow_All(b *testing.B) {
	b.Run("HighConcurrency", BenchmarkSlidingWindow_HighConcurrency)
	b.Run("HighTraffic", BenchmarkSlidingWindow_HighTraffic)
	b.Run("HighErrorRate", BenchmarkSlidingWindow_HighErrorRate)
	b.Run("ExpiredRequests", BenchmarkSlidingWindow_ExpiredRequests)
}
func BenchmarkSlidingWindow_HighErrorRate(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := 5
	errorRateThreshold := 0.6
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, errorRateThreshold, minRequestCount)

	// 每个请求都为错误请求，模拟错误率超过阈值的情况
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timestamp := time.Now().Unix() + int64(i)
		sw.ShouldTrigger(timestamp, true) // 每个请求都是错误请求
	}
}

func BenchmarkSlidingWindow_HighConcurrency(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := 5
	errorRateThreshold := 0.6
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, errorRateThreshold, minRequestCount)

	// 使用并发模拟大量请求
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			timestamp := time.Now().Unix()
			isError := i%2 == 0 // 偶数为错误请求，奇数为正常请求
			sw.ShouldTrigger(timestamp, isError)
		}(i)
	}
	wg.Wait() // 等待所有 goroutine 完成
}

func BenchmarkSlidingWindow_HighTraffic(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := 5
	errorRateThreshold := 0.6
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, errorRateThreshold, minRequestCount)

	// 每秒模拟 1000 个请求
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timestamp := time.Now().Unix() + int64(i)
		isError := i%2 == 0 // 偶数请求为错误请求
		sw.ShouldTrigger(timestamp, isError)
	}
}

func BenchmarkSlidingWindow_ExpiredRequests(b *testing.B) {
	// 设置窗口大小为 5 秒，错误率阈值为 60%，最小请求数阈值为 100
	windowSize := 5
	errorRateThreshold := 0.6
	minRequestCount := int64(100)

	sw := NewSlidingWindow(windowSize, errorRateThreshold, minRequestCount)

	// 每秒请求一个正常的请求
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timestamp := time.Now().Unix() + int64(i)
		isError := i%2 == 0 // 偶数为错误请求
		sw.ShouldTrigger(timestamp, isError)
	}
}
