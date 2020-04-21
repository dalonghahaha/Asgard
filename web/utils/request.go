package utils

import (
	"Asgard/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
