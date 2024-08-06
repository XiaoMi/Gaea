package server

import (
	"fmt"
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util"
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 检测 namespace changed 时， ks conn 是否被清除
func TestSessionRunNamespaceChangedWithMock(t *testing.T) {
	globalSe, _ := newDefaultSessionExecutor(nil)
	namespaceName := globalSe.namespace
	manager := globalSe.manager
	var NewNamespace = func() *models.Namespace {
		return &models.Namespace{
			Name: namespaceName,
			Slices: []*models.Slice{
				{Name: "slice-0", Master: "127.0.0.1:3306"},
			},
			DefaultSlice: "slice-0",
		}
	}

	mockey.PatchConvey("test", t, func() {
		mockey.Mock((*mysql.Conn).SetSequence).Return().Build()
		mockey.Mock((*mysql.Conn).Close).Return().Build()
		mockey.Mock((*mysql.Conn).ReadEphemeralPacket).To(func() ([]byte, error) {
			time.Sleep(time.Millisecond * 10)
			return []byte{0, 0, 0, 0}, nil
		}).Build()
		mockey.Mock((*mysql.Conn).RecycleReadPacket).Return().Build()
		mockey.Mock((*mysql.Conn).GetConnectionID).Return(0).Build()
		mockey.Mock((*SessionExecutor).ExecuteCommand).Return(Response{}).Build()
		mockey.Mock((*Session).writeResponse).Return(nil).Build()
		mockey.Mock((*Manager).GetStatisticManager).Return(&StatisticManager{}).Build()
		mockey.Mock((*StatisticManager).IncrSessionCount).Return().Build()
		mockey.Mock((*StatisticManager).IncrConnectionCount).Return().Build()
		mockey.Mock((*StatisticManager).DescSessionCount).Return().Build()
		mockey.Mock((*StatisticManager).DescConnectionCount).Return().Build()
		mockey.Mock((*StatisticManager).AddReadFlowCount).Return().Build()
		mockey.Mock((*StatisticManager).initBackend).Return(nil).Build()
		mockey.Mock((*StatisticManager).Init).Return(nil).Build()
		mockey.Mock((*backend.Slice).CheckStatus).Return().Build()
		mockey.Mock((*util.TimeWheel).Remove).Return(nil).Build()
		mockey.Mock((*backend.MockPooledConnect).Recycle).Return().Build()
		mockey.Mock((*backend.MockPooledConnect).Rollback).Return(nil).Build()
		mockey.Mock((*backend.MockPooledConnect).Close).Return().Build()
		mockey.Mock(util.GetHostDatacenter).Return("", nil).Build()

		sessionCount := 100
		var g sync.WaitGroup
		seList := make([]*SessionExecutor, sessionCount)
		for i := 0; i < sessionCount; i++ {
			se, _ := newDefaultSessionExecutor(nil)
			se.manager = manager
			se.session.manager = manager
			se.session.namespace = namespaceName
			seList[i] = se
		}
		for _, se := range seList {
			go func(se *SessionExecutor) {
				defer g.Done()
				g.Add(1)
				se.ksConns = make(map[string]backend.PooledConnect)
				v := &backend.MockPooledConnect{}
				se.ksConns["slice-0"] = v
				se.session.executor = se
				se.session.executor.keepSession = true
				closeStatus := atomic.Value{}
				closeStatus.Store(false)
				se.session.closed = closeStatus
				se.session.c = &ClientConn{}
				se.SetContextNamespace()
				se.session.Run()
			}(se)
		}
		time.Sleep(time.Millisecond * 50)
		hasKsConnSeNum := 0
		for _, se := range seList {
			// 初始连接需不在事务中
			assert.Equal(t, se.isInTransaction(), false, "not in tx")
			if len(se.ksConns) > 0 {
				hasKsConnSeNum++
			}
		}
		assert.Equal(t, sessionCount, hasKsConnSeNum, "begin test")
		manager.ReloadNamespacePrepare(NewNamespace())
		manager.ReloadNamespaceCommit(namespaceName)

		time.Sleep(time.Millisecond * 50)
		for _, se := range seList {
			if len(se.ksConns) == 0 {
				hasKsConnSeNum--
			}
		}
		// 连接必须全部回收
		assert.Equal(t, hasKsConnSeNum, 0, "hasKsConnSeNum")
		closedNum := 0
		for _, se := range seList {
			if se.session.IsClosed() {
				closedNum += 1
			}
		}
		// 非事务连接不关闭
		assert.Equal(t, closedNum, 0, "closedNum")
		for _, se := range seList {
			se.session.Close()
		}
		g.Wait()
	})
}

// 检测 namespace changed 且 有事务时， session 关闭，且 ks conn 被清除
func TestSessionRunInTxNamespaceChangedWithMock(t *testing.T) {
	globalSe, _ := newDefaultSessionExecutor(nil)
	namespaceName := globalSe.namespace
	manager := globalSe.manager
	var NewNamespace = func() *models.Namespace {
		return &models.Namespace{
			Name: namespaceName,
			Slices: []*models.Slice{
				{Name: "slice-0", Master: "127.0.0.1:3306"},
			},
			DefaultSlice: "slice-0",
		}
	}

	mockey.PatchConvey("test", t, func() {
		mockey.Mock((*mysql.Conn).SetSequence).Return().Build()
		mockey.Mock((*mysql.Conn).Close).Return().Build()
		mockey.Mock((*mysql.Conn).ReadEphemeralPacket).To(func() ([]byte, error) {
			time.Sleep(time.Millisecond * 10)
			return []byte{0, 0, 0, 0}, nil
		}).Build()
		mockey.Mock((*mysql.Conn).RecycleReadPacket).Return().Build()
		mockey.Mock((*mysql.Conn).GetConnectionID).Return(0).Build()
		mockey.Mock((*SessionExecutor).ExecuteCommand).Return(Response{}).Build()
		mockey.Mock((*Session).writeResponse).Return(nil).Build()
		mockey.Mock((*Manager).GetStatisticManager).Return(&StatisticManager{}).Build()
		mockey.Mock((*StatisticManager).IncrSessionCount).Return().Build()
		mockey.Mock((*StatisticManager).IncrConnectionCount).Return().Build()
		mockey.Mock((*StatisticManager).DescSessionCount).Return().Build()
		mockey.Mock((*StatisticManager).DescConnectionCount).Return().Build()
		mockey.Mock((*StatisticManager).AddReadFlowCount).Return().Build()
		mockey.Mock((*StatisticManager).initBackend).Return(nil).Build()
		mockey.Mock((*StatisticManager).Init).Return(nil).Build()
		mockey.Mock((*backend.Slice).CheckStatus).Return().Build()
		mockey.Mock((*util.TimeWheel).Remove).Return(nil).Build()
		mockey.Mock((*backend.MockPooledConnect).Recycle).Return().Build()
		mockey.Mock((*backend.MockPooledConnect).Rollback).Return(nil).Build()
		mockey.Mock((*backend.MockPooledConnect).Close).Return().Build()
		mockey.Mock(util.GetHostDatacenter).Return("", nil).Build()

		sessionCount := 100
		var g sync.WaitGroup
		seList := make([]*SessionExecutor, sessionCount)
		for i := 0; i < sessionCount; i++ {
			se, _ := newDefaultSessionExecutor(nil)
			se.manager = manager
			se.session.manager = manager
			se.session.namespace = namespaceName
			seList[i] = se
		}
		for _, se := range seList {
			go func(se *SessionExecutor) {
				g.Add(1)
				defer g.Done()
				se.ksConns = make(map[string]backend.PooledConnect)
				se.status = se.status | mysql.ServerStatusInTrans
				v := &backend.MockPooledConnect{}
				se.ksConns["slice-0"] = v
				se.session.executor = se
				se.session.executor.keepSession = true
				closeStatus := atomic.Value{}
				closeStatus.Store(false)
				se.session.closed = closeStatus
				se.session.c = &ClientConn{}
				se.SetContextNamespace()
				se.session.Run()
			}(se)
		}
		time.Sleep(time.Millisecond * 50)
		hasKsConnSeNum := 0
		for _, se := range seList {
			// 初始连接需要在事务中
			assert.Equal(t, se.isInTransaction(), true)
			if len(se.ksConns) > 0 {
				hasKsConnSeNum++
			}
		}
		assert.Equal(t, sessionCount, hasKsConnSeNum)
		manager.ReloadNamespacePrepare(NewNamespace())
		manager.ReloadNamespaceCommit(namespaceName)
		time.Sleep(time.Millisecond * 50)
		for _, se := range seList {
			if len(se.ksConns) == 0 {
				hasKsConnSeNum--
			}
		}
		// 连接必须全部回收
		assert.Equal(t, hasKsConnSeNum, 0)
		closedNum := 0
		for _, se := range seList {
			if se.session.IsClosed() {
				closedNum += 1
			}
		}
		// 连接必须全部关闭
		assert.Equal(t, sessionCount, closedNum)
		for _, se := range seList {
			se.session.Close()
		}
		g.Wait()
	})
}

func TestNamespaceChangeTimesMock(t *testing.T) {
	se, _ := newDefaultSessionExecutor(nil)
	namespaceName := se.namespace

	var NewNamespace = func() *models.Namespace {
		return &models.Namespace{
			Name: namespaceName,
			Slices: []*models.Slice{
				{Name: "slice-0", Master: "127.0.0.1:3306"},
			},
			DefaultSlice: "slice-0",
		}
	}

	for i := 0; i <= 3; i += 1 {
		oldNmaspace := se.GetManagerNamespace()
		se.manager.ReloadNamespacePrepare(NewNamespace())
		se.manager.ReloadNamespaceCommit(namespaceName)
		fmt.Println("change index old : ", oldNmaspace.namespaceChangeIndex)
		fmt.Println("change index new: ", se.GetManagerNamespace().namespaceChangeIndex)
		assert.Equal(t, se.GetManagerNamespace().namespaceChangeIndex > oldNmaspace.namespaceChangeIndex, true)
	}
}
