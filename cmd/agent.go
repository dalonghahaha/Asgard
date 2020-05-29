package cmd

import (
	"context"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"Asgard/applications"
	"Asgard/constants"
	"Asgard/providers"
	"Asgard/rpc"
	"Asgard/server"
)

func init() {
	agentCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	listCommonCmd.PersistentFlags().StringP("port", "p", "27147", "agent port")
	agentCommonCmd.AddCommand(listCommonCmd)
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
		go StartAgent()
		go StartAgentRpcServer()
		go MoniterAgent()
		NotityKill(StopAgent)
	},
}

var listCommonCmd = &cobra.Command{
	Use:   "list",
	Short: "show agent running applications",
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("port").Value.String()
		addr := fmt.Sprintf("127.0.0.1:%s", port)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		option := grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*1024),
			grpc.MaxCallSendMsgSize(1024*1024*1024),
		)
		conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), option)
		if err != nil {
			fmt.Printf("can't connect to agent!:%s\n", err.Error())
		}
		client := rpc.NewAgentClient(conn)
		fmt.Println(client)
		fmt.Println()
		fmt.Println("app list")
		fmt.Println()
		fmt.Println("job list")
		fmt.Println()
		fmt.Println("timing list")
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
	err := providers.RegisterMaster()
	if err != nil {
		panic("register master failed:" + err.Error())
	}
}

func StartAgent() {
	err := providers.MasterClient.AgentRegister()
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
	applications.AppStartAll(false)
	applications.JobStartAll(false)
	applications.TimingStartAll(false)
	logger.Info("Agent Server Started!")
	applications.MoniterStart()
}

func StopAgent() {
	constants.AGENT_MONITER_TICKER.Stop()
	applications.AppStopAll()
	applications.JobStopAll()
	applications.TimingStopAll()
	maxWait := 10
	countWait := 0
	for {
		if providers.MasterClient.IsRunning() && countWait <= maxWait {
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
	rpc.RegisterAgentServer(s, &server.AgentServer{})
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
	agentMonitor := applications.AgentMonitor{
		Ip:      constants.AGENT_IP,
		Port:    constants.AGENT_PORT,
		Monitor: applications.BuildMonitor(info),
	}
	providers.MasterClient.AgentMonitorChan <- agentMonitor
}

func AppsRegister() error {
	apps, err := providers.MasterClient.GetAppList()
	if err != nil {
		return err
	}
	for _, app := range apps {
		logger.Debug("app register: ", app.GetName())
		config := rpc.BuildAppConfig(app)
		err := applications.AppRegister(
			app.GetId(),
			config,
			providers.MasterClient.Reports,
			providers.MasterClient.AppMonitorChan,
			providers.MasterClient.AppArchiveChan,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func JobsRegister() error {
	jobs, err := providers.MasterClient.GetJobList()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		logger.Debugf("job register: %s %s", job.GetName(), job.GetSpec())
		config := rpc.BuildJobConfig(job)
		err := applications.JobRegister(
			job.GetId(),
			config,
			providers.MasterClient.JobMonitorChan,
			providers.MasterClient.JobArchiveChan,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func TimingsRegister() error {
	timings, err := providers.MasterClient.GetTimingList()
	if err != nil {
		return err
	}
	for _, timing := range timings {
		logger.Debugf("timing register: %s %s", timing.GetName(), time.Unix(timing.GetTime(), 0).Format("2006-01-02 15:04:05"))
		config := rpc.BuildTimingConfig(timing)
		err := applications.TimingRegister(
			timing.GetId(),
			config,
			providers.MasterClient.TimingMonitorChan,
			providers.MasterClient.TimingArchiveChan,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
