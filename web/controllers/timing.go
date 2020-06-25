package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/web/utils"
)

type TimingController struct {
}

func NewTimingController() *TimingController {
	return &TimingController{}
}

func (c *TimingController) List(ctx *gin.Context) {
	groupID := utils.DefaultInt(ctx, "group_id", 0)
	agentID := utils.DefaultInt(ctx, "agent_id", 0)
	status := utils.DefaultInt(ctx, "status", -99)
	name := ctx.Query("name")
	page := utils.DefaultInt(ctx, "page", 1)
	user := utils.GetUser(ctx)
	where := map[string]interface{}{
		"status": status,
	}
	querys := []string{}
	if user.Role != constants.USER_ROLE_ADMIN {
		where["creator"] = user.ID
	}
	if groupID != 0 {
		where["group_id"] = groupID
		querys = append(querys, "group_id="+strconv.Itoa(groupID))
	}
	if agentID != 0 {
		where["agent_id"] = agentID
		querys = append(querys, "agent_id="+strconv.Itoa(agentID))
	}
	if status != -99 {
		querys = append(querys, "status="+strconv.Itoa(status))
	}
	if name != "" {
		where["name"] = name
		querys = append(querys, "name="+name)
	}
	timingList, total := providers.TimingService.GetTimingPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
	if timingList == nil {
		utils.APIError(ctx, "获取定时任务列表失败")
	}
	list := []map[string]interface{}{}
	for _, timing := range timingList {
		list = append(list, utils.TimingFormat(&timing))
	}
	mpurl := "/timing/list"
	if len(querys) > 0 {
		mpurl = "/timing/list?" + strings.Join(querys, "&")
	}
	utils.Render(ctx, "timing/list", gin.H{
		"Subtitle":   "定时任务列表",
		"List":       list,
		"Total":      total,
		"GroupList":  providers.GroupService.GetUsageGroup(),
		"AgentList":  providers.AgentService.GetUsageAgent(),
		"StatusList": constants.TIMING_STATUS,
		"GroupID":    groupID,
		"AgentID":    agentID,
		"Name":       name,
		"Status":     status,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *TimingController) Show(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	utils.Render(ctx, "timing/show", gin.H{
		"Subtitle": "查看定时任务",
		"Timing":   utils.TimingFormat(timing),
	})
}

func (c *TimingController) Add(ctx *gin.Context) {
	utils.Render(ctx, "timing/add", gin.H{
		"Subtitle":   "添加定时任务",
		"OutBaseDir": constants.WEB_OUT_DIR + "timer/",
		"GroupList":  providers.GroupService.GetUsageGroup(),
		"AgentList":  providers.AgentService.GetUsageAgent(),
	})
}

func (c *TimingController) Create(ctx *gin.Context) {
	if utils.FormDefaultInt64(ctx, "agent_id", 0) == 0 {
		utils.APIError(ctx, "运行实例未选择")
		return
	}
	timing := new(models.Timing)
	timing.GroupID = utils.FormDefaultInt64(ctx, "group_id", 0)
	timing.AgentID = utils.FormDefaultInt64(ctx, "agent_id", 0)
	timing.Name = ctx.PostForm("name")
	timing.Dir = ctx.PostForm("dir")
	timing.Program = ctx.PostForm("program")
	timing.Args = ctx.PostForm("args")
	timing.StdOut = ctx.PostForm("std_out")
	timing.StdErr = ctx.PostForm("std_err")
	timing.Time, _ = utils.ParseTime(ctx.PostForm("time"))
	timing.Timeout = utils.FormDefaultInt64(ctx, "timeout", -1)
	timing.Status = constants.TIMING_STATUS_PAUSE
	timing.Creator = utils.GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		timing.IsMonitor = 1
	}
	ok := providers.TimingService.CreateTiming(timing)
	if !ok {
		utils.APIError(ctx, "创建定时任务失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_CREATE)
	utils.APIOK(ctx)
}

func (c *TimingController) Edit(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	utils.Render(ctx, "timing/edit", gin.H{
		"Subtitle":  "编辑定时任务",
		"BackUrl":   utils.GetReferer(ctx),
		"Info":      utils.TimingFormat(timing),
		"GroupList": providers.GroupService.GetUsageGroup(),
		"AgentList": providers.AgentService.GetUsageAgent(),
	})
}

func (c *TimingController) Update(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	timing.GroupID = utils.FormDefaultInt64(ctx, "group_id", 0)
	timing.Name = ctx.PostForm("name")
	timing.Dir = ctx.PostForm("dir")
	timing.Program = ctx.PostForm("program")
	timing.Args = ctx.PostForm("args")
	timing.StdOut = ctx.PostForm("std_out")
	timing.StdErr = ctx.PostForm("std_err")
	timing.Time, _ = utils.ParseTime(ctx.PostForm("time"))
	timing.Timeout = utils.FormDefaultInt64(ctx, "timeout", -1)
	timing.Creator = utils.GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		timing.IsMonitor = 1
	}
	if utils.FormDefaultInt64(ctx, "agent_id", 0) != 0 {
		timing.AgentID = utils.FormDefaultInt64(ctx, "agent_id", 0)
	}
	ok := providers.TimingService.UpdateTiming(timing)
	if !ok {
		utils.APIError(ctx, "更新定时任务失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_UPDATE)
	utils.APIOK(ctx)
}

func (c *TimingController) Copy(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	_timing := new(models.Timing)
	_timing.GroupID = timing.GroupID
	_timing.Name = timing.Name + "_copy"
	_timing.AgentID = timing.AgentID
	_timing.Dir = timing.Dir
	_timing.Program = timing.Program
	_timing.Args = timing.Args
	_timing.StdOut = timing.StdOut
	_timing.StdErr = timing.StdErr
	_timing.Time = timing.Time
	_timing.Timeout = timing.Timeout
	_timing.IsMonitor = timing.IsMonitor
	_timing.Status = constants.TIMING_STATUS_PAUSE
	_timing.Creator = utils.GetUserID(ctx)
	ok := providers.TimingService.CreateTiming(_timing)
	if !ok {
		utils.APIError(ctx, "复制定时任务失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_COPY)
	utils.APIOK(ctx)
}

func (c *TimingController) Start(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	if timing.Status == constants.TIMING_STATUS_RUNNING {
		utils.APIError(ctx, "定时任务已经启动")
		return
	}
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_timing, err := client.GetTiming(timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing == nil {
		err = client.AddTiming(timing)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("添加定时任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_RUNNING, utils.GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新定时任务状态失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_START)
	utils.APIOK(ctx)
}

func (c *TimingController) ReStart(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_timing, err := client.GetTiming(timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing == nil {
		err = client.AddTiming(timing)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("重启定时任务异常:%s", err.Error()))
			return
		}
	} else {
		err = client.UpdateTiming(timing)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("重启定时任务异常:%s", err.Error()))
			return
		}
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_RESTART)
	utils.APIOK(ctx)
}

func (c *TimingController) Pause(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_timing, err := client.GetTiming(timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing != nil {
		err = client.RemoveTiming(timing.ID)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("停止定时任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_PAUSE, utils.GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新定时任务状态失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_PAUSE)
	utils.APIOK(ctx)
}

func (c *TimingController) Delete(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	if timing.Status != constants.TIMING_STATUS_PAUSE {
		utils.APIError(ctx, "定时任务启动状态不能删除")
		return
	}
	client, err := providers.GetAgent(agent)
	if err != nil {
		utils.APIError(ctx, "初始化RPC客户端异常:\n"+err.Error())
		return
	}
	_timing, err := client.GetTiming(timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing != nil {
		err = client.RemoveTiming(timing.ID)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("停止定时任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_DELETED, utils.GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新定时任务状态失败")
		return
	}
	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_DELETE)
	utils.APIOK(ctx)
}

func (c *TimingController) BatchStart(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		if timing.Status == constants.TIMING_STATUS_RUNNING {
			continue
		}
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("Timing BatchStart GetAgent Error:[%d][%s]", timing.ID, err.Error())
			continue
		}
		_timing, err := client.GetTiming(timing.ID)
		if err != nil {
			logger.Errorf("Timing BatchStart GetAgentTiming Error:[%d][%s]", timing.ID, err.Error())
			continue
		}
		if _timing == nil {
			err = client.AddTiming(timing)
			if err != nil {
				logger.Errorf("Timing BatchStart AddAgentTiming Error:%s", err.Error())
			}
		}
		providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_RUNNING, utils.GetUserID(ctx))
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_START)
	}
	utils.APIOK(ctx)
}

func (c *TimingController) BatchReStart(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("Timing BatchReStart GetAgent Error:[%d][%s]", timing.ID, err.Error())
			continue
		}
		_timing, err := client.GetTiming(timing.ID)
		if err != nil {
			logger.Errorf("Timing BatchReStart GetAgentTiming Error:[%d][%s]", timing.ID, err.Error())
			continue
		}
		if _timing == nil {
			err = client.AddTiming(timing)
			if err != nil {
				logger.Errorf("Timing BatchReStart AddAgentTiming Error:[%d][%s]", timing.ID, err.Error())
			}
		} else {
			err = client.UpdateTiming(timing)
			if err != nil {
				logger.Errorf("Timing BatchReStart UpdateAgentJob Error:[%d][%s]", timing.ID, err.Error())
			}
		}
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_RESTART)
	}
	utils.APIOK(ctx)
}

func (c *TimingController) BatchPause(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("Timing BatchPause GetAgent Error:[%d][%s]", timing.ID, err.Error())
		}
		_timing, err := client.GetTiming(timing.ID)
		if err != nil {
			logger.Errorf("Timing BatchPause GetAgentTiming Error:[%d][%s]", timing.ID, err.Error())
			continue
		}
		if _timing != nil {
			err = client.RemoveTiming(timing.ID)
			if err != nil {
				logger.Errorf("Timing BatchPause RemoveAgentTiming Error:[%d][%s]", timing.ID, err.Error())
				return
			}
		}
		providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_PAUSE, utils.GetUserID(ctx))
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_PAUSE)
	}
	utils.APIOK(ctx)
}

func (c *TimingController) BatchDelete(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		if timing.Status == constants.TIMING_STATUS_RUNNING {
			continue
		}
		client, err := providers.GetAgent(agent)
		if err != nil {
			logger.Errorf("Timing BatchDelete GetAgent Error:[%d][%s]", timing.ID, err.Error())
		}
		_timing, err := client.GetTiming(timing.ID)
		if err != nil {
			logger.Errorf("Timing BatchDelete GetAgentTiming Error:[%d][%s]", timing.ID, err.Error())
			continue
		}
		if _timing != nil {
			err = client.RemoveTiming(timing.ID)
			if err != nil {
				logger.Errorf("Timing BatchDelete RemoveAgentTiming Error:[%d][%s]", timing.ID, err.Error())
				return
			}
		}
		providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_DELETED, utils.GetUserID(ctx))
		utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_TIMING, timing.ID, constants.ACTION_DELETE)
	}
	utils.APIOK(ctx)
}
