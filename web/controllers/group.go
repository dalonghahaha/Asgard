package controllers

import (
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
	where := map[string]interface{}{}
	list, total := providers.GroupService.GetGroupPageList(where, page, PageSize)
	mpurl := "/group/list"
	ctx.HTML(StatusOK, "group/list", gin.H{
		"Subtitle":   "分组列表",
		"List":       list,
		"Total":      total,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *GroupController) Add(ctx *gin.Context) {
	ctx.HTML(StatusOK, "group/add", gin.H{
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
	group.Creator = GetUserID(ctx)
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
	ctx.HTML(StatusOK, "group/edit", gin.H{
		"Subtitle": "编辑分组",
		"Group":    group,
	})
}

func (c *GroupController) Update(ctx *gin.Context) {
	id := utils.FormDefaultInt64(ctx, "id", 0)
	name := ctx.PostForm("name")
	status := ctx.PostForm("status")
	if id == 0 {
		utils.APIBadRequest(ctx, "ID格式错误")
		return
	}
	group := providers.GroupService.GetGroupByID(id)
	if group == nil {
		utils.APIBadRequest(ctx, "分组不存在")
		return
	}
	group.Name = name
	group.Updator = GetUserID(ctx)
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
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.APIBadRequest(ctx, "ID格式错误")
		return
	}
	group := providers.GroupService.GetGroupByID(id)
	if group == nil {
		utils.APIBadRequest(ctx, "分组不存在")
		return
	}
	ok := providers.GroupService.DeleteGroupByID(id)
	if !ok {
		utils.APIError(ctx, "删除分组失败")
		return
	}
	utils.APIOK(ctx)
}
