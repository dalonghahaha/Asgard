package server

import (
	"context"
	"time"

	"Asgard/models"
	"Asgard/rpc"
	"Asgard/services"
)

type MasterServer struct {
	baseServer
}

func (s *MasterServer) Register(ctx context.Context, request *rpc.AgentInfo) (*rpc.Response, error) {
	agentService := services.NewAgentService()
	agent := agentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent != nil {
		agent.Status = 1
		agentService.UpdateAgent(agent)
		return s.OK()
	}
	agent = new(models.Agent)
	agent.IP = request.GetIp()
	agent.Port = request.GetPort()
	agent.Status = 1
	ok := agentService.CreateAgent(agent)
	if !ok {
		return s.Error("CreateAgent Failed")
	}
	return s.OK()
}

func (s *MasterServer) AppList(ctx context.Context, request *rpc.AgentInfo) (*rpc.AppListResponse, error) {
	agentService := services.NewAgentService()
	appService := services.NewAppService()
	agent := agentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.AppListResponse{Code: 404, Apps: nil}, nil
	}
	apps := appService.GetAppByAgentID(agent.ID)
	list := []*rpc.App{}
	for _, app := range apps {
		_app := new(rpc.App)
		_app.Id = app.ID
		_app.Name = app.Name
		_app.Dir = app.Dir
		_app.Program = app.Program
		_app.Args = app.Args
		_app.StdOut = app.StdOut
		_app.StdErr = app.StdErr
		if app.AutoRestart == 1 {
			_app.AutoRestart = true
		} else {
			_app.AutoRestart = false
		}
		if app.IsMonitor == 1 {
			_app.IsMonitor = true
		} else {
			_app.IsMonitor = false
		}
		list = append(list, _app)
	}
	return &rpc.AppListResponse{Code: 200, Apps: list}, nil
}

func (s *MasterServer) JobList(ctx context.Context, request *rpc.AgentInfo) (*rpc.JobListResponse, error) {
	agentService := services.NewAgentService()
	jobService := services.NewJobService()
	agent := agentService.GetAgentByIPAndPort(request.GetIp(), request.GetPort())
	if agent == nil {
		return &rpc.JobListResponse{Code: 404, Jobs: nil}, nil
	}
	jobs := jobService.GetJobByAgentID(agent.ID)
	list := []*rpc.Job{}
	for _, job := range jobs {
		_job := new(rpc.Job)
		_job.Id = job.ID
		_job.Name = job.Name
		_job.Dir = job.Dir
		_job.Program = job.Program
		_job.Args = job.Args
		_job.StdOut = job.StdOut
		_job.StdErr = job.StdErr
		_job.Spec = job.Spec
		_job.Timeout = job.Timeout
		if job.IsMonitor == 1 {
			_job.IsMonitor = true
		} else {
			_job.IsMonitor = false
		}
		list = append(list, _job)
	}
	return &rpc.JobListResponse{Code: 200, Jobs: list}, nil
}

func (s *MasterServer) AgentMonitorReport(ctx context.Context, request *rpc.AgentMonitor) (*rpc.Response, error) {
	agentService := services.NewAgentService()
	agent := agentService.GetAgentByIPAndPort(request.GetAgent().GetIp(), request.GetAgent().GetPort())
	if agent == nil {
		return s.Error("no such agent!")
	}
	monitorService := services.NewMonitorService()
	monitor := &models.Monitor{
		Type:      models.TYPE_AGENT,
		RelatedID: agent.ID,
		UUID:      request.GetMonitor().GetUuid(),
		PID:       int64(request.GetMonitor().GetPid()),
		CPU:       float64(request.GetMonitor().GetCpu()),
		Memory:    float64(request.GetMonitor().GetMemory()),
		CreatedAt: time.Now(),
	}
	ok := monitorService.CreateMonitor(monitor)
	if !ok {
		return s.Error("add agent monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) AppMonitorReport(ctx context.Context, request *rpc.AppMonitor) (*rpc.Response, error) {
	monitorService := services.NewMonitorService()
	monitor := &models.Monitor{
		Type:      models.TYPE_APP,
		RelatedID: request.GetApp().GetId(),
		UUID:      request.GetMonitor().GetUuid(),
		PID:       int64(request.GetMonitor().GetPid()),
		CPU:       float64(request.GetMonitor().GetCpu()),
		Memory:    float64(request.GetMonitor().GetMemory()),
		CreatedAt: time.Now(),
	}
	ok := monitorService.CreateMonitor(monitor)
	if !ok {
		return s.Error("add app monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) JobMoniorReport(ctx context.Context, request *rpc.JobMonior) (*rpc.Response, error) {
	monitorService := services.NewMonitorService()
	monitor := &models.Monitor{
		Type:      models.TYPE_JOB,
		RelatedID: request.GetJob().GetId(),
		UUID:      request.GetMonitor().GetUuid(),
		PID:       int64(request.GetMonitor().GetPid()),
		CPU:       float64(request.GetMonitor().GetCpu()),
		Memory:    float64(request.GetMonitor().GetMemory()),
		CreatedAt: time.Now(),
	}
	ok := monitorService.CreateMonitor(monitor)
	if !ok {
		return s.Error("add job monitor failed")
	}
	return s.OK()
}

func (s *MasterServer) AppArchiveReport(ctx context.Context, request *rpc.AppArchive) (*rpc.Response, error) {
	archiveService := services.NewArchiveService()
	archive := &models.Archive{
		Type:      models.TYPE_APP,
		RelatedID: request.GetApp().GetId(),
		UUID:      request.GetArchive().GetUuid(),
		PID:       int64(request.GetArchive().GetPid()),
		BeginTime: time.Unix(request.GetArchive().GetBeginTime(), 0),
		EndTime:   time.Unix(request.GetArchive().GetEndTime(), 0),
		Status:    int64(request.GetArchive().GetStatus()),
		Signal:    request.GetArchive().GetSignal(),
		CreatedAt: time.Now(),
	}
	ok := archiveService.CreateArchive(archive)
	if !ok {
		return s.Error("add app archive failed")
	}
	return s.OK()
}

func (s *MasterServer) JobArchiveReport(ctx context.Context, request *rpc.JobArchive) (*rpc.Response, error) {
	archiveService := services.NewArchiveService()
	archive := &models.Archive{
		Type:      models.TYPE_JOB,
		RelatedID: request.GetJob().GetId(),
		UUID:      request.GetArchive().GetUuid(),
		PID:       int64(request.GetArchive().GetPid()),
		BeginTime: time.Unix(request.GetArchive().GetBeginTime(), 0),
		EndTime:   time.Unix(request.GetArchive().GetEndTime(), 0),
		Status:    int64(request.GetArchive().GetStatus()),
		Signal:    request.GetArchive().GetSignal(),
		CreatedAt: time.Now(),
	}
	ok := archiveService.CreateArchive(archive)
	if !ok {
		return s.Error("add job archive failed")
	}
	return s.OK()
}
