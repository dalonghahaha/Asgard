package runtimes

import (
	"fmt"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
)

type Timing struct {
	Command
	ID       int64
	Time     time.Time
	TimeOut  time.Duration
	Executed bool
}

func (t *Timing) Run() {
	err := t.build()
	if err != nil {
		return
	}
	stop := make(chan bool, 1)
	if t.TimeOut > 0 {
		go t.timer(stop)
	}
	t.Executed = true
	err = t.start()
	if err != nil {
		logger.Errorf("%s start fail: %s", t.Name, err)
		if t.ExceptionReport != nil {
			t.ExceptionReport(fmt.Sprintf("start fail: %+v", err))
		}
		stop <- true
		return
	}
	t.wait(t.record)
	if t.TimeOut > 0 {
		stop <- true
	}
}

func (j *Timing) timer(ch chan bool) {
	timer := time.NewTimer(time.Second * j.TimeOut)
	for {
		select {
		case <-timer.C:
			err := j.Cmd.Process.Kill()
			if err != nil {
				logger.Error("app Kill Error:", err)
			}
		case <-ch:
			timer.Stop()
		}
	}
}

func (t *Timing) record() {
	logger.Infof("%s finished with %.2f seconds", t.Name, t.End.Sub(t.Begin).Seconds())
	if !t.Successed {
		if t.ExceptionReport != nil {
			t.ExceptionReport(fmt.Sprintf("finished unseccessed with status: %d", t.Status))
		}
	}
}
