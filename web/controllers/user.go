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
	list := c.useService.GetAllUser()
	total := len(list)
	mpurl := "/user/list"
	ctx.HTML(200, "user/list", gin.H{
		"Subtitle":   "用户列表",
		"List":       list,
		"Total":      total,
		"Pagination": PagerHtml(total, page, mpurl),
	})
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
	user.Salt = salt
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
