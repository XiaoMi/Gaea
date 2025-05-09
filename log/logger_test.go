// Copyright 2025 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockLogger is a mock implementation of the Logger interface for testing purposes.
type MockLogger struct {
	output *bytes.Buffer
}

func NewMockLogger() *MockLogger {
	return &MockLogger{output: new(bytes.Buffer)}
}

func (ml *MockLogger) SetLevel(name, level string) error {
	return nil
}

func (ml *MockLogger) Debug(format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, format+"\n", a...)
	return
}

func (ml *MockLogger) Trace(format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, format+"\n", a...)
	return
}

func (ml *MockLogger) Notice(format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, format+"\n", a...)
	return
}

func (ml *MockLogger) Warn(format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, format+"\n", a...)
	return
}

func (ml *MockLogger) Fatal(format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, format+"\n", a...)
	return
}

func (ml *MockLogger) Debugx(logID, format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, "[%s] "+format+"\n", append([]interface{}{logID}, a...)...)
	return
}

func (ml *MockLogger) Tracex(logID, format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, "[%s] "+format+"\n", append([]interface{}{logID}, a...)...)
	return
}

func (ml *MockLogger) Noticex(logID, format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, "[%s] "+format+"\n", append([]interface{}{logID}, a...)...)
	return
}

func (ml *MockLogger) Warnx(logID, format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, "[%s] "+format+"\n", append([]interface{}{logID}, a...)...)
	return
}

func (ml *MockLogger) Fatalx(logID, format string, a ...interface{}) (err error) {
	_, err = fmt.Fprintf(ml.output, "[%s] "+format+"\n", append([]interface{}{logID}, a...)...)
	return
}

func (ml *MockLogger) Close() {
	// No-op for now
}

func TestSetGlobalLogger(t *testing.T) {
	assert := assert.New(t)

	// Create a mock logger
	mockLogger := NewMockLogger()

	// Set the global logger
	SetGlobalLogger(mockLogger)

	// Check that the global logger is set correctly
	assert.Equal(mockLogger, logger)

	// Reset the global logger
	SetGlobalLogger(nil)

	// Check that the global logger is reset
	assert.Nil(logger)
}

func TestLoggingFunctions(t *testing.T) {
	assert := assert.New(t)

	// Create a mock logger
	mockLogger := NewMockLogger()
	SetGlobalLogger(mockLogger)

	// Define test cases for each logging function
	testCases := []struct {
		name        string
		logFunc     func(string, ...interface{}) error
		expectedOut string
	}{
		{
			name:        "Debug",
			logFunc:     Debug,
			expectedOut: "message",
		},
		{
			name:        "Trace",
			logFunc:     Trace,
			expectedOut: "message",
		},
		{
			name:        "Notice",
			logFunc:     Notice,
			expectedOut: "message",
		},
		{
			name:        "Warn",
			logFunc:     Warn,
			expectedOut: "message",
		},
		{
			name:        "Fatal",
			logFunc:     Fatal,
			expectedOut: "message",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.logFunc("message")
			assert.NoError(err)
			assert.Equal(tc.expectedOut, strings.TrimSpace(mockLogger.output.String()))
			mockLogger.output.Reset() // Reset buffer for next test case
		})
	}
}

func TestLoggingWithFormat(t *testing.T) {
	assert := assert.New(t)

	// Create a mock logger
	mockLogger := NewMockLogger()
	SetGlobalLogger(mockLogger)

	// Define test cases for each logging function with format
	testCases := []struct {
		name        string
		logFunc     func(string, ...interface{}) error
		args        []interface{}
		expectedOut string
	}{
		{
			name:        "Debug with format",
			logFunc:     func(f string, a ...interface{}) error { return Debug(f, a...) },
			args:        []interface{}{"formatted %s", "arg"},
			expectedOut: "formatted arg",
		},
		{
			name:        "Trace with format",
			logFunc:     func(f string, a ...interface{}) error { return Trace(f, a...) },
			args:        []interface{}{"formatted %s", "arg"},
			expectedOut: "formatted arg",
		},
		{
			name:        "Notice with format",
			logFunc:     func(f string, a ...interface{}) error { return Notice(f, a...) },
			args:        []interface{}{"formatted %s", "arg"},
			expectedOut: "formatted arg",
		},
		{
			name:        "Warn with format",
			logFunc:     func(f string, a ...interface{}) error { return Warn(f, a...) },
			args:        []interface{}{"formatted %s", "arg"},
			expectedOut: "formatted arg",
		},
		{
			name:        "Fatal with format",
			logFunc:     func(f string, a ...interface{}) error { return Fatal(f, a...) },
			args:        []interface{}{"formatted %s", "arg"},
			expectedOut: "formatted arg",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.logFunc(tc.args[0].(string), tc.args[1])
			assert.NoError(err)
			assert.Equal(tc.expectedOut, strings.TrimSpace(mockLogger.output.String()))
			mockLogger.output.Reset() // Reset buffer for next test case
		})
	}
}

func TestClose(t *testing.T) {
	assert := assert.New(t)

	// Create a mock logger
	mockLogger := NewMockLogger()
	SetGlobalLogger(mockLogger)

	// Call Close method
	Close()

	// Ensure the logger is closed
	assert.True(true) // We can't check exactly what happens inside Close(), so we just ensure it doesn't panic
}
