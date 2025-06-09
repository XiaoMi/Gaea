package xlog

import (
	"os"
	"sync"
	"testing"
)

func BenchmarkXLoggerWriterWithWaitGroup(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	l := &XFileLog{}
	l.file = f
	g := sync.WaitGroup{}
	logMsg := "ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		g.Add(1)
		go func() {
			defer g.Done()
			l.Notice(logMsg)
		}()
	}
	g.Wait()
	b.ReportAllocs()
}

func BenchmarkXLogger(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	l := &XFileLog{}
	l.file = f
	logMsg := "ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)"
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Notice(logMsg)
		}
	})

	b.ReportAllocs() // 必须放在测试逻辑之后
}
