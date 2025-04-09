package backend

import (
	"fmt"
	"sync/atomic"
	"time"
)

// 更新 `NodeInfo` 结构
type NodeInfo struct {
	// sync.RWMutex          // 保护 `Status`
	Address    string         // 节点地址
	Datacenter string         // 节点所属的数据中心
	Weight     int            // 该节点的负载均衡权重
	ConnPool   ConnectionPool // 该节点的连接池
	Status     StatusCode     // 该节点状态`status` 只能通过`GetStatus` 和 `SetStatus` 访问

	FuseStrategy     FuseStrategy     // 该节点使用的熔断策略
	RecoveryStrategy RecoveryStrategy // 该节点使用的恢复策略
}

// GetStatus 原子读取状态
func (n *NodeInfo) GetStatus() StatusCode {
	return n.Status
}

// IsStatusUp 原子读取判断是否为 UP
func (n *NodeInfo) IsStatusUp() bool {
	return n.Status == StatusUp
}

// IsStatusDown 原子读取判断是否为 Down
func (n *NodeInfo) IsStatusDown() bool {
	return n.Status == StatusDown
}

// SetStatusUp 原子设置为 UP，返回 true 表示状态发生变更
func (n *NodeInfo) SetStatusUp() bool {
	old := atomic.SwapUint32((*uint32)(&n.Status), uint32(StatusUp))
	return old != uint32(StatusUp)
}

// SetStatusDown 原子设置为 Down，返回 true 表示状态发生变更
func (n *NodeInfo) SetStatusDown() bool {
	old := atomic.SwapUint32((*uint32)(&n.Status), uint32(StatusDown))
	return old != uint32(StatusDown)
}

// NodeInfo 封装获取 PooledConnect 的方法
func (n *NodeInfo) GetPooledConnectWithHealthCheck(name string, healthCheckSql string) (PooledConnect, error) {
	pc, err := checkInstanceStatus(name, n.ConnPool, healthCheckSql)
	if err != nil {
		return nil, err
	}
	if pc == nil {
		return nil, fmt.Errorf("get nil check conn")
	}
	return pc, nil
}

// 检查是否超过下线阈值
// bool 表示是否需要将节点设为 StatusDown, int64 表示自 LastChecked 以来经过的时间，以便在外部直接用于日志记录
func (n *NodeInfo) ShouldDownAfterNoAlive(downAfterNoAlive int) (bool, int64) {
	elapsed := time.Now().Unix() - n.ConnPool.GetLastChecked()
	return elapsed >= int64(downAfterNoAlive), elapsed
}
