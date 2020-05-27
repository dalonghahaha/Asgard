package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/providers"
	"Asgard/web/utils"
)

func GroupInit(ctx *gin.Context) {
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
	group := providers.GroupService.GetGroupByID(id)
	if group == nil {
		utils.APIError(ctx, "分组不存在")
		ctx.Abort()
		return
	}
	ctx.Set("group", group)
	ctx.Next()
}
