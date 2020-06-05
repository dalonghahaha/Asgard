package clients

import (
	"context"
	"fmt"
	"sync"

	"github.com/dalonghahaha/avenger/components/logger"
	"google.golang.org/grpc"

	"Asgard/constants"
	"Asgard/rpc"
	"Asgard/runtimes"
)

type Master struct {
	Reports   *sync.Map
	agent     *rpc.AgentInfo
	rpcClient rpc.MasterClient

	//monitor chans
	AgentMonitorChan  chan runtimes.AgentMonitor
	AppMonitorChan    chan runtimes.AppMonitor
	JobMonitorChan    chan runtimes.JobMonitor
	TimingMonitorChan chan runtimes.TimingMonitor

	//archive chans
	AppArchiveChan    chan runtimes.AppArchive
	JobArchiveChan    chan runtimes.JobArchive
	TimingArchiveChan chan runtimes.TimingArchive

	//exception chans
	AppExceptionChan    chan runtimes.AppException
	JobExceptionChan    chan runtimes.JobException
	TimingExceptionChan chan runtimes.TimingException
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
		rpcClient: rpc.NewMasterClient(conn),
		Reports:   new(sync.Map),

		AgentMonitorChan:  make(chan runtimes.AgentMonitor, 100),
		JobMonitorChan:    make(chan runtimes.JobMonitor, 100),
		TimingMonitorChan: make(chan runtimes.TimingMonitor, 100),

		AppMonitorChan:    make(chan runtimes.AppMonitor, 100),
		AppArchiveChan:    make(chan runtimes.AppArchive, 100),
		JobArchiveChan:    make(chan runtimes.JobArchive, 100),
		TimingArchiveChan: make(chan runtimes.TimingArchive, 100),

		AppExceptionChan:    make(chan runtimes.AppException, 100),
		JobExceptionChan:    make(chan runtimes.JobException, 100),
		TimingExceptionChan: make(chan runtimes.TimingException, 100),
	}
	return &master, nil
}

func (c *Master) IsRunning() bool {
	count := 0
	c.Reports.Range(func(k, v interface{}) bool {
		count += 1
		return true
	})
	return count > 0
}

func (c *Master) Report() {
	logger.Debug("master report litsen!")
	for {
		select {

		//monitor report
		case m := <-c.AgentMonitorChan:
			c.agentMonitorReport(rpc.BuildAgentMonitor(m.Ip, m.Port, m.Monitor))
			c.Reports.Delete(m.UUID)
		case m := <-c.AppMonitorChan:
			c.appMonitorReport(rpc.BuildAppMonitor(m.App, m.Monitor))
			c.Reports.Delete(m.UUID)
		case m := <-c.JobMonitorChan:
			c.jobMonitorReport(rpc.BuildJobMonior(m.Job, m.Monitor))
			c.Reports.Delete(m.UUID)
		case m := <-c.TimingMonitorChan:
			c.timingMonitorReport(rpc.BuildTimingMonior(m.Timing, m.Monitor))
			c.Reports.Delete(m.UUID)

		//archive report
		case a := <-c.AppArchiveChan:
			c.appArchiveReport(rpc.BuildAppArchive(a.App, a.Archive))
			c.Reports.Delete(a.UUID)
		case a := <-c.JobArchiveChan:
			c.jobArchiveReport(rpc.BuildJobArchive(a.Job, a.Archive))
			c.Reports.Delete(a.UUID)
		case a := <-c.TimingArchiveChan:
			c.timingArchiveReport(rpc.BuildTimingArchive(a.Timing, a.Archive))
			c.Reports.Delete(a.UUID)

		//exception report
		case e := <-c.AppExceptionChan:
			c.appExceptionReport(rpc.BuildAppException(e))
			c.Reports.Delete(e.UUID)
		case e := <-c.JobExceptionChan:
			c.jobExceptionReport(rpc.BuildJobException(e))
			c.Reports.Delete(e.UUID)
		case e := <-c.TimingExceptionChan:
			c.timingExceptionReport(rpc.BuildTimingException(e))
			c.Reports.Delete(e.UUID)
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

func (c *Master) GetTimingList() ([]*rpc.Timing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.TimingList(ctx, c.agent)
	if err != nil {
		return nil, fmt.Errorf("get job list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get job list error: agent error")
	}
	return response.GetTimings(), nil
}

func (c *Master) agentMonitorReport(agentMonitor *rpc.AgentMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.AgentMonitorReport(ctx, agentMonitor)
	if err != nil {
		logger.Errorf("agent moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("agent moniter report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) appMonitorReport(appMonitor *rpc.AppMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.AppMonitorReport(ctx, appMonitor)
	if err != nil {
		logger.Errorf("app moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("app moniter report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) jobMonitorReport(jobMonitor *rpc.JobMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.JobMoniorReport(ctx, jobMonitor)
	if err != nil {
		logger.Errorf("job moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("job moniter report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) timingMonitorReport(timingMonitor *rpc.TimingMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.TimingMoniorReport(ctx, timingMonitor)
	if err != nil {
		logger.Errorf("timing moniter report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("timing moniter report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) appArchiveReport(appArchive *rpc.AppArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.AppArchiveReport(ctx, appArchive)
	if err != nil {
		logger.Errorf("app archive report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("app archive report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) jobArchiveReport(jobArchive *rpc.JobArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.JobArchiveReport(ctx, jobArchive)
	if err != nil {
		logger.Errorf("job archive report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("job archive report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) timingArchiveReport(timingArchive *rpc.TimingArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.TimingArchiveReport(ctx, timingArchive)
	if err != nil {
		logger.Errorf("timing archive report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("timing archive report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) appExceptionReport(appException *rpc.AppException) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.AppExceptionReport(ctx, appException)
	if err != nil {
		logger.Errorf("app exception report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("app exception report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) jobExceptionReport(jobException *rpc.JobException) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.JobExceptionReport(ctx, jobException)
	if err != nil {
		logger.Errorf("job exception report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("job exception report failed：%s", response.GetMessage())
		return
	}
}

func (c *Master) timingExceptionReport(timingException *rpc.TimingException) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.RPC_TIMEOUT)
	defer cancel()
	response, err := c.rpcClient.TimingExceptionReport(ctx, timingException)
	if err != nil {
		logger.Errorf("timing exception report failed：%s", err.Error())
		return
	}
	if response.GetCode() != 200 {
		logger.Errorf("timing exception report failed：%s", response.GetMessage())
		return
	}
}
