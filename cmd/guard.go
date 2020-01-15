package cmd

import (
	"fmt"
	"os"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/applications"
)

func init() {
	guardCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(guardCommonCmd)
}

var guardCommonCmd = &cobra.Command{
	Use:    "guard",
	Short:  "guard apps",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		StartGuard()
		NotityKill(applications.AppStopAll)
	},
}

func StartGuard() {
	configs := viper.Get("app")
	if configs == nil {
		fmt.Println("no apps!")
		return
	}
	_configs, ok := configs.([]interface{})
	if !ok {
		fmt.Println("apps config wrong!")
		return
	}
	for index, v := range _configs {
		_v, ok := v.(map[interface{}]interface{})
		if !ok {
			fmt.Println("apps config wrong!")
			return
		}
		config := map[string]interface{}{}
		for k, v := range _v {
			_k, ok := k.(string)
			if !ok {
				fmt.Println("apps config wrong!")
				return
			}
			config[_k] = v
		}
		_, err := applications.AppRegister(int64(index), config)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	logger.Info("guard started at ", os.Getpid())
	applications.AppStartAll(true)
}
