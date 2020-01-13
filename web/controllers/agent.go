package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/models"
	"Asgard/services"
)

type AgentController struct {
	agentService *services.AgentService
}

func NewAgentController() *AgentController {
	return &AgentController{
		agentService: services.NewAgentService(),
	}
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
	page := DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	agentList, total := c.agentService.GetAgentPageList(where, page, PageSize)
	mpurl := "/agent/list"
	ctx.HTML(200, "agent/list", gin.H{
		"Subtitle":   "实例列表",
		"List":       agentList,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
}
