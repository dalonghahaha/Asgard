package web

import "github.com/gin-gonic/gin"

// setupAPIRouter 注册 Phase 1 的 JSON API 子路由。
// 当前仅有 /health 占位接口，后续 T-105~T-119 逐任务在此追加。
// 公共中间件（CORS / APIAuth）由 setupAPIRouter 内部按需挂载。
func setupAPIRouter(api *gin.RouterGroup) {
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 200, "message": "ok", "data": gin.H{"status": "up"}})
	})
}
