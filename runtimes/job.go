package runtimes

import (
	"fmt"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/robfig/cron/v3"
)

type Job struct {
	Command
	ID      int64
	Spec    string
	TimeOut time.Duration
	CronID  cron.EntryID
}

func (j *Job) Run() {
	err := j.build()
	if err != nil {
		return
	}
	stop := make(chan bool, 1)
	if j.TimeOut > 0 {
		go j.timer(stop)
	}
	err = j.start()
	if err != nil {
		stop <- true
		return
	}
	j.wait(j.record)
	if j.TimeOut > 0 {
		stop <- true
	}
}

func (j *Job) timer(ch chan bool) {
	timer := time.NewTimer(time.Second * j.TimeOut)
	for {
		select {
		case <-timer.C:
			err := j.Cmd.Process.Kill()
			if err != nil {
				logger.Error("job Kill Error:", err)
			}
		case <-ch:
			timer.Stop()
		}
	}
}

func (j *Job) record() {
	info := fmt.Sprintf("%s finished with %.2f seconds", j.Name, j.End.Sub(j.Begin).Seconds())
	logger.Info(info)
}
