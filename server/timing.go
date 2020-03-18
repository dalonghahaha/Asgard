package server

import (
	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
	"fmt"
)

func AddTiming(id int64, request *rpc.Timing) error {
	timing, err := applications.TimingRegister(id, rpc.BuildTimingConfig(request))
	if err != nil {
		return err
	}
	timing.MonitorReport = func(monitor *applications.Monitor) {
		client.TimingMonitorReport(rpc.BuildTimingMonior(timing, monitor))
	}
	timing.ArchiveReport = func(command *applications.Command) {
		client.TimingArchiveReport(rpc.BuildTimingArchive(timing, command))
	}
	applications.TimingAdd(id, timing)
	return nil
}

func UpdateTiming(id int64, timing *applications.Timing, request *rpc.Timing) error {
	err := DeleteTiming(id, timing)
	if err != nil {
		return err
	}
	return AddTiming(id, request)
}

func DeleteTiming(id int64, timing *applications.Timing) error {
	ok := applications.TimingStopByID(id)
	if !ok {
		return fmt.Errorf("app %d stop failed", id)
	}
	delete(applications.Timings, id)
	return nil
}

func GetTimingOutLog(id int64) []string {
	if timing, ok := applications.Timings[id]; ok {
		return timing.GetOutLog()
	}
	return []string{"无记录"}
}

func GetTimingErrLog(id int64) []string {
	if timing, ok := applications.Timings[id]; ok {
		return timing.GetErrLog()
	}
	return []string{"无记录"}
}
