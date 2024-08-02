// Copyright 2024 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysql

import (
	"bufio"
	"testing"

	"github.com/XiaoMi/Gaea/util/mocks/pipeTest"
	"github.com/stretchr/testify/require"
)

// TestMariadbConnWithoutDB 为用来测试数据库一开始连线的详细流程，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestMariadbConnWithoutDB(t *testing.T) {
	// 函数测试开始
	t.Run("MariaDB连接 的抽换缓存测试", func(t *testing.T) {
		// 开始模拟
		mockClient, mockServer := pipeTest.NewDcServerClient(t, pipeTest.TestReplyMsgFunc) // 产生 Gaea 和 MariaDB 模拟物件

		// 针对这次测试进行临时修改
		err := mockClient.OverwriteConnBufWrite(nil, writersPool.Get().(*bufio.Writer))
		mockClient.GetBufWriter().Reset(mockClient.GetConnWrite())
		require.Equal(t, err, nil)

		// 产生一开始的讯息和预期讯息
		msg0 := []uint8{0}  // 起始传送讯息
		correct := uint8(0) // 预期之后的正确讯息

		// 开始进行讯息操作

		// 写入部份
		mockClient.SendOrReceiveMsg(msg0) // 模拟客户端传送讯息
		require.Equal(t, msg0[0], correct)

		// 读取部份
		msg1 := mockClient.ReplyMsg(mockServer) // 模拟服务端接收讯息
		correct++
		require.Equal(t, msg1[0], correct)
	})
}

func TestInitNetBufferSize(t *testing.T) {
	connBufferSize = 128
	InitNetBufferSize(0)
	require.Equal(t, connBufferSize, 128)
	InitNetBufferSize(128)
	require.Equal(t, connBufferSize, 128)
	InitNetBufferSize(512)
	require.Equal(t, connBufferSize, 512)
	InitNetBufferSize(16*1024 + 1)
	require.Equal(t, connBufferSize, 16*1024)
}
