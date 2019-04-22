/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     controller
 * @date        2018-02-03 12:12
 */
package blog

import (
	"github.com/liuyongshuai/runtofu/model"
)

type IndexController struct {
	RunToFuBaseController
}

//提取首页的文章列表
func (bc *IndexController) Run() {
	bc.TplName = "index.tpl"
	//接收两个参数
	page, _ := bc.GetParam("page", 1).ToInt()

	//提取总数
	cond := make(map[string]interface{})
	cond["is_publish"] = 1
	total := model.MArticle.GetArticleTotal(cond)
	articleList, err := model.MArticle.GetArticleList(cond, page, int(gPageSize))
	if err != nil {
		bc.TplData["err"] = err.Error()
		bc.RenderHtml()
		return
	}

	//生成html文本
	htmlStr, err := bc.articleListAndTags(articleList, total)
	if err != nil {
		bc.TplData["err"] = err.Error()
		bc.RenderHtml()
		return
	}

	//赋文章列表和话题列表的值
	bc.TplData["indexContent"] = htmlStr
	bc.RenderHtml()
}
