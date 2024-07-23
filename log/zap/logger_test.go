package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"testing"
)

func BenchmarkSyncLoggerWriter(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	encoder := &ZapEncoder{}
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(f), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return true
		})),
	)
	l := zap.New(core)
	g := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		g.Add(1)
		go func() {
			defer g.Done()
			l.Info("ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)")
		}()
	}
	g.Wait()
	l.Sync()
}

func BenchmarkAsyncLoggerWriter(b *testing.B) {
	f, _ := os.OpenFile("/dev/null", os.O_RDWR|os.O_CREATE, 0666)
	encoder := &ZapEncoder{}
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, NewAsyncWriter(f), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return true
		})),
	)
	l := zap.New(core)
	g := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		g.Add(1)
		go func() {
			defer g.Done()
			l.Info("ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)")
		}()
	}
	g.Wait()
	l.Sync()
}
