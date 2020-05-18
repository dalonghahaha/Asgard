package constants

import "net/http"

var (
	WEB_PORT        = 12345
	WEB_MODE        = "release"
	WEB_DOMAIN      = "localhost"
	WEB_COOKIE_SALT = "sdswqeqx"
	StatusOK        = http.StatusOK
	StatusFound     = http.StatusFound
	StatusForbidden = http.StatusForbidden
	Lang            = "cn"
	TimeLocation    = "Asia/Shanghai"
	TimeLayout      = "2006-01-02 15:04"
)
