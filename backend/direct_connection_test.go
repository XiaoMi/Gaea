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
	mysqlInitHandShakeResponse = []uint8{
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
)

// TestDirectConnWithoutDB 为测试数据库的后端连线流程，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestDirectConnWithoutDB(t *testing.T) {
	// 开始正式测试
	t.Run("测试数据库后端连线的初始交握", func(t *testing.T) {
		// 开始模拟
		mockGaea, mockMariaDB := pipeTest.NewDcServerClient(t, pipeTest.TestReplyFunc) // 产生 Gaea 和 mockServer 模拟物件
		mockGaea.SendOrReceive(mysqlInitHandShakeResponse)                             // 模拟数据库开始交握

		// 產生 Mysql dc 直連物件 (用以下内容取代 reply() 函数 !)
		var dc DirectConnection
		var mysqlConn = mysql.NewConn(mockMariaDB.GetConnRead())
		dc.conn = mysqlConn
		err := dc.readInitialHandshake()
		require.Equal(t, err, nil)

		// 等待和确认资料已经写入 pipe 并单方向重置模拟物件
		err = mockGaea.WaitAndReset(mockMariaDB)
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
}
