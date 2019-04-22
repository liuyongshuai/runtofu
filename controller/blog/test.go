package blog

type TestController struct {
	RunToFuBaseController
}

//校验是否登录
func (bc *TestController) Run() {
	bc.TplName = "test.tpl"
	bc.RenderHtml()
}
