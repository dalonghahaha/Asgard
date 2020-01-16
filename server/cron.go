package server

import (
	"context"
	"fmt"

	"Asgard/applications"
	"Asgard/rpc"
)

type CronServer struct {
	baseServer
}

func (s *CronServer) List(ctx context.Context, request *rpc.Empty) (*rpc.JobListResponse, error) {
	list := []*rpc.Job{}
	for _, job := range applications.Jobs {
		list = append(list, rpc.BuildJob(job))
	}
	return &rpc.JobListResponse{Code: rpc.OK, Jobs: list}, nil
}

func (s *CronServer) Get(ctx context.Context, request *rpc.Name) (*rpc.JobResponse, error) {
	for _, job := range applications.Jobs {
		if request.GetName() == job.Name {
			return &rpc.JobResponse{Code: rpc.OK, Job: rpc.BuildJob(job)}, nil
		}
	}
	return &rpc.JobResponse{Code: rpc.Nofound, Job: nil}, nil
}

func (s *CronServer) Add(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
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

func (s *CronServer) Update(ctx context.Context, request *rpc.Job) (*rpc.Response, error) {
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

func (s *CronServer) Remove(ctx context.Context, request *rpc.Name) (*rpc.Response, error) {
	for _, job := range applications.Jobs {
		if request.GetName() == job.Name {
			ok := applications.JobStopByID(job.ID)
			if !ok {
				return s.Error(fmt.Sprintf("job %s stop failed", request.GetName()))
			}
		}
	}
	return s.OK()
}
