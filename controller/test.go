package controller

type TestController struct {
	BaseController
}

//校验是否登录
func (bc *TestController) Run() {
	bc.TplName = "test.tpl"
	bc.RenderHtml()
}
