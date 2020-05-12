package cmd

import (
	"net"
	"os"
	"time"

	"github.com/dalonghahaha/avenger/components/cache"
	"github.com/dalonghahaha/avenger/components/db"
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/client"
	"Asgard/constants"
	"Asgard/rpc"
	"Asgard/server"
	"Asgard/services"
	"Asgard/web"
)

var (
	agentService        *services.AgentService
	appService          *services.AppService
	jobService          *services.JobService
	timingService       *services.TimingService
	masterMoniterTicker *time.Ticker
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
		err := db.Register()
		if err != nil {
			panic(err)
		}
		err = cache.Register()
		if err != nil {
			panic(err)
		}
		agentService = services.NewAgentService()
		appService = services.NewAppService()
		jobService = services.NewJobService()
		timingService = services.NewTimingService()
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
	s := server.NewRPCServer()
	rpc.RegisterMasterServer(s, &server.MasterServer{})
	reflection.Register(s)
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func StopMaster() {
	masterMoniterTicker.Stop()
}

func MoniterMaster() {
	duration := viper.GetInt("system.moniter")
	if duration == 0 {
		duration = 10
	}
	masterMoniterTicker = time.NewTicker(time.Second * time.Duration(duration))
	for range masterMoniterTicker.C {
		CheckOnlineAgent()
		CheckOfflineAgent()
	}
}

func CheckOnlineAgent() {
	agentList := agentService.GetOnlineAgent()
	for _, agent := range agentList {
		apps, err := client.GetAgentAppList(&agent)
		if err != nil {
			agent.Status = constants.AGENT_OFFLINE
			agentService.UpdateAgent(&agent)
		} else {
			for _, app := range apps {
				_app := appService.GetAppByID(app.GetId())
				if _app != nil {
					_app.Status = constants.APP_STATUS_RUNNING
					appService.UpdateApp(_app)
				}
			}
		}
		jobs, err := client.GetAgentJobList(&agent)
		if err != nil {
			agent.Status = constants.AGENT_OFFLINE
			agentService.UpdateAgent(&agent)
		} else {
			for _, job := range jobs {
				_job := jobService.GetJobByID(job.GetId())
				if _job != nil {
					_job.Status = constants.JOB_STATUS_RUNNING
					jobService.UpdateJob(_job)
				}
			}
		}
		timings, err := client.GetAgentTimingList(&agent)
		if err != nil {
			agent.Status = constants.AGENT_OFFLINE
			agentService.UpdateAgent(&agent)
		} else {
			for _, timing := range timings {
				_timing := timingService.GetTimingByID(timing.GetId())
				if _timing != nil {
					_timing.Status = constants.TIMING_STATUS_RUNNING
					timingService.UpdateTiming(_timing)
				}
			}
		}
	}
}

func CheckOfflineAgent() {
	agentList := agentService.GetOfflineAgent()
	for _, agent := range agentList {
		_, err := client.GetAgentStat(&agent)
		if err == nil {
			agent.Status = constants.AGENT_ONLINE
			agentService.UpdateAgent(&agent)
		} else {
			apps := appService.GetAppByAgentID(agent.ID)
			for _, app := range apps {
				if app.Status != constants.APP_STATUS_PAUSE && app.Status != constants.APP_STATUS_DELETED {
					app.Status = constants.APP_STATUS_STOP
					appService.UpdateApp(&app)
				}
			}
			jobs := jobService.GetJobByAgentID(agent.ID)
			for _, job := range jobs {
				if job.Status != constants.JOB_STATUS_PAUSE && job.Status != constants.JOB_STATUS_DELETED {
					job.Status = constants.JOB_STATUS_STOP
					jobService.UpdateJob(&job)
				}
			}
			timings := timingService.GetTimingByAgentID(agent.ID)
			for _, timing := range timings {
				if timing.Status != constants.TIMING_STATUS_PAUSE && timing.Status != constants.TIMING_STATUS_DELETED && timing.Status != constants.TIMING_STATUS_FINISHED {
					timing.Status = constants.TIMING_STATUS_STOP
					timingService.UpdateTiming(&timing)
				}
			}
		}
	}
}
