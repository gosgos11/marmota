// Copyright 2022 Teamgram Authors
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: teamgramio (teamgram.io@gmail.com)
//

package commands

import (
	"errors"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

////////////////////////////////////////////////////////////////
var (
	GMainInst MainInstance
	GSignal   chan os.Signal
)

type MainInstance interface {
	Initialize() error
	RunLoop()
	Destroy()
}

func Run(inst MainInstance) {
	flag.Parse()
	// if err := paladin.Init(); err != nil {
	//	panic(err)
	//}
	// logx.Init(nil) // debug flag: log.dir={path}
	defer logx.Close()

	if inst == nil {
		panic(errors.New("inst is nil, exit"))
		return
	}

	//
	rand.Seed(time.Now().UTC().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())

	//if enableZipkin {
	//	zipkin.Init(&zipkin.Config{
	//		Endpoint: "http://localhost:9411/api/v2/spans",
	//	})
	//} else if enableJaeger {
	//	jaeger.Init()
	//}

	//log.SetFormat("[%D %T] [%L] [%S] %M")
	logx.Info("instance initialize...")
	err := inst.Initialize()
	logx.Info("inited")
	if err != nil {
		panic(err)
		return
	}

	// global
	GMainInst = inst

	logx.Info("instance run_loop...")
	go inst.RunLoop()

	GSignal = make(chan os.Signal, 1)
	signal.Notify(GSignal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-GSignal
		logx.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logx.Infof("instance exit...")
			inst.Destroy()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
