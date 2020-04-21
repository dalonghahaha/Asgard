package web

import (
	"Asgard/web/controllers"
	"Asgard/web/middlewares"
)

func setupRouter() {
	server.GET("/ping", controllers.Ping)
	server.GET("/UI", controllers.UI)
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
		agent.GET("/monitor", agentController.Monitor)
		agent.GET("/edit", agentController.Edit)
		agent.POST("/update", agentController.Update)
	}
	group := server.Group("/group")
	group.Use(middlewares.Login)
	{
		group.GET("/list", groupController.List)
		group.GET("/add", groupController.Add)
		group.POST("/create", groupController.Create)
		group.GET("/edit", groupController.Edit)
		group.POST("/update", groupController.Update)
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
		app.GET("/monitor", middlewares.AppInit, appController.Monitor)
		app.GET("/archive", middlewares.AppInit, appController.Archive)
		app.GET("/copy", middlewares.AppInit, appController.Copy)
		app.GET("/delete", middlewares.AppInit, appController.Delete)

		app.GET("/start", middlewares.AppAgentInit, appController.Start)
		app.GET("/restart", middlewares.AppAgentInit, appController.ReStart)
		app.GET("/pause", middlewares.AppAgentInit, appController.Pause)
		app.GET("/out_log", middlewares.AppAgentInit, appController.OutLog)
		app.GET("/err_log", middlewares.AppAgentInit, appController.ErrLog)
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
		job.GET("/monitor", middlewares.JobInit, jobController.Monitor)
		job.GET("/archive", middlewares.JobInit, jobController.Archive)
		job.GET("/copy", middlewares.JobInit, jobController.Copy)
		job.GET("/delete", middlewares.JobInit, jobController.Delete)

		job.GET("/start", middlewares.JobAgentInit, jobController.Start)
		job.GET("/restart", middlewares.JobAgentInit, jobController.ReStart)
		job.GET("/pause", middlewares.JobAgentInit, jobController.Pause)
		job.GET("/out_log", middlewares.JobAgentInit, jobController.OutLog)
		job.GET("/err_log", middlewares.JobAgentInit, jobController.ErrLog)
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
		timing.GET("/monitor", middlewares.TimingInit, timingController.Monitor)
		timing.GET("/archive", middlewares.TimingInit, timingController.Archive)
		timing.GET("/copy", middlewares.TimingInit, timingController.Copy)
		timing.GET("/delete", middlewares.TimingInit, timingController.Delete)

		timing.GET("/start", middlewares.TimingAgentInit, timingController.Start)
		timing.GET("/restart", middlewares.TimingAgentInit, timingController.ReStart)
		timing.GET("/pause", middlewares.TimingAgentInit, timingController.Pause)
		timing.GET("/out_log", middlewares.TimingAgentInit, timingController.OutLog)
		timing.GET("/err_log", middlewares.TimingAgentInit, timingController.ErrLog)
	}
}
