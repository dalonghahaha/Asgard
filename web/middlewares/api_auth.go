 package middlewares

 import (
 	"net/http"
 	"strconv"
 	"strings"

 	"github.com/dalonghahaha/avenger/tools/coding"
 	"github.com/gin-gonic/gin"

 	"Asgard/constants"
 	"Asgard/providers"
 	"Asgard/web/utils"
 )

 // APIAuth 同时接受 Authorization: Bearer <jwt> 与现有 DES cookie（双轨过渡期）。
 // 解析成功后将 user 挂到 gin.Context，复用 GetUser/GetUserID。
 func APIAuth(ctx *gin.Context) {
 	userID, ok := extractUserID(ctx)
 	if !ok {
 		abortUnauthorized(ctx, "未登录或登录已失效")
 		return
 	}
 	user := providers.UserService.GetUserByID(userID)
 	if user == nil {
 		abortUnauthorized(ctx, "用户不存在")
 		return
 	}
 	if user.Status == constants.USER_STATUS_FORBIDDEN {
 		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
 			"code":    http.StatusForbidden,
 			"message": "用户已被禁用",
 		})
 		return
 	}
 	ctx.Set("user", user)
 	ctx.Next()
 }

 // APIAuthAdmin 校验当前 user 是否为管理员，复用 APIAuth 注入的 user。
 func APIAuthAdmin(ctx *gin.Context) {
 	user := utils.GetUser(ctx)
 	if user == nil {
 		abortUnauthorized(ctx, "未登录或登录已失效")
 		return
 	}
 	if user.Role != constants.USER_ROLE_ADMIN {
 		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
 			"code":    http.StatusForbidden,
 			"message": "只有管理员才能进行此操作",
 		})
 		return
 	}
 	ctx.Next()
 }

 func extractUserID(ctx *gin.Context) (int64, bool) {
 	auth := ctx.GetHeader("Authorization")
 	if auth != "" {
 		// 形如 "Bearer xxx" 或直接 "xxx"
 		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
 		if token == "" {
 			token = strings.TrimSpace(auth)
 		}
 		if token != "" {
 			uid, err := utils.ParseToken(token, constants.WEB_JWT_SECRET)
 			if err == nil {
 				return uid, true
 			}
 		}
 	}
 	// 回落到 DES cookie，便于过渡期旧 HTML 路由与 API 共存
 	if cookie, err := ctx.Cookie("token"); err == nil && cookie != "" {
 		plain, derr := coding.DesDecrypt(cookie, constants.WEB_COOKIE_SALT)
 		if derr == nil {
 			if uid, perr := strconv.ParseInt(plain, 10, 64); perr == nil {
 				return uid, true
 			}
 		}
 	}
 	return 0, false
 }

 func abortUnauthorized(ctx *gin.Context, message string) {
 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
 		"code":    http.StatusUnauthorized,
 		"message": message,
 	})
 }
