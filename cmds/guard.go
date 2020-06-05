package cmds

import (
	"fmt"
	"net"
	"os"
	"runtime/debug"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/clients"
	"Asgard/managers"
	"Asgard/rpc"
	"Asgard/runtimes"
	"Asgard/server"
)

func init() {
	guardCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	guardCommonCmd.PersistentFlags().StringP("sock", "s", "runtime/asgard_guard", "socket file path")
	guardCommonCmd.AddCommand(GuardStatusCommonCmd)
	rootCmd.AddCommand(guardCommonCmd)
}

var guardServerPath string

var GuardStatusCommonCmd = &cobra.Command{
	Use:   "status",
	Short: "show guard running status",
	Run: func(cmd *cobra.Command, args []string) {
		serverFile := cmd.Flag("sock").Value.String()
		client, err := clients.NewGuard(serverFile)
		if err != nil {
			fmt.Printf("fail connect to guard:%s\n", err.Error())
			return
		}
		apps, err := client.GetList()
		if err != nil {
			fmt.Printf("fail to get app list:%s\n", err.Error())
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
	},
}

var guardCommonCmd = &cobra.Command{
	Use:    "guard",
	Short:  "guard apps",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		guardServerPath = cmd.Flag("sock").Value.String()
		go StartGuard()
		go StartGuardRpcServer()
		Wait(StopGuard)
	},
}

func StartGuard() {
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
	appManager = managers.NewAppManager()
	for index, v := range _configs {
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
		err := appManager.Register(int64(index), config)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	logger.Info("guard started at ", os.Getpid())
	appManager.StartAll(false)
}

func StartGuardRpcServer() {
	defer func() {
		if err := recover(); err != nil {
			appManager.StopAll()
			fmt.Println("panic:", err)
			fmt.Println("stack:", string(debug.Stack()))
			return
		}
	}()
	serverAddr, err := net.ResolveUnixAddr("unix", guardServerPath)
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
	logger.Info("local rpc server started!")
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func StopGuard() {
	runtimes.Exit()
	appManager.StopAll()
	err := file.Remove(guardServerPath)
	if err != nil {
		logger.Errorf("remove guard server path failed:%s", err.Error())
	}
}
