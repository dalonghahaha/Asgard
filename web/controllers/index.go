package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"Asgard/services"
)

type IndexController struct {
	agentService  *services.AgentService
	groupService  *services.GroupService
	appService    *services.AppService
	jobService    *services.JobService
	timingService *services.TimingService
}

func NewIndexController() *IndexController {
	return &IndexController{
		appService:    services.NewAppService(),
		groupService:  services.NewGroupService(),
		agentService:  services.NewAgentService(),
		jobService:    services.NewJobService(),
		timingService: services.NewTimingService(),
	}
}

func (c *IndexController) Index(ctx *gin.Context) {
	agentList := []string{}
	groupList := []string{}
	agentApps := []int{}
	agentJobs := []int{}
	agentTimings := []int{}
	groupApps := []int{}
	groupJobs := []int{}
	groupTimings := []int{}
	agents := c.agentService.GetUsageAgent()
	groups := c.groupService.GetUsageGroup()
	for _, agent := range agents {
		agentList = append(agentList, fmt.Sprintf("%s:%s", agent.IP, agent.Port))
		where := map[string]interface{}{"agent_id": agent.ID}
		apps := c.appService.GetAppCount(where)
		jobs := c.jobService.GetJobCount(where)
		timings := c.timingService.GetTimingCount(where)
		agentApps = append(agentApps, apps)
		agentJobs = append(agentJobs, jobs)
		agentTimings = append(agentTimings, timings)
	}
	for _, group := range groups {
		groupList = append(groupList, group.Name)
		where := map[string]interface{}{"group_id": group.ID}
		apps := c.appService.GetAppCount(where)
		jobs := c.jobService.GetJobCount(where)
		timings := c.timingService.GetTimingCount(where)
		groupApps = append(groupApps, apps)
		groupJobs = append(groupJobs, jobs)
		groupTimings = append(groupTimings, timings)
	}
	where := map[string]interface{}{}
	appCount := c.appService.GetAppCount(where)
	jobCount := c.jobService.GetJobCount(where)
	timingCount := c.timingService.GetTimingCount(where)
	ctx.HTML(StatusOK, "index", gin.H{
		"Subtitle":     "首页",
		"Agents":       len(agents),
		"Apps":         appCount,
		"Jobs":         jobCount,
		"Timings":      timingCount,
		"AgentList":    agentList,
		"GroupList":    groupList,
		"AgentApps":    agentApps,
		"AgentJobs":    agentJobs,
		"AgentTimings": agentTimings,
		"GroupApps":    groupApps,
		"GroupJobs":    groupJobs,
		"GroupTimings": groupTimings,
	})
}

func UI(c *gin.Context) {
	c.HTML(StatusOK, "UI", gin.H{
		"Subtitle": "布局",
	})
}

func Nologin(c *gin.Context) {
	c.HTML(StatusOK, "error/no_login.html", gin.H{
		"Subtitle": "未登录提示页",
	})
}

func AuthFail(c *gin.Context) {
	c.HTML(StatusOK, "error/auth_fail.html", gin.H{
		"Subtitle": "登录验证失败页",
	})
}

func Error(c *gin.Context) {
	c.HTML(StatusOK, "error/err.html", gin.H{
		"Subtitle": "服务器异常",
	})
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
