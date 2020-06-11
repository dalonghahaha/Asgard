package cmds

import (
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"syscall"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
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
	RootCmd.AddCommand(agentCommonCmd)
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
			fmt.Printf("fail to get app list:%s\n", err.Error())
			return
		}
		jobs, err := client.GetJobList()
		if err != nil {
			fmt.Printf("fail to get job list:%s\n", err.Error())
			return
		}
		timings, err := client.GetTimingList()
		if err != nil {
			fmt.Printf("fail to get timing list:%s\n", err.Error())
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

var agentCommonCmd = &cobra.Command{
	Use:    "agent",
	Short:  "run as agent",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("panic:", err)
				fmt.Println("stack:", string(debug.Stack()))
				return
			}
		}()
		InitAgent()
		go agentManager.StartAll()
		go StartAgentRpcServer()
		runtimes.Wait(agentManager.StopAll)
	},
}

func InitAgent() {
	agentIP := viper.GetString("agent.rpc.ip")
	agentPort := viper.GetString("agent.rpc.port")
	if agentIP == "" && agentPort == "" {
		panic("agent rpc config error")
	}
	constants.AGENT_IP = agentIP
	constants.AGENT_PORT = agentPort
	masterIP := viper.GetString("agent.master.ip")
	masterPort := viper.GetString("agent.master.port")
	if masterIP == "" && masterPort == "" {
		panic("agent master config error")
	}
	constants.MASTER_IP = masterIP
	constants.MASTER_PORT = masterPort
	constants.AGENT_PID = os.Getpid()
	constants.AGENT_UUID = uuid.GenerateV4()
	duration := viper.GetInt("agent.moniter")
	if duration != 0 {
		constants.AGENT_MONITER = duration
	}
	var err error
	agentManager, err = managers.NewAgentManager()
	if err != nil {
		panic("init agentManager failed:" + err.Error())
	}
}

func StartAgentRpcServer() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("agent rpc server start failed:", err)
			runtimes.ExitSinal <- syscall.SIGTERM
			return
		}
	}()
	listen, err := net.Listen("tcp", ":"+constants.AGENT_PORT)
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	s := server.NewRPCServer()
	rpc.RegisterAgentServer(s, server.NewAgentServer(agentManager))
	reflection.Register(s)
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}
