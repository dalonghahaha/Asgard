package server

import (
	"context"
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"

	"Asgard/models"
	"Asgard/rpc"
	"Asgard/services"
)

type MasterServer struct {
	baseServer
}

func (s *MasterServer) Register(ctx context.Context, request *rpc.Agent) (*rpc.Response, error) {
	logger.Debug(fmt.Sprintf("agent %s:%s joined!", request.GetIp(), request.GetPort()))
	agentService := services.NewAgentService()
	agent := agentService.GetAgentByIP(request.GetIp())
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
