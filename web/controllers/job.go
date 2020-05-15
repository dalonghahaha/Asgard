package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"Asgard/client"
	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/web/utils"
)

type JobController struct {
}

func NewJobController() *JobController {
	return &JobController{}
}

func (c *JobController) List(ctx *gin.Context) {
	groupID := utils.DefaultInt(ctx, "group_id", 0)
	agentID := utils.DefaultInt(ctx, "agent_id", 0)
	status := utils.DefaultInt(ctx, "status", -99)
	name := ctx.Query("name")
	page := utils.DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{
		"status": status,
	}
	querys := []string{}
	if groupID != 0 {
		where["group_id"] = groupID
		querys = append(querys, "group_id="+strconv.Itoa(groupID))
	}
	if agentID != 0 {
		where["agent_id"] = agentID
		querys = append(querys, "agent_id="+strconv.Itoa(agentID))
	}
	if status != -99 {
		querys = append(querys, "status="+strconv.Itoa(status))
	}
	if name != "" {
		where["name"] = name
		querys = append(querys, "name="+name)
	}
	jobList, total := providers.JobService.GetJobPageList(where, page, PageSize)
	if jobList == nil {
		utils.APIError(ctx, "获取计划任务列表失败")
	}
	list := []gin.H{}
	for _, job := range jobList {
		list = append(list, utils.JobFormat(&job))
	}
	mpurl := "/job/list"
	if len(querys) > 0 {
		mpurl = "/job/list?" + strings.Join(querys, "&")
	}
	ctx.HTML(StatusOK, "job/list", gin.H{
		"Subtitle":   "计划任务列表",
		"List":       list,
		"Total":      total,
		"GroupList":  providers.GroupService.GetUsageGroup(),
		"AgentList":  providers.AgentService.GetUsageAgent(),
		"StatusList": constants.JOB_STATUS,
		"GroupID":    groupID,
		"AgentID":    agentID,
		"Name":       name,
		"Status":     status,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *JobController) Show(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	ctx.HTML(StatusOK, "job/show", gin.H{
		"Subtitle": "查看计划任务",
		"Job":      utils.JobFormat(job),
	})
}

func (c *JobController) Monitor(ctx *gin.Context) {
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

func (c *JobController) Archive(ctx *gin.Context) {
	page := utils.DefaultInt(ctx, "page", 1)
	job := utils.GetJob(ctx)
	where := map[string]interface{}{
		"type":       constants.TYPE_JOB,
		"related_id": job.ID,
	}
	archiveList, total := providers.ArchiveService.GetArchivePageList(where, page, PageSize)
	if archiveList == nil {
		utils.APIError(ctx, "获取归档列表失败")
	}
	list := []map[string]interface{}{}
	for _, archive := range archiveList {
		list = append(list, formatArchive(&archive))
	}
	mpurl := fmt.Sprintf("/job/archive?id=%d", job.ID)
	ctx.HTML(StatusOK, "archive/list", gin.H{
		"Subtitle":   "计划任务归档列表——" + job.Name,
		"BackUrl":    GetReferer(ctx),
		"List":       list,
		"Total":      total,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *JobController) OutLog(ctx *gin.Context) {
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
		"Path":     "/job/out_log",
		"BackUrl":  GetReferer(ctx),
		"ID":       job.ID,
		"Name":     job.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *JobController) ErrLog(ctx *gin.Context) {
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
		"Path":     "/job/err_log",
		"BackUrl":  GetReferer(ctx),
		"ID":       job.ID,
		"Name":     job.Name,
		"Agent":    agent,
		"Lines":    lines,
		"Content":  content,
	})
}

func (c *JobController) Add(ctx *gin.Context) {
	ctx.HTML(StatusOK, "job/add", gin.H{
		"Subtitle":   "添加计划任务",
		"OutBaseDir": OutDir + "cron/",
		"GroupList":  providers.GroupService.GetUsageGroup(),
		"AgentList":  providers.AgentService.GetUsageAgent(),
	})
}

func (c *JobController) Create(ctx *gin.Context) {
	if !utils.Required(ctx, ctx.PostForm("spec"), "运行配置不能为空") {
		return
	}
	job := new(models.Job)
	job.GroupID = utils.FormDefaultInt64(ctx, "group_id", 0)
	job.AgentID = utils.FormDefaultInt64(ctx, "agent_id", 0)
	job.Name = ctx.PostForm("name")
	job.Dir = ctx.PostForm("dir")
	job.Program = ctx.PostForm("program")
	job.Args = ctx.PostForm("args")
	job.StdOut = ctx.PostForm("std_out")
	job.StdErr = ctx.PostForm("std_err")
	job.Spec = ctx.PostForm("spec")
	job.Timeout = utils.FormDefaultInt64(ctx, "timeout", -1)
	job.Status = constants.JOB_STATUS_STOP
	job.Creator = GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		job.IsMonitor = 1
	}
	ok := providers.JobService.CreateJob(job)
	if !ok {
		utils.APIError(ctx, "创建计划任务失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *JobController) Edit(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	ctx.HTML(StatusOK, "job/edit", gin.H{
		"Subtitle":  "编辑计划任务",
		"BackUrl":   GetReferer(ctx),
		"Info":      utils.JobFormat(job),
		"GroupList": providers.GroupService.GetUsageGroup(),
		"AgentList": providers.AgentService.GetUsageAgent(),
	})
}

func (c *JobController) Update(ctx *gin.Context) {
	if !utils.Required(ctx, ctx.PostForm("spec"), "运行配置不能为空") {
		return
	}
	job := utils.GetJob(ctx)
	job.GroupID = utils.FormDefaultInt64(ctx, "group_id", 0)
	job.AgentID = utils.FormDefaultInt64(ctx, "agent_id", 0)
	job.Name = ctx.PostForm("name")
	job.Dir = ctx.PostForm("dir")
	job.Program = ctx.PostForm("program")
	job.Args = ctx.PostForm("args")
	job.StdOut = ctx.PostForm("std_out")
	job.StdErr = ctx.PostForm("std_err")
	job.Spec = ctx.PostForm("spec")
	job.Timeout = utils.FormDefaultInt64(ctx, "timeout", -1)
	job.Updator = GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		job.IsMonitor = 1
	}
	ok := providers.JobService.UpdateJob(job)
	if !ok {
		utils.APIError(ctx, "更新计划任务失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *JobController) Copy(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	_job := new(models.Job)
	_job.GroupID = job.GroupID
	_job.Name = job.Name + "_copy"
	_job.AgentID = job.AgentID
	_job.Dir = job.Dir
	_job.Program = job.Program
	_job.Args = job.Args
	_job.StdOut = job.StdOut
	_job.StdErr = job.StdErr
	_job.Spec = job.Spec
	_job.Timeout = job.Timeout
	_job.IsMonitor = job.IsMonitor
	_job.Status = constants.JOB_STATUS_STOP
	_job.Creator = GetUserID(ctx)
	ok := providers.JobService.CreateJob(_job)
	if !ok {
		utils.APIError(ctx, "复制计划任务失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *JobController) Delete(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	if job.Status == 1 {
		utils.APIError(ctx, "计划任务正在运行不能删除")
		return
	}
	job.Status = constants.JOB_STATUS_DELETED
	job.Updator = GetUserID(ctx)
	ok := providers.JobService.UpdateJob(job)
	if !ok {
		utils.APIError(ctx, "删除计划任务失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *JobController) Start(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	if job.Status == constants.JOB_STATUS_RUNNING {
		utils.APIError(ctx, "计划任务已经启动")
		return
	}
	_job, err := client.GetAgentJob(agent, job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取计划任务情况异常:%s", err.Error()))
		return
	}
	if _job == nil {
		err = client.AddAgentJob(agent, job)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("添加计划任务异常:%s", err.Error()))
			return
		}
		job.Status = constants.JOB_STATUS_RUNNING
		job.Updator = GetUserID(ctx)
		providers.JobService.UpdateJob(job)
		utils.APIOK(ctx)
		return
	}
	job.Status = constants.JOB_STATUS_RUNNING
	job.Updator = GetUserID(ctx)
	providers.JobService.UpdateJob(job)
	utils.APIOK(ctx)
}

func (c *JobController) ReStart(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	_job, err := client.GetAgentJob(agent, job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取计划任务情况异常:%s", err.Error()))
		return
	}
	if _job == nil {
		err = client.AddAgentJob(agent, job)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("重启异常:%s", err.Error()))
			return
		}
		utils.APIOK(ctx)
	}
	err = client.UpdateAgentJob(agent, job)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("重启异常:%s", err.Error()))
		return
	}
	utils.APIOK(ctx)
}

func (c *JobController) Pause(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	_job, err := client.GetAgentJob(agent, job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取计划任务情况异常:%s", err.Error()))
		return
	}
	if _job == nil {
		utils.APIOK(ctx)
		return
	}
	err = client.RemoveAgentJob(agent, job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("停止计划任务异常:%s", err.Error()))
		return
	}
	job.Status = constants.JOB_STATUS_PAUSE
	job.Updator = GetUserID(ctx)
	providers.JobService.UpdateJob(job)
	utils.APIOK(ctx)
}
