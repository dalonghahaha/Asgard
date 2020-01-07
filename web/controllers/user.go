package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/services"
)

type UserController struct {
	useService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		useService: services.NewUserService(),
	}
}

func (c *UserController) Register(ctx *gin.Context) {
	ctx.HTML(StatusOK, "register.html", gin.H{
		"Subtitle": "用户注册",
	})
}

func (c *UserController) DoRegister(ctx *gin.Context) {
	ctx.JSON(StatusOK, gin.H{
		"code": 200,
	})
}
