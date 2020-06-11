package managers

import (
	"Asgard/clients"
	"Asgard/runtimes"
	"fmt"
	"sync"
	"time"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"
	"github.com/robfig/cron/v3"
)

type JobManager struct {
	lock         sync.Mutex
	crontab      *cron.Cron
	jobs         map[int64]*runtimes.Job
	masterClient *clients.Master
	monitor      *runtimes.Monitor
}

func NewJobManager() *JobManager {
	manager := &JobManager{
		jobs:    make(map[int64]*runtimes.Job),
		monitor: runtimes.NewMonitor("job"),
	}
	return manager
}

func (m *JobManager) SetMaster(masterClient *clients.Master) {
	m.masterClient = masterClient
}

func (m *JobManager) StartMonitor() {
	go m.monitor.Start()
}

func (m *JobManager) NewJob(config map[string]interface{}) (*runtimes.Job, error) {
	job := new(runtimes.Job)
	err := job.Configure(config)
	if err != nil {
		return nil, err
	}
	spec, ok := config["spec"].(string)
	if !ok {
		return nil, fmt.Errorf("config spec type wrong")
	}
	job.Spec = spec
	timeout1, ok1 := config["timeout"].(int64)
	timeout2, ok2 := config["timeout"].(int)
	if !ok1 && !ok2 {
		return nil, fmt.Errorf("config timeout type wrong")
	}
	if ok1 {
		job.TimeOut = time.Duration(timeout1)
	}
	if ok2 {
		job.TimeOut = time.Duration(timeout2)
	}
	return job, nil
}

func (m *JobManager) Register(id int64, config map[string]interface{}) error {
	job, err := m.NewJob(config)
	if err != nil {
		return err
	}
	job.ID = id
	job.Monitor = m.monitor
	job.ExceptionReport = func(message string) {
		jobException := runtimes.JobException{
			UUID:  uuid.GenerateV4(),
			JobID: job.ID,
			Desc:  message,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(jobException.UUID, 1)
			m.masterClient.JobExceptionChan <- jobException
		}
	}
	job.MonitorReport = func(monitor *runtimes.MonitorInfo) {
		jobMonitor := runtimes.JobMonitor{
			UUID:    uuid.GenerateV4(),
			Job:     job,
			Monitor: monitor,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(jobMonitor.UUID, 1)
			m.masterClient.JobMonitorChan <- jobMonitor
		}
	}
	job.ArchiveReport = func(archive *runtimes.Archive) {
		jobArchive := runtimes.JobArchive{
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

func (m *JobManager) Count() int {
	return len(m.jobs)
}

func (m *JobManager) GetList() []*runtimes.Job {
	list := []*runtimes.Job{}
	for _, job := range m.jobs {
		list = append(list, job)
	}
	return list
}

func (m *JobManager) Get(id int64) *runtimes.Job {
	app, ok := m.jobs[id]
	if !ok {
		return nil
	}
	return app
}

func (m *JobManager) GetByName(name string) *runtimes.Job {
	for _, job := range m.jobs {
		if job.Name == name {
			return job
		}
	}
	return nil
}

func (m *JobManager) StartAll(monitor bool) {
	if monitor {
		m.StartMonitor()
	}
	m.crontab = cron.New()
	for _, job := range m.jobs {
		m.Create(job)
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

func (m *JobManager) Create(job *runtimes.Job) bool {
	id, err := m.crontab.AddJob(job.Spec, job)
	if err != nil {
		logger.Error(job.Name+" add fail:", err)
		return false
	}
	job.CronID = id
	return true
}

func (m *JobManager) Stop(id int64) bool {
	job := m.Get(id)
	if job == nil {
		return true
	}
	if job.Running {
		job.Kill()
	}
	m.crontab.Remove(job.CronID)
	return true
}

func (m *JobManager) Add(id int64, config map[string]interface{}) error {
	app := m.Get(id)
	if app != nil {
		return nil
	}
	err := m.Register(id, config)
	if err != nil {
		return fmt.Errorf("job register failed:%s", err.Error())
	}
	ok := m.Create(m.jobs[id])
	if !ok {
		return fmt.Errorf("job create failed!")
	}
	return nil
}

func (m *JobManager) Update(id int64, config map[string]interface{}) error {
	err := m.Remove(id)
	if err != nil {
		return err
	}
	return m.Add(id, config)
}

func (m *JobManager) Remove(id int64) error {
	ok := m.Stop(id)
	if !ok {
		return fmt.Errorf("job stop failed!")
	}
	m.UnRegister(id)
	return nil
}
