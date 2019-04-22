/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     controller
 * @date        2018-02-03 15:48
 */
package blog

import (
	"github.com/liuyongshuai/runtofu/model"
)

type TabListController struct {
	RunToFuBaseController
}

//这些标签话题列表是全部显示出来的
func (bc *TabListController) Run() {
	bc.TplName = "taglist.tpl"
	total := model.MTag.GetTagTotal()
	tagList, _ := model.MTag.GetTagList(1, int(total))
	bc.TplData["tagList"] = tagList
	bc.RenderHtml()
}
