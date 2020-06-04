package server

import (
	"context"

	"Asgard/managers"
	"Asgard/rpc"
)

type TimerServer struct {
	baseServer
	timingManager *managers.TimingManager
}

func (s *TimerServer) SetTimingManager(timingManager *managers.TimingManager) {
	s.timingManager = timingManager
}

func (s *TimerServer) List(ctx context.Context, request *rpc.Empty) (*rpc.TimingListResponse, error) {
	timings := s.timingManager.GetList()
	list := []*rpc.Timing{}
	for _, timing := range timings {
		if timing.Executed {
			continue
		}
		list = append(list, rpc.BuildTiming(timing))
	}
	return &rpc.TimingListResponse{Code: rpc.OK, Timings: list}, nil
}

func (s *TimerServer) Get(ctx context.Context, request *rpc.Name) (*rpc.TimingResponse, error) {
	timing := s.timingManager.GetByName(request.GetName())
	if timing != nil {
		return &rpc.TimingResponse{Code: rpc.OK, Timing: rpc.BuildTiming(timing)}, nil
	}

	return &rpc.TimingResponse{Code: rpc.Nofound, Timing: nil}, nil
}

func (s *TimerServer) Add(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
	err := s.timingManager.Register(request.GetId(), rpc.BuildTimingConfig(request))
	if err == nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *TimerServer) Update(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
	err := s.timingManager.Update(request.GetId(), rpc.BuildTimingConfig(request))
	if err == nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *TimerServer) Remove(ctx context.Context, request *rpc.Name) (*rpc.Response, error) {
	timing := s.timingManager.GetByName(request.GetName())
	if timing == nil {
		return s.OK()
	}
	err := s.timingManager.Remove(timing.ID)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}
