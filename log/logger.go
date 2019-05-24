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

package log

import (
	"fmt"

	"github.com/XiaoMi/Gaea/log/xlog"
)

var logger Logger

// Logger is the log interface
type Logger interface {
	SetLevel(name, level string) error

	Debug(format string, a ...interface{}) (err error)
	Trace(format string, a ...interface{}) (err error)
	Notice(format string, a ...interface{}) (err error)
	Warn(format string, a ...interface{}) (err error)
	Fatal(format string, a ...interface{}) (err error)

	Debugx(logID, format string, a ...interface{}) (err error)
	Tracex(logID, format string, a ...interface{}) (err error)
	Noticex(logID, format string, a ...interface{}) (err error)
	Warnx(logID, format string, a ...interface{}) (err error)
	Fatalx(logID, format string, a ...interface{}) (err error)

	Close()
}

func init() {
	cfg := make(map[string]string)
	cfg["level"] = "debug"
	lg, err := xlog.CreateLogManager("console", cfg)
	if err != nil {
		panic(fmt.Errorf("init global logger error: %v", err))
	}
	logger = lg
}

// SetGlobalLogger set global logger, if global logger already exists, close it and set to new.
func SetGlobalLogger(lg Logger) {
	if logger != nil {
		logger.Close()
	}
	logger = lg
}

// Debug log debug message.
func Debug(format string, a ...interface{}) (err error) {
	return logger.Debug(format, a...)
}

// Debugx log debug message with logID.
func Debugx(logID, format string, a ...interface{}) (err error) {
	return logger.Debugx(logID, format, a...)
}

// Trace log trace message.
func Trace(format string, a ...interface{}) (err error) {
	return logger.Trace(format, a...)
}

// Tracex log trace message with logID.
func Tracex(logID, format string, a ...interface{}) (err error) {
	return logger.Tracex(logID, format, a...)
}

// Notice log notice message.
func Notice(format string, a ...interface{}) (err error) {
	return logger.Notice(format, a...)
}

// Noticex log notice message with logID.
func Noticex(logID, format string, a ...interface{}) (err error) {
	return logger.Noticex(logID, format, a...)
}

// Warn log warn message.
func Warn(format string, a ...interface{}) (err error) {
	return logger.Warn(format, a...)
}

// Warnx log warn message with logID.
func Warnx(logID, format string, a ...interface{}) (err error) {
	return logger.Warnx(logID, format, a...)
}

// Fatal log fatal message.
func Fatal(format string, a ...interface{}) (err error) {
	return logger.Fatal(format, a...)
}

// Fatalx log fatal message with logID.
func Fatalx(logID, format string, a ...interface{}) (err error) {
	return logger.Fatalx(logID, format, a...)
}

// Close close the global logger
func Close() {
	logger.Close()
}
