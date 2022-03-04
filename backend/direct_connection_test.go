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
		mockMariaDB.SendOrReceive(mysqlInitHandShakeFirstResponseFromMaraiadbToGaea) // mockMariaDB 会回传数据库资讯给 mockGaea，而讯息内容为 mysqlInitHandShakeFirstResponseFromMaraiadbToGaea，内容包含数据库资讯

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

		// 以下代码正确性待确认
		customFunc := func() {
			err := dc.writeHandshakeResponse41() // 模拟 Gaea 进行交握解析
			require.Equal(t, err, nil)
			err = dc.conn.Flush()
			require.Equal(t, err, nil)
			err = mockGaea.GetConnWrite().Close()
			require.Equal(t, err, nil)
		}

		fmt.Println(mockGaea.CustomSend(customFunc).ArrivedMsg(mockMariaDB))
	})
}
