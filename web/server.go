package web

import (
	"fmt"
	"html/template"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/dalonghahaha/avenger/components/db"
	"github.com/dalonghahaha/avenger/components/logger"
	common_middlewares "github.com/dalonghahaha/avenger/middlewares/gin"

	"Asgard/web/controllers"
	"Asgard/web/middlewares"
)

var (
	server          *gin.Engine
	appController   *controllers.AppController
	jobController   *controllers.JobController
	agentController *controllers.AgentController
	useController   *controllers.UserController
	groupController *controllers.GroupController
)

func Init() error {
	//初始化日志组件
	err := logger.Register()
	if err != nil {
		return fmt.Errorf("日志组件初始化错误：%e", err)
	}
	//初始化数据库组件
	err = db.Register()
	if err != nil {
		return fmt.Errorf("数据库组件初始化错误：%e", err)
	}
	if viper.GetString("server.mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	cookieSalt := viper.GetString("server.cookie_salt")
	if cookieSalt != "" {
		controllers.CookieSalt = cookieSalt
	}
	domain := viper.GetString("server.domain")
	if cookieSalt != "" {
		controllers.Domain = domain
	}
	server = gin.New()
	server.Use(common_middlewares.Logger)
	server.Use(common_middlewares.Recover)
	viewConfig := goview.DefaultConfig
	viewConfig.Root = "web/views"
	viewConfig.DisableCache = true
	viewConfig.Funcs = template.FuncMap{
		"unescaped": controllers.Unescaped,
	}
	server.HTMLRender = ginview.New(viewConfig)
	server.Static("/assets", "web/assets")
	appController = controllers.NewAppController()
	agentController = controllers.NewAgentController()
	useController = controllers.NewUserController()
	groupController = controllers.NewGroupController()
	jobController = controllers.NewJobController()
	return nil
}

func setupRouter() {
	server.GET("/ping", controllers.Ping)
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
		user.GET("/list", useController.List)
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
	}
	agent := server.Group("/agent")
	agent.Use(middlewares.Login)
	{
		agent.GET("/list", agentController.List)
	}
}

func Run() {
	setupRouter()
	addr := fmt.Sprintf(":%s", viper.GetString("master.web.port"))
	err := server.Run(addr)
	if err != nil {
		panic("web服务启动失败!")
	}
}
