package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
	"Asgard/web/utils"
)

type LogController struct {
}

func NewLogController() *LogController {
	return &LogController{}
}

func (c *LogController) AppOutLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, app.StdOut, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取失败:"+err.Error())
		return
	}
	utils.Render(ctx, "log/list", gin.H{
		"Subtitle": "应用正常日志查看",
		"Path":     "/out_log/app",
		"BackUrl":  utils.GetReferer(ctx),
		"ID":       app.ID,
		"Name":     app.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) AppErrLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, app.StdErr, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	utils.Render(ctx, "log/list", gin.H{
		"Subtitle": "应用错误日志查看",
		"Path":     "/err_log/app",
		"BackUrl":  utils.GetReferer(ctx),
		"ID":       app.ID,
		"Name":     app.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) AppOutLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, app.StdOut, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) AppErrLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	app := utils.GetApp(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, app.StdErr, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) JobOutLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, job.StdOut, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	utils.Render(ctx, "log/list", gin.H{
		"Subtitle": "计划任务日志查看",
		"Path":     "/out_log/job/",
		"BackUrl":  utils.GetReferer(ctx),
		"ID":       job.ID,
		"Name":     job.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) JobErrLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, job.StdErr, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	utils.Render(ctx, "log/list", gin.H{
		"Subtitle": "计划任务错误日志查看",
		"Path":     "/err_log/job",
		"BackUrl":  utils.GetReferer(ctx),
		"ID":       job.ID,
		"Name":     job.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) JobOutLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, job.StdOut, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) JobErrLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, job.StdErr, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) TimingOutLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, timing.StdOut, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	utils.Render(ctx, "log/list", gin.H{
		"Subtitle": "定时任务日志查看",
		"Path":     "/out_log/timing",
		"BackUrl":  utils.GetReferer(ctx),
		"ID":       timing.ID,
		"Name":     timing.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *LogController) TimingErrLog(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, timing.StdErr, lines)
	if err != nil {
		utils.JumpWarning(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	utils.Render(ctx, "log/list", gin.H{
		"Subtitle": "定时任务错误日志查看",
		"Path":     "/err_log/timing",
		"BackUrl":  utils.GetReferer(ctx),
		"ID":       timing.ID,
		"Name":     timing.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Type":     "err_log",
		"Content":  content,
	})
}

func (c *LogController) TimingOutLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, timing.StdOut, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}

func (c *LogController) TimingErrLogData(ctx *gin.Context) {
	lines := utils.DefaultInt64(ctx, "lines", constants.WEB_LOG_SIZE)
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "获取日志失败:\n"+err.Error())
		return
	}
	content, err := client.GetLog(agent, timing.StdErr, lines)
	if err != nil {
		utils.APIError(ctx, err.Error())
		return
	}
	utils.APIData(ctx, content)
}
