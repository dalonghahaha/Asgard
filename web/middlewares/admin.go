package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/web/utils"
)

func Admin(ctx *gin.Context) {
	user := utils.GetUser(ctx)
	if user == nil {
		ctx.Redirect(constants.StatusFound, "/error")
		ctx.Abort()
		return
	}
	if user.Role != constants.USER_ROLE_ADMIN {
		ctx.Redirect(constants.StatusFound, "/admin_only")
		ctx.Abort()
	}
	ctx.Next()
}
