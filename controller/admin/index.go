/**
 * @author      Liu Yongshuai
 * @package     admin
 * @date        2018-02-11 21:47
 */
package admin

type IndexController struct {
	AdminBaseController
}

//校验是否登录
func (bc *IndexController) Run() {
	bc.TplName = "index.tpl"
	bc.RenderHtml()
}
