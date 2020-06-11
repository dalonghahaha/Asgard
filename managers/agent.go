package managers

import (
	"Asgard/clients"
	"Asgard/constants"
	"Asgard/rpc"
	"Asgard/runtimes"
	"fmt"
	"syscall"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/shirou/gopsutil/process"
)

type AgentManager struct {
	appManager    *AppManager
	jobManager    *JobManager
	timingManager *TimingManager
	masterClient  *clients.Master
	ticker        *time.Ticker
	exitSingel    chan bool
}

func NewAgentManager() (*AgentManager, error) {
	masterClient, err := clients.NewMaster(constants.MASTER_IP, constants.MASTER_PORT)
	if err != nil {
		return nil, fmt.Errorf("init master client failed:%s", err.Error())
	}
	appManager := NewAppManager()
	appManager.SetMaster(masterClient)
	jobManager := NewJobManager()
	jobManager.SetMaster(masterClient)
	timingManager := NewTimingManager()
	timingManager.SetMaster(masterClient)
	manager := &AgentManager{
		appManager:    appManager,
		jobManager:    jobManager,
		timingManager: timingManager,
		masterClient:  masterClient,
		ticker:        time.NewTicker(time.Second * time.Duration(constants.AGENT_MONITER)),
		exitSingel:    make(chan bool, 1),
	}
	return manager, nil
}

func (a *AgentManager) GetAppManager() *AppManager {
	return a.appManager
}

func (a *AgentManager) GetJobManager() *JobManager {
	return a.jobManager
}

func (a *AgentManager) GetTimingManager() *TimingManager {
	return a.timingManager
}

func (a *AgentManager) StartMonitor() {
	logger.Debug("agent monitor ticker start!")
	runtimes.SubscribeExit(a.exitSingel)
	for {
		select {
		case <-a.exitSingel:
			logger.Debug("agent monitor ticker stop!")
			a.ticker.Stop()
			break
		case <-a.ticker.C:
			go a.AgentMonitorReport()
		}
	}
}

func (a *AgentManager) AgentMonitorReport() {
	info, err := process.NewProcess(int32(constants.AGENT_PID))
	if err != nil {
		logger.Error("get process failed:", err)
		return
	}
	agentMonitor := runtimes.AgentMonitor{
		Ip:      constants.AGENT_IP,
		Port:    constants.AGENT_PORT,
		Monitor: runtimes.BuildMonitorInfo(info),
	}
	a.masterClient.AgentMonitorChan <- agentMonitor
}

func (a *AgentManager) StartAll() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("anget manager start failed :", err)
			runtimes.ExitSinal <- syscall.SIGTERM
			return
		}
	}()
	go a.masterClient.Report()
	go a.StartMonitor()
	err := a.masterClient.AgentRegister()
	if err != nil {
		panic(err)
	}
	err = a.AppsRegister()
	if err != nil {
		panic(err)
	}
	err = a.JobsRegister()
	if err != nil {
		panic(err)
	}
	err = a.TimingsRegister()
	if err != nil {
		panic(err)
	}
	a.appManager.StartAll(true)
	a.jobManager.StartAll(true)
	a.timingManager.StartAll(true)
	logger.Info("Agent Started!")
	logger.Debugf("Agent Master: %s:%s", constants.MASTER_IP, constants.MASTER_PORT)
	logger.Debugf("Agent Address: %s:%s", constants.AGENT_IP, constants.AGENT_PORT)
	logger.Debugf("Agent Loop:%d", constants.AGENT_MONITER)
}

func (a *AgentManager) StopAll() {
	runtimes.Exit()
	a.appManager.StopAll()
	a.jobManager.StopAll()
	a.timingManager.StopAll()
	//make sure all data report to master
	maxWait := 10
	countWait := 0
	for {
		if a.masterClient.IsRunning() && countWait <= maxWait {
			time.Sleep(time.Second * 1)
			countWait += 1
		} else {
			break
		}
	}
	logger.Info("Agent Server Stop!")
}

func (a *AgentManager) AppsRegister() error {
	apps, err := a.masterClient.GetAppList()
	if err != nil {
		return err
	}
	for _, app := range apps {
		err := a.appManager.Register(app.GetId(), rpc.BuildAppConfig(app))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AgentManager) JobsRegister() error {
	jobs, err := a.masterClient.GetJobList()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		err := a.jobManager.Register(job.GetId(), rpc.BuildJobConfig(job))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AgentManager) TimingsRegister() error {
	timings, err := a.masterClient.GetTimingList()
	if err != nil {
		return err
	}
	for _, timing := range timings {
		err := a.timingManager.Register(timing.GetId(), rpc.BuildTimingConfig(timing))
		if err != nil {
			return err
		}
	}
	return nil
}
