package cmds

import (
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/clients"
	"Asgard/constants"
	"Asgard/managers"
	"Asgard/rpc"
	"Asgard/runtimes"
	"Asgard/server"
)

func init() {
	agentCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	statusCommonCmd.PersistentFlags().StringP("port", "p", "27149", "agent port")
	agentCommonCmd.AddCommand(statusCommonCmd)
	rootCmd.AddCommand(agentCommonCmd)
}

var agentCommonCmd = &cobra.Command{
	Use:    "agent",
	Short:  "run as agent",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				NotityKill(StopAgent)
				fmt.Println("panic:", err)
				fmt.Println("stack:", string(debug.Stack()))
				return
			}
		}()
		InitAgent()
		go masterClient.Report()
		go StartAgent()
		go StartAgentRpcServer()
		go MoniterAgent()
		NotityKill(StopAgent)
	},
}

var statusCommonCmd = &cobra.Command{
	Use:   "status",
	Short: "show agent running status",
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("port").Value.String()
		client, err := clients.NewAgent("127.0.0.1", port)
		if err != nil {
			fmt.Printf("fail connect to agent:%s\n", err.Error())
			return
		}
		apps, err := client.GetAppList()
		if err != nil {
			fmt.Printf("fail connect to get app list:%s\n", err.Error())
			return
		}
		jobs, err := client.GetJobList()
		if err != nil {
			fmt.Printf("fail connect to get app list:%s\n", err.Error())
			return
		}
		timings, err := client.GetTimingList()
		if err != nil {
			fmt.Printf("fail connect to get app list:%s\n", err.Error())
			return
		}
		titleFormat := "%-5s %-50s %-50s %-30s\n"
		contentFormat := "%-5d %-50s %-50s %-30s\n"
		fmt.Println()
		fmt.Println("app total:", len(apps))
		fmt.Println()
		fmt.Printf(titleFormat, "ID", "Dir", "Program", "Name")
		fmt.Println()
		for _, app := range apps {
			program := fmt.Sprintf("%s %s", app.GetProgram(), app.GetArgs())
			fmt.Printf(contentFormat, app.GetId(), app.GetDir(), program, app.GetName())
		}
		fmt.Println()
		fmt.Println("job list:", len(jobs))
		fmt.Println()
		fmt.Printf(titleFormat, "ID", "Dir", "Program", "Name")
		fmt.Println()
		for _, job := range jobs {
			program := fmt.Sprintf("%s %s", job.GetProgram(), job.GetArgs())
			fmt.Printf(contentFormat, job.GetId(), job.GetDir(), program, job.GetName())
		}
		fmt.Println()
		fmt.Println("timing list:", len(timings))
		fmt.Println()
		fmt.Printf(titleFormat, "ID", "Dir", "Program", "Name")
		fmt.Println()
		for _, timing := range timings {
			program := fmt.Sprintf("%s %s", timing.GetProgram(), timing.GetArgs())
			fmt.Printf(contentFormat, timing.GetId(), timing.GetDir(), program, timing.GetName())
		}
	},
}

func InitAgent() {
	agentIP := viper.GetString("agent.rpc.ip")
	agentPort := viper.GetString("agent.rpc.port")
	if agentIP == "" && agentPort == "" {
		panic("agent config error")
	}
	constants.AGENT_IP = agentIP
	constants.AGENT_PORT = agentPort
	masterIP := viper.GetString("agent.master.ip")
	masterPort := viper.GetString("agent.master.port")
	if masterIP == "" && masterPort == "" {
		panic("agent config error")
	}
	constants.MASTER_IP = masterIP
	constants.MASTER_PORT = masterPort
	constants.AGENT_PID = os.Getpid()
	constants.AGENT_UUID = uuid.GenerateV4()
	duration := viper.GetInt("agent.moniter")
	if duration != 0 {
		constants.AGENT_MONITER = duration
	}
	constants.AGENT_MONITER_TICKER = time.NewTicker(time.Second * time.Duration(constants.AGENT_MONITER))

	_masterClient, err := clients.NewMaster(constants.MASTER_IP, constants.MASTER_PORT)
	if err != nil {
		panic("init master client failed:" + err.Error())
	}
	masterClient = _masterClient
	if err != nil {
		panic("register master failed:" + err.Error())
	}
	appManager = managers.NewAppManager()
	appManager.SetMaster(masterClient)
	jobManager = managers.NewJobManager()
	jobManager.SetMaster(masterClient)
	timingManager = managers.NewTimingManager()
	timingManager.SetMaster(masterClient)
}

func StartAgent() {
	err := masterClient.AgentRegister()
	if err != nil {
		panic(err)
	}
	err = AppsRegister()
	if err != nil {
		panic(err)
	}
	err = JobsRegister()
	if err != nil {
		panic(err)
	}
	err = TimingsRegister()
	if err != nil {
		panic(err)
	}
	appManager.StartAll()
	appManager.StartMonitor()
	jobManager.StartAll()
	jobManager.StartMonitor()
	timingManager.StartAll()
	timingManager.StartMonitor()
	logger.Info("Agent Started!")
	logger.Debugf("Agent Master: %s:%s", constants.MASTER_IP, constants.MASTER_PORT)
	logger.Debugf("Agent Address: %s:%s", constants.AGENT_IP, constants.AGENT_PORT)
	logger.Debugf("Agent Loop:%d", constants.AGENT_MONITER)
}

func StopAgent() {
	runtimes.Exit()
	constants.AGENT_MONITER_TICKER.Stop()
	appManager.StopAll()
	appManager.StopMonitor()
	jobManager.StopAll()
	jobManager.StopMonitor()
	timingManager.StopAll()
	timingManager.StopMonitor()
	//make sure all data report to master
	maxWait := 10
	countWait := 0
	for {
		if masterClient.IsRunning() && countWait <= maxWait {
			time.Sleep(time.Second * 1)
			countWait += 1
		} else {
			break
		}
	}
	logger.Info("Agent Server Stop!")
}

func StartAgentRpcServer() {
	listen, err := net.Listen("tcp", ":"+constants.AGENT_PORT)
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	s := server.NewRPCServer()
	agentServer := &server.AgentServer{}
	agentServer.SetAppManager(appManager)
	agentServer.SetJobManager(jobManager)
	agentServer.SetTimingManager(timingManager)
	rpc.RegisterAgentServer(s, agentServer)
	reflection.Register(s)
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func MoniterAgent() {
	for range constants.AGENT_MONITER_TICKER.C {
		go AgentMonitorReport()
	}
}

func AgentMonitorReport() {
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
	masterClient.AgentMonitorChan <- agentMonitor
}

func AppsRegister() error {
	apps, err := masterClient.GetAppList()
	if err != nil {
		return err
	}
	for _, app := range apps {
		err := appManager.Register(app.GetId(), rpc.BuildAppConfig(app))
		if err != nil {
			return err
		}
	}
	return nil
}

func JobsRegister() error {
	jobs, err := masterClient.GetJobList()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		err := jobManager.Register(job.GetId(), rpc.BuildJobConfig(job))
		if err != nil {
			return err
		}
	}
	return nil
}

func TimingsRegister() error {
	timings, err := masterClient.GetTimingList()
	if err != nil {
		return err
	}
	for _, timing := range timings {
		err := timingManager.Register(timing.GetId(), rpc.BuildTimingConfig(timing))
		if err != nil {
			return err
		}
	}
	return nil
}
