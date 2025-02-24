package backend

import (
	"sync"
	"time"
)

// 定义 NowFunc 为返回 Unix 时间戳的函数
var nowFunc = func() int64 {
	return time.Now().Unix()
}

type slideBucket struct {
	second       int64 // 时间戳（秒）
	errorCount   int64 // 该秒内的错误次数
	requestCount int64 // 该秒内的总请求次数
}

type SlidingWindow struct {
	mu                   sync.Mutex
	buckets              []*slideBucket
	size                 int     // 窗口大小（秒）
	errorRateThreshold   float64 // 错误率阈值（0.0~1.0）
	enabled              bool    // 是否启用熔断
	currentErrors        int64   // 当前窗口内的总错误数
	currentRequests      int64   // 当前窗口内的总请求数
	startSec             int64   // 窗口起始时间（最旧的秒）
	fueseMinRequestCount int64   // 窗口内总的错误数
}

func NewSlidingWindow(windowSize int, errorRateThreshold float64, fueseMinRequestCount int64) *SlidingWindow {
	// 强制校验参数，非法时禁用
	if windowSize <= 0 || errorRateThreshold <= 0 || errorRateThreshold > 1 || fueseMinRequestCount <= 0 {
		return &SlidingWindow{enabled: false}
	}
	buckets := make([]*slideBucket, windowSize)
	for i := range buckets {
		buckets[i] = &slideBucket{
			second:       nowFunc(),
			errorCount:   0,
			requestCount: 0,
		}
	}

	return &SlidingWindow{
		buckets:              buckets,
		size:                 windowSize,
		errorRateThreshold:   errorRateThreshold,
		enabled:              true,
		fueseMinRequestCount: fueseMinRequestCount,
	}
}

// New method to record request and evaluate fuse condition
func (sw *SlidingWindow) ShouldTrigger(timestamp int64, isError bool) bool {
	if !sw.enabled {
		return false
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()

	// Only slide the window if the timestamp has changed
	if timestamp != sw.startSec {
		sw.slideWindowAt(timestamp)
	}

	// 获取当前秒的桶索引
	idx := timestamp % int64(sw.size)
	bucket := sw.buckets[idx]

	// 处理新秒的桶
	if bucket.second != timestamp {
		if bucket.second >= sw.startSec {
			sw.currentErrors -= bucket.errorCount
			sw.currentRequests -= bucket.requestCount
		}
		bucket.second = timestamp
		bucket.errorCount = 0
		bucket.requestCount = 0
	}

	// 更新请求统计
	bucket.requestCount++
	sw.currentRequests++
	if isError {
		bucket.errorCount++
		sw.currentErrors++
	}

	// 如果请求数低于最小阈值，不考虑熔断
	if sw.currentRequests < sw.fueseMinRequestCount {
		return false
	}

	// 计算错误率并判断是否触发熔断
	if sw.currentRequests == 0 {
		return false
	}

	errorRate := float64(sw.currentErrors) / float64(sw.currentRequests)
	return errorRate >= sw.errorRateThreshold
}

// New method to handle sliding window logic
func (sw *SlidingWindow) slideWindowAt(timestamp int64) {
	newStartSec := timestamp - int64(sw.size-1)
	if newStartSec > sw.startSec {
		delta := newStartSec - sw.startSec
		if delta >= int64(sw.size) {
			// All buckets expired, reset everything
			for i := range sw.buckets {
				sw.buckets[i].second = 0
				sw.buckets[i].errorCount = 0
				sw.buckets[i].requestCount = 0
			}
			sw.currentErrors = 0
			sw.currentRequests = 0
		} else {
			// Remove expired buckets
			for s := sw.startSec; s < newStartSec; s++ {
				idx := s % int64(sw.size)
				bucket := sw.buckets[idx]
				if bucket.second == s {
					sw.currentErrors -= bucket.errorCount
					sw.currentRequests -= bucket.requestCount
					bucket.second = 0
					bucket.errorCount = 0
					bucket.requestCount = 0
				}
			}
		}
		sw.startSec = newStartSec
	}
}

func (sw *SlidingWindow) SetEnabled(e bool) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.enabled = e
}

func (sw *SlidingWindow) IsEnabled() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.enabled
}
