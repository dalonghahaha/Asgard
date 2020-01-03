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
	cronCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(cronCommonCmd)
}

var cronCommonCmd = &cobra.Command{
	Use:    "cron",
	Short:  "cron jobs",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		StartCron()
		NotityKill(applications.JobStopAll)
	},
}

func StartCron() {
	configs := viper.Get("cron")
	if configs == nil {
		fmt.Println("no crons!")
		return
	}
	_configs, ok := configs.([]interface{})
	if !ok {
		fmt.Println("crons config wrong!")
		return
	}
	for _, v := range _configs {
		_v, ok := v.(map[interface{}]interface{})
		if !ok {
			fmt.Println("crons config wrong!")
			return
		}
		config := map[string]interface{}{}
		for k, v := range _v {
			_k, ok := k.(string)
			if !ok {
				fmt.Println("crons config wrong!")
				return
			}
			config[_k] = v
		}
		err := applications.JobRegister(config)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	logger.Info("cron started at ", os.Getpid())
	applications.JobStartAll(true)
}
