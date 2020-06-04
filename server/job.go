package server

import (
	"fmt"

	"Asgard/applications"
	"Asgard/providers"
	"Asgard/rpc"
)

func AddJob(id int64, request *rpc.Job) error {
	_, ok := applications.Jobs[id]
	if ok {
		return nil
	}
	err := applications.JobRegister(
		id,
		rpc.BuildJobConfig(request),
		providers.MonitorMamager,
		providers.MasterClient.Reports,
		providers.MasterClient.JobMonitorChan,
		providers.MasterClient.JobArchiveChan,
	)
	if err != nil {
		return err
	}
	applications.JobAdd(applications.Jobs[id])
	return nil
}

func UpdateJob(id int64, job *applications.Job, request *rpc.Job) error {
	err := DeleteJob(id, job)
	if err != nil {
		return err
	}
	return AddJob(id, request)
}

func DeleteJob(id int64, job *applications.Job) error {
	_, ok := applications.Jobs[id]
	if !ok {
		return nil
	}
	ok = applications.JobStop(id)
	if !ok {
		return fmt.Errorf("app %d stop failed", id)
	}
	applications.JobUnRegister(id)
	return nil
}
