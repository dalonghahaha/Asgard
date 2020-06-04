package applications

import (
	"math"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"

	"Asgard/constants"
)

var (
	cupUnit    = 0.001
	memoryUnit = 0.0001
)

var moniters = map[int]func(monitor *Monitor){}

type Monitor struct {
	CPUPercent float64
	Memory     float64
	NumThreads int
}

type AgentMonitor struct {
	UUID    string
	Ip      string
	Port    string
	Monitor *Monitor
}

type AppMonitor struct {
	UUID    string
	App     *App
	Monitor *Monitor
}

type JobMonitor struct {
	UUID    string
	Job     *Job
	Monitor *Monitor
}

type TimingMonitor struct {
	UUID    string
	Timing  *Timing
	Monitor *Monitor
}

func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

func BuildMonitor(info *process.Process) *Monitor {
	monitor := new(Monitor)
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

func MoniterAdd(pid int, callback func(monitor *Monitor)) {
	lock.Lock()
	moniters[pid] = callback
	lock.Unlock()
}

func MoniterRemove(pid int) {
	lock.Lock()
	delete(moniters, pid)
	lock.Unlock()
}

func MoniterStart() {
	constants.SYSTEM_MONITER_TICKER = time.NewTicker(time.Second * time.Duration(constants.SYSTEM_MONITER))
	for range constants.SYSTEM_MONITER_TICKER.C {
		for pid, function := range moniters {
			info, err := process.NewProcess(int32(pid))
			if err != nil {
				continue
			}
			function(BuildMonitor(info))
		}
	}
}

func MoniterStop() {
	if constants.SYSTEM_MONITER_TICKER != nil {
		constants.SYSTEM_MONITER_TICKER.Stop()
	}
}

type MonitorMamager struct {
	lock     sync.Mutex
	moniters map[int]func(monitor *Monitor)
}

func NewMonitorMamager() *MonitorMamager {
	return &MonitorMamager{
		moniters: make(map[int]func(monitor *Monitor)),
	}
}

func (m *MonitorMamager) Add(pid int, callback func(monitor *Monitor)) {
	m.lock.Lock()
	m.moniters[pid] = callback
	m.lock.Unlock()
}

func (m *MonitorMamager) Remove(pid int) {
	m.lock.Lock()
	delete(m.moniters, pid)
	m.lock.Unlock()
}

func (m *MonitorMamager) Start() {
	constants.SYSTEM_MONITER_TICKER = time.NewTicker(time.Second * time.Duration(constants.SYSTEM_MONITER))
	for range constants.SYSTEM_MONITER_TICKER.C {
		for pid, function := range moniters {
			info, err := process.NewProcess(int32(pid))
			if err != nil {
				continue
			}
			if function != nil {
				function(BuildMonitor(info))
			}
		}
	}
}

func (m *MonitorMamager) Stop() {
	if constants.SYSTEM_MONITER_TICKER != nil {
		constants.SYSTEM_MONITER_TICKER.Stop()
	}
}
