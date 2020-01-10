package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Asgard/models"
	"Asgard/services"
)

type AppController struct {
	appService   *services.AppService
	agentService *services.AgentService
	groupService *services.GroupService
}

func NewAppController() *AppController {
	return &AppController{
		appService:   services.NewAppService(),
		agentService: services.NewAgentService(),
		groupService: services.NewGroupService(),
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
	if group != nil {
		data["AgentName"] = agent.IP + ":" + agent.Port
	} else {
		data["AgentName"] = ""
	}
	return data
}

func (c *AppController) List(ctx *gin.Context) {
	page := DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	appList, total := c.appService.GetAppPageList(where, page, PageSize)
	if appList == nil {
		APIError(ctx, "获取应用列表失败")
	}
	list := []map[string]interface{}{}
	for _, app := range appList {
		list = append(list, c.formatApp(&app))
	}
	mpurl := "/app/list"
	ctx.HTML(StatusOK, "app/list", gin.H{
		"Subtitle":   "应用列表",
		"List":       list,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
}

func (c *AppController) Add(ctx *gin.Context) {
	groups := c.groupService.GetAllGroup()
	agents := c.agentService.GetAllAgent()
	ctx.HTML(StatusOK, "app/add", gin.H{
		"Subtitle":  "添加应用",
		"GroupList": groups,
		"AgentList": agents,
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
	if !Required(ctx, &name, "名称不能为空") {
		return
	}
	if !Required(ctx, &dir, "执行目录不能为空") {
		return
	}
	if !Required(ctx, &program, "执行程序不能为空") {
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
	groups := c.groupService.GetAllGroup()
	agents := c.agentService.GetAllAgent()
	ctx.HTML(StatusOK, "app/edit", gin.H{
		"Subtitle":  "编辑分组",
		"App":       c.formatApp(app),
		"GroupList": groups,
		"AgentList": agents,
	})
}

func (c *AppController) Update(ctx *gin.Context) {
	id := FormDefaultInt(ctx, "id", 0)
	name := ctx.PostForm("name")
	status := ctx.PostForm("status")
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	if name == "" && status == "" {
		APIBadRequest(ctx, "请求数据格式错误")
		return
	}
	app := c.appService.GetAppByID(int64(id))
	if app == nil {
		APIBadRequest(ctx, "分组不存在")
		return
	}
	if name != "" {
		app.Name = name
	}
	if status != "" {
		_status, err := strconv.ParseInt(status, 10, 64)
		if err != nil {
			APIBadRequest(ctx, "status格式错误")
			return
		}
		app.Status = _status
	}
	app.Updator = GetUserID(ctx)
	ok := c.appService.UpdateApp(app)
	if !ok {
		APIError(ctx, "更新应用失败")
		return
	}
	APIOK(ctx)
}

func (c *AppController) Delete(ctx *gin.Context) {
	id := FormDefaultInt(ctx, "id", 0)
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	ok := c.appService.DeleteAppByID(int64(id))
	if !ok {
		APIError(ctx, "删除应用失败")
		return
	}
	APIOK(ctx)
}
