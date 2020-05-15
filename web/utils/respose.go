package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrorUrl         = "/error"
	StatusOK         = http.StatusOK
	StatusBadRequest = http.StatusBadRequest
	StatusError      = http.StatusInternalServerError
	StatusFound      = http.StatusFound
)

func APIOK(ctx *gin.Context) {
	ctx.JSON(StatusOK, gin.H{"code": StatusOK})
}

func APIData(ctx *gin.Context, data interface{}) {
	ctx.JSON(StatusOK, gin.H{"code": StatusOK, "data": data})
}

func APIBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(StatusOK, gin.H{"code": StatusBadRequest, "message": message})
}

func APIError(ctx *gin.Context, message string) {
	ctx.JSON(StatusOK, gin.H{"code": StatusError, "message": message})
}

func APIErrorByCode(ctx *gin.Context, code int) {
	ctx.JSON(StatusOK, gin.H{"code": code, "message": GetErrorMessage(code)})
}

func JumpWarning(ctx *gin.Context, message string) {
	ctx.HTML(StatusOK, "warning", gin.H{
		"Message": message,
	})
}

func JumpWarningByCode(ctx *gin.Context, code int) {
	ctx.HTML(StatusOK, "warning", gin.H{
		"Message": GetErrorMessage(code),
	})
}

func JumpError(ctx *gin.Context) {
	ctx.Redirect(StatusError, ErrorUrl)
}
