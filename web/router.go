 package web

 import (
 	"github.com/gin-gonic/gin"

 	"Asgard/web/routers"
 )

 // 前后端分离后，web 层只挂 /api/v1 子路由；HTML 路由已下线（移至 web/legacy/）。
 func setupRouter() {
 	api := server.Group("/api/v1")
 	routers.SetupAPIRouter(api)

 	// 兜底：所有未匹配的 GET 请求返回 JSON 404（前端 SPA 自己处理路由）
 	server.NoRoute(func(ctx *gin.Context) {
 		if ctx.Request.Method == "GET" {
 			ctx.JSON(404, gin.H{
 				"code":    404,
 				"message": "API 不存在或前端路由由 SPA 处理",
 			})
 			return
 		}
 		ctx.JSON(404, gin.H{"code": 404, "message": "not found"})
 	})
 }
