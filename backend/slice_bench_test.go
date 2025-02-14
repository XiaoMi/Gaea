package backend

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

// 本 benchmark 测试使用 LocalSlaveReadForce 策略，模拟并发访问同一 Slice
// 其中期望所有请求均选中健康的本地节点
func BenchmarkGetSlaveConnNoLock(b *testing.B) {
	mockCtl := gomock.NewController(b)
	defer mockCtl.Finish()

	// 模拟测试用例：proxy dc 为 "c4"，三个节点，其中只有 "c4" 的节点健康
	slaveAddrs := []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"}
	// 指定各节点数据中心（第1个为 c3，其余为 c4）

	slaveStatus := []StatusCode{StatusUp, StatusUp, StatusDown}
	weights := []int{1, 1, 1}
	proxyDc := "c4"

	// 构造 DBInfo
	dbInfo, err := generateDBInfoWithWeights(mockCtl, slaveAddrs, slaveStatus, weights)
	if err != nil {
		b.Fatalf("generateDBInfoWithWeights failed: %v", err)
	}
	// 初始化 balancers
	if err := dbInfo.InitBalancers(proxyDc); err != nil {
		b.Fatal(err)
	}
	s := &Slice{Slave: dbInfo, ProxyDatacenter: proxyDc}

	b.ResetTimer()
	// 使用 b.RunParallel 模拟高并发调用
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 调用无加锁版本（内部并发调用时可能因状态竞争出现少量错误）
			_, err := s.GetSlaveConnNoLock(dbInfo, LocalSlaveReadClosed)
			if err != nil {
				b.Error(err)
				// 为 benchmark 忽略错误计数，可通过 b.Error(err) 输出，但这里不做处理
			}
		}
	})
}
func BenchmarkGetSlaveConnLock(b *testing.B) {
	mockCtl := gomock.NewController(b)
	defer mockCtl.Finish()

	// 模拟测试用例：proxy dc 为 "c4"，三个节点，其中只有 "c4" 的节点健康
	slaveAddrs := []string{"c3-mysql-test00.bj:3306", "c4-mysql-test02.bj:3310", "c4-mysql-test02.bj:3310"}
	// 指定各节点数据中心（第1个为 c3，其余为 c4）

	slaveStatus := []StatusCode{StatusUp, StatusUp, StatusDown}
	weights := []int{1, 1, 1}
	proxyDc := "c4"

	// 构造 DBInfo
	dbInfo, err := generateDBInfoWithWeights(mockCtl, slaveAddrs, slaveStatus, weights)
	if err != nil {
		b.Fatalf("generateDBInfoWithWeights failed: %v", err)
	}
	// 初始化 balancers
	if err := dbInfo.InitBalancers(proxyDc); err != nil {
		b.Fatal(err)
	}
	s := &Slice{Slave: dbInfo, ProxyDatacenter: proxyDc}

	b.ResetTimer()
	// 使用 b.RunParallel 模拟高并发调用
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 调用无加锁版本（内部并发调用时可能因状态竞争出现少量错误）
			_, err := s.GetSlaveConnWithLock(dbInfo, LocalSlaveReadClosed)
			if err != nil {
				b.Error(err)
				// 为 benchmark 忽略错误计数，可通过 b.Error(err) 输出，但这里不做处理
			}
		}
	})
}

// GetSlaveConn 根据读取策略选择使用 LocalBalancer 或 GlobalBalancer
func (s *Slice) GetSlaveConnNoLock(dbInfo *DBInfo, localSlaveReadPriority int) (PooledConnect, error) {
	if len(dbInfo.Nodes) == 0 {
		return nil, fmt.Errorf("no available slave DB")
	}
	switch localSlaveReadPriority {
	case LocalSlaveReadForce:
		if dbInfo.LocalBalancer == nil {
			return nil, fmt.Errorf("no primary balancer available")
		}
		// 这里调用不同版本，由上层选择
		return s.getConnFromBalancerNoLock(dbInfo, dbInfo.LocalBalancer)
	case LocalSlaveReadClosed:
		if dbInfo.GlobalBalancer == nil {
			return nil, fmt.Errorf("no secondary balancer available")
		}
		return s.getConnFromBalancerNoLock(dbInfo, dbInfo.GlobalBalancer)
	case LocalSlaveReadPrefer:
		if dbInfo.LocalBalancer != nil {
			if conn, err := s.getConnFromBalancerNoLock(dbInfo, dbInfo.LocalBalancer); err == nil {
				return conn, nil
			}
		}
		if dbInfo.GlobalBalancer == nil {
			return nil, fmt.Errorf("no secondary balancer available")
		}
		return s.getConnFromBalancerNoLock(dbInfo, dbInfo.GlobalBalancer)
	default:
		return nil, fmt.Errorf("invalid localSlaveReadPriority: %d", localSlaveReadPriority)
	}
}

func (s *Slice) GetSlaveConnWithLock(dbInfo *DBInfo, localSlaveReadPriority int) (PooledConnect, error) {
	if len(dbInfo.Nodes) == 0 {
		return nil, fmt.Errorf("no available slave DB")
	}
	switch localSlaveReadPriority {
	case LocalSlaveReadForce:
		if dbInfo.LocalBalancer == nil {
			return nil, fmt.Errorf("no primary balancer available")
		}
		// 这里调用不同版本，由上层选择
		return s.getConnFromBalancerLock(dbInfo, dbInfo.LocalBalancer)
	case LocalSlaveReadClosed:
		if dbInfo.GlobalBalancer == nil {
			return nil, fmt.Errorf("no secondary balancer available")
		}
		return s.getConnFromBalancerLock(dbInfo, dbInfo.GlobalBalancer)
	case LocalSlaveReadPrefer:
		if dbInfo.LocalBalancer != nil {
			if conn, err := s.getConnFromBalancerLock(dbInfo, dbInfo.LocalBalancer); err == nil {
				return conn, nil
			}
		}
		if dbInfo.GlobalBalancer == nil {
			return nil, fmt.Errorf("no secondary balancer available")
		}
		return s.getConnFromBalancerLock(dbInfo, dbInfo.GlobalBalancer)
	default:
		return nil, fmt.Errorf("invalid localSlaveReadPriority: %d", localSlaveReadPriority)
	}
}

func (s *Slice) getConnFromBalancerNoLock(slavesInfo *DBInfo, bal *balancer) (PooledConnect, error) {
	for i := 0; i < len(bal.roundRobinQ); i++ {
		index, err := bal.next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next index from balancer: %v", err)
		}
		node := slavesInfo.Nodes[index]
		// 直接检查 `NodeInfo.Status`
		if !node.IsStatusUp() {
			continue
		}
		// 返回健康节点的连接
		return s.getConnWithFuse(node)
	}
	return nil, fmt.Errorf("no healthy connection available from selected balancer")
}

func (s *Slice) getConnFromBalancerLock(slavesInfo *DBInfo, bal *balancer) (PooledConnect, error) {
	// 加锁保证同一时刻只有一个 Session 在使用该 balancer
	s.Lock()
	defer s.Unlock()
	for i := 0; i < len(bal.roundRobinQ); i++ {
		index, err := bal.next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next index from balancer: %v", err)
		}
		node := slavesInfo.Nodes[index]
		// 直接检查 `NodeInfo.Status`
		if !node.IsStatusUp() {
			continue
		}
		// 返回健康节点的连接
		return s.getConnWithFuse(node)
	}
	return nil, fmt.Errorf("no healthy connection available from selected balancer")
}
