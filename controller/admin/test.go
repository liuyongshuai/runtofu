package admin

type TestController struct {
	AdminBaseController
}

//校验是否登录
func (bc *TestController) Run() {
	bc.TplName = "test.tpl"
	bc.RenderHtml()
}
