package middlewares

import (
	"net/http"
	"strconv"

	"github.com/dalonghahaha/avenger/tools/coding"
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
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
	userID, err := coding.DesDecrypt(token, constants.CookieSalt)
	if err != nil {
		ctx.SetCookie("token", "", 0, "/", constants.Domain, false, true)
		ctx.Redirect(constants.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	_userID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ctx.SetCookie("token", "", 0, "/", constants.Domain, false, true)
		ctx.Redirect(constants.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	user := providers.UserService.GetUserByID(_userID)
	if user == nil {
		ctx.SetCookie("token", "", 0, "/", constants.Domain, false, true)
		ctx.Redirect(constants.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	if user.Status == constants.USER_STATUS_FORBIDDEN {
		ctx.SetCookie("token", "", 0, "/", constants.Domain, false, true)
		ctx.Redirect(constants.StatusForbidden, "/forbidden")
		ctx.Abort()
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}
