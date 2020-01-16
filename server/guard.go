package server

import (
	"context"
	"fmt"

	"Asgard/applications"
	"Asgard/rpc"
)

type GuardServer struct {
	baseServer
}

func (s *GuardServer) List(ctx context.Context, request *rpc.Empty) (*rpc.AppListResponse, error) {
	list := []*rpc.App{}
	for _, app := range applications.APPs {
		list = append(list, rpc.BuildApp(app))
	}
	return &rpc.AppListResponse{Code: rpc.OK, Apps: list}, nil
}

func (s *GuardServer) Get(ctx context.Context, request *rpc.Name) (*rpc.AppResponse, error) {
	for _, app := range applications.APPs {
		if request.GetName() == app.Name {
			return &rpc.AppResponse{Code: rpc.OK, App: rpc.BuildApp(app)}, nil
		}
	}
	return &rpc.AppResponse{Code: rpc.Nofound, App: nil}, nil
}

func (s *GuardServer) Add(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	id := request.GetId()
	_, ok := applications.APPs[id]
	if ok {
		return s.OK()
	}
	err := AddApp(id, request)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *GuardServer) Update(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	id := request.GetId()
	app, ok := applications.APPs[id]
	if !ok {
		return s.Error(fmt.Sprintf("no app %d", id))
	}
	err := UpdateApp(id, app, request)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *GuardServer) Remove(ctx context.Context, request *rpc.Name) (*rpc.Response, error) {
	for _, app := range applications.APPs {
		if request.GetName() == app.Name {
			ok := applications.AppStopByID(app.ID)
			if !ok {
				return s.Error(fmt.Sprintf("job %s stop failed", request.GetName()))
			}
		}
	}
	return s.OK()
}
