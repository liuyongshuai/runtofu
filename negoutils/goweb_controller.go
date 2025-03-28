// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        goweb_controller.go
// @date        2018-02-02 17:51

package negoutils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type RuntofuControllerInterface interface {
	Init(ct *RuntofuContext, app interface{}, tpl *TplBuilder, tplInitData map[interface{}]interface{})
	Prepare() error //做一些预处理工作，如登录校验、提取用户信息等
	Run()           //正儿八经的业务逻辑处理
	Finish()        //结束时的清理工作，一般不用实现
}

// 控制层基类
type RuntofuController struct {
	Ctx           *RuntofuContext
	AppController interface{}
	Tpl           *TplBuilder                 //模板对象类型
	TplData       map[interface{}]interface{} //赋给tpl模板的变量
	TplName       string                      //模板名称，如“index.tpl”
	TplSections   map[string]string           //页面上各个块
	MainContent   string                      //当前模板的主内容
}

/*
*
**********************************************
初始化相关操作，自动执行，一般不用
**********************************************
*/
func (c *RuntofuController) Init(ctx *RuntofuContext, app interface{}, tpl *TplBuilder, tplInitData map[interface{}]interface{}) {
	c.Ctx = ctx
	c.AppController = app
	c.Tpl = tpl
	c.TplData = make(map[interface{}]interface{})

	//往模板上赋一些公共的数据
	c.TplData["SERVER_REQUEST_URL"] = ctx.Input.URL()
	c.TplData["SERVER_REQUEST_URI"] = ctx.Input.URI()
	c.TplData["REQUEST_DOMAIN"] = ctx.Input.Domain()
	c.TplData["REQUEST_SITE"] = ctx.Input.Site()
	for k, v := range tplInitData {
		c.TplData[k] = v
	}
}

/*
*
**********************************************
预先执行的方法，一般用作统一校验是否登录、打开资源等
**********************************************
*/
func (c *RuntofuController) Prepare() error {
	return nil
}

/*
*
**********************************************
具体的业务逻辑
**********************************************
*/
func (c *RuntofuController) Run() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

/*
*
**********************************************
所有最后的清理工作
**********************************************
*/
func (c *RuntofuController) Finish() {}

// 添加输出的响应头信息
func (c *RuntofuController) AddHeader(key, val string) {
	c.Ctx.Output.AddHeader(key, val)
}

// 设置输出的body
func (c *RuntofuController) SetBody(body []byte) {
	c.Ctx.Output.SetBody(body)
}

// 添加模板数据
func (c *RuntofuController) AddTplData(k interface{}, v interface{}) {
	c.TplData[k] = v
}

// 批量添加模板数据
func (c *RuntofuController) AddTplDatas(d map[interface{}]interface{}) {
	for k, v := range d {
		c.TplData[k] = v
	}
}

// 返回json数据
func (c *RuntofuController) RenderJson(data interface{}) error {
	return c.Ctx.Output.RenderJson(data)
}

// 响应jsonp数据，要求传个callback参数
func (c *RuntofuController) RenderJsonp(data interface{}, callback ...string) error {
	return c.Ctx.Output.RenderJsonp(data, callback...)
}

// 渲染html模板
func (c *RuntofuController) RenderHtml() error {
	buf := new(bytes.Buffer)
	err := c.Tpl.ExecuteTpl(buf, c.TplName, c.TplData)
	if err != nil {
		fmt.Println("controller.Tpl.ExecuteTpl", err)
		return err
	}
	c.SetBody(buf.Bytes())
	return err
}

// 设置响应的状态值
func (c *RuntofuController) SetStatus(status int) {
	c.Ctx.Output.SetStatus(status)
}

// 设置cookie值
func (c *RuntofuController) AddCookie(name string, value string, others ...interface{}) {
	c.Ctx.Output.AddCookie(name, value, others...)
}

// 重定向
func (c *RuntofuController) Redirect(url string, code int) {
	c.Ctx.Redirect(url, code)
}

// 提取表单数据
func (c *RuntofuController) FormParam() url.Values {
	if c.Ctx.Request.Form == nil {
		c.Ctx.Request.ParseForm()
	}
	return c.Ctx.Request.Form
}

// 提取参数
func (c *RuntofuController) GetParam(key string, defaultVal ...interface{}) ElemType {
	if v := c.Ctx.Input.Query(key); v != "" {
		return MakeElemType(v)
	}
	if len(defaultVal) > 0 {
		return MakeElemType(defaultVal[0])
	}
	return MakeElemType("")
}

// 获取上传文件
func (c *RuntofuController) GetFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Ctx.Request.FormFile(key)
}

// 获取所有的上传文件
func (c *RuntofuController) GetFiles(key string) ([]*multipart.FileHeader, error) {
	if files, ok := c.Ctx.Request.MultipartForm.File[key]; ok {
		return files, nil
	}
	return nil, http.ErrMissingFile
}

// 将上传的文件保存到本地
func (c *RuntofuController) SaveToFile(fromfile, tofile string) error {
	file, _, err := c.Ctx.Request.FormFile(fromfile)
	if err != nil {
		return err
	}
	defer file.Close()
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}

// 是否为异步请求
func (c *RuntofuController) IsAjax() bool {
	return c.Ctx.Input.IsAjax()
}

// 返回UserAgent
func (c *RuntofuController) GetUserAgent() string {
	return c.Ctx.Input.UserAgent()
}

// 返回Referer信息
func (c *RuntofuController) GetReferer() string {
	return c.Ctx.Input.Referer()
}

// 客户端
func (c *RuntofuController) GetRemoteIP() string {
	return c.Ctx.Input.IP()
}

// 域名信息
func (c *RuntofuController) GetDomain() string {
	return c.Ctx.Input.Domain()
}

// 获取请求的URI信息
func (c *RuntofuController) GetURI() string {
	return c.Ctx.Input.URI()
}

// 获取请求的URL信息
func (c *RuntofuController) GetURL() string {
	return c.Ctx.Input.URL()
}
