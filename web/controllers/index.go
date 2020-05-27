package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"Asgard/providers"
	"Asgard/web/utils"
)

type IndexController struct {
}

func NewIndexController() *IndexController {
	return &IndexController{}
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
	agents := providers.AgentService.GetUsageAgent()
	groups := providers.GroupService.GetUsageGroup()
	for _, agent := range agents {
		if agent.Alias == "" {
			agentList = append(agentList, fmt.Sprintf("%s:%s", agent.IP, agent.Port))
		} else {
			agentList = append(agentList, agent.Alias)
		}
		where := map[string]interface{}{"agent_id": agent.ID}
		apps := providers.AppService.GetAppCount(where)
		jobs := providers.JobService.GetJobCount(where)
		timings := providers.TimingService.GetTimingCount(where)
		agentApps = append(agentApps, apps)
		agentJobs = append(agentJobs, jobs)
		agentTimings = append(agentTimings, timings)
	}
	for _, group := range groups {
		groupList = append(groupList, group.Name)
		where := map[string]interface{}{"group_id": group.ID}
		apps := providers.AppService.GetAppCount(where)
		jobs := providers.JobService.GetJobCount(where)
		timings := providers.TimingService.GetTimingCount(where)
		groupApps = append(groupApps, apps)
		groupJobs = append(groupJobs, jobs)
		groupTimings = append(groupTimings, timings)
	}
	where := map[string]interface{}{}
	appCount := providers.AppService.GetAppCount(where)
	jobCount := providers.JobService.GetJobCount(where)
	timingCount := providers.TimingService.GetTimingCount(where)
	utils.Render(ctx, "index", gin.H{
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

func Nologin(ctx *gin.Context) {
	utils.Render(ctx, "error/no_login.html", gin.H{})
}

func AuthFail(ctx *gin.Context) {
	utils.Render(ctx, "error/auth_fail.html", gin.H{})
}

func AdminOnly(ctx *gin.Context) {
	utils.Render(ctx, "error/admin_only.html", gin.H{})
}

func Error(ctx *gin.Context) {
	utils.Render(ctx, "error/err.html", gin.H{
		"Subtitle": "服务器异常",
	})
}

func Forbidden(ctx *gin.Context) {
	utils.Render(ctx, "error/forbidden.html", gin.H{})
}
