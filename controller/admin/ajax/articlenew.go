/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-13 21:13
 */
package ajax

import (
	"github.com/liuyongshuai/negoutils"
	"github.com/liuyongshuai/runtofu/model"
	"strings"
)

type AdminAjaxArticleNewController struct {
	AdminAjaxBaseController
}

// 返回数据信息
func (bc *AdminAjaxArticleNewController) Run() {
	title := bc.GetParam("title", "").ToString()
	tags := bc.GetParam("tags", "").ToString()
	isOrigin, _ := bc.GetParam("is_origin", 0).ToInt()
	tmp := strings.Split(tags, ",")
	var tids []int
	for _, t := range tmp {
		td, _ := negoutils.MakeElemType(t).ToInt()
		tids = append(tids, td)
	}
	isOri := false
	if isOrigin > 0 {
		isOri = true
	}
	articleId, b, e := model.MArticle.AddArticleInfo(title, "", tids, isOri)
	if e != nil {
		bc.Notice(nil, 100100, e.Error())
		return
	}
	if !b {
		bc.Notice(nil, 100100, "创建文章失败")
		return
	}
	bc.Notice(articleId)
	return
}
