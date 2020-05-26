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

func (c *MonitorController) App(ctx *gin.Context) {
	app := utils.GetApp(ctx)
	moniters := providers.MoniterService.GetAppMonitor(app.ID, 100)
	cpus, memorys, times := utils.MonitorFormat(moniters)
	ctx.HTML(StatusOK, "monitor/list", gin.H{
		"Subtitle": "应用监控信息——" + app.Name,
		"BackUrl":  GetReferer(ctx),
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}

func (c *MonitorController) Job(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	moniters := providers.MoniterService.GetJobMonitor(job.ID, 100)
	cpus, memorys, times := utils.MonitorFormat(moniters)
	ctx.HTML(StatusOK, "monitor/list", gin.H{
		"Subtitle": "计划任务监控信息——" + job.Name,
		"BackUrl":  GetReferer(ctx),
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}

func (c *MonitorController) Timing(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	moniters := providers.MoniterService.GetTimingMonitor(timing.ID, 100)
	cpus, memorys, times := utils.MonitorFormat(moniters)
	ctx.HTML(StatusOK, "monitor/list", gin.H{
		"Subtitle": "定时任务监控信息——" + timing.Name,
		"BackUrl":  GetReferer(ctx),
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}
