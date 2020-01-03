package server

import (
	"context"

	"Asgard/applications"
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
	return s.OK()
}

func (s *CronServer) Update(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	return s.OK()
}

func (s *CronServer) Remove(ctx context.Context, request *rpc.JobNameRequest) (*rpc.Response, error) {
	return s.OK()
}
