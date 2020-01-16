package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"Asgard/client"
	"Asgard/models"
	"Asgard/services"
)

type AppController struct {
	appService     *services.AppService
	agentService   *services.AgentService
	groupService   *services.GroupService
	moniterService *services.MonitorService
	archiveService *services.ArchiveService
}

func NewAppController() *AppController {
	return &AppController{
		appService:     services.NewAppService(),
		agentService:   services.NewAgentService(),
		groupService:   services.NewGroupService(),
		moniterService: services.NewMonitorService(),
		archiveService: services.NewArchiveService(),
	}
}

func (c *AppController) formatApp(info *models.App) map[string]interface{} {
	data := map[string]interface{}{
		"ID":          info.ID,
		"Name":        info.Name,
		"GroupID":     info.GroupID,
		"AgentID":     info.AgentID,
		"Dir":         info.Dir,
		"Program":     info.Program,
		"Args":        info.Args,
		"StdOut":      info.StdOut,
		"StdErr":      info.StdErr,
		"AutoRestart": info.AutoRestart,
		"IsMonitor":   info.IsMonitor,
		"Status":      info.Status,
	}
	group := c.groupService.GetGroupByID(info.GroupID)
	if group != nil {
		data["GroupName"] = group.Name
	} else {
		data["GroupName"] = ""
	}
	agent := c.agentService.GetAgentByID(info.AgentID)
	if agent != nil {
		data["AgentName"] = agent.IP + ":" + agent.Port
	} else {
		data["AgentName"] = ""
	}
	return data
}

func (c *AppController) List(ctx *gin.Context) {
	agent := DefaultInt(ctx, "agent", 0)
	page := DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	if agent != 0 {
		where["agent_id"] = agent
	}
	appList, total := c.appService.GetAppPageList(where, page, PageSize)
	if appList == nil {
		APIError(ctx, "获取应用列表失败")
	}
	list := []map[string]interface{}{}
	for _, app := range appList {
		list = append(list, c.formatApp(&app))
	}
	mpurl := "/app/list"
	if agent != 0 {
		mpurl = "/app/list?agent=" + strconv.Itoa(agent)
	}
	ctx.HTML(StatusOK, "app/list", gin.H{
		"Subtitle":   "应用列表",
		"List":       list,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
}

func (c *AppController) Show(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		JumpError(ctx)
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "app/show", gin.H{
		"Subtitle": "查看应用",
		"App":      c.formatApp(app),
	})
}

func (c *AppController) Monitor(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		JumpError(ctx)
		return
	}
	cpus := []string{}
	memorys := []string{}
	times := []string{}
	moniters := c.moniterService.GetAppMonitor(id, 100)
	for _, moniter := range moniters {
		cpus = append(cpus, FormatFloat(moniter.CPU))
		memorys = append(memorys, FormatFloat(moniter.Memory))
		times = append(times, FormatTime(moniter.CreatedAt))
	}
	ctx.HTML(StatusOK, "app/monitor", gin.H{
		"Subtitle": "监控信息",
		"CPU":      cpus,
		"Memory":   memorys,
		"Time":     times,
	})
}

func (c *AppController) Archive(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	page := DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{
		"type":       models.TYPE_APP,
		"related_id": id,
	}
	if id == 0 {
		JumpError(ctx)
		return
	}
	archiveList, total := c.archiveService.GetArchivePageList(where, page, PageSize)
	if archiveList == nil {
		APIError(ctx, "获取归档列表失败")
	}
	list := []map[string]interface{}{}
	for _, archive := range archiveList {
		list = append(list, formatArchive(&archive))
	}
	mpurl := fmt.Sprintf("/app/archive?id=%d", id)
	ctx.HTML(StatusOK, "app/archive", gin.H{
		"Subtitle":   "归档列表",
		"List":       list,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
}

func (c *AppController) Add(ctx *gin.Context) {
	ctx.HTML(StatusOK, "app/add", gin.H{
		"Subtitle":  "添加应用",
		"GroupList": c.groupService.GetUsageGroup(),
		"AgentList": c.agentService.GetUsageAgent(),
	})
}

func (c *AppController) Create(ctx *gin.Context) {
	groupID := FormDefaultInt(ctx, "group_id", 0)
	name := ctx.PostForm("name")
	agentID := FormDefaultInt(ctx, "agent_id", 0)
	dir := ctx.PostForm("dir")
	program := ctx.PostForm("program")
	args := ctx.PostForm("args")
	stdOut := ctx.PostForm("std_out")
	stdErr := ctx.PostForm("std_err")
	autoRestart := ctx.PostForm("auto_restart")
	isMonitor := ctx.PostForm("is_monitor")
	if !Required(ctx, &name, "名称不能为空") {
		return
	}
	if !Required(ctx, &dir, "执行目录不能为空") {
		return
	}
	if !Required(ctx, &program, "执行程序不能为空") {
		return
	}
	if !Required(ctx, &stdOut, "标准输出路径不能为空") {
		return
	}
	if !Required(ctx, &stdErr, "错误输出路径不能为空") {
		return
	}
	if agentID == 0 {
		APIBadRequest(ctx, "运行实例不能为空")
		return
	}
	app := new(models.App)
	app.GroupID = int64(groupID)
	app.Name = name
	app.AgentID = int64(agentID)
	app.Dir = dir
	app.Program = program
	app.Args = args
	app.StdOut = stdOut
	app.StdErr = stdErr
	app.Status = 0
	app.Creator = GetUserID(ctx)
	if autoRestart != "" {
		app.AutoRestart = 1
	}
	if isMonitor != "" {
		app.IsMonitor = 1
	}
	ok := c.appService.CreateApp(app)
	if !ok {
		APIError(ctx, "创建应用失败")
		return
	}
	APIOK(ctx)
}

func (c *AppController) Edit(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		JumpError(ctx)
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "app/edit", gin.H{
		"Subtitle":  "编辑分组",
		"App":       c.formatApp(app),
		"GroupList": c.groupService.GetUsageGroup(),
		"AgentList": c.agentService.GetUsageAgent(),
	})
}

func (c *AppController) Update(ctx *gin.Context) {
	id := FormDefaultInt(ctx, "id", 0)
	groupID := FormDefaultInt(ctx, "group_id", 0)
	name := ctx.PostForm("name")
	agentID := FormDefaultInt(ctx, "agent_id", 0)
	dir := ctx.PostForm("dir")
	program := ctx.PostForm("program")
	args := ctx.PostForm("args")
	stdOut := ctx.PostForm("std_out")
	stdErr := ctx.PostForm("std_err")
	autoRestart := ctx.PostForm("auto_restart")
	isMonitor := ctx.PostForm("is_monitor")
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	if !Required(ctx, &name, "名称不能为空") {
		return
	}
	if !Required(ctx, &dir, "执行目录不能为空") {
		return
	}
	if !Required(ctx, &program, "执行程序不能为空") {
		return
	}
	if !Required(ctx, &stdOut, "标准输出路径不能为空") {
		return
	}
	if !Required(ctx, &stdErr, "错误输出路径不能为空") {
		return
	}
	if agentID == 0 {
		APIBadRequest(ctx, "运行实例不能为空")
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		APIBadRequest(ctx, "应用不存在")
		return
	}
	app.GroupID = int64(groupID)
	app.Name = name
	app.AgentID = int64(agentID)
	app.Dir = dir
	app.Program = program
	app.Args = args
	app.StdOut = stdOut
	app.StdErr = stdErr
	app.Updator = GetUserID(ctx)
	if autoRestart != "" {
		app.AutoRestart = 1
	}
	if isMonitor != "" {
		app.IsMonitor = 1
	}
	ok := c.appService.UpdateApp(app)
	if !ok {
		APIError(ctx, "更新应用失败")
		return
	}
	APIOK(ctx)
}

func (c *AppController) Delete(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		APIError(ctx, "应用不存在")
		return
	}
	if app.Status == 1 {
		APIError(ctx, "应用正在运行不能删除")
		return
	}
	app.Status = -1
	app.Updator = GetUserID(ctx)
	ok := c.appService.UpdateApp(app)
	if !ok {
		APIError(ctx, "删除应用失败")
		return
	}
	APIOK(ctx)
}

func (c *AppController) Start(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		APIError(ctx, "应用不存在")
		return
	}
	agent := c.agentService.GetAgentByID(app.AgentID)
	if agent == nil {
		APIError(ctx, "应用对应实例获取异常")
		return
	}
	_app, err := client.GetAgentApp(agent, int64(id))
	if err != nil {
		APIError(ctx, fmt.Sprintf("获取应用情况异常:%s", err.Error()))
		return
	}
	if _app == nil {
		err = client.AddAgentApp(agent, app)
		if err != nil {
			APIError(ctx, fmt.Sprintf("添加应用异常:%s", err.Error()))
			return
		}
		APIOK(ctx)
	}
	APIOK(ctx)
}

func (c *AppController) ReStart(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		APIError(ctx, "应用不存在")
		return
	}
	agent := c.agentService.GetAgentByID(app.AgentID)
	if agent == nil {
		APIError(ctx, "应用对应实例获取异常")
		return
	}
	_app, err := client.GetAgentApp(agent, int64(id))
	if err != nil {
		APIError(ctx, fmt.Sprintf("获取应用情况异常:%s", err.Error()))
		return
	}
	if _app == nil {
		err = client.AddAgentApp(agent, app)
		if err != nil {
			APIError(ctx, fmt.Sprintf("重启异常:%s", err.Error()))
			return
		}
		APIOK(ctx)
	}
	err = client.UpdateAgentApp(agent, app)
	if err != nil {
		APIError(ctx, fmt.Sprintf("重启异常:%s", err.Error()))
		return
	}
	APIOK(ctx)
}

func (c *AppController) Stop(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		APIError(ctx, "应用不存在")
		return
	}
	agent := c.agentService.GetAgentByID(app.AgentID)
	if agent == nil {
		APIError(ctx, "应用对应实例获取异常")
		return
	}
	_app, err := client.GetAgentApp(agent, int64(id))
	if err != nil {
		APIError(ctx, fmt.Sprintf("获取应用情况异常:%s", err.Error()))
		return
	}
	if _app == nil {
		APIOK(ctx)
		return
	}
	err = client.RemoveAgentApp(agent, int64(id))
	if err != nil {
		APIError(ctx, fmt.Sprintf("停止应用异常:%s", err.Error()))
		return
	}
	app.Status = 0
	app.Updator = GetUserID(ctx)
	c.appService.UpdateApp(app)
	APIOK(ctx)
}
