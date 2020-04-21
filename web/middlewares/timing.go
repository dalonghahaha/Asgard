package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
	"Asgard/web/utils"
)

func TimingInit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	_id := utils.FormDefaultInt64(ctx, "id", 0)
	if id == 0 && _id == 0 {
		utils.JumpWarning(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	if id != 0 {
		timing := providers.TimingService.GetTimingByID(id)
		if timing == nil {
			utils.JumpWarning(ctx, "定时任务不存在")
			ctx.Abort()
			return
		}
		ctx.Set("timing", timing)
	} else {
		timing := providers.TimingService.GetTimingByID(_id)
		if timing == nil {
			utils.JumpWarning(ctx, "定时任务不存在")
			ctx.Abort()
			return
		}
		ctx.Set("timing", timing)
	}
	ctx.Next()
}

func TimingAgentInit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.APIBadRequest(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	timing := providers.TimingService.GetTimingByID(id)
	if timing == nil {
		utils.APIError(ctx, "定时任务不存在")
		ctx.Abort()
		return
	}
	ctx.Set("timing", timing)
	agent := providers.AgentService.GetAgentByID(timing.AgentID)
	if agent == nil {
		utils.APIError(ctx, "实例获取失败")
		ctx.Abort()
		return
	}
	if agent.Status != constants.AGENT_ONLINE {
		utils.APIError(ctx, "实例不处于运行状态")
		ctx.Abort()
		return
	}
	ctx.Set("agent", agent)
	ctx.Next()
}
