package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/web/utils"
)

type AgentController struct{}

func NewAgentController() *AgentController {
	return &AgentController{}
}

func (c *AgentController) formatAgent(info *models.Agent) map[string]interface{} {
	data := map[string]interface{}{
		"ID":     info.ID,
		"IP":     info.IP,
		"Port":   info.Port,
		"Status": info.Status,
	}
	return data
}

func (c *AgentController) List(ctx *gin.Context) {
	page := utils.DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	agentList, total := providers.AgentService.GetAgentPageList(where, page, PageSize)
	mpurl := "/agent/list"
	ctx.HTML(200, "agent/list", gin.H{
		"Subtitle":   "实例列表",
		"List":       agentList,
		"Total":      total,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *AgentController) Monitor(ctx *gin.Context) {
	id := utils.DefaultInt(ctx, "id", 0)
	if id == 0 {
		utils.JumpError(ctx)
		return
	}
	cpus := []string{}
	memorys := []string{}
	times := []string{}
	moniters := providers.MoniterService.GetAgentMonitor(id, 100)
	for _, moniter := range moniters {
		cpus = append(cpus, utils.FormatFloat(moniter.CPU))
		memorys = append(memorys, utils.FormatFloat(moniter.Memory))
		times = append(times, utils.FormatTime(moniter.CreatedAt))
	}
	ctx.HTML(StatusOK, "monitor/list", gin.H{
		"Subtitle": "实例监控信息",
		"BackUrl":  "/agent/list",
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
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
	ctx.HTML(StatusOK, "agent/edit", gin.H{
		"Subtitle": "编辑别名",
		"Agent":    agent,
		"BackUrl":  GetReferer(ctx),
	})
}

func (c *AgentController) Update(ctx *gin.Context) {
	id := utils.FormDefaultInt64(ctx, "id", 0)
	alias := ctx.PostForm("alias")
	status := ctx.PostForm("status")
	if id == 0 {
		utils.APIBadRequest(ctx, "ID格式错误")
		return
	}
	if alias == "" {
		utils.APIBadRequest(ctx, "别名不能为空")
		return
	}
	agent := providers.AgentService.GetAgentByID(id)
	if agent == nil {
		utils.APIBadRequest(ctx, "实例不存在")
		return
	}
	agent.Alias = alias
	if status != "" {
		agent.Status = constants.AGENT_FORBIDDEN
	}
	ok := providers.AgentService.UpdateAgent(agent)
	if !ok {
		utils.APIError(ctx, "更新失败")
		return
	}
	utils.APIOK(ctx)
}
