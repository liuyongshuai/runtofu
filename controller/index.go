/**
 * @author      Liu Yongshuai
 * @package     admin
 * @date        2018-02-11 21:47
 */
package controller

type IndexController struct {
	BaseController
}

//校验是否登录
func (bc *IndexController) Run() {
	bc.TplName = "index.tpl"
	bc.RenderHtml()
}
