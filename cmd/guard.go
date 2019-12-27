package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgarde/structs"
)

func init() {
	guardCommonCmd.PersistentFlags().StringP("conf", "c", "conf", "config path")
	rootCmd.AddCommand(guardCommonCmd)
}

var guardCommonCmd = &cobra.Command{
	Use:   "guard",
	Short: "guard apps",
	Run: func(cmd *cobra.Command, args []string) {
		confPath := cmd.Flag("conf").Value.String()
		viper.SetConfigName("app")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(confPath)
		//init config
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("init config error:", err.Error())
			return
		}
		configs := viper.GetStringMap("app")
		if len(configs) == 0 {
			fmt.Println("no apps!")
			return
		}
		apps := []*structs.App{}
		for key := range configs {
			config := viper.GetStringMapString("app." + key)
			app := &structs.App{
				Name:    config["name"],
				Dir:     config["dir"],
				Program: config["program"],
				Args:    config["args"],
				Stdout:  config["stdout"],
				Stderr:  config["stderr"],
			}
			apps = append(apps, app)
			go run(app)
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
		for s := range c {
			switch s {
			case os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				//kill all app before exit
				killAll(apps)
				os.Exit(0)
			}
		}
	},
}

func killAll(apps []*structs.App) {
	//TODO
	fmt.Println("killAll")
}

func run(app *structs.App) {
	app.Cmd = exec.Command(app.Program, app.Args)
	app.Cmd.Dir = app.Dir
	stdout, err := os.OpenFile(app.Stdout, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("open stdout error:", err)
		return
	}
	stderr, err := os.OpenFile(app.Stderr, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("open stderr error:", err)
		return
	}
	app.Cmd.Stdout = stdout
	app.Cmd.Stderr = stderr
	err = app.Cmd.Start()
	if err != nil {
		fmt.Println(app.Name+" start fail:", err)
		return
	}
	app.Begin = time.Now()
	app.Pid = app.Cmd.Process.Pid
	fmt.Println(app.Name+" start on ", app.Pid)
	go wait(app)
	go moniter(app)
}

func wait(app *structs.App) {
	err := app.Cmd.Wait()
	status := app.Cmd.ProcessState.Sys().(syscall.WaitStatus)
	signaled := status.Signaled()
	signal := status.Signal()
	if signaled {
		fmt.Println(app.Name+" signaled:", signal.String())
	} else if err != nil {
		fmt.Println(app.Name+" wait fail:", err)
	}
	app.End = time.Now()
	app.Finished = true
	fmt.Println(app.Name+" exit:", app.Cmd.ProcessState.ExitCode())
}

func moniter(app *structs.App) {
	loop := time.Second * time.Duration(5)
	ticker := time.NewTicker(loop)
	defer func() {
		ticker.Stop()
	}()
	for range ticker.C {
		//app finish exit ticker
		if app.Finished {
			return
		}
		info, err := process.NewProcess(int32(app.Pid))
		if err != nil {
			fmt.Println(app.Name+" process info fail:", err)
			continue
		}
		data := map[string]interface{}{}
		data["memory_percent"], _ = info.MemoryPercent()
		data["cpu_percent"], _ = info.CPUPercent()
		data["threads"], _ = info.NumThreads()
		//fmt.Println(app.Name+" process info:", data)
	}
}
