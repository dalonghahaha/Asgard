 package controllers

 import (
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 type APIArchiveController struct{}

 func NewAPIArchiveController() *APIArchiveController {
 	return &APIArchiveController{}
 }

 func (c *APIArchiveController) App(ctx *gin.Context) {
 	id := utils.QueryInt64(ctx, "app_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "app_id 不能为空")
 		return
 	}
 	if providers.AppService.GetAppByID(id) == nil {
 		utils.APIBadRequest(ctx, "应用不存在")
 		return
 	}
 	page := utils.QueryInt(ctx, "page", 1)
 	where := map[string]interface{}{
 		"type":       constants.TYPE_APP,
 		"related_id": id,
 	}
 	list, total := providers.ArchiveService.GetArchivePageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.ArchiveFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 func (c *APIArchiveController) Job(ctx *gin.Context) {
 	id := utils.QueryInt64(ctx, "job_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "job_id 不能为空")
 		return
 	}
 	if providers.JobService.GetJobByID(id) == nil {
 		utils.APIBadRequest(ctx, "计划任务不存在")
 		return
 	}
 	page := utils.QueryInt(ctx, "page", 1)
 	where := map[string]interface{}{
 		"type":       constants.TYPE_JOB,
 		"related_id": id,
 	}
 	list, total := providers.ArchiveService.GetArchivePageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.ArchiveFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 func (c *APIArchiveController) Timing(ctx *gin.Context) {
 	id := utils.QueryInt64(ctx, "timing_id", 0)
 	if id == 0 {
 		utils.APIBadRequest(ctx, "timing_id 不能为空")
 		return
 	}
 	if providers.TimingService.GetTimingByID(id) == nil {
 		utils.APIBadRequest(ctx, "定时任务不存在")
 		return
 	}
 	page := utils.QueryInt(ctx, "page", 1)
 	where := map[string]interface{}{
 		"type":       constants.TYPE_TIMING,
 		"related_id": id,
 	}
 	list, total := providers.ArchiveService.GetArchivePageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.ArchiveFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 // APIExceptionController 异常记录查询。
 type APIExceptionController struct{}

 func NewAPIExceptionController() *APIExceptionController {
 	return &APIExceptionController{}
 }

 func (c *APIExceptionController) List(ctx *gin.Context) {
 	page := utils.QueryInt(ctx, "page", 1)
 	_type := utils.QueryInt64(ctx, "type", 0)
 	where := map[string]interface{}{}
 	if _type > 0 {
 		where["type"] = _type
 	}
 	list, total := providers.ExceptionService.GetExceptionPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.ExceptionFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 // APIOperationController 操作日志查询。
 type APIOperationController struct{}

 func NewAPIOperationController() *APIOperationController {
 	return &APIOperationController{}
 }

 func (c *APIOperationController) List(ctx *gin.Context) {
 	page := utils.QueryInt(ctx, "page", 1)
 	userID := utils.QueryInt64(ctx, "user_id", 0)
 	_type := utils.QueryInt64(ctx, "type", 0)
 	where := map[string]interface{}{}
 	if userID > 0 {
 		where["user_id"] = userID
 	}
 	if _type > 0 {
 		where["type"] = _type
 	}
 	list, total := providers.OperationService.GetOperationPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		out = append(out, utils.OperationFormat(&list[i]))
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }
