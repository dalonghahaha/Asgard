 package controllers

 import (
 	"fmt"

 	"github.com/dalonghahaha/avenger/components/logger"
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 // apiAppReq 是 App 创建/更新通用请求。
 type apiAppReq struct {
 	GroupID     int64  `form:"group_id" json:"group_id"`
 	AgentID     int64  `form:"agent_id" json:"agent_id"`
 	Name        string `form:"name" json:"name" binding:"required"`
 	Dir         string `form:"dir" json:"dir" binding:"required"`
 	Program     string `form:"program" json:"program" binding:"required"`
 	Args        string `form:"args" json:"args"`
 	StdOut      string `form:"std_out" json:"std_out" binding:"required"`
 	StdErr      string `form:"std_err" json:"std_err" binding:"required"`
 	AutoRestart int    `form:"auto_restart" json:"auto_restart"`
 	IsMonitor   int    `form:"is_monitor" json:"is_monitor"`
 }

 type APIAppController struct{}

 func NewAPIAppController() *APIAppController {
 	return &APIAppController{}
 }

 func (c *APIAppController) List(ctx *gin.Context) {
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
 	list, total := providers.AppService.GetAppPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.AppFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 func (c *APIAppController) Show(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	app := providers.AppService.GetAppByID(id)
 	if app == nil {
 		utils.APIBadRequest(ctx, "应用不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此应用的权限")
 		return
 	}
 	utils.APIData(ctx, utils.AppFormat(app))
 }

 func (c *APIAppController) Create(ctx *gin.Context) {
 	var req apiAppReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	if req.AgentID == 0 {
 		utils.APIBadRequest(ctx, "运行实例未选择")
 		return
 	}
 	app := new(models.App)
 	app.GroupID = req.GroupID
 	app.AgentID = req.AgentID
 	app.Name = req.Name
 	app.Dir = req.Dir
 	app.Program = req.Program
 	app.Args = req.Args
 	app.StdOut = req.StdOut
 	app.StdErr = req.StdErr
 	app.Status = constants.APP_STATUS_PAUSE
 	app.Creator = utils.GetUserID(ctx)
 	app.AutoRestart = int64(req.AutoRestart)
 	app.IsMonitor = int64(req.IsMonitor)
 	if !providers.AppService.CreateApp(app) {
 		utils.APIError(ctx, "创建应用失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_APP, app.ID, constants.ACTION_CREATE)
 	utils.APIData(ctx, gin.H{"id": app.ID})
 }

 func (c *APIAppController) Update(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	app := providers.AppService.GetAppByID(id)
 	if app == nil {
 		utils.APIBadRequest(ctx, "应用不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此应用的权限")
 		return
 	}
 	var req apiAppReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	app.GroupID = req.GroupID
 	app.Name = req.Name
 	app.Dir = req.Dir
 	app.Program = req.Program
 	app.Args = req.Args
 	app.StdOut = req.StdOut
 	app.StdErr = req.StdErr
 	app.Updator = utils.GetUserID(ctx)
 	app.AutoRestart = int64(req.AutoRestart)
 	app.IsMonitor = int64(req.IsMonitor)
 	if req.AgentID != 0 {
 		app.AgentID = req.AgentID
 	}
 	if !providers.AppService.UpdateApp(app) {
 		utils.APIError(ctx, "更新应用失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_APP, app.ID, constants.ACTION_UPDATE)
 	utils.APIOK(ctx)
 }

 func (c *APIAppController) Copy(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	app := providers.AppService.GetAppByID(id)
 	if app == nil {
 		utils.APIBadRequest(ctx, "应用不存在")
 		return
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此应用的权限")
 		return
 	}
 	_app := new(models.App)
 	_app.GroupID = app.GroupID
 	_app.Name = app.Name + "_copy"
 	_app.AgentID = app.AgentID
 	_app.Dir = app.Dir
 	_app.Program = app.Program
 	_app.Args = app.Args
 	_app.StdOut = app.StdOut
 	_app.StdErr = app.StdErr
 	_app.AutoRestart = app.AutoRestart
 	_app.IsMonitor = app.IsMonitor
 	_app.Status = constants.APP_STATUS_PAUSE
 	_app.Creator = utils.GetUserID(ctx)
 	if !providers.AppService.CreateApp(_app) {
 		utils.APIError(ctx, "复制应用失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_APP, app.ID, constants.ACTION_COPY)
 	utils.APIData(ctx, gin.H{"id": _app.ID})
 }

 // —— 控制类（start/restart/pause/delete） ——

 func (c *APIAppController) Start(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	app, agent, ok := loadAppAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if app.Status == constants.APP_STATUS_RUNNING {
 		utils.APIError(ctx, "应用已经启动")
 		return
 	}
 	if err := agentStart(agent, app); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.AppService.ChangeAPPStatus(app, constants.APP_STATUS_RUNNING, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "更新应用状态失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_APP, app.ID, constants.ACTION_START)
 	utils.APIOK(ctx)
 }

 func (c *APIAppController) ReStart(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	app, agent, ok := loadAppAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if err := agentRestart(agent, app); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_APP, app.ID, constants.ACTION_RESTART)
 	utils.APIOK(ctx)
 }

 func (c *APIAppController) Pause(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	app, agent, ok := loadAppAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if err := agentStop(agent, app.ID); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.AppService.ChangeAPPStatus(app, constants.APP_STATUS_PAUSE, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "更新应用状态失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_APP, app.ID, constants.ACTION_PAUSE)
 	utils.APIOK(ctx)
 }

 func (c *APIAppController) Delete(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	app, agent, ok := loadAppAgentForUser(ctx, id)
 	if !ok {
 		return
 	}
 	if app.Status == constants.APP_STATUS_RUNNING {
 		utils.APIError(ctx, "应用启动状态不能删除")
 		return
 	}
 	if err := agentStop(agent, app.ID); err != nil {
 		utils.APIError(ctx, err.Error())
 		return
 	}
 	if !providers.AppService.ChangeAPPStatus(app, constants.APP_STATUS_DELETED, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "删除应用失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_APP, app.ID, constants.ACTION_DELETE)
 	utils.APIOK(ctx)
 }

 // —— 批量 ——

 type apiBatchReq struct {
 	IDs []int64 `form:"ids" json:"ids" binding:"required"`
 }

 func (c *APIAppController) BatchAction(ctx *gin.Context, action string) {
 	var req apiBatchReq
 	if err := ctx.ShouldBind(&req); err != nil || len(req.IDs) == 0 {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	user := utils.GetUser(ctx)
 	for _, id := range req.IDs {
 		app := providers.AppService.GetAppByID(id)
 		if app == nil {
 			logger.Errorf("App Batch%s app missing id=%d", action, id)
 			continue
 		}
 		if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
 			logger.Errorf("App Batch%s permission id=%d", action, id)
 			continue
 		}
 		agent := providers.AgentService.GetAgentByID(app.AgentID)
 		if agent == nil || agent.Status != constants.AGENT_ONLINE {
 			logger.Errorf("App Batch%s agent unavailable id=%d", action, id)
 			continue
 		}
 		switch action {
 		case "start":
 			if app.Status == constants.APP_STATUS_RUNNING {
 				continue
 			}
 			if err := agentStart(agent, app); err != nil {
 				logger.Errorf("App BatchStart err:%s", err.Error())
 				continue
 			}
 			providers.AppService.ChangeAPPStatus(app, constants.APP_STATUS_RUNNING, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_APP, app.ID, constants.ACTION_START)
 		case "restart":
 			if err := agentRestart(agent, app); err != nil {
 				logger.Errorf("App BatchRestart err:%s", err.Error())
 				continue
 			}
 			utils.OpetationLog(user.ID, constants.TYPE_APP, app.ID, constants.ACTION_RESTART)
 		case "pause":
 			if err := agentStop(agent, app.ID); err != nil {
 				logger.Errorf("App BatchPause err:%s", err.Error())
 				continue
 			}
 			providers.AppService.ChangeAPPStatus(app, constants.APP_STATUS_PAUSE, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_APP, app.ID, constants.ACTION_PAUSE)
 		case "delete":
 			if app.Status == constants.APP_STATUS_RUNNING {
 				continue
 			}
 			if err := agentStop(agent, app.ID); err != nil {
 				logger.Errorf("App BatchDelete err:%s", err.Error())
 				continue
 			}
 			providers.AppService.ChangeAPPStatus(app, constants.APP_STATUS_DELETED, user.ID)
 			utils.OpetationLog(user.ID, constants.TYPE_APP, app.ID, constants.ACTION_DELETE)
 		}
 	}
 	utils.APIOK(ctx)
 }

 func (c *APIAppController) BatchStart(ctx *gin.Context)    { c.BatchAction(ctx, "start") }
 func (c *APIAppController) BatchReStart(ctx *gin.Context)  { c.BatchAction(ctx, "restart") }
 func (c *APIAppController) BatchPause(ctx *gin.Context)    { c.BatchAction(ctx, "pause") }
 func (c *APIAppController) BatchDelete(ctx *gin.Context)  { c.BatchAction(ctx, "delete") }

 // —— 共享 helper ——

 func loadAppAgentForUser(ctx *gin.Context, id int64) (*models.App, *models.Agent, bool) {
 	app := providers.AppService.GetAppByID(id)
 	if app == nil {
 		utils.APIBadRequest(ctx, "应用不存在")
 		return nil, nil, false
 	}
 	user := utils.GetUser(ctx)
 	if user.Role != constants.USER_ROLE_ADMIN && app.Creator != user.ID {
 		utils.APIBadRequest(ctx, "对不起，您没有操作此应用的权限")
 		return nil, nil, false
 	}
 	agent := providers.AgentService.GetAgentByID(app.AgentID)
 	if agent == nil {
 		utils.APIBadRequest(ctx, "实例获取失败")
 		return nil, nil, false
 	}
 	if agent.Status != constants.AGENT_ONLINE {
 		utils.APIBadRequest(ctx, "实例不处于运行状态")
 		return nil, nil, false
 	}
 	return app, agent, true
 }

 func agentStart(agent *models.Agent, app *models.App) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetApp(app.ID)
 	if err != nil {
 		return fmt.Errorf("获取应用情况异常:%s", err.Error())
 	}
 	if exist == nil {
 		if err := client.AddApp(app); err != nil {
 			return fmt.Errorf("添加应用异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func agentRestart(agent *models.Agent, app *models.App) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetApp(app.ID)
 	if err != nil {
 		return fmt.Errorf("获取应用情况异常:%s", err.Error())
 	}
 	if exist == nil {
 		if err := client.AddApp(app); err != nil {
 			return fmt.Errorf("重启应用异常:%s", err.Error())
 		}
 	} else {
 		if err := client.UpdateApp(app); err != nil {
 			return fmt.Errorf("重启应用异常:%s", err.Error())
 		}
 	}
 	return nil
 }

 func agentStop(agent *models.Agent, appID int64) error {
 	client, err := providers.GetAgent(agent)
 	if err != nil {
 		return fmt.Errorf("初始化RPC客户端异常:%s", err.Error())
 	}
 	exist, err := client.GetApp(appID)
 	if err != nil {
 		return fmt.Errorf("获取应用情况异常:%s", err.Error())
 	}
 	if exist != nil {
 		if err := client.RemoveApp(appID); err != nil {
 			return fmt.Errorf("停止应用异常:%s", err.Error())
 		}
 	}
 	return nil
 }
