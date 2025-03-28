// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        utils_pager.go
// @date        2025-03-25 10:25

package negoutils

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
	"time"
)

// 分页操作相关，从请求的URL里提取当前页码
// curUriPath为当前请求的路径
// curRawQuery为当前的请求参数，如"a=2&b=3"
// totalNum为总数量
// pagesize为页大小
// pageFieldName为分页时的页码参数名，默认为"page"
func Pagination(curPath, curRawQuery string, tnum, psize interface{}, pageFieldName string) template.HTML {
	totalNum, _ := MakeElemType(tnum).ToInt64()
	pagesize, _ := MakeElemType(psize).ToInt64()
	if len(pageFieldName) <= 0 {
		pageFieldName = "page"
	}
	pager := ""

	var curPage int64 = 1

	//原始请求的参数
	var curQuery []string
	if len(curRawQuery) > 0 {
		tmp := strings.Split(curRawQuery, "&")
		for _, ttt := range tmp {
			tmptmp := strings.Split(ttt, "=")
			if len(tmptmp) != 2 {
				continue
			}
			if tmptmp[0] == pageFieldName {
				curPage, _ = strconv.ParseInt(tmptmp[1], 10, 64)
				continue
			}
			curQuery = append(curQuery, ttt)
		}
	}
	if len(curQuery) <= 0 {
		curQuery = append(curQuery, "_t="+strconv.FormatInt(time.Now().Unix(), 10))
	}
	t1 := float64(totalNum)
	t2 := float64(pagesize)
	totalPage := int64(math.Ceil(t1 / t2))
	if totalPage <= 1 {
		return template.HTML(pager)
	}
	prePage := curPage - 1
	if prePage < 0 {
		prePage = 0
	}
	nextPage := curPage + 1
	if curPage == totalPage {
		nextPage = 0
	}

	curURL := curPath + "?" + strings.Join(curQuery, "&")

	//计算开始页，当前页的前面，最多跟着2个页按钮
	s := curPage - 2
	if s < 1 {
		s = 1
	}

	//计算结束页，当前页的后面，最多跟着4个页按钮
	end := s + 4
	if curPage <= 3 {
		end++
	}
	if end > 0 && end < 5 {
		end = 5
	}

	if end > totalPage {
		end = totalPage
	}

	//结束页跟总页数相等时
	if end == totalPage && s > 1 {
		tmp := s - (3 - (totalPage - curPage))
		if tmp < 1 {
			s = 1
		} else {
			s = tmp
		}
	}

	pager = "<nav class=\"pagination-center\"><ul class=\"pagination\">"

	//上一页
	if prePage > 0 {
		pager = fmt.Sprintf("%s<li><a href=\"%s&%s=%d\">上一页</a></li>", pager, curURL, pageFieldName, prePage)
	} else {
		pager = fmt.Sprintf("%s<li class=\"disabled\"><span>上一页</span></li>", pager)
	}

	//第一页
	if s > 1 {
		pager = fmt.Sprintf("%s<li><a href=\"%s&%s=%d\">1</a></li>", pager, curURL, pageFieldName, 1)
		if s > 2 {
			pager = fmt.Sprintf("%s<li class=\"disabled\"><span>...</span></li>", pager)
		}
	}

	//循环显示中间部分
	for i := s; i <= end; i++ {
		if i == curPage {
			pager = fmt.Sprintf("%s<li class=\"active\"><span>%d</span></li>", pager, curPage)
		} else {
			pager = fmt.Sprintf("%s<li><a href=\"%s&%s=%d\">%d</a></li>", pager, curURL, pageFieldName, i, i)
		}
	}

	//最后一页
	if end < totalPage {
		if end+1 < totalPage {
			pager = fmt.Sprintf("%s<li class=\"disabled\"><span>...</span></li>", pager)
		}
		pager = fmt.Sprintf("%s<li><a href=\"%s&%s=%d\">%d</a></li>", pager, curURL, pageFieldName, totalPage, totalPage)
	}

	//下一页
	if nextPage > 0 {
		pager = fmt.Sprintf("%s<li><a href=\"%s&%s=%d\">下一页</a></li>", pager, curURL, pageFieldName, nextPage)
	} else {
		pager = fmt.Sprintf("%s<li class=\"disabled\"><span>下一页</span></li>", pager)
	}
	pager += "</ul></nav>"

	return template.HTML(pager)
}
