package web

import (
	"fmt"
	"os"

	"github.com/dalonghahaha/avenger/components/cache"
	"github.com/dalonghahaha/avenger/components/db"
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/constants"
	"Asgard/runtimes"
	"Asgard/web"
)

func GetCmd() *cobra.Command {
	webCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	return webCmd
}

func PreRun(cmd *cobra.Command, args []string) {
	confPath := cmd.Flag("conf").Value.String()
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = logger.Register()
	if err != nil {
		panic(err)
	}
	systemMoniter := viper.GetInt("system.moniter")
	if systemMoniter > 0 {
		constants.SYSTEM_MONITER = systemMoniter
	}
	systemTimer := viper.GetInt("system.timer")
	if systemMoniter > 0 {
		constants.SYSTEM_TIMER = systemTimer
	}
}

var webCmd = &cobra.Command{
	Use:    "web",
	Short:  "run as web server",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		if err := InitWebServer(); err != nil {
			fmt.Println(err)
			return
		}
		go StartWebServer()
		runtimes.Wait(StopWebServer)
	},
}

func InitWebServer() error {
	err := db.Register()
	if err != nil {
		return err
	}
	err = cache.Register()
	if err != nil {
		return err
	}
	port := viper.GetInt("web.port")
	if port != 0 {
		constants.WEB_PORT = port
	}
	mode := viper.GetString("web.mode")
	if mode != "" {
		constants.WEB_MODE = mode
	}
	domain := viper.GetString("web.domain")
	if domain != "" {
		constants.WEB_DOMAIN = domain
	}
	cookieSalt := viper.GetString("web.cookie_salt")
	if cookieSalt != "" {
		constants.WEB_COOKIE_SALT = cookieSalt
	}
	return nil
}

func StartWebServer() {
	err := web.Init()
	if err != nil {
		logger.Error("web init error:", err)
		os.Exit(1)
	}
	logger.Info("Web Server Started!")
	logger.Debugf("Server Port:%d", constants.WEB_PORT)
	logger.Debugf("Server Pid:%d", os.Getpid())
	web.Run()
}

func StopWebServer() {
	logger.Info("Web Server Stop!")
}
