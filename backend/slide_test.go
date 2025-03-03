package backend

import (
	"testing"
	"time"

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

func TestNodeRecovery(t *testing.T) {
	node := &NodeInfo{
		Status: StatusUp,
		FuseWindow: &SlidingWindow{
			buckets:       nil,
			windowSizeSec: 3,
			enabled:       true,
		},
	}
	// 重复设置 up
	node.SetStatus(StatusUp)
	assert.Equal(t, StatusUp, node.Status)
	assert.Equal(t, int(node.RecoveryTime), 0)
	assert.Equal(t, int(node.ErrorRecoveryCount), 0)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)

	// up 转为 down
	node.SetStatus(StatusDown)
	assert.Equal(t, StatusDown, node.Status)
	assert.Equal(t, int(node.RecoveryTime), 0)
	assert.Equal(t, int(node.ErrorRecoveryCount), 0)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)

	// 重复设置为 down
	node.SetStatus(StatusDown)
	assert.Equal(t, StatusDown, node.Status)
	assert.Equal(t, int(node.RecoveryTime), 0)
	assert.Equal(t, int(node.ErrorRecoveryCount), 0)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)

	// down 转为 up
	node.SetStatus(StatusUp)
	assert.Equal(t, StatusUp, node.Status)
	assert.Equal(t, node.RecoveryTime > 0, true)
	assert.Equal(t, int(node.ErrorRecoveryCount), 0)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)
}

// 测试误恢复的情况
func TestNodeRecoveryError(t *testing.T) {
	node := &NodeInfo{
		Status: StatusDown,
		FuseWindow: &SlidingWindow{
			buckets:       nil,
			windowSizeSec: 3,
			enabled:       true,
		},
	}
	// down 转为 up, 标记为恢复
	node.SetStatus(StatusUp)
	assert.Equal(t, StatusUp, node.Status)
	assert.Equal(t, int(node.ErrorRecoveryCount), 0)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)
	assert.Equal(t, node.RecoveryTime > 0, true)

	// 又被置为 down
	node.SetStatus(StatusDown)
	assert.Equal(t, StatusDown, node.Status)
	assert.Equal(t, int(node.ErrorRecoveryCount), 1)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 1)
	assert.Equal(t, node.RecoveryTime > 0, true)

	// 重复置为 down
	assert.Equal(t, StatusDown, node.Status)
	assert.Equal(t, int(node.ErrorRecoveryCount), 1)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 1)
	assert.Equal(t, node.RecoveryTime > 0, true)

	// 测试是否被跳过
	node.SetStatus(StatusUp)
	assert.Equal(t, StatusDown, node.Status)
	assert.Equal(t, int(node.ErrorRecoveryCount), 1)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)
	assert.Equal(t, node.RecoveryTime > 0, true)
	node.SetStatus(StatusUp)
	assert.Equal(t, StatusUp, node.Status)
	assert.Equal(t, int(node.ErrorRecoveryCount), 1)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)
	assert.Equal(t, node.RecoveryTime > 0, true)

	// 测试计数器被重置
	node.RecoveryTime = time.Now().Unix() - int64(node.FuseWindow.windowSizeSec*2) - 1
	node.SetStatus(StatusDown)
	assert.Equal(t, StatusDown, node.Status)
	assert.Equal(t, int(node.ErrorRecoveryCount), 0)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)
}

// 测试误恢复次数越多，越难恢复
func TestNodeRecoveryMulti(t *testing.T) {
	node := &NodeInfo{
		Status: StatusDown,
		FuseWindow: &SlidingWindow{
			buckets:       nil,
			windowSizeSec: 3,
			enabled:       true,
		},
	}
	for i := 0; i <= 10; i++ {
		skipCount := (1 + i) * i / 2
		assert.Equal(t, int(node.ErrorRecoveryCount), i)
		assert.Equal(t, int(node.RecoveryShouldSkipCount), skipCount)
		for k := 0; k < skipCount; k++ {
			assert.Equal(t, StatusDown, node.Status)
			assert.Equal(t, int(node.ErrorRecoveryCount), i)
			assert.Equal(t, int(node.RecoveryShouldSkipCount), skipCount-k)
			node.SetStatus(StatusUp)
			assert.Equal(t, StatusDown, node.Status)
			assert.Equal(t, int(node.ErrorRecoveryCount), i)
			assert.Equal(t, int(node.RecoveryShouldSkipCount), skipCount-k-1)
		}
		node.SetStatus(StatusUp)
		assert.Equal(t, StatusUp, node.Status)
		assert.Equal(t, int(node.ErrorRecoveryCount), i)
		assert.Equal(t, int(node.RecoveryShouldSkipCount), 0)
		assert.Equal(t, node.RecoveryTime > 0, true)
		node.SetStatus(StatusDown)
		assert.Equal(t, StatusDown, node.Status)
		assert.Equal(t, int(node.ErrorRecoveryCount), i+1)
		assert.Equal(t, int(node.RecoveryShouldSkipCount), skipCount+i+1)
	}
}

// 测试最大恢复计数
func TestNodeCheckBadRecovery(t *testing.T) {
	node := &NodeInfo{
		Status: StatusUp,
		FuseWindow: &SlidingWindow{
			buckets:       nil,
			windowSizeSec: 3,
			enabled:       true,
		},
		RecoveryTime: time.Now().Unix(),
	}
	recoveryCount := 10
	for i := 0; i < recoveryCount; i++ {
		node.checkBadRecovery()
	}
	skipCount := (1 + recoveryCount) * recoveryCount / 2
	assert.Equal(t, int(node.ErrorRecoveryCount), recoveryCount)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), skipCount)

	node.ErrorRecoveryCount = 0
	node.RecoveryShouldSkipCount = 0
	recoveryCount = 20
	for i := 0; i < recoveryCount; i++ {
		node.checkBadRecovery()
	}
	skipCount = 150
	assert.Equal(t, int(node.ErrorRecoveryCount), recoveryCount)
	assert.Equal(t, int(node.RecoveryShouldSkipCount), skipCount)
}
