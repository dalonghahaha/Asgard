package providers

import (
	"Asgard/runtimes"
)

var (
	MonitorMamager *runtimes.Monitor
)

func RegisterMonitorMamager() {
	MonitorMamager = runtimes.NewMonitor()
	go MonitorMamager.Start()
}
