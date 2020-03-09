package applications

import (
	"fmt"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/spf13/viper"
)

var Timings = map[int64]*Timing{}

func TimingStopAll() {
	MoniterStop()
	processExit = true
	for _, timing := range Timings {
		if !timing.Finished {
			timing.stop()
		}
	}
}

func TimingStartAll(moniter bool) {
	if moniter {
		go MoniterStart()
	}
	duration := viper.GetInt("system.timer")
	ticker = time.NewTicker(time.Second * time.Duration(duration))
	for range ticker.C {
		now := time.Now().Unix()
		for _, timing := range Timings {
			if timing.Time.Unix() < now && !timing.Executed {
				go timing.Run()
			}
		}
	}
}

func TimingAdd(id int64, timing *Timing) {
	Timings[id] = timing
}

func TimingStart(name string) bool {
	for _, timing := range Timings {
		if timing.Name == name {
			go timing.Run()
			return true
		}
	}
	return false
}

func TimingStartByID(id int64) bool {
	timing, ok := Timings[id]
	if !ok {
		return false
	}
	go timing.Run()
	return true
}

func TimingStop(name string) bool {
	for _, timing := range Timings {
		if timing.Name == name {
			timing.stop()
			return true
		}
	}
	return false
}

func TimingStopByID(id int64) bool {
	timing, ok := Timings[id]
	if !ok {
		return false
	}
	timing.stop()
	return true
}

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
	info := fmt.Sprintf("%s executed with %.2f seconds", t.Name, t.End.Sub(t.Begin).Seconds())
	logger.Info(info)
	delete(Timings, t.ID)
}

func NewTiming(config map[string]interface{}) (*Timing, error) {
	timing := new(Timing)
	err := timing.configure(config)
	if err != nil {
		return nil, err
	}
	_time, ok := config["time"].(int64)
	if !ok {
		return nil, fmt.Errorf("config timeout type wrong")
	}
	timing.Time = time.Unix(_time, 0)
	timeout, ok := config["timeout"].(int64)
	if !ok {
		return nil, fmt.Errorf("config timeout type wrong")
	}
	timing.TimeOut = time.Duration(timeout)
	return timing, nil
}

func TimingRegister(id int64, config map[string]interface{}) (*Timing, error) {
	timing, err := NewTiming(config)
	if err != nil {
		return nil, err
	}
	timing.ID = id
	Timings[id] = timing
	return timing, nil
}
