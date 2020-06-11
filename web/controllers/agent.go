package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
	"Asgard/web/utils"
)

type AgentController struct{}

func NewAgentController() *AgentController {
	return &AgentController{}
}

func (c *AgentController) List(ctx *gin.Context) {
	user := utils.GetUser(ctx)
	page := utils.DefaultInt(ctx, "page", 1)
	status := utils.DefaultInt(ctx, "status", -99)
	alias := ctx.Query("alias")
	where := map[string]interface{}{
		"status": status,
	}
	querys := []string{}
	if alias != "" {
		where["alias"] = alias
		querys = append(querys, "alias="+alias)
	}
	if status != -99 {
		querys = append(querys, "status="+strconv.Itoa(status))
	}
	agentList, total := providers.AgentService.GetAgentPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
	mpurl := "/agent/list"
	if len(querys) > 0 {
		mpurl = "/agent/list?" + strings.Join(querys, "&")
	}
	utils.Render(ctx, "agent/list", gin.H{
		"Subtitle":   "实例列表",
		"List":       agentList,
		"Total":      total,
		"StatusList": constants.AGENT_STATUS,
		"Alias":      alias,
		"Status":     status,
		"Role":       user.Role,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *AgentController) Edit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.JumpWarning(ctx, "id参数异常")
		return
	}
	agent := providers.AgentService.GetAgentByID(id)
	if agent == nil {
		utils.JumpWarning(ctx, "获取实例信息异常")
		return
	}
	utils.Render(ctx, "agent/edit", gin.H{
		"Subtitle": "编辑别名",
		"Agent":    agent,
		"BackUrl":  utils.GetReferer(ctx),
	})
}

func (c *AgentController) Update(ctx *gin.Context) {
	alias := ctx.PostForm("alias")
	if alias == "" {
		utils.APIBadRequest(ctx, "别名不能为空")
		return
	}
	agent := utils.GetAgent(ctx)
	agent.Alias = alias
	ok := providers.AgentService.UpdateAgent(agent)
	if !ok {
		utils.APIError(ctx, "实例更新失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_AGENT, agent.ID, constants.ACTION_UPDATE)
	utils.APIOK(ctx)
}

func (c *AgentController) Forbidden(ctx *gin.Context) {
	agent := utils.GetAgent(ctx)
	agent.Status = constants.AGENT_FORBIDDEN
	ok := providers.AgentService.UpdateAgent(agent)
	if !ok {
		utils.APIError(ctx, "实例更新失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_AGENT, agent.ID, constants.ACTION_DELETE)
	apps := providers.AppService.GetAppByAgentID(agent.ID)
	for _, app := range apps {
		providers.AppService.ChangeAPPStatus(&app, constants.APP_STATUS_DELETED, utils.GetUserID(ctx))
	}
	jobs := providers.JobService.GetJobByAgentID(agent.ID)
	for _, job := range jobs {
		providers.JobService.ChangeJobStatus(&job, constants.JOB_STATUS_DELETED, utils.GetUserID(ctx))
	}
	timings := providers.TimingService.GetTimingByAgentID(agent.ID)
	for _, timing := range timings {
		providers.TimingService.ChangeTimingStatus(&timing, constants.TIMING_STATUS_DELETED, utils.GetUserID(ctx))
	}
	utils.APIOK(ctx)
}
