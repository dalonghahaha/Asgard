package server

import (
	"context"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/rpc"
)

type MasterServer struct {
	baseServer
}

func NewMasterServer() *MasterServer {
	return &MasterServer{}
}

func (s *MasterServer) Register(ctx context.Context, request *rpc.AgentInfo) (*rpc.Response, error) {
	agent := providers.AgentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent != nil {
		//禁止状态的实例直接忽略注册
		if agent.Status == constants.AGENT_FORBIDDEN {
			return s.OK()
		}
		agent.Status = constants.AGENT_ONLINE
		providers.AgentService.UpdateAgent(agent)
		return s.OK()
	}
	agent = new(models.Agent)
	agent.IP = request.GetIp()
	agent.Port = request.GetPort()
	agent.Master = constants.MASTER_IP
	agent.Status = constants.AGENT_ONLINE
	ok := providers.AgentService.CreateAgent(agent)
	if !ok {
		return s.Error("CreateAgent Failed")
	}
	return s.OK()
}

func (s *MasterServer) AppList(ctx context.Context, request *rpc.AgentInfo) (*rpc.AppListResponse, error) {
	agent := providers.AgentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.AppListResponse{Code: rpc.Nofound, Apps: nil}, nil
	}
	apps := providers.AppService.GetUsageAppByAgentID(agent.ID)
	list := []*rpc.App{}
	for _, app := range apps {
		list = append(list, rpc.FormatApp(&app))
	}
	return &rpc.AppListResponse{Code: rpc.OK, Apps: list}, nil
}

func (s *MasterServer) JobList(ctx context.Context, request *rpc.AgentInfo) (*rpc.JobListResponse, error) {
	agent := providers.AgentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.JobListResponse{Code: rpc.Nofound, Jobs: nil}, nil
	}
	jobs := providers.JobService.GetUsageJobByAgentID(agent.ID)
	list := []*rpc.Job{}
	for _, job := range jobs {
		list = append(list, rpc.FormatJob(&job))
	}
	return &rpc.JobListResponse{Code: rpc.OK, Jobs: list}, nil
}

func (s *MasterServer) TimingList(ctx context.Context, request *rpc.AgentInfo) (*rpc.TimingListResponse, error) {
	agent := providers.AgentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.TimingListResponse{Code: rpc.Nofound, Timings: nil}, nil
	}
	timings := providers.TimingService.GetUsageTimingByAgentID(agent.ID)
	list := []*rpc.Timing{}
	for _, timing := range timings {
		list = append(list, rpc.FormatTiming(&timing))
	}
	return &rpc.TimingListResponse{Code: rpc.OK, Timings: list}, nil
}

func (s *MasterServer) AgentMonitorReport(ctx context.Context, request *rpc.AgentMonitor) (*rpc.Response, error) {
	agent := providers.AgentService.GetAgentByIPAndPort(request.GetAgent().GetIp(), request.GetAgent().GetPort())
	if agent == nil {
		return s.Error("no such agent!")
	}
	ok := providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_AGENT, agent.ID, request.GetMonitor()))
	if !ok {
		return s.Error("add agent monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) AppMonitorReport(ctx context.Context, request *rpc.AppMonitor) (*rpc.Response, error) {
	ok := providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_APP, request.GetAppId(), request.GetMonitor()))
	if !ok {
		return s.Error("add app monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) JobMoniorReport(ctx context.Context, request *rpc.JobMonior) (*rpc.Response, error) {
	ok := providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_JOB, request.GetJobId(), request.GetMonitor()))
	if !ok {
		return s.Error("add job monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) TimingMoniorReport(ctx context.Context, request *rpc.TimingMonior) (*rpc.Response, error) {
	ok := providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_TIMING, request.GetTimingId(), request.GetMonitor()))
	if !ok {
		return s.Error("add timing monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) AgentMonitorBatchReport(ctx context.Context, request *rpc.AgentMonitorList) (*rpc.Response, error) {
	agent := providers.AgentService.GetAgentByIPAndPort(request.GetAgent().GetIp(), request.GetAgent().GetPort())
	if agent == nil {
		return s.Error("no such agent!")
	}
	for _, monitor := range request.GetMonitors() {
		providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_AGENT, agent.ID, monitor))
	}
	return s.OK()
}

