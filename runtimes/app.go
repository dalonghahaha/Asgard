package runtimes

import (
	"github.com/dalonghahaha/avenger/components/logger"
)

type App struct {
	Command
	ID          int64
	AutoRestart bool
	RestartTime int
}

func (a *App) Run() {
	err := a.build()
	if err != nil {
		return
	}
	err = a.start()
	if err != nil {
		return
	}
	go a.wait(a.restart)
}

func (a *App) restart() {
	if a.AutoRestart && !processExit {
		//最多重启5次
		if !a.Successed && a.RestartTime <= 5 {
			a.RestartTime += 1
			logger.Infof("%s Restart.....", a.Name)
			go a.Run()
		} else {
			logger.Warnf("%s Restart Reach Max", a.Name)
		}
	}
}
