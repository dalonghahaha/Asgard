package middlewares

import (
	"net/http"
	"strconv"

	"github.com/dalonghahaha/avenger/tools/coding"
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
	"Asgard/web/utils"
)

func Login(ctx *gin.Context) {
	token, err := ctx.Cookie("token")
	if err != nil && err != http.ErrNoCookie {
		ctx.Redirect(constants.StatusFound, "/error")
		ctx.Abort()
		return
	}
	if token == "" || err == http.ErrNoCookie {
		ctx.Redirect(constants.StatusFound, "/no_login")
		ctx.Abort()
		return
	}
	userID, err := coding.DesDecrypt(token, constants.WEB_COOKIE_SALT)
	if err != nil {
		utils.CleanTokenCookie(ctx)
		ctx.Redirect(constants.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	_userID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		utils.CleanTokenCookie(ctx)
		ctx.Redirect(constants.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	user := providers.UserService.GetUserByID(_userID)
	if user == nil {
		utils.CleanTokenCookie(ctx)
		ctx.Redirect(constants.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	if user.Status == constants.USER_STATUS_FORBIDDEN {
		utils.CleanTokenCookie(ctx)
		ctx.Redirect(constants.StatusForbidden, "/forbidden")
		ctx.Abort()
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}
