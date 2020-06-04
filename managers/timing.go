package managers

import (
	"fmt"
	"sync"
	"time"

	"Asgard/applications"
	"Asgard/client"
	"Asgard/constants"

	"github.com/dalonghahaha/avenger/tools/uuid"
)

type TimingManager struct {
	lock           sync.Mutex
	timings        map[int64]*applications.Timing
	masterClient   *client.Master
	monitorManager *applications.MonitorMamager
}

func NewTimingManager() (*TimingManager, error) {
	manager := &TimingManager{
		timings:        make(map[int64]*applications.Timing),
		monitorManager: applications.NewMonitorMamager(),
	}
	return manager, nil
}

func (m *TimingManager) SetMaster(masterClient *client.Master) {
	m.masterClient = masterClient
}

func (m *TimingManager) StartMonitor() {
	m.monitorManager.Start()
}

func (m *TimingManager) StopMonitor() {
	m.monitorManager.Stop()
}

func (m *TimingManager) Count() int {
	return len(m.timings)
}

func (m *TimingManager) GetTiming(id int64) *applications.Timing {
	app, ok := m.timings[id]
	if !ok {
		return nil
	}
	return app
}

func (m *TimingManager) StartAll(moniter bool) {
	go m.Run()
}

func (m *TimingManager) Run() {
	constants.SYSTEM_TIMER_TICKER = time.NewTicker(time.Second * time.Duration(constants.SYSTEM_TIMER))
	for range constants.SYSTEM_TIMER_TICKER.C {
		now := time.Now().Unix()
		for _, timing := range m.timings {
			if timing.Time.Unix() < now && !timing.Executed {
				go timing.Run()
			}
		}
	}
}

func (m *TimingManager) StopAll() {
	if constants.SYSTEM_TIMER_TICKER != nil {
		constants.SYSTEM_TIMER_TICKER.Stop()
	}
	for _, timing := range m.timings {
		if timing.Running {
			timing.Kill()
		}
	}
}

func (m *TimingManager) Remove(id int64) bool {
	timing := m.GetTiming(id)
	if timing == nil {
		return true
	}
	if timing.Running {
		timing.Kill()
	}
	return true
}

func (m *TimingManager) NewTiming(config map[string]interface{}) (*applications.Timing, error) {
	timing := new(applications.Timing)
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

func (m *TimingManager) TimingRegister(id int64, config map[string]interface{}) error {
	timing, err := m.NewTiming(config)
	if err != nil {
		return err
	}
	timing.ID = id
	timing.MonitorMamager = m.monitorManager
	timing.MonitorReport = func(monitor *applications.Monitor) {
		timingMonitor := applications.TimingMonitor{
			UUID:    uuid.GenerateV4(),
			Timing:  timing,
			Monitor: monitor,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(timingMonitor.UUID, 1)
			m.masterClient.TimingMonitorChan <- timingMonitor
		}
	}
	timing.ArchiveReport = func(archive *applications.Archive) {
		timingArchive := applications.TimingArchive{
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

func (m *TimingManager) TimingUnRegister(id int64) {
	m.lock.Lock()
	delete(m.timings, id)
	m.lock.Unlock()
}
