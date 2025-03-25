/**
 * @author      Liu Yongshuai
 * @package     admin
 * @date        2018-02-09 11:14
 */
package admin

import (
	"github.com/liuyongshuai/runtofu/model"
)

type TagListController struct {
	AdminBaseController
}

//校验是否登录
func (bc *TagListController) Run() {
	bc.TplName = "taglist.tpl"
	total := model.MTag.GetTagTotal()
	tagList, _ := model.MTag.GetTagList(1, int(total))
	bc.TplData["tagList"] = tagList
	bc.RenderHtml()
}
