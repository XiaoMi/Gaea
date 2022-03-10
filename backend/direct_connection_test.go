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
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util/mocks/pipeTest"
	"github.com/stretchr/testify/require"
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
	// 准备数据库的回应资料

	// 第一个交握讯息，由 MariaDB 传送欢迎讯息到 Gaea
	mysqlInitHandShakeFirstResponseFromMaraiadbToGaea = []uint8{
		// 资料长度
		93, 0, 0,
		// 自增序列号码
		0,
		// 以下 93 笔数据
		// 数据库的版本号 version
		10, 53, 46, 53, 46,
		53, 45, 49, 48, 46,
		53, 46, 49, 50, 45,
		77, 97, 114, 105, 97,
		68, 66, 45, 108, 111,
		103,
		// 数据库的版本结尾
		0,
		// 连线编号 connection id
		16, 0, 0, 0,
		// Salt
		81, 64, 43, 85, 76, 90, 97, 91,
		// filter
		0,
		// 取得功能标志 capability
		254, 247,
		// 數據庫編碼 charset
		33, // 可以用 SHOW CHARACTER SET LIKE 'utf8'; 查询
		// 服务器状态，在 Gaea/mysql/constants.go 的 Server information
		2, 0,
		// 延伸的功能标志 capability
		255, 129,
		// Auth 资料和保留值
		21, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0,
		// 延伸的 Salt
		34, 53, 36, 85,
		93, 86, 117, 105, 49,
		87, 65, 125,
		// 其他未用到的资料
		0, 109, 121, 115, 113, 108, 95, 110, 97, 116,
		105, 118, 101, 95, 112, 97, 115, 115, 119, 111,
		114, 100, 0,
	}

	// 第二个交握讯息，由 Gaea 回应欢迎讯息给 MariaDB
	mysqlInitHandShakeSecondResponseFromGaeaToMaraiadb = []uint8{
		//
	}
)

