package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/services"
)

type IndexController struct {
	agentService  *services.AgentService
	appService    *services.AppService
	jobService    *services.JobService
	timingService *services.TimingService
}

func NewIndexController() *IndexController {
	return &IndexController{
		appService:    services.NewAppService(),
		agentService:  services.NewAgentService(),
		jobService:    services.NewJobService(),
		timingService: services.NewTimingService(),
	}
}

func (c *IndexController) Index(ctx *gin.Context) {
	where := map[string]interface{}{}
	agentCount := c.agentService.GetAgentCount(where)
	appCount := c.appService.GetAppCount(where)
	jobCount := c.jobService.GetJobCount(where)
	timingCount := c.timingService.GetTimingCount(where)
	ctx.HTML(StatusOK, "index", gin.H{
		"Subtitle": "首页",
		"Agents":   agentCount,
		"Apps":     appCount,
		"Jobs":     jobCount,
		"Timings":  timingCount,
	})
}

func UI(c *gin.Context) {
	c.HTML(StatusOK, "UI", gin.H{
		"Subtitle": "布局",
	})
}

func Nologin(c *gin.Context) {
	c.HTML(StatusOK, "error/nologin.html", gin.H{
		"Subtitle": "未登录提示页",
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
