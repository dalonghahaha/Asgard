package client

import (
	"context"
	"fmt"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"Asgard/applications"
	"Asgard/rpc"
)

var (
	MasterClient rpc.MasterClient
	TimeOut      = time.Second * 10
)

func InitMasterClient() {
	masterIP := viper.GetString("agent.master.ip")
	masterPort := viper.GetString("agent.master.port")
	addr := fmt.Sprintf("%s:%s", masterIP, masterPort)
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeOut)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		panic("Can't connect: " + addr)
	}
	MasterClient = rpc.NewMasterClient(conn)
}

func AgentRegister(agentIP, agentPort string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	logger.Debug(fmt.Sprintf("agent register：%s:%s", agentIP, agentPort))
	response, err := MasterClient.Register(ctx, &rpc.AgentInfo{Ip: agentIP, Port: agentPort})
	if err != nil {
		return fmt.Errorf("agent register fail: %v", err.Error())
	}
	if response.GetCode() != 200 {
		return fmt.Errorf("agent register fail: %s", response.GetMessage())
	}
	return nil
}

func GetAppList(agentIP, agentPort string) ([]*rpc.App, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AppList(ctx, &rpc.AgentInfo{Ip: agentIP, Port: agentPort})
	if err != nil {
		return nil, fmt.Errorf("get app list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get app list error: agent error")
	}
	return response.GetApps(), nil
}

func GetJobList(agentIP, agentPort string) ([]*rpc.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.JobList(ctx, &rpc.AgentInfo{Ip: agentIP, Port: agentPort})
	if err != nil {
		return nil, fmt.Errorf("get job list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get job list error: agent error")
	}
	return response.GetJobs(), nil
}

func GetTimingList(agentIP, agentPort string) ([]*rpc.Timing, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.TimingList(ctx, &rpc.AgentInfo{Ip: agentIP, Port: agentPort})
	if err != nil {
		return nil, fmt.Errorf("get job list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get job list error: agent error")
	}
	return response.GetTimings(), nil
}

func AgentMonitorReport(agentIP string, agentPort string, Pid int, UUID string, monitor *applications.Monitor) {
	agentMonitor := rpc.AgentMonitor{
		Agent: &rpc.AgentInfo{
			Ip:   agentIP,
			Port: agentPort,
		},
		Monitor: &rpc.Monitor{
			Uuid:    UUID,
			Pid:     int32(Pid),
			Cpu:     float32(monitor.CPUPercent),
			Memory:  monitor.MemoryPercent,
			Threads: int32(monitor.NumThreads),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AgentMonitorReport(ctx, &agentMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("agent moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("agent moniter report failed：%s", response.GetMessage()))
		return
	}
}

func AppMonitorReport(app *applications.App, monitor *applications.Monitor) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AppMonitorReport(ctx, rpc.BuildAppMonitor(app, monitor))
	if err != nil {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", response.GetMessage()))
		return
	}
}

func JobMonitorReport(job *applications.Job, monitor *applications.Monitor) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.JobMoniorReport(ctx, rpc.BuildJobMonior(job, monitor))
	if err != nil {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", response.GetMessage()))
		return
	}
}

func TimingMonitorReport(timing *applications.Timing, monitor *applications.Monitor) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.TimingMoniorReport(ctx, rpc.BuildTimingMonior(timing, monitor))
	if err != nil {
		logger.Error(fmt.Sprintf("timing moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("timing moniter report failed：%s", response.GetMessage()))
		return
	}
}

func AppArchiveReport(app *applications.App, command *applications.Command) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AppArchiveReport(ctx, rpc.BuildAppArchive(app, command))
	if err != nil {
		logger.Error(fmt.Sprintf("app archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app archive report failed：%s", response.GetMessage()))
		return
	}
}

func JobArchiveReport(job *applications.Job, command *applications.Command) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.JobArchiveReport(ctx, rpc.BuildJobArchive(job, command))
	if err != nil {
		logger.Error(fmt.Sprintf("job archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job archive report failed：%s", response.GetMessage()))
		return
	}
}

func TimingArchiveReport(job *applications.Timing, command *applications.Command) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.TimingArchiveReport(ctx, rpc.BuildTimingArchive(job, command))
	if err != nil {
		logger.Error(fmt.Sprintf("timing archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("timing archive report failed：%s", response.GetMessage()))
		return
	}
}
