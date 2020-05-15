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
	"Asgard/models"
	"Asgard/providers"
	"Asgard/rpc"
	"Asgard/server"
)

func init() {
	masterCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(masterCmd)
}

var masterMoniterTicker *time.Ticker

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
		go StartMasterRpcServer()
		go MoniterMaster()
		NotityKill(StopMaster)
	},
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
	logger.Info("Master Rpc Server Started! Pid:", os.Getpid())
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func StopMaster() {
	logger.Info("Master Rpc Server Stop!")
	masterMoniterTicker.Stop()
}

func MoniterMaster() {
	duration := viper.GetInt("system.moniter")
	if duration == 0 {
		duration = 10
	}
	masterMoniterTicker = time.NewTicker(time.Second * time.Duration(duration))
	for range masterMoniterTicker.C {
		agentList := providers.AgentService.GetUsageAgent()
		for _, agent := range agentList {
			go checkAgent(agent)
		}
	}
}

func checkAgent(agent models.Agent) {
	usageApps := providers.AppService.GetUsageAppByAgentID(agent.ID)
	usageJobs := providers.JobService.GetUsageJobByAgentID(agent.ID)
	usageTimings := providers.TimingService.GetUsageTimingByAgentID(agent.ID)
	_, err := client.GetAgentStat(&agent)
	if err != nil {
		//标记实例状态为离线
		agent.Status = constants.AGENT_OFFLINE
		providers.AgentService.UpdateAgent(&agent)
		//标记应用状态为未知
		for _, app := range usageApps {
			app.Status = constants.APP_STATUS_UNKNOWN
			providers.AppService.UpdateApp(&app)
		}
		//标记计划任务状态为未知
		for _, job := range usageJobs {
			job.Status = constants.APP_STATUS_UNKNOWN
			providers.JobService.UpdateJob(&job)
		}
		//标记定时任务状态为未知
		for _, timing := range usageTimings {
			timing.Status = constants.APP_STATUS_UNKNOWN
			providers.TimingService.UpdateTiming(&timing)
		}
		return
	} else {
		//标记实例状态为在线
		agent.Status = constants.AGENT_ONLINE
		providers.AgentService.UpdateAgent(&agent)
		//更新实例应用运行状态
		apps, err := client.GetAgentAppList(&agent)
		if err != nil {
			runningApps := map[int64]string{}
			for _, app := range apps {
				runningApps[app.GetId()] = app.GetName()
			}
			for _, app := range usageApps {
				_, ok := runningApps[app.ID]
				if ok {
					app.Status = constants.APP_STATUS_RUNNING
				} else {
					app.Status = constants.APP_STATUS_STOP
				}
				providers.AppService.UpdateApp(&app)
			}
		} else {
			logger.Error("checkOnlineAgent GetAgentAppList Error:", err)
		}
		//更新实例计划任务运行状态
		jobs, err := client.GetAgentJobList(&agent)
		if err != nil {
			runningJobs := map[int64]string{}
			for _, job := range jobs {
				runningJobs[job.GetId()] = job.GetName()
			}
			for _, job := range usageJobs {
				_, ok := runningJobs[job.ID]
				if ok {
					job.Status = constants.JOB_STATUS_RUNNING
				} else {
					job.Status = constants.JOB_STATUS_STOP
				}
				providers.JobService.UpdateJob(&job)
			}
		} else {
			logger.Error("checkOnlineAgent GetAgentJobList Error:", err)
		}
		//更新实例计划任务运行状态
		timings, err := client.GetAgentTimingList(&agent)
		if err != nil {
			runningTimings := map[int64]string{}
			for _, timing := range timings {
				runningTimings[timing.GetId()] = timing.GetName()
			}
			for _, timing := range usageTimings {
				_, ok := runningTimings[timing.ID]
				if ok {
					timing.Status = constants.JOB_STATUS_RUNNING
				} else {
					timing.Status = constants.JOB_STATUS_STOP
				}
				providers.TimingService.UpdateTiming(&timing)
			}
		} else {
			logger.Error("checkOnlineAgent GetAgentTimingList Error:", err)
		}
	}
}
