package cmds

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/clients"
	"Asgard/constants"
	"Asgard/managers"
)

var (
	jobManager    *managers.JobManager
	appManager    *managers.AppManager
	timingManager *managers.TimingManager
	masterClient  *clients.Master
)

var rootCmd = &cobra.Command{
	Use:   "Asgard",
	Short: "welcome to use Asgard!",
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

func NotityKill(function func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range c {
		switch s {
		case os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			function()
			os.Exit(0)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
