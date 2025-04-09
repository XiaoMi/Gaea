package backend

import (
	"sync"
	"testing"
	"time"
)

func TestNodeInfo_AtomicMethods(t *testing.T) {
	node := &NodeInfo{
		Status: StatusUp,
	}

	t.Run("IsStatusUp/Down", func(t *testing.T) {
		// 初始状态应为默认值（假设默认是 UP）
		if !node.IsStatusUp() || node.IsStatusDown() {
			t.Fatal("Initial status mismatch")
		}

		// 设置为 Down 后验证
		node.SetStatusDown()
		if node.IsStatusUp() || !node.IsStatusDown() {
			t.Fatal("After SetStatusDown: status mismatch")
		}

		// 重新设置为 Up 后验证
		node.SetStatusUp()
		if !node.IsStatusUp() || node.IsStatusDown() {
			t.Fatal("After SetStatusUp: status mismatch")
		}
	})

	t.Run("ConcurrentMustSet", func(t *testing.T) {
		const goroutines = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)

		// 并发交替设置状态
		for i := 0; i < goroutines; i++ {
			go func(i int) {
				defer wg.Done()
				if i%2 == 0 {
					node.SetStatusUp()
				} else {
					node.SetStatusDown()
				}
			}(i)
		}
		wg.Wait()

		// 最终状态应为最后一次操作的结果（非确定性，但应无竞争）
		finalStatusUp := node.IsStatusUp()
		finalStatusDown := node.IsStatusDown()
		if finalStatusUp && finalStatusDown {
			t.Error("Status cannot be both UP and Down")
		}
		if !finalStatusUp && !finalStatusDown {
			t.Error("Status must be either UP or Down")
		}
	})

}

func TestNodeInfo_AtomicSet(t *testing.T) {
	node := &NodeInfo{
		Status: StatusUp,
	}

	t.Run("ConcurrentSet", func(t *testing.T) {
		tag1 := node.SetStatusDown()
		tag2 := node.SetStatusDown()
		tag3 := node.SetStatusDown()
		tag4 := node.SetStatusDown()
		tag5 := node.SetStatusDown()
		if !tag1 {
			t.Error("Expected Status Change to true")
		}
		if tag2 {
			t.Error("Expected Status Change to false")
		}
		if tag3 {
			t.Error("Expected Status Change to false")
		}
		if tag4 {
			t.Error("Expected Status Change to false")
		}
		if tag5 {
			t.Error("Expected Status Change to false")
		}

	})
}

func TestHardCoolDownStrategy(t *testing.T) {
	t.Run("AllowRecovery within cooling period", func(t *testing.T) {
		h := NewHardCoolDown(10)
		now := time.Now().Unix()
		h.UpdateFuseTime(now - 5) // 5秒前熔断，冷却期还剩5秒
		if h.AllowRecovery() {
			t.Error("Expected recovery to be disallowed within cooling period")
		}
	})

	t.Run("AllowRecovery after cooling period", func(t *testing.T) {
		h := NewHardCoolDown(10)
		h.UpdateFuseTime(time.Now().Unix() - 15) // 15秒前熔断，已过冷却期
		if !h.AllowRecovery() {
			t.Error("Expected recovery to be allowed after cooling period")
		}
	})

	t.Run("UpdateFuseTime correctly sets last fuse time", func(t *testing.T) {
		h := NewHardCoolDown(10)
		now := time.Now().Unix()
		h.UpdateFuseTime(now)
		if h.lastFuseTime.Get() != now {
			t.Errorf("Expected last fuse time %d, got %d", now, h.lastFuseTime.Get())
		}
	})
}

var TestPingPeriod int64 = 4

func gradualRecoveryFunc(g *GradualRecoveryStrategy, now int64) {
	g.UpdateFuseTime(now) // 记录一次误恢复
	isBad := g.IsBadRecovery(now)
	if isBad {
		g.UpdateCoolDownCount()
	} else {
		g.ResetBadRecovery(now)
	}
}

