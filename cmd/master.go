package cmd

import (
	"net"
	"os"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/rpc"
	"Asgard/server"
	"Asgard/web"
)

func init() {
	masterCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(masterCmd)
}

var masterCmd = &cobra.Command{
	Use:    "master",
	Short:  "run as master",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		go StartWebServer()
		go StartMasterRpcServer()
		NotityKill(StopMaster)
	},
}

func StartWebServer() {
	err := web.Init()
	if err != nil {
		logger.Error("web init error:", err)
		os.Exit(1)
	}
	web.Run()
}

func StartMasterRpcServer() {
	port := viper.GetString("master.rpc.port")
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	s := server.DefaultServer()
	rpc.RegisterMasterServer(s, &server.MasterServer{})
	reflection.Register(s)
	logger.Info("master rpc started at ", port)
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func StopMaster() {

}
