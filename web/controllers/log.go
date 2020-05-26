package controllers

import (
	"Asgard/client"
	"Asgard/web/utils"

	"github.com/gin-gonic/gin"
)

type LogController struct {
}

func NewLogController() *LogController {
	return &LogController{}
}

func (c *LogController) AppOutLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, app.StdOut, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取失败:"+err.Error())
		return
	}
	ctx.HTML(StatusOK, "log/list", gin.H{
		"Subtitle": "应用正常日志查看",
		"Path":     "/out_log/app",
		"BackUrl":  GetReferer(ctx),
		"ID":       app.ID,
		"Name":     app.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) AppErrLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, app.StdErr, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取失败:"+err.Error())
		return
	}
	ctx.HTML(StatusOK, "log/list", gin.H{
		"Subtitle": "应用错误日志查看",
		"Path":     "/err_log/app",
		"BackUrl":  GetReferer(ctx),
		"ID":       app.ID,
		"Name":     app.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) AppOutLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, app.StdOut, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) AppErrLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, app.StdErr, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) JobOutLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, job.StdOut, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取失败:"+err.Error())
		return
	}
	ctx.HTML(StatusOK, "log/list", gin.H{
		"Subtitle": "计划任务日志查看",
		"Path":     "/out_log/job/",
		"BackUrl":  GetReferer(ctx),
		"ID":       job.ID,
		"Name":     job.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) JobErrLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, job.StdErr, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取失败:"+err.Error())
		return
	}
	ctx.HTML(StatusOK, "log/list", gin.H{
		"Subtitle": "计划任务错误日志查看",
		"Path":     "/err_log/job",
		"BackUrl":  GetReferer(ctx),
		"ID":       job.ID,
		"Name":     job.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) JobOutLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, job.StdOut, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) JobErrLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, job.StdErr, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) TimingOutLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, timing.StdOut, lines)
	if err != nil {
		utils.JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "log/list", gin.H{
		"Subtitle": "定时任务日志查看",
		"Path":     "/out_log/timing",
		"BackUrl":  GetReferer(ctx),
		"ID":       timing.ID,
		"Name":     timing.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) TimingErrLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, timing.StdErr, lines)
	if err != nil {
		utils.JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "log/list", gin.H{
		"Subtitle": "定时任务错误日志查看",
		"Path":     "/err_log/timing",
		"BackUrl":  GetReferer(ctx),
		"ID":       timing.ID,
		"Name":     timing.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Type":     "err_log",
		"Content":  content,
	})
}

func (c *LogController) TimingOutLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, timing.StdOut, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) TimingErrLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", LogSize)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	content, err := client.GetAgentLog(agent, timing.StdErr, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}