func (s *MasterServer) AppMonitorBatchReport(ctx context.Context, request *rpc.AppMonitorList) (*rpc.Response, error) {
	for _, monitor := range request.GetMonitors() {
		providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_AGENT, request.GetAppId(), monitor))
	}
	return s.OK()
}

func (s *MasterServer) JobMoniorBatchReport(ctx context.Context, request *rpc.JobMonitorList) (*rpc.Response, error) {
	for _, monitor := range request.GetMonitors() {
		providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_AGENT, request.GetJobId(), monitor))
	}
	return s.OK()
}

func (s *MasterServer) TimingMoniorBatchReport(ctx context.Context, request *rpc.TimingMoniorList) (*rpc.Response, error) {
	for _, monitor := range request.GetMonitors() {
		providers.MoniterService.CreateMonitor(rpc.ParseMonitor(constants.TYPE_AGENT, request.GetTimingId(), monitor))
	}
	return s.OK()
}

func (s *MasterServer) AppArchiveReport(ctx context.Context, request *rpc.AppArchive) (*rpc.Response, error) {
	ok := providers.ArchiveService.CreateArchive(rpc.ParseArchive(constants.TYPE_APP, request.GetAppId(), request.GetArchive()))
	if !ok {
		return s.Error("add app archive failed")
	}
	if constants.MASTER_NOTIFY && request.GetArchive().GetStatus() != 0 {
		app := providers.AppService.GetAppByID(request.GetAppId())
		if app != nil {
			agent := providers.AgentService.GetAgentByID(app.AgentID)
			if agent != nil {
				go providers.NoticeService.AppUnsuccessNotify(app, agent, request)
			}
		}
	}
	return s.OK()
}

func (s *MasterServer) JobArchiveReport(ctx context.Context, request *rpc.JobArchive) (*rpc.Response, error) {
	ok := providers.ArchiveService.CreateArchive(rpc.ParseArchive(constants.TYPE_JOB, request.GetJobId(), request.GetArchive()))
	if !ok {
		return s.Error("add job archive failed")
	}
	if constants.MASTER_NOTIFY && request.GetArchive().GetStatus() != 0 {
		job := providers.JobService.GetJobByID(request.GetJobId())
		if job != nil {
			agent := providers.AgentService.GetAgentByID(job.AgentID)
			if agent != nil {
				go providers.NoticeService.JobUnsuccessNotify(job, agent, request)
			}
		}
	}
	return s.OK()
}

func (s *MasterServer) TimingArchiveReport(ctx context.Context, request *rpc.TimingArchive) (*rpc.Response, error) {
	ok := providers.ArchiveService.CreateArchive(rpc.ParseArchive(constants.TYPE_TIMING, request.GetTimingId(), request.GetArchive()))
	if !ok {
		return s.Error("add timing archive failed")
	}
	timing := providers.TimingService.GetTimingByID(request.GetTimingId())
	if timing != nil {
		providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_FINISHED, 0)
	}
	if constants.MASTER_NOTIFY && request.GetArchive().GetStatus() != 0 {
		agent := providers.AgentService.GetAgentByID(timing.AgentID)
		if agent != nil {
			go providers.NoticeService.TimingUnsuccessNotify(timing, agent, request)
		}
	}
	return s.OK()
}

func (s *MasterServer) AppExceptionReport(ctx context.Context, request *rpc.AppException) (*rpc.Response, error) {
	ok := providers.ExceptionService.CreateException(rpc.ParseException(constants.TYPE_APP, request.GetAppId(), request.GetDesc()))
	if !ok {
		return s.Error("add app exception failed")
	}
	return s.OK()
}

func (s *MasterServer) JobExceptionReport(ctx context.Context, request *rpc.JobException) (*rpc.Response, error) {
	ok := providers.ExceptionService.CreateException(rpc.ParseException(constants.TYPE_JOB, request.GetJobId(), request.GetDesc()))
	if !ok {
		return s.Error("add job exception failed")
	}
	return s.OK()
}

func (s *MasterServer) TimingExceptionReport(ctx context.Context, request *rpc.TimingException) (*rpc.Response, error) {
	ok := providers.ExceptionService.CreateException(rpc.ParseException(constants.TYPE_TIMING, request.GetTimingId(), request.GetDesc()))
	if !ok {
		return s.Error("add timing exception failed")
	}
	return s.OK()
}
