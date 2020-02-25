package web

import (
	"Asgard/web/controllers"
	"Asgard/web/middlewares"
)

func setupRouter() {
	server.GET("/ping", controllers.Ping)
	server.GET("/UI", controllers.UI)
	server.GET("/", middlewares.Login, controllers.Index)
	server.GET("/nologin", controllers.Nologin)
	server.GET("/error", controllers.Error)
	server.GET("/register", useController.Register)
	server.POST("/register", useController.DoRegister)
	server.GET("/login", useController.Login)
	server.POST("/login", useController.DoLogin)
	user := server.Group("/user")
	user.Use(middlewares.Login)
	{
		user.GET("/info", useController.Info)
		user.GET("/list", useController.List)
		user.GET("/setting", useController.Setting)
		user.POST("/setting", useController.Update)
		user.GET("/change_password", useController.ChangePassword)
		user.POST("/change_password", useController.DoChangePassword)
	}
	agent := server.Group("/agent")
	agent.Use(middlewares.Login)
	{
		agent.GET("/list", agentController.List)
		agent.GET("/monitor", agentController.Monitor)
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
		app.GET("/show", appController.Show)
		app.GET("/add", appController.Add)
		app.POST("/create", appController.Create)
		app.GET("/edit", appController.Edit)
		app.POST("/update", appController.Update)
		app.GET("/monitor", appController.Monitor)
		app.GET("/archive", appController.Archive)
		app.GET("/delete", appController.Delete)
		app.GET("/start", appController.Start)
		app.GET("/restart", appController.ReStart)
		app.GET("/pause", appController.Pause)
	}
	job := server.Group("/job")
	job.Use(middlewares.Login)
	{
		job.GET("/list", jobController.List)
		job.GET("/show", jobController.Show)
		job.GET("/add", jobController.Add)
		job.POST("/create", jobController.Create)
		job.GET("/edit", jobController.Edit)
		job.POST("/update", jobController.Update)
		job.GET("/monitor", jobController.Monitor)
		job.GET("/archive", jobController.Archive)
		job.GET("/delete", jobController.Delete)
		job.GET("/start", jobController.Start)
		job.GET("/restart", jobController.ReStart)
		job.GET("/pause", jobController.Pause)
	}
	timing := server.Group("/timing")
	timing.Use(middlewares.Login)
	{
		timing.GET("/list", timingController.List)
		timing.GET("/show", timingController.Show)
		timing.GET("/add", timingController.Add)
		timing.POST("/create", timingController.Create)
		timing.GET("/edit", timingController.Edit)
		timing.POST("/update", timingController.Update)
		timing.GET("/monitor", timingController.Monitor)
		timing.GET("/archive", timingController.Archive)
		timing.GET("/delete", timingController.Delete)
	}
}
