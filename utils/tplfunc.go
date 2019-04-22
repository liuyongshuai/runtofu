/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     helper
 * @date        2018-02-02 22:26
 */
package utils

import (
	"fmt"
	"github.com/liuyongshuai/goUtils"
	"github.com/liuyongshuai/runtofuAdmin/config"
	"html/template"
	"regexp"
	"strings"
	"time"
)

var TplFuncs template.FuncMap
var conf *config.WeGoAdminConfig

func init() {
	conf = config.GetConfiger()
	TplFuncs = make(template.FuncMap)
	TplFuncs["static_css"] = StaticCSS
	TplFuncs["static_js"] = StaticJS
	TplFuncs["static_image"] = StaticImage
	TplFuncs["editor_js"] = StaticEditorMdJS
	TplFuncs["editor_css"] = StaticEditorMdCSS
	TplFuncs["ftime"] = FormatCTime
	TplFuncs["strtotime"] = StrToTime
}

//插入js引入文件标签，js形如“base.js”
func StaticJS(js string) template.HTML {
	js = conf.Common.StaticPrefix + "/js/" + js
	js = "<script charset=\"utf-8\" type=\"text/javascript\" src=\"" + js + "\"></script>"
	return template.HTML(js)
}

//插入css引入文件标签,css形如“base.css”
func StaticCSS(css string) template.HTML {
	css = conf.Common.StaticPrefix + "/css/" + css
	css = "<link type=\"text/css\" rel=\"stylesheet\" href=\"" + css + "\" />"
	return template.HTML(css)
}

//插入图片路径
func StaticImage(image string) template.HTML {
	image = conf.Common.StaticPrefix + "/images/" + image
	return template.HTML(image)
}

//插入editor.md编辑器
func StaticEditorMdJS(p string) template.HTML {
	p = conf.Common.StaticPrefix + "/editor.md/" + p
	js := "<script charset=\"utf-8\" type=\"text/javascript\" src=\"" + p + "\"></script>"
	return template.HTML(js)
}

//插入editor.md编辑器
func StaticEditorMdCSS(p string) template.HTML {
	p = conf.Common.StaticPrefix + "/editor.md/" + p
	css := "<link type=\"text/css\" rel=\"stylesheet\" href=\"" + p + "\" />"
	return template.HTML(css)
}

//格式化时间
//t为时间戳，秒数
//format为要格式化的格式，如"Y-m-d H:i:s"
func FormatCTime(t interface{}, format string) template.HTML {
	tm, err := goUtils.MakeElemType(t).ToInt64()
	if err != nil {
		return template.HTML("")
	}
	local, err2 := time.LoadLocation("Asia/Chongqing")
	if err2 != nil {
		fmt.Println(err2)
	}
	replacer := strings.NewReplacer(
		"Y", "2006",
		"m", "01",
		"d", "02",
		"H", "15",
		"i", "04",
		"s", "05",
	)
	format = replacer.Replace(format)
	ret := time.Unix(tm, 0).In(local).Format(format)
	return template.HTML(ret)
}

//将字符串日期转为时间戳，秒
//t要转化的时间戳，如"2018-04-09"、"2018-08-04 12:34"
func StrToTime(t string) int64 {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	reg, _ := regexp.Compile(`[\-:\s]+`)
	tmp := reg.Split(t, -1)
	if len(tmp) <= 0 {
		return 0
	}
	if len(tmp) > 6 {
		tmp = tmp[0:6]
	}
	for k, v := range tmp {
		vv, _ := goUtils.MakeElemType(v).ToInt()
		if k == 0 {
			tmp[k] = fmt.Sprintf("%04d", vv)
		} else {
			tmp[k] = fmt.Sprintf("%02d", vv)
		}
	}
	tl := []string{"", "01", "01", "00", "00", "00"}
	if len(tmp) < 6 {
		tmp = append(tmp, tl[len(tmp):]...)
	}
	tstr := fmt.Sprintf("%s-%s-%s %s:%s:%s", tmp[0], tmp[1], tmp[2], tmp[3], tmp[4], tmp[5])
	utm, err := time.ParseInLocation("2006-01-02 15:04:05", tstr, loc)
	if err != nil {
		return 0
	}
	return utm.Unix()
}
