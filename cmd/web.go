package cmd

import (
	"os"

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
	Short:  "web服务模式",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		err := web.Init()
		if err != nil {
			logger.Error("web init error:", err)
			os.Exit(1)
		}
		web.Run()
	},
}
