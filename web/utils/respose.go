 package utils

 import (
 	"net/http"

 	"github.com/gin-gonic/gin"
 )

 // 前后端分离后，web 层只剩 JSON 响应工具；HTML 渲染（Render/JumpWarning 等）已下线。

 // APIOK 简单成功响应。
 func APIOK(ctx *gin.Context) {
 	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK})
 }

 // APIData 带 data 的成功响应。
 func APIData(ctx *gin.Context, data interface{}) {
 	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data})
 }

 // APIPage 分页数据响应：{code, message, data:{list, total, page, page_size}}
 func APIPage(ctx *gin.Context, list interface{}, total int, page int, pageSize int) {
 	ctx.JSON(http.StatusOK, gin.H{
 		"code":    http.StatusOK,
 		"message": "ok",
 		"data": gin.H{
 			"list":      list,
 			"total":     total,
 			"page":      page,
 			"page_size": pageSize,
 		},
 	})
 }

 // APIBadRequest 业务参数错误。
 func APIBadRequest(ctx *gin.Context, message string) {
 	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusBadRequest, "message": message})
 }

 // APIError 业务异常。
 func APIError(ctx *gin.Context, message string) {
 	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusInternalServerError, "message": message})
 }
