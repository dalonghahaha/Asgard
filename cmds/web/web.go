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

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "run as web server",
	Run: func(cmd *cobra.Command, args []string) {
		confPath := cmd.Flag("conf").Value.String()
		runtimes.ParseConfig(confPath)
		if err := InitWebServer(); err != nil {
			fmt.Println(err)
			return
		}
		go StartWebServer()
		runtimes.Wait(StopWebServer)
	},
}

func InitWebServer() error {
	if err := logger.Register(); err != nil {
		return fmt.Errorf("init logger failed:%+v", err)
	}
	if err := db.Register(); err != nil {
		return fmt.Errorf("init db failed:%+v", err)
	}
	if err := cache.Register(); err != nil {
		return fmt.Errorf("init cache failed:%+v", err)
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
