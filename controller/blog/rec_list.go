/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     controller
 * @date        2018-02-03 15:40
 */
package blog

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/model"
)

type RecListController struct {
	RunToFuBaseController
}

//校验是否登录
func (bc *RecListController) Run() {
	bc.TplName = "rec_list.tpl"
	//接收页码参数
	page, _ := bc.GetParam("page", 1).ToInt()

	//提取当前页的文章ID列表
	cond := make(map[string]interface{})
	cond["is_rec"] = 1
	articleList, err := model.MArticle.GetArticleList(cond, page, int(gPageSize))

	//提了总数，分页用的
	total := model.MArticle.GetArticleTotal(cond)

	//生成html文本
	htmlStr, err := bc.articleListAndTags(articleList, total)
	bc.TplData["recListContent"] = htmlStr
	bc.TplData["err"] = fmt.Sprintf("%v", err)
	bc.RenderHtml()
}
