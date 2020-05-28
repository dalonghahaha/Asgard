package applications

import (
	"fmt"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/robfig/cron/v3"
)

var crontab *cron.Cron

var Jobs = map[int64]*Job{}

func JobStopAll() {
	MoniterStop()
	processExit = true
	for _, job := range Jobs {
		if !job.Finished {
			job.stop()
		}
	}
}

func JobStartAll(moniter bool) {
	crontab = cron.New()
	for _, job := range Jobs {
		JobAdd(job)
	}
	crontab.Start()
	if moniter {
		MoniterStart()
	}
}

func JobAdd(job *Job) {
	id, err := crontab.AddJob(job.Spec, job)
	if err != nil {
		logger.Error(job.Name+" add fail:", err)
	}
	logger.Info(job.Name + " add seccess!")
	job.CronID = id
}

func JobStart(name string) bool {
	for _, job := range Jobs {
		if job.Name == name {
			go job.Run()
			return true
		}
	}
	return false
}

func JobStartByID(id int64) bool {
	job, ok := Jobs[id]
	if !ok {
		return false
	}
	go job.Run()
	return true
}

func JobStop(name string) error {
	for _, job := range Jobs {
		if job.Name == name {
			job.stop()
			crontab.Remove(job.CronID)
		}
	}
	return nil
}

func JobStopByID(id int64) bool {
	job, ok := Jobs[id]
	if !ok {
		return false
	}
	job.stop()
	crontab.Remove(job.CronID)
	return true
}

type Job struct {
	Command
	ID      int64
	Spec    string
	TimeOut time.Duration
	CronID  cron.EntryID
}

func (j *Job) Run() {
	err := j.build()
	if err != nil {
		return
	}
	stop := make(chan bool, 1)
	if j.TimeOut > 0 {
		go j.timer(stop)
	}
	err = j.start()
	if err != nil {
		stop <- true
		return
	}
	j.wait(j.record)
	if j.TimeOut > 0 {
		stop <- true
	}
}

func (j *Job) timer(ch chan bool) {
	timer := time.NewTimer(time.Second * j.TimeOut)
	for {
		select {
		case <-timer.C:
			err := j.Cmd.Process.Kill()
			if err != nil {
				logger.Error("app Kill Error:", err)
			}
		case <-ch:
			timer.Stop()
		}
	}
}

func (j *Job) record() {
	info := fmt.Sprintf("%s finished with %.2f seconds", j.Name, j.End.Sub(j.Begin).Seconds())
	logger.Info(info)
}

func NewJob(config map[string]interface{}) (*Job, error) {
	job := new(Job)
	err := job.configure(config)
	if err != nil {
		return nil, err
	}
	spec, ok := config["spec"].(string)
	if !ok {
		return nil, fmt.Errorf("config spec type wrong")
	}
	job.Spec = spec
	timeout, ok := config["timeout"].(int64)
	if !ok {
		return nil, fmt.Errorf("config timeout type wrong")
	}
	job.TimeOut = time.Duration(timeout)
	return job, nil
}

func JobAppend(id int64, config map[string]interface{}) (*Job, error) {
	job, err := NewJob(config)
	if err != nil {
		return nil, err
	}
	Jobs[id] = job
	JobAdd(job)
	return job, nil
}

func JobRegister(id int64, config map[string]interface{}, jobMonitorChan chan JobMonitor, jobArchiveChan chan JobArchive) error {
	job, err := NewJob(config)
	if err != nil {
		return err
	}
	job.ID = id
	job.MonitorReport = func(monitor *Monitor) {
		jobMonitor := JobMonitor{
			Job:     job,
			Monitor: monitor,
		}
		jobMonitorChan <- jobMonitor
	}
	job.ArchiveReport = func(archive *Archive) {
		jobArchive := JobArchive{
			Job:     job,
			Archive: archive,
		}
		jobArchiveChan <- jobArchive
	}
	Jobs[id] = job
	return nil
}
