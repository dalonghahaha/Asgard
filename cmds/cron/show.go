package cron

import (
	"fmt"

	"github.com/spf13/cobra"

	"Asgard/clients"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show job status",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("args worng!")
			return
		}
		serverFile := cmd.Flag("socket").Value.String()
		client, err := clients.NewCron(serverFile)
		if err != nil {
			fmt.Printf("fail connect to guard:%s\n", err.Error())
			return
		}
		job, err := client.Get(args[0])
		if err != nil {
			fmt.Printf("fail to get job:%s\n", err.Error())
			return
		}
		if job == nil {
			fmt.Println("job no exist")
			return
		}
		fmt.Println("job info:")
		fmt.Println("job spec:", job.GetSpec())
		fmt.Println("job name:", job.GetName())
		fmt.Println("job timeout:", job.GetTimeout())
		fmt.Println("job dir:", job.GetDir())
		fmt.Println("job cmd:", job.GetProgram(), job.GetArgs())
		fmt.Println("job std_out:", job.GetStdOut())
		fmt.Println("job err_out:", job.GetStdErr())
	},
}
