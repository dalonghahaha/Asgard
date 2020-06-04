package applications

import (
	"fmt"
	"sync"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
	"github.com/robfig/cron/v3"
)

var (
	crontab *cron.Cron
	Jobs    = map[int64]*Job{}
	jobLock sync.Mutex
)

func JobStopAll() {
	for _, job := range Jobs {
		if job.Running {
			job.Kill()
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

func JobStop(id int64) bool {
	job, ok := Jobs[id]
	if !ok {
		return false
	}
	if job.Running {
		job.Kill()
	}
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
	err := job.Configure(config)
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

func JobRegister(id int64, config map[string]interface{}, monitorMamager *MonitorMamager, reports *sync.Map, mc chan JobMonitor, ac chan JobArchive) error {
	job, err := NewJob(config)
	if err != nil {
		return err
	}
	job.ID = id
	job.MonitorMamager = monitorMamager
	job.MonitorReport = func(monitor *Monitor) {
		jobMonitor := JobMonitor{
			UUID:    uuid.GenerateV4(),
			Job:     job,
			Monitor: monitor,
		}
		if reports != nil {
			reports.Store(jobMonitor.UUID, 1)
		}
		if mc != nil {
			mc <- jobMonitor
		}
	}
	job.ArchiveReport = func(archive *Archive) {
		jobArchive := JobArchive{
			UUID:    uuid.GenerateV4(),
			Job:     job,
			Archive: archive,
		}
		if reports != nil {
			reports.Store(jobArchive.UUID, 1)
		}
		if ac != nil {
			ac <- jobArchive
		}
	}
	jobLock.Lock()
	Jobs[id] = job
	jobLock.Unlock()
	return nil
}

func JobUnRegister(id int64) {
	jobLock.Lock()
	delete(Jobs, id)
	jobLock.Unlock()
}
