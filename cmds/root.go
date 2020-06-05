package cmds

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/cmds/cron"
	"Asgard/constants"
	"Asgard/managers"
)

var (
	agentManager *managers.AgentManager
	appManager   *managers.AppManager
)

func init() {
	RootCmd.AddCommand(cron.GetCmd())
}

var RootCmd = &cobra.Command{
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

func Wait(function func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGTSTP)
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGTSTP:
			function()
			os.Exit(0)
		}
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
