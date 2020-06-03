package applications

import (
	"fmt"
	"sync"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
)

var (
	APPs    = map[int64]*App{}
	appLock sync.Mutex
)

func AppStopAll() {
	for _, app := range APPs {
		if app.Running {
			app.stop()
		}
	}
}

func AppStartAll(moniter bool) {
	for _, app := range APPs {
		go app.Run()
	}
	if moniter {
		MoniterStart()
	}
}

func AppStart(id int64) bool {
	app, ok := APPs[id]
	if !ok {
		return false
	}
	if app.Running {
		return true
	}
	go app.Run()
	return true
}

func AppStop(id int64) bool {
	app, ok := APPs[id]
	if !ok {
		return false
	}
	app.stop()
	return true
}

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

func NewApp(config map[string]interface{}) (*App, error) {
	app := new(App)
	err := app.configure(config)
	if err != nil {
		return nil, err
	}
	autoRestart, ok := config["auto_restart"].(bool)
	if !ok {
		return nil, fmt.Errorf("config auto_restart type wrong")
	}
	app.AutoRestart = autoRestart
	return app, nil
}

func AppRegister(id int64, config map[string]interface{}, reports *sync.Map, mc chan AppMonitor, ac chan AppArchive) error {
	app, err := NewApp(config)
	if err != nil {
		return err
	}
	app.ID = id
	app.MonitorReport = func(monitor *Monitor) {
		appMonitor := AppMonitor{
			UUID:    uuid.GenerateV4(),
			App:     app,
			Monitor: monitor,
		}
		if reports != nil {
			reports.Store(appMonitor.UUID, 1)
		}
		if mc != nil {
			mc <- appMonitor
		}
	}
	app.ArchiveReport = func(archive *Archive) {
		appArchive := AppArchive{
			UUID:    uuid.GenerateV4(),
			App:     app,
			Archive: archive,
		}
		if reports != nil {
			reports.Store(appArchive.UUID, 1)
		}
		if ac != nil {
			ac <- appArchive
		}
	}
	appLock.Lock()
	APPs[id] = app
	appLock.Unlock()
	return nil
}

func AppUnRegister(id int64) {
	appLock.Lock()
	delete(APPs, id)
	appLock.Unlock()
}
