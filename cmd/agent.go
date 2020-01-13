package cmd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"Asgard/applications"
	"Asgard/rpc"
	"Asgard/server"
)

var (
	masterClient rpc.MasterClient
	agentIP      string
	agentPort    string
)

func init() {
	agentCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(agentCommonCmd)
}

var agentCommonCmd = &cobra.Command{
	Use:    "agent",
	Short:  "run as agent",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		go StartAgent()
		go StartAgentRpcServer()
		NotityKill(StopAgent)
	},
}

func StartAgent() {
	agentIP = viper.GetString("agent.rpc.ip")
	agentPort = viper.GetString("agent.rpc.port")
	if agentIP == "" && agentPort == "" {
		panic("agent config error")
	}
	InitMasterClient()
	AgentRegister()
	AppsRegister()
	JobsRegister()
	applications.AppStartAll(true)
	applications.JobStartAll(true)
	applications.MoniterStart()
}

func StopAgent() {
	applications.AppStopAll()
	applications.JobStopAll()
}

func StartAgentRpcServer() {
	port := viper.GetString("agent.rpc.port")
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	s := server.DefaultServer()
	rpc.RegisterGuardServer(s, &server.GuardServer{})
	rpc.RegisterCronServer(s, &server.CronServer{})
	reflection.Register(s)
	logger.Info("agent rpc started at ", port)
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func InitMasterClient() {
	masterIP := viper.GetString("agent.master.ip")
	masterPort := viper.GetString("agent.master.port")
	addr := fmt.Sprintf("%s:%s", masterIP, masterPort)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic("Can't connect: " + addr)
	}
	masterClient = rpc.NewMasterClient(conn)
}

func AgentRegister() {
	timeout := time.Second * 30
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := masterClient.Register(ctx, &rpc.Agent{Ip: agentIP, Port: agentPort})
	if err != nil {
		panic("agent register fail: " + err.Error())
	}
	if response.GetCode() != 200 {
		panic("agent register fail: " + response.GetMessage())
	}
}

func AppsRegister() {
	timeout := time.Second * 30
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := masterClient.AppList(ctx, &rpc.Agent{Ip: agentIP, Port: agentPort})
	if err != nil {
		panic("get app list error: " + err.Error())
	}
	if response.GetCode() == 404 {
		panic("get app list error: agent error")
	}
	apps := response.GetApps()
	for _, app := range apps {
		config := map[string]interface{}{
			"id":           app.GetId(),
			"name":         app.GetName(),
			"dir":          app.GetDir(),
			"program":      app.GetProgram(),
			"args":         app.GetArgs(),
			"stdout":       app.GetStdOut(),
			"stderr":       app.GetStdErr(),
			"auto_restart": app.GetAutoRestart(),
			"is_monitor":   app.GetIsMonitor(),
		}
		err := applications.AppRegister(app.GetId(), config)
		if err != nil {
			logger.Error("app register failed:"+err.Error(), config)
			return
		}
	}
}

func JobsRegister() {
	timeout := time.Second * 30
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := masterClient.JobList(ctx, &rpc.Agent{Ip: agentIP, Port: agentPort})
	if err != nil {
		panic("get app list error: " + err.Error())
	}
	if response.GetCode() == 404 {
		panic("get app list error: agent error")
	}
	jobs := response.GetJobs()
	for _, job := range jobs {
		config := map[string]interface{}{
			"id":         job.GetId(),
			"name":       job.GetName(),
			"dir":        job.GetDir(),
			"program":    job.GetProgram(),
			"args":       job.GetArgs(),
			"stdout":     job.GetStdOut(),
			"stderr":     job.GetStdErr(),
			"spec":       job.GetSpec(),
			"timeout":    job.GetTimeout(),
			"is_monitor": job.GetIsMonitor(),
		}
		err := applications.JobRegister(job.GetId(), config)
		if err != nil {
			logger.Error("job register failed:"+err.Error(), config)
			return
		}
	}
}
