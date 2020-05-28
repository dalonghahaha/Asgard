package applications

import (
	"math"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/spf13/viper"
)

var (
	cupUnit    = 0.001
	memoryUnit = 0.0001
	lock       sync.Mutex
	ticker     *time.Ticker
)

var moniters = map[int]func(info *process.Process){}

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

func MoniterAdd(pid int, callback func(info *process.Process)) {
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
	duration := viper.GetInt("system.moniter")
	ticker = time.NewTicker(time.Second * time.Duration(duration))
	for range ticker.C {
		for pid, function := range moniters {
			info, err := process.NewProcess(int32(pid))
			if err != nil {
				continue
			}
			function(info)
		}
	}
}

func MoniterStop() {
	if ticker != nil {
		ticker.Stop()
	}
}
