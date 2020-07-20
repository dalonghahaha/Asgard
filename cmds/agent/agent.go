package agent

import (
	"fmt"
	"net"
	"os"
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

var agentManager *managers.AgentManager

func GetCmd() *cobra.Command {
	agentCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	statusCmd.PersistentFlags().StringP("port", "p", "27149", "agent port")
	agentCmd.AddCommand(statusCmd)
	return agentCmd
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "run as agent",
	Run: func(cmd *cobra.Command, args []string) {
		confPath := cmd.Flag("conf").Value.String()
		runtimes.ParseConfig(confPath)
		if err := InitAgent(); err != nil {
			fmt.Println(err)
			return
		}
		go agentManager.StartAll()
		go StartRpcServer()
		runtimes.Wait(agentManager.StopAll)
	},
}

func InitAgent() error {
	if err := logger.Register(); err != nil {
		return fmt.Errorf("init logger failed:%+v", err)
	}
	agentIP := viper.GetString("agent.rpc.ip")
	agentPort := viper.GetString("agent.rpc.port")
	if agentIP == "" && agentPort == "" {
		return fmt.Errorf("agent rpc config error")
	}
	constants.AGENT_IP = agentIP
	constants.AGENT_PORT = agentPort
	constants.AGENT_PID = os.Getpid()
	constants.AGENT_UUID = uuid.GenerateV4()
	duration := viper.GetInt("agent.moniter")
	if duration != 0 {
		constants.AGENT_MONITER = duration
	}
	var master *clients.Master
	var err error
	cluster := viper.GetBool("agent.master.cluster")
	if cluster {
		constants.MASTER_CLUSTER_REGISTRY = viper.GetStringSlice("agent.master.cluster_registry")
		constants.MASTER_CLUSTER_NAME = viper.GetString("agent.master.cluster_name")
		master, err = clients.NewClusterMaster(constants.MASTER_CLUSTER_REGISTRY)
		if err != nil {
			return fmt.Errorf("init cluster master client failed:%s", err.Error())
		}
	} else {
		masterIP := viper.GetString("agent.master.ip")
		masterPort := viper.GetString("agent.master.port")
		if masterIP == "" && masterPort == "" {
			return fmt.Errorf("agent master config error")
		}
		constants.MASTER_IP = masterIP
		constants.MASTER_PORT = masterPort
		master, err = clients.NewMaster(constants.MASTER_IP, constants.MASTER_PORT)
		if err != nil {
			return fmt.Errorf("init master client failed:%s", err.Error())
		}
	}
	agentManager, err = managers.NewAgentManager(master)
	if err != nil {
		return fmt.Errorf("init agentManager failed:" + err.Error())
	}
	return nil
}

func StartRpcServer() {
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
