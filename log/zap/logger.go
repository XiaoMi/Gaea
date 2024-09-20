// Copyright 2024 The Gaea Authors. All Rights Reserved.
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

package zap

import (
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ZapLogDefaultlogID = "900000001"
)

type ZapLoggerManager struct {
	logger  *zap.Logger
	writers []io.WriteCloser
}

// CreateLogManager create log manager from configs.
func CreateLogManager(config map[string]string) (*ZapLoggerManager, error) {
	logDir, ok := config["path"]
	if !ok {
		return nil, fmt.Errorf("init XFileLog failed, not found path")
	}
	filename, ok := config["filename"]
	if !ok {
		return nil, fmt.Errorf("init XFileLog failed, not found filename")
	}

	level, ok := config["level"]
	if !ok {
		return nil, fmt.Errorf("init XFileLog failed, not found level")
	}
	logKeepDays := 0
	if value, ok := config["log_keep_days"]; ok {
		logKeepDays, _ = strconv.Atoi(value)
	}
	logKeepCounts := 0
	if value, ok := config["log_keep_counts"]; ok {
		logKeepCounts, _ = strconv.Atoi(value)
	}

	encoder := &ZapEncoder{}

	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= getZapLevelFromStr(level)
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= getZapLevelFromStr(level)
	})
	logFile := path.Join(logDir, filename+".log")

	// 获取 info、warn日志文件的 io.WriteCloser 抽象 getWriter() 在下方实现
	infoWriter := NewAsyncWriter(getInfoWriter(logFile, logKeepDays, logKeepCounts))
	warnWriter := NewAsyncWriter(getWarnWriter(logFile, logKeepDays, logKeepCounts))

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoWriter, infoLevel),
		zapcore.NewCore(encoder, warnWriter, warnLevel),
	)
	l := zap.New(core)
	return &ZapLoggerManager{
		logger:  l,
		writers: []io.WriteCloser{infoWriter, warnWriter},
	}, nil
}

func getLogWriter(baseName string, filename string, logKeepDays, logKeepCounts int) io.WriteCloser {
	// rotatelogs 不允许全部设置 > 0
	if logKeepDays > 0 && logKeepCounts > 0 {
		if logKeepDays*24 > logKeepCounts {
			logKeepDays = 0
		} else {
			logKeepCounts = 0
		}
	}
	hook, err := rotatelogs.New(
		filename,
		rotatelogs.WithLinkName(baseName),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(logKeepDays)),
		rotatelogs.WithRotationCount(uint(logKeepCounts)),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func getInfoWriter(filename string, logKeepDays, logKeepCounts int) io.WriteCloser {
	return getLogWriter(filename, filename+"-%Y%m%d%H.log", logKeepDays, logKeepCounts)
}

func getWarnWriter(filename string, logKeepDays, logKeepCounts int) io.WriteCloser {
	return getLogWriter(filename+".wf", filename+".wf-%Y%m%d%H.log.wf", logKeepDays, logKeepCounts)
}

// LevelFromStr get log level from level string
func getZapLevelFromStr(level string) zapcore.Level {
	resultLevel := zap.DebugLevel
	levelLower := strings.ToLower(level)
	switch levelLower {
	case "debug":
		resultLevel = zap.DebugLevel
	case "trace":
		resultLevel = zap.DebugLevel
	case "notice":
		resultLevel = zap.InfoLevel
	case "warn":
		resultLevel = zap.WarnLevel
	case "fatal":
		resultLevel = zap.FatalLevel
	case "none":
		resultLevel = 99
	default:
		resultLevel = zap.InfoLevel
	}
	return resultLevel
}

// SetLevel implements XLogger
func (l *ZapLoggerManager) SetLevel(name, level string) (err error) {
	return nil
}

// Warn implements XLogger
func (l *ZapLoggerManager) Notice(format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.Noticex(ZapLogDefaultlogID, format, a...)
	return
}

// Warn implements XLogger
func (l *ZapLoggerManager) Warn(format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.Warnx(ZapLogDefaultlogID, format, a...)
	return
}

// Fatal implements XLogger
func (l *ZapLoggerManager) Fatal(format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.Fatalx(ZapLogDefaultlogID, format, a...)
	return
}

// Trace implements XLogger
func (l *ZapLoggerManager) Trace(format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.Tracex(ZapLogDefaultlogID, format, a...)
	return
}

// Debug implements XLogger
func (l *ZapLoggerManager) Debug(format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.Debugx(ZapLogDefaultlogID, format, a...)
	return
}

// Warnx implements XLogger
func (l *ZapLoggerManager) Warnx(logID, format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.logger.Warn("[" + logID + "] " + fmt.Sprintf(format, a...))
	return
}

// Fatalx implements XLogger, 不使用 Fatal，会导致进程退出
func (l *ZapLoggerManager) Fatalx(logID, format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.logger.Error("[" + logID + "] " + fmt.Sprintf(format, a...))
	return
}

// Noticex implements XLogger
func (l *ZapLoggerManager) Noticex(logID, format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.logger.Info("[" + logID + "] " + fmt.Sprintf(format, a...))
	return
}

// Tracex implements XLogger
func (l *ZapLoggerManager) Tracex(logID, format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.logger.Debug("[" + logID + "] " + fmt.Sprintf(format, a...))
	return
}

// Debugx implements XLogger
func (l *ZapLoggerManager) Debugx(logID, format string, a ...interface{}) (err error) {
	if l.logger == nil {
		return
	}
	l.logger.Debug("[" + logID + "] " + fmt.Sprintf(format, a...))
	return
}

// Close implements XLogger
func (l *ZapLoggerManager) Close() {
	if l.logger == nil {
		return
	}
	l.logger.Sync()
	for _, writer := range l.writers {
		writer.Close()
	}
	l.logger = nil
}
