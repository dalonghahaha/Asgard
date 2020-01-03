package server

import (
	"context"
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"

	"Asgard/rpc"
)

type MasterServer struct {
	baseServer
}

func (s *MasterServer) Register(ctx context.Context, request *rpc.Agent) (*rpc.Response, error) {
	logger.Debug(fmt.Sprintf("agent %s:%s joined!", request.GetIp(), request.GetPort()))
	return s.OK()
}

func (s *MasterServer) AppMonitorReport(ctx context.Context, request *rpc.AppMonitor) (*rpc.Response, error) {
	return s.OK()
}

func (s *MasterServer) JobMoniorReport(ctx context.Context, request *rpc.JobMonior) (*rpc.Response, error) {
	return s.OK()
}

func (s *MasterServer) AppExceptionReport(ctx context.Context, request *rpc.AppException) (*rpc.Response, error) {
	return s.OK()
}
func (s *MasterServer) JobExceptionReport(ctx context.Context, request *rpc.JobException) (*rpc.Response, error) {
	return s.OK()
}
