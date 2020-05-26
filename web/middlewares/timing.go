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
	id := utils.DefaultInt64(ctx, "id", 0)
	_id := utils.FormDefaultInt64(ctx, "id", 0)
	if id == 0 && _id == 0 {
		utils.JumpWarning(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	if id == 0 {
		id = _id
	}
	timing := providers.TimingService.GetTimingByID(id)
	if timing == nil {
		utils.JumpWarning(ctx, "定时任务不存在")
		ctx.Abort()
		return
	}
	ctx.Set("timing", timing)
	ctx.Next()
}

func TimingAgentInit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	_id := utils.FormDefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.APIBadRequest(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	if id == 0 {
		id = _id
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

func BatchTimingAgentInit(ctx *gin.Context) {
	ids := ctx.PostForm("ids")
	if ids == "" {
		utils.APIBadRequest(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	timingAgents := map[*models.Timing]*models.Agent{}
	idList := strings.Split(ids, ",")
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
