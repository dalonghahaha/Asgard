package applications

import (
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
)

var APPs = []*App{}

func AppStopAll() {
	MoniterStop()
	processExit = true
	for _, app := range APPs {
		if !app.Finished {
			app.stop()
		}
	}
}

func AppStartAll() {
	for _, app := range APPs {
		go app.Run()
	}
	MoniterStart()
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

func AppStop(name string) bool {
	for _, app := range APPs {
		if app.Name == name {
			app.stop()
			return true
		}
	}
	return false
}

type App struct {
	Command
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

func AppRegister(config map[string]interface{}) error {
	app, err := NewApp(config)
	if err != nil {
		return err
	}
	APPs = append(APPs, app)
	return nil
}
