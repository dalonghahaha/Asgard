package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/applications"
	"Asgard/rpc"
	"Asgard/server"
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
		port := viper.GetString("agent.port")
		listen, err := net.Listen("tcp", ":"+port)
		if err != nil {
			logger.Error("failed to listen:", err)
			return
		}
		s := server.DefaultServer()
		rpc.RegisterGuardServer(s, &server.GuardServer{})
		reflection.Register(s)
		go test()
		logger.Info("agent started at ", port)
		err = s.Serve(listen)
		if err != nil {
			logger.Error("failed to serve:", err)
			return
		}
	},
}

func test() {
	logger.Info("guard started at ", os.Getpid())
	configs := viper.Get("app")
	if configs == nil {
		fmt.Println("no apps!")
		return
	}
	_configs, ok := configs.([]interface{})
	if !ok {
		fmt.Println("apps config wrong!")
		return
	}
	for _, v := range _configs {
		_v, ok := v.(map[interface{}]interface{})
		if !ok {
			fmt.Println("apps config wrong!")
			return
		}
		config := map[string]interface{}{}
		for k, v := range _v {
			_k, ok := k.(string)
			if !ok {
				fmt.Println("apps config wrong!")
				return
			}
			config[_k] = v
		}
		err := applications.AppRegister(config)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	applications.AppStartAll()
	NotityKill(applications.AppStopAll)
}
