package middlewares

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/web/utils"
)

func TimingInit(ctx *gin.Context) {
	id, ok := utils.GetID(ctx)
	if !ok {
		ctx.Abort()
		return
	}
	timing := providers.TimingService.GetTimingByID(id)
	if timing == nil {
		utils.Warning(ctx, "定时任务不存在")
		ctx.Abort()
		return
	}
	user := utils.GetUser(ctx)
	if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
		utils.Warning(ctx, "对不起，您没有操作此应用的权限")
		ctx.Abort()
		return
	}
	ctx.Set("timing", timing)
	ctx.Next()
}

func TimingAgentInit(ctx *gin.Context) {
	id, ok := utils.GetID(ctx)
	if !ok {
		ctx.Abort()
		return
	}
	timing := providers.TimingService.GetTimingByID(id)
	if timing == nil {
		utils.Warning(ctx, "定时任务不存在")
		ctx.Abort()
		return
	}
	user := utils.GetUser(ctx)
	if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
		utils.Warning(ctx, "对不起，您没有操作此应用的权限")
		ctx.Abort()
		return
	}
	ctx.Set("timing", timing)
	agent := providers.AgentService.GetAgentByID(timing.AgentID)
	if agent == nil {
		utils.Warning(ctx, "实例获取失败")
		ctx.Abort()
		return
	}
	if agent.Status != constants.AGENT_ONLINE {
		utils.Warning(ctx, "实例不处于运行状态")
		ctx.Abort()
		return
	}
	ctx.Set("agent", agent)
	ctx.Next()
}

func BatchTimingAgentInit(ctx *gin.Context) {
	ids := ctx.PostForm("ids")
	if ids == "" {
		utils.APIBadRequest(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	timingAgents := map[*models.Timing]*models.Agent{}
	idList := strings.Split(ids, ",")
	user := utils.GetUser(ctx)
	for _, id := range idList {
		_id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.APIBadRequest(ctx, "请求参数异常")
			ctx.Abort()
			return
		}
		timing := providers.TimingService.GetTimingByID(_id)
		if timing == nil {
			utils.APIError(ctx, "应用不存在")
			ctx.Abort()
			return
		}
		if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
			utils.APIError(ctx, "对不起，您没有操作此应用的权限")
			ctx.Abort()
			return
		}
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
		timingAgents[timing] = agent
	}
	ctx.Set("timing_agent", timingAgents)
	ctx.Next()
}
