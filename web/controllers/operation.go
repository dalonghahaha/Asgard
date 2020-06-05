package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
	"Asgard/web/utils"
)

type OperationController struct {
}

func NewOperationController() *OperationController {
	return &OperationController{}
}

func (c *OperationController) List(ctx *gin.Context) {
	page := utils.DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	operationList, total := providers.OperationService.GetOperationPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
	if operationList == nil {
		utils.JumpWarning(ctx, "获取操作记录失败")
	}
	list := []gin.H{}
	for _, operatio := range operationList {
		list = append(list, utils.OperationFormat(&operatio))
	}
	mpurl := "/operation/list"
	utils.Render(ctx, "operation/list", gin.H{
		"Subtitle":   "操作记录",
		"List":       list,
		"Total":      total,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}
