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
	TimeOut      = time.Second * 30
)

func InitMasterClient() {
	masterIP := viper.GetString("agent.master.ip")
	masterPort := viper.GetString("agent.master.port")
	addr := fmt.Sprintf("%s:%s", masterIP, masterPort)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic("Can't connect: " + addr)
	}
	MasterClient = rpc.NewMasterClient(conn)
}

func buildApp(app *applications.App) *rpc.App {
	return &rpc.App{
		Id:          app.ID,
		Name:        app.Name,
		Dir:         app.Dir,
		Program:     app.Program,
		Args:        app.Args,
		StdOut:      app.Stdout,
		StdErr:      app.Stderr,
		AutoRestart: app.AutoRestart,
		IsMonitor:   app.IsMonitor,
	}
}

func buildJob(job *applications.Job) *rpc.Job {
	return &rpc.Job{
		Id:        job.ID,
		Name:      job.Name,
		Dir:       job.Dir,
		Program:   job.Program,
		Args:      job.Args,
		StdOut:    job.Stdout,
		StdErr:    job.Stderr,
		Spec:      job.Spec,
		Timeout:   int64(job.TimeOut),
		IsMonitor: job.IsMonitor,
	}
}

func buildArchive(command *applications.Command) *rpc.Archive {
	return &rpc.Archive{
		Uuid:      command.UUID,
		Pid:       int32(command.Pid),
		BeginTime: command.Begin.Unix(),
		EndTime:   command.End.Unix(),
		Status:    int32(command.Status),
		Signal:    command.Signal,
	}
}

func AgentRegister(agentIP, agentPort string) error {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	logger.Debug(fmt.Sprintf("agent register：%s:%s", agentIP, agentPort))
	response, err := MasterClient.Register(ctx, &rpc.Agent{Ip: agentIP, Port: agentPort})
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
	logger.Debug(fmt.Sprintf("get app list：%s:%s", agentIP, agentPort))
	response, err := MasterClient.AppList(ctx, &rpc.Agent{Ip: agentIP, Port: agentPort})
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
	logger.Debug(fmt.Sprintf("get job list：%s:%s", agentIP, agentPort))
	response, err := MasterClient.JobList(ctx, &rpc.Agent{Ip: agentIP, Port: agentPort})
	if err != nil {
		return nil, fmt.Errorf("get job list error: %v", err.Error())
	}
	if response.GetCode() == 404 {
		return nil, fmt.Errorf("get job list error: agent error")
	}
	return response.GetJobs(), nil
}

func AppMonitorReport(app *applications.App, monitor *applications.Monitor) {
	appMonitor := rpc.AppMonitor{
		App: buildApp(app),
		Monitor: &rpc.Monitor{
			Uuid:    app.UUID,
			Pid:     int32(app.Pid),
			Cpu:     float32(monitor.CPUPercent),
			Memory:  monitor.MemoryPercent,
			Threads: int32(monitor.NumThreads),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AppMonitorReport(ctx, &appMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app moniter report failed：%s", response.GetMessage()))
		return
	}
	logger.Debug("app moniter report success")
}

func JobMonitorReport(job *applications.Job, monitor *applications.Monitor) {
	jobMonitor := rpc.JobMonior{
		Job: buildJob(job),
		Monitor: &rpc.Monitor{
			Uuid:    job.UUID,
			Pid:     int32(job.Pid),
			Cpu:     float32(monitor.CPUPercent),
			Memory:  monitor.MemoryPercent,
			Threads: int32(monitor.NumThreads),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.JobMoniorReport(ctx, &jobMonitor)
	if err != nil {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job moniter report failed：%s", response.GetMessage()))
		return
	}
	logger.Debug("job moniter report success")
}

func AppArchiveReport(app *applications.App, command *applications.Command) {
	appArchive := rpc.AppArchive{
		App:     buildApp(app),
		Archive: buildArchive(command),
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.AppArchiveReport(ctx, &appArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("app archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("app archive report failed：%s", response.GetMessage()))
		return
	}
	logger.Debug("app archive report success")
}

func JobArchiveReport(job *applications.Job, command *applications.Command) {
	jobArchive := rpc.JobArchive{
		Job:     buildJob(job),
		Archive: buildArchive(command),
	}
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	response, err := MasterClient.JobArchiveReport(ctx, &jobArchive)
	if err != nil {
		logger.Error(fmt.Sprintf("job archive report failed：%s", err.Error()))
		return
	}
	if response.GetCode() != 200 {
		logger.Error(fmt.Sprintf("job archive report failed：%s", response.GetMessage()))
		return
	}
	logger.Debug("job archive report success")
}
