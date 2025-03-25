/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-10 16:03
 */
package ajax

import (
	"github.com/liuyongshuai/runtofu/model"
)

type AdminAjaxArticleModifyController struct {
	AdminAjaxBaseController
}

//返回数据信息
func (bc *AdminAjaxArticleModifyController) Run() {
	action := bc.GetParam("action", "modify").ToString()
	aid, _ := bc.GetParam("article_id", 0).ToInt64()
	title := bc.GetParam("title", "").ToString()
	tags := bc.GetParam("tags", "").ToString()
	content := bc.GetParam("content", "").ToString()
	isPublish, _ := bc.GetParam("is_publish", -1).ToInt64()
	isOrigin, _ := bc.GetParam("is_origin", -1).ToInt64()
	isRec, _ := bc.GetParam("is_rec", -1).ToInt64()
	if aid <= 0 {
		bc.Notice(nil, 100100, "文章ID有误")
		return
	}
	if action == "modify" {
		data := make(map[string]interface{})
		if len(title) > 0 {
			data["title"] = title
		}
		if len(tags) > 0 {
			data["tags"] = tags
		}
		if len(content) > 0 {
			data["content"] = content
		}
		if isPublish >= 0 {
			data["is_publish"] = isPublish
		}
		if isRec >= 0 {
			data["is_rec"] = isRec
		}
		if isOrigin >= 0 {
			data["is_origin"] = isOrigin
		}
		_, e := model.MArticle.UpdateArticleInfo(aid, data)
		if e != nil {
			bc.Notice(nil, 100100, "修改文章失败,"+e.Error())
			return
		}
	} else if action == "delete" {
		_, e := model.MArticle.DeleteArticleInfo(aid)
		if e != nil {
			bc.Notice(nil, 100100, "删除文章失败,"+e.Error())
			return
		}
	}

	bc.Notice(nil)
}
