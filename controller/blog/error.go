/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     controller
 * @date        2018-02-03 11:13
 */
package blog

import (
	"net/http"
)

type ErrorController struct {
	RunToFuBaseController
}

//校验是否登录
func (bc *ErrorController) Run() {
	bc.TplName = "404.tpl"
	bc.SetStatus(http.StatusNotFound)
	bc.RenderHtml()
}
