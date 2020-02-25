package web

import (
	"fmt"
	"html/template"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	common_middlewares "github.com/dalonghahaha/avenger/middlewares/gin"

	"Asgard/web/controllers"
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
	if viper.GetString("server.mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
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
	return nil
}

func setupController() {
	cookieSalt := viper.GetString("server.cookie_salt")
	if cookieSalt != "" {
		controllers.CookieSalt = cookieSalt
	}
	domain := viper.GetString("server.domain")
	if cookieSalt != "" {
		controllers.Domain = domain
	}
	appController = controllers.NewAppController()
	agentController = controllers.NewAgentController()
	useController = controllers.NewUserController()
	groupController = controllers.NewGroupController()
	jobController = controllers.NewJobController()
}

func Run() {
	setupController()
	setupRouter()
	addr := fmt.Sprintf(":%s", viper.GetString("master.web.port"))
	err := server.Run(addr)
	if err != nil {
		panic("web服务启动失败!")
	}
}
