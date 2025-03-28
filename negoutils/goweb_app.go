// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        goweb_app.go
// @date        2018-02-02 17:51

package negoutils

import (
	"fmt"
	"html/template"
	"net/http"
)

// 新建一个APP
func NewRuntofuAPP() *RuntofuApp {
	app := &RuntofuApp{
		Handlers: NewRuntofuHandler(),
	}
	return app
}

// APP结构体
type RuntofuApp struct {
	Handlers *RuntofuHandler //处理句柄
}

// 开始运行
func (app *RuntofuApp) Run() {
	app.Handlers.Tpl.SetRootPathDir(app.Handlers.TplDir).SetTplExt(app.Handlers.TplExt)
	err := http.ListenAndServe(":"+app.Handlers.Port, app.Handlers)
	if err != nil {
		fmt.Println(err)
	}
}

// 设置监听端口
func (app *RuntofuApp) SetPort(port string) *RuntofuApp {
	app.Handlers.Port = port
	return app
}

// 设置错误信息提示
func (app *RuntofuApp) SetErrController(c RuntofuControllerInterface) *RuntofuApp {
	c = c.(RuntofuControllerInterface)
	app.Handlers.SetErrController(c)
	return app
}

// 设置POST最大内存
func (app *RuntofuApp) SetMaxMemory(n int64) *RuntofuApp {
	app.Handlers.SetMaxMemory(n)
	return app
}

// 设置模板路径
func (app *RuntofuApp) SetTplDir(dir string) *RuntofuApp {
	app.Handlers.SetTplDir(dir)
	return app
}

// 设置模板扩展名称
func (app *RuntofuApp) SetTplExt(ext string) *RuntofuApp {
	app.Handlers.SetTplExt(ext)
	return app
}

// 设置给模板的公共参数
func (app *RuntofuApp) SetTplCommonData(data map[interface{}]interface{}) *RuntofuApp {
	app.Handlers.SetTplCommonData(data)
	return app
}

// 设置给模板的公共参数
func (app *RuntofuApp) AddTplCommonData(k interface{}, v interface{}) *RuntofuApp {
	app.Handlers.AddTplCommonData(k, v)
	return app
}

// 添加一个插件
func (app *RuntofuApp) AddHooks(when int, hk HooksFunc) *RuntofuApp {
	app.Handlers.AddHooks(when, hk)
	return app
}

// 添加模板函数
func (app *RuntofuApp) AddTplFuncMap(fm template.FuncMap) *RuntofuApp {
	app.Handlers.Tpl.AddTplFuncs(fm)
	return app
}

// 添加模板函数
func (app *RuntofuApp) AddTplFunc(name string, fn interface{}) *RuntofuApp {
	app.Handlers.Tpl.AddTplFunc(name, fn)
	return app
}

// 添加一个路由
func (app *RuntofuApp) AddRouter(r *RuntofuRouterItem) *RuntofuApp {
	app.Handlers.AddRouter(r)
	return app
}

// 批量添加路由
func (app *RuntofuApp) AddRouters(rs ...*RuntofuRouterItem) *RuntofuApp {
	app.Handlers.AddRouters(rs...)
	return app
}

// 设置发生错误时的处理函数
func (app *RuntofuApp) SetRecoverFunc(fn RecoverFunc) *RuntofuApp {
	app.Handlers.SetRecoverFunc(fn)
	return app
}
