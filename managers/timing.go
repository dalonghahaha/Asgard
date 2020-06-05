package managers

import (
	"fmt"
	"sync"
	"time"

	"Asgard/clients"
	"Asgard/constants"
	"Asgard/runtimes"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
)

type TimingManager struct {
	lock         sync.Mutex
	ticker       *time.Ticker
	exitSingel   chan bool
	timings      map[int64]*runtimes.Timing
	masterClient *clients.Master
	monitor      *runtimes.Monitor
}

func NewTimingManager() *TimingManager {
	manager := &TimingManager{
		ticker:     time.NewTicker(time.Second * time.Duration(constants.SYSTEM_TIMER)),
		exitSingel: make(chan bool, 1),
		timings:    make(map[int64]*runtimes.Timing),
		monitor:    runtimes.NewMonitor("timing"),
	}
	return manager
}

func (m *TimingManager) SetMaster(masterClient *clients.Master) {
	m.masterClient = masterClient
}

func (m *TimingManager) StartMonitor() {
	go m.monitor.Start()
}

func (m *TimingManager) NewTiming(config map[string]interface{}) (*runtimes.Timing, error) {
	timing := new(runtimes.Timing)
	err := timing.Configure(config)
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

func (m *TimingManager) Register(id int64, config map[string]interface{}) error {
	timing, err := m.NewTiming(config)
	if err != nil {
		return err
	}
	timing.ID = id
	timing.Monitor = m.monitor
	timing.ExceptionReport = func(message string) {
		logger.Infof("%s ExceptionReport", timing.Name)
		timingException := runtimes.TimingException{
			UUID:     uuid.GenerateV4(),
			TimingID: timing.ID,
			Desc:     message,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(timingException.UUID, 1)
			m.masterClient.TimingExceptionChan <- timingException
		}
	}
	timing.MonitorReport = func(monitor *runtimes.MonitorInfo) {
		timingMonitor := runtimes.TimingMonitor{
			UUID:    uuid.GenerateV4(),
			Timing:  timing,
			Monitor: monitor,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(timingMonitor.UUID, 1)
			m.masterClient.TimingMonitorChan <- timingMonitor
		}
	}
	timing.ArchiveReport = func(archive *runtimes.Archive) {
		timingArchive := runtimes.TimingArchive{
			UUID:    uuid.GenerateV4(),
			Timing:  timing,
			Archive: archive,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(timingArchive.UUID, 1)
			m.masterClient.TimingArchiveChan <- timingArchive
		}
	}
	m.lock.Lock()
	m.timings[id] = timing
	m.lock.Unlock()
	return nil
}

func (m *TimingManager) UnRegister(id int64) {
	m.lock.Lock()
	delete(m.timings, id)
	m.lock.Unlock()
}

func (m *TimingManager) Count() int {
	return len(m.timings)
}

func (m *TimingManager) GetList() []*runtimes.Timing {
	list := []*runtimes.Timing{}
	for _, timing := range m.timings {
		list = append(list, timing)
	}
	return list
}

func (m *TimingManager) Get(id int64) *runtimes.Timing {
	app, ok := m.timings[id]
	if !ok {
		return nil
	}
	return app
}

func (m *TimingManager) GetByName(name string) *runtimes.Timing {
	for _, timing := range m.timings {
		if timing.Name == name {
			return timing
		}
	}
	return nil
}

func (m *TimingManager) StartAll(monitor bool) {
	if monitor {
		m.StartMonitor()
	}
	go m.Run()
}

func (m *TimingManager) Run() {
	logger.Debug("timing manager ticker start!")
	runtimes.SubscribeExit(m.exitSingel)
	for {
		select {
		case <-m.exitSingel:
			logger.Debug("timing manager ticker stop!")
			m.ticker.Stop()
			break
		case <-m.ticker.C:
			now := time.Now().Unix()
			for _, timing := range m.timings {
				if timing.Executed {
					m.UnRegister(timing.ID)
				}
				if timing.Time.Unix() < now && !timing.Executed {
					go timing.Run()
				}
			}
		}
	}
}

func (m *TimingManager) StopAll() {
	for _, timing := range m.timings {
		if timing.Running {
			timing.Kill()
		}
	}
}

func (m *TimingManager) Stop(id int64) bool {
	timing := m.Get(id)
	if timing == nil {
		return true
	}
	if timing.Running {
		timing.Kill()
	}
	return true
}

func (m *TimingManager) Update(id int64, config map[string]interface{}) error {
	err := m.Remove(id)
	if err != nil {
		return err
	}
	return m.Register(id, config)
}

func (m *TimingManager) Remove(id int64) error {
	ok := m.Stop(id)
	if !ok {
		return fmt.Errorf("timing stop failed!")
	}
	m.UnRegister(id)
	return nil
}
