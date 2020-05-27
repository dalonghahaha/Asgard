package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/providers"
	"Asgard/web/utils"
)

type MonitorController struct {
}

func NewMonitorController() *MonitorController {
	return &MonitorController{}
}

func (c *MonitorController) Agent(ctx *gin.Context) {
	agent := utils.GetAgent(ctx)
	cpus := []string{}
	memorys := []string{}
	times := []string{}
	moniters := providers.MoniterService.GetAgentMonitor(agent.ID, 100)
	for _, moniter := range moniters {
		cpus = append(cpus, utils.FormatFloat(moniter.CPU))
		memorys = append(memorys, utils.FormatFloat(moniter.Memory))
		times = append(times, utils.FormatTime(moniter.CreatedAt))
	}
	utils.Render(ctx, "monitor/list", gin.H{
		"Subtitle": "实例监控信息" + agent.Alias,
		"BackUrl":  utils.GetReferer(ctx),
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}

func (c *MonitorController) App(ctx *gin.Context) {
	app := utils.GetApp(ctx)
	moniters := providers.MoniterService.GetAppMonitor(app.ID, 100)
	cpus, memorys, times := utils.MonitorFormat(moniters)
	utils.Render(ctx, "monitor/list", gin.H{
		"Subtitle": "应用监控信息——" + app.Name,
		"BackUrl":  utils.GetReferer(ctx),
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}

func (c *MonitorController) Job(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	moniters := providers.MoniterService.GetJobMonitor(job.ID, 100)
	cpus, memorys, times := utils.MonitorFormat(moniters)
	utils.Render(ctx, "monitor/list", gin.H{
		"Subtitle": "计划任务监控信息——" + job.Name,
		"BackUrl":  utils.GetReferer(ctx),
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}

func (c *MonitorController) Timing(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	moniters := providers.MoniterService.GetTimingMonitor(timing.ID, 100)
	cpus, memorys, times := utils.MonitorFormat(moniters)
	utils.Render(ctx, "monitor/list", gin.H{
		"Subtitle": "定时任务监控信息——" + timing.Name,
		"BackUrl":  utils.GetReferer(ctx),
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}
