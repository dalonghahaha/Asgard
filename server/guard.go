package server

import (
	"context"

	"Asgard/rpc"
)

type GuardServer struct {
	baseServer
}

func (s *GuardServer) List(ctx context.Context, request *rpc.Empty) (*rpc.AppListResponse, error) {
	return &rpc.AppListResponse{Code: rpc.OK, Apps: GetAppList()}, nil
}

func (s *GuardServer) Get(ctx context.Context, request *rpc.Name) (*rpc.AppResponse, error) {
	app := GetAppByName(request.GetName())
	if app != nil {
		return &rpc.AppResponse{Code: rpc.OK, App: app}, nil
	}
	return &rpc.AppResponse{Code: rpc.Nofound, App: app}, nil
}

func (s *GuardServer) Add(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	if err := AddApp(request.GetId(), request); err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *GuardServer) Update(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	if err := UpdateApp(request.GetId(), request); err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *GuardServer) Remove(ctx context.Context, request *rpc.Name) (*rpc.Response, error) {
	if err := DeleteAppByName(request.GetName()); err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}
