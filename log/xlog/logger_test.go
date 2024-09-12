package xlog

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

const TEST_FILE = "/Users/yangming/project/gaea.xiaomi/log/xlog/test.log"

func BenchmarkXLoggerWriter(b *testing.B) {
	f, _ := os.OpenFile(TEST_FILE, os.O_RDWR|os.O_CREATE, 0666)
	l := &XFileLog{}
	l.file = f
	g := sync.WaitGroup{}
	for n := 0; n < b.N; n++ {
		g.Add(1)
		go func() {
			defer g.Done()
			l.Notice("ns=test_namespace_1, root@127.0.0.1:61855->10.38.164.125:3308/, mysql_connect_id=1637760|select sleep(3)")
		}()
	}
	g.Wait()
}

// MockLogger 是 XLogger 的一个模拟实现
type MockLogger struct {
	lock sync.RWMutex
}

func newMockLogger() XLogger {
	return &MockLogger{
		lock: sync.RWMutex{},
	}
}

func (m *MockLogger) Init(config map[string]string) error { return nil }
func (m *MockLogger) ReOpen() error                       { return nil }
func (m *MockLogger) SetLevel(level string)               {}
func (m *MockLogger) Warn(format string, a ...interface{}) error {
	fmt.Println("Warn:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Fatal(format string, a ...interface{}) error {
	fmt.Println("Fatal:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Notice(format string, a ...interface{}) error {
	fmt.Println("Notice:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Trace(format string, a ...interface{}) error {
	fmt.Println("Trace:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Debug(format string, a ...interface{}) error {
	fmt.Println("Debug:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Warnx(logID, format string, a ...interface{}) error {
	fmt.Println("Warnx:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Fatalx(logID, format string, a ...interface{}) error {
	fmt.Println("Fatalx:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Noticex(logID, format string, a ...interface{}) error {
	fmt.Println("Noticex:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Tracex(logID, format string, a ...interface{}) error {
	fmt.Println("Tracex:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Debugx(logID, format string, a ...interface{}) error {
	fmt.Println("Debugx:", fmt.Sprintf(format, a...))
	return nil
}
func (m *MockLogger) Close()      {}
func (m *MockLogger) SetSkip(int) {}

func newMockConfig() map[string]string {
	return map[string]string{
		"level": "notice",
	}
}

func TestCreateLogManager(t *testing.T) {
	mgr, err := CreateLogManager("console", newMockConfig())
	if err != nil {
		t.Fatalf("Failed to create log manager: %v", err)
	}
	if mgr == nil {
		t.Error("Expected non-nil LogManager")
	}
}

func TestRegisterLogger(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	err := mgr.RegisterLogger("testLogger", newMockLogger())
	if err != nil {
		t.Fatalf("Failed to register logger: %v", err)
	}
}

func TestEnableLogger(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	err := mgr.EnableLogger("testLogger", true)
	if err != nil {
		t.Fatalf("Failed to enable logger: %v", err)
	}
}

func TestGetLogger(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	logger, err := mgr.GetLogger("testLogger")
	if err != nil {
		t.Fatalf("Failed to get logger: %v", err)
	}
	if logger == nil {
		t.Error("Expected non-nil Logger")
	}
}

func TestReOpen(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	err := mgr.EnableLogger("testLogger", true)
	if err != nil {
		t.Fatalf("Failed to enable logger: %v", err)
	}
	err = mgr.ReOpen()
	if err != nil && err.Error() != "" {
		t.Fatalf("Failed to re-open logger: %v", err)
	}
}

func TestSetLevelAll(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	mgr.EnableLogger("testLogger", true)
	mgr.SetLevelAll("Debug")
}

func TestWarn(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	mgr.EnableLogger("testLogger", true)
	err := mgr.Warn("Test warn message")
	if err != nil {
		t.Fatalf("Failed to call Warn: %v", err)
	}
}

func TestFatal(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	mgr.EnableLogger("testLogger", true)
	err := mgr.Fatal("Test fatal message")
	if err != nil {
		t.Fatalf("Failed to call Fatal: %v", err)
	}
}

func TestNotice(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	mgr.EnableLogger("testLogger", true)
	err := mgr.Notice("Test notice message")
	if err != nil {
		t.Fatalf("Failed to call Notice: %v", err)
	}
}

func TestTrace(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	mgr.EnableLogger("testLogger", true)
	err := mgr.Trace("Test trace message")
	if err != nil {
		t.Fatalf("Failed to call Trace: %v", err)
	}
}

func TestDebug(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	mgr.EnableLogger("testLogger", true)
	err := mgr.Debug("Test debug message")
	if err != nil {
		t.Fatalf("Failed to call Debug: %v", err)
	}
}

func TestClose(t *testing.T) {
	mgr, _ := CreateLogManager("console", newMockConfig())
	mgr.RegisterLogger("testLogger", newMockLogger())
	mgr.EnableLogger("testLogger", true)
	mgr.Close()
}
