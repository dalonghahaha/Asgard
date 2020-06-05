package server

import (
	"context"

	"Asgard/managers"
	"Asgard/rpc"
)

type CronServer struct {
	baseServer
	jobManager *managers.JobManager
}

func NewCronServer(jobManager *managers.JobManager) *CronServer {
	return &CronServer{
		jobManager: jobManager,
	}
}

func (s *CronServer) SetJobManager(jobManager *managers.JobManager) {
	s.jobManager = jobManager
}

func (s *CronServer) List(ctx context.Context, request *rpc.Empty) (*rpc.JobListResponse, error) {
	jobs := s.jobManager.GetList()
	list := []*rpc.Job{}
	for _, job := range jobs {
		list = append(list, rpc.BuildJob(job))
	}
	return &rpc.JobListResponse{Code: rpc.OK, Jobs: list}, nil
}

func (s *CronServer) Get(ctx context.Context, request *rpc.Name) (*rpc.JobResponse, error) {
	job := s.jobManager.GetByName(request.GetName())
	if job != nil {
		return &rpc.JobResponse{Code: rpc.OK, Job: rpc.BuildJob(job)}, nil
	}
	return &rpc.JobResponse{Code: rpc.Nofound, Job: nil}, nil
}

func (s *CronServer) Add(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	err := s.jobManager.Add(request.GetId(), rpc.BuildJobConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *CronServer) Update(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
	err := s.jobManager.Update(request.GetId(), rpc.BuildJobConfig(request))
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *CronServer) Remove(ctx context.Context, request *rpc.Name) (*rpc.Response, error) {
	job := s.jobManager.GetByName(request.GetName())
	if job == nil {
		return s.OK()
	}
	err := s.jobManager.Remove(job.ID)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}
