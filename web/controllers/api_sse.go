 package controllers

 import (
 	"fmt"
 	"io"
 	"time"

 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 // APISSELogController 实时日志流（SSE）：周期性 tail agent 上的日志文件。
 // 默认每 3s 拉一次最近 50 行。
 type APISSELogController struct{}

 func NewAPISSELogController() *APISSELogController {
 	return &APISSELogController{}
 }

 func (c *APISSELogController) AppOutLogStream(ctx *gin.Context) {
 	c.streamLog(ctx, "app", "out")
 }

 func (c *APISSELogController) AppErrLogStream(ctx *gin.Context) {
 	c.streamLog(ctx, "app", "err")
 }

 func (c *APISSELogController) JobOutLogStream(ctx *gin.Context) {
 	c.streamLog(ctx, "job", "out")
 }

 func (c *APISSELogController) JobErrLogStream(ctx *gin.Context) {
 	c.streamLog(ctx, "job", "err")
 }

 func (c *APISSELogController) TimingOutLogStream(ctx *gin.Context) {
 	c.streamLog(ctx, "timing", "out")
 }

 func (c *APISSELogController) TimingErrLogStream(ctx *gin.Context) {
 	c.streamLog(ctx, "timing", "err")
 }

 // APISSEMonitorController 实时监控数据流。
 type APISSEMonitorController struct{}

 func NewAPISSEMonitorController() *APISSEMonitorController {
 	return &APISSEMonitorController{}
 }

 func (c *APISSEMonitorController) AgentStream(ctx *gin.Context) {
 	c.streamMonitor(ctx, "agent")
 }

 func (c *APISSEMonitorController) AppStream(ctx *gin.Context) {
 	c.streamMonitor(ctx, "app")
 }

 func (c *APISSEMonitorController) JobStream(ctx *gin.Context) {
 	c.streamMonitor(ctx, "job")
 }

 func (c *APISSEMonitorController) TimingStream(ctx *gin.Context) {
 	c.streamMonitor(ctx, "timing")
 }

 // —— 共享 SSE 工具 ——

 // writeSSEHeaders 写入 SSE 必要的响应头。
 func writeSSEHeaders(ctx *gin.Context) {
 	h := ctx.Writer.Header()
 	h.Set("Content-Type", "text/event-stream")
 	h.Set("Cache-Control", "no-cache")
 	h.Set("Connection", "keep-alive")
 	h.Set("X-Accel-Buffering", "no")
 }

 // writeSSEEvent 写一个 data 事件，payload 会被格式化为 data: <payload>\n\n。
 func writeSSEEvent(w io.Writer, event string, payload []byte) {
 	if event != "" {
 		fmt.Fprintf(w, "event: %s\n", event)
 	}
 	fmt.Fprintf(w, "data: %s\n\n", string(payload))
 	if f, ok := w.(interface{ Flush() }); ok {
 		f.Flush()
 	}
 }

 func (c *APISSELogController) streamLog(ctx *gin.Context, kind, logType string) {
 	var (
 		entityID int64
 		path     string
 		creator  int64
 		agent    *models.Agent
 	)
 	switch kind {
 	case "app":
 		id := utils.QueryInt64(ctx, "app_id", 0)
 		app := providers.AppService.GetAppByID(id)
 		if app == nil {
 			utils.APIError(ctx, "应用不存在")
 			return
 		}
 		entityID, path, creator = app.ID, app.StdOut, app.Creator
 		if logType == "err" {
 			path = app.StdErr
 		}
 		agent = providers.AgentService.GetAgentByID(app.AgentID)
 	case "job":
 		id := utils.QueryInt64(ctx, "job_id", 0)
 		job := providers.JobService.GetJobByID(id)
 		if job == nil {
 			utils.APIError(ctx, "计划任务不存在")
 			return
 		}
 		entityID, path, creator = job.ID, job.StdOut, job.Creator
 		if logType == "err" {
 			path = job.StdErr
 		}
 		agent = providers.AgentService.GetAgentByID(job.AgentID)
 	case "timing":
 		id := utils.QueryInt64(ctx, "timing_id", 0)
 		timing := providers.TimingService.GetTimingByID(id)
 		if timing == nil {
 			utils.APIError(ctx, "定时任务不存在")
 			return
 		}
 		entityID, path, creator = timing.ID, timing.StdOut, timing.Creator
 		if logType == "err" {
 			path = timing.StdErr
 		}
 		agent = providers.AgentService.GetAgentByID(timing.AgentID)
 	}
 	if agent == nil || agent.Status != constants.AGENT_ONLINE {
 		utils.APIError(ctx, "实例不在线")
 		return
 	}
 	if !checkCreator(ctx, creator) {
 		return
 	}

 	writeSSEHeaders(ctx)
 	ctx.Writer.WriteHeaderNow()
 	interval := time.Duration(utils.QueryInt(ctx, "interval", 3)) * time.Second
 	ticker := time.NewTicker(interval)
 	defer ticker.Stop()
 	clientGone := ctx.Request.Context().Done()
 	for {
 		select {
 		case <-clientGone:
 			return
 		case <-ticker.C:
 			content, err := fetchLog(agent, path, constants.WEB_LOG_SIZE)
 			if err != nil {
 				writeSSEEvent(ctx.Writer, "error", []byte(err.Error()))
 				continue
 			}
 			payload := []byte("{\"id\":" + itoa(entityID) + ",\"lines\":" + itoa(int64(len(content))) + "}")
 			writeSSEEvent(ctx.Writer, "meta", payload)
 			for _, line := range content {
 				writeSSEEvent(ctx.Writer, "log", []byte(line))
 			}
 			writeSSEEvent(ctx.Writer, "ping", []byte("{}"))
 		}
 	}
 }

 func (c *APISSEMonitorController) streamMonitor(ctx *gin.Context, kind string) {
 	var (
 		id     int64
 		getter func() []models.Monitor
 	)
 	switch kind {
 	case "agent":
 		id = utils.QueryInt64(ctx, "agent_id", 0)
 		if providers.AgentService.GetAgentByID(id) == nil {
 			utils.APIError(ctx, "实例不存在")
 			return
 		}
 		getter = func() []models.Monitor {
 			return providers.MoniterService.GetAgentMonitor(id, 60)
 		}
 	case "app":
 		id = utils.QueryInt64(ctx, "app_id", 0)
 		if providers.AppService.GetAppByID(id) == nil {
 			utils.APIError(ctx, "应用不存在")
 			return
 		}
 		getter = func() []models.Monitor {
 			return providers.MoniterService.GetAppMonitor(id, 60)
 		}
 	case "job":
 		id = utils.QueryInt64(ctx, "job_id", 0)
 		if providers.JobService.GetJobByID(id) == nil {
 			utils.APIError(ctx, "计划任务不存在")
 			return
 		}
 		getter = func() []models.Monitor {
 			return providers.MoniterService.GetJobMonitor(id, 60)
 		}
 	case "timing":
 		id = utils.QueryInt64(ctx, "timing_id", 0)
 		if providers.TimingService.GetTimingByID(id) == nil {
 			utils.APIError(ctx, "定时任务不存在")
 			return
 		}
 		getter = func() []models.Monitor {
 			return providers.MoniterService.GetTimingMonitor(id, 60)
 		}
 	}

 	writeSSEHeaders(ctx)
 	ctx.Writer.WriteHeaderNow()
 	interval := time.Duration(utils.QueryInt(ctx, "interval", 5)) * time.Second
 	ticker := time.NewTicker(interval)
 	defer ticker.Stop()
 	clientGone := ctx.Request.Context().Done()
 	for {
 		select {
 		case <-clientGone:
 			return
 		case <-ticker.C:
 			list := getter()
 			for i := range list {
 				row := []byte(fmt.Sprintf(`{"cpu":%f,"memory":%f,"time":"%s"}`,
 					list[i].CPU, list[i].Memory, list[i].CreatedAt.Format("2006-01-02 15:04:05")))
 				writeSSEEvent(ctx.Writer, "point", row)
 			}
 			writeSSEEvent(ctx.Writer, "ping", []byte("{}"))
 		}
 	}
 }