func TestGradualRecoveryStrategy(t *testing.T) {
	now := time.Now().Unix()
	t.Run("AllowRecoval when disabled until passed", func(t *testing.T) {
		g := NewGradualRecovery()
		g.consecutiveSuccessCheckCount.Set(0)
		if !g.AllowRecovery() {
			t.Error("Expected recovery to be allowed")
		} else {
			g.UpdateLastRecoveryTime()
		}
		if g.lastRecoveryTime.Get() < now {
			t.Error("Expected last recovery time to be updated")
		}
	})

	t.Run("Disallow recovery when disabled until not passed", func(t *testing.T) {
		g := NewGradualRecovery()
		g.consecutiveSuccessCheckCount.Set(6)
		if g.AllowRecovery() {
			t.Error("Expected recovery to be disallowed")
		}
	})

	t.Run("Bad recovery increases penalty", func(t *testing.T) {
		g := NewGradualRecovery()
		// 模拟在恢复窗口内再次熔断
		g.lastRecoveryTime.Set(now - 5) // 在8秒窗口内,5s前恢复
		gradualRecoveryFunc(g, now)

		if count := g.errorRecoveryCount.Get(); count != initErrorRecoveryCount+1 {
			t.Errorf("Expected error count 1, got %d", count)
		}
		//((1+n)*n/2)
		expectedPenalty := int64((1 + initErrorRecoveryCount + 1) * (initErrorRecoveryCount + 1) / 2)
		if penalty := g.consecutiveSuccessCheckCount.Get(); penalty != expectedPenalty {
			t.Errorf("Expected penalty %d, got %d", expectedPenalty, penalty)
		}
	})

	t.Run("Multiple bad recoveries compound penalty", func(t *testing.T) {
		g := NewGradualRecovery()
		// 第一次误恢复
		// now-5，now-4，now-3，now-2，now-1，now
		// 恢复                              熔断
		g.lastRecoveryTime.Set(now - 5)
		gradualRecoveryFunc(g, now)
		// 第二次误恢复
		// now-5，now-4，now-3，now-2，now-1，now，now+1，now+2，now+3，now+4，now+5，now+6，now+7，now+8，now+9，now+10，now+11，now+12，
		// 恢复                              熔断                                                       恢复                   熔断
		g.lastRecoveryTime.Set(now + TestPingPeriod*2 + 1) // 模拟恢复后立即再次熔断
		gradualRecoveryFunc(g, now)                        // 记录一次误恢复

		if count := g.errorRecoveryCount.Get(); count != initErrorRecoveryCount+2 {
			t.Errorf("Expected error count 2, got %d", count)
		}
		expectedPenalty := int64((1 + initErrorRecoveryCount + 2) * (initErrorRecoveryCount + 2) / 2)
		if penalty := g.consecutiveSuccessCheckCount.Get(); penalty != expectedPenalty {
			t.Errorf("Expected penalty %d, got %d", expectedPenalty, penalty)
		}
	})

	t.Run("Penalty capped at maximum", func(t *testing.T) {
		g := NewGradualRecovery()
		g.errorRecoveryCount.Set(100) // 足够大的错误次数
		g.lastRecoveryTime.Set(now)
		gradualRecoveryFunc(g, now)

		expectedPenalty := int64(maxPenalty)
		if penalty := g.consecutiveSuccessCheckCount.Get(); penalty != expectedPenalty {
			t.Errorf("Expected penalty capped at %d, got %d", expectedPenalty, penalty)
		}
	})

	t.Run("Recovery outside window resets counter", func(t *testing.T) {
		g := NewGradualRecovery()
		// 上一次恢复在25秒前（超过2*TestPingPeriod）
		g.lastRecoveryTime.Set(now - 25)
		gradualRecoveryFunc(g, now) // 记录一次误恢复

		if count := g.errorRecoveryCount.Get(); count != initErrorRecoveryCount {
			t.Errorf("Expected error count reset to 0, got %d", count)
		}
	})
}
