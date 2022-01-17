// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

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

package etcdclientv3

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_EtcdV3(t *testing.T) {
	// 设定颜色输出
	colorReset := "\033[0m"
	colorRed := "\033[31m"

	remote, err := New("http://127.0.0.1:2379", 3*time.Second, "", "", "")
	if err != nil {
		// 如果 etcd 连线失败
		fmt.Println(colorRed, "目前找不到可实验的 Etcd 服务器", colorReset)
	}
	if err == nil {
		// >>>>> >>>>> >>>>> 如果 etcd 连线成功，执行以下成功，执行以下工作

		// >>>>> 先进行 (1) 新增测试 (新增 key1 和 key2)

		err = remote.Update("key1", []byte("value1")) // 写入 key1
		require.Equal(t, err, nil)

		err = remote.Update("key2", []byte("value2")) // 写入 key2
		require.Equal(t, err, nil)

		keys, err := remote.List("key") // 先列出所有的 key 值
		require.Equal(t, err, nil)
		require.Equal(t, keys, []string{"key1", "key2"})

		byte1, err := remote.Read("key1") // 检查 key1 值
		require.Equal(t, err, nil)
		require.Equal(t, string(byte1), "value1")

		byte2, err := remote.Read("key2") // 检查 key2 值
		require.Equal(t, err, nil)
		require.Equal(t, string(byte2), "value2")

		// return // 先中断，截图 (1) 新增测试

		// >>>>> 接着进行 (2) 删除测试 (删除 key1)

		err = remote.Delete("key1") // 删除 key1
		require.Equal(t, err, nil)

		keys, err = remote.List("key") // 列出所有的 key 值
		require.Equal(t, err, nil)
		require.Equal(t, keys, []string{"key2"})

		// return // 先中断，截图 (2) 删除测试

		// >>>>> 接着进行 (3) 到时删除测试 (测试到时删除 key 值)

		err = remote.UpdateWithTTL("key3", []byte("value3"), 5*time.Second) // key3 只保留 5 秒
		require.Equal(t, err, nil)

		keys, err = remote.List("key") // 先列出所有的 key 值
		require.Equal(t, err, nil)
		require.Equal(t, keys, []string{"key2", "key3"}) // 列出所有的 key 值

		byte3, err := remote.Read("key3") // 检查 key3 值
		require.Equal(t, err, nil)
		require.Equal(t, string(byte3), "value3")

		time.Sleep(6 * time.Second) // 停留 6 秒

		keys, err = remote.List("key") // 先列出所有的 key 值
		require.Equal(t, err, nil)
		require.Equal(t, keys, []string{"key2"}) // 列出所有的 key 值

		// return // 先中断，截图 (3) 到时删除测试

		// >>>>> 接着进行 (4) 追踪测试

		channel := make(chan string, 2)
		go func(ch chan string) error {
			err := remote.Watch("key", ch) // 追踪 key 值，只要 key1 key2 等值增加时，就会经由通道去通知
			return err
		}(channel)

		time.Sleep(1 * time.Second)

		err = remote.Update("key5", []byte("value5-1")) // 目前为永久保存 key5，值为 value5-1
		require.Equal(t, err, nil)

		msg := <-channel
		require.Equal(t, msg, "key5") // 将会收到新增 key5 的讯息

		err = remote.UpdateWithTTL("key5", []byte("value5-2"), 1*time.Second) // 更新 key5 只保存 1 秒，值为 value5-2
		require.Equal(t, err, nil)

		msg = <-channel
		require.Equal(t, msg, "key5") // 将会收到新增 key5 的讯息

		time.Sleep(4 * time.Second) // 停留 4 秒

		// return // 先中断，截图 (4) 追踪测试

		// >>>>> 接着进行 (5) 租约测试

		leaseID, err := remote.Lease(5 * time.Second) // 先产生 5 秒的租约
		require.Equal(t, err, nil)

		// 利用 5 秒租约产生 key6 和 key7
		err = remote.UpdateWithLease("key6", []byte("value6"), leaseID) // key5 只保留 2 秒
		require.Equal(t, err, nil)
		err = remote.UpdateWithLease("key7", []byte("value7"), leaseID) // key5 只保留 2 秒
		require.Equal(t, err, nil)

		keys, err = remote.List("key") // 先列出所有的 key 值
		require.Equal(t, err, nil)
		require.Equal(t, keys, []string{"key2", "key6", "key7"}) // 列出所有的 key 值

		time.Sleep(6 * time.Second) // 停留 6 秒

		keys, err = remote.List("key") // 先列出所有的 key 值
		require.Equal(t, err, nil)
		require.Equal(t, keys, []string{"key2"}) // 列出所有的 key 值

		// return // 先中断，截图 (5) 租约测试

		// >>>>> 最后 (6) 复原测试环境
		// (一共产生 key1 key2 key3 key 5 ,key3 和 key5 为定时删除，key1 之前用删除了，只剩 key2)

		err = remote.Delete("key2") // 删除 key2
		require.Equal(t, err, nil)

		keys, err = remote.List("key") // 列出所有的 key 值
		require.Equal(t, err, nil)
		require.Equal(t, len(keys), 0) // 最后所有的 key 都被清空

		// return // 先中断，截图 (6) 复原测试环境
	}
}
