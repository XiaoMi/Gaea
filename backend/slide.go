package backend

import (
	"sync"
)

type SlideBucket struct {
	StartTime  int64 // 桶的起始时间
	ErrorCount int64
}

type SlidingWindow struct {
	mu                sync.Mutex
	windowSizeSec     int64 // 窗口总时长（秒）
	buckets           []*SlideBucket
	enabled           bool
	startSec          int64
	allErrorCount     int64
	fuseMinErrorCount int64
}

func NewSlidingWindow(windowSec int64, fuseMinErrorCount int64) *SlidingWindow {
	// 强制校验参数，非法时禁用
	if windowSec <= 0 || fuseMinErrorCount <= 0 {
		return &SlidingWindow{enabled: false}
	}
	return &SlidingWindow{
		windowSizeSec:     windowSec,
		buckets:           make([]*SlideBucket, windowSec),
		fuseMinErrorCount: fuseMinErrorCount,
		enabled:           true,
		startSec:          0,
		allErrorCount:     0,
	}
}

func (sw *SlidingWindow) Trigger(now int64) bool {
	if !sw.enabled {
		return false
	}

	sw.mu.Lock()
	defer sw.mu.Unlock()

	var (
		newStartSec = now - sw.windowSizeSec + 1
		index       = now % sw.windowSizeSec
	)

	// Only slide the window if the timestamp has changed
	if newStartSec > sw.startSec {
		sw.slide(newStartSec)
	}
	if sw.buckets[index] == nil {
		sw.buckets[index] = &SlideBucket{
			StartTime:  now,
			ErrorCount: 0,
		}
	}

	sw.buckets[index].ErrorCount++
	sw.allErrorCount++
	return sw.allErrorCount >= sw.fuseMinErrorCount
}

func (sw *SlidingWindow) slide(newStartSec int64) {
	delta := newStartSec - sw.startSec
	if delta >= int64(sw.windowSizeSec) {
		// All buckets expired, reset everything
		sw.buckets = make([]*SlideBucket, sw.windowSizeSec)
		sw.allErrorCount = 0
	} else {
		// Remove expired buckets
		for s := sw.startSec; s < newStartSec; s++ {
			idx := s % int64(sw.windowSizeSec)
			if sw.buckets[idx] == nil {
				continue
			}
			sw.allErrorCount -= sw.buckets[idx].ErrorCount
			sw.buckets[idx] = nil
		}
	}
	sw.startSec = newStartSec
}
