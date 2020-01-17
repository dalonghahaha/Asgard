package server

import (
	"fmt"

	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
)

func AddApp(id int64, request *rpc.App) error {
	app, err := applications.AppRegister(id, rpc.BuildAppConfig(request))
	if err != nil {
		return err
	}
	app.MonitorReport = func(monitor *applications.Monitor) {
		client.AppMonitorReport(app, monitor)
	}
	app.ArchiveReport = func(command *applications.Command) {
		client.AppArchiveReport(app, command)
	}
	ok := applications.AppStartByID(id)
	if !ok {
		return fmt.Errorf("app %d start failed", id)
	}
	return nil
}

func UpdateApp(id int64, app *applications.App, request *rpc.App) error {
	err := DeleteApp(id, app)
	if err != nil {
		return err
	}
	return AddApp(id, request)
}

func DeleteApp(id int64, app *applications.App) error {
	app.AutoRestart = false
	ok := applications.AppStopByID(id)
	if !ok {
		return fmt.Errorf("app %d stop failed", id)
	}
	delete(applications.APPs, id)
	return nil
}
