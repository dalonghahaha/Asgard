package middlewares

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
