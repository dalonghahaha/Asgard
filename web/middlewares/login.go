package middlewares

import (
	"net/http"
	"strconv"

	"github.com/dalonghahaha/avenger/tools/coding"
	"github.com/gin-gonic/gin"

	"Asgard/services"
	"Asgard/web/controllers"
)

func Login(ctx *gin.Context) {
	token, err := ctx.Cookie("token")
	if err != nil && err != http.ErrNoCookie {
		ctx.Redirect(http.StatusFound, "/error")
		ctx.Abort()
		return
	} else if token == "" || err == http.ErrNoCookie {
		ctx.Redirect(http.StatusFound, "/no_login")
		ctx.Abort()
		return
	}
	userID, err := coding.DesDecrypt(token, controllers.CookieSalt)
	if err != nil {
		ctx.SetCookie("token", "", 0, "/", controllers.Domain, false, true)
		ctx.Redirect(http.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	_userID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ctx.SetCookie("token", "", 0, "/", controllers.Domain, false, true)
		ctx.Redirect(http.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	useService := services.NewUserService()
	user := useService.GetUserByID(_userID)
	if user == nil {
		ctx.SetCookie("token", "", 0, "/", controllers.Domain, false, true)
		ctx.Redirect(http.StatusFound, "/auth_fail")
		ctx.Abort()
		return
	}
	ctx.Next()
}
