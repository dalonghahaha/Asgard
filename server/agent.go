package server

import (
	"context"
	"fmt"

	"Asgard/applications"
	"Asgard/rpc"
)

type AgentServer struct {
	baseServer
}

func (s *AgentServer) Stat(ctx context.Context, request *rpc.Empty) (*rpc.AgentStatResponse, error) {
	stat := &rpc.AgentStatResponse{
		Code: rpc.OK,
		AgentStat: &rpc.AgentStat{
			Apps: int64(len(applications.APPs)),
			Jobs: int64(len(applications.Jobs)),
		},
	}
	return stat, nil
}

func (s *AgentServer) AppList(ctx context.Context, request *rpc.Empty) (*rpc.AppListResponse, error) {
	list := []*rpc.App{}
	for _, app := range applications.APPs {
		list = append(list, rpc.BuildApp(app))
	}
	return &rpc.AppListResponse{Code: rpc.OK, Apps: list}, nil
}

func (s *AgentServer) AppGet(ctx context.Context, request *rpc.ID) (*rpc.AppResponse, error) {
	id := request.GetId()
	app, ok := applications.APPs[id]
	if ok {
		return &rpc.AppResponse{Code: rpc.OK, App: rpc.BuildApp(app)}, nil
	}
	return &rpc.AppResponse{Code: rpc.Nofound, App: nil}, nil
}

func (s *AgentServer) AppAdd(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
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

func (s *AgentServer) AppUpdate(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
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

func (s *AgentServer) AppRemove(ctx context.Context, request *rpc.ID) (*rpc.Response, error) {
	id := request.GetId()
	app, ok := applications.APPs[id]
	if !ok {
		return s.OK()
	}
	err := DeleteApp(id, app)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) JobList(ctx context.Context, request *rpc.Empty) (*rpc.JobListResponse, error) {
	list := []*rpc.Job{}
	for _, job := range applications.Jobs {
		list = append(list, rpc.BuildJob(job))
	}
	return &rpc.JobListResponse{Code: rpc.OK, Jobs: list}, nil
}

func (s *AgentServer) JobGet(ctx context.Context, request *rpc.ID) (*rpc.JobResponse, error) {
	id := request.GetId()
	job, ok := applications.Jobs[id]
	if ok {
		return &rpc.JobResponse{Code: rpc.OK, Job: rpc.BuildJob(job)}, nil
	}
	return &rpc.JobResponse{Code: rpc.Nofound, Job: nil}, nil
}

func (s *AgentServer) JobAdd(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	id := request.GetId()
	_, ok := applications.Jobs[id]
	if ok {
		return s.OK()
	}
	err := AddJob(id, request)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) JobUpdate(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	id := request.GetId()
	job, ok := applications.Jobs[id]
	if !ok {
		return s.Error(fmt.Sprintf("no job %d", id))
	}
	err := UpdateJob(id, job, request)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) JobRemove(ctx context.Context, request *rpc.ID) (*rpc.Response, error) {
	id := request.GetId()
	job, ok := applications.Jobs[id]
	if !ok {
		return s.OK()
	}
	err := DeleteJob(id, job)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}
