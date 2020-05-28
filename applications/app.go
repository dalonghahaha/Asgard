package applications

import (
	"fmt"
	"sync"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
)

var APPs = map[int64]*App{}

func AppStopAll() {
	MoniterStop()
	processExit = true
	for _, app := range APPs {
		if !app.Finished {
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

func AppStart(name string) bool {
	for _, app := range APPs {
		if app.Name == name {
			go app.Run()
			return true
		}
	}
	return false
}

func AppStartByID(id int64) bool {
	app, ok := APPs[id]
	if !ok {
		return false
	}
	go app.Run()
	return true
}

func AppStop(name string) bool {
	for _, app := range APPs {
		if app.Name == name {
			app.stop()
			return true
		}
	}
	return false
}

func AppStopByID(id int64) bool {
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
		logger.Info(a.Name + " Restart.....")
		go a.Run()
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

func AppRegister(id int64, config map[string]interface{}, reports *sync.Map, appMonitorChan chan AppMonitor, appArchiveChan chan AppArchive) error {
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
		reports.Store(appMonitor.UUID, 1)
		appMonitorChan <- appMonitor
	}
	app.ArchiveReport = func(archive *Archive) {
		appArchive := AppArchive{
			UUID:    uuid.GenerateV4(),
			App:     app,
			Archive: archive,
		}
		reports.Store(appArchive.UUID, 1)
		appArchiveChan <- appArchive
	}
	APPs[id] = app
	return nil
}
