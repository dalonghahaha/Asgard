package guard

import (
	"fmt"
	"net"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/managers"
	"Asgard/rpc"
	"Asgard/runtimes"
	"Asgard/server"
)

var (
	appManager *managers.AppManager
	serverPath string
)

func GetCmd() *cobra.Command {
	guardCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	guardCmd.PersistentFlags().StringP("socket", "s", "runtime/asgard_guard", "socket file path")
	guardCmd.AddCommand(statusCmd)
	guardCmd.AddCommand(showCmd)
	return guardCmd
}

var guardCmd = &cobra.Command{
	Use:   "guard",
	Short: "guard apps",
	Run: func(cmd *cobra.Command, args []string) {
		confPath := cmd.Flag("conf").Value.String()
		serverPath = cmd.Flag("socket").Value.String()
		runtimes.ParseConfig(confPath)
		err := logger.Register()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := InitGuard(); err != nil {
			fmt.Println(err)
			return
		}
		go StartGuard()
		go StartRpcServer()
		runtimes.Wait(StopGuard)
	},
}

func InitGuard() error {
	configs := viper.Get("app")
	if configs == nil {
		return fmt.Errorf("no apps!")
	}
	_configs, ok := configs.([]interface{})
	if !ok {
		return fmt.Errorf("apps config wrong!")
	}
	appManager = managers.NewAppManager()
	for index, v := range _configs {
		_v, ok := v.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("apps config wrong!")
		}
		config := map[string]interface{}{}
		for k, v := range _v {
			_k, ok := k.(string)
			if !ok {
				return fmt.Errorf("apps config wrong!")
			}
			config[_k] = v
		}
		err := appManager.Register(int64(index), config)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartGuard() {
	defer func() {
		if err := recover(); err != nil {
			StopGuard()
			fmt.Println("guard server panic:", err)
			return
		}
	}()
	logger.Info("Guard Server start")
	logger.Infof("apps:%d", appManager.Count())
	appManager.StartAll(false)
}

func StopGuard() {
	runtimes.Exit()
	appManager.StopAll()
	logger.Info("Guard Server stop")
	err := file.Remove(serverPath)
	if err != nil {
		logger.Errorf("remove guard server path failed:%s", err.Error())
	}
}

func StartRpcServer() {
	defer func() {
		if err := recover(); err != nil {
			StopGuard()
			fmt.Println("rpc server panic:", err)
			return
		}
	}()
	serverAddr, err := net.ResolveUnixAddr("unix", serverPath)
	if err != nil {
		logger.Error("fialed to resolve unix addr")
		panic(err)
	}
	listen, err := net.ListenUnix("unix", serverAddr)
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	s := server.NewRPCServer()
	rpc.RegisterGuardServer(s, server.NewGuardServer(appManager))
	reflection.Register(s)
	logger.Info("rpc server start!")
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}
