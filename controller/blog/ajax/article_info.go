/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     ajax
 * @date        2018-02-03 13:39
 */
package ajax

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/model"
)

type RunToFuAjaxArticleInfoController struct {
	RunToFuAjaxBaseController
}

//返回数据信息
func (bc *RunToFuAjaxArticleInfoController) Run() {
	aid, _ := bc.GetParam("aid", 0).ToInt64()
	if aid <= 0 {
		bc.Notice(nil, 100100, "文章ID非法")
		return
	}
	ainfo, err := model.MArticle.GetArticleInfo(aid)
	if err != nil {
		bc.Notice(nil, 100100, "查询出错")
		fmt.Println(err)
		return
	}
	bc.Notice(ainfo)
}
