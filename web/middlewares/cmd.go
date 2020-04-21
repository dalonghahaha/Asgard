package middlewares

import (
	"github.com/gin-gonic/gin"

	"Asgard/web/utils"
)

func CmdConfigVerify(ctx *gin.Context) {
	agentID := utils.FormDefaultInt64(ctx, "agent_id", 0)
	name := ctx.PostForm("name")
	dir := ctx.PostForm("dir")
	program := ctx.PostForm("program")
	stdOut := ctx.PostForm("std_out")
	stdErr := ctx.PostForm("std_err")
	if !utils.Required(ctx, name, "名称不能为空") {
		ctx.Abort()
		return
	}
	if !utils.Required(ctx, dir, "执行目录不能为空") {
		ctx.Abort()
		return
	}
	if !utils.Required(ctx, program, "执行程序不能为空") {
		ctx.Abort()
		return
	}
	if !utils.Required(ctx, stdOut, "标准输出路径不能为空") {
		ctx.Abort()
		return
	}
	if !utils.Required(ctx, stdErr, "错误输出路径不能为空") {
		ctx.Abort()
		return
	}
	if agentID == 0 {
		utils.APIBadRequest(ctx, "运行实例不能为空")
		ctx.Abort()
		return
	}
	ctx.Next()
}
