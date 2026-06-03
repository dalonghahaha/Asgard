 package controllers

 import (
 	"fmt"
 	"time"

 	"github.com/dalonghahaha/avenger/components/logger"
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 type apiTimingReq struct {
 	GroupID   int64     `form:"group_id" json:"group_id"`
 	AgentID   int64     `form:"agent_id" json:"agent_id" binding:"required"`
 	Name      string    `form:"name" json:"name" binding:"required"`
 	Dir       string    `form:"dir" json:"dir" binding:"required"`
 	Program   string    `form:"program" json:"program" binding:"required"`
 	Args      string    `form:"args" json:"args"`
 	StdOut    string    `form:"std_out" json:"std_out" binding:"required"`
 	StdErr    string    `form:"std_err" json:"std_err" binding:"required"`
 	Time      time.Time `form:"time" json:"time" binding:"required"`
 	Timeout   int64     `form:"timeout" json:"timeout"`
 	IsMonitor int       `form:"is_monitor" json:"is_monitor"`
 }

 type APITimingController struct{}

 func NewAPITimingController() *APITimingController {
 	return &APITimingController{}
 }

 func (c *APITimingController) List(ctx *gin.Context) {
 	page := utils.QueryInt(ctx, "page", 1)
 	groupID := utils.QueryInt64(ctx, "group_id", 0)
 	agentID := utils.QueryInt64(ctx, "agent_id", 0)
 	status := utils.QueryInt(ctx, "status", -99)
 	name := ctx.Query("name")
 	user := utils.GetUser(ctx)
 	where := map[string]interface{}{"status": status}
 	if user.Role != constants.USER_ROLE_ADMIN {
 		where["creator"] = user.ID
 	}
 	if groupID != 0 {
 		where["group_id"] = groupID
 	}
 	if agentID != 0 {
 		where["agent_id"] = agentID
 	}
 	if name != "" {
 		where["name"] = name
 	}
 	list, total := providers.TimingService.GetTimingPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.TimingFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 func (c *APITimingController) Show(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	timing := providers.TimingService.GetTimingByID(id)
 	if timing == nil {
 		utils.APIBadRequest(ctx, "定时任务不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此定时任务的权限")
 		return
 	}
 	utils.APIData(ctx, utils.TimingFormat(timing))
 }

 func (c *APITimingController) Create(ctx *gin.Context) {
 	var req apiTimingReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	t := new(models.Timing)
 	t.GroupID = req.GroupID
 	t.AgentID = req.AgentID
 	t.Name = req.Name
 	t.Dir = req.Dir
 	t.Program = req.Program
 	t.Args = req.Args
 	t.StdOut = req.StdOut
 	t.StdErr = req.StdErr
 	t.Time = req.Time
 	t.Timeout = req.Timeout
 	t.IsMonitor = int64(req.IsMonitor)
 	t.Status = constants.TIMING_STATUS_PAUSE
 	t.Creator = utils.GetUserID(ctx)
 	if !providers.TimingService.CreateTiming(t) {
 		utils.APIError(ctx, "创建定时任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, t.ID, constants.ACTION_CREATE)
 	utils.APIData(ctx, gin.H{"id": t.ID})
 }

 func (c *APITimingController) Update(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	timing := providers.TimingService.GetTimingByID(id)
 	if timing == nil {
 		utils.APIBadRequest(ctx, "定时任务不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此定时任务的权限")
 		return
 	}
 	var req apiTimingReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	timing.GroupID = req.GroupID
 	timing.Name = req.Name
 	timing.Dir = req.Dir
 	timing.Program = req.Program
 	timing.Args = req.Args
 	timing.StdOut = req.StdOut
 	timing.StdErr = req.StdErr
 	timing.Time = req.Time
 	timing.Timeout = req.Timeout
 	timing.IsMonitor = int64(req.IsMonitor)
 	timing.Updator = utils.GetUserID(ctx)
 	if req.AgentID != 0 {
 		timing.AgentID = req.AgentID
 	}
 	if !providers.TimingService.UpdateTiming(timing) {
 		utils.APIError(ctx, "更新定时任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_UPDATE)
 	utils.APIOK(ctx)
 }

 func (c *APITimingController) Copy(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	timing := providers.TimingService.GetTimingByID(id)
 	if timing == nil {
 		utils.APIBadRequest(ctx, "定时任务不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此定时任务的权限")
 		return
 	}
 	_t := new(models.Timing)
 	_t.GroupID = timing.GroupID
 	_t.Name = timing.Name + "_copy"
 	_t.AgentID = timing.AgentID
 	_t.Dir = timing.Dir
 	_t.Program = timing.Program
 	_t.Args = timing.Args
 	_t.StdOut = timing.StdOut
 	_t.StdErr = timing.StdErr
 	_t.Time = timing.Time
 	_t.Timeout = timing.Timeout
 	_t.IsMonitor = timing.IsMonitor
 	_t.Status = constants.TIMING_STATUS_PAUSE
 	_t.Creator = utils.GetUserID(ctx)
 	if !providers.TimingService.CreateTiming(_t) {
 		utils.APIError(ctx, "复制定时任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_COPY)
 	utils.APIData(ctx, gin.H{"id": _t.ID})
 }

 func loadTimingAgentForUser(ctx *gin.Context, id int64) (*models.Timing, *models.Agent, bool) {
 	timing := providers.TimingService.GetTimingByID(id)
 	if timing == nil {
 		utils.APIBadRequest(ctx, "定时任务不存在")
 		return nil, nil, false
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此定时任务的权限")
 		return nil, nil, false
 	}
 	agent := providers.AgentService.GetAgentByID(timing.AgentID)
 	if agent == nil {
 		utils.APIBadRequest(ctx, "实例获取失败")
 		return nil, nil, false
 	}
 	if agent.Status != constants.AGENT_ONLINE {
 		utils.APIBadRequest(ctx, "实例不处于运行状态")
 		return nil, nil, false
 	}
 	return timing, agent, true
 }

 func timingStart(agent *models.Agent, t *models.Timing) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetTiming(t.ID)
 	if err != nil {
 		return fmt.Errorf("获取定时任务情况异常:%s", err.Error())
 	}
 	if exist == nil {
 		if err := client.AddTiming(t); err != nil {
 			return fmt.Errorf("添加定时任务异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func timingRestart(agent *models.Agent, t *models.Timing) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetTiming(t.ID)
 	if err != nil {
 		return fmt.Errorf("获取定时任务情况异常:%s", err.Error())
 	}
 	if exist == nil {
 		if err := client.AddTiming(t); err != nil {
 			return fmt.Errorf("重启定时任务异常:%s", err.Error())
 		}
 	} else {
 		if err := client.UpdateTiming(t); err != nil {
 			return fmt.Errorf("重启定时任务异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func timingStop(agent *models.Agent, timingID int64) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetTiming(timingID)
 	if err != nil {
 		return fmt.Errorf("获取定时任务情况异常:%s", err.Error())
 	}
 	if exist != nil {
 		if err := client.RemoveTiming(timingID); err != nil {
 			return fmt.Errorf("停止定时任务异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func (c *APITimingController) Start(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	timing, agent, ok := loadTimingAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if timing.Status == constants.TIMING_STATUS_RUNNING {
 		utils.APIError(ctx, "定时任务已经启动")
 		return
 	}
 	if err := timingStart(agent, timing); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_RUNNING, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "更新定时任务状态失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_START)
 	utils.APIOK(ctx)
 }

 func (c *APITimingController) ReStart(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	timing, agent, ok := loadTimingAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if err := timingRestart(agent, timing); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_RESTART)
 	utils.APIOK(ctx)
 }

 func (c *APITimingController) Pause(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	timing, agent, ok := loadTimingAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if err := timingStop(agent, timing.ID); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_PAUSE, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "更新定时任务状态失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_PAUSE)
 	utils.APIOK(ctx)
 }

 func (c *APITimingController) Delete(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	timing, agent, ok := loadTimingAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if timing.Status == constants.TIMING_STATUS_RUNNING {
 		utils.APIError(ctx, "定时任务启动状态不能删除")
 		return
 	}
 	if err := timingStop(agent, timing.ID); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_DELETED, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "删除定时任务失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_DELETE)
 	utils.APIOK(ctx)
 }

 func (c *APITimingController) BatchAction(ctx *gin.Context, action string) {
 	var req apiBatchReq
 	if err := ctx.ShouldBind(&req); err != nil || len(req.IDs) == 0 {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	user := utils.GetUser(ctx)
 	for _, id := range req.IDs {
 		timing := providers.TimingService.GetTimingByID(id)
 		if timing == nil {
 			continue
 		}
 		if user.Role != constants.USER_ROLE_ADMIN && timing.Creator != user.ID {
 			continue
 		}
 		agent := providers.AgentService.GetAgentByID(timing.AgentID)
 		if agent == nil || agent.Status != constants.AGENT_ONLINE {
 			continue
 		}
 		switch action {
 		case "start":
 			if timing.Status == constants.TIMING_STATUS_RUNNING {
 				continue
 			}
 			if err := timingStart(agent, timing); err != nil {
 				logger.Errorf("Timing BatchStart err:%s", err.Error())
 				continue
 			}
 			providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_RUNNING, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_TIMING, timing.ID, constants.ACTION_START)
 		case "restart":
 			if err := timingRestart(agent, timing); err != nil {
 				logger.Errorf("Timing BatchRestart err:%s", err.Error())
 				continue
 			}
 			utils.OpetationLog(user.ID, constants.TYPE_TIMING, timing.ID, constants.ACTION_RESTART)
 		case "pause":
 			if err := timingStop(agent, timing.ID); err != nil {
 				logger.Errorf("Timing BatchPause err:%s", err.Error())
 				continue
 			}
 			providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_PAUSE, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_TIMING, timing.ID, constants.ACTION_PAUSE)
 		case "delete":
 			if timing.Status == constants.TIMING_STATUS_RUNNING {
 				continue
 			}
 			if err := timingStop(agent, timing.ID); err != nil {
 				logger.Errorf("Timing BatchDelete err:%s", err.Error())
 				continue
 			}
 			providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_DELETED, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_TIMING, timing.ID, constants.ACTION_DELETE)
 		}
 	}
 	utils.APIOK(ctx)
 }

 func (c *APITimingController) BatchStart(ctx *gin.Context)   { c.BatchAction(ctx, "start") }
 func (c *APITimingController) BatchReStart(ctx *gin.Context) { c.BatchAction(ctx, "restart") }
 func (c *APITimingController) BatchPause(ctx *gin.Context)   { c.BatchAction(ctx, "pause") }
 func (c *APITimingController) BatchDelete(ctx *gin.Context) { c.BatchAction(ctx, "delete") }
