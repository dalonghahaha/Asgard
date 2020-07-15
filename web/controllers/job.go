package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/gin-gonic/gin"

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
	user := utils.GetUser(ctx)
	where := map[string]interface{}{
		"status": status,
	}
	querys := []string{}
	if user.Role != constants.USER_ROLE_ADMIN {
		where["creator"] = user.ID
	}
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
	jobList, total := providers.JobService.GetJobPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
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
	utils.Render(ctx, "job/list", gin.H{
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
	utils.Render(ctx, "job/show", gin.H{
		"Subtitle": "查看计划任务",
		"Job":      utils.JobFormat(job),
	})
}

func (c *JobController) Add(ctx *gin.Context) {
	utils.Render(ctx, "job/add", gin.H{
		"Subtitle":   "添加计划任务",
		"OutBaseDir": constants.WEB_OUT_DIR + "cron/",
		"GroupList":  providers.GroupService.GetUsageGroup(),
		"AgentList":  providers.AgentService.GetUsageAgent(),
	})
}

func (c *JobController) Create(ctx *gin.Context) {
	if utils.FormDefaultInt64(ctx, "agent_id", 0) == 0 {
		utils.APIError(ctx, "运行实例未选择")
		return
	}
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
	job.Status = constants.JOB_STATUS_PAUSE
	job.Creator = utils.GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		job.IsMonitor = 1
	}
	ok := providers.JobService.CreateJob(job)
	if !ok {
		utils.APIError(ctx, "创建计划任务失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_CREATE)
	utils.APIOK(ctx)
}

func (c *JobController) Edit(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	utils.Render(ctx, "job/edit", gin.H{
		"Subtitle":  "编辑计划任务",
		"BackUrl":   utils.GetReferer(ctx),
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
	job.Name = ctx.PostForm("name")
	job.Dir = ctx.PostForm("dir")
	job.Program = ctx.PostForm("program")
	job.Args = ctx.PostForm("args")
	job.StdOut = ctx.PostForm("std_out")
	job.StdErr = ctx.PostForm("std_err")
	job.Spec = ctx.PostForm("spec")
	job.Timeout = utils.FormDefaultInt64(ctx, "timeout", -1)
	job.Updator = utils.GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		job.IsMonitor = 1
	} else {
		job.IsMonitor = 0
	}
	if utils.FormDefaultInt64(ctx, "agent_id", 0) != 0 {
		job.AgentID = utils.FormDefaultInt64(ctx, "agent_id", 0)
	}
	ok := providers.JobService.UpdateJob(job)
	if !ok {
		utils.APIError(ctx, "更新计划任务失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_UPDATE)
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
	_job.Status = constants.JOB_STATUS_PAUSE
	_job.Creator = utils.GetUserID(ctx)
	ok := providers.JobService.CreateJob(_job)
	if !ok {
		utils.APIError(ctx, "复制计划任务失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_COPY)
	utils.APIOK(ctx)
}

func (c *JobController) Start(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	if job.Status == constants.JOB_STATUS_RUNNING {
		utils.APIError(ctx, "计划任务已经启动")
		return
	}
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_job, err := client.GetJob(job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取计划任务情况异常:%s", err.Error()))
		return
	}
	if _job == nil {
		err = client.AddJob(job)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("添加计划任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_RUNNING, utils.GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新计划任务状态失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_START)
	utils.APIOK(ctx)
}

func (c *JobController) ReStart(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_job, err := client.GetJob(job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取计划任务情况异常:%s", err.Error()))
		return
	}
	if _job == nil {
		err = client.AddJob(job)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("重启计划任务异常:%s", err.Error()))
			return
		}
	} else {
		err = client.UpdateJob(job)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("重启计划任务异常:%s", err.Error()))
			return
		}
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_RESTART)
	utils.APIOK(ctx)
}

func (c *JobController) Pause(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_job, err := client.GetJob(job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取计划任务情况异常:%s", err.Error()))
		return
	}
	if _job != nil {
		err = client.RemoveJob(job.ID)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("停止计划任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_PAUSE, utils.GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新计划任务状态失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_PAUSE)
	utils.APIOK(ctx)
}

func (c *JobController) Delete(ctx *gin.Context) {
	job := utils.GetJob(ctx)
	agent := utils.GetAgent(ctx)
	if job.Status != constants.JOB_STATUS_PAUSE {
		utils.APIError(ctx, "计划任务启动状态不能删除")
		return
	}
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_job, err := client.GetJob(job.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取计划任务情况异常:%s", err.Error()))
		return
	}
	if _job != nil {
		err = client.RemoveJob(job.ID)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("停止计划任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_DELETED, utils.GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新计划任务状态失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_DELETE)
	utils.APIOK(ctx)
}

func (c *JobController) BatchStart(ctx *gin.Context) {
	jobAgent := utils.GetJobAgent(ctx)
	for job, agent := range jobAgent {
		if job.Status == constants.JOB_STATUS_RUNNING {
			continue
		}
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("App BatchStart GetAgent Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		_job, err := client.GetJob(job.ID)
		if err != nil {
			logger.Errorf("Job BatchStart GetJob Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		if _job == nil {
			err = client.AddJob(job)
			if err != nil {
				logger.Errorf("Job BatchStart AddJob Error:%s", err.Error())
			}
		}
		providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_RUNNING, utils.GetUserID(ctx))
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_START)
	}
	utils.APIOK(ctx)
}

func (c *JobController) BatchReStart(ctx *gin.Context) {
	jobAgent := utils.GetJobAgent(ctx)
	for job, agent := range jobAgent {
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("Job BatchReStart GetAgent Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		_job, err := client.GetJob(job.ID)
		if err != nil {
			logger.Errorf("Job BatchReStart GetAgentJob Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		if _job == nil {
			err = client.AddJob(job)
			if err != nil {
				logger.Errorf("Job BatchReStart AddAgentJob Error:[%d][%s]", job.ID, err.Error())
			}
		} else {
			err = client.UpdateJob(job)
			if err != nil {
				logger.Errorf("Job BatchReStart UpdateAgentJob Error:[%d][%s]", job.ID, err.Error())
			}
		}
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_RESTART)
	}
	utils.APIOK(ctx)
}

func (c *JobController) BatchPause(ctx *gin.Context) {
	jobAgent := utils.GetJobAgent(ctx)
	for job, agent := range jobAgent {
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("Job BatchPause GetAgent Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		_job, err := client.GetJob(job.ID)
		if err != nil {
			logger.Errorf("Job BatchPause GetAgentJob Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		if _job != nil {
			err = client.RemoveJob(job.ID)
			if err != nil {
				logger.Errorf("Job BatchPause RemoveAgentJob Error:[%d][%s]", job.ID, err.Error())
				return
			}
		}
		providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_PAUSE, utils.GetUserID(ctx))
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_PAUSE)
	}
	utils.APIOK(ctx)
}

func (c *JobController) BatchDelete(ctx *gin.Context) {
	jobAgent := utils.GetJobAgent(ctx)
	for job, agent := range jobAgent {
		if job.Status == constants.JOB_STATUS_RUNNING {
			continue
		}
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("Job BatchDelete GetAgent Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		_job, err := client.GetJob(job.ID)
		if err != nil {
			logger.Errorf("Job BatchDelete GetJob Error:[%d][%s]", job.ID, err.Error())
			continue
		}
		if _job != nil {
			err = client.RemoveJob(job.ID)
			if err != nil {
				logger.Errorf("Job BatchDelete RemoveJob Error:[%d][%s]", job.ID, err.Error())
				return
			}
		}
		providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_DELETED, utils.GetUserID(ctx))
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_DELETE)
	}
	utils.APIOK(ctx)
}
