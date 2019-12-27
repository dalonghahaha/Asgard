package applications

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/shirou/gopsutil/process"
)

var APPs = []*App{}

type App struct {
	Name        string
	Dir         string
	Program     string
	Args        string
	Stdout      string
	Stderr      string
	Pid         int
	Cmd         *exec.Cmd
	Finished    bool
	AutoRestart bool
	Begin       time.Time
	End         time.Time
}

func KillAll() {
	for _, app := range APPs {
		if !app.Finished {
			err := app.Cmd.Process.Kill()
			if err != nil {
				logger.Error(err)
			} else {
				logger.Debug(app.Name + " killed!")
			}
		}
	}
}

func StartAll() {
	for _, app := range APPs {
		go app.Run()
	}
}

func Start(name string) bool {
	for _, app := range APPs {
		if app.Name == name {
			go app.Run()
			return true
		}
	}
	return false
}

func Register(config map[string]string) *App {
	app := App{
		Name:    config["name"],
		Dir:     config["dir"],
		Program: config["program"],
		Args:    config["args"],
		Stdout:  config["stdout"],
		Stderr:  config["stderr"],
	}
	APPs = append(APPs, &app)
	return &app
}

func (a *App) Run() {
	a.Cmd = exec.Command(a.Program, a.Args)
	a.Cmd.Dir = a.Dir
	stdout, err := os.OpenFile(a.Stdout, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logger.Error("open stdout error:", err)
		return
	}
	stderr, err := os.OpenFile(a.Stderr, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logger.Error("open stderr error:", err)
		return
	}
	a.Cmd.Stdout = stdout
	a.Cmd.Stderr = stderr
	err = a.Cmd.Start()
	if err != nil {
		logger.Error(a.Name+" start fail:", err)
		return
	}
	a.Begin = time.Now()
	a.Pid = a.Cmd.Process.Pid
	logger.Debug(a.Name+" started at ", a.Pid)
	go a.wait()
	go a.moniter()
}

func (a *App) wait() {
	err := a.Cmd.Wait()
	status := a.Cmd.ProcessState.Sys().(syscall.WaitStatus)
	signaled := status.Signaled()
	signal := status.Signal()
	if signaled {
		logger.Info(a.Name+" signaled:", signal.String())
	} else if err != nil {
		logger.Error(a.Name+" wait fail:", err)
	}
	a.End = time.Now()
	a.Finished = true
	logger.Debug(a.Name+" exit with ", a.Cmd.ProcessState.ExitCode())
}

func (a *App) moniter() {
	loop := time.Second * time.Duration(5)
	ticker := time.NewTicker(loop)
	defer func() {
		ticker.Stop()
	}()
	for range ticker.C {
		//app finish exit ticker
		if a.Finished {
			return
		}
		info, err := process.NewProcess(int32(a.Pid))
		if err != nil {
			logger.Error(a.Name+" process info err:", err)
			continue
		}
		data := map[string]interface{}{}
		data["memory_percent"], _ = info.MemoryPercent()
		data["cpu_percent"], _ = info.CPUPercent()
		data["threads"], _ = info.NumThreads()
		logger.Debug(a.Name+" process info:", data)
	}
}
