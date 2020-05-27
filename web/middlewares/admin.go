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
		ctx.Redirect(http.StatusFound, "/error")
		ctx.Abort()
		return
	}
	if user.Role != constants.USER_ROLE_ADMIN {
		ctx.Redirect(http.StatusFound, "/admin_only")
		ctx.Abort()
	}
	ctx.Next()
}
