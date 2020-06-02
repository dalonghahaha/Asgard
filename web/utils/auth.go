package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"Asgard/models"
)

func GetUser(ctx *gin.Context) *models.User {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.Redirect(http.StatusFound, "/error")
		ctx.Abort()
		return nil
	}
	_user, ok := user.(*models.User)
	if !ok {
		ctx.Redirect(http.StatusFound, "/error")
		ctx.Abort()
		return nil
	}
	return _user
}

func GetUserID(ctx *gin.Context) int64 {
	user := GetUser(ctx)
	if user == nil {
		return 0
	}
	return user.ID
}
