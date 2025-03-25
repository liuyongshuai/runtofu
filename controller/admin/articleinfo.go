/**
 * @author      Liu Yongshuai
 * @package     admin
 * @date        2018-02-09 16:24
 */
package admin

import (
	"net/http"
)

type ArticleInfoController struct {
	AdminBaseController
}

//校验是否登录
func (bc *ArticleInfoController) Run() {
	bc.TplName = "articleinfo.tpl"
	aid, _ := bc.GetParam("article_id", 0).ToInt64()
	if aid <= 0 {
		bc.SetStatus(http.StatusNotFound)
		return
	}
	bc.TplData["articleId"] = aid
	bc.RenderHtml()
}
