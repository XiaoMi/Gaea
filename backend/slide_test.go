package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlidingWindow_Trigger_Simple(t *testing.T) {
	// 控制时间
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 3)

	// 使用Trigger插入特定时间的请求
	sw.Trigger(currentTime)     // 请求
	sw.Trigger(currentTime + 1) // 错误请求
	sw.Trigger(currentTime + 2) // 请求
	sw.Trigger(currentTime + 3) // 请求
	sw.Trigger(currentTime + 4) // 请求

	assert.Equal(t, int64(1), sw.buckets[0].ErrorCount) // 0秒的错误数

	assert.Equal(t, int64(1), sw.buckets[1].ErrorCount) // 1秒的错误数

	assert.Equal(t, int64(1), sw.buckets[2].ErrorCount) // 2秒的错误数

	assert.Equal(t, int64(1), sw.buckets[3].ErrorCount) // 3秒的错误数

	assert.Equal(t, int64(1), sw.buckets[4].ErrorCount) // 3秒的错误数

}

func TestSlidingWindowTrigger(t *testing.T) {
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 6)

	sw.Trigger(currentTime) // 请求
	sw.Trigger(currentTime) // 请求
	sw.Trigger(currentTime) // 请求

	sw.Trigger(currentTime + 1) // 请求
	sw.Trigger(currentTime + 1) // 请求
	assert.True(t, sw.Trigger(currentTime+1), "Should break due to high error rate")

}

func TestSlidingWindowTriggerTimeSkip(t *testing.T) {
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 10)

	for i := 0; i < 5; i++ {
		sw.Trigger(currentTime) // 请求
	}

	assert.False(t, sw.Trigger(currentTime+5), "Should not break due to high error rate")

}
func TestSlidingWindow_TriggerOnlyOneRequest(t *testing.T) {
	// 控制时间
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 4)

	// 使用RecordAt插入特定时间的请求
	assert.False(t, sw.Trigger(currentTime), "Should break due to high error rate")

}

func TestSlidingWindow_Trigger_NumCount(t *testing.T) {
	// 控制时间
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 5)

	// 使用RecordAt插入特定时间的请求
	sw.Trigger(currentTime)
	assert.Equal(t, sw.allErrorCount, int64(1))
	sw.Trigger(currentTime + 1)
	assert.Equal(t, sw.allErrorCount, int64(2))
	sw.Trigger(currentTime + 2)
	assert.Equal(t, sw.allErrorCount, int64(3))
	sw.Trigger(currentTime + 3)
	assert.Equal(t, sw.allErrorCount, int64(4))
	sw.Trigger(currentTime + 4)
	assert.Equal(t, sw.allErrorCount, int64(5))

	sw.Trigger(currentTime + 5)
	assert.Equal(t, sw.allErrorCount, int64(5))
	sw.Trigger(currentTime + 6)
	assert.Equal(t, sw.allErrorCount, int64(5))
	sw.Trigger(currentTime + 7)
	assert.Equal(t, sw.allErrorCount, int64(5))
	sw.Trigger(currentTime + 8)
	assert.Equal(t, sw.allErrorCount, int64(5))
	sw.Trigger(currentTime + 9)
	assert.Equal(t, sw.allErrorCount, int64(5))

}

func TestSlidingWindow_AllBucketExpired(t *testing.T) {
	// 控制时间
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 5)

	sw.Trigger(currentTime + 5)
	assert.Equal(t, sw.allErrorCount, int64(1))
	sw.Trigger(currentTime + 6)
	assert.Equal(t, sw.allErrorCount, int64(2))
	sw.Trigger(currentTime + 7)
	assert.Equal(t, sw.allErrorCount, int64(3))
	sw.Trigger(currentTime + 8)
	assert.Equal(t, sw.allErrorCount, int64(4))
	sw.Trigger(currentTime + 9)
	assert.Equal(t, sw.allErrorCount, int64(5))

	// [5 6 7 8 9] [10 11 12 13 14]
	sw.Trigger(currentTime + 14)
	assert.Equal(t, sw.allErrorCount, int64(1))
}

func TestSlidingWindow_SomeBucketExpired(t *testing.T) {
	// 控制时间
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 5)

	sw.Trigger(currentTime + 5)
	assert.Equal(t, sw.allErrorCount, int64(1))
	sw.Trigger(currentTime + 6)
	assert.Equal(t, sw.allErrorCount, int64(2))
	sw.Trigger(currentTime + 7)
	assert.Equal(t, sw.allErrorCount, int64(3))
	sw.Trigger(currentTime + 8)
	assert.Equal(t, sw.allErrorCount, int64(4))
	sw.Trigger(currentTime + 9)
	assert.Equal(t, sw.allErrorCount, int64(5))

	// [5 6 7 8 9] [10 11 12 13 14]
	sw.Trigger(currentTime + 12)
	assert.Equal(t, sw.allErrorCount, int64(3))
}

func TestSlidingWindow_ShouldBreak(t *testing.T) {
	// 控制时间
	currentTime := int64(0)      // 自定义的开始时间（2021-01-01 00:00:00）
	sw := NewSlidingWindow(5, 5) // 5 second window, error rate threshold 50%, 1 min request per sec

	// 插入请求
	assert.False(t, sw.Trigger(currentTime), "Should not break at time 0")

	assert.False(t, sw.Trigger(currentTime+1), "Should not break at time 1")

	assert.False(t, sw.Trigger(currentTime+2), "Should not break at time 2")

	assert.False(t, sw.Trigger(currentTime+3), "Should not break at time 3")

	assert.True(t, sw.Trigger(currentTime+4), "Should break at time 4 due to high error rate")

	// 插入接下来的5秒的请求，测试总共10秒的时间范围
	assert.True(t, sw.Trigger(currentTime+5), "Should break at time 5 due to high error rate")

	assert.True(t, sw.Trigger(currentTime+6), "Should break at time 6 due to high error rate")

	assert.True(t, sw.Trigger(currentTime+7), "Should break at time 7 due to high error rate")

	assert.True(t, sw.Trigger(currentTime+8), "Should not break at time 8")

	assert.True(t, sw.Trigger(currentTime+9), "Should not break at time 9 after 10 seconds")

}

func TestSlidingWindow_ShouldNotBreak(t *testing.T) {
	// 控制时间
	currentTime := int64(1609459200) // 自定义的开始时间（2021-01-01 00:00:00）

	sw := NewSlidingWindow(5, 5)

	// 插入请求
	sw.Trigger(currentTime)
	sw.Trigger(currentTime + 1)
	sw.Trigger(currentTime + 2)
	// 计算是否触发熔断
	assert.False(t, sw.Trigger(currentTime+4), "Should not break due to low error rate")
}

func TestSlidingWindow_SomeBucketExpiredAndRequestMore(t *testing.T) {
	// 控制时间
	currentTime := int64(0)
	sw := NewSlidingWindow(5, 5)

	for i := 0; i < 10000; i++ {
		sw.Trigger(currentTime + 6)
	}

	for i := 0; i < 10000; i++ {
		sw.Trigger(currentTime + 7)
	}

	assert.Equal(t, sw.allErrorCount, int64(20000))

	// [5 6 7 8 9] [10 11 12 13 14]
	sw.Trigger(currentTime + 11)
	assert.Equal(t, sw.allErrorCount, int64(10001))
}
