package agent

import (
	"Asgard/clients"
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show agent running status",
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("port").Value.String()
		client, err := clients.NewAgent("127.0.0.1", port)
		if err != nil {
			fmt.Printf("fail connect to agent:%s\n", err.Error())
			return
		}
		apps, err := client.GetAppList()
		if err != nil {
			fmt.Printf("fail to get app list:%s\n", err.Error())
			return
		}
		jobs, err := client.GetJobList()
		if err != nil {
			fmt.Printf("fail to get job list:%s\n", err.Error())
			return
		}
		timings, err := client.GetTimingList()
		if err != nil {
			fmt.Printf("fail to get timing list:%s\n", err.Error())
			return
		}
		titleFormat := "%-5s %-50s %-50s %-30s\n"
		contentFormat := "%-5d %-50s %-50s %-30s\n"
		fmt.Println()
		fmt.Println("app total:", len(apps))
		fmt.Println()
		fmt.Printf(titleFormat, "ID", "Dir", "Program", "Name")
		fmt.Println()
		for _, app := range apps {
			program := fmt.Sprintf("%s %s", app.GetProgram(), app.GetArgs())
			fmt.Printf(contentFormat, app.GetId(), app.GetDir(), program, app.GetName())
		}
		fmt.Println()
		fmt.Println("job list:", len(jobs))
		fmt.Println()
		fmt.Printf(titleFormat, "ID", "Dir", "Program", "Name")
		fmt.Println()
		for _, job := range jobs {
			program := fmt.Sprintf("%s %s", job.GetProgram(), job.GetArgs())
			fmt.Printf(contentFormat, job.GetId(), job.GetDir(), program, job.GetName())
		}
		fmt.Println()
		fmt.Println("timing list:", len(timings))
		fmt.Println()
		fmt.Printf(titleFormat, "ID", "Dir", "Program", "Name")
		fmt.Println()
		for _, timing := range timings {
			program := fmt.Sprintf("%s %s", timing.GetProgram(), timing.GetArgs())
			fmt.Printf(contentFormat, timing.GetId(), timing.GetDir(), program, timing.GetName())
		}
	},
}
