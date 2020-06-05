package web

import (
	"Asgard/web/controllers"
	"Asgard/web/middlewares"
)

func setupRouter() {
	server.GET("/", middlewares.Login, indexController.Index)
	server.GET("/no_login", controllers.Nologin)
	server.GET("/auth_fail", controllers.AuthFail)
	server.GET("/forbidden", controllers.Forbidden)
	server.GET("/admin_only", controllers.AdminOnly)
	server.GET("/error", controllers.Error)
	server.GET("/register", useController.Register)
	server.POST("/register", useController.DoRegister)
	server.GET("/login", useController.Login)
	server.POST("/login", useController.DoLogin)
	user := server.Group("/user")
	user.Use(middlewares.Login)
	{
		user.GET("/info", useController.Info)
		user.GET("/setting", useController.Setting)
		user.POST("/setting", useController.DoSetting)
		user.GET("/change_password", useController.ChangePassword)
		user.POST("/change_password", useController.DoChangePassword)
		user.GET("/list", useController.List)
		user.GET("/add", middlewares.Admin, useController.Add)
		user.POST("/create", middlewares.Admin, useController.Create)
		user.GET("/edit", middlewares.Admin, useController.Edit)
		user.POST("/update", middlewares.Admin, useController.Update)
		user.GET("/forbidden", middlewares.Admin, useController.Forbidden)
		user.GET("/reset_password", middlewares.Admin, useController.ResetPassword)
		user.POST("/reset_password", middlewares.Admin, useController.DoResetPassword)
	}
	agent := server.Group("/agent")
	agent.Use(middlewares.Login)
	{
		agent.GET("/list", agentController.List)
		agent.GET("/edit", agentController.Edit)
		agent.POST("/update", middlewares.AgentInit, agentController.Update)
	}
	group := server.Group("/group")
	group.Use(middlewares.Login)
	{
		group.GET("/list", groupController.List)
		group.GET("/add", groupController.Add)
		group.POST("/create", groupController.Create)
		group.GET("/edit", groupController.Edit)
		group.POST("/update", middlewares.GroupInit, groupController.Update)
		group.POST("/delete", middlewares.GroupInit, groupController.Delete)
	}
	app := server.Group("/app")
	app.Use(middlewares.Login)
	{
		app.GET("/list", appController.List)
		app.GET("/show", middlewares.AppInit, appController.Show)
		app.GET("/add", appController.Add)
		app.POST("/create", middlewares.CmdConfigVerify, appController.Create)
		app.GET("/edit", middlewares.AppInit, appController.Edit)
		app.POST("/update", middlewares.AppInit, middlewares.CmdConfigVerify, appController.Update)
		app.POST("/copy", middlewares.AppInit, appController.Copy)
		//control
		app.POST("/start", middlewares.AppAgentInit, appController.Start)
		app.POST("/restart", middlewares.AppAgentInit, appController.ReStart)
		app.POST("/pause", middlewares.AppAgentInit, appController.Pause)
		app.POST("/delete", middlewares.AppAgentInit, appController.Delete)
		//batch-control
		app.POST("/batch-start", middlewares.BatchAppAgentInit, appController.BatchStart)
		app.POST("/batch-restart", middlewares.BatchAppAgentInit, appController.BatchReStart)
		app.POST("/batch-pause", middlewares.BatchAppAgentInit, appController.BatchPause)
		app.POST("/batch-delete", middlewares.BatchAppAgentInit, appController.BatchDelete)
	}
	job := server.Group("/job")
	job.Use(middlewares.Login)
	{
		job.GET("/list", jobController.List)
		job.GET("/show", middlewares.JobInit, jobController.Show)
		job.GET("/add", jobController.Add)
		job.POST("/create", middlewares.CmdConfigVerify, jobController.Create)
		job.GET("/edit", middlewares.JobInit, jobController.Edit)
		job.POST("/update", middlewares.JobInit, middlewares.CmdConfigVerify, jobController.Update)
		job.POST("/copy", middlewares.JobInit, jobController.Copy)
		//control
		job.POST("/start", middlewares.JobAgentInit, jobController.Start)
		job.POST("/restart", middlewares.JobAgentInit, jobController.ReStart)
		job.POST("/pause", middlewares.JobAgentInit, jobController.Pause)
		job.POST("/delete", middlewares.JobAgentInit, jobController.Delete)
		//batch-control
		job.POST("/batch-start", middlewares.BatchJobAgentInit, jobController.BatchStart)
		job.POST("/batch-restart", middlewares.BatchJobAgentInit, jobController.BatchReStart)
		job.POST("/batch-pause", middlewares.BatchJobAgentInit, jobController.BatchPause)
		job.POST("/batch-delete", middlewares.BatchJobAgentInit, jobController.BatchDelete)
	}
	timing := server.Group("/timing")
	timing.Use(middlewares.Login)
	{
		timing.GET("/list", timingController.List)
		timing.GET("/show", middlewares.TimingInit, timingController.Show)
		timing.GET("/add", timingController.Add)
		timing.POST("/create", middlewares.CmdConfigVerify, timingController.Create)
		timing.GET("/edit", middlewares.TimingInit, timingController.Edit)
		timing.POST("/update", middlewares.TimingInit, middlewares.CmdConfigVerify, timingController.Update)
		timing.POST("/copy", middlewares.TimingInit, timingController.Copy)
		//control
		timing.POST("/start", middlewares.TimingAgentInit, timingController.Start)
		timing.POST("/restart", middlewares.TimingAgentInit, timingController.ReStart)
		timing.POST("/pause", middlewares.TimingAgentInit, timingController.Pause)
		timing.POST("/delete", middlewares.TimingAgentInit, timingController.Delete)
		//batch-control
		timing.POST("/batch-start", middlewares.BatchTimingAgentInit, timingController.BatchStart)
		timing.POST("/batch-restart", middlewares.BatchTimingAgentInit, timingController.BatchReStart)
		timing.POST("/batch-pause", middlewares.BatchTimingAgentInit, timingController.BatchPause)
		timing.POST("/batch-delete", middlewares.BatchTimingAgentInit, timingController.BatchDelete)
	}
	monitor := server.Group("/monitor")
	monitor.Use(middlewares.Login)
	{
		monitor.GET("/agent", middlewares.AgentInit, monitorController.Agent)
		monitor.GET("/app", middlewares.AppInit, monitorController.App)
		monitor.GET("/job", middlewares.JobInit, monitorController.Job)
		monitor.GET("/timing", middlewares.TimingInit, monitorController.Timing)
	}
	archive := server.Group("/archive")
	archive.Use(middlewares.Login)
	{
		archive.GET("/app", middlewares.AppInit, archiveController.App)
		archive.GET("/job", middlewares.JobInit, archiveController.Job)
		archive.GET("/timing", middlewares.TimingInit, archiveController.Timing)
	}
	outLog := server.Group("/out_log")
	outLog.Use(middlewares.Login)
	{
		outLog.GET("/app", middlewares.AppAgentInit, logController.AppOutLog)
		outLog.GET("/job", middlewares.JobAgentInit, logController.JobOutLog)
		outLog.GET("/timing", middlewares.TimingAgentInit, logController.TimingOutLog)
		outLog.GET("/app/data", middlewares.AppAgentInit, logController.AppOutLogData)
		outLog.GET("/job/data", middlewares.JobAgentInit, logController.JobOutLogData)
		outLog.GET("/timing/data", middlewares.TimingAgentInit, logController.TimingOutLogData)
	}
	errLog := server.Group("/err_log")
	errLog.Use(middlewares.Login)
	{
		errLog.GET("/app", middlewares.AppAgentInit, logController.AppErrLog)
		errLog.GET("/job", middlewares.JobAgentInit, logController.JobErrLog)
		errLog.GET("/timing", middlewares.TimingAgentInit, logController.TimingErrLog)
		errLog.GET("/app/data", middlewares.AppAgentInit, logController.AppErrLogData)
		errLog.GET("/job/data", middlewares.JobAgentInit, logController.JobErrLogData)
		errLog.GET("/timing/data", middlewares.TimingAgentInit, logController.TimingErrLogData)
	}
	exception := server.Group("/exception")
	exception.Use(middlewares.Login)
	{
		exception.GET("/list", exceptionController.List)
	}
	operation := server.Group("/operation")
	operation.Use(middlewares.Login)
	{
		operation.GET("/list", operationController.List)
	}
}
