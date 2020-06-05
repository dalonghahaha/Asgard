package cron

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
	jobManager *managers.JobManager
	serverPath string
)

func GetCmd() *cobra.Command {
	cronCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	cronCommonCmd.PersistentFlags().StringP("socket", "s", "runtime/asgard_cron", "socket file path")
	cronCommonCmd.AddCommand(statusCommonCmd)
	cronCommonCmd.AddCommand(showCommonCmd)
	return cronCommonCmd
}

var cronCommonCmd = &cobra.Command{
	Use:   "cron",
	Short: "cron jobs",
	Run: func(cmd *cobra.Command, args []string) {
		confPath := cmd.Flag("conf").Value.String()
		serverPath = cmd.Flag("socket").Value.String()
		runtimes.ParseConfig(confPath)
		err := logger.Register()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := InitCron(); err != nil {
			fmt.Println(err)
			return
		}
		go StartCron()
		go StartRpcServer()
		runtimes.Wait(StopCron)
	},
}

func InitCron() error {
	configs := viper.Get("cron")
	if configs == nil {
		return fmt.Errorf("no crons!")
	}
	_configs, ok := configs.([]interface{})
	if !ok {
		return fmt.Errorf("crons config wrong!")
	}
	jobManager = managers.NewJobManager()
	for index, v := range _configs {
		_v, ok := v.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("crons config wrong!")
		}
		config := map[string]interface{}{}
		for k, v := range _v {
			_k, ok := k.(string)
			if !ok {
				return fmt.Errorf("crons config wrong!")
			}
			config[_k] = v
		}
		err := jobManager.Register(int64(index), config)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartCron() {
	defer func() {
		if err := recover(); err != nil {
			StopCron()
			fmt.Println("cron server panic:", err)
			return
		}
	}()
	logger.Info("Cron Server start")
	logger.Infof("Jobs:%d", jobManager.Count())
	jobManager.StartAll(false)
}

func StopCron() {
	runtimes.Exit()
	jobManager.StopAll()
	logger.Info("cron Server stop")
	err := file.Remove(serverPath)
	if err != nil {
		logger.Errorf("remove guard server path failed:%s", err.Error())
	}
}

func StartRpcServer() {
	defer func() {
		if err := recover(); err != nil {
			StopCron()
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
	rpc.RegisterCronServer(s, server.NewCronServer(jobManager))
	reflection.Register(s)
	logger.Info("rpc server start!")
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}
