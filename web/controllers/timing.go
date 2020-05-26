package controllers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/gin-gonic/gin"

	"Asgard/client"
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
	where := map[string]interface{}{
		"status": status,
	}
	querys := []string{}
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
	timingList, total := providers.TimingService.GetTimingPageList(where, page, PageSize)
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
	ctx.HTML(StatusOK, "timing/list", gin.H{
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
	ctx.HTML(StatusOK, "timing/show", gin.H{
		"Subtitle": "查看定时任务",
		"Timing":   utils.TimingFormat(timing),
	})
}

func (c *TimingController) Add(ctx *gin.Context) {
	ctx.HTML(StatusOK, "timing/add", gin.H{
		"Subtitle":   "添加定时任务",
		"OutBaseDir": OutDir + "timer/",
		"GroupList":  providers.GroupService.GetUsageGroup(),
		"AgentList":  providers.AgentService.GetUsageAgent(),
	})
}

func (c *TimingController) Create(ctx *gin.Context) {
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
	timing.Creator = GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		timing.IsMonitor = 1
	}
	ok := providers.TimingService.CreateTiming(timing)
	if !ok {
		utils.APIError(ctx, "创建定时任务失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *TimingController) Edit(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	ctx.HTML(StatusOK, "timing/edit", gin.H{
		"Subtitle":  "编辑定时任务",
		"BackUrl":   GetReferer(ctx),
		"Info":      utils.TimingFormat(timing),
		"GroupList": providers.GroupService.GetUsageGroup(),
		"AgentList": providers.AgentService.GetUsageAgent(),
	})
}

func (c *TimingController) Update(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
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
	timing.Creator = GetUserID(ctx)
	if ctx.PostForm("is_monitor") != "" {
		timing.IsMonitor = 1
	}
	ok := providers.TimingService.UpdateTiming(timing)
	if !ok {
		utils.APIError(ctx, "更新定时任务失败")
		return
	}
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
	_timing.Creator = GetUserID(ctx)
	ok := providers.TimingService.CreateTiming(_timing)
	if !ok {
		utils.APIError(ctx, "复制定时任务失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *TimingController) Start(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	if timing.Status == constants.TIMING_STATUS_RUNNING {
		utils.APIError(ctx, "定时任务已经启动")
		return
	}
	_timing, err := client.GetAgentTiming(agent, timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing == nil {
		err = client.AddAgentTiming(agent, timing)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("添加定时任务异常:%s", err.Error()))
			return
		}
		utils.APIOK(ctx)
		return
	}
	ok := providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_RUNNING, GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新定时任务状态失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *TimingController) ReStart(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	_timing, err := client.GetAgentTiming(agent, timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing == nil {
		err = client.AddAgentTiming(agent, timing)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("重启定时任务异常:%s", err.Error()))
			return
		}
	} else {
		err = client.UpdateAgentTiming(agent, timing)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("重启定时任务异常:%s", err.Error()))
			return
		}
	}
	utils.APIOK(ctx)
}

func (c *TimingController) Pause(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	_timing, err := client.GetAgentTiming(agent, timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing != nil {
		err = client.RemoveAgentTiming(agent, timing.ID)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("停止定时任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_PAUSE, GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新定时任务状态失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *TimingController) Delete(ctx *gin.Context) {
	timing := utils.GetTiming(ctx)
	agent := utils.GetAgent(ctx)
	if timing.Status != constants.TIMING_STATUS_PAUSE {
		utils.APIError(ctx, "定时任务启动状态不能删除")
		return
	}
	_timing, err := client.GetAgentTiming(agent, timing.ID)
	if err != nil {
		utils.APIError(ctx, fmt.Sprintf("获取定时任务情况异常:%s", err.Error()))
		return
	}
	if _timing != nil {
		err = client.RemoveAgentTiming(agent, timing.ID)
		if err != nil {
			utils.APIError(ctx, fmt.Sprintf("停止定时任务异常:%s", err.Error()))
			return
		}
	}
	ok := providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_DELETED, GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新定时任务状态失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *TimingController) BatchStart(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		if timing.Status == constants.TIMING_STATUS_RUNNING {
			continue
		}
		_timing, err := client.GetAgentTiming(agent, timing.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("Timing BatchStart GetAgentTiming Error:[%d][%s]", timing.ID, err.Error()))
			continue
		}
		if _timing == nil {
			err = client.AddAgentTiming(agent, timing)
			if err != nil {
				logger.Error(fmt.Sprintf("Timing BatchStart AddAgentTiming Error:%s", err.Error()))
			}
		}
		providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_RUNNING, GetUserID(ctx))
	}
	utils.APIOK(ctx)
}

func (c *TimingController) BatchReStart(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		_timing, err := client.GetAgentTiming(agent, timing.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("Timing BatchReStart GetAgentTiming Error:[%d][%s]", timing.ID, err.Error()))
			continue
		}
		if _timing == nil {
			err = client.AddAgentTiming(agent, timing)
			if err != nil {
				logger.Error(fmt.Sprintf("Timing BatchReStart AddAgentTiming Error:[%d][%s]", timing.ID, err.Error()))
			}
		} else {
			err = client.UpdateAgentTiming(agent, timing)
			if err != nil {
				logger.Error(fmt.Sprintf("Timing BatchReStart UpdateAgentJob Error:[%d][%s]", timing.ID, err.Error()))
			}
		}
	}
	utils.APIOK(ctx)
}

func (c *TimingController) BatchPause(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		_timing, err := client.GetAgentTiming(agent, timing.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("Timing BatchPause GetAgentTiming Error:[%d][%s]", timing.ID, err.Error()))
			continue
		}
		if _timing != nil {
			err = client.RemoveAgentTiming(agent, timing.ID)
			if err != nil {
				logger.Error(fmt.Sprintf("Timing BatchPause RemoveAgentTiming Error:[%d][%s]", timing.ID, err.Error()))
				return
			}
		}
		providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_PAUSE, GetUserID(ctx))
	}
	utils.APIOK(ctx)
}

func (c *TimingController) BatchDelete(ctx *gin.Context) {
	timingAgent := utils.GetTimingAgent(ctx)
	for timing, agent := range timingAgent {
		if timing.Status == constants.TIMING_STATUS_RUNNING {
			continue
		}
		_timing, err := client.GetAgentTiming(agent, timing.ID)
		if err != nil {
			logger.Error(fmt.Sprintf("Timing BatchDelete GetAgentTiming Error:[%d][%s]", timing.ID, err.Error()))
			continue
		}
		if _timing != nil {
			err = client.RemoveAgentTiming(agent, timing.ID)
			if err != nil {
				logger.Error(fmt.Sprintf("Timing BatchDelete RemoveAgentTiming Error:[%d][%s]", timing.ID, err.Error()))
				return
			}
		}
		providers.TimingService.ChangeTimingStatus(timing, constants.TIMING_STATUS_DELETED, GetUserID(ctx))
	}
	utils.APIOK(ctx)
}
