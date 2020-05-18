package controllers

import (
	"net/http"
	"strconv"

	"github.com/dalonghahaha/avenger/tools/coding"
	"github.com/gin-gonic/gin"

	"Asgard/constants"
)

var (
	OutDir       = "runtime/"
	StatusOK     = http.StatusOK
	StatusFound  = http.StatusFound
	PageSize     = 10
	LogSize      = int64(20)
	TimeLocation = "Asia/Shanghai"
	TimeLayout   = "2006-01-02 15:04"
)

func GetReferer(ctx *gin.Context) string {
	return ctx.Request.Header.Get("Referer")
}

func GetUserID(ctx *gin.Context) int64 {
	token, err := ctx.Cookie("token")
	if err != nil {
		return 0
	}
	_token, err := coding.DesDecrypt(token, constants.WEB_COOKIE_SALT)
	if err != nil {
		return 0
	}
	userID, err := strconv.Atoi(_token)
	if err != nil {
		return 0
	}
	return int64(userID)
}
