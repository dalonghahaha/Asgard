package applications

import (
	"fmt"
	"sync"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"

	"Asgard/constants"
)

var (
	Timings    = map[int64]*Timing{}
	timingLock sync.Mutex
)

func TimingStopAll() {
	if constants.SYSTEM_TIMER_TICKER != nil {
		constants.SYSTEM_TIMER_TICKER.Stop()
	}
	for _, timing := range Timings {
		if timing.Running {
			timing.stop()
		}
	}
}

func TimingStartAll(moniter bool) {
	if moniter {
		go MoniterStart()
	}
	go TimingRun()
}

func TimingRun() {
	constants.SYSTEM_TIMER_TICKER = time.NewTicker(time.Second * time.Duration(constants.SYSTEM_TIMER))
	for range constants.SYSTEM_TIMER_TICKER.C {
		now := time.Now().Unix()
		for _, timing := range Timings {
			if timing.Time.Unix() < now && !timing.Executed {
				go timing.Run()
			}
		}
	}
}

func TimingStop(id int64) bool {
	timing, ok := Timings[id]
	if !ok {
		return false
	}
	if timing.Running {
		timing.stop()
	}
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

func TimingRegister(id int64, config map[string]interface{}, reports *sync.Map, mc chan TimingMonitor, ac chan TimingArchive) error {
	timing, err := NewTiming(config)
	if err != nil {
		return err
	}
	timing.ID = id
	timing.MonitorReport = func(monitor *Monitor) {
		timingMonitor := TimingMonitor{
			UUID:    uuid.GenerateV4(),
			Timing:  timing,
			Monitor: monitor,
		}
		if reports != nil {
			reports.Store(timingMonitor.UUID, 1)
		}
		if mc != nil {
			mc <- timingMonitor
		}
	}
	timing.ArchiveReport = func(archive *Archive) {
		timingArchive := TimingArchive{
			UUID:    uuid.GenerateV4(),
			Timing:  timing,
			Archive: archive,
		}
		if reports != nil {
			reports.Store(timingArchive.UUID, 1)
		}
		if ac != nil {
			ac <- timingArchive
		}
	}
	timingLock.Lock()
	Timings[id] = timing
	timingLock.Unlock()
	return nil
}

func TimingUnRegister(id int64) {
	timingLock.Lock()
	delete(Timings, id)
	timingLock.Unlock()
}
