package runtimes

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/file"
	"github.com/dalonghahaha/avenger/tools/uuid"
)

var (
	processExit = false
)

type Command struct {
	lock            sync.Mutex
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
	Running         bool
	Successed       bool
	Status          int
	Signal          string
	Cmd             *exec.Cmd
	Monitor         *Monitor
	ExceptionReport func(message string)
	MonitorReport   func(monitor *MonitorInfo)
	ArchiveReport   func(archive *Archive)
}

func Exit() {
	processExit = true
	logger.Info("exit!")
	time.Sleep(time.Millisecond * 100)
}

func (c *Command) Configure(config map[string]interface{}) error {
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
	args := strings.Split(c.Args, " ")
	c.Cmd = exec.Command(c.Program, args...)
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
	c.UUID = uuid.GenerateV4()
	c.Begin = time.Now()
	err := c.Cmd.Start()
	if err != nil {
		logger.Errorf("%s start fail: %s", c.Name, err)
		c.finish()
		return err
	}
	c.runing()
	c.Pid = c.Cmd.Process.Pid
	logger.Infof("%s started at %d", c.Name, c.Pid)
	if c.IsMonitor && c.Monitor != nil {
		c.Monitor.Add(c.Pid, c.MonitorReport)
	}
	return nil
}

func (c *Command) wait(callback func()) {
	_ = c.Cmd.Wait()
	if !c.Running {
		return
	}
	c.finish()
	if c.Cmd == nil || c.Cmd.ProcessState == nil {
		c.Status = -3
		c.Signal = "unknow"
	} else {
		c.Successed = c.Cmd.ProcessState.Success()
		c.Status = c.Cmd.ProcessState.ExitCode()
		status, ok := c.Cmd.ProcessState.Sys().(syscall.WaitStatus)
		if ok && status.Signaled() {
			c.Signal = status.Signal().String()
		}
	}
	if c.ArchiveReport != nil {
		c.ArchiveReport(buildArchive(c))
	}
	callback()
}

func (c *Command) Kill() {
	if !c.Running {
		return
	}
	c.finish()
	if c.Cmd == nil || c.Cmd.Process == nil {
		c.Status = -3
		c.Signal = "unknow"
	} else {
		//try kill
		err := c.Cmd.Process.Kill()
		if err == nil {
			c.Status = -2
			c.Signal = "kill"
		} else {
			fmt.Println(err)
			if c.Cmd.ProcessState != nil {
				c.Status = c.Cmd.ProcessState.ExitCode()
				status, ok := c.Cmd.ProcessState.Sys().(syscall.WaitStatus)
				if ok && status.Signaled() {
					c.Signal = status.Signal().String()
				}
			}
		}
	}
	if c.ArchiveReport != nil {
		c.ArchiveReport(buildArchive(c))
	}
}

func (c *Command) finish() {
	c.lock.Lock()
	logger.Info(c.Name + " finish")
	if c.IsMonitor && c.Monitor != nil {
		c.Monitor.Remove(c.Pid)
	}
	c.End = time.Now()
	c.Running = false
	c.lock.Unlock()
}

func (c *Command) runing() {
	c.lock.Lock()
	c.Running = true
	c.lock.Unlock()
}
