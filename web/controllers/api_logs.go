 package controllers

 import (
 	"fmt"
 	"strings"

 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 // APILogController 提供 out_log/err_log 的内容查询与分页拉取。
 // 数据来源是 agent 端的本地日志文件，通过 gRPC GetLog 拉回。
 type APILogController struct{}

 func NewAPILogController() *APILogController {
 	return &APILogController{}
 }

 func (c *APILogController) appLogData(ctx *gin.Context, errLog bool) {
 	id := utils.QueryInt64(ctx, "app_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "app_id 不能为空")
 		return
 	}
 	app := providers.AppService.GetAppByID(id)
 	if app == nil {
 		utils.APIBadRequest(ctx, "应用不存在")
 		return
 	}
 	if !checkCreator(ctx, app.Creator) {
 		return
 	}
 	lines := utils.QueryInt64(ctx, "lines", constants.WEB_LOG_SIZE)
 	agent := providers.AgentService.GetAgentByID(app.AgentID)
 	if agent == nil || agent.Status != constants.AGENT_ONLINE {
 		utils.APIBadRequest(ctx, "实例不在线")
 		return
 	}
 	path := app.StdOut
 	if errLog {
 		path = app.StdErr
 	}
 	content, err := fetchLog(agent, path, lines)
 	if err != nil {
 		utils.APIError(ctx, "获取日志失败:"+err.Error())
 		return
 	}
 	utils.APIData(ctx, gin.H{
 		"app_id":  app.ID,
 		"path":    path,
 		"content": content,
 	})
 }

 func (c *APILogController) jobLogData(ctx *gin.Context, errLog bool) {
 	id := utils.QueryInt64(ctx, "job_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "job_id 不能为空")
 		return
 	}
 	job := providers.JobService.GetJobByID(id)
 	if job == nil {
 		utils.APIBadRequest(ctx, "计划任务不存在")
 		return
 	}
 	if !checkCreator(ctx, job.Creator) {
 		return
 	}
 	lines := utils.QueryInt64(ctx, "lines", constants.WEB_LOG_SIZE)
 	agent := providers.AgentService.GetAgentByID(job.AgentID)
 	if agent == nil || agent.Status != constants.AGENT_ONLINE {
 		utils.APIBadRequest(ctx, "实例不在线")
 		return
 	}
 	path := job.StdOut
 	if errLog {
 		path = job.StdErr
 	}
 	content, err := fetchLog(agent, path, lines)
 	if err != nil {
 		utils.APIError(ctx, "获取日志失败:"+err.Error())
 		return
 	}
 	utils.APIData(ctx, gin.H{
 		"job_id":  job.ID,
 		"path":    path,
 		"content": content,
 	})
 }

 func (c *APILogController) timingLogData(ctx *gin.Context, errLog bool) {
 	id := utils.QueryInt64(ctx, "timing_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "timing_id 不能为空")
 		return
 	}
 	timing := providers.TimingService.GetTimingByID(id)
 	if timing == nil {
 		utils.APIBadRequest(ctx, "定时任务不存在")
 		return
 	}
 	if !checkCreator(ctx, timing.Creator) {
 		return
 	}
 	lines := utils.QueryInt64(ctx, "lines", constants.WEB_LOG_SIZE)
 	agent := providers.AgentService.GetAgentByID(timing.AgentID)
 	if agent == nil || agent.Status != constants.AGENT_ONLINE {
 		utils.APIBadRequest(ctx, "实例不在线")
 		return
 	}
 	path := timing.StdOut
 	if errLog {
 		path = timing.StdErr
 	}
 	content, err := fetchLog(agent, path, lines)
 	if err != nil {
 		utils.APIError(ctx, "获取日志失败:"+err.Error())
 		return
 	}
 	utils.APIData(ctx, gin.H{
 		"timing_id": timing.ID,
 		"path":      path,
 		"content":   content,
 	})
 }

 func (c *APILogController) AppOutLogData(ctx *gin.Context)  { c.appLogData(ctx, false) }
 func (c *APILogController) AppErrLogData(ctx *gin.Context)  { c.appLogData(ctx, true) }
 func (c *APILogController) JobOutLogData(ctx *gin.Context)  { c.jobLogData(ctx, false) }
 func (c *APILogController) JobErrLogData(ctx *gin.Context)  { c.jobLogData(ctx, true) }
 func (c *APILogController) TimingOutLogData(ctx *gin.Context) {
 	c.timingLogData(ctx, false)
 }
 func (c *APILogController) TimingErrLogData(ctx *gin.Context) {
 	c.timingLogData(ctx, true)
 }

 // —— 共享 helper ——

 func fetchLog(agent *models.Agent, path string, lines int64) ([]string, error) {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return nil, err
 	}
 	return client.GetLog(path, lines)
 }

 func checkCreator(ctx *gin.Context, creator int64) bool {
 	user := utils.GetUser(ctx)
 	if user == nil {
 		utils.APIBadRequest(ctx, "未登录")
 		return false
 	}
 	if user.Role != constants.USER_ROLE_ADMIN && creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此资源的权限")
 		return false
 	}
 	return true
 }

 // 避免 fmt/strings 引入告警（如果未来扩展）
 var _ = fmt.Sprintf
 var _ = strings.Join
