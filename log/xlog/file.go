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
	"sync"
	"time"
)

// XFileLog is the file logger
type XFileLog struct {
	filename string
	path     string
	level    int

	skip     int
	runtime  bool
	file     *os.File
	errFile  *os.File
	hostname string
	service  string
	split    sync.Once
	mu       sync.Mutex
}

// constants of XFileLog
const (
	XFileLogDefaultLogID = "900000001"
	SpliterDelay         = 5
	CleanDays            = -3
)

// NewXFileLog is the constructor of XFileLog
//生成一个日志实例，service用来标识业务的服务名。
//比如：logger := xlog.NewXFileLog("gaea")
func NewXFileLog() XLogger {
	return &XFileLog{
		skip: XLogDefSkipNum,
	}
}

// Init implements XLogger
func (p *XFileLog) Init(config map[string]string) (err error) {

	path, ok := config["path"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found path")
		return
	}

	filename, ok := config["filename"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found filename")
		return
	}

	level, ok := config["level"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found level")
		return
	}

	service, _ := config["service"]
	if len(service) > 0 {
		p.service = service
	}

	runtime, ok := config["runtime"]
	if !ok || runtime == "true" || runtime == "TRUE" {
		p.runtime = true
	} else {
		p.runtime = false
	}

	skip, _ := config["skip"]
	if len(skip) > 0 {
		skipNum, err := strconv.Atoi(skip)
		if err == nil {
			p.skip = skipNum
		}
	}

	isDir, err := isDir(path)
	if err != nil || !isDir {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return newError("Mkdir failed, err:%v", err)
		}
	}

	p.path = path
	p.filename = filename
	p.level = LevelFromStr(level)

	hostname, _ := os.Hostname()
	p.hostname = hostname
	body := func() {
		go p.spliter()
	}
	doSplit, ok := config["dosplit"]
	if !ok {
		doSplit = "true"
	}
	if doSplit == "true" {
		p.split.Do(body)
	}
	return p.ReOpen()
}

// split the log file
func (p *XFileLog) spliter() {
	preHour := time.Now().Hour()
	splitTime := time.Now().Format("2006010215")
	defer p.Close()
	for {
		time.Sleep(time.Second * SpliterDelay)
		if time.Now().Hour() != preHour {
			p.clean()
			p.rename(splitTime)
			preHour = time.Now().Hour()
			splitTime = time.Now().Format("2006010215")
		}
	}
}

// SetLevel implements XLogger
func (p *XFileLog) SetLevel(level string) {
	p.level = LevelFromStr(level)
}

// SetSkip implements XLogger
func (p *XFileLog) SetSkip(skip int) {
	p.skip = skip
}

func (p *XFileLog) openFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return nil, newError("open %s failed, err:%v", filename, err)
	}

	return file, err
}

func delayClose(fp *os.File) {
	if fp == nil {
		return
	}
	time.Sleep(1000 * time.Millisecond)
	fp.Close()
}

func (p *XFileLog) clean() (err error) {
	deadline := time.Now().AddDate(0, 0, CleanDays)
	var files []string
	files, err = filepath.Glob(fmt.Sprintf("%s/%s.log*", p.path, p.filename))
	if err != nil {
		return
	}
	var fileInfo os.FileInfo
	for _, file := range files {
		if filepath.Base(file) == fmt.Sprintf("%s.log", p.filename) {
			continue
		}
		if filepath.Base(file) == fmt.Sprintf("%s.log.wf", p.filename) {
			continue
		}
		if fileInfo, err = os.Stat(file); err == nil {
			if fileInfo.ModTime().Before(deadline) {
				os.Remove(file)
			} else if fileInfo.Size() == 0 {
				os.Remove(file)
			}
		}
	}
	return
}

func (p *XFileLog) rename(shuffix string) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	defer p.ReOpen()
	if p.file == nil {
		return
	}
	var fileInfo os.FileInfo
	normalLog := p.path + "/" + p.filename + ".log"
	warnLog := normalLog + ".wf"
	newLog := fmt.Sprintf("%s/%s.log-%s.log", p.path, p.filename, shuffix)
	newWarnLog := fmt.Sprintf("%s/%s.log.wf-%s.log.wf", p.path, p.filename, shuffix)
	if fileInfo, err = os.Stat(normalLog); err == nil && fileInfo.Size() == 0 {
		return
	}
	if _, err = os.Stat(newLog); err == nil {
		return
	}
	if err = os.Rename(normalLog, newLog); err != nil {
		return
	}
	if fileInfo, err = os.Stat(warnLog); err == nil && fileInfo.Size() == 0 {
		return
	}
	if _, err = os.Stat(newWarnLog); err == nil {
		return
	}
	if err = os.Rename(warnLog, newWarnLog); err != nil {
		return
	}
	return
}