// TestDirectConnWithoutDB 为测试数据库的后端连线流程，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestDirectConnWithoutDB(t *testing.T) {
	// 开始正式测试，把每次的交握连接进行解测试

	// 先产生模拟对象
	mockGaea, mockMariaDB := pipeTest.NewDcServerClient(t, nil) // 产生 Gaea 和 MariaDB 模拟物件
	var dc DirectConnection                                     // dc 对象会产生回应欢迎讯息给数据库

	// 交握第一步 Step1
	t.Run("测试数据库后端连线的初始交握", func(t *testing.T) {
		// 开始进行模拟，方向为 MariaDB 欢迎 Gaea
		mockMariaDB.SendOrReceiveMsg(mysqlInitHandShakeFirstResponseFromMaraiadbToGaea) // mockMariaDB 会回传数据库资讯给 mockGaea，而讯息内容为 mysqlInitHandShakeFirstResponseFromMaraiadbToGaea，内容包含数据库资讯

		// 產生 Mysql dc 直連物件 (用以下内容取代 reply() 函数 !)
		var connForReceivingMsgFromMariadb = mysql.NewConn(mockGaea.GetConnRead()) // 等一下 MariaDB 数据库会把交握讯息传送到这
		// mysql.NewConn 会同时初始化 读取连接 net.conn 和 读取缓存 bufferedReader
		dc.conn = connForReceivingMsgFromMariadb // 这时 dc 的连接 接通 整个测试环境
		err := dc.readInitialHandshake()         // 模拟 Gaea 进行交握解析
		require.Equal(t, err, nil)

		// 等待和确认资料已经写入 pipe 并单方向重置模拟物件
		err = mockMariaDB.WaitAndReset(mockGaea)
		require.Equal(t, err, nil)

		// 开始计算

		/* 功能标志 capability 的计算
		先把所有的功能标志 capability 的数据收集起来，包含延伸部份
		数值分别为 254, 247, 255, 129
		并反向排列
		数值分别为 129, 255, 247, 254
		全部 十进制 转成 二进制
		254 的二进制为 1111 1110
		247 的二进制为 1111 0111
		255 的二进制为 1111 1111
		129 的二进制为 1000 0001
		把全部二进制的数值合并
		二进制数值分别为 1000 0001 1111 1111 1111 0111 1111 1110 (转成十进制数值为 2181036030)
		再用文档 https://mariadb.com/kb/en/connection/ 进行对照
		比如，功能标志 capability 的第一个值为 0，意思为 CLIENT_MYSQL 值为 0，代表是由服务器发出的讯息 */

		/* 连线编号 connection id 的计算
		先把所有的连线编号 connection id 的数据收集起来，包含延伸部份
		数值分别为 16, 0, 0, 0
		并反向排列
		数值分别为 0, 0, 0, 16
		全部 十进制 转成 二进制
		  0 的二进制为 0000 0000
		 16 的二进制为 0001 0000
		把全部二进制的数值合并
		二进制数值分别为 0000 0000 0001 0000 (转成十进制数值为 16) */

		// 先把所有 Salt 的数据收集起来，包含延伸部份
		// 数值分别为 81,64,43,85,76,90,97,91,34,53,36,85,93,86,117,105,49,87,65,125

		/* 服务器状态 status 的计算
		先把所有的服务器状态 的数据收集起来，包含延伸部份
		数值分别为 2, 0
		并反向排列
		数值分别为 0, 2
		全部 十进制 转成 二进制
		2 的二进制为 0000 0010
		0 的二进制为 0000 0000
		把全部二进制的数值合并
		二进制数值分别为 0000 0000 0000 0010 (转成十进制数值为 2)
		再用代码 Gaea/mysql/constants.go 里的 Server information 进行对照
		功能标志 capability 的第一个值为 0，意思为 CLIENT_MYSQL 值为 0，代表是由服务器发出的讯息 */

		// 计算后的检查
		require.Equal(t, dc.capability, uint32(2181036030))                                                                   // 检查功能标志 capability
		require.Equal(t, dc.conn.ConnectionID, uint32(16))                                                                    // 检查连线编号 connection id
		require.Equal(t, dc.salt, []uint8{81, 64, 43, 85, 76, 90, 97, 91, 34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}) // 检查 Salt
		require.Equal(t, dc.status, mysql.ServerStatusAutocommit)                                                             // 检查服务器状态
	})

	// 交握第二步 Step2
	t.Run("测试数据库后端连线初始交握后的回应", func(t *testing.T) {
		// 开始进行模拟，方向为 Gaea 回应 MariaDB 欢迎
		// 產生 Mysql dc 直連物件 (用以下内容取代 reply() 函数 !)
		var connForSengingMsgToMariadb = mysql.NewConn(mockGaea.GetConnWrite()) // 等一下会把要回应给 MariaDB 数据库回应的欢迎讯息写入到这里
		// mysql.NewConn 会同时初始化 写入连接 net.conn 和 写入缓存 bufferReader，但这里 bufferReader 将不会用到，所以不会影响测试
		dc.conn = connForSengingMsgToMariadb // dc 对象和整个测试进行连接
		dc.conn.StartWriterBuffering()       // 初始化 Gaea 的 写入缓存 bufferReader
		// 最好的状况是 Gaea 的 写入缓存 bufferReader 和这个测试整个分离，但目前现有代码的函数不支援，又不想修改现在代码，所以就把 写入缓存 捉进来一起测试

		// 填入 Gaea 用户的资讯，包含 密码
		dc.password = "12345"

		// 印出传送到达对方的讯息
		fmt.Println(
			// 使用支援使用匿名函式传送讯息
			mockGaea.UseAnonymousFuncSendMsg(
				// 自订的匿名函式开始
				func() {
					err := dc.writeHandshakeResponse41() // 模拟 Gaea 进行交握解析
					require.Equal(t, err, nil)
					err = dc.conn.Flush()
					require.Equal(t, err, nil)
					err = mockGaea.GetConnWrite().Close()
					require.Equal(t, err, nil)
				},
				// 自订的匿名函式结束
			).CheckArrivedMsg(mockMariaDB)) // 对传送到达对方的讯息取出进行确认
	})

	// 开始计算

	/*
		先计算安全密码
		参考官方文件，关于安全码的说明
		https://dev.mysql.com/doc/internals/en/secure-password-authentication.html

		安全密码产生规则如下
		SHA1( password ) XOR SHA1( "20-bytes random data from server" <concat> SHA1( SHA1( password ) ) )

		由上述公式可以看出
		自己输入的密码 12345，要经过两次 Sha1sum 的转换
		由服务器传送来的 Salt 乱数最后也要经过一次的 Sha1sum 的转换

		在 Linux 执行命令，获得第一次 Sha1sum 转换的结果
		[Xiaomi@Gaea]:~ $ echo -n 12345 | sha1sum | head -c 40
		8cb2237d0679ca88db6464eac60da96345513964 # 会得到第一次 Sha1sum 转换的结果

		只不过 intellj 或者是 Goland 之类工具会把 Sha1sum 验证码用人类最熟悉的十进位显示出来，内容如下
			stage1 = {[]uint8} len:20, cap:24
				0 = {uint8} 140
				1 = {uint8} 178
				2 = {uint8} 35
				3 = {uint8} 125
				4 = {uint8} 6
				5 = {uint8} 121
				6 = {uint8} 202
				7 = {uint8} 136
				8 = {uint8} 219
				9 = {uint8} 100
				10 = {uint8} 100
				11 = {uint8} 234
				12 = {uint8} 198
				13 = {uint8} 13
				14 = {uint8} 169
				15 = {uint8} 99
				16 = {uint8} 69
				17 = {uint8} 81
				18 = {uint8} 57
				19 = {uint8} 100

		只要把十六进位 两个 两个一组 转成 十进位，就会发现其结果都为相同的，只是一个为十六进位显示，另一个为十进位显示而已

		先把十六进位的 Sha1sum 验证码 一个byte 一个byte 分隔开来，如下所示
		8c b2 23 7d 06 79 ca 88 db 64 64 ea c6 0d a9 63 45 51 39 64 <= 十六进位 一个byte 一个byte一组

		也把十进位的 Sha1sum 验证码 一个byte 一个byte 分隔开来，如下所示
		140 178 35 125 6 121 202 136 219 100 100 234 198 13 169 99 69 81 57 100

		再用下表进行对照，可以发现其结果为相同的

		十六进位的 8c 为十进位的 140 (也为二进位的 10001100)
		十六进位的 b2 为十进位的 178 (也为二进位的 10110010)
		十六进位的 23 为十进位的  35 (也为二进位的 00100011)
		十六进位的 7d 为十进位的 125 (也为二进位的 01111101)
		十六进位的 06 为十进位的   6 (也为二进位的 00000110)
		十六进位的 79 为十进位的 121 (也为二进位的 01111001)
		十六进位的 ca 为十进位的 202 (也为二进位的 11001010)
		十六进位的 88 为十进位的 136 (也为二进位的 10001000)
		十六进位的 db 为十进位的 219 (也为二进位的 11011011)
		十六进位的 64 为十进位的 100 (也为二进位的 01100100)
		十六进位的 64 为十进位的 100 (也为二进位的 01100100)
		十六进位的 ea 为十进位的 234 (也为二进位的 11101010)
		十六进位的 c6 为十进位的 198 (也为二进位的 11000110)
		十六进位的 0d 为十进位的  13 (也为二进位的 00001101)
		十六进位的 a9 为十进位的 169 (也为二进位的 10101001)
		十六进位的 63 为十进位的  99 (也为二进位的 01100011)
		十六进位的 45 为十进位的  69 (也为二进位的 01000101)
		十六进位的 51 为十进位的  81 (也为二进位的 01010001)
		十六进位的 39 为十进位的  57 (也为二进位的 00111001)
		十六进位的 64 为十进位的 100 (也为二进位的 01100100)

		在 Linux 执行命令，获得连续二次 Sha1sum 转换的结果
		[Xiaomi@Gaea]:~ $ echo -n 12345 | sha1sum | xxd -r -p | sha1sum | head -c 40
		00a51f3f48415c7d4e8908980d443c29c69b60c9 # 会得到连续二次 Sha1sum 转换的结果

		只不过 intellj 或者是 Goland 之类工具同样也会把 Sha1sum 验证码用人类最熟悉的十进位显示出来，内容如下
		hash = {[]uint8} len:20, cap:24
				0 = {uint8} 0
				1 = {uint8} 165
				2 = {uint8} 31
				3 = {uint8} 63
				4 = {uint8} 72
				5 = {uint8} 65
				6 = {uint8} 92
				7 = {uint8} 125
				8 = {uint8} 78
				9 = {uint8} 137
				10 = {uint8} 8
				11 = {uint8} 152
				12 = {uint8} 13
				13 = {uint8} 68
				14 = {uint8} 60
				15 = {uint8} 41
				16 = {uint8} 198
				17 = {uint8} 155
				18 = {uint8} 96
				19 = {uint8} 201

		先把十六进位的 Sha1sum 验证码 一个byte 一个byte 分隔开来，如下所示
		00 a5 1f 3f 48 41 5c 7d 4e 89 08 98 0d 44 3c 29 c6 9b 60 c9 <= 十六进位 一个byte 一个byte一组

		也把十进位的 Sha1sum 验证码 一个byte 一个byte 分隔开来，如下所示
		0 165 31 63 72 65 92 125 78 137 8 152 13 68 60 41 198 155 96 201

		再用下表进行对照，可以发现其结果为相同的

		十六进位的 00 为十进位的   0 (也为二进位的 00000000)
		十六进位的 a5 为十进位的 165 (也为二进位的 10100101)
		十六进位的 1f 为十进位的  31 (也为二进位的 00011111)
		十六进位的 3f 为十进位的  63 (也为二进位的 00111111)
		十六进位的 48 为十进位的  72 (也为二进位的 01001000)
		十六进位的 41 为十进位的  65 (也为二进位的 01000001)
		十六进位的 5c 为十进位的  92 (也为二进位的 01011100)
		十六进位的 7d 为十进位的 125 (也为二进位的 01111101)
		十六进位的 4 为十进位的  78  (也为二进位的 01001110)
		十六进位的 89 为十进位的 137 (也为二进位的 10001001)
		十六进位的 08 为十进位的   8 (也为二进位的 00001000)
		十六进位的 98 为十进位的 152 (也为二进位的 10011000)
		十六进位的 0d 为十进位的  13 (也为二进位的 00001101)
		十六进位的 44 为十进位的  68 (也为二进位的 01000100)
		十六进位的 3c 为十进位的  60 (也为二进位的 00111100)
		十六进位的 29 为十进位的  41 (也为二进位的 00101001)
		十六进位的 c6 为十进位的 198 (也为二进位的 11000110)
		十六进位的 9b 为十进位的 155 (也为二进位的 10011011)
		十六进位的 60 为十进位的  96 (也为二进位的 01100000)
		十六进位的 c9 为十进位的 201 (也为二进位的 11001001)

		最后发现其结果也相符合
	*/

}
