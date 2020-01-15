package cmd

import (
	"net"
	"os"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/client"
	"Asgard/rpc"
	"Asgard/server"
	"Asgard/services"
	"Asgard/web"
)

var (
	agentService *services.AgentService
	appService   *services.AppService
	jobService   *services.JobService
	ticker       *time.Ticker
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
		agentService = services.NewAgentService()
		appService = services.NewAppService()
		jobService = services.NewJobService()
		go StartWebServer()
		go StartMasterRpcServer()
		go MoniterMaster()
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
	ticker.Stop()
}

func MoniterMaster() {
	duration := viper.GetInt("system.moniter")
	if duration == 0 {
		duration = 10
	}
	ticker = time.NewTicker(time.Second * time.Duration(duration))
	for range ticker.C {
		CheckOnlineAgent()
		CheckOfflineAgent()
	}
}

func CheckOnlineAgent() {
	agentList := agentService.GetOnlineAgent()
	logger.Debug("online agent: ", len(agentList))
	for _, agent := range agentList {
		apps, err := client.GetGuardList(&agent)
		if err != nil {
			agent.Status = 0
			agentService.UpdateAgent(&agent)
		} else {
			for _, app := range apps {
				_app := appService.GetAppByID(app.GetId())
				if _app != nil {
					_app.Status = 1
					appService.UpdateApp(_app)
				}
			}
		}
		jobs, err := client.GetCronList(&agent)
		if err != nil {
			agent.Status = 0
			agentService.UpdateAgent(&agent)
		} else {
			for _, job := range jobs {
				_job := jobService.GetJobByID(job.GetId())
				if _job != nil {
					_job.Status = 1
					jobService.UpdateJob(_job)
				}
			}
		}
	}
}

func CheckOfflineAgent() {
	agentList := agentService.GetOfflineAgent()
	logger.Debug("offline agent: ", len(agentList))
	for _, agent := range agentList {
		_, err := client.GetCronClient(&agent)
		if err != nil {
			agent.Status = 1
			agentService.UpdateAgent(&agent)
		}
	}
}
