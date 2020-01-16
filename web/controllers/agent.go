package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/models"
	"Asgard/services"
)

type AgentController struct {
	agentService   *services.AgentService
	moniterService *services.MonitorService
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

func (c *AgentController) Monitor(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		JumpError(ctx)
		return
	}
	cpus := []string{}
	memorys := []string{}
	times := []string{}
	moniters := c.moniterService.GetAgentMonitor(id, 100)
	for _, moniter := range moniters {
		cpus = append(cpus, FormatFloat(moniter.CPU))
		memorys = append(memorys, FormatFloat(moniter.Memory))
		times = append(times, FormatTime(moniter.CreatedAt))
	}
	ctx.HTML(StatusOK, "agent/monitor", gin.H{
		"Subtitle": "监控信息",
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}
