package server

import (
	"context"
	"fmt"

	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
)

type GuardServer struct {
	baseServer
}

func (s *GuardServer) List(ctx context.Context, request *rpc.Empty) (*rpc.AppListResponse, error) {
	apps := applications.APPs
	list := []*rpc.App{}
	for _, app := range apps {
		_app := new(rpc.App)
		_app.Id = app.ID
		_app.Name = app.Name
		_app.Dir = app.Dir
		_app.Program = app.Program
		_app.Args = app.Args
		_app.StdOut = app.Stdout
		_app.StdErr = app.Stderr
		_app.AutoRestart = app.AutoRestart
		list = append(list, _app)
	}
	return &rpc.AppListResponse{Code: 200, Apps: list}, nil
}

func (s *GuardServer) Get(ctx context.Context, request *rpc.AppNameRequest) (*rpc.AppResponse, error) {
	apps := applications.APPs
	name := request.GetName()
	for _, app := range apps {
		if name == app.Name {
			_app := new(rpc.App)
			_app.Name = app.Name
			_app.Dir = app.Dir
			_app.Program = app.Program
			_app.Args = app.Args
			_app.StdOut = app.Stdout
			_app.StdErr = app.Stderr
			_app.AutoRestart = app.AutoRestart
			return &rpc.AppResponse{Code: 200, App: _app}, nil
		}
	}
	return &rpc.AppResponse{Code: 0, App: nil}, nil
}

func (s *GuardServer) Add(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	id := request.GetId()
	_, ok := applications.APPs[id]
	if ok {
		return s.OK()
	}
	config := map[string]interface{}{
		"id":           request.GetId(),
		"name":         request.GetName(),
		"dir":          request.GetDir(),
		"program":      request.GetProgram(),
		"args":         request.GetArgs(),
		"stdout":       request.GetStdOut(),
		"stderr":       request.GetStdErr(),
		"auto_restart": request.GetAutoRestart(),
		"is_monitor":   request.GetIsMonitor(),
	}
	app, err := applications.AppRegister(id, config)
	if err != nil {
		return s.Error(err.Error())
	}
	app.MonitorReport = func(monitor *applications.Monitor) {
		client.AppMonitorReport(app, monitor)
	}
	app.ArchiveReport = func(command *applications.Command) {
		client.AppArchiveReport(app, command)
	}
	ok = applications.AppStartByID(id)
	if !ok {
		return s.Error(fmt.Sprintf("app %d start failed", id))
	}
	return s.OK()
}

func (s *GuardServer) Update(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	id := request.GetId()
	app, ok := applications.APPs[id]
	if !ok {
		return s.Error(fmt.Sprintf("no app %d", id))
	}
	ok = applications.AppStopByID(id)
	if !ok {
		return s.Error(fmt.Sprintf("app %d stop failed", id))
	}
	app.Name = request.GetName()
	app.Dir = request.GetDir()
	app.Program = request.GetProgram()
	app.Args = request.GetArgs()
	app.Stdout = request.GetStdOut()
	app.Stderr = request.GetStdErr()
	app.AutoRestart = request.GetAutoRestart()
	app.IsMonitor = request.GetIsMonitor()
	ok = applications.AppStartByID(id)
	if !ok {
		return s.Error(fmt.Sprintf("app %d start failed", id))
	}
	return s.OK()
}

func (s *GuardServer) Remove(ctx context.Context, request *rpc.AppNameRequest) (*rpc.Response, error) {
	return s.OK()
}
