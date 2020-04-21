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

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) List(ctx *gin.Context) {
	userID := GetUserID(ctx)
	if userID == 0 {
		utils.APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := providers.UserService.GetUserByID(userID)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	page := utils.DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	list, total := providers.UserService.GetUserPageList(where, page, PageSize)
	mpurl := "/user/list"
	ctx.HTML(200, "user/list", gin.H{
		"Subtitle":   "用户列表",
		"List":       list,
		"Total":      total,
		"Role":       user.Role,
		"Pagination": utils.PagerHtml(total, page, mpurl),
	})
}

func (c *UserController) Add(ctx *gin.Context) {
	ctx.HTML(StatusOK, "user/add", gin.H{
		"Subtitle": "添加用户",
	})
}

func (c *UserController) Create(ctx *gin.Context) {
	nickname := ctx.PostForm("nickname")
	email := ctx.PostForm("email")
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	if !utils.Required(ctx, nickname, "昵称不能为空") {
		return
	}
	if !utils.Required(ctx, email, "邮箱不能为空") {
		return
	}
	if !utils.EmailFormat(email) {
		utils.APIBadRequest(ctx, "邮箱格式不正确")
		return
	}
	if !utils.Required(ctx, mobile, "手机号不能为空") {
		return
	}
	if !utils.MobileFormat(mobile) {
		utils.APIBadRequest(ctx, "邮箱格式不正确")
		return
	}
	if !utils.Required(ctx, password, "密码不能为空") {
		return
	}
	usercheck := providers.UserService.GetUserByNickName(nickname)
	if usercheck != nil {
		utils.APIBadRequest(ctx, "该昵称已经注册")
		return
	}
	usercheck = providers.UserService.GetUserByEmail(email)
	if usercheck != nil {
		utils.APIBadRequest(ctx, "该邮箱已经注册")
		return
	}
	usercheck = providers.UserService.GetUserByMobile(mobile)
	if usercheck != nil {
		utils.APIBadRequest(ctx, "该手机号已经注册")
		return
	}
	salt := random.Letters(8)
	password, err := coding.MD5(password + "|" + salt)
	if err != nil {
		utils.APIError(ctx, "生产密码失败")
		return
	}
	user := new(models.User)
	user.NickName = nickname
	user.Email = email
	user.Mobile = mobile
	user.Salt = salt
	user.Password = password
	user.Role = constants.USER_ROLE_NORMAL
	user.Status = constants.USER_STATUS_NORMAL
	ok := providers.UserService.CreateUser(user)
	if !ok {
		utils.APIError(ctx, "创建用户")
		return
	}
	utils.APIOK(ctx)
}

func (c *UserController) Info(ctx *gin.Context) {
	userID := GetUserID(ctx)
	if userID == 0 {
		utils.APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := providers.UserService.GetUserByID(userID)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	utils.APIData(ctx, gin.H{
		"id":       user.ID,
		"nickname": user.NickName,
		"avatar":   user.Avatar,
		"role":     user.Role,
	})
}

func (c *UserController) Edit(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	user := providers.UserService.GetUserByID(id)
	if user == nil {
		utils.JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "user/edit", gin.H{
		"Subtitle": "用户信息修改",
		"User":     user,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	nickname := ctx.PostForm("nickname")
	email := ctx.PostForm("email")
	mobile := ctx.PostForm("mobile")
	if id == 0 {
		utils.APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := providers.UserService.GetUserByID(id)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	user.NickName = nickname
	user.Email = email
	user.Mobile = mobile
	//处理头像
	avatarFile, err := ctx.FormFile("avatar")
	if err == nil {
		fileName, err := coding.MD5(avatarFile.Filename)
		if err != nil {
			utils.APIBadRequest(ctx, "生成文件名失败")
			return
		}
		avatarPath := "web/assets/upload/" + fileName + ".jpg"
		err = ctx.SaveUploadedFile(avatarFile, avatarPath)
		if err != nil {
			utils.APIBadRequest(ctx, "保存文件失败")
			return
		}
		user.Avatar = "/assets/upload/" + fileName + ".jpg"
	}
	ok := providers.UserService.UpdateUser(user)
	if !ok {
		utils.APIError(ctx, "保存设置失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *UserController) Setting(ctx *gin.Context) {
	userID := GetUserID(ctx)
	if userID == 0 {
		utils.JumpError(ctx)
		return
	}
	user := providers.UserService.GetUserByID(userID)
	if user == nil {
		utils.JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "user/setting", gin.H{
		"Subtitle": "用户信息设置",
		"User":     user,
	})
}

func (c *UserController) DoSetting(ctx *gin.Context) {
	nickname := ctx.PostForm("nickname")
	email := ctx.PostForm("email")
	mobile := ctx.PostForm("mobile")
	userID := GetUserID(ctx)
	if userID == 0 {
		utils.APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := providers.UserService.GetUserByID(userID)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	user.NickName = nickname
	user.Email = email
	user.Mobile = mobile
	//处理头像
	avatarFile, err := ctx.FormFile("avatar")
	if err == nil {
		fileName, err := coding.MD5(avatarFile.Filename)
		if err != nil {
			utils.APIBadRequest(ctx, "生成文件名失败")
			return
		}
		avatarPath := "web/assets/upload/" + fileName + ".jpg"
		err = ctx.SaveUploadedFile(avatarFile, avatarPath)
		if err != nil {
			utils.APIBadRequest(ctx, "保存文件失败")
			return
		}
		user.Avatar = "/assets/upload/" + fileName + ".jpg"
	}
	ok := providers.UserService.UpdateUser(user)
	if !ok {
		utils.APIError(ctx, "保存设置失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *UserController) Verify(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.APIBadRequest(ctx, "ID格式错误")
		return
	}
	user := providers.UserService.GetUserByID(id)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	user.Status = constants.USER_STATUS_NORMAL
	ok := providers.UserService.UpdateUser(user)
	if !ok {
		utils.APIError(ctx, "禁用用户失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *UserController) Forbidden(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.APIBadRequest(ctx, "ID格式错误")
		return
	}
	user := providers.UserService.GetUserByID(id)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	user.Status = constants.USER_STATUS_FORBIDDEN
	ok := providers.UserService.UpdateUser(user)
	if !ok {
		utils.APIError(ctx, "禁用用户失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	if id == 0 {
		utils.JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "user/reset_password", gin.H{
		"Subtitle": "重置密码",
		"ID":       id,
	})
}

func (c *UserController) DoResetPassword(ctx *gin.Context) {
	id := utils.DefaultInt64(ctx, "id", 0)
	password := ctx.PostForm("password")
	if id == 0 {
		utils.APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := providers.UserService.GetUserByID(id)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	salt := random.Letters(8)
	password, err := coding.MD5(password + "|" + salt)
	if err != nil {
		utils.APIError(ctx, "生成密码失败")
		return
	}
	user.Salt = salt
	user.Password = password
	ok := providers.UserService.UpdateUser(user)
	if !ok {
		utils.APIError(ctx, "重置密码失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := GetUserID(ctx)
	if userID == 0 {
		utils.JumpError(ctx)
		return
	}
	user := providers.UserService.GetUserByID(userID)
	if user == nil {
		utils.JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "user/change_password", gin.H{
		"Subtitle": "修改密码",
	})
}

func (c *UserController) DoChangePassword(ctx *gin.Context) {
	password := ctx.PostForm("password")
	userID := GetUserID(ctx)
	if userID == 0 {
		utils.APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := providers.UserService.GetUserByID(userID)
	if user == nil {
		utils.APIBadRequest(ctx, "用户不存在")
		return
	}
	salt := random.Letters(8)
	password, err := coding.MD5(password + "|" + salt)
	if err != nil {
		utils.APIError(ctx, "生产密码失败")
		return
	}
	user.Salt = salt
	user.Password = password
	ok := providers.UserService.UpdateUser(user)
	if !ok {
		utils.APIError(ctx, "修改密码失败")
		return
	}
	ctx.SetCookie("token", "", 0, "/", Domain, false, true)
	utils.APIOK(ctx)
}

func (c *UserController) Register(ctx *gin.Context) {
	ctx.HTML(StatusOK, "user/register.html", gin.H{
		"Subtitle": "用户注册",
	})
}

func (c *UserController) DoRegister(ctx *gin.Context) {
	nickname := ctx.PostForm("nickname")
	email := ctx.PostForm("email")
	mobile := ctx.PostForm("mobile")
	password := ctx.PostForm("password")
	passwordConfirm := ctx.PostForm("password-confirm")
	if !utils.Required(ctx, nickname, "昵称不能为空") {
		return
	}
	if !utils.Required(ctx, email, "邮箱不能为空") {
		return
	}
	if !utils.EmailFormat(email) {
		utils.APIBadRequest(ctx, "邮箱格式不正确")
		return
	}
	if !utils.Required(ctx, mobile, "手机号不能为空") {
		return
	}
	if !utils.MobileFormat(mobile) {
		utils.APIBadRequest(ctx, "邮箱格式不正确")
		return
	}
	if !utils.Required(ctx, password, "密码不能为空") {
		return
	}
	if !utils.Required(ctx, passwordConfirm, "确认密码不能为空") {
		return
	}
	if password != passwordConfirm {
		utils.APIBadRequest(ctx, "两次输入的密码不一致")
		return
	}
	usercheck := providers.UserService.GetUserByNickName(nickname)
	if usercheck != nil {
		utils.APIBadRequest(ctx, "该昵称已经注册")
		return
	}
	usercheck = providers.UserService.GetUserByEmail(email)
	if usercheck != nil {
		utils.APIBadRequest(ctx, "该邮箱已经注册")
		return
	}
	usercheck = providers.UserService.GetUserByMobile(mobile)
	if usercheck != nil {
		utils.APIBadRequest(ctx, "该手机号已经注册")
		return
	}
	salt := random.Letters(8)
	password, err := coding.MD5(password + "|" + salt)
	if err != nil {
		utils.APIError(ctx, "生产密码失败")
		return
	}
	user := new(models.User)
	user.NickName = nickname
	user.Email = email
	user.Mobile = mobile
	user.Salt = salt
	user.Password = password
	user.Role = constants.USER_ROLE_NORMAL
	user.Status = constants.USER_STATUS_UNVERIFIED
	ok := providers.UserService.CreateUser(user)
	if !ok {
		utils.APIError(ctx, "注册失败")
		return
	}
	utils.APIOK(ctx)
}

func (c *UserController) Login(ctx *gin.Context) {
	ctx.HTML(StatusOK, "user/login.html", gin.H{
		"Subtitle": "用户登录",
	})
}

func (c *UserController) DoLogin(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if !utils.Required(ctx, username, "用户名不能为空") {
		return
	}
	if !utils.Required(ctx, password, "密码不能为空") {
		return
	}
	var user *models.User
	if utils.EmailFormat(username) {
		user = providers.UserService.GetUserByEmail(username)
	} else if utils.MobileFormat(username) {
		user = providers.UserService.GetUserByMobile(username)
	} else {
		user = providers.UserService.GetUserByNickName(username)
	}
	if user == nil {
		utils.APIError(ctx, "用户不存在")
		return
	}
	passwordCheck, err := coding.MD5(password + "|" + user.Salt)
	if err != nil {
		utils.APIError(ctx, "密码不正确")
		return
	}
	if passwordCheck != user.Password {
		utils.APIError(ctx, "密码不正确")
		return
	}
	cookie, err := coding.DesEncrypt(strconv.Itoa(int(user.ID)), CookieSalt)
	if err != nil {
		utils.APIError(ctx, "登录失败")
	}
	//add cookie
	ctx.SetCookie("token", cookie, 3600, "/", Domain, false, true)
	utils.APIOK(ctx)
}

func (c *UserController) Logout(ctx *gin.Context) {
	//remove cookie
	ctx.SetCookie("token", "", 0, "/", Domain, false, true)
	ctx.Redirect(StatusFound, "/login")
}
