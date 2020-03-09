package rpc

import (
	"Asgard/applications"
	"Asgard/models"
	"time"
)

var (
	OK      = int32(200)
	Nofound = int32(404)
	Error   = int32(500)
)

func FormatApp(app *models.App) *App {
	return &App{
		Id:          app.ID,
		Name:        app.Name,
		Dir:         app.Dir,
		Program:     app.Program,
		Args:        app.Args,
		StdOut:      app.StdOut,
		StdErr:      app.StdErr,
		AutoRestart: app.AutoRestart == 1,
		IsMonitor:   app.IsMonitor == 1,
	}
}

func FormatJob(job *models.Job) *Job {
	return &Job{
		Id:        job.ID,
		Name:      job.Name,
		Dir:       job.Dir,
		Program:   job.Program,
		Args:      job.Args,
		StdOut:    job.StdOut,
		StdErr:    job.StdErr,
		Spec:      job.Spec,
		Timeout:   job.Timeout,
		IsMonitor: job.IsMonitor == 1,
	}
}

func FormatTiming(job *models.Timing) *Timing {
	return &Timing{
		Id:        job.ID,
		Name:      job.Name,
		Dir:       job.Dir,
		Program:   job.Program,
		Args:      job.Args,
		StdOut:    job.StdOut,
		StdErr:    job.StdErr,
		Time:      job.Time.Unix(),
		Timeout:   job.Timeout,
		IsMonitor: job.IsMonitor == 1,
	}
}

func BuildApp(app *applications.App) *App {
	return &App{
		Id:          app.ID,
		Name:        app.Name,
		Dir:         app.Dir,
		Program:     app.Program,
		Args:        app.Args,
		StdOut:      app.Stdout,
		StdErr:      app.Stderr,
		AutoRestart: app.AutoRestart,
		IsMonitor:   app.IsMonitor,
	}
}

func BuildJob(job *applications.Job) *Job {
	return &Job{
		Id:        job.ID,
		Name:      job.Name,
		Dir:       job.Dir,
		Program:   job.Program,
		Args:      job.Args,
		StdOut:    job.Stdout,
		StdErr:    job.Stderr,
		Spec:      job.Spec,
		Timeout:   int64(job.TimeOut),
		IsMonitor: job.IsMonitor,
	}
}

func BuildTiming(job *applications.Timing) *Timing {
	return &Timing{
		Id:        job.ID,
		Name:      job.Name,
		Dir:       job.Dir,
		Program:   job.Program,
		Args:      job.Args,
		StdOut:    job.Stdout,
		StdErr:    job.Stderr,
		Time:      job.Time.Unix(),
		Timeout:   int64(job.TimeOut),
		IsMonitor: job.IsMonitor,
	}
}

func BuildArchive(command *applications.Command) *Archive {
	return &Archive{
		Uuid:      command.UUID,
		Pid:       int32(command.Pid),
		BeginTime: command.Begin.Unix(),
		EndTime:   command.End.Unix(),
		Status:    int32(command.Status),
		Signal:    command.Signal,
	}
}

func BuildAgentMonitor(agentIP string, agentPort string, Pid int, UUID string, monitor *applications.Monitor) *AgentMonitor {
	return &AgentMonitor{
		Agent: &AgentInfo{
			Ip:   agentIP,
			Port: agentPort,
		},
		Monitor: &Monitor{
			Uuid:    UUID,
			Pid:     int32(Pid),
			Cpu:     float32(monitor.CPUPercent),
			Memory:  monitor.MemoryPercent,
			Threads: int32(monitor.NumThreads),
		},
	}
}

func BuildAppMonitor(app *applications.App, monitor *applications.Monitor) *AppMonitor {
	return &AppMonitor{
		App: BuildApp(app),
		Monitor: &Monitor{
			Uuid:    app.UUID,
			Pid:     int32(app.Pid),
			Cpu:     float32(monitor.CPUPercent),
			Memory:  monitor.MemoryPercent,
			Threads: int32(monitor.NumThreads),
		},
	}
}

