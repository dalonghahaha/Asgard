package guard

import (
	"Asgard/clients"
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show guard running status",
	Run: func(cmd *cobra.Command, args []string) {
		serverFile := cmd.Flag("socket").Value.String()
		client, err := clients.NewGuard(serverFile)
		if err != nil {
			fmt.Printf("fail connect to guard:%s\n", err.Error())
			return
		}
		apps, err := client.GetList()
		if err != nil {
			fmt.Printf("fail to get app list:%s\n", err.Error())
			return
		}
		titleFormat := "%-5s %-50s %-50s %-30s\n"
		contentFormat := "%-5d %-50s %-50s %-30s\n"
		fmt.Println("app total:", len(apps))
		fmt.Printf(titleFormat, "ID", "Dir", "Program", "Name")
		for _, app := range apps {
			program := fmt.Sprintf("%s %s", app.GetProgram(), app.GetArgs())
			fmt.Printf(contentFormat, app.GetId(), app.GetDir(), program, app.GetName())
		}
	},
}
