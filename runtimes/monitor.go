package runtimes

import (
	"math"
	"sync"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/shirou/gopsutil/process"

	"Asgard/constants"
)

var (
	cupUnit    = 0.001
	memoryUnit = 0.0001
)

type MonitorInfo struct {
	CPUPercent float64
	Memory     float64
	NumThreads int
}

type AgentMonitor struct {
	UUID    string
	Ip      string
	Port    string
	Monitor *MonitorInfo
}

type AppMonitor struct {
	UUID    string
	App     *App
	Monitor *MonitorInfo
}

type JobMonitor struct {
	UUID    string
	Job     *Job
	Monitor *MonitorInfo
}

type TimingMonitor struct {
	UUID    string
	Timing  *Timing
	Monitor *MonitorInfo
}

func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

func BuildMonitorInfo(info *process.Process) *MonitorInfo {
	monitor := new(MonitorInfo)
	memoryInfo, err := info.MemoryInfo()
	if err == nil {
		monitor.Memory = math.Trunc(bytesToMB(memoryInfo.RSS)/memoryUnit) * memoryUnit
	}
	cpuPercent, err := info.CPUPercent()
	if err == nil {
		monitor.CPUPercent = math.Trunc(cpuPercent/cupUnit) * cupUnit
	}
	threads, err := info.NumThreads()
	if err == nil {
		monitor.NumThreads = int(threads)
	}
	return monitor
}

type Monitor struct {
	name       string
	lock       sync.Mutex
	ticker     *time.Ticker
	exitSingel chan bool
	monitors   map[int]func(monitor *MonitorInfo)
}

func NewMonitor(name string) *Monitor {
	return &Monitor{
		name:       name,
		ticker:     time.NewTicker(time.Second * time.Duration(constants.SYSTEM_MONITER)),
		exitSingel: make(chan bool, 1),
		monitors:   make(map[int]func(monitor *MonitorInfo)),
	}
}

func (m *Monitor) Add(pid int, callback func(monitor *MonitorInfo)) {
	m.lock.Lock()
	m.monitors[pid] = callback
	m.lock.Unlock()
}

func (m *Monitor) Remove(pid int) {
	m.lock.Lock()
	delete(m.monitors, pid)
	m.lock.Unlock()
}

func (m *Monitor) Start() {
	logger.Debugf("%s monitor ticker start!", m.name)
	SubscribeExit(m.exitSingel)
	for {
		select {
		case <-m.exitSingel:
			logger.Debugf("%s monitor ticker stop!", m.name)
			m.ticker.Stop()
			break
		case <-m.ticker.C:
			for pid, function := range m.monitors {
				info, err := process.NewProcess(int32(pid))
				if err != nil {
					continue
				}
				if function != nil {
					function(BuildMonitorInfo(info))
				}
			}

		}
	}
}
