package cmd

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

	"Asgard/applications"
	"Asgard/client"
	"Asgard/constants"
	"Asgard/rpc"
	"Asgard/server"
)

func init() {
	agentCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(agentCommonCmd)
}

var master *client.Master

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

func InitAgent() {
	agentIP := viper.GetString("agent.rpc.ip")
	agentPort := viper.GetString("agent.rpc.port")
	if agentIP == "" && agentPort == "" {
		panic("agent config error")
	}
	constants.AGENT_IP = agentIP
	constants.AGENT_PORT = agentPort
	masterIP := viper.GetString("agent.master.ip")
	masterPort := viper.GetInt("agent.master.port")
	if masterIP == "" && masterPort == 0 {
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
	master = client.NewMaster()
}

func StartAgent() {
	err := master.AgentRegister()
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
	applications.MoniterStart()
}

func StopAgent() {
	constants.AGENT_MONITER_TICKER.Stop()
	applications.AppStopAll()
	applications.JobStopAll()
	applications.TimingStopAll()
}

func StartAgentRpcServer() {
	port := viper.GetString("agent.rpc.port")
	listen, err := net.Listen("tcp", ":"+port)
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
	monitor := applications.BuildMonitor(info)
	master.AgentMonitorReport(rpc.BuildAgentMonitor(monitor))
}

func AppsRegister() error {
	apps, err := master.GetAppList()
	if err != nil {
		return err
	}
	for _, app := range apps {
		logger.Debug("app register: ", app.GetName())
		config := rpc.BuildAppConfig(app)
		app, err := applications.AppRegister(app.GetId(), config)
		if err != nil {
			return err
		}
		app.MonitorReport = func(monitor *applications.Monitor) {
			master.AppMonitorReport(rpc.BuildAppMonitor(app, monitor))
		}
		app.ArchiveReport = func(command *applications.Command) {
			master.AppArchiveReport(rpc.BuildAppArchive(app, command))
		}
	}
	return nil
}

func JobsRegister() error {
	jobs, err := master.GetJobList()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		logger.Debug("job register: ", job.GetName())
		config := rpc.BuildJobConfig(job)
		job, err := applications.JobRegister(job.GetId(), config)
		if err != nil {
			return err
		}
		job.MonitorReport = func(monitor *applications.Monitor) {
			master.JobMonitorReport(rpc.BuildJobMonior(job, monitor))
		}
		job.ArchiveReport = func(command *applications.Command) {
			master.JobArchiveReport(rpc.BuildJobArchive(job, command))
		}
	}
	return nil
}

func TimingsRegister() error {
	timings, err := master.GetTimingList()
	if err != nil {
		return err
	}
	for _, timing := range timings {
		logger.Debug(fmt.Sprintf("timing register: %s %v", timing.GetName(), time.Unix(timing.GetTime(), 0).Format("2006-01-02 15:04:05")))
		config := rpc.BuildTimingConfig(timing)
		timing, err := applications.TimingRegister(timing.GetId(), config)
		if err != nil {
			return err
		}
		timing.MonitorReport = func(monitor *applications.Monitor) {
			master.TimingMonitorReport(rpc.BuildTimingMonior(timing, monitor))
		}
		timing.ArchiveReport = func(command *applications.Command) {
			master.TimingArchiveReport(rpc.BuildTimingArchive(timing, command))
		}
	}
	return nil
}
