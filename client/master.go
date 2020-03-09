package client

import (
	"context"
	"fmt"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

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

func AgentMonitorReport(agentMonitor *rpc.AgentMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AgentMonitorReport(ctx, agentMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("agent moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("agent moniter report failed：%s", response.GetMessage()))
		return
	}
}

func AppMonitorReport(appMonitor *rpc.AppMonitor) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AppMonitorReport(ctx, appMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", response.GetMessage()))
		return
	}
}

func JobMonitorReport(jobMonitor *rpc.JobMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.JobMoniorReport(ctx, jobMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", response.GetMessage()))
		return
	}
}

func TimingMonitorReport(timingMonitor *rpc.TimingMonior) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.TimingMoniorReport(ctx, timingMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("timing moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("timing moniter report failed：%s", response.GetMessage()))
		return
	}
}

func AppArchiveReport(appArchive *rpc.AppArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AppArchiveReport(ctx, appArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("app archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app archive report failed：%s", response.GetMessage()))
		return
	}
}

func JobArchiveReport(jobArchive *rpc.JobArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.JobArchiveReport(ctx, jobArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("job archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job archive report failed：%s", response.GetMessage()))
		return
	}
}

func TimingArchiveReport(timingArchive *rpc.TimingArchive) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.TimingArchiveReport(ctx, timingArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("timing archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("timing archive report failed：%s", response.GetMessage()))
		return
	}
}
