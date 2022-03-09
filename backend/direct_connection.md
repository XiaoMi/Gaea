# backend 包的直连 direct_connection

以下内容都还在考量中，未定案

## 代码



## 测试

> 以下会说明在写测试时考量的点

### 匿名函数的考量

在以下代码内有一个 测试函数 t.Run("测试数据库后端连线初始交握后的回应", func(t *testing.T)

此测试函数内含 匿名函数 customFunc

以下代码，匿名函数 customFunc 将会取 dc 对象的內存位置，考量后，觉得可以这样写

```go
	// 交握第二步 Step2
	t.Run("测试数据库后端连线初始交握后的回应", func(t *testing.T) {
		var connForSengingMsgToMariadb = mysql.NewConn(mockGaea.GetConnWrite())
		dc.conn = connForSengingMsgToMariadb
		dc.conn.StartWriterBuffering()
        
		customFunc := func() {
			err := dc.writeHandshakeResponse41()
			require.Equal(t, err, nil)
			err = dc.conn.Flush()
			require.Equal(t, err, nil)
			err = mockGaea.GetConnWrite().Close()
			require.Equal(t, err, nil)
		}

		fmt.Println(mockGaea.CustomSend(customFunc).ArrivedMsg(mockMariaDB))
	})
```







