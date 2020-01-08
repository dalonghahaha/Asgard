package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	token, err := ctx.Cookie("token")
	if err != nil && err != http.ErrNoCookie {
		ctx.Redirect(http.StatusFound, "/error")
	} else if token == "" || err == http.ErrNoCookie {
		ctx.Redirect(http.StatusFound, "/nologin")
	} else {
		ctx.Next()
	}
}
