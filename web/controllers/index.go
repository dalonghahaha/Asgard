package controllers

import (
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(StatusOK, "index", gin.H{
		"Subtitle": "首页",
	})
}

func UI(c *gin.Context) {
	c.HTML(StatusOK, "UI", gin.H{
		"Subtitle": "布局",
	})
}

func Nologin(c *gin.Context) {
	c.HTML(StatusOK, "error/nologin.html", gin.H{
		"Subtitle": "未登录提示页",
	})
}

func Error(c *gin.Context) {
	c.HTML(StatusOK, "error/err.html", gin.H{
		"Subtitle": "服务器异常",
	})
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
