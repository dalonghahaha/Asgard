package utils

import (
	"Asgard/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetReferer(ctx *gin.Context) string {
	return ctx.Request.Header.Get("Referer")
}

func DefaultInt(ctx *gin.Context, key string, defaultVal int) int {
	val := ctx.Query(key)
	if val == "" {
		return defaultVal
	}
	_val, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return _val
}

func DefaultInt64(ctx *gin.Context, key string, defaultVal int64) int64 {
	val := ctx.Query(key)
	if val == "" {
		return defaultVal
	}
	_val, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return _val
}

func FormDefaultInt(ctx *gin.Context, key string, defaultVal int) int {
	val := ctx.PostForm(key)
	if val == "" {
		return defaultVal
	}
	_val, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return _val
}

func FormDefaultInt64(ctx *gin.Context, key string, defaultVal int64) int64 {
	val := ctx.PostForm(key)
	if val == "" {
		return defaultVal
	}
	_val, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return _val
}

func Required(ctx *gin.Context, val string, message string) bool {
	if val == "" {
		APIBadRequest(ctx, message)
		return false
	}
	return true
}

func GetAgent(ctx *gin.Context) *models.Agent {
	agent, ok := ctx.Get("agent")
	if !ok {
		return nil
	}
	_agent, ok := agent.(*models.Agent)
	if !ok {
		return nil
	}
	return _agent
}

func GetGroup(ctx *gin.Context) *models.Group {
	group, ok := ctx.Get("group")
	if !ok {
		return nil
	}
	_group, ok := group.(*models.Group)
	if !ok {
		return nil
	}
	return _group
}

func GetApp(ctx *gin.Context) *models.App {
	app, ok := ctx.Get("app")
	if !ok {
		return nil
	}
	_app, ok := app.(*models.App)
	if !ok {
		return nil
	}
	return _app
}

func GetJob(ctx *gin.Context) *models.Job {
	job, ok := ctx.Get("job")
	if !ok {
		return nil
	}
	_job, ok := job.(*models.Job)
	if !ok {
		return nil
	}
	return _job
}

func GetTiming(ctx *gin.Context) *models.Timing {
	timing, ok := ctx.Get("timing")
	if !ok {
		return nil
	}
	_timing, ok := timing.(*models.Timing)
	if !ok {
		return nil
	}
	return _timing
}

func GetAppAgent(ctx *gin.Context) map[*models.App]*models.Agent {
	appAgents, ok := ctx.Get("app_agent")
	if !ok {
		return map[*models.App]*models.Agent{}
	}
	_appAgents, ok := appAgents.(map[*models.App]*models.Agent)
	if !ok {
		return map[*models.App]*models.Agent{}
	}
	return _appAgents
}

func GetJobAgent(ctx *gin.Context) map[*models.Job]*models.Agent {
	jobAgents, ok := ctx.Get("job_agent")
	if !ok {
		return map[*models.Job]*models.Agent{}
	}
	_jobAgents, ok := jobAgents.(map[*models.Job]*models.Agent)
	if !ok {
		return map[*models.Job]*models.Agent{}
	}
	return _jobAgents
}

func GetTimingAgent(ctx *gin.Context) map[*models.Timing]*models.Agent {
	timingAgents, ok := ctx.Get("timing_agent")
	if !ok {
		return map[*models.Timing]*models.Agent{}
	}
	_timingAgents, ok := timingAgents.(map[*models.Timing]*models.Agent)
	if !ok {
		return map[*models.Timing]*models.Agent{}
	}
	return _timingAgents
}
