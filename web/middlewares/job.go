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

func JobInit(ctx *gin.Context) {
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
	job := providers.JobService.GetJobByID(id)
	if job == nil {
		utils.APIError(ctx, "计划任务不存在")
		ctx.Abort()
		return
	}
	ctx.Set("job", job)
	ctx.Next()
}

func JobAgentInit(ctx *gin.Context) {
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

func BatchJobAgentInit(ctx *gin.Context) {
	ids := ctx.PostForm("ids")
	if ids == "" {
		utils.APIBadRequest(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	jobAgents := map[*models.Job]*models.Agent{}
	idList := strings.Split(ids, ",")
	for _, id := range idList {
		_id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.APIBadRequest(ctx, "请求参数异常")
			ctx.Abort()
			return
		}
		job := providers.JobService.GetJobByID(_id)
		if job == nil {
			utils.APIError(ctx, "应用不存在")
			ctx.Abort()
			return
		}
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
		jobAgents[job] = agent
	}
	ctx.Set("job_agent", jobAgents)
	ctx.Next()
}
