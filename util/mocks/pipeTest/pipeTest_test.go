package pipeTest

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPipeTestWorkable 为验证测试 PipeTest 是否能正常运作，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestPipeTestWorkable(t *testing.T) {
	t.Run("此为测试 PipeTest 的验证测试，主要是用来确认整个测试流程没有问题", func(t *testing.T) {
		// 开始模拟物件
		mockClient, mockServer := NewDcServerClient(t, TestReplyMsgFunc) // 产生 mockClient 和 mockServer 模拟物件

		// 产生一开始的讯息和之后的预期正确讯息
		msg0 := []uint8{0}  // 起始传送讯息
		correct := uint8(0) // 之后的预期正确讯息

		// 产生一连串的接收和回应的操作
		for i := 0; i < 5; i++ {
			msg1 := mockClient.SendOrReceiveMsg(msg0).ReplyMsg(mockServer) // 接收和回应
			correct++                                                      // 每经过一个reply() 函数时，回应讯息会加1
			require.Equal(t, msg1[0], correct)
			msg0 = mockServer.SendOrReceiveMsg(msg1).ReplyMsg(mockClient) // 接收和回应
			correct++                                                     // 每经过一个reply() 函数时，回应讯息会加1
			require.Equal(t, msg0[0], correct)
		}
	})
}
