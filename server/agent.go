package server

import (
	"context"

	"Asgard/managers"
	"Asgard/rpc"
)

type AgentServer struct {
	baseServer
	agentManager *managers.AgentManager
}

func NewAgentServer(agentManager *managers.AgentManager) *AgentServer {
	return &AgentServer{
		agentManager: agentManager,
	}
}
func (s *AgentServer) Stat(ctx context.Context, request *rpc.Empty) (*rpc.AgentStatResponse, error) {
	stat := &rpc.AgentStatResponse{
		Code: rpc.OK,
		AgentStat: &rpc.AgentStat{
			Apps:    int64(s.agentManager.GetAppManager().Count()),
			Jobs:    int64(s.agentManager.GetJobManager().Count()),
			Timings: int64(s.agentManager.GetTimingManager().Count()),
		},
	}
	return stat, nil
}

func (s *AgentServer) Log(ctx context.Context, request *rpc.LogRuquest) (*rpc.LogResponse, error) {
	return &rpc.LogResponse{Content: GetLog(request.GetDir(), int(request.GetLines()))}, nil
}

func (s *AgentServer) AppList(ctx context.Context, request *rpc.Empty) (*rpc.AppListResponse, error) {
	apps := s.agentManager.GetAppManager().GetList()
	list := []*rpc.App{}
	for _, app := range apps {
		list = append(list, rpc.BuildApp(app))
	}
	return &rpc.AppListResponse{Code: rpc.OK, Apps: list}, nil
}

func (s *AgentServer) AppGet(ctx context.Context, request *rpc.ID) (*rpc.AppResponse, error) {
	app := s.agentManager.GetAppManager().Get(request.GetId())
	if app != nil {
		return &rpc.AppResponse{Code: rpc.OK, App: rpc.BuildApp(app)}, nil
	}
	return &rpc.AppResponse{Code: rpc.Nofound, App: nil}, nil
}

func (s *AgentServer) AppAdd(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	err := s.agentManager.GetAppManager().Add(request.GetId(), rpc.BuildAppConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) AppUpdate(ctx context.Context, request *rpc.App) (*rpc.Response, error) {
	err := s.agentManager.GetAppManager().Update(request.GetId(), rpc.BuildAppConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) AppRemove(ctx context.Context, request *rpc.ID) (*rpc.Response, error) {
	err := s.agentManager.GetAppManager().Remove(request.GetId())
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) JobList(ctx context.Context, request *rpc.Empty) (*rpc.JobListResponse, error) {
	jobs := s.agentManager.GetJobManager().GetList()
	list := []*rpc.Job{}
	for _, job := range jobs {
		list = append(list, rpc.BuildJob(job))
	}
	return &rpc.JobListResponse{Code: rpc.OK, Jobs: list}, nil
}

func (s *AgentServer) JobGet(ctx context.Context, request *rpc.ID) (*rpc.JobResponse, error) {
	job := s.agentManager.GetJobManager().Get(request.GetId())
	if job != nil {
		return &rpc.JobResponse{Code: rpc.OK, Job: rpc.BuildJob(job)}, nil
	}
	return &rpc.JobResponse{Code: rpc.Nofound, Job: nil}, nil
}

func (s *AgentServer) JobAdd(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	err := s.agentManager.GetJobManager().Add(request.GetId(), rpc.BuildJobConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) JobUpdate(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	err := s.agentManager.GetJobManager().Update(request.GetId(), rpc.BuildJobConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) JobRemove(ctx context.Context, request *rpc.ID) (*rpc.Response, error) {
	err := s.agentManager.GetJobManager().Remove(request.GetId())
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) TimingList(ctx context.Context, request *rpc.Empty) (*rpc.TimingListResponse, error) {
	timings := s.agentManager.GetTimingManager().GetList()
	list := []*rpc.Timing{}
	for _, timing := range timings {
		if timing.Executed {
			continue
		}
		list = append(list, rpc.BuildTiming(timing))
	}
	return &rpc.TimingListResponse{Code: rpc.OK, Timings: list}, nil
}

func (s *AgentServer) TimingGet(ctx context.Context, request *rpc.ID) (*rpc.TimingResponse, error) {
	timing := s.agentManager.GetTimingManager().Get(request.GetId())
	if timing != nil {
		return &rpc.TimingResponse{Code: rpc.OK, Timing: rpc.BuildTiming(timing)}, nil
	}
	return &rpc.TimingResponse{Code: rpc.Nofound, Timing: nil}, nil
}

func (s *AgentServer) TimingAdd(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
	err := s.agentManager.GetTimingManager().Register(request.GetId(), rpc.BuildTimingConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) TimingUpdate(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
	err := s.agentManager.GetTimingManager().Update(request.GetId(), rpc.BuildTimingConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *AgentServer) TimingRemove(ctx context.Context, request *rpc.ID) (*rpc.Response, error) {
	err := s.agentManager.GetTimingManager().Remove(request.GetId())
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}
