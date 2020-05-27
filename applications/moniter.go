package applications

import (
	"math"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
	"github.com/spf13/viper"
)

var (
	unit   = 0.0001
	lock   sync.Mutex
	ticker *time.Ticker
)

var moniters = map[int]func(info *process.Process){}

type Monitor struct {
	CPUPercent float64
	Memory     float64
	NumThreads int
}

func bytesToMB(bytes uint64) float64 {
	mb := float64(bytes) / 1024 / 1024
	return math.Trunc(mb/unit) * unit
}

func BuildMonitor(info *process.Process) *Monitor {
	monitor := new(Monitor)
	memoryPercent, err := info.MemoryInfo()
	if err == nil {
		monitor.Memory = bytesToMB(memoryPercent.RSS)
	}
	cpuPercent, err := info.CPUPercent()
	if err == nil {
		monitor.CPUPercent = cpuPercent
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
