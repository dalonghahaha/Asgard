package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/providers"
	"Asgard/web/utils"
)

func AgentInit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	_id := utils.FormDefaultInt64(ctx, "id", 0)
	if id == 0 && _id == 0 {
		utils.APIError(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	if id == 0 {
		id = _id
	}
	agent := providers.AgentService.GetAgentByID(id)
	if agent == nil {
		utils.APIError(ctx, "实例不存在")
		ctx.Abort()
		return
	}
	ctx.Set("agent", agent)
	ctx.Next()
}
