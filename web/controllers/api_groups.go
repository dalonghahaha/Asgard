 package controllers

 import (
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 type APIGroupController struct{}

 func NewAPIGroupController() *APIGroupController {
 	return &APIGroupController{}
 }

 func (c *APIGroupController) List(ctx *gin.Context) {
 	page := utils.QueryInt(ctx, "page", 1)
 	status := utils.QueryInt(ctx, "status", -99)
 	name := ctx.Query("name")
 	where := map[string]interface{}{"status": status}
 	if name != "" {
 		where["name"] = name
 	}
 	list, total := providers.GroupService.GetGroupPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		g := list[i]
 		out = append(out, gin.H{
 			"id":         g.ID,
 			"name":       g.Name,
 			"status":     g.Status,
 			"creator":    g.Creator,
 			"created_at": g.CreatedAt,
 		})
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 type apiGroupCreateReq struct {
 	Name   string `form:"name" json:"name" binding:"required"`
 	Status int64  `form:"status" json:"status"`
 }

 func (c *APIGroupController) Create(ctx *gin.Context) {
 	var req apiGroupCreateReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	group := new(models.Group)
 	group.Name = req.Name
 	group.Creator = utils.GetUserID(ctx)
 	if req.Status == constants.GROUP_STATUS_USAGE {
 		group.Status = constants.GROUP_STATUS_USAGE
 	} else {
 		group.Status = constants.GROUP_STATUS_UNUSAGE
 	}
 	if !providers.GroupService.CreateGroup(group) {
 		utils.APIError(ctx, "创建分组失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_GROUP, group.ID, constants.ACTION_CREATE)
 	utils.APIOK(ctx)
 }

 type apiGroupUpdateReq struct {
 	Name   string `form:"name" json:"name"`
 	Status int64  `form:"status" json:"status"`
 }

 func (c *APIGroupController) Update(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	group := providers.GroupService.GetGroupByID(id)
 	if group == nil {
 		utils.APIBadRequest(ctx, "分组不存在")
 		return
 	}
 	var req apiGroupUpdateReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	if req.Name != "" {
 		group.Name = req.Name
 	}
 	group.Updator = utils.GetUserID(ctx)
 	if req.Status == constants.GROUP_STATUS_USAGE {
 		group.Status = constants.GROUP_STATUS_USAGE
 	} else if req.Status == constants.GROUP_STATUS_DELETED {
 		group.Status = constants.GROUP_STATUS_DELETED
 	} else {
 		group.Status = constants.GROUP_STATUS_UNUSAGE
 	}
 	if !providers.GroupService.UpdateGroup(group) {
 		utils.APIError(ctx, "更新分组失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_GROUP, group.ID, constants.ACTION_UPDATE)
 	utils.APIOK(ctx)
 }

 func (c *APIGroupController) Delete(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	group := providers.GroupService.GetGroupByID(id)
 	if group == nil {
 		utils.APIBadRequest(ctx, "分组不存在")
 		return
 	}
 	if !providers.GroupService.ChangeGroupStatus(group, constants.GROUP_STATUS_DELETED, utils.GetUserID(ctx)) {
 		utils.APIError(ctx, "更新分组状态失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_GROUP, group.ID, constants.ACTION_DELETE)
 	utils.APIOK(ctx)
 }
