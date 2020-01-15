package applications

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/file"
	"github.com/dalonghahaha/avenger/tools/uuid"
	"github.com/shirou/gopsutil/process"
)

var processExit = false

type Command struct {
	Name            string
	Dir             string
	Program         string
	Args            string
	Stdout          string
	Stderr          string
	IsMonitor       bool
	Pid             int
	UUID            string
	Begin           time.Time
	End             time.Time
	Finished        bool
	Status          int
	Cmd             *exec.Cmd
	ExceptionReport func(message string)
	MonitorReport   func(monitor *Monitor)
	ArchiveReport   func(command *Command)
}

func (c *Command) configure(config map[string]interface{}) error {
	name, ok := config["name"].(string)
	if !ok {
		return fmt.Errorf("config name type wrong")
	}
	c.Name = name
	dir, ok := config["dir"].(string)
	if !ok {
		return fmt.Errorf("config dir type wrong")
	}
	c.Dir = dir
	program, ok := config["program"].(string)
	if !ok {
		return fmt.Errorf("config program type wrong")
	}
	c.Program = program
	args, ok := config["args"].(string)
	if !ok {
		return fmt.Errorf("config args type wrong")
	}
	c.Args = args
	stdout, ok := config["stdout"].(string)
	if !ok {
		return fmt.Errorf("config stdout type wrong")
	}
	c.Stdout = stdout
	stderr, ok := config["stderr"].(string)
	if !ok {
		return fmt.Errorf("config stderr type wrong")
	}
	c.Stderr = stderr
	isMonitor, ok := config["is_monitor"].(bool)
	if !ok {
		return fmt.Errorf("config is_monitor type wrong")
	}
	c.IsMonitor = isMonitor
	return nil
}

func (c *Command) build() error {
	c.Cmd = exec.Command(c.Program, c.Args)
	c.Cmd.Dir = c.Dir
	if !file.Exists(c.Stdout) {
		err := file.Mkdir(filepath.Dir(c.Stdout))
		if err != nil {
			logger.Error("mkdir stdout error:", err)
			return err
		}
		_, err = os.Create(c.Stdout)
		if err != nil {
			logger.Error("create stdout error:", err)
			return err
		}
	}
	stdout, err := os.OpenFile(c.Stdout, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logger.Error("open stdout error:", err)
		return err
	}
	c.Cmd.Stdout = stdout
	if !file.Exists(c.Stderr) {
		err := file.Mkdir(filepath.Dir(c.Stderr))
		if err != nil {
			logger.Error("mkdir stderr error:", err)
			return err
		}
		_, err = os.Create(c.Stderr)
		if err != nil {
			logger.Error("create stdout error:", err)
			return err
		}
	}
	stderr, err := os.OpenFile(c.Stderr, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		logger.Error("open stderr error:", err)
		return err
	}

	c.Cmd.Stderr = stderr
	return nil
}

func (c *Command) start() error {
	err := c.Cmd.Start()
	if err != nil {
		logger.Error(c.Name+" start fail:", err)
		c.Finished = true
		return err
	}
	c.Begin = time.Now()
	c.Finished = false
	c.UUID = uuid.GenerateV1()
	c.Pid = c.Cmd.Process.Pid
	logger.Info(c.Name+" started at ", c.Pid)
	if c.IsMonitor {
		MoniterAdd(c.Pid, c.monitor)
	}
	return nil
}

func (c *Command) wait(callback func()) {
	_ = c.Cmd.Wait()
	if c.IsMonitor {
		MoniterRemove(c.Pid)
	}
	status := c.Cmd.ProcessState.Sys().(syscall.WaitStatus)
	signaled := status.Signaled()
	signal := status.Signal()
	if signaled {
		logger.Info(c.Name+" signaled:", signal.String())
	}
	c.End = time.Now()
	c.Finished = true
	if c.Cmd.ProcessState.ExitCode() != 0 {
		logger.Error(c.Name+" exit with status ", c.Cmd.ProcessState.ExitCode())
	} else {
		logger.Info(c.Name + " finished")
	}
	c.Status = c.Cmd.ProcessState.ExitCode()
	if c.ArchiveReport != nil {
		c.ArchiveReport(c)
	}
	callback()
}

func (c *Command) stop() {
	if !c.Finished {
		if c.Cmd == nil {
			return
		}
		if c.Cmd.Process == nil {
			return
		}
		err := c.Cmd.Process.Kill()
		if err != nil {
			logger.Error(c.Name+" kill fail:", err)
		}
		logger.Info(c.Name + " killed!")
		c.Status = -2
		if c.ArchiveReport != nil {
			c.ArchiveReport(c)
		}
	}
}

func (c *Command) monitor(info *process.Process) {
	monitor := new(Monitor)
	memoryPercent, err := info.MemoryPercent()
	if err == nil {
		monitor.MemoryPercent = memoryPercent
	}
	cpuPercent, err := info.CPUPercent()
	if err == nil {
		monitor.CPUPercent = cpuPercent
	}
	threads, err := info.NumThreads()
	if err == nil {
		monitor.NumThreads = int(threads)
	}
	if c.MonitorReport != nil {
		c.MonitorReport(monitor)
	}
}
