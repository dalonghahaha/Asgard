 package controllers

 import (
 	"fmt"
 	"time"

 	"github.com/dalonghahaha/avenger/components/logger"
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 type apiJobReq struct {
 	GroupID   int64  `form:"group_id" json:"group_id"`
 	AgentID   int64  `form:"agent_id" json:"agent_id" binding:"required"`
 	Name      string `form:"name" json:"name" binding:"required"`
 	Dir       string `form:"dir" json:"dir" binding:"required"`
 	Program   string `form:"program" json:"program" binding:"required"`
 	Args      string `form:"args" json:"args"`
 	StdOut    string `form:"std_out" json:"std_out" binding:"required"`
 	StdErr    string `form:"std_err" json:"std_err" binding:"required"`
 	Spec      string `form:"spec" json:"spec" binding:"required"`
 	Timeout   int64  `form:"timeout" json:"timeout"`
 	IsMonitor int    `form:"is_monitor" json:"is_monitor"`
 }

 type APIJobController struct{}

 func NewAPIJobController() *APIJobController {
 	return &APIJobController{}
 }

 func (c *APIJobController) List(ctx *gin.Context) {
 	page := utils.QueryInt(ctx, "page", 1)
 	groupID := utils.QueryInt64(ctx, "group_id", 0)
 	agentID := utils.QueryInt64(ctx, "agent_id", 0)
 	status := utils.QueryInt(ctx, "status", -99)
 	name := ctx.Query("name")
 	user := utils.GetUser(ctx)
 	where := map[string]interface{}{"status": status}
 	if user.Role != constants.USER_ROLE_ADMIN {
 		where["creator"] = user.ID
 	}
 	if groupID != 0 {
 		where["group_id"] = groupID
 	}
 	if agentID != 0 {
 		where["agent_id"] = agentID
 	}
 	if name != "" {
 		where["name"] = name
 	}
 	list, total := providers.JobService.GetJobPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.JobFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 func (c *APIJobController) Show(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	job := providers.JobService.GetJobByID(id)
 	if job == nil {
 		utils.APIBadRequest(ctx, "计划任务不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && job.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此计划任务的权限")
 		return
 	}
 	utils.APIData(ctx, utils.JobFormat(job))
 }

 func (c *APIJobController) Create(ctx *gin.Context) {
 	var req apiJobReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	job := new(models.Job)
 	job.GroupID = req.GroupID
 	job.AgentID = req.AgentID
 	job.Name = req.Name
 	job.Dir = req.Dir
 	job.Program = req.Program
 	job.Args = req.Args
 	job.StdOut = req.StdOut
 	job.StdErr = req.StdErr
 	job.Spec = req.Spec
 	job.Timeout = req.Timeout
 	job.IsMonitor = int64(req.IsMonitor)
 	job.Status = constants.JOB_STATUS_PAUSE
 	job.Creator = utils.GetUserID(ctx)
 	if !providers.JobService.CreateJob(job) {
 		utils.APIError(ctx, "创建计划任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_CREATE)
 	utils.APIData(ctx, gin.H{"id": job.ID})
 }

 func (c *APIJobController) Update(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	job := providers.JobService.GetJobByID(id)
 	if job == nil {
 		utils.APIBadRequest(ctx, "计划任务不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && job.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此计划任务的权限")
 		return
 	}
 	var req apiJobReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	job.GroupID = req.GroupID
 	job.Name = req.Name
 	job.Dir = req.Dir
 	job.Program = req.Program
 	job.Args = req.Args
 	job.StdOut = req.StdOut
 	job.StdErr = req.StdErr
 	job.Spec = req.Spec
 	job.Timeout = req.Timeout
 	job.IsMonitor = int64(req.IsMonitor)
 	job.Updator = utils.GetUserID(ctx)
 	if req.AgentID != 0 {
 		job.AgentID = req.AgentID
 	}
 	if !providers.JobService.UpdateJob(job) {
 		utils.APIError(ctx, "更新计划任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_UPDATE)
 	utils.APIOK(ctx)
 }

 func (c *APIJobController) Copy(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	job := providers.JobService.GetJobByID(id)
 	if job == nil {
 		utils.APIBadRequest(ctx, "计划任务不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && job.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此计划任务的权限")
 		return
 	}
 	_j := new(models.Job)
 	_j.GroupID = job.GroupID
 	_j.Name = job.Name + "_copy"
 	_j.AgentID = job.AgentID
 	_j.Dir = job.Dir
 	_j.Program = job.Program
 	_j.Args = job.Args
 	_j.StdOut = job.StdOut
 	_j.StdErr = job.StdErr
 	_j.Spec = job.Spec
 	_j.Timeout = job.Timeout
 	_j.IsMonitor = job.IsMonitor
 	_j.Status = constants.JOB_STATUS_PAUSE
 	_j.Creator = utils.GetUserID(ctx)
 	if !providers.JobService.CreateJob(_j) {
 		utils.APIError(ctx, "复制计划任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_COPY)
 	utils.APIData(ctx, gin.H{"id": _j.ID})
 }

 func loadJobAgentForUser(ctx *gin.Context, id int64) (*models.Job, *models.Agent, bool) {
 	job := providers.JobService.GetJobByID(id)
 	if job == nil {
 		utils.APIBadRequest(ctx, "计划任务不存在")
 		return nil, nil, false
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && job.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此计划任务的权限")
 		return nil, nil, false
 	}
 	agent := providers.AgentService.GetAgentByID(job.AgentID)
 	if agent == nil {
 		utils.APIBadRequest(ctx, "实例获取失败")
 		return nil, nil, false
 	}
 	if agent.Status != constants.AGENT_ONLINE {
 		utils.APIBadRequest(ctx, "实例不处于运行状态")
 		return nil, nil, false
 	}
 	return job, agent, true
 }

 func jobStart(agent *models.Agent, job *models.Job) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetJob(job.ID)
 	if err != nil {
 		return fmt.Errorf("获取计划任务情况异常:%s", err.Error())
 	}
 	if exist == nil {
 		if err := client.AddJob(job); err != nil {
 			return fmt.Errorf("添加计划任务异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func jobRestart(agent *models.Agent, job *models.Job) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetJob(job.ID)
 	if err != nil {
 		return fmt.Errorf("获取计划任务情况异常:%s", err.Error())
 	}
 	if exist == nil {
 		if err := client.AddJob(job); err != nil {
 			return fmt.Errorf("重启计划任务异常:%s", err.Error())
 		}
 	} else {
 		if err := client.UpdateJob(job); err != nil {
 			return fmt.Errorf("重启计划任务异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func jobStop(agent *models.Agent, jobID int64) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetJob(jobID)
 	if err != nil {
 		return fmt.Errorf("获取计划任务情况异常:%s", err.Error())
 	}
 	if exist != nil {
 		if err := client.RemoveJob(jobID); err != nil {
 			return fmt.Errorf("停止计划任务异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func (c *APIJobController) Start(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	job, agent, ok := loadJobAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if job.Status == constants.JOB_STATUS_RUNNING {
 		utils.APIError(ctx, "计划任务已经启动")
 		return
 	}
 	if err := jobStart(agent, job); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_RUNNING, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "更新计划任务状态失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_START)
 	utils.APIOK(ctx)
 }

 func (c *APIJobController) ReStart(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	job, agent, ok := loadJobAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if err := jobRestart(agent, job); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_RESTART)
 	utils.APIOK(ctx)
 }

 func (c *APIJobController) Pause(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	job, agent, ok := loadJobAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if err := jobStop(agent, job.ID); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_PAUSE, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "更新计划任务状态失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_PAUSE)
 	utils.APIOK(ctx)
 }

 func (c *APIJobController) Delete(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	job, agent, ok := loadJobAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if job.Status == constants.JOB_STATUS_RUNNING {
 		utils.APIError(ctx, "计划任务启动状态不能删除")
 		return
 	}
 	if err := jobStop(agent, job.ID); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_DELETED, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "删除计划任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_JOB, job.ID, constants.ACTION_DELETE)
 	utils.APIOK(ctx)
 }

 func (c *APIJobController) BatchAction(ctx *gin.Context, action string) {
 	var req apiBatchReq
 	if err := ctx.ShouldBind(&req); err != nil || len(req.IDs) == 0 {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	user := utils.GetUser(ctx)
 	for _, id := range req.IDs {
 		job := providers.JobService.GetJobByID(id)
 		if job == nil {
 			continue
 		}
 		if user.Role != constants.USER_ROLE_ADMIN && job.Creator != user.ID {
 			continue
 		}
 		agent := providers.AgentService.GetAgentByID(job.AgentID)
 		if agent == nil || agent.Status != constants.AGENT_ONLINE {
 			continue
 		}
 		switch action {
 		case "start":
 			if job.Status == constants.JOB_STATUS_RUNNING {
 				continue
 			}
 			if err := jobStart(agent, job); err != nil {
 				logger.Errorf("Job BatchStart err:%s", err.Error())
 				continue
 			}
 			providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_RUNNING, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_JOB, job.ID, constants.ACTION_START)
 		case "restart":
 			if err := jobRestart(agent, job); err != nil {
 				logger.Errorf("Job BatchRestart err:%s", err.Error())
 				continue
 			}
 			utils.OpetationLog(user.ID, constants.TYPE_JOB, job.ID, constants.ACTION_RESTART)
 		case "pause":
 			if err := jobStop(agent, job.ID); err != nil {
 				logger.Errorf("Job BatchPause err:%s", err.Error())
 				continue
 			}
 			providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_PAUSE, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_JOB, job.ID, constants.ACTION_PAUSE)
 		case "delete":
 			if job.Status == constants.JOB_STATUS_RUNNING {
 				continue
 			}
 			if err := jobStop(agent, job.ID); err != nil {
 				logger.Errorf("Job BatchDelete err:%s", err.Error())
 				continue
 			}
 			providers.JobService.ChangeJobStatus(job, constants.JOB_STATUS_DELETED, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_JOB, job.ID, constants.ACTION_DELETE)
 		}
 	}
 	utils.APIOK(ctx)
 }

 func (c *APIJobController) BatchStart(ctx *gin.Context)   { c.BatchAction(ctx, "start") }
 func (c *APIJobController) BatchReStart(ctx *gin.Context) { c.BatchAction(ctx, "restart") }
 func (c *APIJobController) BatchPause(ctx *gin.Context)   { c.BatchAction(ctx, "pause") }
 func (c *APIJobController) BatchDelete(ctx *gin.Context) { c.BatchAction(ctx, "delete") }

 // 防止 time 引入告警
 var _ = time.Now
