package cmd

import (
	"fmt"
	"os"

	"github.com/dalonghahaha/avenger/components/cache"
	"github.com/dalonghahaha/avenger/components/db"
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/constants"
	"Asgard/web"
)

func init() {
	webCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:    "web",
	Short:  "run as web server",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		InitWebServer()
		go StartWebServer()
		NotityKill(StopWebServer)
	},
}

func InitWebServer() {
	err := db.Register()
	if err != nil {
		panic(err)
	}
	err = cache.Register()
	if err != nil {
		panic(err)
	}
	port := viper.GetInt("web.port")
	if port != 0 {
		constants.WEB_PORT = port
	}
	domain := viper.GetString("web.domain")
	if domain != "" {
		constants.WEB_DOMAIN = domain
	}
	cookieSalt := viper.GetString("web.cookie_salt")
	if cookieSalt != "" {
		constants.WEB_COOKIE_SALT = cookieSalt
	}
}

func StartWebServer() {
	err := web.Init()
	if err != nil {
		logger.Error("web init error:", err)
		os.Exit(1)
	}
	logger.Info("Web Server Started!", os.Getpid())
	logger.Debug(fmt.Sprintf("Server Port:%d", constants.WEB_PORT))
	logger.Debug(fmt.Sprintf("Server Pid:%d", os.Getpid()))
	web.Run()
}

func StopWebServer() {
	logger.Info("Web Server Stop!")
}
