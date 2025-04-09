package backend

import (
	"github.com/XiaoMi/Gaea/util/sync2"

	"time"
)

// 组合熔断策略和恢复策略
type FuseStrategy interface {
	// 触发熔断判断，返回是否触发熔断
	Trigger(int64) bool
}

type RecoveryStrategy interface {
	// 检查是否满足恢复条件（硬冷却期优先/渐进式恢复次之）
	AllowRecovery() bool
}

// 硬冷却期恢复策略
type HardCoolDownStrategy struct {
	coolingPeriod int64             // 冷却时长（秒）
	lastFuseTime  sync2.AtomicInt64 // 原子存储上次熔断的时间
}

func NewHardCoolDown(coolingSec int64) *HardCoolDownStrategy {
	return &HardCoolDownStrategy{
		coolingPeriod: coolingSec,
	}
}

// 实现 RecoveryStrategy 接口
func (s *HardCoolDownStrategy) AllowRecovery() bool {
	now := time.Now().Unix()
	line := s.lastFuseTime.Get() + s.coolingPeriod
	return now >= line
}

func (s *HardCoolDownStrategy) UpdateFuseTime(now int64) {
	s.lastFuseTime.Set(now)
}

// 渐进式恢复策略
type GradualRecoveryStrategy struct {
	errorRecoveryCount           sync2.AtomicInt64 // 原子存储最近误恢复次数
	consecutiveSuccessCheckCount sync2.AtomicInt64 // 原子存储需要连续探活成功次数
	lastRecoveryTime             sync2.AtomicInt64 // 原子存储最近一次恢复的时间
	lastFuseTime                 sync2.AtomicInt64 // 原子存储最近一次熔断的时间
}

const (
	maxPenalty             = 120 // 最大惩罚跳过 120
	initErrorRecoveryCount = 3
)

func NewGradualRecovery() *GradualRecoveryStrategy {
	g := &GradualRecoveryStrategy{
		errorRecoveryCount:           sync2.NewAtomicInt64(initErrorRecoveryCount),
		consecutiveSuccessCheckCount: sync2.NewAtomicInt64(0),
		lastRecoveryTime:             sync2.NewAtomicInt64(time.Now().Unix()),
	}
	return g
}

func min(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// 熔断节点(状态从 StatusUp -> StatusDown)时调用
func (s *GradualRecoveryStrategy) UpdateFuseTime(fuseTime int64) {
	// 更新熔断时间
	s.lastFuseTime.Set(fuseTime)
}

func (s *GradualRecoveryStrategy) UpdateCoolDownCount() {
	// 误恢复：在恢复窗口内再次熔断
	n := s.errorRecoveryCount.Add(1)
	newPenalty := (1 + n) * n / 2
	skip := min(newPenalty, maxPenalty)
	s.consecutiveSuccessCheckCount.Set(skip)
}

func (s *GradualRecoveryStrategy) RefreshCoolDownCount() {
	// 误恢复：在恢复窗口内再次熔断
	n := s.errorRecoveryCount.Get()
	newPenalty := (1 + n) * n / 2
	skip := min(newPenalty, maxPenalty)
	s.consecutiveSuccessCheckCount.Set(skip)
}

func (s *GradualRecoveryStrategy) ResetBadRecovery(fuseTime int64) {
	s.errorRecoveryCount.Set(initErrorRecoveryCount)
}

func (g *GradualRecoveryStrategy) IsBadRecovery(fuseTime int64) bool {
	// 获取上次探活成功时间
	lastRecovery := g.lastRecoveryTime.Get()
	return fuseTime-lastRecovery <= PingPeriod*2
}

// 探活节点（状态从 StatusDown -> StatusUp)时调用
func (g *GradualRecoveryStrategy) AllowRecovery() bool {
	if c := g.consecutiveSuccessCheckCount.Get(); c > 0 {
		g.consecutiveSuccessCheckCount.Set(c - 1)
		return false
	}
	return true
}

func (g *GradualRecoveryStrategy) UpdateLastRecoveryTime() {
	g.lastRecoveryTime.Set(time.Now().Unix())
}
