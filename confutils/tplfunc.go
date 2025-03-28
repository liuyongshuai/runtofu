// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        tplfunc.go
// @date        2025-03-25 15:48

package confutils

import (
	"github.com/liuyongshuai/runtofu/negoutils"
	"html/template"
)

var TplFuncs template.FuncMap
var conf *RuntofuConfig

func init() {
	conf = GetConfiger()
	TplFuncs = make(template.FuncMap)
	TplFuncs["static_css"] = StaticCSS
	TplFuncs["static_js"] = StaticJS
	TplFuncs["static_image"] = StaticImage
	TplFuncs["editor_js"] = StaticEditorMdJS
	TplFuncs["editor_css"] = StaticEditorMdCSS
	TplFuncs["ftime"] = negoutils.FormatCTime
	TplFuncs["strtotime"] = negoutils.StrToTimestamp
	for k, v := range negoutils.CommonTplFuncs {
		TplFuncs[k] = v
	}
}

// 插入js引入文件标签，js形如“base.js”
func StaticJS(js string) template.HTML {
	js = conf.Common.StaticPrefix + "/js/" + js
	js = "<script charset=\"utf-8\" type=\"text/javascript\" src=\"" + js + "\"></script>"
	return template.HTML(js)
}

// 插入css引入文件标签,css形如“base.css”
func StaticCSS(css string) template.HTML {
	css = conf.Common.StaticPrefix + "/css/" + css
	css = "<link type=\"text/css\" rel=\"stylesheet\" href=\"" + css + "\" />"
	return template.HTML(css)
}

// 插入图片路径
func StaticImage(image string) template.HTML {
	image = conf.Common.StaticPrefix + "/images/" + image
	return template.HTML(image)
}

// 插入editor.md编辑器
func StaticEditorMdJS(p string) template.HTML {
	p = conf.Common.StaticPrefix + "/editor.md/" + p
	js := "<script charset=\"utf-8\" type=\"text/javascript\" src=\"" + p + "\"></script>"
	return template.HTML(js)
}

// 插入editor.md编辑器
func StaticEditorMdCSS(p string) template.HTML {
	p = conf.Common.StaticPrefix + "/editor.md/" + p
	css := "<link type=\"text/css\" rel=\"stylesheet\" href=\"" + p + "\" />"
	return template.HTML(css)
}
