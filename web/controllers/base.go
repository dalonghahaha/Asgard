package controllers

import (
	"net/http"
	"strconv"

	"github.com/dalonghahaha/avenger/tools/coding"
	"github.com/gin-gonic/gin"
)

var (
	OutDir       = "runtime/"
	StatusOK     = http.StatusOK
	StatusFound  = http.StatusFound
	PageSize     = 10
	CookieSalt   = "sdswqeqx"
	Domain       = "localhost"
	TimeLocation = "Asia/Shanghai"
	TimeLayout   = "2006-01-02 15:04"
)

func GetUserID(ctx *gin.Context) int64 {
	token, err := ctx.Cookie("token")
	if err != nil {
		return 0
	}
	_token, err := coding.DesDecrypt(token, CookieSalt)
	if err != nil {
		return 0
	}
	userID, err := strconv.Atoi(_token)
	if err != nil {
		return 0
	}
	return int64(userID)
}
