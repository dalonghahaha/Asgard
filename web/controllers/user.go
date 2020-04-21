package controllers

import (
	"strconv"

	"github.com/dalonghahaha/avenger/tools/coding"
	"github.com/dalonghahaha/avenger/tools/random"
	"github.com/gin-gonic/gin"

	"Asgard/models"
	"Asgard/services"
)

type UserController struct {
	useService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		useService: services.NewUserService(),
	}
}

func (c *UserController) List(ctx *gin.Context) {
	page := DefaultInt(ctx, "page", 1)
	where := map[string]interface{}{}
	list, total := c.useService.GetUserPageList(where, page, PageSize)
	mpurl := "/user/list"
	ctx.HTML(200, "user/list", gin.H{
		"Subtitle":   "用户列表",
		"List":       list,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
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
	salt := random.Letters(8)
	password, err := coding.MD5(password + "|" + salt)
	if err != nil {
		APIError(ctx, "生产密码失败")
		return
	}
	user := new(models.User)
	user.NickName = nickname
	user.Email = email
	user.Mobile = mobile
	user.Salt = salt
	user.Password = password
	ok := c.useService.CreateUser(user)
	if !ok {
		APIError(ctx, "创建用户")
		return
	}
	APIOK(ctx)
}

func (c *UserController) Info(ctx *gin.Context) {
	userID := GetUserID(ctx)
	if userID == 0 {
		APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := c.useService.GetUserByID(userID)
	if user == nil {
		APIBadRequest(ctx, "用户不存在")
		return
	}
	APIData(ctx, gin.H{
		"id":       user.ID,
		"nickname": user.NickName,
		"avatar":   user.Avatar,
		"role":     "Administrator",
	})
}

func (c *UserController) Setting(ctx *gin.Context) {
	userID := GetUserID(ctx)
	if userID == 0 {
		JumpError(ctx)
		return
	}
	user := c.useService.GetUserByID(userID)
	if user == nil {
		JumpError(ctx)
		return
	}
	ctx.HTML(StatusOK, "user/setting", gin.H{
		"Subtitle": "用户信息设置",
		"User":     user,
	})
}

func (c *UserController) Update(ctx *gin.Context) {
	nickname := ctx.PostForm("nickname")
	email := ctx.PostForm("email")
	mobile := ctx.PostForm("mobile")
	userID := GetUserID(ctx)
	if userID == 0 {
		APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := c.useService.GetUserByID(userID)
	if user == nil {
		APIBadRequest(ctx, "用户不存在")
		return
	}
	user.NickName = nickname
	user.Email = email
	user.Mobile = mobile
	APIOK(ctx)
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := GetUserID(ctx)
	if userID == 0 {
		JumpError(ctx)
		return
	}
	user := c.useService.GetUserByID(userID)
	if user == nil {
		JumpError(ctx)
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
		APIBadRequest(ctx, "用户ID错误")
		return
	}
	user := c.useService.GetUserByID(userID)
	if user == nil {
		APIBadRequest(ctx, "用户不存在")
		return
	}
	salt := random.Letters(8)
	password, err := coding.MD5(password + "|" + salt)
	if err != nil {
		APIError(ctx, "生产密码失败")
		return
	}
	user.Salt = salt
	user.Password = password
	ok := c.useService.UpdateUser(user)
	if !ok {
		APIError(ctx, "修改密码失败")
		return
	}
	ctx.SetCookie("token", "", 0, "/", Domain, false, true)
	APIOK(ctx)
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
	if !Required(ctx, &email, "邮箱不能为空") {
		return
	}
	if !Required(ctx, &mobile, "手机号不能为空") {
		return
	}
	if !Required(ctx, &password, "密码不能为空") {
		return
	}
	if !Required(ctx, &passwordConfirm, "确认密码不能为空") {
		return
	}
	if password != passwordConfirm {
		APIBadRequest(ctx, "两次输入的密码不一致")
		return
	}
	salt := random.Letters(8)
	password, err := coding.MD5(password + "|" + salt)
	if err != nil {
		APIError(ctx, "生产密码失败")
		return
	}
	user := new(models.User)
	user.NickName = nickname
	user.Email = email
	user.Mobile = mobile
	user.Salt = salt
	user.Password = password
	ok := c.useService.CreateUser(user)
	if !ok {
		APIError(ctx, "注册失败")
		return
	}
	APIOK(ctx)
}

func (c *UserController) Login(ctx *gin.Context) {
	ctx.HTML(StatusOK, "user/login.html", gin.H{
		"Subtitle": "用户登录",
	})
}

func (c *UserController) DoLogin(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if !Required(ctx, &username, "用户名不能为空") {
		return
	}
	if !Required(ctx, &password, "密码不能为空") {
		return
	}
	var user *models.User
	if EmailFormat(username) {
		user = c.useService.GetUserByEmail(username)
	} else if MobileFormat(username) {
		user = c.useService.GetUserByMobile(username)
	} else {
		user = c.useService.GetUserByNickName(username)
	}
	if user == nil {
		APIError(ctx, "用户不存在")
		return
	}
	passwordCheck, err := coding.MD5(password + "|" + user.Salt)
	if err != nil {
		APIError(ctx, "密码不正确")
		return
	}
	if passwordCheck != user.Password {
		APIError(ctx, "密码不正确")
		return
	}
	cookie, err := coding.DesEncrypt(strconv.Itoa(int(user.ID)), CookieSalt)
	if err != nil {
		APIError(ctx, "登录失败")
	}
	//add cookie
	ctx.SetCookie("token", cookie, 3600, "/", Domain, false, true)
	APIOK(ctx)
}

func (c *UserController) Logout(ctx *gin.Context) {
	//remove cookie
	ctx.SetCookie("token", "", 0, "/", Domain, false, true)
	ctx.Redirect(StatusFound, "/login")
}
