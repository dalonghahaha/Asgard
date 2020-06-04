package providers

import (
	"Asgard/applications"
)

var (
	MonitorMamager *applications.MonitorMamager
)

func RegisterMonitorMamager() {
	MonitorMamager = applications.NewMonitorMamager()
	go MonitorMamager.Start()
}
