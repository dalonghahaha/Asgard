package utils

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
)

func SetTokenCookie(ctx *gin.Context, value string) {
	ctx.SetCookie("token", value, 7200, "/", constants.WEB_DOMAIN, false, true)
}

func CleanTokenCookie(ctx *gin.Context) {
	ctx.SetCookie("token", "", 0, "/", constants.WEB_DOMAIN, false, true)
}
