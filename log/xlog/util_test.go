package xlog

import (
	"strings"
	"testing"
)

func TestLevelFromStr(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"debug", DebugLevel},
		{"trace", TraceLevel},
		{"notice", NoticeLevel},
		{"warn", WarnLevel},
		{"fatal", FatalLevel},
		{"none", NoneLevel},
		{"unknown", NoticeLevel}, // unknown level should default to NOTICE
	}

	for _, test := range tests {
		result := LevelFromStr(test.input)
		if result != test.expected {
			t.Errorf("LevelFromStr(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestGetRuntimeInfo(t *testing.T) {
	_, _, lineno := getRuntimeInfo(XLogDefSkipNum)
	if lineno < 0 {
		t.Errorf("getRuntimeInfo returned invalid line number: %d", lineno)
	}
}

func TestFormatLog(t *testing.T) {
	body := "Test log message"
	fields := []string{"field1", "field2"}
	logMessage := formatLog(&body, fields[0], fields[1])
	expected := "[field1] [field2] Test log message\n"
	if logMessage != expected {
		t.Errorf("formatLog returned %q; want %q", logMessage, expected)
	}
}

func TestFormatValue(t *testing.T) {
	testCases := []struct {
		format   string
		args     []interface{}
		expected string
	}{
		{"Hello, %s!", []interface{}{"World"}, "Hello, World!"},
		{"Number: %d", []interface{}{42}, "Number: 42"},
		{"No args", nil, "No args"},
	}

	for _, tc := range testCases {
		result := formatValue(tc.format, tc.args...)
		if result != tc.expected {
			t.Errorf("formatValue(%q, %v) = %q; want %q", tc.format, tc.args, result, tc.expected)
		}
	}
}

func TestFormatLineInfo(t *testing.T) {
	functionName := "main.main"
	filename := "main.go"
	logText := "This is a test log"
	lineno := 42
	formatted := formatLineInfo(true, functionName, filename, logText, lineno)
	expected := "[main/main.go:42] This is a test log"
	if formatted != expected {
		t.Errorf("formatLineInfo returned %q; want %q", formatted, expected)
	}
}

func TestNewError(t *testing.T) {
	err := newError("An error occurred")
	expectedMsg := "An error occurred"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("newError returned %q; want %q", err.Error(), expectedMsg)
	}
}
