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

	"Asgard/constants"
	"Asgard/web/controllers"
	"Asgard/web/utils"
)

var (
	server            *gin.Engine
	appController     *controllers.AppController
	jobController     *controllers.JobController
	agentController   *controllers.AgentController
	useController     *controllers.UserController
	groupController   *controllers.GroupController
	timingController  *controllers.TimingController
	monitorController *controllers.MonitorController
	archiveController *controllers.ArchiveController
	logController     *controllers.LogController
	indexController   *controllers.IndexController
)

func Server() *gin.Engine {
	return server
}

func Init() error {
	if constants.WEB_MODE == "release" {
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
		"unescaped": utils.Unescaped,
	}
	server.HTMLRender = ginview.New(viewConfig)
	server.Static("/assets", "web/assets")
	return nil
}

func setupController() {
	outDir := viper.GetString("log.dir")
	if outDir != "" {
		constants.WEB_OUT_DIR = outDir
	}
	useController = controllers.NewUserController()
	agentController = controllers.NewAgentController()
	groupController = controllers.NewGroupController()
	appController = controllers.NewAppController()
	jobController = controllers.NewJobController()
	timingController = controllers.NewTimingController()
	indexController = controllers.NewIndexController()
	monitorController = controllers.NewMonitorController()
}

func Run() {
	setupController()
	setupRouter()
	addr := fmt.Sprintf(":%d", constants.WEB_PORT)
	err := server.Run(addr)
	if err != nil {
		panic("web服务启动失败!")
	}
}
