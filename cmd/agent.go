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

var masterClient rpc.MasterClient

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
		go RegisterAgent()
		NotityKill(StopAgent)
	},
}

func StartAgent() {
	applications.AppStartAll(false)
	applications.JobStartAll(false)
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

func RegisterAgent() {
	masterIP := viper.GetString("agent.master.ip")
	masterPort := viper.GetString("agent.master.port")
	agentIP := viper.GetString("agent.rpc.ip")
	agentPort := viper.GetString("agent.rpc.port")
	addr := fmt.Sprintf("%s:%s", masterIP, masterPort)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic("Can't connect: " + addr)
	}
	masterClient = rpc.NewMasterClient(conn)
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
