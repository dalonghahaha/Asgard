package cmd

import (
	"os"

	"github.com/dalonghahaha/avenger/components/cache"
	"github.com/dalonghahaha/avenger/components/db"
	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"

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
		err := db.Register()
		if err != nil {
			panic(err)
		}
		err = cache.Register()
		if err != nil {
			panic(err)
		}
		go StartWebServer()
		NotityKill(StopWebServer)
	},
}

func StartWebServer() {
	err := web.Init()
	if err != nil {
		logger.Error("web init error:", err)
		os.Exit(1)
	}
	logger.Info("Web Server Started! Pid:", os.Getpid())
	web.Run()
}

func StopWebServer() {
	logger.Info("Web Server Stop!")
}
