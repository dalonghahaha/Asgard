package server

import (
	"context"

	"Asgard/managers"
	"Asgard/rpc"
)

type GuardServer struct {
	baseServer
	appManager *managers.AppManager
}

func (s *GuardServer) SetAppManager(appManager *managers.AppManager) {
	s.appManager = appManager
}

func (s *GuardServer) List(ctx context.Context, request *rpc.Empty) (*rpc.AppListResponse, error) {
	apps := s.appManager.GetList()
	list := []*rpc.App{}
	for _, app := range apps {
		list = append(list, rpc.BuildApp(app))
	}
	return &rpc.AppListResponse{Code: rpc.OK, Apps: list}, nil
}

func (s *GuardServer) Get(ctx context.Context, request *rpc.Name) (*rpc.AppResponse, error) {
	app := s.appManager.GetByName(request.GetName())
	if app != nil {
		return &rpc.AppResponse{Code: rpc.OK, App: rpc.BuildApp(app)}, nil
	}
	return &rpc.AppResponse{Code: rpc.Nofound, App: nil}, nil
}

func (s *GuardServer) Add(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	err := s.appManager.Add(request.GetId(), rpc.BuildAppConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *GuardServer) Update(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	err := s.appManager.Update(request.GetId(), rpc.BuildAppConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *GuardServer) Remove(ctx context.Context, request *rpc.Name) (*rpc.Response, error) {
	app := s.appManager.GetByName(request.GetName())
	if app == nil {
		return s.OK()
	}
	err := s.appManager.Remove(app.ID)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}
