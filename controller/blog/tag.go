/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     controller
 * @date        2018-02-05 12:51
 */
package blog

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/model"
)

type TagController struct {
	RunToFuBaseController
}

//校验是否登录
func (bc *TagController) Run() {
	bc.TplName = "tag.tpl"
	tagId, _ := bc.GetParam("tag_id", 0).ToInt()
	if tagId <= 0 {
		fmt.Println(tagId)
	}

	//接收页码参数
	page, _ := bc.GetParam("page", 1).ToInt()

	//提取文章ID列表
	aids, _ := model.MArticleTag.GetArticleList(tagId, page, int(gPageSize))

	//提取对应的文章内容信息
	alist, _ := model.MArticle.GetArticleInfos(aids)
	var articleList []model.ArticleInfo
	for _, aid := range aids {
		if ainfo, ok := alist[aid]; ok {
			articleList = append(articleList, ainfo)
		}
	}

	//获取总数，分页用的
	total := model.MTag.GetTagTotal()
	htmlStr, err := bc.articleListAndTags(articleList, total)
	bc.TplData["articleTagContent"] = htmlStr
	bc.TplData["err"] = fmt.Sprintf("%v", err)
	bc.RenderHtml()
}
