 package controllers

 import (
 	"strconv"

 	"github.com/dalonghahaha/avenger/tools/coding"
 	"github.com/dalonghahaha/avenger/tools/random"
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 type APIUserController struct{}

 func NewAPIUserController() *APIUserController {
 	return &APIUserController{}
 }

 // GET /api/v1/users
 func (c *APIUserController) List(ctx *gin.Context) {
 	page := utils.QueryInt(ctx, "page", 1)
 	status := utils.QueryInt(ctx, "status", -99)
 	nickname := ctx.Query("nickname")
 	phone := ctx.Query("phone")
 	email := ctx.Query("email")
 	where := map[string]interface{}{"status": status}
 	if nickname != "" {
 		where["nickname"] = nickname
 	}
 	if phone != "" {
 		where["phone"] = phone
 	}
 	if email != "" {
 		where["email"] = email
 	}
 	list, total := providers.UserService.GetUserPageList(where, page, constants.WEB_LIST_PAGE_SIZE)
 	out := make([]gin.H, 0, len(list))
 	for i := range list {
 		u := list[i]
 		out = append(out, gin.H{
 			"id":         u.ID,
 			"nickname":   u.NickName,
 			"email":      u.Email,
 			"mobile":     u.Mobile,
 			"role":       u.Role,
 			"status":     u.Status,
 			"avatar":     u.Avatar,
 			"created_at": u.CreatedAt,
 		})
 	}
 	utils.APIPage(ctx, out, total, page, constants.WEB_LIST_PAGE_SIZE)
 }

 // GET /api/v1/users/:id
 func (c *APIUserController) Show(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	user := providers.UserService.GetUserByID(id)
 	if user == nil {
 		utils.APIBadRequest(ctx, "用户不存在")
 		return
 	}
 	utils.APIData(ctx, gin.H{
 		"id":         user.ID,
 		"nickname":   user.NickName,
 		"email":      user.Email,
 		"mobile":     user.Mobile,
 		"role":       user.Role,
 		"status":     user.Status,
 		"avatar":     user.Avatar,
 		"created_at": user.CreatedAt,
 	})
 }

 type apiUserCreateReq struct {
 	Nickname string `form:"nickname" json:"nickname" binding:"required"`
 	Email    string `form:"email" json:"email" binding:"required"`
 	Mobile   string `form:"mobile" json:"mobile" binding:"required"`
 	Password string `form:"password" json:"password" binding:"required"`
 	Role     string `form:"role" json:"role"`
 	Status   int64  `form:"status" json:"status"`
 }

 // POST /api/v1/users
 func (c *APIUserController) Create(ctx *gin.Context) {
 	var req apiUserCreateReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	if !utils.EmailFormat(req.Email) {
 		utils.APIBadRequest(ctx, "邮箱格式不正确")
 		return
 	}
 	if !utils.MobileFormat(req.Mobile) {
 		utils.APIBadRequest(ctx, "手机号格式不正确")
 		return
 	}
 	if providers.UserService.GetUserByNickName(req.Nickname) != nil {
 		utils.APIBadRequest(ctx, "该昵称已经注册")
 		return
 	}
 	if providers.UserService.GetUserByEmail(req.Email) != nil {
 		utils.APIBadRequest(ctx, "该邮箱已经注册")
 		return
 	}
 	if providers.UserService.GetUserByMobile(req.Mobile) != nil {
 		utils.APIBadRequest(ctx, "该手机号已经注册")
 		return
 	}
 	salt := random.Letters(8)
 	pwd, err := coding.MD5(req.Password + "|" + salt)
 	if err != nil {
 		utils.APIError(ctx, "生成密码失败")
 		return
 	}
 	user := new(models.User)
 	user.NickName = req.Nickname
 	user.Email = req.Email
 	user.Mobile = req.Mobile
 	user.Salt = salt
 	user.Password = pwd
 	if req.Role == constants.USER_ROLE_ADMIN {
 		user.Role = constants.USER_ROLE_ADMIN
 	} else {
 		user.Role = constants.USER_ROLE_NORMAL
 	}
 	if req.Status == constants.USER_STATUS_FORBIDDEN {
 		user.Status = constants.USER_STATUS_FORBIDDEN
 	} else {
 		user.Status = constants.USER_STATUS_NORMAL
 	}
 	if !providers.UserService.CreateUser(user) {
 		utils.APIError(ctx, "创建用户失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_USER, user.ID, constants.ACTION_CREATE)
 	utils.APIOK(ctx)
 }

 type apiUserUpdateReq struct {
 	Nickname string `form:"nickname" json:"nickname"`
 	Email    string `form:"email" json:"email"`
 	Mobile   string `form:"mobile" json:"mobile"`
 	Role     string `form:"role" json:"role"`
 	Status   int64  `form:"status" json:"status"`
 }

 // PUT /api/v1/users/:id
 func (c *APIUserController) Update(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	user := providers.UserService.GetUserByID(id)
 	if user == nil {
 		utils.APIBadRequest(ctx, "用户不存在")
 		return
 	}
 	var req apiUserUpdateReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	if req.Nickname != "" {
 		user.NickName = req.Nickname
 	}
 	if req.Email != "" {
 		user.Email = req.Email
 	}
 	if req.Mobile != "" {
 		user.Mobile = req.Mobile
 	}
 	if req.Role == constants.USER_ROLE_ADMIN || req.Role == constants.USER_ROLE_NORMAL {
 		user.Role = req.Role
 	}
 	if req.Status != 0 {
 		user.Status = req.Status
 	}
 	if !providers.UserService.UpdateUser(user) {
 		utils.APIError(ctx, "更新用户失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_USER, user.ID, constants.ACTION_UPDATE)
 	utils.APIOK(ctx)
 }

 // POST /api/v1/users/:id/forbidden
 func (c *APIUserController) Forbidden(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	user := providers.UserService.GetUserByID(id)
 	if user == nil {
 		utils.APIBadRequest(ctx, "用户不存在")
 		return
 	}
 	user.Status = constants.USER_STATUS_FORBIDDEN
 	if !providers.UserService.UpdateUser(user) {
 		utils.APIError(ctx, "禁用用户失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_USER, user.ID, constants.ACTION_DELETE)
 	utils.APIOK(ctx)
 }

 type apiResetPasswordReq struct {
 	Password string `form:"password" json:"password" binding:"required"`
 }

 // POST /api/v1/users/:id/reset_password
 func (c *APIUserController) ResetPassword(ctx *gin.Context) {
 	id := utils.PathInt64(ctx, "id")
 	user := providers.UserService.GetUserByID(id)
 	if user == nil {
 		utils.APIBadRequest(ctx, "用户不存在")
 		return
 	}
 	var req apiResetPasswordReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	salt := random.Letters(8)
 	pwd, err := coding.MD5(req.Password + "|" + salt)
 	if err != nil {
 		utils.APIError(ctx, "生成密码失败")
 		return
 	}
 	user.Salt = salt
 	user.Password = pwd
 	if !providers.UserService.UpdateUser(user) {
 		utils.APIError(ctx, "重置密码失败")
 		return
 	}
 	utils.OpetationLog(utils.GetUserID(ctx), constants.TYPE_USER, user.ID, constants.ACTION_UPDATE)
 	utils.APIOK(ctx)
 }

 // 避免 strconv 未使用告警
 var _ = strconv.Itoa
