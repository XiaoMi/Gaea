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
	"github.com/XiaoMi/Gaea/core"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/proxy/server"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var configFile = flag.String("config", "etc/gaea.ini", "gaea config file")
var info = flag.Bool("info", false, "show info of gaea")
var numCPU = flag.Int("num-cpu", 0, "how many operating systems threads attempt to execute simultaneously")

const (
	defaultCPUNum = 4 //best practise
)

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

	var cpuNums int
	if *numCPU > 0 && *numCPU <= runtime.NumCPU() {
		// best practice -num-cpu=4
		fmt.Printf("max use cpu : %d(from command-line)", *numCPU)
		cpuNums = *numCPU
	} else if cfg.NumCPU > 0 && *numCPU <= runtime.NumCPU() {
		fmt.Printf("max use cpu : %d(from cfg)", cfg.NumCPU)
		cpuNums = cfg.NumCPU
	} else {
		fmt.Printf("max use cpu : %d(default)", defaultCPUNum)
		cpuNums = defaultCPUNum
	}
	runtime.GOMAXPROCS(cpuNums)
	cfg.NumCPU = cpuNums

	if err = models.InitXLog(cfg.LogOutput, cfg.LogPath, cfg.LogFileName, cfg.LogLevel, cfg.Service, cfg.LogKeepDays, cfg.LogKeepCounts); err != nil {
		fmt.Printf("init xlog error: %v\n", err.Error())
		return
	}
	defer log.Close()

	// init manager
	mgr, err := server.LoadAndCreateManager(cfg)
	if err != nil {
		log.Fatal("init manager failed, error: %v", err)
		fmt.Printf("init manager failed, error: %v\n", err)
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
				log.Notice("got signal %d, quit", sig)
				svr.Close()
				break
			} else if sig == syscall.SIGPIPE {
				log.Notice("ignore broken pipe signal")
			} else if sig == syscall.SIGUSR1 {
				log.Notice("got update config signal,old config:%#v", cfg)
				if err := models.ReloadProxyConfig(*configFile); err != nil {
					log.Notice("reload proxy confi error:%s", err)
				} else {
					log.Notice("reload proxy config success")
				}
			}
		}
	}()
	svr.Run()
	wg.Wait()
}
