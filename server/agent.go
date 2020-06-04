package server

import (
	"context"
	"fmt"

	"Asgard/applications"
	"Asgard/managers"
	"Asgard/rpc"
)

type AgentServer struct {
	baseServer
	appManager    *managers.AppManager
	jobManager    *managers.JobManager
	timingManager *managers.TimingManager
}

func (s *AgentServer) SetAppManager(appManager *managers.AppManager) {
	s.appManager = appManager
}

func (s *AgentServer) SetJobManager(jobManager *managers.JobManager) {
	s.jobManager = jobManager
}

func (s *AgentServer) SetTimingManager(timingManager *managers.TimingManager) {
	s.timingManager = timingManager
}

func (s *AgentServer) Stat(ctx context.Context, request *rpc.Empty) (*rpc.AgentStatResponse, error) {
	stat := &rpc.AgentStatResponse{
		Code: rpc.OK,
		AgentStat: &rpc.AgentStat{
			Apps:    int64(s.appManager.Count()),
			Jobs:    int64(s.jobManager.Count()),
			Timings: int64(s.timingManager.Count()),
		},
	}
	return stat, nil
}

func (s *AgentServer) Log(ctx context.Context, request *rpc.LogRuquest) (*rpc.LogResponse, error) {
	return &rpc.LogResponse{Content: GetLog(request.GetDir(), int(request.GetLines()))}, nil
}

func (s *AgentServer) AppList(ctx context.Context, request *rpc.Empty) (*rpc.AppListResponse, error) {
	apps := s.appManager.GetAppList()
	list := []*rpc.App{}
	for _, app := range apps {
		list = append(list, rpc.BuildApp(app))
	}
	return &rpc.AppListResponse{Code: rpc.OK, Apps: list}, nil
}

func (s *AgentServer) AppGet(ctx context.Context, request *rpc.ID) (*rpc.AppResponse, error) {
	app := s.appManager.GetApp(request.GetId())
	if app != nil {
		return &rpc.AppResponse{Code: rpc.OK, App: rpc.BuildApp(app)}, nil
	}
	return &rpc.AppResponse{Code: rpc.Nofound, App: nil}, nil
}

func (s *AgentServer) AppAdd(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	app := s.appManager.GetApp(request.GetId())
	if app != nil {
		if !app.Running {
			ok := s.appManager.Start(request.GetId())
			if !ok {
				return s.Error("app start failed!")
			}
		}
		return s.OK()
	}
	err := s.appManager.Register(request.GetId(), rpc.BuildAppConfig(request))
	if err != nil {
		return s.Error(fmt.Sprintf("app register failed:%s", err.Error()))
	}
	ok := s.appManager.Start(request.GetId())
	if !ok {
		return s.Error("app start failed!")
	}
	return s.OK()
}

func (s *AgentServer) AppUpdate(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	app := s.appManager.GetApp(request.GetId())
	if app == nil {
		return s.Error("no such app in agent")
	}
	if app.Running {
		ok := s.appManager.Stop(request.GetId())
		if !ok {
			return s.Error("app stop failed!")
		}
	}
	s.appManager.UnRegister(request.GetId())
	err := s.appManager.Register(request.GetId(), rpc.BuildAppConfig(request))
	if err != nil {
		return s.Error(fmt.Sprintf("app register failed:%s", err.Error()))
	}
	ok := s.appManager.Start(request.GetId())
	if !ok {
		return s.Error("app start failed!")
	}
	return s.OK()
}

func (s *AgentServer) AppRemove(ctx context.Context, request *rpc.ID) (*rpc.Response, error) {
	app := s.appManager.GetApp(request.GetId())
	if app == nil {
		return s.OK()
	}
	if app.Running {
		ok := s.appManager.Stop(request.GetId())
		if !ok {
			return s.Error("app stop failed!")
		}
	}
	s.appManager.UnRegister(request.GetId())
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

func (s *AgentServer) TimingList(ctx context.Context, request *rpc.Empty) (*rpc.TimingListResponse, error) {
	list := []*rpc.Timing{}
	for _, timing := range applications.Timings {
		if timing.Executed {
			continue
		}
		list = append(list, rpc.BuildTiming(timing))
	}
	return &rpc.TimingListResponse{Code: rpc.OK, Timings: list}, nil
}

func (s *AgentServer) TimingGet(ctx context.Context, request *rpc.ID) (*rpc.TimingResponse, error) {
	id := request.GetId()
	timing, ok := applications.Timings[id]
	if ok {
		return &rpc.TimingResponse{Code: rpc.OK, Timing: rpc.BuildTiming(timing)}, nil
	}
	return &rpc.TimingResponse{Code: rpc.Nofound, Timing: nil}, nil
}

func (s *AgentServer) TimingAdd(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
	id := request.GetId()
	_, ok := applications.Timings[id]
	if ok {
		return s.OK()
	}
	err := AddTiming(id, request)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) TimingUpdate(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
	id := request.GetId()
	timing, ok := applications.Timings[id]
	if !ok {
		return s.Error(fmt.Sprintf("no timing %d", id))
	}
	err := UpdateTiming(id, timing, request)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) TimingRemove(ctx context.Context, request *rpc.ID) (*rpc.Response, error) {
	id := request.GetId()
	timing, ok := applications.Timings[id]
	if !ok {
		return s.OK()
	}
	err := DeleteTiming(id, timing)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}
