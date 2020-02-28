package server

import (
	"context"

	"Asgard/models"
	"Asgard/rpc"
	"Asgard/services"
)

type MasterServer struct {
	baseServer
	agentService   *services.AgentService
	appService     *services.AppService
	jobService     *services.JobService
	timingService  *services.TimingService
	monitorService *services.MonitorService
	archiveService *services.ArchiveService
}

func NewMasterServer() *MasterServer {
	return &MasterServer{
		agentService:   services.NewAgentService(),
		appService:     services.NewAppService(),
		jobService:     services.NewJobService(),
		timingService:  services.NewTimingService(),
		monitorService: services.NewMonitorService(),
		archiveService: services.NewArchiveService(),
	}
}

func (s *MasterServer) Register(ctx context.Context, request *rpc.AgentInfo) (*rpc.Response, error) {
	agent := s.agentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent != nil {
		agent.Status = 1
		s.agentService.UpdateAgent(agent)
		return s.OK()
	}
	agent = new(models.Agent)
	agent.IP = request.GetIp()
	agent.Port = request.GetPort()
	agent.Status = 1
	ok := s.agentService.CreateAgent(agent)
	if !ok {
		return s.Error("CreateAgent Failed")
	}
	return s.OK()
}

func (s *MasterServer) AppList(ctx context.Context, request *rpc.AgentInfo) (*rpc.AppListResponse, error) {
	agent := s.agentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.AppListResponse{Code: rpc.Nofound, Apps: nil}, nil
	}
	apps := s.appService.GetAppByAgentID(agent.ID)
	list := []*rpc.App{}
	for _, app := range apps {
		list = append(list, rpc.FormatApp(&app))
	}
	return &rpc.AppListResponse{Code: rpc.OK, Apps: list}, nil
}

func (s *MasterServer) JobList(ctx context.Context, request *rpc.AgentInfo) (*rpc.JobListResponse, error) {
	agent := s.agentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.JobListResponse{Code: rpc.Nofound, Jobs: nil}, nil
	}
	jobs := s.jobService.GetJobByAgentID(agent.ID)
	list := []*rpc.Job{}
	for _, job := range jobs {
		list = append(list, rpc.FormatJob(&job))
	}
	return &rpc.JobListResponse{Code: rpc.OK, Jobs: list}, nil
}

func (s *MasterServer) TimingList(ctx context.Context, request *rpc.AgentInfo) (*rpc.TimingListResponse, error) {
	agent := s.agentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.TimingListResponse{Code: rpc.Nofound, Timings: nil}, nil
	}
	timings := s.timingService.GetTimingByAgentID(agent.ID)
	list := []*rpc.Timing{}
	for _, timing := range timings {
		list = append(list, rpc.FormatTiming(&timing))
	}
	return &rpc.TimingListResponse{Code: rpc.OK, Timings: list}, nil
}

func (s *MasterServer) AgentMonitorReport(ctx context.Context, request *rpc.AgentMonitor) (*rpc.Response, error) {
	agent := s.agentService.GetAgentByIPAndPort(request.GetAgent().GetIp(), request.GetAgent().GetPort())
	if agent == nil {
		return s.Error("no such agent!")
	}
	ok := s.monitorService.CreateMonitor(rpc.ParseMonitor(models.TYPE_AGENT, agent.ID, request.GetMonitor()))
	if !ok {
		return s.Error("add agent monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) AppMonitorReport(ctx context.Context, request *rpc.AppMonitor) (*rpc.Response, error) {
	ok := s.monitorService.CreateMonitor(rpc.ParseMonitor(models.TYPE_APP, request.GetApp().GetId(), request.GetMonitor()))
	if !ok {
		return s.Error("add app monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) JobMoniorReport(ctx context.Context, request *rpc.JobMonior) (*rpc.Response, error) {
	ok := s.monitorService.CreateMonitor(rpc.ParseMonitor(models.TYPE_JOB, request.GetJob().GetId(), request.GetMonitor()))
	if !ok {
		return s.Error("add job monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) TimingMoniorReport(ctx context.Context, request *rpc.TimingMonior) (*rpc.Response, error) {
	ok := s.monitorService.CreateMonitor(rpc.ParseMonitor(models.TYPE_TIMING, request.GetTiming().GetId(), request.GetMonitor()))
	if !ok {
		return s.Error("add timing monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) AppArchiveReport(ctx context.Context, request *rpc.AppArchive) (*rpc.Response, error) {
	ok := s.archiveService.CreateArchive(rpc.ParseArchive(models.TYPE_APP, request.GetApp().GetId(), request.GetArchive()))
	if !ok {
		return s.Error("add app archive failed")
	}
	return s.OK()
}

func (s *MasterServer) JobArchiveReport(ctx context.Context, request *rpc.JobArchive) (*rpc.Response, error) {
	archiveService := services.NewArchiveService()
	ok := archiveService.CreateArchive(rpc.ParseArchive(models.TYPE_JOB, request.GetJob().GetId(), request.GetArchive()))
	if !ok {
		return s.Error("add job archive failed")
	}
	return s.OK()
}

func (s *MasterServer) TimingArchiveReport(ctx context.Context, request *rpc.TimingArchive) (*rpc.Response, error) {
	archiveService := services.NewArchiveService()
	ok := archiveService.CreateArchive(rpc.ParseArchive(models.TYPE_TIMING, request.GetTiming().GetId(), request.GetArchive()))
	if !ok {
		return s.Error("add timing archive failed")
	}
	return s.OK()
}
