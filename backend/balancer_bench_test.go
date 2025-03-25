package backend

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 定义加锁版本的 balancer 结构
type lockedBalancer struct {
	mu          sync.Mutex
	nextIndex   int
	roundRobinQ []int
}

func (b *lockedBalancer) next() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	idx := b.roundRobinQ[b.nextIndex]
	b.nextIndex = (b.nextIndex + 1) % len(b.roundRobinQ)
	return idx
}

type atomicBalancer struct {
	nextIndex   uint32 // 改为 uint32
	roundRobinQ []int
}

func (b *atomicBalancer) next() int {
	// 直接原子自增（自动处理溢出）
	newIndex := atomic.AddUint32(&b.nextIndex, 1)
	return b.roundRobinQ[int(newIndex)%len(b.roundRobinQ)]
}

// 初始化测试用的 roundRobinQ
func initBalancerData() []int {
	const size = 1000 // 队列长度
	q := make([]int, size)
	for i := 0; i < size; i++ {
		q[i] = i % 10 // 假设有 10 个实际节点
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(size, func(i, j int) { q[i], q[j] = q[j], q[i] })
	return q
}

// 基准测试：加锁版本
func BenchmarkLockedBalancer(b *testing.B) {
	q := initBalancerData()
	lb := &lockedBalancer{
		roundRobinQ: q,
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = lb.next()
		}
	})
}

// 基准测试：原子操作版本
func BenchmarkAtomicBalancer(b *testing.B) {
	q := initBalancerData()
	ab := &atomicBalancer{
		roundRobinQ: q,
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = ab.next()
		}
	})
}