// ReOpen implements XLogger
func (p *XFileLog) ReOpen() error {
	go delayClose(p.file)
	go delayClose(p.errFile)

	normalLog := p.path + "/" + p.filename + ".log"
	file, err := p.openFile(normalLog)
	if err != nil {
		return err
	}

	p.file = file
	warnLog := normalLog + ".wf"
	p.errFile, err = p.openFile(warnLog)
	if err != nil {
		p.file.Close()
		p.file = nil
		return err
	}

	return nil
}

// Warn implements XLogger
func (p *XFileLog) Warn(format string, a ...interface{}) error {
	if p.level > WarnLevel {
		return nil
	}

	return p.warnx(XFileLogDefaultLogID, format, a...)
}

// Warnx implements XLogger
func (p *XFileLog) Warnx(logID, format string, a ...interface{}) error {
	if p.level > WarnLevel {
		return nil
	}

	return p.warnx(logID, format, a...)
}

func (p *XFileLog) warnx(logID, format string, a ...interface{}) error {
	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText)

	return p.write(WarnLevel, &logText, logID)
}

// Fatal implements XLogger
func (p *XFileLog) Fatal(format string, a ...interface{}) error {
	if p.level > FatalLevel {
		return nil
	}

	return p.fatalx(XFileLogDefaultLogID, format, a...)
}

// Fatalx implements XLogger
func (p *XFileLog) Fatalx(logID, format string, a ...interface{}) error {
	if p.level > FatalLevel {
		return nil
	}

	return p.fatalx(logID, format, a...)
}

func (p *XFileLog) fatalx(logID, format string, a ...interface{}) error {
	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText)

	return p.write(FatalLevel, &logText, logID)
}

// Notice implements XLogger
func (p *XFileLog) Notice(format string, a ...interface{}) error {
	if p.level > NoticeLevel {
		return nil
	}
	return p.noticex(XFileLogDefaultLogID, format, a...)
}

// Noticex implements XLogger
func (p *XFileLog) Noticex(logID, format string, a ...interface{}) error {
	if p.level > NoticeLevel {
		return nil
	}
	return p.noticex(logID, format, a...)
}

func (p *XFileLog) noticex(logID, format string, a ...interface{}) error {
	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)

	return p.write(NoticeLevel, &logText, logID)
}

// Trace implements XLogger
func (p *XFileLog) Trace(format string, a ...interface{}) error {
	return p.tracex(XFileLogDefaultLogID, format, a...)
}

// Tracex implements XLogger
func (p *XFileLog) Tracex(logID, format string, a ...interface{}) error {
	return p.tracex(logID, format, a...)
}

func (p *XFileLog) tracex(logID, format string, a ...interface{}) error {
	if p.level > TraceLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(filename), lineno, logText)

	return p.write(TraceLevel, &logText, logID)
}

// Debug implements XLogger
func (p *XFileLog) Debug(format string, a ...interface{}) error {
	return p.debugx(XFileLogDefaultLogID, format, a...)
}

func (p *XFileLog) debugx(logID, format string, a ...interface{}) error {
	if p.level > DebugLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)

	return p.write(DebugLevel, &logText, logID)
}

// Debugx implements XLogger
func (p *XFileLog) Debugx(logID, format string, a ...interface{}) error {
	return p.debugx(logID, format, a...)
}

// Close implements XLogger
func (p *XFileLog) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.file != nil {
		p.file.Close()
		p.file = nil
	}

	if p.errFile != nil {
		p.errFile.Close()
		p.errFile = nil
	}
}

// GetHost getter of hostname
func (p *XFileLog) GetHost() string {
	return p.hostname
}

func (p *XFileLog) write(level int, msg *string, logID string) error {
	levelText := levelTextArray[level]
	time := time.Now().Format("2006-01-02 15:04:05")

	logText := formatLog(msg, time, p.service, p.hostname, levelText, logID)
	file := p.file
	if level >= WarnLevel {
		file = p.errFile
	}

	file.Write([]byte(logText))
	return nil
}

func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}
