package server

import (
	"Asgard/applications"
	"Asgard/providers"
	"Asgard/rpc"
	"fmt"
)

func AddTiming(id int64, request *rpc.Timing) error {
	err := applications.TimingRegister(
		id,
		rpc.BuildTimingConfig(request),
		providers.MasterClient.TimingMonitorChan,
		providers.MasterClient.TimingArchiveChan,
	)
	if err != nil {
		return err
	}
	applications.TimingAdd(id, applications.Timings[id])
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
