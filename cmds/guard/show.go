package guard

import (
	"fmt"

	"github.com/spf13/cobra"

	"Asgard/clients"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show app status",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("args worng!")
			return
		}
		serverFile := cmd.Flag("socket").Value.String()
		client, err := clients.NewGuard(serverFile)
		if err != nil {
			fmt.Printf("fail connect to guard:%s\n", err.Error())
			return
		}
		app, err := client.Get(args[0])
		if err != nil {
			fmt.Printf("fail to get app:%s\n", err.Error())
			return
		}
		if app == nil {
			fmt.Println("app no exist")
			return
		}
		fmt.Println("app info:")
		fmt.Println("app name:", app.GetName())
		fmt.Println("app dir:", app.GetDir())
		fmt.Println("app cmd:", app.GetProgram(), app.GetArgs())
		fmt.Println("app std_out:", app.GetStdOut())
		fmt.Println("app err_out:", app.GetStdErr())
		fmt.Println("app auto_restart:", app.GetAutoRestart())
	},
}
