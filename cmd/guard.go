package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/applications"
)

func init() {
	guardCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(guardCommonCmd)
}

var guardCommonCmd = &cobra.Command{
	Use:   "guard",
	Short: "guard apps",
	PreRun: PreRun,
	Run: func(cmd *cobra.Command, args []string) {
		configs := viper.GetStringMap("app")
		if len(configs) == 0 {
			fmt.Println("no apps!")
			return
		}
		for key := range configs {
			config := viper.GetStringMapString("app." + key)
			applications.RegisterApp(config)
		}
		applications.StartAll()
		NotityKill(applications.KillAll)
	},
}
