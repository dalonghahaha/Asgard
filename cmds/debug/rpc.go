package debug

import (
	"context"
	"fmt"
	"net"
	"syscall"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/resolver"

	"Asgard/managers"
	"Asgard/registry"
	"Asgard/rpc"
	"Asgard/runtimes"
	"Asgard/server"
)

var (
	serverID     string
	agentManager *managers.AgentManager
	rpcServer    *grpc.Server
)

var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "debug agent server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("this cmd use for debug rpc server")
	},
}

var serverCmd = &cobra.Command{
	Use:    "server",
	Short:  "debug rpc server",
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		agentManager, err = managers.NewAgentManager(nil)
		if err != nil {
			fmt.Printf("init agentManager failed:%+v", err)
			return
		}
		err = registry.RegisterRegistry([]string{"http://localhost:2379"})
		if err != nil {
			fmt.Printf("init serviceCenter failed:%+v", err)
			return
		}
		id := cmd.Flag("id").Value.String()
		if id == "" {
			fmt.Printf("id can not be empty!")
			return
		} else {
			serverID = id
		}
		port := cmd.Flag("port").Value.String()
		if port == "" {
			fmt.Printf("port can not be empty!")
			return
		}
		go RegisterServer(port)
		go StartServer(port)
		runtimes.Wait(StopServer)
	},
}

var clientCmd = &cobra.Command{
	Use:    "client",
	Short:  "debug rpc client",
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := registry.NewResolver([]string{"http://localhost:2379"})
		if err != nil {
			fmt.Printf("registry failed:%+v", err)
			return
		}
		resolver.Register(r)
		conn, err := grpc.Dial(
			r.Scheme()+"://author/inter",
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, roundrobin.Name)),
			grpc.WithInsecure(),
		)
		if err != nil {
			fmt.Printf("grpc dial failed:%+v", err)
			return
		}
		client := rpc.NewAgentClient(conn)
		for {
			result, err := client.Stat(context.Background(), &rpc.Empty{})
			if err != nil {
				fmt.Printf("grpc call failed:%+v", err)
				return
			}
			fmt.Println(result.GetCode())
			time.Sleep(time.Second)
		}
	},
}

func preRun(cmd *cobra.Command, args []string) {
	confPath := cmd.Flag("conf").Value.String()
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = logger.Register()
	if err != nil {
		panic(err)
	}
}

func RegisterServer(port string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("register rpc server panic:", err)
			runtimes.ExitSinal <- syscall.SIGTERM
			return
		}
	}()
	logger.Infof("register rpc server:[%s][%s]", serverID, port)
	registry.Register("inter", serverID, "127.0.0.1", port)
}

func StartServer(port string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("rpc server start failed:", err)
			runtimes.ExitSinal <- syscall.SIGTERM
			return
		}
	}()
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	rpcServer = server.NewRPCServer()
	rpc.RegisterAgentServer(rpcServer, server.NewAgentServer(agentManager))
	reflection.Register(rpcServer)
	logger.Info("rpc server start at : ", port)
	err = rpcServer.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func StopServer() {
	if rpcServer != nil {
		rpcServer.GracefulStop()
	}
	registry.UnRegister("inter", serverID)
	logger.Info("rpc server stop!")
}
