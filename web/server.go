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
	middlewares "github.com/dalonghahaha/avenger/middlewares/gin"

	"Asgard/web/controllers"
)

var (
	server          *gin.Engine
	appController   *controllers.AppController
	agentController *controllers.AgentController
	useController   *controllers.UserController
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
	server = gin.New()
	server.Use(middlewares.Logger)
	server.Use(middlewares.Recover)
	viewConfig := goview.DefaultConfig
	viewConfig.Root = "web/views"
	viewConfig.Funcs = template.FuncMap{
		"unescaped": controllers.Unescaped,
	}
	server.HTMLRender = ginview.New(viewConfig)
	server.Static("/assets", "web/assets")
	appController = controllers.NewAppController()
	agentController = controllers.NewAgentController()
	useController = controllers.NewUserController()
	return nil
}

func setupRouter() {
	server.GET("/ping", controllers.Ping)
	server.GET("/", controllers.Index)
	server.GET("/register", useController.Register)
	server.POST("/register", useController.DoRegister)
	app := server.Group("/app")
	{
		app.GET("/list", appController.List)
	}
	agent := server.Group("/agent")
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
