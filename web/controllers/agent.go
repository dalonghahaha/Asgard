package controllers

import (
	"github.com/gin-gonic/gin"

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

func (c *AgentController) List(ctx *gin.Context) {
	page := DefaultInt(ctx, "page", 1)
	agentList := c.agentService.GetAllAgent()
	total := len(agentList)
	mpurl := "/agent/list"
	ctx.HTML(200, "agent_list", gin.H{
		"Subtitle":   "实例列表",
		"List":       agentList,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
}
