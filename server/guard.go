package server

import (
	"context"

	"Asgard/applications"
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
	return s.OK()
}

func (s *GuardServer) Update(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	return s.OK()
}

func (s *GuardServer) Remove(ctx context.Context, request *rpc.AppNameRequest) (*rpc.Response, error) {
	return s.OK()
}
