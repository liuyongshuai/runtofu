/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-09 17:32
 */
package ajax

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/model"
)

type AdminAjaxArticleInfoController struct {
	AdminAjaxBaseController
}

//返回数据信息
func (bc *AdminAjaxArticleInfoController) Run() {
	aid, _ := bc.GetParam("article_id", 0).ToInt64()
	if aid <= 0 {
		bc.Notice(nil)
		return
	}
	ainfo, err := model.MArticle.GetArticleInfo(aid)
	if err != nil {
		bc.Notice(nil, 100100, "查询出错 %v", err)
		fmt.Println(err)
		return
	}
	bc.Notice(ainfo)
}
