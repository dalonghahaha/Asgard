 package routers

import (
	"github.com/gin-gonic/gin"

	"Asgard/web/controllers"
	"Asgard/web/middlewares"
)

 var (
 	apiAuthCtl    = controllers.NewAPIAuthController()
 	apiUserCtl    = controllers.NewAPIUserController()
 	apiAgentCtl   = controllers.NewAPIAgentController()
 	apiGroupCtl   = controllers.NewAPIGroupController()
 	apiAppCtl     = controllers.NewAPIAppController()
 	apiJobCtl     = controllers.NewAPIJobController()
 	apiTimingCtl  = controllers.NewAPITimingController()
 	apiMonitorCtl = controllers.NewAPIMonitorController()
 	apiArchiveCtl = controllers.NewAPIArchiveController()
 	apiLogCtl     = controllers.NewAPILogController()
 	apiSSELogCtl  = controllers.NewAPISSELogController()
 	apiSSEMonCtl  = controllers.NewAPISSEMonitorController()
 	apiExcCtl     = controllers.NewAPIExceptionController()
 	apiOpCtl      = controllers.NewAPIOperationController()
 )

 // SetupAPIRouter 注册 Phase 1 的 JSON API 子路由。
// 公共中间件：CORS（开发期允许任意来源）+ APIAuth（除 /auth/login 与 /health 外都需要鉴权）。
func SetupAPIRouter(api *gin.RouterGroup) {
 	// 全局 CORS（覆盖 /api/v1 下所有子路由）
 	api.Use(middlewares.CORS)

 	// 公共：健康检查 + 登录（不需要鉴权）
 	api.GET("/health", func(c *gin.Context) {
 		c.JSON(200, gin.H{"code": 200, "message": "ok", "data": gin.H{"status": "up"}})
 	})
 	auth := api.Group("/auth")
 	{
 		auth.POST("/login", apiAuthCtl.Login)
 	}

 	// 需要鉴权的子路由
 	authed := api.Group("")
 	authed.Use(middlewares.APIAuth)
 	{
 		// 鉴权
 		authed.GET("/auth/info", apiAuthCtl.Info)
 		authed.POST("/auth/logout", apiAuthCtl.Logout)
 		authed.POST("/auth/change_password", apiAuthCtl.ChangePassword)

 		// 用户
 		users := authed.Group("/users")
 		users.GET("", apiUserCtl.List)
 		users.GET("/:id", apiUserCtl.Show)
 		users.POST("", middlewares.APIAuthAdmin, apiUserCtl.Create)
 		users.PUT("/:id", middlewares.APIAuthAdmin, apiUserCtl.Update)
 		users.POST("/:id/forbidden", middlewares.APIAuthAdmin, apiUserCtl.Forbidden)
 		users.POST("/:id/reset_password", middlewares.APIAuthAdmin, apiUserCtl.ResetPassword)

 		// 实例
 		agents := authed.Group("/agents")
 		agents.GET("", apiAgentCtl.List)
 		agents.GET("/:id", apiAgentCtl.Show)
 		agents.PUT("/:id", middlewares.APIAuthAdmin, apiAgentCtl.Update)
 		agents.POST("/:id/forbidden", middlewares.APIAuthAdmin, apiAgentCtl.Forbidden)

 		// 分组
 		groups := authed.Group("/groups")
 		groups.GET("", apiGroupCtl.List)
 		groups.POST("", apiGroupCtl.Create)
 		groups.PUT("/:id", apiGroupCtl.Update)
 		groups.DELETE("/:id", apiGroupCtl.Delete)

 		// 应用
 		apps := authed.Group("/apps")
 		apps.GET("", apiAppCtl.List)
 		apps.GET("/:id", apiAppCtl.Show)
 		apps.POST("", apiAppCtl.Create)
 		apps.PUT("/:id", apiAppCtl.Update)
 		apps.POST("/:id/copy", apiAppCtl.Copy)
 		apps.POST("/:id/start", apiAppCtl.Start)
 		apps.POST("/:id/restart", apiAppCtl.ReStart)
 		apps.POST("/:id/pause", apiAppCtl.Pause)
 		apps.DELETE("/:id", apiAppCtl.Delete)
 		apps.POST("/batch_start", apiAppCtl.BatchStart)
 		apps.POST("/batch_restart", apiAppCtl.BatchReStart)
 		apps.POST("/batch_pause", apiAppCtl.BatchPause)
 		apps.POST("/batch_delete", apiAppCtl.BatchDelete)

 		// 计划任务
 		jobs := authed.Group("/jobs")
 		jobs.GET("", apiJobCtl.List)
 		jobs.GET("/:id", apiJobCtl.Show)
 		jobs.POST("", apiJobCtl.Create)
 		jobs.PUT("/:id", apiJobCtl.Update)
 		jobs.POST("/:id/copy", apiJobCtl.Copy)
 		jobs.POST("/:id/start", apiJobCtl.Start)
 		jobs.POST("/:id/restart", apiJobCtl.ReStart)
 		jobs.POST("/:id/pause", apiJobCtl.Pause)
 		jobs.DELETE("/:id", apiJobCtl.Delete)
 		jobs.POST("/batch_start", apiJobCtl.BatchStart)
 		jobs.POST("/batch_restart", apiJobCtl.BatchReStart)
 		jobs.POST("/batch_pause", apiJobCtl.BatchPause)
 		jobs.POST("/batch_delete", apiJobCtl.BatchDelete)

 		// 定时任务
 		timings := authed.Group("/timings")
 		timings.GET("", apiTimingCtl.List)
 		timings.GET("/:id", apiTimingCtl.Show)
 		timings.POST("", apiTimingCtl.Create)
 		timings.PUT("/:id", apiTimingCtl.Update)
 		timings.POST("/:id/copy", apiTimingCtl.Copy)
 		timings.POST("/:id/start", apiTimingCtl.Start)
 		timings.POST("/:id/restart", apiTimingCtl.ReStart)
 		timings.POST("/:id/pause", apiTimingCtl.Pause)
 		timings.DELETE("/:id", apiTimingCtl.Delete)
 		timings.POST("/batch_start", apiTimingCtl.BatchStart)
 		timings.POST("/batch_restart", apiTimingCtl.BatchReStart)
 		timings.POST("/batch_pause", apiTimingCtl.BatchPause)
 		timings.POST("/batch_delete", apiTimingCtl.BatchDelete)

 		// 监控
 		monitor := authed.Group("/monitor")
 		monitor.GET("/agent", apiMonitorCtl.Agent)
 		monitor.GET("/app", apiMonitorCtl.App)
 		monitor.GET("/job", apiMonitorCtl.Job)
 		monitor.GET("/timing", apiMonitorCtl.Timing)

 		// 归档
 		archives := authed.Group("/archives")
 		archives.GET("/app", apiArchiveCtl.App)
 		archives.GET("/job", apiArchiveCtl.Job)
 		archives.GET("/timing", apiArchiveCtl.Timing)

 		// 日志（数据接口；HTML 视图由前端实现）
 		outLogs := authed.Group("/out_logs")
 		outLogs.GET("/app/data", apiLogCtl.AppOutLogData)
 		outLogs.GET("/job/data", apiLogCtl.JobOutLogData)
 		outLogs.GET("/timing/data", apiLogCtl.TimingOutLogData)
 		errLogs := authed.Group("/err_logs")
 		errLogs.GET("/app/data", apiLogCtl.AppErrLogData)
 		errLogs.GET("/job/data", apiLogCtl.JobErrLogData)
 		errLogs.GET("/timing/data", apiLogCtl.TimingErrLogData)

 		// 异常 + 操作日志
 		authed.GET("/exceptions", apiExcCtl.List)
 		authed.GET("/operations", apiOpCtl.List)

 		// SSE：实时日志/监控
 		sse := authed.Group("/sse")
 		sse.GET("/out_log/app", apiSSELogCtl.AppOutLogStream)
 		sse.GET("/err_log/app", apiSSELogCtl.AppErrLogStream)
 		sse.GET("/out_log/job", apiSSELogCtl.JobOutLogStream)
 		sse.GET("/err_log/job", apiSSELogCtl.JobErrLogStream)
 		sse.GET("/out_log/timing", apiSSELogCtl.TimingOutLogStream)
 		sse.GET("/err_log/timing", apiSSELogCtl.TimingErrLogStream)
 		sse.GET("/monitor/agent", apiSSEMonCtl.AgentStream)
 		sse.GET("/monitor/app", apiSSEMonCtl.AppStream)
 		sse.GET("/monitor/job", apiSSEMonCtl.JobStream)
 		sse.GET("/monitor/timing", apiSSEMonCtl.TimingStream)
 	}
}
