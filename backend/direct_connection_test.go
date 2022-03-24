// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package backend

import (
	"bytes"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util/mocks/pipeTest"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestAppendSetVariable(t *testing.T) {
	var buf bytes.Buffer
	appendSetVariable(&buf, "charset", "utf8")
	t.Log(buf.String())
	appendSetVariable(&buf, "autocommit", 1)
	t.Log(buf.String())
	appendSetVariableToDefault(&buf, "sql_mode")
	t.Log(buf.String())
}

func TestAppendSetVariable2(t *testing.T) {
	var buf bytes.Buffer
	appendSetCharset(&buf, "utf8", "utf8_general_ci")
	t.Log(buf.String())
	appendSetVariable(&buf, "autocommit", 1)
	t.Log(buf.String())
	appendSetVariableToDefault(&buf, "sql_mode")
	t.Log(buf.String())
}

var (
	// preparation 准备数据库的回应资料

	// The initial handshake packet from MariaDB to Gaea
	// 第一个交握讯息，由 MariaDB 传送欢迎讯息到 Gaea
	mysqlInitHandShakeFirstResponseFromMaraiadbToGaea = []uint8{
		// length. 资料长度
		93, 0, 0,
		// the increment numbers. 自增串行号码
		0,

		// 93 bytes of data. 以下 93 笔数据

		// protocol 协定版本
		10,
		// version 数据库 版本
		53, 46, 53, 46, 53,
		45, 49, 48, 46, 53,
		46, 49, 50, 45, 77,
		97, 114, 105, 97, 68,
		66, 45, 108, 111, 103,
		// terminated 数据库的版本结尾
		0,
		// connection id 连线编号
		16, 0, 0, 0,
		// The first scramble. 第一部份的 scramble
		81, 64, 43, 85, 76, 90, 97, 91,
		// reserved byte 保留数据
		0,
		// capability 取得功能标志
		254, 247,
		// charset 数据库编码
		33,
		// status 服务器状态
		2, 0,
		// capability 延伸的功能标志
		255, 129,
		// auth 资料和保留值
		21, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0,
		// the second scramble. 延申的 scramble
		34, 53, 36, 85,
		93, 86, 117, 105,
		49, 87, 65, 125,
		// unused data. 其他未用到的资料
		0, 109, 121, 115, 113, 108, 95, 110, 97, 116,
		105, 118, 101, 95, 112, 97, 115, 115, 119, 111,
		114, 100, 0,
	}
)

