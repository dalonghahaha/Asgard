package structs

import (
	"os/exec"
	"time"
)

type App struct {
	Name     string
	Dir      string
	Program  string
	Args     string
	Stdout   string
	Stderr   string
	Pid      int
	Cmd      *exec.Cmd
	Finished bool
	Begin    time.Time
	End      time.Time
}
