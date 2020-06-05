package cron

import (
	"fmt"

	"github.com/spf13/cobra"

	"Asgard/clients"
)

var statusCommonCmd = &cobra.Command{
	Use:   "status",
	Short: "show cron runing jobs",
	Run: func(cmd *cobra.Command, args []string) {
		serverFile := cmd.Flag("socket").Value.String()
		client, err := clients.NewCron(serverFile)
		if err != nil {
			fmt.Printf("fail connect to guard:%s\n", err.Error())
			return
		}
		jobs, err := client.GetList()
		if err != nil {
			fmt.Printf("fail to get app list:%s\n", err.Error())
			return
		}
		titleFormat := "%-4s %-15s %-50s %-50s %-30s\n"
		contentFormat := "%-4d %-15s %-50s %-50s %-30s\n"
		fmt.Println("app total:", len(jobs))
		fmt.Printf(titleFormat, "ID", "Spec", "Dir", "Program", "Name")
		for _, job := range jobs {
			program := fmt.Sprintf("%s %s", job.GetProgram(), job.GetArgs())
			fmt.Printf(contentFormat, job.GetId(), job.GetSpec(), job.GetDir(), program, job.GetName())
		}
	},
}