func BuildJobMonior(job *applications.Job, monitor *applications.Monitor) *JobMonior {
	return &JobMonior{
		Job: BuildJob(job),
		Monitor: &Monitor{
			Uuid:    job.UUID,
			Pid:     int32(job.Pid),
			Cpu:     float32(monitor.CPUPercent),
			Memory:  monitor.MemoryPercent,
			Threads: int32(monitor.NumThreads),
		},
	}
}

func BuildTimingMonior(timing *applications.Timing, monitor *applications.Monitor) *TimingMonior {
	return &TimingMonior{
		Timing: BuildTiming(timing),
		Monitor: &Monitor{
			Uuid:    timing.UUID,
			Pid:     int32(timing.Pid),
			Cpu:     float32(monitor.CPUPercent),
			Memory:  monitor.MemoryPercent,
			Threads: int32(monitor.NumThreads),
		},
	}
}

func BuildAppArchive(app *applications.App, command *applications.Command) *AppArchive {
	return &AppArchive{
		App:     BuildApp(app),
		Archive: BuildArchive(command),
	}
}

func BuildJobArchive(job *applications.Job, command *applications.Command) *JobArchive {
	return &JobArchive{
		Job:     BuildJob(job),
		Archive: BuildArchive(command),
	}
}

func BuildTimingArchive(timing *applications.Timing, command *applications.Command) *TimingArchive {
	return &TimingArchive{
		Timing:  BuildTiming(timing),
		Archive: BuildArchive(command),
	}
}

func ParseArchive(tp int64, relatedID int64, archive *Archive) *models.Archive {
	return &models.Archive{
		Type:      tp,
		RelatedID: relatedID,
		UUID:      archive.GetUuid(),
		PID:       int64(archive.GetPid()),
		BeginTime: time.Unix(archive.GetBeginTime(), 0),
		EndTime:   time.Unix(archive.GetEndTime(), 0),
		Status:    int64(archive.GetStatus()),
		Signal:    archive.GetSignal(),
		CreatedAt: time.Now(),
	}
}

func ParseMonitor(tp int64, relatedID int64, moniter *Monitor) *models.Monitor {
	return &models.Monitor{
		Type:      tp,
		RelatedID: relatedID,
		UUID:      moniter.GetUuid(),
		PID:       int64(moniter.GetPid()),
		CPU:       float64(moniter.GetCpu()),
		Memory:    float64(moniter.GetMemory()),
		CreatedAt: time.Now(),
	}
}

func BuildAppConfig(app *App) map[string]interface{} {
	return map[string]interface{}{
		"id":           app.GetId(),
		"name":         app.GetName(),
		"dir":          app.GetDir(),
		"program":      app.GetProgram(),
		"args":         app.GetArgs(),
		"stdout":       app.GetStdOut(),
		"stderr":       app.GetStdErr(),
		"auto_restart": app.GetAutoRestart(),
		"is_monitor":   app.GetIsMonitor(),
	}
}

func BuildJobConfig(job *Job) map[string]interface{} {
	return map[string]interface{}{
		"id":         job.GetId(),
		"name":       job.GetName(),
		"dir":        job.GetDir(),
		"program":    job.GetProgram(),
		"args":       job.GetArgs(),
		"stdout":     job.GetStdOut(),
		"stderr":     job.GetStdErr(),
		"spec":       job.GetSpec(),
		"timeout":    job.GetTimeout(),
		"is_monitor": job.GetIsMonitor(),
	}
}

func BuildTimingConfig(timing *Timing) map[string]interface{} {
	return map[string]interface{}{
		"id":         timing.GetId(),
		"name":       timing.GetName(),
		"dir":        timing.GetDir(),
		"program":    timing.GetProgram(),
		"args":       timing.GetArgs(),
		"stdout":     timing.GetStdOut(),
		"stderr":     timing.GetStdErr(),
		"time":       timing.GetTime(),
		"timeout":    timing.GetTimeout(),
		"is_monitor": timing.GetIsMonitor(),
	}
}
