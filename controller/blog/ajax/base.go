/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     ajax
 * @date        2018-02-03 13:31
 */
package ajax

import (
	"fmt"
	"github.com/liuyongshuai/goUtils"
	"github.com/liuyongshuai/runtofu/controller/blog"
	"net/http"
)

//ajax响应信息里result部分
type AjaxResponseResult struct {
	Errno  int64  `json:"errno"`  //返回的错误码
	ErrMsg string `json:"errmsg"` //返回的错误信息
}

//ajax响应总体返回的数据格式
type AjaxResponseData struct {
	Result AjaxResponseResult `json:"result"`
	Data   interface{}        `json:"data"`
}

//Ajax层的基类
type RunToFuAjaxBaseController struct {
	blog.RunToFuBaseController
}

//校验是否为Ajax请求
func (bc *RunToFuAjaxBaseController) Prepare() error {
	if !bc.IsAjax() {
		bc.Notice(nil, 100404, "非法请求，请确认后重试")
		bc.SetStatus(http.StatusForbidden)
		return fmt.Errorf("invalid request")
	}
	return nil
}

//返回json数据
func (bc *RunToFuAjaxBaseController) Notice(d interface{}, ret ...interface{}) {
	var errno int64 = 0
	var errmsg = "ok"
	if len(ret) > 0 {
		errno, _ = goUtils.MakeElemType(ret[0]).ToInt64()
	}
	if len(ret) > 1 {
		errmsg = goUtils.MakeElemType(ret[1]).ToString()
	}
	bc.RenderJson(AjaxResponseData{
		Result: AjaxResponseResult{
			Errno:  errno,
			ErrMsg: errmsg,
		},
		Data: d,
	})
}
