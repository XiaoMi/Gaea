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

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/XiaoMi/Gaea/core"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/log/xlog"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/proxy/server"
)

var configFile = flag.String("config", "etc/gaea.ini", "gaea config file")
var info = flag.Bool("info", false, "show info of gaea")

func main() {
	flag.Parse()
	if *info {
		fmt.Printf("Build Version Information:%s\n", core.Info.LongForm())
		return
	}

	fmt.Printf("Build Version Information:%s\n", core.Info.LongForm())

	// init config of gaea proxy
	cfg, err := models.ParseProxyConfigFromFile(*configFile)
	if err != nil {
		fmt.Printf("parse config file error:%v\n", err.Error())
		return
	}

	if err = initXLog(cfg.LogOutput, cfg.LogPath, cfg.LogFileName, cfg.LogLevel, cfg.Service); err != nil {
		fmt.Printf("init xlog error: %v\n", err.Error())
		return
	}
	defer log.Close()

	// init manager
	mgr, err := server.LoadAndCreateManager(cfg)
	if err != nil {
		log.Fatal("init manager failed, error: %v", err)
		return
	}

	svr, err := server.NewServer(cfg, mgr)
	if err != nil {
		log.Fatal("NewServer error, quit. error: %s", err.Error())
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
		syscall.SIGUSR1,
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			sig := <-sc
			if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				log.Notice("Got signal %d, quit", sig)
				svr.Close()
				break
			} else if sig == syscall.SIGPIPE {
				log.Notice("Ignore broken pipe signal")
			} else if sig == syscall.SIGUSR1 {
				log.Notice("Got update config signal")
			}
		}
	}()
	svr.Run()
	wg.Wait()
}

func initXLog(output, path, filename, level, service string) error {
	cfg := make(map[string]string)
	cfg["path"] = path
	cfg["filename"] = filename
	cfg["level"] = level
	cfg["service"] = service
	cfg["skip"] = "5" // 设置xlog打印方法堆栈需要跳过的层数, 5目前为调用log.Debug()等方法的方法名, 比xlog默认值多一层.

	logger, err := xlog.CreateLogManager(output, cfg)
	if err != nil {
		return err
	}

	log.SetGlobalLogger(logger)
	return nil
}
