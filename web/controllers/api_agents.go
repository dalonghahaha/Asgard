 package controllers

 import (
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 type APIAgentController struct{}

 func NewAPIAgentController() *APIAgentController {
 	return &APIAgentController{}
 }

 func (c *APIAgentController) List(ctx *gin.Context) {
 	page := utils.QueryInt(ctx, "page", 1)
 	status := utils.QueryInt(ctx, "status", -99)
 	alias := ctx.Query("alias")
 	where := map[string]interface{}{"status": status}
 	if alias != "" {
 		where["alias"] = alias
 	}
 	list, total := providers.AgentService.GetAgentPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		a := list[i]
 		out = append(out, gin.H{
 			"id":         a.ID,
 			"alias":      a.Alias,
 			"ip":         a.IP,
 			"port":       a.Port,
 			"master":     a.Master,
 			"status":     a.Status,
 			"created_at": a.CreatedAt,
 		})
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 func (c *APIAgentController) Show(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	agent := providers.AgentService.GetAgentByID(id)
 	if agent == nil {
 		utils.APIBadRequest(ctx, "实例不存在")
 		return
 	}
 	utils.APIData(ctx, gin.H{
 		"id":     agent.ID,
 		"alias":  agent.Alias,
 		"ip":     agent.IP,
 		"port":   agent.Port,
 		"master": agent.Master,
 		"status": agent.Status,
 	})
 }

 type apiAgentUpdateReq struct {
 	Alias string `form:"alias" json:"alias" binding:"required"`
 }

 func (c *APIAgentController) Update(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	agent := providers.AgentService.GetAgentByID(id)
 	if agent == nil {
 		utils.APIBadRequest(ctx, "实例不存在")
 		return
 	}
 	var req apiAgentUpdateReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	agent.Alias = req.Alias
 	if !providers.AgentService.UpdateAgent(agent) {
 		utils.APIError(ctx, "实例更新失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_AGENT, agent.ID, constants.ACTION_UPDATE)
 	utils.APIOK(ctx)
 }

 func (c *APIAgentController) Forbidden(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	agent := providers.AgentService.GetAgentByID(id)
 	if agent == nil {
 		utils.APIBadRequest(ctx, "实例不存在")
 		return
 	}
 	agent.Status = constants.AGENT_FORBIDDEN
 	if !providers.AgentService.UpdateAgent(agent) {
 		utils.APIError(ctx, "实例更新失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_AGENT, agent.ID, constants.ACTION_DELETE)
 	// 级联停用 agent 上所有 app/job/timing
 	apps := providers.AppService.GetAppByAgentID(agent.ID)
 	for i := range apps {
 		providers.AppService.ChangeAPPStatus(&apps[i], constants.APP_STATUS_DELETED, utils.GetUserID(ctx))
 	}
 	jobs := providers.JobService.GetJobByAgentID(agent.ID)
 	for i := range jobs {
 		providers.JobService.ChangeJobStatus(&jobs[i], constants.JOB_STATUS_DELETED, utils.GetUserID(ctx))
 	}
 	timings := providers.TimingService.GetTimingByAgentID(agent.ID)
 	for i := range timings {
 		providers.TimingService.ChangeTimingStatus(&timings[i], constants.TIMING_STATUS_DELETED, utils.GetUserID(ctx))
 	}
 	utils.APIOK(ctx)
 }
