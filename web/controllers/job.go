package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"

	"Asgard/models"
	"Asgard/services"
)

type JobController struct {
	jobService   *services.JobService
	agentService *services.AgentService
	groupService *services.GroupService
}

func NewJobController() *JobController {
	return &JobController{
		jobService:   services.NewJobService(),
		agentService: services.NewAgentService(),
		groupService: services.NewGroupService(),
	}
}

func (c *JobController) formatJob(info *models.Job) map[string]interface{} {
	data := map[string]interface{}{
		"ID":        info.ID,
		"Name":      info.Name,
		"GroupID":   info.GroupID,
		"AgentID":   info.AgentID,
		"Dir":       info.Dir,
		"Program":   info.Program,
		"Args":      info.Args,
		"StdOut":    info.StdOut,
		"StdErr":    info.StdErr,
		"Spec":      info.Spec,
		"Timeout":   info.Timeout,
		"IsMonitor": info.IsMonitor,
		"Status":    info.Status,
	}
	group := c.groupService.GetGroupByID(info.GroupID)
	if group != nil {
		data["GroupName"] = group.Name
	} else {
		data["GroupName"] = ""
	}
	agent := c.agentService.GetAgentByID(info.AgentID)
	if group != nil {
		data["AgentName"] = agent.IP + ":" + agent.Port
	} else {
		data["AgentName"] = ""
	}
	return data
}

func (c *JobController) List(ctx *gin.Context) {
	page := DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	jobList, total := c.jobService.GetJobPageList(where, page, PageSize)
	if jobList == nil {
		APIError(ctx, "获取计划任务列表失败")
	}
	list := []map[string]interface{}{}
	for _, job := range jobList {
		list = append(list, c.formatJob(&job))
	}
	mpurl := "/job/list"
	ctx.HTML(StatusOK, "job/list", gin.H{
		"Subtitle":   "计划任务列表",
		"List":       list,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
}

func (c *JobController) Add(ctx *gin.Context) {
	groups := c.groupService.GetAllGroup()
	agents := c.agentService.GetAllAgent()
	ctx.HTML(StatusOK, "job/add", gin.H{
		"Subtitle":  "添加计划任务",
		"GroupList": groups,
		"AgentList": agents,
	})
}

func (c *JobController) Create(ctx *gin.Context) {
	groupID := FormDefaultInt(ctx, "group_id", 0)
	name := ctx.PostForm("name")
	agentID := FormDefaultInt(ctx, "agent_id", 0)
	dir := ctx.PostForm("dir")
	program := ctx.PostForm("program")
	args := ctx.PostForm("args")
	stdOut := ctx.PostForm("std_out")
	stdErr := ctx.PostForm("std_err")
	if !Required(ctx, &name, "名称不能为空") {
		return
	}
	if !Required(ctx, &dir, "执行目录不能为空") {
		return
	}
	if !Required(ctx, &program, "执行程序不能为空") {
		return
	}
	if agentID == 0 {
		APIBadRequest(ctx, "运行实例不能为空")
		return
	}
	job := new(models.Job)
	job.GroupID = int64(groupID)
	job.Name = name
	job.AgentID = int64(agentID)
	job.Dir = dir
	job.Program = program
	job.Args = args
	job.StdOut = stdOut
	job.StdErr = stdErr
	job.Status = 0
	job.Creator = GetUserID(ctx)
	ok := c.jobService.CreateJob(job)
	if !ok {
		APIError(ctx, "创建计划任务失败")
		return
	}
	APIOK(ctx)
}

func (c *JobController) Edit(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		JumpError(ctx)
		return
	}
	job := c.jobService.GetJobByID(int64(id))
	if job == nil {
		JumpError(ctx)
		return
	}
	groups := c.groupService.GetAllGroup()
	agents := c.agentService.GetAllAgent()
	ctx.HTML(StatusOK, "job/edit", gin.H{
		"Subtitle":  "编辑计划任务",
		"Job":       c.formatJob(job),
		"GroupList": groups,
		"AgentList": agents,
	})
}

func (c *JobController) Update(ctx *gin.Context) {
	id := FormDefaultInt(ctx, "id", 0)
	name := ctx.PostForm("name")
	status := ctx.PostForm("status")
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	if name == "" && status == "" {
		APIBadRequest(ctx, "请求数据格式错误")
		return
	}
	job := c.jobService.GetJobByID(int64(id))
	if job == nil {
		APIBadRequest(ctx, "分组不存在")
		return
	}
	if name != "" {
		job.Name = name
	}
	if status != "" {
		_status, err := strconv.ParseInt(status, 10, 64)
		if err != nil {
			APIBadRequest(ctx, "status格式错误")
			return
		}
		job.Status = _status
	}
	job.Updator = GetUserID(ctx)
	ok := c.jobService.UpdateJob(job)
	if !ok {
		APIError(ctx, "更新计划任务失败")
		return
	}
	APIOK(ctx)
}

func (c *JobController) Delete(ctx *gin.Context) {
	id := FormDefaultInt(ctx, "id", 0)
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	ok := c.jobService.DeleteJobByID(int64(id))
	if !ok {
		APIError(ctx, "删除计划任务失败")
		return
	}
	APIOK(ctx)
}
