package cmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"Asgard/cmds/agent"
	"Asgard/cmds/cron"
	"Asgard/cmds/debug"
	"Asgard/cmds/guard"
	"Asgard/cmds/master"
	"Asgard/cmds/web"
)

func init() {
	RootCmd.AddCommand(agent.GetCmd())
	RootCmd.AddCommand(cron.GetCmd())
	RootCmd.AddCommand(debug.GetCmd())
	RootCmd.AddCommand(guard.GetCmd())
	RootCmd.AddCommand(master.GetCmd())
	RootCmd.AddCommand(web.GetCmd())
}

var RootCmd = &cobra.Command{
	Use:   "Asgard",
	Short: "welcome to use Asgard!",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
