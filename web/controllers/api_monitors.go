 package controllers

 import (
 	"github.com/gin-gonic/gin"

 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 // APIMonitorController 提供 agent/app/job/timing 4 类监控的查询接口。
 type APIMonitorController struct{}

 func NewAPIMonitorController() *APIMonitorController {
 	return &APIMonitorController{}
 }

 func (c *APIMonitorController) Agent(ctx *gin.Context) {
 	id := utils.QueryInt64(ctx, "agent_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "agent_id 不能为空")
 		return
 	}
 	if providers.AgentService.GetAgentByID(id) == nil {
 		utils.APIBadRequest(ctx, "实例不存在")
 		return
 	}
 	size := utils.QueryInt(ctx, "size", 100)
 	list := providers.MoniterService.GetAgentMonitor(id, size)
 	utils.APIData(ctx, monitorPoints(list))
 }

 func (c *APIMonitorController) App(ctx *gin.Context) {
 	id := utils.QueryInt64(ctx, "app_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "app_id 不能为空")
 		return
 	}
 	if providers.AppService.GetAppByID(id) == nil {
 		utils.APIBadRequest(ctx, "应用不存在")
 		return
 	}
 	size := utils.QueryInt(ctx, "size", 100)
 	list := providers.MoniterService.GetAppMonitor(id, size)
 	utils.APIData(ctx, monitorPoints(list))
 }

 func (c *APIMonitorController) Job(ctx *gin.Context) {
 	id := utils.QueryInt64(ctx, "job_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "job_id 不能为空")
 		return
 	}
 	if providers.JobService.GetJobByID(id) == nil {
 		utils.APIBadRequest(ctx, "计划任务不存在")
 		return
 	}
 	size := utils.QueryInt(ctx, "size", 100)
 	list := providers.MoniterService.GetJobMonitor(id, size)
 	utils.APIData(ctx, monitorPoints(list))
 }

 func (c *APIMonitorController) Timing(ctx *gin.Context) {
 	id := utils.QueryInt64(ctx, "timing_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "timing_id 不能为空")
 		return
 	}
 	if providers.TimingService.GetTimingByID(id) == nil {
 		utils.APIBadRequest(ctx, "定时任务不存在")
 		return
 	}
 	size := utils.QueryInt(ctx, "size", 100)
 	list := providers.MoniterService.GetTimingMonitor(id, size)
 	utils.APIData(ctx, monitorPoints(list))
 }

 func monitorPoints(list []models.Monitor) []gin.H {
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, gin.H{
 			"cpu":        list[i].CPU,
 			"memory":     list[i].Memory,
 			"created_at": list[i].CreatedAt,
 		})
 	}
 	return out
 }
