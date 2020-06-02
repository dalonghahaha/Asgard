package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/providers"
	"Asgard/web/utils"
)

func GroupInit(ctx *gin.Context) {
	id, ok := utils.GetID(ctx)
	if !ok {
		ctx.Abort()
		return
	}
	group := providers.GroupService.GetGroupByID(id)
	if group == nil {
		utils.Warning(ctx, "分组不存在")
		ctx.Abort()
		return
	}
	ctx.Set("group", group)
	ctx.Next()
}
