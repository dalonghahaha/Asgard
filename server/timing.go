package server

import (
	"Asgard/applications"
	"Asgard/providers"
	"Asgard/rpc"
	"fmt"
)

func AddTiming(id int64, request *rpc.Timing) error {
	_, ok := applications.Timings[id]
	if ok {
		return nil
	}
	err := applications.TimingRegister(
		id,
		rpc.BuildTimingConfig(request),
		providers.MasterClient.Reports,
		providers.MasterClient.TimingMonitorChan,
		providers.MasterClient.TimingArchiveChan,
	)
	if err != nil {
		return err
	}
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
	_, ok := applications.Timings[id]
	if !ok {
		return nil
	}
	ok = applications.TimingStop(id)
	if !ok {
		return fmt.Errorf("app %d stop failed", id)
	}
	applications.TimingUnRegister(id)
	return nil
}
