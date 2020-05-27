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
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *AgentController) Edit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.JumpError(ctx)
		return
	}
	agent := providers.AgentService.GetAgentByID(id)
	if agent == nil {
		utils.JumpError(ctx)
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
	status := ctx.PostForm("status")
	if alias == "" {
		utils.APIBadRequest(ctx, "别名不能为空")
		return
	}
	agent := utils.GetAgent(ctx)
	agent.Alias = alias
	if status != "" {
		agent.Status = constants.AGENT_FORBIDDEN
	}
	ok := providers.AgentService.UpdateAgent(agent)
	if !ok {
		utils.APIError(ctx, "实例更新失败")
		return
	}
	utils.APIOK(ctx)
}
