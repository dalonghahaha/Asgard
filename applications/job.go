package applications

import (
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/robfig/cron/v3"
)

var crontab *cron.Cron

var Jobs = []*Job{}

func buildJob(config map[string]string) (*Job, error) {
	job := Job{
		Spec: config["spec"],
		App: &App{
			Name:    config["name"],
			Dir:     config["dir"],
			Program: config["program"],
			Args:    config["args"],
			Stdout:  config["stdout"],
			Stderr:  config["stderr"],
		},
	}
	return &job, nil
}

func AddJob(config map[string]string) error {
	job, err := buildJob(config)
	if job != nil {
		logger.Error("build Job Error:", err)
		return err
	}
	Jobs = append(Jobs, job)
	id, err := crontab.AddJob(job.Spec, job)
	if err != nil {
		logger.Error("AddJob Error:", err)
		return err
	}
	job.ID = id
	return nil
}

func RegisterJob(config map[string]string) error {
	job, err := buildJob(config)
	if job != nil {
		logger.Error("build Job Error:", err)
		return err
	}
	Jobs = append(Jobs, job)
	return nil
}

func CronAll() {
	crontab = cron.New()
	for _, job := range Jobs {
		id, err := crontab.AddJob(job.Spec, job)
		if err != nil {
			logger.Error("AddJob Error:", err)
		}
		job.ID = id
	}
	crontab.Start()
}

func StopJob(name string) error {
	for _, job := range Jobs {
		if job.App.Name == name {
			err := job.App.Cmd.Process.Kill()
			if err != nil {
				logger.Error("app Kill Error:", err)
			}
			crontab.Remove(job.ID)
		}
	}
	return nil
}

type Job struct {
	ID      cron.EntryID
	Spec    string
	TimeOut time.Duration
	App     *App
}

func (j *Job) Run() {
	stop := make(chan bool, 1)
	go j.timer(stop)
	j.App.Exec()
	stop <- true
	go j.record()
}


func (j *Job) timer(ch chan bool) {
	timer := time.NewTimer(j.TimeOut)
	for{
		select{
		case <- timer.C:
			err := j.App.Cmd.Process.Kill()
			if err != nil {
				logger.Error("app Kill Error:", err)
			}
		case <- ch:
			timer.Stop()
		}
	}
}

func (j *Job) record() {
	//TODO record
}
