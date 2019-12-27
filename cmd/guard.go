package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgarde/applications"
)

func init() {
	guardCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(guardCommonCmd)
}

var guardCommonCmd = &cobra.Command{
	Use:   "guard",
	Short: "guard apps",
	PreRun: func(cmd *cobra.Command, args []string) {
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
	},
	Run: func(cmd *cobra.Command, args []string) {
		configs := viper.GetStringMap("app")
		if len(configs) == 0 {
			fmt.Println("no apps!")
			return
		}
		for key := range configs {
			config := viper.GetStringMapString("app." + key)
			applications.Register(config)
		}
		applications.StartAll()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
		for s := range c {
			switch s {
			case os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				//kill all app before exit
				applications.KillAll()
				os.Exit(0)
			}
		}
	},
}
