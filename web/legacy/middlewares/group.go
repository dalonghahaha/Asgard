//go:build ignore
// +build ignore

// DEPRECATED: 仅作历史参考；前后端分离后已下线，详见 doc/TASKS.md Phase 5。
// 原 import 路径已变；如需恢复，请同步修改 import 路径与 package 名。

package legacy

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
