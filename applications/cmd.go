package applications

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
	"github.com/shirou/gopsutil/process"
)

var processExit = false

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
	Finished        bool
	Status          int
	Signal          string
	Cmd             *exec.Cmd
	ExceptionReport func(message string)
	MonitorReport   func(monitor *Monitor)
	ArchiveReport   func(archive *Archive)
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
	err := c.Cmd.Start()
	if err != nil {
		logger.Error(c.Name+" start fail:", err)
		c.Finished = true
		return err
	}
	c.Begin = time.Now()
	c.Finished = false
	c.UUID = uuid.GenerateV4()
	c.Pid = c.Cmd.Process.Pid
	logger.Info(c.Name+" started at ", c.Pid)
	if c.IsMonitor {
		MoniterAdd(c.Pid, c.monitor)
	}
	return nil
}

func (c *Command) wait(callback func()) {
	_ = c.Cmd.Wait()
	if c.Finished {
		return
	}
	c.finish()
	if c.Cmd == nil || c.Cmd.ProcessState == nil {
		c.Status = -3
		c.Signal = "unknow"
	} else {
		c.Status = c.Cmd.ProcessState.ExitCode()
		status, ok := c.Cmd.ProcessState.Sys().(syscall.WaitStatus)
		if ok && status.Signaled() {
			c.Signal = status.Signal().String()
		}
	}
	if c.ArchiveReport != nil {
		logger.Debug(fmt.Sprintf("appArchive Send from wait:[%s][%d][%s]", c.Name, c.Status, c.Signal))
		c.ArchiveReport(buildArchive(c))
	}
	callback()
}

func (c *Command) stop() {
	if c.Finished {
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
		logger.Debug(fmt.Sprintf("appArchive Send from stop:[%s][%d][%s]", c.Name, c.Status, c.Signal))
		c.ArchiveReport(buildArchive(c))
	}
}

func (c *Command) finish() {
	c.lock.Lock()
	if c.IsMonitor {
		MoniterRemove(c.Pid)
	}
	c.End = time.Now()
	c.Finished = true
	c.lock.Unlock()
}

func (c *Command) monitor(info *process.Process) {
	monitor := BuildMonitor(info)
	if c.MonitorReport != nil {
		c.MonitorReport(monitor)
	}
}
