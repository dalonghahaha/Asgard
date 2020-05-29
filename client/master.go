package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/dalonghahaha/avenger/components/logger"
	"google.golang.org/grpc"

	"Asgard/applications"
	"Asgard/constants"
	"Asgard/rpc"
)

type Master struct {
	Reports           *sync.Map
	agent             *rpc.AgentInfo
	rpcClient         rpc.MasterClient
	AgentMonitorChan  chan applications.AgentMonitor
	AppMonitorChan    chan applications.AppMonitor
	AppArchiveChan    chan applications.AppArchive
	JobMonitorChan    chan applications.JobMonitor
	JobArchiveChan    chan applications.JobArchive
	TimingMonitorChan chan applications.TimingMonitor
	TimingArchiveChan chan applications.TimingArchive
}

func NewMaster(ip, port string) (*Master, error) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	option := grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(constants.RPC_MESSAGE_SIZE),
		grpc.MaxCallSendMsgSize(constants.RPC_MESSAGE_SIZE),
	)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), option)
	if err != nil {
		return nil, err
	}
	master := Master{
		agent: &rpc.AgentInfo{
			Ip:   constants.AGENT_IP,
			Port: constants.AGENT_PORT,
		},
		rpcClient:         rpc.NewMasterClient(conn),
		Reports:           new(sync.Map),
		AgentMonitorChan:  make(chan applications.AgentMonitor, 100),
		AppMonitorChan:    make(chan applications.AppMonitor, 100),
		AppArchiveChan:    make(chan applications.AppArchive, 100),
		JobMonitorChan:    make(chan applications.JobMonitor, 100),
		JobArchiveChan:    make(chan applications.JobArchive, 100),
		TimingMonitorChan: make(chan applications.TimingMonitor, 100),
		TimingArchiveChan: make(chan applications.TimingArchive, 100),
	}
	return &master, nil
}

func (m *Master) IsRunning() bool {
	count := 0
	m.Reports.Range(func(k, v interface{}) bool {
		count += 1
		return true
	})
	return count > 0
}

func (m *Master) Report() {
	logger.Debug("master report litsen!")
	for {
		select {
		case agentMonitor := <-m.AgentMonitorChan:
			m.agentMonitorReport(rpc.BuildAgentMonitor(agentMonitor.Ip, agentMonitor.Port, agentMonitor.Monitor))
		case appMonitor := <-m.AppMonitorChan:
			m.appMonitorReport(rpc.BuildAppMonitor(appMonitor.App, appMonitor.Monitor))
		case appArchive := <-m.AppArchiveChan:
			m.appArchiveReport(rpc.BuildAppArchive(appArchive.App, appArchive.Archive))
			m.Reports.Delete(appArchive.UUID)
			logger.Debug("appArchive Report: ", appArchive.App.Name)
		case jobMonitor := <-m.JobMonitorChan:
			m.jobMonitorReport(rpc.BuildJobMonior(jobMonitor.Job, jobMonitor.Monitor))
		case jobArchive := <-m.JobArchiveChan:
			m.jobArchiveReport(rpc.BuildJobArchive(jobArchive.Job, jobArchive.Archive))
		case timingMonitor := <-m.TimingMonitorChan:
			m.timingMonitorReport(rpc.BuildTimingMonior(timingMonitor.Timing, timingMonitor.Monitor))
		case timingArchive := <-m.TimingArchiveChan:
			m.timingArchiveReport(rpc.BuildTimingArchive(timingArchive.Timing, timingArchive.Archive))
		}
	}
}

func (m *Master) AgentRegister() error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
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
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
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
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
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
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
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

func (m *Master) agentMonitorReport(agentMonitor *rpc.AgentMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.AgentMonitorReport(ctx, agentMonitor)
	if err != nil {
		logger.Errorf("agent moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("agent moniter report failed：%s", response.GetMessage())
		return
	}
}

func (m *Master) appMonitorReport(appMonitor *rpc.AppMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.AppMonitorReport(ctx, appMonitor)
	if err != nil {
		logger.Errorf("app moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("app moniter report failed：%s", response.GetMessage())
		return
	}
}

func (m *Master) jobMonitorReport(jobMonitor *rpc.JobMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.JobMoniorReport(ctx, jobMonitor)
	if err != nil {
		logger.Errorf("job moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Error("job moniter report failed：%s", response.GetMessage())
		return
	}
}

func (m *Master) timingMonitorReport(timingMonitor *rpc.TimingMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.TimingMoniorReport(ctx, timingMonitor)
	if err != nil {
		logger.Errorf("timing moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("timing moniter report failed：%s", response.GetMessage())
		return
	}
}

func (m *Master) appArchiveReport(appArchive *rpc.AppArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.AppArchiveReport(ctx, appArchive)
	if err != nil {
		logger.Errorf("app archive report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Error("app archive report failed：%s", response.GetMessage())
		return
	}
}

func (m *Master) jobArchiveReport(jobArchive *rpc.JobArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.JobArchiveReport(ctx, jobArchive)
	if err != nil {
		logger.Errorf("job archive report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("job archive report failed：%s", response.GetMessage())
		return
	}
}

func (m *Master) timingArchiveReport(timingArchive *rpc.TimingArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := m.rpcClient.TimingArchiveReport(ctx, timingArchive)
	if err != nil {
		logger.Errorf("timing archive report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("timing archive report failed：%s", response.GetMessage())
		return
	}
}
