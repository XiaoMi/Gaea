package xlog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestNewXConsoleLog(t *testing.T) {
	logger := NewXConsoleLog()
	if logger == nil {
		t.Error("Expected logger to be non-nil")
	}
}

func TestInit(t *testing.T) {
	config := map[string]string{
		"level": "info",
	}
	logger := &XConsoleLog{}
	err := logger.Init(config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if logger.level != NoticeLevel {
		t.Errorf("Expected level to be %d, but got %d", NoticeLevel, logger.level)
	}
}

func TestSetLevel(t *testing.T) {
	logger := &XConsoleLog{}
	logger.SetLevel("fatal")
	if logger.level != FatalLevel {
		t.Errorf("Expected level to be %d, but got %d", FatalLevel, logger.level)
	}
}

func TestSetSkip(t *testing.T) {
	logger := &XConsoleLog{}
	logger.SetSkip(2)
	if logger.skip != 2 {
		t.Errorf("Expected skip to be 2, but got %d", logger.skip)
	}
}

func TestConoleReOpen(t *testing.T) {
	logger := &XConsoleLog{}
	err := logger.ReOpen()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestConsoleWarn(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	defer func() {
		os.Stdout = old
		os.Stderr = oldErr
	}()
	outC := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(r)
		outC <- string(data)
	}()

	logger := &XConsoleLog{}
	logger.Warn("test warning message")
	os.Stdout.Close()
	buf.WriteString(<-outC)
	expected := "test warning message"
	if !strings.Contains(strings.TrimSpace(buf.String()), strings.TrimSpace(expected)) {
		t.Errorf("Expected output '%s', but got '%s'", expected, buf.String())
	}
}

func TestConsoleFatal(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	defer func() {
		os.Stdout = old
		os.Stderr = oldErr
	}()
	outC := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(r)
		outC <- string(data)
	}()

	logger := &XConsoleLog{}
	logger.Fatal("test fatal message")
	os.Stdout.Close()
	buf.WriteString(<-outC)
	expected := "test fatal message"
	if !strings.Contains(strings.TrimSpace(buf.String()), strings.TrimSpace(expected)) {
		t.Errorf("Expected output '%s', but got '%s'", expected, buf.String())
	}
}

func TestConsoleNotice(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	defer func() {
		os.Stdout = old
		os.Stderr = oldErr
	}()
	outC := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(r)
		outC <- string(data)
	}()

	logger := &XConsoleLog{}
	logger.Notice("test notice message")
	os.Stdout.Close()
	buf.WriteString(<-outC)
	expected := "test notice message"
	if !strings.Contains(strings.TrimSpace(buf.String()), strings.TrimSpace(expected)) {
		t.Errorf("Expected output '%s', but got '%s'", expected, buf.String())
	}
}

func TestConsoleTrace(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	defer func() {
		os.Stdout = old
		os.Stderr = oldErr
	}()
	outC := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(r)
		outC <- string(data)
	}()

	logger := &XConsoleLog{}
	logger.Trace("test trace message")
	os.Stdout.Close()
	buf.WriteString(<-outC)
	expected := "test trace message"
	if !strings.Contains(strings.TrimSpace(buf.String()), strings.TrimSpace(expected)) {
		t.Errorf("Expected output '%s', but got '%s'", expected, buf.String())
	}
}

func TestConsoleDebug(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	defer func() {
		os.Stdout = old
		os.Stderr = oldErr
	}()
	outC := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(r)
		outC <- string(data)
	}()

	logger := &XConsoleLog{}
	logger.Debug("test debug message")
	os.Stdout.Close()
	buf.WriteString(<-outC)
	expected := "test debug message"
	if !strings.Contains(strings.TrimSpace(buf.String()), strings.TrimSpace(expected)) {
		t.Errorf("Expected output '%s', but got '%s'", expected, buf.String())
	}
}

func TestGetHost(t *testing.T) {
	logger := &XConsoleLog{}
	logger.hostname = "localhost"
	if logger.GetHost() != "localhost" {
		t.Errorf("Expected host to be 'localhost', but got '%s'", logger.GetHost())
	}
}

func TestNoticeWrite(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	defer func() {
		os.Stdout = old
		os.Stderr = oldErr
	}()
	outC := make(chan string)
	go func() {
		data, _ := ioutil.ReadAll(r)
		outC <- string(data)
	}()

	logger := &XConsoleLog{}
	logger.write(NoticeLevel, new(string), "")
	expected := ""
	if !strings.Contains(strings.TrimSpace(buf.String()), strings.TrimSpace(expected)) {
		t.Errorf("Expected output '%s', but got '%s'", expected, buf.String())
	}
}

func TestMain(m *testing.M) {
	// Initialize any necessary resources here if needed
	fmt.Println("Running tests...")
	result := m.Run()
	// Clean up any resources after all tests have run
	os.Exit(result)
}
