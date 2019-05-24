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

/*
小米网golang日志库

小米网golang日志库支持6种日志级别：

 1）Debug
 2）Trace
 3）Notice
 4）Warn
 5）Fatal
 6）None

支持两种输出格式：

 1）json格式
 2）自定义输出

日志级别优先级：

 Debug < Trace < Notice < Warn < Fatal < None

即如果定义日志级别为Debug：则Trace、Notice、Warn、Fatal等级别的日志都会输出；反之，如果定义日志级别为Trace，则Debug不会输出，其它级别的日志都会输出。当日志级别为None时，将不会有任何日志输出。

*/

package xlog

import (
	"errors"
	"fmt"
	"sync"
)

// LogInstance wraps a XLogger
type LogInstance struct {
	logger  XLogger
	enable  bool
	initial bool
	source  string
}

// LogManager is the manager that hold different kind of xlog loggers
type LogManager struct {
	loggers map[string]*LogInstance
	lock    sync.RWMutex
}

// CreateLogManager create log manager from configs.
func CreateLogManager(name string, config map[string]string, source ...string) (*LogManager, error) {
	mgr := &LogManager{
		loggers: make(map[string]*LogInstance),
	}

	if err := mgr.Init(name, config, source...); err != nil {
		return nil, err
	}

	return mgr, nil
}

// Init LogManager with config. The name can be set to "console" or "file".
// Please reference to XConsoleLog and XFileLog.
func (l *LogManager) Init(name string, config map[string]string, source ...string) (err error) {
	if err := l.RegisterLogger("console", NewXConsoleLog()); err != nil {
		return err
	}

	if err := l.RegisterLogger("file", NewXFileLog()); err != nil {
		return err
	}

	err = l.initLogger(name, config)
	//关闭自动注入的logger
	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}

		if v.source == "auto" {
			v.enable = false
		}
	}

	return
}

func (l *LogManager) initLogger(name string, config map[string]string, source ...string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	instance, ok := l.loggers[name]
	if !ok {
		err = fmt.Errorf("not found logger[%s]", name)
		return
	}

	err = instance.logger.Init(config)
	if err != nil {
		return
	}

	if len(source) > 0 {
		instance.source = source[0]
	} else {
		instance.source = ""
	}

	instance.enable = true
	instance.initial = true

	return
}

// RegisterLogger register a logger
func (l *LogManager) RegisterLogger(name string, logger XLogger) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	_, ok := l.loggers[name]
	if ok {
		err = fmt.Errorf("duplicate logger[%s]", name)
		return
	}

	l.loggers[name] = &LogInstance{
		logger:  logger,
		enable:  false,
		initial: false,
	}

	return
}

// UnregisterLogger unregister a logger
func (l *LogManager) UnregisterLogger(name string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	v, ok := l.loggers[name]
	if !ok {
		err = fmt.Errorf("not found logger[%s]", name)
		return
	}

	if v != nil {
		v.logger.Close()
	}

	delete(l.loggers, name)
	return
}

// EnableLogger enable or disable a logger
func (l *LogManager) EnableLogger(name string, enable bool) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	instance, ok := l.loggers[name]
	if !ok {
		err = fmt.Errorf("not found logger[%s]", name)
		return
	}

	if !instance.initial {
		instance.enable = false
		return
	}

	instance.enable = enable
	return
}

// GetLogger get logger by name
func (l *LogManager) GetLogger(name string) (logger XLogger, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	instance, ok := l.loggers[name]
	if !ok {
		err = fmt.Errorf("not found logger[%s]", name)
		return
	}

	logger = instance.logger
	return
}

// ReOpen reopen all enabled loggers
func (l *LogManager) ReOpen() (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	var errorMsg string
	for k, v := range l.loggers {

		if v.logger == nil || !v.enable {
			continue
		}

		errRet := v.logger.ReOpen()
		if errRet != nil {
			errorMsg += fmt.Sprintf("logger[%s] reload failed, err[%v]\n", k, errRet)
			continue
		}
	}

	err = errors.New(errorMsg)
	return
}

// SetLevelAll set level to all loggers
func (l *LogManager) SetLevelAll(level string) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if !v.enable || v.logger == nil {
			continue
		}

		v.logger.SetLevel(level)
	}
}

// SetLevel implements XLogger
func (l *LogManager) SetLevel(name, level string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	v, ok := l.loggers[name]
	if !ok || v.logger == nil {
		err = fmt.Errorf("not found logger[%s]", name)
		return
	}

	v.logger.SetLevel(level)
	return
}

// Warn implements XLogger
func (l *LogManager) Warn(format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Warn(format, a...)
	}

	return
}

// Fatal implements XLogger
func (l *LogManager) Fatal(format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Fatal(format, a...)
	}

	return
}

// Notice implements XLogger
func (l *LogManager) Notice(format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Notice(format, a...)
	}

	return
}

// Trace implements XLogger
func (l *LogManager) Trace(format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Trace(format, a...)
	}

	return
}

// Debug implements XLogger
func (l *LogManager) Debug(format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Debug(format, a...)
	}

	return
}

// Warnx implements XLogger
func (l *LogManager) Warnx(logID, format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Warnx(logID, format, a...)
	}

	return
}

// Fatalx implements XLogger
func (l *LogManager) Fatalx(logID, format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Fatalx(logID, format, a...)
	}

	return
}

// Noticex implements XLogger
func (l *LogManager) Noticex(logID, format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Noticex(logID, format, a...)
	}

	return
}

// Tracex implements XLogger
func (l *LogManager) Tracex(logID, format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Tracex(logID, format, a...)
	}

	return
}

// Debugx implements XLogger
func (l *LogManager) Debugx(logID, format string, a ...interface{}) (err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Debugx(logID, format, a...)
	}

	return
}

// Close implements XLogger
func (l *LogManager) Close() {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, v := range l.loggers {
		if v.logger == nil || !v.enable {
			continue
		}
		v.logger.Close()
	}

	return
}
