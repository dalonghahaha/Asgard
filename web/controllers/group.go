package controllers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/providers"
	"Asgard/web/utils"
)

type GroupController struct {
}

func NewGroupController() *GroupController {
	return &GroupController{}
}

func (c *GroupController) List(ctx *gin.Context) {
	page := utils.DefaultInt(ctx, "page", 1)
	status := utils.DefaultInt(ctx, "status", -99)
	name := ctx.Query("name")
	where := map[string]interface{}{
		"status": status,
	}
	querys := []string{}
	if name != "" {
		where["name"] = name
		querys = append(querys, "name="+name)
	}
	if status != -99 {
		querys = append(querys, "status="+strconv.Itoa(status))
	}
	list, total := providers.GroupService.GetGroupPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
	mpurl := "/group/list"
	if len(querys) > 0 {
		mpurl = "/group/list?" + strings.Join(querys, "&")
	}
	utils.Render(ctx, "group/list", gin.H{
		"Subtitle":   "分组列表",
		"List":       list,
		"Total":      total,
		"StatusList": constants.GROUP_STATUS,
		"Name":       name,
		"Status":     status,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *GroupController) Add(ctx *gin.Context) {
	utils.Render(ctx, "group/add", gin.H{
		"Subtitle": "添加分组",
	})
}

func (c *GroupController) Create(ctx *gin.Context) {
	name := ctx.PostForm("name")
	status := ctx.PostForm("status")
	if !utils.Required(ctx, name, "名称不能为空") {
		return
	}
	group := new(models.Group)
	group.Name = name
	group.Creator = utils.GetUserID(ctx)
	if status != "" {
		group.Status = constants.GROUP_STATUS_USAGE
	} else {
		group.Status = constants.GROUP_STATUS_UNUSAGE
	}
	ok := providers.GroupService.CreateGroup(group)
	if !ok {
		utils.APIError(ctx, "创建分组失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *GroupController) Edit(ctx *gin.Context) {
	id := utils.DefaultInt(ctx, "id", 0)
	if id == 0 {
		utils.JumpError(ctx)
		return
	}
	group := providers.GroupService.GetGroupByID(int64(id))
	if group == nil {
		utils.JumpError(ctx)
		return
	}
	utils.Render(ctx, "group/edit", gin.H{
		"Subtitle": "编辑分组",
		"Group":    group,
	})
}

func (c *GroupController) Update(ctx *gin.Context) {
	group := utils.GetGroup(ctx)
	name := ctx.PostForm("name")
	status := ctx.PostForm("status")
	group.Name = name
	group.Updator = utils.GetUserID(ctx)
	if status != "" {
		group.Status = constants.GROUP_STATUS_USAGE
	} else {
		group.Status = constants.GROUP_STATUS_UNUSAGE
	}
	ok := providers.GroupService.UpdateGroup(group)
	if !ok {
		utils.APIError(ctx, "更新分组失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *GroupController) Delete(ctx *gin.Context) {
	group := utils.GetGroup(ctx)
	ok := providers.GroupService.ChangeGroupStatus(group, constants.GROUP_STATUS_DELETED, utils.GetUserID(ctx))
	if !ok {
		utils.APIError(ctx, "更新分组状态失败")
		return
	}
	utils.APIOK(ctx)
}
