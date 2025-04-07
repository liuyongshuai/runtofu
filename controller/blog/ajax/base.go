/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     ajax
 * @date        2018-02-03 13:31
 */
package ajax

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/controller"
	"github.com/liuyongshuai/runtofu/controller/blog"
	"github.com/liuyongshuai/runtofu/negoutils"
	"net/http"
)

// Ajax层的基类
type RunToFuAjaxBaseController struct {
	blog.RunToFuBaseController
}

// 校验是否为Ajax请求
func (bc *RunToFuAjaxBaseController) Prepare() error {
	if !bc.IsAjax() {
		bc.Notice(nil, 100404, "非法请求，请确认后重试")
		bc.SetStatus(http.StatusForbidden)
		return fmt.Errorf("invalid request")
	}
	return nil
}

// 返回json数据
func (bc *RunToFuAjaxBaseController) Notice(d interface{}, ret ...interface{}) {
	var errno int64 = 0
	var errmsg = "ok"
	if len(ret) > 0 {
		errno, _ = negoutils.MakeElemType(ret[0]).ToInt64()
	}
	if len(ret) > 1 {
		errmsg = negoutils.MakeElemType(ret[1]).ToString()
	}
	bc.RenderJson(controller.AjaxResponseResult{
		Errcode: errno,
		ErrMsg:  errmsg,
		Data:    d,
	})
}
