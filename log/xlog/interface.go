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

// XLogger declares method that log instance should implement.
type XLogger interface {
	// init logger
	Init(config map[string]string) error

	// reopen logger
	ReOpen() error

	//设置日志级别, 级别如下: "Debug", "Trace", "Notice", "Warn", "Fatal", "None"
	SetLevel(level string)

	// set skip
	SetSkip(skip int)

	// 打印Debug日志. 当日志级别大于Debug时, 不会输出任何日志
	Debug(format string, a ...interface{}) error

	// 打印Trace日志. 当日志级别大于Trace时, 不会输出任何日志
	Trace(format string, a ...interface{}) error

	// 打印Notice日志. 当日志级别大于Notice时, 不会输出任何日志
	Notice(format string, a ...interface{}) error

	// 打印Warn日志. 当日志级别大于Warn时, 不会输出任何日志
	Warn(format string, a ...interface{}) error

	// 打印Fatal日志. 当日志级别大于Fatal时, 不会输出任何日志
	Fatal(format string, a ...interface{}) error

	// 打印Debug日志, 需要传入logID. 当日志级别大于Debug时, 不会输出任何日志
	Debugx(logID, format string, a ...interface{}) error

	// 打印Trace日志, 需要传入logID. 当日志级别大于Trace时, 不会输出任何日志
	Tracex(logID, format string, a ...interface{}) error

	// 打印Notice日志, 需要传入logID. 当日志级别大于Notice时, 不会输出任何日志
	Noticex(logID, format string, a ...interface{}) error

	// 打印Warn日志, 需要传入logID. 当日志级别大于Warn时, 不会输出任何日志
	Warnx(logID, format string, a ...interface{}) error

	// 打印Fatal日志, 需要传入logID. 当日志级别大于Fatal时, 不会输出任何日志
	Fatalx(logID, format string, a ...interface{}) error

	// 关闭日志库. 注意: 如果没有调用Close()关闭日志库的话, 将会造成文件句柄泄露
	Close()
}
