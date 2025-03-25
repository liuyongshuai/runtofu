/**
 * @author      Liu Yongshuai
 * @package     admin
 * @date        2018-02-09 11:14
 */
package admin

import (
	"github.com/liuyongshuai/runtofu/model"
	"github.com/liuyongshuai/runtofu/utils"
)

type ArticleListController struct {
	AdminBaseController
}

//校验是否登录
func (bc *ArticleListController) Run() {
	bc.TplName = "articlelist.tpl"
	//接收两个参数
	page, _ := bc.GetParam("page", 1).ToInt()
	is_publish, _ := bc.GetParam("is_publish", -1).ToInt()
	is_rec, _ := bc.GetParam("is_rec", -1).ToInt()
	sctime := bc.GetParam("sctime", "").ToString()
	ectime := bc.GetParam("ectime", "").ToString()
	st := utils.StrToTime(sctime)
	et := utils.StrToTime(ectime)

	cond := make(map[string]interface{})
	if is_publish >= 0 {
		cond["is_publish"] = is_publish
	}
	if is_rec >= 0 {
		cond["is_rec"] = is_rec
	}
	if st > 0 {
		cond["create_time:gte"] = st
	}
	if et > 0 {
		cond["create_time:lte"] = et
	}

	//提取总数
	total := model.MArticle.GetArticleTotal(cond)
	articleList, _ := model.MArticle.GetArticleList(cond, page, int(gPageSize))
	bc.TplData["articleList"] = articleList
	bc.TplData["pagination"] = bc.getPagination(total)
	bc.TplData["is_publish"] = is_publish
	bc.TplData["is_rec"] = is_rec
	bc.TplData["sctime"] = sctime
	bc.TplData["ectime"] = ectime

	//所有的话题标签
	tagList, _ := model.MTag.GetTagList(1, int(model.MTag.GetTagTotal()))
	bc.TplData["allTagList"] = tagList

	bc.RenderHtml()
}
