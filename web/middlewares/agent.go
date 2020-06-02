package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/providers"
	"Asgard/web/utils"
)

func AgentInit(ctx *gin.Context) {
	id, ok := utils.GetID(ctx)
	if !ok {
		ctx.Abort()
		return
	}
	agent := providers.AgentService.GetAgentByID(id)
	if agent == nil {
		utils.Warning(ctx, "实例不存在")
		ctx.Abort()
		return
	}
	ctx.Set("agent", agent)
	ctx.Next()
}
