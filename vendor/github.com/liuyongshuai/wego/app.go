package wego

import (
	"fmt"
	"github.com/liuyongshuai/wego/controller"
	"github.com/liuyongshuai/wego/router"
	"html/template"
	"net/http"
)

//新建一个APP
func NewWeGoAPP() *WeGoApp {
	app := &WeGoApp{
		Handlers: NewWeGoHandler(),
	}
	return app
}

//APP结构体
type WeGoApp struct {
	Handlers *WeGoHandler //处理句柄
}

//开始运行
func (app *WeGoApp) Run() {
	app.Handlers.Tpl.SetRootPathDir(app.Handlers.TplDir).SetTplExt(app.Handlers.TplExt)
	err := http.ListenAndServe(":"+app.Handlers.Port, app.Handlers)
	if err != nil {
		fmt.Println(err)
	}
}

//设置监听端口
func (app *WeGoApp) SetPort(port string) *WeGoApp {
	app.Handlers.Port = port
	return app
}

//设置错误信息提示
func (app *WeGoApp) SetErrController(c controller.WeGoControllerInterface) *WeGoApp {
	c = c.(controller.WeGoControllerInterface)
	app.Handlers.SetErrController(c)
	return app
}

//设置POST最大内存
func (app *WeGoApp) SetMaxMemory(n int64) *WeGoApp {
	app.Handlers.SetMaxMemory(n)
	return app
}

//设置模板路径
func (app *WeGoApp) SetTplDir(dir string) *WeGoApp {
	app.Handlers.SetTplDir(dir)
	return app
}

//设置模板扩展名称
func (app *WeGoApp) SetTplExt(ext string) *WeGoApp {
	app.Handlers.SetTplExt(ext)
	return app
}

//设置给模板的公共参数
func (app *WeGoApp) SetTplCommonData(data map[interface{}]interface{}) *WeGoApp {
	app.Handlers.SetTplCommonData(data)
	return app
}

//设置给模板的公共参数
func (app *WeGoApp) AddTplCommonData(k interface{}, v interface{}) *WeGoApp {
	app.Handlers.AddTplCommonData(k, v)
	return app
}

//添加一个插件
func (app *WeGoApp) AddHooks(when int, hk HooksFunc) *WeGoApp {
	app.Handlers.AddHooks(when, hk)
	return app
}

//添加模板函数
func (app *WeGoApp) AddTplFuncMap(fm template.FuncMap) *WeGoApp {
	app.Handlers.Tpl.AddTplFuncs(fm)
	return app
}

//添加模板函数
func (app *WeGoApp) AddTplFunc(name string, fn interface{}) *WeGoApp {
	app.Handlers.Tpl.AddTplFunc(name, fn)
	return app
}

//添加一个路由
func (app *WeGoApp) AddRouter(r *router.WeGoRouterItem) *WeGoApp {
	app.Handlers.AddRouter(r)
	return app
}

//批量添加路由
func (app *WeGoApp) AddRouters(rs ...*router.WeGoRouterItem) *WeGoApp {
	app.Handlers.AddRouters(rs...)
	return app
}

//设置发生错误时的处理函数
func (app *WeGoApp) SetRecoverFunc(fn RecoverFunc) *WeGoApp {
	app.Handlers.SetRecoverFunc(fn)
	return app
}
