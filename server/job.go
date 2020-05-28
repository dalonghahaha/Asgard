package server

import (
	"Asgard/applications"
	"Asgard/providers"
	"Asgard/rpc"
	"fmt"
)

func AddJob(id int64, request *rpc.Job) error {
	err := applications.JobRegister(
		id,
		rpc.BuildJobConfig(request),
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
	ok := applications.JobStopByID(id)
	if !ok {
		return fmt.Errorf("app %d stop failed", id)
	}
	delete(applications.Jobs, id)
	return nil
}
