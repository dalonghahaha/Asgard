 package middlewares

 import (
 	"net/http"

 	"github.com/gin-gonic/gin"
 )

 // CORS 简单 CORS 中间件：允许任意来源（开发期足够；生产应由 Nginx 收敛）。
 // 同时处理 OPTIONS 预检。
 func CORS(ctx *gin.Context) {
 	origin := ctx.GetHeader("Origin")
 	if origin == "" {
 		origin = "*"
 	}
 	ctx.Header("Access-Control-Allow-Origin", origin)
 	ctx.Header("Access-Control-Allow-Credentials", "true")
 	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
 	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
 	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Content-Disposition")
 	ctx.Header("Access-Control-Max-Age", "86400")
 	if ctx.Request.Method == http.MethodOptions {
 		ctx.AbortWithStatus(http.StatusNoContent)
 		return
 	}
 	ctx.Next()
 }
