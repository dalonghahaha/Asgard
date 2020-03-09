package server

import (
	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
	"fmt"
)

func AddJob(id int64, request *rpc.Job) error {
	job, err := applications.JobRegister(id, rpc.BuildJobConfig(request))
	if err != nil {
		return err
	}
	job.MonitorReport = func(monitor *applications.Monitor) {
		client.JobMonitorReport(rpc.BuildJobMonior(job, monitor))
	}
	job.ArchiveReport = func(command *applications.Command) {
		client.JobArchiveReport(rpc.BuildJobArchive(job, command))
	}
	applications.JobAdd(job)
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
