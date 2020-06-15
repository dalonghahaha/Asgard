package master

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dalonghahaha/avenger/components/cache"
	"github.com/dalonghahaha/avenger/components/db"
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/components/mail"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/reflection"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/rpc"
	"Asgard/runtimes"
	"Asgard/server"
)

func GetCmd() *cobra.Command {
	masterCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	return masterCmd
}

func PreRun(cmd *cobra.Command, args []string) {
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
	systemMoniter := viper.GetInt("system.moniter")
	if systemMoniter > 0 {
		constants.SYSTEM_MONITER = systemMoniter
	}
	systemTimer := viper.GetInt("system.timer")
	if systemMoniter > 0 {
		constants.SYSTEM_TIMER = systemTimer
	}
}

var masterCmd = &cobra.Command{
	Use:    "master",
	Short:  "run as master",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		InitMaster()
		go StartMasterRpcServer()
		go MoniterMaster()
		runtimes.Wait(StopMaster)
	},
}

func InitMaster() {
	err := db.Register()
	if err != nil {
		panic(err)
	}
	err = cache.Register()
	if err != nil {
		panic(err)
	}
	port := viper.GetString("master.port")
	if port != "" {
		constants.MASTER_PORT = port
	}
	moniter := viper.GetInt("master.moniter")
	if moniter != 0 {
		constants.MASTER_MONITER = moniter
	}
	notify := viper.GetBool("master.notify")
	if notify {
		err = mail.Register()
		if err != nil {
			panic(err)
		}
		constants.MASTER_NOTIFY = notify
		receiver := viper.GetString("master.receiver")
		if receiver == "" {
			panic("receiver can not be empty when receiver enable!")
		}
		constants.MASTER_RECEIVER = receiver
		mailUser := viper.GetString("component.mail." + constants.MAIL_NAME + ".user")
		if mailUser == "" {
			panic(fmt.Errorf("mail user can not be empty!"))
		}
		constants.MAIL_USER = mailUser
	}
}

func StartMasterRpcServer() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", constants.MASTER_PORT))
	if err != nil {
		logger.Error("failed to listen:", err)
		panic(err)
	}
	s := server.NewRPCServer()
	rpc.RegisterMasterServer(s, &server.MasterServer{})
	reflection.Register(s)
	logger.Info("Master Rpc Server Started!")
	logger.Debugf("Server Port:%s", constants.MASTER_PORT)
	logger.Debugf("Server Pid:%d", os.Getpid())
	logger.Debugf("Moniter Notify:%v", constants.MASTER_NOTIFY)
	logger.Debugf("Moniter Loop:%d", constants.MASTER_MONITER)
	err = s.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		panic(err)
	}
}

func StopMaster() {
	logger.Info("Master Rpc Server Stop!")
	constants.MASTER_TICKER.Stop()
}

func MoniterMaster() {
	constants.MASTER_TICKER = time.NewTicker(time.Second * time.Duration(constants.MASTER_MONITER))
	for range constants.MASTER_TICKER.C {
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
	client, err := providers.GetAgent(&agent)
	if err != nil {
		return
	}
	_, err = client.GetStat()
	if err != nil {
		//标记实例状态为离线
		agent.Status = constants.AGENT_OFFLINE
		providers.AgentService.UpdateAgent(&agent)
		//标记应用状态为未知
		for _, app := range usageApps {
			providers.AppService.ChangeAPPStatus(&app, constants.APP_STATUS_UNKNOWN, 0)
		}
		//标记计划任务状态为未知
		for _, job := range usageJobs {
			providers.JobService.ChangeJobStatus(&job, constants.JOB_STATUS_UNKNOWN, 0)
		}
		//标记定时任务状态为未知
		for _, timing := range usageTimings {
			providers.TimingService.ChangeTimingStatus(&timing, constants.TIMING_STATUS_UNKNOWN, 0)
		}
		return
	} else {
		//标记实例状态为在线
		agent.Status = constants.AGENT_ONLINE
		providers.AgentService.UpdateAgent(&agent)
		//更新实例应用运行状态
		apps, err := client.GetAppList()
		//logger.Debugf("%s runing apps:%d", agent.Alias, len(apps))
		if err != nil {
			logger.Error("checkOnlineAgent GetAgentAppList Error:", err)
		} else {
			markAppStatus(apps, usageApps)
		}
		//更新实例计划任务运行状态
		jobs, err := client.GetJobList()
		//logger.Debugf("%s runing jobs:%d", agent.Alias, len(jobs))
		if err != nil {
			logger.Error("checkOnlineAgent GetAgentJobList Error:", err)
		} else {
			markJobStatus(jobs, usageJobs)
		}
		//更新实例计划任务运行状态
		timings, err := client.GetTimingList()
		//logger.Debugf("%s runing timings:%d", agent.Alias, len(timings))
		if err != nil {
			logger.Error("checkOnlineAgent GetAgentTimingList Error:", err)
		} else {
			markTimigStatus(timings, usageTimings)
		}
	}
}

func markAppStatus(apps []*rpc.App, usageApps []models.App) {
	runningApps := map[int64]string{}
	for _, app := range apps {
		runningApps[app.GetId()] = app.GetName()
	}
	for _, app := range usageApps {
		_, ok := runningApps[app.ID]
		if ok {
			providers.AppService.ChangeAPPStatus(&app, constants.APP_STATUS_RUNNING, 0)
		} else {
			providers.AppService.ChangeAPPStatus(&app, constants.APP_STATUS_STOP, 0)
		}
	}
}

func markJobStatus(jobs []*rpc.Job, usageJobs []models.Job) {
	runningJobs := map[int64]string{}
	for _, job := range jobs {
		runningJobs[job.GetId()] = job.GetName()
	}
	for _, job := range usageJobs {
		_, ok := runningJobs[job.ID]
		if ok {
			providers.JobService.ChangeJobStatus(&job, constants.JOB_STATUS_RUNNING, 0)
		} else {
			providers.JobService.ChangeJobStatus(&job, constants.JOB_STATUS_STOP, 0)
		}
	}
}

func markTimigStatus(timings []*rpc.Timing, usageTimings []models.Timing) {
	runningTimings := map[int64]string{}
	for _, timing := range timings {
		runningTimings[timing.GetId()] = timing.GetName()
	}
	for _, timing := range usageTimings {
		_, ok := runningTimings[timing.ID]
		if ok {
			providers.TimingService.ChangeTimingStatus(&timing, constants.TIMING_STATUS_RUNNING, 0)
		} else {
			providers.TimingService.ChangeTimingStatus(&timing, constants.TIMING_STATUS_STOP, 0)
		}
	}
}
