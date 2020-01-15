package server

import (
	"context"
	"fmt"
	"time"

	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
)

type CronServer struct {
	baseServer
}

func (s *CronServer) List(ctx context.Context, request *rpc.Empty) (*rpc.JobListResponse, error) {
	jobs := applications.Jobs
	list := []*rpc.Job{}
	for _, job := range jobs {
		_job := new(rpc.Job)
		_job.Id = job.ID
		_job.Name = job.Name
		_job.Dir = job.Dir
		_job.Program = job.Program
		_job.Args = job.Args
		_job.StdOut = job.Stdout
		_job.StdErr = job.Stderr
		_job.Spec = job.Spec
		_job.Timeout = job.TimeOut.Milliseconds()
		list = append(list, _job)
	}
	return &rpc.JobListResponse{Code: 200, Jobs: list}, nil
}

func (s *CronServer) Get(ctx context.Context, request *rpc.JobNameRequest) (*rpc.JobResponse, error) {
	jobs := applications.Jobs
	name := request.GetName()
	for _, job := range jobs {
		if name == job.Name {
			_job := new(rpc.Job)
			_job.Name = job.Name
			_job.Dir = job.Dir
			_job.Program = job.Program
			_job.Args = job.Args
			_job.StdOut = job.Stdout
			_job.StdErr = job.Stderr
			_job.Spec = job.Spec
			_job.Timeout = job.TimeOut.Milliseconds()
			return &rpc.JobResponse{Code: 200, Job: _job}, nil
		}
	}
	return &rpc.JobResponse{Code: 0, Job: nil}, nil
}

func (s *CronServer) Add(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	id := request.GetId()
	_, ok := applications.Jobs[id]
	if ok {
		return s.OK()
	}
	config := map[string]interface{}{
		"id":      request.GetId(),
		"name":    request.GetName(),
		"dir":     request.GetDir(),
		"program": request.GetProgram(),
		"args":    request.GetArgs(),
		"stdout":  request.GetStdOut(),
		"stderr":  request.GetStdErr(),
		"spec":    request.GetSpec(),
		"timeout": request.GetTimeout(),
	}
	job, err := applications.JobRegister(id, config)
	if err != nil {
		return s.Error(err.Error())
	}
	job.MonitorReport = func(monitor *applications.Monitor) {
		client.JobMonitorReport(job, monitor)
	}
	job.ArchiveReport = func(command *applications.Command) {
		client.JobArchiveReport(job, command)
	}
	ok = applications.JobStartByID(id)
	if !ok {
		return s.Error(fmt.Sprintf("job %d start failed", id))
	}
	return s.OK()
}

func (s *CronServer) Update(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	id := request.GetId()
	job, ok := applications.Jobs[id]
	if !ok {
		return s.Error(fmt.Sprintf("no job %d", id))
	}
	ok = applications.JobStopByID(id)
	if !ok {
		return s.Error(fmt.Sprintf("job %d stop failed", id))
	}
	job.Name = request.GetName()
	job.Dir = request.GetDir()
	job.Program = request.GetProgram()
	job.Args = request.GetArgs()
	job.Stdout = request.GetStdOut()
	job.Stderr = request.GetStdErr()
	job.Spec = request.GetSpec()
	job.TimeOut = time.Duration(request.GetTimeout())
	job.IsMonitor = request.GetIsMonitor()
	ok = applications.JobStartByID(id)
	if !ok {
		return s.Error(fmt.Sprintf("job %d start failed", id))
	}
	return s.OK()
}

func (s *CronServer) Remove(ctx context.Context, request *rpc.JobNameRequest) (*rpc.Response, error) {
	return s.OK()
}
