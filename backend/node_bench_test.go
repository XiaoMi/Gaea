package backend

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkNodeSetStatus(b *testing.B) {
	// 测试单线程下的锁和原子操作
	b.Run("SingleThread", func(b *testing.B) {
		b.Run("Lock_SetStatusUp", func(b *testing.B) {
			var l LockVersion
			for i := 0; i < b.N; i++ {
				l.SetStatusUp()
			}
		})
		b.Run("Atomic_SetStatusUp", func(b *testing.B) {
			var a AtomicVersion
			for i := 0; i < b.N; i++ {
				a.SetStatusUp()
			}
		})
	})

	// 测试高并发下的锁和原子操作（8线程）
	b.Run("HighConcurrency", func(b *testing.B) {
		b.Run("Lock_SetStatusUp", func(b *testing.B) {
			var l LockVersion
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					l.SetStatusUp()
				}
			})
		})
		b.Run("Atomic_SetStatusUp", func(b *testing.B) {
			var a AtomicVersion
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					a.SetStatusUp()
				}
			})
		})
	})
}

// LockVersion 使用锁保护状态
type LockVersion struct {
	mu     sync.RWMutex
	status int32
}

func (l *LockVersion) SetStatusUp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.status = 1
}

func (l *LockVersion) SetStatusDown() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.status = 0
}

// AtomicVersion 使用原子操作保护状态
type AtomicVersion struct {
	status int32
}

func (a *AtomicVersion) SetStatusUp() {
	atomic.StoreInt32(&a.status, 1)
}

func (a *AtomicVersion) SetStatusDown() {
	atomic.StoreInt32(&a.status, 0)
}

// 基准测试：锁版本的 SetStatusUp
func BenchmarkLockVersion_SetStatusUp(b *testing.B) {
	var l LockVersion
	for i := 0; i < b.N; i++ {
		l.SetStatusUp()
	}
}

// 基准测试：原子版本的 SetStatusUp
func BenchmarkAtomicVersion_SetStatusUp(b *testing.B) {
	var a AtomicVersion
	for i := 0; i < b.N; i++ {
		a.SetStatusUp()
	}
}

// 基准测试：锁版本的 SetStatusDown
func BenchmarkLockVersion_SetStatusDown(b *testing.B) {
	var l LockVersion
	for i := 0; i < b.N; i++ {
		l.SetStatusDown()
	}
}

// 基准测试：原子版本的 SetStatusDown
func BenchmarkAtomicVersion_SetStatusDown(b *testing.B) {
	var a AtomicVersion
	for i := 0; i < b.N; i++ {
		a.SetStatusDown()
	}
}

// 多线程并发测试：锁版本的 SetStatusUp（8 线程）
func BenchmarkLockVersion_SetStatusUp_Parallel(b *testing.B) {
	var l LockVersion
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.SetStatusUp()
		}
	})
}

// 多线程并发测试：原子版本的 SetStatusUp（8 线程）
func BenchmarkAtomicVersion_SetStatusUp_Parallel(b *testing.B) {
	var a AtomicVersion
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			a.SetStatusUp()
		}
	})
}
