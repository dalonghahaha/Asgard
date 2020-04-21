package utils

import (
	"strconv"

	"github.com/dalonghahaha/avenger/tools/coding"
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/models"
)

func GetUserID(ctx *gin.Context) int64 {
	token, err := ctx.Cookie("token")
	if err != nil {
		return 0
	}
	id, err := coding.DesDecrypt(token, constants.CookieSalt)
	if err != nil {
		return 0
	}
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}
	return userID
}

func GetUser(ctx *gin.Context) *models.User {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.Redirect(constants.StatusFound, "/error")
		ctx.Abort()
		return nil
	}
	_user, ok := user.(*models.User)
	if !ok {
		ctx.Redirect(constants.StatusFound, "/error")
		ctx.Abort()
		return nil
	}
	return _user
}
