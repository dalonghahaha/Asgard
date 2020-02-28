package server

import (
	"context"
	"fmt"

	"Asgard/applications"
	"Asgard/rpc"
)

type TimerServer struct {
	baseServer
}

func (s *TimerServer) List(ctx context.Context, request *rpc.Empty) (*rpc.TimingListResponse, error) {
	list := []*rpc.Timing{}
	for _, timing := range applications.Timings {
		list = append(list, rpc.BuildTiming(timing))
	}
	return &rpc.TimingListResponse{Code: rpc.OK, Timings: list}, nil
}

func (s *TimerServer) Get(ctx context.Context, request *rpc.Name) (*rpc.TimingResponse, error) {
	for _, timing := range applications.Timings {
		if request.GetName() == timing.Name {
			return &rpc.TimingResponse{Code: rpc.OK, Timing: rpc.BuildTiming(timing)}, nil
		}
	}
	return &rpc.TimingResponse{Code: rpc.Nofound, Timing: nil}, nil
}

func (s *TimerServer) Add(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
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

func (s *TimerServer) Update(ctx context.Context, request *rpc.Timing) (*rpc.Response, error) {
	id := request.GetId()
	timing, ok := applications.Timings[id]
	if !ok {
		return s.Error(fmt.Sprintf("no job %d", id))
	}
	err := UpdateTiming(id, timing, request)
	if err != nil {
		return s.Error(err.Error())
	}
	return s.OK()
}

func (s *TimerServer) Remove(ctx context.Context, request *rpc.Name) (*rpc.Response, error) {
	for _, timing := range applications.Timings {
		if request.GetName() == timing.Name {
			err := DeleteTiming(timing.ID, timing)
			if err != nil {
				return s.Error(err.Error())
			}
			return s.OK()
		}
	}
	return s.OK()
}
