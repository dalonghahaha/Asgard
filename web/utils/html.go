package utils

import (
	"html/template"
	"strconv"
	"strings"
)

var PageSize = 10

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
