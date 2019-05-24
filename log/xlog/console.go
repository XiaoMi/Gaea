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

package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// constants of XConsoleLog
const (
	XConsoleLogDefaultlogID = "800000001"
)

// XConsoleLog is the console logger
type XConsoleLog struct {
	level int

	skip     int
	hostname string
	service  string
}

// Brush is pretty string
type Brush func(string) string

// NewBrush is the constructor of Brush
func NewBrush(color string) Brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;37"), // white
	NewBrush("1;36"), // debug			cyan
	NewBrush("1;35"), // trace   magenta
	NewBrush("1;31"), // notice      red
	NewBrush("1;33"), // warn    yellow
	NewBrush("1;32"), // fatal			green
	NewBrush("1;34"), //
	NewBrush("1;34"), //
}

// NewXConsoleLog is the constructor of XConsoleLog
func NewXConsoleLog() XLogger {
	return &XConsoleLog{
		skip: XLogDefSkipNum,
	}
}

// Init init XConsoleLog
func (p *XConsoleLog) Init(config map[string]string) (err error) {
	level, ok := config["level"]
	if !ok {
		err = fmt.Errorf("init XConsoleLog failed, not found level")
		return
	}

	service, _ := config["service"]
	if len(service) > 0 {
		p.service = service
	}

	skip, _ := config["skip"]
	if len(skip) > 0 {
		skipNum, err := strconv.Atoi(skip)
		if err == nil {
			p.skip = skipNum
		}
	}

	p.level = LevelFromStr(level)
	hostname, _ := os.Hostname()
	p.hostname = hostname

	return
}

// SetLevel implements XLogger
func (p *XConsoleLog) SetLevel(level string) {
	p.level = LevelFromStr(level)
}

// SetSkip implements XLogger
func (p *XConsoleLog) SetSkip(skip int) {
	p.skip = skip
}

// ReOpen implements XLogger
func (p *XConsoleLog) ReOpen() error {
	return nil
}

// Warn implements XLogger
func (p *XConsoleLog) Warn(format string, a ...interface{}) error {
	return p.warnx(XConsoleLogDefaultlogID, format, a...)
}

// Warnx implements XLogger
func (p *XConsoleLog) Warnx(logID, format string, a ...interface{}) error {
	return p.warnx(logID, format, a...)
}

func (p *XConsoleLog) warnx(logID, format string, a ...interface{}) error {
	if p.level > WarnLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)

	color := colors[WarnLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	return p.write(WarnLevel, &logText, logID)
}

// Fatal implements XLogger
func (p *XConsoleLog) Fatal(format string, a ...interface{}) error {
	return p.fatalx(XConsoleLogDefaultlogID, format, a...)
}

// Fatalx implements XLogger
func (p *XConsoleLog) Fatalx(logID, format string, a ...interface{}) error {
	return p.fatalx(logID, format, a...)
}

func (p *XConsoleLog) fatalx(logID, format string, a ...interface{}) error {
	if p.level > FatalLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)

	color := colors[FatalLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	return p.write(FatalLevel, &logText, logID)
}

// Notice implements XLogger
func (p *XConsoleLog) Notice(format string, a ...interface{}) error {
	return p.Noticex(XConsoleLogDefaultlogID, format, a...)
}

// Noticex implements XLogger
func (p *XConsoleLog) Noticex(logID, format string, a ...interface{}) error {
	if p.level > NoticeLevel {
		return nil
	}

	logText := formatValue(format, a...)
	return p.write(NoticeLevel, &logText, logID)
}

// Trace implements XLogger
func (p *XConsoleLog) Trace(format string, a ...interface{}) error {
	return p.tracex(XConsoleLogDefaultlogID, format, a...)
}

// Tracex implements XLogger
func (p *XConsoleLog) Tracex(logID, format string, a ...interface{}) error {
	return p.tracex(logID, format, a...)
}

func (p *XConsoleLog) tracex(logID, format string, a ...interface{}) error {
	if p.level > TraceLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)

	color := colors[TraceLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	return p.write(TraceLevel, &logText, logID)
}

// Debug implements XLogger
func (p *XConsoleLog) Debug(format string, a ...interface{}) error {
	return p.debugx(XConsoleLogDefaultlogID, format, a...)
}

// Debugx implements XLogger
func (p *XConsoleLog) Debugx(logID, format string, a ...interface{}) error {
	return p.debugx(logID, format, a...)
}

func (p *XConsoleLog) debugx(logID, format string, a ...interface{}) error {
	if p.level > DebugLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)

	color := colors[DebugLevel]
	logText = color(fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText))

	return p.write(DebugLevel, &logText, logID)
}

// Close implements XLogger
//关闭日志库。注意：如果没有调用Close()关闭日志库的话，将会造成文件句柄泄露
func (p *XConsoleLog) Close() {
}

// GetHost getter of hostname
func (p *XConsoleLog) GetHost() string {
	return p.hostname
}

func (p *XConsoleLog) write(level int, msg *string, logID string) error {
	color := colors[level]
	levelText := color(levelTextArray[level])
	time := time.Now().Format("2006-01-02 15:04:05")

	logText := formatLog(msg, time, p.service, p.hostname, levelText, logID)
	file := os.Stdout
	if level >= WarnLevel {
		file = os.Stderr
	}

	file.Write([]byte(logText))
	return nil
}
