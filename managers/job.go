package managers

import (
	"Asgard/applications"
	"Asgard/client"
	"fmt"
	"sync"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
	"github.com/robfig/cron/v3"
)

type JobManager struct {
	lock           sync.Mutex
	crontab        *cron.Cron
	jobs           map[int64]*applications.Job
	masterClient   *client.Master
	monitorManager *applications.MonitorMamager
}

func NewJobManager() (*JobManager, error) {
	manager := &JobManager{
		jobs:           make(map[int64]*applications.Job),
		monitorManager: applications.NewMonitorMamager(),
	}
	return manager, nil
}

func (m *JobManager) SetMaster(masterClient *client.Master) {
	m.masterClient = masterClient
}

func (m *JobManager) StartMonitor() {
	m.monitorManager.Start()
}

func (m *JobManager) StopMonitor() {
	m.monitorManager.Stop()
}

func (m *JobManager) Count() int {
	return len(m.jobs)
}

func (m *JobManager) GetJob(id int64) *applications.Job {
	app, ok := m.jobs[id]
	if !ok {
		return nil
	}
	return app
}

func (m *JobManager) StartAll() {
	m.crontab = cron.New()
	for _, job := range m.jobs {
		m.Add(job)
	}
	m.crontab.Start()
}

func (m *JobManager) StopAll() {
	for _, job := range m.jobs {
		if job.Running {
			job.Kill()
		}
	}
}

func (m *JobManager) Add(job *applications.Job) {
	id, err := m.crontab.AddJob(job.Spec, job)
	if err != nil {
		logger.Error(job.Name+" add fail:", err)
	}
	logger.Info(job.Name + " add seccess!")
	job.CronID = id
}

func (m *JobManager) Remove(id int64) bool {
	job := m.GetJob(id)
	if job == nil {
		return true
	}
	if job.Running {
		job.Kill()
	}
	m.crontab.Remove(job.CronID)
	return true
}

func (m *JobManager) NewJob(config map[string]interface{}) (*applications.Job, error) {
	job := new(applications.Job)
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

func (m *JobManager) Register(id int64, config map[string]interface{}) error {
	job, err := m.NewJob(config)
	if err != nil {
		return err
	}
	job.ID = id
	job.MonitorMamager = m.monitorManager
	job.MonitorReport = func(monitor *applications.Monitor) {
		jobMonitor := applications.JobMonitor{
			UUID:    uuid.GenerateV4(),
			Job:     job,
			Monitor: monitor,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(jobMonitor.UUID, 1)
			m.masterClient.JobMonitorChan <- jobMonitor
		}
	}
	job.ArchiveReport = func(archive *applications.Archive) {
		jobArchive := applications.JobArchive{
			UUID:    uuid.GenerateV4(),
			Job:     job,
			Archive: archive,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(jobArchive.UUID, 1)
			m.masterClient.JobArchiveChan <- jobArchive
		}
	}
	m.lock.Lock()
	m.jobs[id] = job
	m.lock.Unlock()
	return nil
}

func (m *JobManager) UnRegister(id int64) {
	m.lock.Lock()
	delete(m.jobs, id)
	m.lock.Unlock()
}
