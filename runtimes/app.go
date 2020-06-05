package runtimes

import (
	"fmt"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/patrickmn/go-cache"
)

var mcache = cache.New(5*time.Second, 10*time.Minute)

type App struct {
	Command
	ID          int64
	AutoRestart bool
	Dead        bool
}

func (a *App) getCacheKey() string {
	return fmt.Sprintf("mcache:app:%d", a.ID)
}

func (a *App) Run() {
	err := a.build()
	if err != nil {
		return
	}
	err = a.start()
	if err != nil {
		logger.Errorf("%s start fail: %s", a.Name, err)
		a.Dead = true
		if a.ExceptionReport != nil {
			a.ExceptionReport(fmt.Sprintf("start fail: %+v", err))
		}
		return
	}
	go a.wait(a.restart)
}

func (a *App) restart() {
	if a.AutoRestart && !processExit {
		_, ok := mcache.Get(a.getCacheKey())
		var restartTime int
		if !ok {
			err := mcache.Add(a.getCacheKey(), 0, cache.DefaultExpiration)
			if err != nil {
				logger.Warnf("add restart time cache failed:%+v", err)
			}
			restartTime = 0
		} else {
			var err error
			restartTime, err = mcache.IncrementInt(a.getCacheKey(), 1)
			if err != nil {
				logger.Warnf("increment restart time failed:%+v", err)
			}
		}
		//最多重启5次
		if restartTime > 5 {
			logger.Warnf("%s Restart Reach Max", a.Name)
			a.Dead = true
			if a.ExceptionReport != nil {
				a.ExceptionReport("restart reach max limited")
			}
		} else {
			go a.Run()
		}
	}
}
