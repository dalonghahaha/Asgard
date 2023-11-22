package master

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/dalonghahaha/avenger/components/cache"
	"github.com/dalonghahaha/avenger/components/db"
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/components/mail"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/registry"
	"Asgard/rpc"
	"Asgard/runtimes"
	"Asgard/server"
)

func GetCmd() *cobra.Command {
	masterCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	return masterCmd
}

var (
	leaderFlag bool
	rpcServer  *grpc.Server
)

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "run as master",
	Run: func(cmd *cobra.Command, args []string) {
		confPath := cmd.Flag("conf").Value.String()
		runtimes.ParseConfig(confPath)
		if err := InitMaster(); err != nil {
			fmt.Println(err)
			return
		}
		if constants.MASTER_CLUSTER {
			go RegisterRpcServer()
			go registry.Campaign("/Asgard/leader", constants.MASTER_CLUSTER_ID)
		}
		go StartRpcServer()
		go MoniterMaster()
		runtimes.Wait(StopMaster)
	},
}

func InitMaster() error {
	if err := logger.Register(); err != nil {
		return fmt.Errorf("init logger failed:%+v", err)
	}
	if err := db.Register(); err != nil {
		return fmt.Errorf("init db failed:%+v", err)
	}
	if err := cache.Register(); err != nil {
		return fmt.Errorf("init cache failed:%+v", err)
	}
	ip := viper.GetString("master.ip")
	if ip != "" {
		constants.MASTER_IP = ip
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
		if err := mail.Register(); err != nil {
			return fmt.Errorf("init mail failed:%+v", err)
		}
		constants.MASTER_NOTIFY = notify
		receiver := viper.GetString("master.receiver")
		if receiver == "" {
			return fmt.Errorf("receiver can not be empty when receiver enable!")
		}
		constants.MASTER_RECEIVER = receiver
		mailUser := viper.GetString("component.mail." + constants.MAIL_NAME + ".user")
		if mailUser == "" {
			return fmt.Errorf("mail user can not be empty!")
		}
		constants.MAIL_USER = mailUser
	}
	cluster := viper.GetBool("master.cluster")
	if cluster {
		constants.MASTER_CLUSTER = true
		constants.MASTER_CLUSTER_REGISTRY = viper.GetStringSlice("master.cluster_registry")
		constants.MASTER_CLUSTER_NAME = viper.GetString("master.cluster_name")
		constants.MASTER_CLUSTER_ID = viper.GetString("master.cluster_id")
		constants.MASTER_CLUSTER_IP = viper.GetString("master.cluster_ip")
		err := registry.RegisterRegistry(constants.MASTER_CLUSTER_REGISTRY)
		if err != nil {
			return fmt.Errorf("init registry failed:%+v", err)
		}
	}
	return nil
}

func RegisterRpcServer() {
	if err := recover(); err != nil {
		logger.Error("RegisterRpcServer panic:", err)
		runtimes.ExitSinal <- syscall.SIGTERM
		return
	}
	logger.Info("Master Rpc Server Registered!")
	registry.Register(
		constants.MASTER_CLUSTER_NAME,
		constants.MASTER_CLUSTER_ID,
		constants.MASTER_CLUSTER_IP,
		constants.MASTER_PORT,
	)
}

func StartRpcServer() {
	if err := recover(); err != nil {
		logger.Error("RegisterRpcServer panic:", err)
		runtimes.ExitSinal <- syscall.SIGTERM
		return
	}
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", constants.MASTER_PORT))
	if err != nil {
		logger.Error("failed to listen:", err)
		runtimes.ExitSinal <- syscall.SIGTERM
		return
	}
	rpcServer := server.NewRPCServer()
	rpc.RegisterMasterServer(rpcServer, &server.MasterServer{})
	reflection.Register(rpcServer)
	logger.Info("Master Rpc Server Started!")
	logger.Debugf("Server Port:%s", constants.MASTER_PORT)
	logger.Debugf("Server Pid:%d", os.Getpid())
	logger.Debugf("Moniter Notify:%v", constants.MASTER_NOTIFY)
	logger.Debugf("Moniter Loop:%d", constants.MASTER_MONITER)
	err = rpcServer.Serve(listen)
	if err != nil {
		logger.Error("failed to serve:", err)
		runtimes.ExitSinal <- syscall.SIGTERM
		return
	}
}

func StopMaster() {
	if rpcServer != nil {
		rpcServer.GracefulStop()
	}
	if constants.MASTER_CLUSTER {
		registry.UnRegister(constants.MASTER_CLUSTER_NAME, constants.MASTER_CLUSTER_ID)
	}
	constants.MASTER_TICKER.Stop()
	logger.Info("Master Rpc Server Stop!")
}

func MoniterMaster() {
	constants.MASTER_TICKER = time.NewTicker(time.Second * time.Duration(constants.MASTER_MONITER))
	for range constants.MASTER_TICKER.C {
		if !constants.MASTER_CLUSTER || registry.IsLeader() {
			logger.Debug("agent checking ......")
			agentList := providers.AgentService.GetUsageAgent()
			for _, agent := range agentList {
				go checkAgent(agent)
			}
		}
	}
}

func checkAgent(agent models.Agent) {
	usageApps := providers.AppService.GetUsageAppByAgentID(agent.ID)
	usageJobs := providers.JobService.GetUsageJobByAgentID(agent.ID)
	usageTimings := providers.TimingService.GetUsageTimingByAgentID(agent.ID)
	// client, err := providers.GetAgent(&agent)
	// if err != nil {
	// 	return
	// }
	// _, err = client.GetStat()

	//检查端口是否开启
	logger.Debugf("check agent:%s[%s:%s]", agent.Alias, agent.IP, agent.Port)
	tong := checkPort(agent.IP, agent.Port)

	if !tong {
		logger.Warnf("agent offline:%s[%s:%s] offline", agent.Alias, agent.IP, agent.Port)
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
		client, err := providers.GetAgent(&agent)
		if err != nil {
			return
		}
		//更新实例应用运行状态
		apps, err := client.GetAppList()
		if err != nil {
			logger.Error("checkOnlineAgent GetAgentAppList Error:", err)
		} else {
			markAppStatus(apps, usageApps)
		}
		//更新实例计划任务运行状态
		jobs, err := client.GetJobList()
		if err != nil {
			logger.Error("checkOnlineAgent GetAgentJobList Error:", err)
		} else {
			markJobStatus(jobs, usageJobs)
		}
		//更新实例计划任务运行状态
		timings, err := client.GetTimingList()
		if err != nil {
			logger.Error("checkOnlineAgent GetAgentTimingList Error:", err)
		} else {
			markTimigStatus(timings, usageTimings)
		}
	}
}

func checkPort(ip string, port string) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), time.Second*3)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
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
