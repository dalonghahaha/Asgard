package server

import (
	"Asgard/applications"
	"Asgard/client"
	"Asgard/rpc"
	"fmt"
	"time"
)

func AddJob(id int64, request *rpc.Job) error {
	job, err := applications.JobRegister(id, rpc.BuildJobConfig(request))
	if err != nil {
		return err
	}
	job.MonitorReport = func(monitor *applications.Monitor) {
		client.JobMonitorReport(job, monitor)
	}
	job.ArchiveReport = func(command *applications.Command) {
		client.JobArchiveReport(job, command)
	}
	ok := applications.JobStartByID(id)
	if !ok {
		return fmt.Errorf("job %d start failed", id)
	}
	return nil
}

func UpdateJob(id int64, job *applications.Job, request *rpc.Job) error {
	ok := applications.AppStopByID(id)
	if !ok {
		return fmt.Errorf("job %d stop failed", id)
	}
	job.Name = request.GetName()
	job.Dir = request.GetDir()
	job.Program = request.GetProgram()
	job.Args = request.GetArgs()
	job.Stdout = request.GetStdOut()
	job.Stderr = request.GetStdErr()
	job.Spec = request.GetSpec()
	job.TimeOut = time.Duration(request.GetTimeout())
	job.IsMonitor = request.GetIsMonitor()
	ok = applications.JobStartByID(id)
	if !ok {
		return fmt.Errorf("job %d start failed", id)
	}
	return nil
}
