package controllers

import (
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	StatusOK         = http.StatusOK
	StatusBadRequest = http.StatusBadRequest
	StatusError      = http.StatusInternalServerError
	PageSize         = 10
	CookieSalt       = "sdswqeqx"
	Domain           = "localhost"
)

func EmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func MobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func Required(ctx *gin.Context, val *string, message string) bool {
	if *val == "" {
		APIBadRequest(ctx, message)
		return false
	}
	return true
}

func APIOK(ctx *gin.Context) {
	ctx.JSON(StatusOK, gin.H{"code": StatusOK})
}

func APIData(ctx *gin.Context, data interface{}) {
	ctx.JSON(StatusOK, gin.H{"code": StatusOK, "data": data})
}

func APIBadRequest(ctx *gin.Context, message string) {
	ctx.JSON(StatusOK, gin.H{"code": StatusBadRequest, "message": message})
}

func APIError(ctx *gin.Context, message string) {
	ctx.JSON(StatusOK, gin.H{"code": StatusError, "message": message})
}

func DefaultInt(ctx *gin.Context, key string, defaultVal int) int {
	page := ctx.Query(key)
	if page == "" {
		return defaultVal
	}
	_page, err := strconv.Atoi(page)
	if err != nil {
		return defaultVal
	}
	return _page
}

func totalPage(total int) int {
	if (total % PageSize) != 0 {
		return (total / PageSize) + 1
	} else {
		return total / PageSize
	}
}

func Unescaped(x string) interface{} {
	return template.HTML(x)
}

func PagerHtml(total int, page int, mpurl string) string {
	totalpage := totalPage(total)
	if total == 0 {
		return ""
	}
	if totalpage == 1 {
		return ""
	}
	var max_page, begin, end int
	html := ""
	if total > page {
		html = "<ul class=\"pagination ml-auto\">"
		html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "\">首页</a></li>"
		if page-1 > 0 {
			if mpurl == "javascript:;" {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "\">上一页</a></li>"
			} else if strings.Contains(mpurl, "?") {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "&page=" + strconv.Itoa(page-1) + "\">上一页</a></li>"
			} else {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "?page=" + strconv.Itoa(page-1) + "\">上一页</a></li>"
			}
		} else {
			html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "\">上一页</a></li>"
		}
		//最大页数
		if totalpage <= 9 {
			max_page = totalpage
		} else {
			max_page = 9
		}
		rank := 4
		if page >= max_page {
			if (page - rank) > 0 {
				begin = page - rank
			} else {
				begin = 1
			}
		} else {
			begin = 1
		}
		if page >= max_page {
			if (page + rank) <= totalpage {
				end = page + rank
			} else {
				end = totalpage
			}
		} else {
			end = max_page
		}
		for i := begin; i <= end; i++ {
			var link string
			if mpurl == "javascript:;" {
				link = "javascript:;"
			} else if strings.Contains(mpurl, "?") {
				link = mpurl + "&page=" + strconv.Itoa(i)
			} else {
				link = mpurl + "?page=" + strconv.Itoa(i)
			}
			class := ""
			if i == page {
				link = "javascript:;"
				class = "page-item active"
			} else {
				class = "page-item"
			}
			html += "<li class=\"" + class + "\"><a class=\"page-link\" href=\"" + link + "\">" + strconv.Itoa(i) + "</a></li>"
		}
		if (page + 1) < totalpage {
			if mpurl == "javascript:;" {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "\">下一页</a></li>"
			} else if strings.Contains(mpurl, "?") {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "&page=" + strconv.Itoa(page+1) + "\">下一页</a></li>"
			} else {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "?page=" + strconv.Itoa(page+1) + "\">下一页</a></li>"
			}
		} else {
			if mpurl == "javascript:;" {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "\">下一页</a></li>"
			} else if strings.Contains(mpurl, "?") {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "&page=" + strconv.Itoa(totalpage) + "\">下一页</a></li>"
			} else {
				html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "?page=" + strconv.Itoa(totalpage) + "\">下一页</a></li>"
			}
		}
		if mpurl == "javascript:;" {
			html += "<li><a href=\"" + mpurl + "\">尾页</a></li>"
		} else if strings.Contains(mpurl, "?") {
			html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "&page=" + strconv.Itoa(totalpage) + "\">尾页</a></li>"
		} else {
			html += "<li class=\"page-item\"><a class=\"page-link\" href=\"" + mpurl + "?page=" + strconv.Itoa(totalpage) + "\">尾页</a></li>"
		}
		html += "<ul>"
	}
	return html
}
