 package controllers

 import (
 	"github.com/dalonghahaha/avenger/tools/coding"
 	"github.com/dalonghahaha/avenger/tools/random"
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/models"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 type APIAuthController struct{}

 func NewAPIAuthController() *APIAuthController {
 	return &APIAuthController{}
 }

 type apiLoginReq struct {
 	Username string `form:"username" json:"username" binding:"required"`
 	Password string `form:"password" json:"password" binding:"required"`
 }

 // Login 登录：同时下发 DES cookie（兼容旧 HTML 路由）和 JWT。
 func (c *APIAuthController) Login(ctx *gin.Context) {
 	var req apiLoginReq
 	if err := ctx.ShouldBind(&req); err != nil {
 		utils.APIBadRequest(ctx, "请求参数异常")
 		return
 	}
 	var user *models.User
 	switch {
 	case utils.EmailFormat(req.Username):
 		user = providers.UserService.GetUserByEmail(req.Username)
 	case utils.MobileFormat(req.Username):
 		user = providers.UserService.GetUserByMobile(req.Username)
 	default:
 		user = providers.UserService.GetUserByNickName(req.Username)
 	}
 	if user == nil {
 		utils.APIError(ctx, "用户不存在")
 		return
 	}
 	passwordCheck, err := coding.MD5(req.Password + "|" + user.Salt)
 	if err != nil || passwordCheck != user.Password {
 		utils.APIError(ctx, "密码不正确")
 		return
 	}
 	if user.Status == constants.USER_STATUS_FORBIDDEN {
 		utils.APIError(ctx, "用户已被禁用")
 		return
 	}
 	token, exp, err := utils.IssueToken(user.ID, constants.WEB_JWT_SECRET, constants.WEB_JWT_TTL)
 	if err != nil {
 		utils.APIError(ctx, "签发 token 失败")
 		return
 	}
 	// 同时下发 DES cookie，保持与旧 HTML 路由兼容
 	cookie, _ := coding.DesEncrypt(toString(user.ID), constants.WEB_COOKIE_SALT)
 	utils.SetTokenCookie(ctx, cookie)
 	utils.APIData(ctx, gin.H{
 		"token":      token,
 		"expires_at": exp,
 		"user":       userInfoMap(user),
 	})
 }

 // Info 返回当前登录用户信息。
 func (c *APIAuthController) Info(ctx *gin.Context) {
 	user := utils.GetUser(ctx)
 	if user == nil {
 		utils.APIBadRequest(ctx, "用户ID错误")
 		return
 	}
 	utils.APIData(ctx, userInfoMap(user))
 }

 // Logout 清理 cookie（前端删除本地 token 即视为失效）。
 func (c *APIAuthController) Logout(ctx *gin.Context) {
 	utils.CleanTokenCookie(ctx)
 	utils.APIOK(ctx)
 }

 type apiChangePasswordReq struct {
 	Password string `form:"password" json:"password" binding:"required"`
 }

 // ChangePassword 修改当前用户密码。
 func (c *APIAuthController) ChangePassword(ctx *gin.Context) {
 	user := utils.GetUser(ctx)
 	if user == nil {
 		utils.APIBadRequest(ctx, "用户ID错误")
 		return
 	}
 	var req apiChangePasswordReq
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
 		utils.APIError(ctx, "修改密码失败")
 		return
 	}
 	utils.CleanTokenCookie(ctx)
 	utils.APIOK(ctx)
 }

 func userInfoMap(user *models.User) gin.H {
 	return gin.H{
 		"id":       user.ID,
 		"nickname": user.NickName,
 		"avatar":   user.Avatar,
 		"email":    user.Email,
 		"mobile":   user.Mobile,
 		"role":     user.Role,
 		"status":   user.Status,
 	}
 }

func toString(i int64) string {
 	return itoa(i)
 }

 // itoa 轻量整数转字符串，仅支持非负。
 func itoa(i int64) string {
 	if i == 0 {
 		return "0"
 	}
 	var buf [20]byte
 	pos := len(buf)
 	for i > 0 {
 		pos--
 		buf[pos] = byte('0' + i%10)
 		i /= 10
 	}
 	return string(buf[pos:])
}
