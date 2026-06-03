//go:build ignore
// +build ignore

// DEPRECATED: 仅作历史参考；前后端分离后已下线，详见 doc/TASKS.md Phase 5。
// 原 import 路径已变；如需恢复，请同步修改 import 路径与 package 名。

package legacy

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/web/utils"
)

func Admin(ctx *gin.Context) {
	user := utils.GetUser(ctx)
	if user == nil {
		if ctx.Request.Method == "GET" {
			ctx.Redirect(http.StatusFound, "/error")
		} else if ctx.Request.Method == "POST" {
			utils.APIBadRequest(ctx, "非法请求")
		}
		ctx.Abort()
		return
	}
	if user.Role != constants.USER_ROLE_ADMIN {
		if ctx.Request.Method == "GET" {
			ctx.Redirect(http.StatusFound, "/admin_only")
		} else if ctx.Request.Method == "POST" {
			utils.APIBadRequest(ctx, "只有管理员才能进行此操作")
		}
		ctx.Abort()
		return
	}
	ctx.Next()
}
