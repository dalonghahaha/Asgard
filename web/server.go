package web

import (
	"fmt"
	"html/template"
	"strings"

	common_middlewares "github.com/dalonghahaha/avenger/middlewares/gin"
	"github.com/dalonghahaha/avenger/tools/file"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"Asgard/web/controllers"
)

var (
	server           *gin.Engine
	appController    *controllers.AppController
	jobController    *controllers.JobController
	agentController  *controllers.AgentController
	useController    *controllers.UserController
	groupController  *controllers.GroupController
	timingController *controllers.TimingController
	indexController  *controllers.IndexController
)

func Server() *gin.Engine {
	return server
}

func Init() error {
	if viper.GetString("server.mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	server = gin.New()
	server.Use(common_middlewares.Logger)
	server.Use(common_middlewares.Recover)
	viewConfig := goview.DefaultConfig
	viewConfig.Root = "web/views"
	viewConfig.Extension = ".html"
	viewConfig.Master = "layouts/master"
	fileList, err := file.ReadDir("web/views/templates")
	if err != nil {
		return err
	}
	partials := []string{}
	for _, file := range fileList {
		partial := "templates/" + strings.Replace(file.Name(), viewConfig.Extension, "", -1)
		partials = append(partials, partial)
	}
	viewConfig.Partials = partials
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
	outDir := viper.GetString("log.dir")
	if outDir != "" {
		controllers.OutDir = outDir
	}
	useController = controllers.NewUserController()
	agentController = controllers.NewAgentController()
	groupController = controllers.NewGroupController()
	appController = controllers.NewAppController()
	jobController = controllers.NewJobController()
	timingController = controllers.NewTimingController()
	indexController = controllers.NewIndexController()
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
