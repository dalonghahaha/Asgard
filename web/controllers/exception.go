package controllers

import (
	"github.com/gin-gonic/gin"

	"Asgard/constants"
	"Asgard/providers"
	"Asgard/web/utils"
)

type ExceptionController struct {
}

func NewExceptionController() *ExceptionController {
	return &ExceptionController{}
}

func (c *ExceptionController) List(ctx *gin.Context) {
	page := utils.DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	exceptionList, total := providers.ExceptionService.GetExceptionPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
	if exceptionList == nil {
		utils.JumpWarning(ctx, "获取异常信息失败")
	}
	list := []gin.H{}
	for _, exception := range exceptionList {
		list = append(list, utils.ExceptionFormat(&exception))
	}
	mpurl := "/exception/list"
	utils.Render(ctx, "exception/list", gin.H{
		"Subtitle":   "异常记录",
		"List":       list,
		"Total":      total,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}
