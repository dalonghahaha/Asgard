package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
	"Asgard/web/utils"
)

func JobInit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	_id := utils.FormDefaultInt64(ctx, "id", 0)
	if id == 0 && _id == 0 {
		utils.JumpWarning(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	if id != 0 {
		job := providers.JobService.GetJobByID(id)
		if job == nil {
			utils.JumpWarning(ctx, "计划任务不存在")
			ctx.Abort()
			return
		}
		ctx.Set("job", job)
	} else {
		job := providers.JobService.GetJobByID(_id)
		if job == nil {
			utils.JumpWarning(ctx, "计划任务不存在")
			ctx.Abort()
			return
		}
		ctx.Set("job", job)
	}
	ctx.Next()
}

func JobAgentInit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.APIBadRequest(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	job := providers.JobService.GetJobByID(id)
	if job == nil {
		utils.APIError(ctx, "计划任务不存在")
		ctx.Abort()
		return
	}
	ctx.Set("job", job)
	agent := providers.AgentService.GetAgentByID(job.AgentID)
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
