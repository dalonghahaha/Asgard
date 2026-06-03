 package web

 import (
 	"fmt"

 	common_middlewares "github.com/dalonghahaha/avenger/middlewares/gin"
 	"github.com/gin-gonic/gin"
 	"github.com/spf13/viper"

 	"Asgard/constants"
 )

 // 前后端分离后 web 层只承载纯 JSON API；HTML 渲染相关（goview、controllers）已下线。
 // 旧 HTML 控制器/中间件/模板/静态资源已迁移到 web/legacy/ 与 doc/legacy-templates/。

 var server *gin.Engine

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
 	return nil
 }

 func Run() {
 	if outDir := viper.GetString("log.dir"); outDir != "" {
 		constants.WEB_OUT_DIR = outDir
 	}
 	setupRouter()
 	addr := fmt.Sprintf(":%d", constants.WEB_PORT)
 	if err := server.Run(addr); err != nil {
 		panic("web服务启动失败!")
 	}
 }
