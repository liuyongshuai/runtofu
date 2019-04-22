/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-10 22:18
 */
package ajax

import (
	"github.com/liuyongshuai/runtofu/model"
)

type AdminAjaxTagController struct {
	AdminAjaxBaseController
}

//返回数据信息
func (bc *AdminAjaxTagController) Run() {
	action := bc.GetParam("action", "").ToString()
	tagId, _ := bc.GetParam("tag_id", 0).ToInt()
	tagName := bc.GetParam("tag_name", "").ToString()
	switch action {
	case "add":
		tid, b, e := model.MTag.AddTagInfo(tagName)
		if tid <= 0 || e != nil || !b {
			bc.Notice(nil, 100100, "创建话题失败,%s", e.Error())
			return
		}
		tInfo, _ := model.MTag.GetTagInfo(tid)
		bc.Notice(tInfo)
	case "delete":
		b, e := model.MTag.DeleteTagInfo(tagId)
		if e != nil || !b {
			bc.Notice(nil, 100100, "删除话题失败 %s", e.Error())
			return
		}
		bc.Notice(nil)
	default:
		bc.Notice(nil, 100100, "非法操作")
	}
}
