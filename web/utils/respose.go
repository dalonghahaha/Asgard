package utils

import (
	"Asgard/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Render(ctx *gin.Context, url string, data gin.H) {
	ctx.HTML(http.StatusOK, url, data)
}

func APIOK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
}

func APIData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
}

func APIBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": message})
}

func APIError(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "message": message})
}

func APIErrorByCode(ctx *gin.Context, code int) {
	ctx.JSON(http.StatusOK, gin.H{"code": code, "message": GetErrorMessage(code)})
}

func JumpWarning(ctx *gin.Context, message string) {
	ctx.HTML(http.StatusOK, "warning", gin.H{"Message": message})
}

func JumpWarningByCode(ctx *gin.Context, code int) {
	ctx.HTML(http.StatusOK, "warning", gin.H{"Message": GetErrorMessage(code)})
}

func JumpError(ctx *gin.Context) {
	ctx.Redirect(http.StatusInternalServerError, constants.WEB_ERROR_URL)
}