// TestDirectConnWithoutDB is to test the initial handshake packet. The test doesn't use MariaDB.
// TestDirectConnWithoutDB 为测试数据库的后端连线流程，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestDirectConnWithoutDB(t *testing.T) {
	// check every step in the handshake process.
	// 开始正式测试，把每次的交握连接进行解测试

	// create the mock object
	// 先产生模拟对象
	mockGaea, mockMariaDB := pipeTest.NewDcServerClient(t, nil) // create the mock object for Game and MariaDB. 产生 Gaea 和 MariaDB 模拟对象
	var dc DirectConnection                                     // create the dc object and receive MariaDB's greet message. 对象 dc 会产生回应欢迎讯息给数据库

	// the first step
	// 交握第一步 Step1 测试数据库后端连线的初始交握
	t.Run("test the initial handshake", func(t *testing.T) {
		// start mocking, send message from MariaDB to Gaea.
		// 开始进行模拟，方向为 MariaDB 欢迎 Gaea
		mockMariaDB.SendOrReceiveMsg(mysqlInitHandShakeFirstResponseFromMaraiadbToGaea) // This message is mysqlInitHandShakeFirstResponseFromMaraiadbToGaea. 对象 mockMariaDB 会回传数据库资讯给 mockGaea，而讯息内容为 mysqlInitHandShakeFirstResponseFromMaraiadbToGaea，内容包含数据库资讯

		// create MariaDB object
		// 产生 MariaDB dc 直连对象 (用以下内容取代 reply() 函数 !)
		var connForReceivingMsgFromMariadb = mysql.NewConn(mockGaea.GetConnRead()) // wait for sending message completely. 等一下 MariaDB 数据库会把交握讯息传送到这
		// mysql.NewConn initializes net.conn and bufferedReader.
		// mysql.NewConn 会同时初始化 读取连接 net.conn 和 读取缓存 bufferedReader
		dc.conn = connForReceivingMsgFromMariadb // dc connects to test environment. 这时 dc 的连接 接通 整个测试环境
		err := dc.readInitialHandshake()         // Mock Gaea object reads the message. 模拟 Gaea 进行交握解析
		require.Equal(t, err, nil)

		// wait for sending messages to the pipe completely and reset the pipe.
		// 等待和确认资料已经写入 pipe 并单方向重置模拟对象
		err = mockMariaDB.WaitAndReset(mockGaea)
		require.Equal(t, err, nil)

		// the final check
		// 计算后的检查
		require.Equal(t, dc.capability, uint32(2181036030))                                                                   // capability 检查功能标志
		require.Equal(t, dc.conn.ConnectionID, uint32(16))                                                                    // connection id 检查连线编号
		require.Equal(t, dc.salt, []uint8{81, 64, 43, 85, 76, 90, 97, 91, 34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}) // Salt 检查
		require.Equal(t, dc.status, mysql.ServerStatusAutocommit)                                                             // Status 检查服务器状态
	})

	// the second step
	// 交握第二步 Step2 测试数据库后端连线初始交握后的回应
	t.Run("test response to the initial handshake", func(t *testing.T) {
		// make a response to MariaDB after receiving the initial handshake packet.
		// 开始进行模拟，方向为 Gaea 回应 MariaDB 欢迎

		// create MariaDB object
		// 产生 Mysql dc 直连对象 (用以下内容取代 reply() 函数 !)
		var connForSengingMsgToMariadb = mysql.NewConn(mockGaea.GetConnWrite()) // make the response to MariaDB's greet message here. 等一下会把要回应给 MariaDB 数据库回应的欢迎讯息写入到这里

		// mysql.NewConn initializes net.conn and bufferReader. However, bufferReader won't affect this result.
		// mysql.NewConn 会同时初始化 写入连接 net.conn 和 写入缓存 bufferReader，但这里 bufferReader 将不会用到，所以不会影响测试
		dc.conn = connForSengingMsgToMariadb // dc connects to test environment. 对象 dc 和整个测试进行连接
		dc.conn.StartWriterBuffering()       // initialize Gaea's bufferReader. 初始化 Gaea 的 写入缓存 bufferReader
		// Isolating bufferReader and no using StartWriterBuffering() in testing will be better. However, I cannot.
		// 最好的状况是 Gaea 的 写入缓存 bufferReader 和这个测试整个分离，但目前现有代码的函数不支持，又不想修改现在代码，所以就把 写入缓存 捉进来一起测试

		// fulfill the dc object, including user, password, etc.
		// 填入 Gaea 用户的资讯，包含 密码
		dc.user = "xiaomi"     // user 帐户名称
		dc.password = "12345"  // password 密码
		dc.charset = "utf8mb4" // charset 数据库编码
		dc.collation = 46      // collation 文本排序

		// use the anonymous function to send the message.
		// 使用支持使用匿名函数传送讯息
		responseMsg := mockGaea.UseAnonymousFuncSendMsg(
			// start using anonymous function
			// 自订的匿名函数开始
			func() {
				err := dc.writeHandshakeResponse41() // write the response message. 写入回应讯息
				require.Equal(t, err, nil)
				err = dc.conn.Flush() // flush dc connection. 移动讯息由缓存到 Pipe
				require.Equal(t, err, nil)
				err = mockGaea.GetConnWrite().Close() // close the pipe. 关闭连线
				require.Equal(t, err, nil)
			},
			// use anonymous function completely.
			// 自订的匿名函数结束
		).CheckArrivedMsg(mockMariaDB) // get the arrived message and check. 对传送到达对方的讯息取出进行确认

		// check result. 确认结果
		require.Equal(t, len(responseMsg), 64) // check length of the packet. 确认封包长度

		require.Equal(t, strings.Contains(string(responseMsg), "xiaomi"), true) // check the existence of the account in the packet. 确认 用户帐户 是否真的写到封包里

		scramble := mysql.CalcPassword(dc.salt, []byte("12345"))                        // calculate token. 计算 token
		require.Equal(t, strings.Contains(string(responseMsg), string(scramble)), true) // check the existence of the token in the packet. 确认 token 是否真的写到封包里

		require.Equal(t, strings.Contains(string(responseMsg), dc.password), false) // check the non-existence of the password in the packet. 确认 密码 不行写到封包里

	})
}
