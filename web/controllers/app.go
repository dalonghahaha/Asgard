package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"Asgard/rpc"
)

type AppController struct {
	appClient rpc.GuardClient
}

func NewAppController() *AppController {
	port := "localhost:27149"
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		panic("Can't connect: " + port)
	}
	client := rpc.NewGuardClient(conn)
	return &AppController{
		appClient: client,
	}
}

func (c *AppController) List(ctx *gin.Context) {
	timeout := time.Second * 30
	_ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := c.appClient.List(_ctx, &rpc.Empty{})
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"list": response.GetApps(),
	})
}

func (c *AppController) Get(ctx *gin.Context) {
	timeout := time.Second * 30
	_ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	name := ctx.Param("name")
	response, err := c.appClient.Get(_ctx, &rpc.AppNameRequest{Name: name})
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(200, response.GetApp())
}

func (c *AppController) Add(ctx *gin.Context) {

}

func (c *AppController) Update(ctx *gin.Context) {

}

func (c *AppController) Remove(ctx *gin.Context) {

}
