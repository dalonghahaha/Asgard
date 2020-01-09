package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"Asgard/models"
	"Asgard/services"
)

type GroupController struct {
	groupService *services.GroupService
}

func NewGroupController() *GroupController {
	return &GroupController{
		groupService: services.NewGroupService(),
	}
}

func (c *GroupController) List(ctx *gin.Context) {
	page := DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	list, total := c.groupService.GetGroupPageList(where, page, PageSize)
	mpurl := "/group/list"
	ctx.HTML(StatusOK, "group/list", gin.H{
		"Subtitle":   "分组列表",
		"List":       list,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
}

func (c *GroupController) Add(ctx *gin.Context) {
	ctx.HTML(StatusOK, "group/add", gin.H{
		"Subtitle": "添加分组",
	})
}

func (c *GroupController) Create(ctx *gin.Context) {
	name := ctx.PostForm("name")
	if !Required(ctx, &name, "名称不能为空") {
		return
	}
	group := new(models.Group)
	group.Name = name
	group.Status = 0
	group.Creator = GetUserID(ctx)
	ok := c.groupService.CreateGroup(group)
	if !ok {
		APIError(ctx, "创建分组失败")
		return
	}
	APIOK(ctx)
}

func (c *GroupController) Edit(ctx *gin.Context) {
	id := DefaultInt(ctx, "id", 0)
	if id == 0 {
		JumpError(ctx)
		return
	}
	group := c.groupService.GetGroupByID(int64(id))
	if group == nil {
		JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "group/edit", gin.H{
		"Subtitle": "编辑分组",
		"Group":    group,
	})
}

func (c *GroupController) Update(ctx *gin.Context) {
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
	group := c.groupService.GetGroupByID(int64(id))
	if group == nil {
		APIBadRequest(ctx, "分组不存在")
		return
	}
	if name != "" {
		group.Name = name
	}
	if status != "" {
		_status, err := strconv.ParseInt(status, 10, 64)
		if err != nil {
			APIBadRequest(ctx, "status格式错误")
			return
		}
		group.Status = _status
	}
	group.Updator = GetUserID(ctx)
	ok := c.groupService.UpdateGroup(group)
	if !ok {
		APIError(ctx, "创建分组失败")
		return
	}
	APIOK(ctx)
}

func (c *GroupController) Delete(ctx *gin.Context) {
	id := FormDefaultInt(ctx, "id", 0)
	if id == 0 {
		APIBadRequest(ctx, "ID格式错误")
		return
	}
	ok := c.groupService.DeleteGroupByID(int64(id))
	if !ok {
		APIError(ctx, "删除分组失败")
		return
	}
	APIOK(ctx)
}
