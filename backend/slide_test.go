package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlidingWindow_RecordAt(t *testing.T) {
	// 控制时间
	currentTime := int64(0) // 自定义的开始时间（2021-01-01 00:00:00）
	sw := NewSlidingWindow(5, 0.5, 3)

	// 使用RecordAt插入特定时间的请求
	sw.ShouldTrigger(currentTime, false)   // 请求
	sw.ShouldTrigger(currentTime+1, true)  // 错误请求
	sw.ShouldTrigger(currentTime+2, false) // 请求
	sw.ShouldTrigger(currentTime+3, true)  // 请求
	sw.ShouldTrigger(currentTime+4, false) // 请求

	assert.Equal(t, int64(1), sw.buckets[0].requestCount) // 0秒的请求数
	assert.Equal(t, int64(0), sw.buckets[0].errorCount)   // 0秒的错误数

	assert.Equal(t, int64(1), sw.buckets[1].requestCount) // 1秒的请求数
	assert.Equal(t, int64(1), sw.buckets[1].errorCount)   // 1秒的错误数

	assert.Equal(t, int64(1), sw.buckets[2].requestCount) // 2秒的请求数
	assert.Equal(t, int64(0), sw.buckets[2].errorCount)   // 2秒的错误数

	assert.Equal(t, int64(1), sw.buckets[3].requestCount) // 3秒的请求数
	assert.Equal(t, int64(1), sw.buckets[3].errorCount)   // 3秒的错误数

	assert.Equal(t, int64(1), sw.buckets[4].requestCount) // 3秒的请求数
	assert.Equal(t, int64(0), sw.buckets[4].errorCount)   // 3秒的错误数

}

func TestSlidingWindow_RecordAtOnlyOneRequest(t *testing.T) {
	// 控制时间
	currentTime := int64(0) // 自定义的开始时间（2021-01-01 00:00:00）
	sw := NewSlidingWindow(5, 0.5, 5)

	// 使用RecordAt插入特定时间的请求
	assert.False(t, sw.ShouldTrigger(currentTime, true), "Should break due to high error rate")
	// 0s 错误请求

}

func TestSlidingWindow_ShouldBreak(t *testing.T) {
	// 控制时间
	currentTime := int64(0)           // 自定义的开始时间（2021-01-01 00:00:00）
	sw := NewSlidingWindow(5, 0.5, 5) // 5 second window, error rate threshold 50%, 1 min request per sec

	// 插入请求
	assert.False(t, sw.ShouldTrigger(currentTime, false), "Should not break at time 0")

	assert.False(t, sw.ShouldTrigger(currentTime+1, true), "Should not break at time 1")

	assert.False(t, sw.ShouldTrigger(currentTime+2, true), "Should not break at time 2")

	assert.False(t, sw.ShouldTrigger(currentTime+3, true), "Should not break at time 3")

	assert.True(t, sw.ShouldTrigger(currentTime+4, false), "Should break at time 4 due to high error rate")

	// 插入接下来的5秒的请求，测试总共10秒的时间范围
	assert.True(t, sw.ShouldTrigger(currentTime+5, false), "Should break at time 5 due to high error rate")

	assert.True(t, sw.ShouldTrigger(currentTime+6, true), "Should break at time 6 due to high error rate")

	assert.True(t, sw.ShouldTrigger(currentTime+7, true), "Should break at time 7 due to high error rate")

	assert.False(t, sw.ShouldTrigger(currentTime+8, false), "Should not break at time 8")

	assert.False(t, sw.ShouldTrigger(currentTime+9, false), "Should not break at time 9 after 10 seconds")

}

func TestSlidingWindow_ShouldNotBreak(t *testing.T) {
	// 控制时间
	currentTime := int64(1609459200) // 自定义的开始时间（2021-01-01 00:00:00）

	sw := NewSlidingWindow(5, 0.5, 5)

	// 插入请求
	sw.ShouldTrigger(currentTime, false)
	sw.ShouldTrigger(currentTime+1, false)
	sw.ShouldTrigger(currentTime+2, false)
	sw.ShouldTrigger(currentTime+3, false)
	// 计算是否触发熔断
	assert.False(t, sw.ShouldTrigger(currentTime+4, false), "Should not break due to low error rate")
}
