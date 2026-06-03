//go:build ignore
// +build ignore

// DEPRECATED: 仅作历史参考；前后端分离后已下线，详见 doc/TASKS.md Phase 5。
// 原 import 路径已变；如需恢复，请同步修改 import 路径与 package 名。

package legacy

import (
	"github.com/gin-gonic/gin"

	"Asgard/web/utils"
)

func CmdConfigVerify(ctx *gin.Context) {
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
	ctx.Next()
}
