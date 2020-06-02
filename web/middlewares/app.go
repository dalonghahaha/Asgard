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

func AppInit(ctx *gin.Context) {
	id, ok := utils.GetID(ctx)
	if !ok {
		ctx.Abort()
		return
	}
	app := providers.AppService.GetAppByID(id)
	if app == nil {
		utils.Warning(ctx, "应用不存在")
		ctx.Abort()
		return
	}
	user := utils.GetUser(ctx)
	if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
		utils.Warning(ctx, "对不起，您没有操作此应用的权限")
		ctx.Abort()
		return
	}
	ctx.Set("app", app)
	ctx.Next()
}

func AppAgentInit(ctx *gin.Context) {
	id, ok := utils.GetID(ctx)
	if !ok {
		ctx.Abort()
		return
	}
	app := providers.AppService.GetAppByID(id)
	if app == nil {
		utils.Warning(ctx, "应用不存在")
		ctx.Abort()
		return
	}
	user := utils.GetUser(ctx)
	if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
		utils.Warning(ctx, "对不起，您没有操作此应用的权限")
		ctx.Abort()
		return
	}
	ctx.Set("app", app)
	agent := providers.AgentService.GetAgentByID(app.AgentID)
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

func BatchAppAgentInit(ctx *gin.Context) {
	ids := ctx.PostForm("ids")
	if ids == "" {
		utils.APIBadRequest(ctx, "请求参数异常")
		ctx.Abort()
		return
	}
	appAgents := map[*models.App]*models.Agent{}
	idList := strings.Split(ids, ",")
	user := utils.GetUser(ctx)
	for _, id := range idList {
		_id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.APIBadRequest(ctx, "请求参数异常")
			ctx.Abort()
			return
		}
		app := providers.AppService.GetAppByID(_id)
		if app == nil {
			utils.APIError(ctx, "应用不存在")
			ctx.Abort()
			return
		}
		if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
			utils.APIError(ctx, "对不起，您没有操作此应用的权限")
			ctx.Abort()
			return
		}
		agent := providers.AgentService.GetAgentByID(app.AgentID)
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
		appAgents[app] = agent
	}
	ctx.Set("app_agent", appAgents)
	ctx.Next()
}
