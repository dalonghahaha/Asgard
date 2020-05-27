package client

import (
	"context"
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"google.golang.org/grpc"

	"Asgard/constants"
	"Asgard/rpc"
)

type Master struct {
	agent     *rpc.AgentInfo
	rpcClient rpc.MasterClient
}

func NewMaster() *Master {
	addr := fmt.Sprintf("%s:%d", constants.MASTER_IP, constants.MASTER_PORT)
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		panic("can't connect master: " + addr)
	}
	logger.Debug("master connected!")
	return &Master{
		agent: &rpc.AgentInfo{
			Ip:   constants.AGENT_IP,
			Port: constants.AGENT_PORT,
		},
		rpcClient: rpc.NewMasterClient(conn),
	}
}

func (m *Master) AgentRegister() error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.Register(ctx, m.agent)
	if err != nil {
		return fmt.Errorf("agent register fail: %s", err.Error())
	}
	if response.GetCode() != 200 {
		return fmt.Errorf("agent register fail: %s", response.GetMessage())
	}
	logger.Debug("agent register success!")
	return nil
}

func (m *Master) GetAppList() ([]*rpc.App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.AppList(ctx, m.agent)
	if err != nil {
		return nil, fmt.Errorf("get app list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get app list error: agent error")
	}
	return response.GetApps(), nil
}

func (m *Master) GetJobList() ([]*rpc.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.JobList(ctx, m.agent)
	if err != nil {
		return nil, fmt.Errorf("get job list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get job list error: agent error")
	}
	return response.GetJobs(), nil
}

func (m *Master) GetTimingList() ([]*rpc.Timing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.TimingList(ctx, m.agent)
	if err != nil {
		return nil, fmt.Errorf("get job list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get job list error: agent error")
	}
	return response.GetTimings(), nil
}

func (m *Master) AgentMonitorReport(agentMonitor *rpc.AgentMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.AgentMonitorReport(ctx, agentMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("agent moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("agent moniter report failed：%s", response.GetMessage()))
		return
	}
}

func (m *Master) AppMonitorReport(appMonitor *rpc.AppMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.AppMonitorReport(ctx, appMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", response.GetMessage()))
		return
	}
}

func (m *Master) JobMonitorReport(jobMonitor *rpc.JobMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.JobMoniorReport(ctx, jobMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", response.GetMessage()))
		return
	}
}

func (m *Master) TimingMonitorReport(timingMonitor *rpc.TimingMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.TimingMoniorReport(ctx, timingMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("timing moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("timing moniter report failed：%s", response.GetMessage()))
		return
	}
}

func (m *Master) AppArchiveReport(appArchive *rpc.AppArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.AppArchiveReport(ctx, appArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("app archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app archive report failed：%s", response.GetMessage()))
		return
	}
}

func (m *Master) JobArchiveReport(jobArchive *rpc.JobArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.JobArchiveReport(ctx, jobArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("job archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job archive report failed：%s", response.GetMessage()))
		return
	}
}

func (m *Master) TimingArchiveReport(timingArchive *rpc.TimingArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.MASTER_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.TimingArchiveReport(ctx, timingArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("timing archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("timing archive report failed：%s", response.GetMessage()))
		return
	}
}
