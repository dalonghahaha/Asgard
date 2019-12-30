package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/applications"
)


func init() {
	cronCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(guardCommonCmd)
}

var cronCommonCmd = &cobra.Command{
	Use:   "cron",
	Short: "cron jobs",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		configs := viper.GetStringMap("cron")
		if len(configs) == 0 {
			fmt.Println("no jobs!")
			return
		}
		for key := range configs {
			config := viper.GetStringMapString("cron." + key)
			applications.RegisterJob(config)
		}
		applications.CronAll()
		NotityKill(func(){})
	},
}
