package runtimes

import (
	"Asgard/constants"

	"github.com/spf13/viper"
)

func ParseConfig(confPath string) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confPath)
	err := viper.ReadInConfig()
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
