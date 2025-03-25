/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     ajax
 * @date        2018-02-03 13:31
 */
package ajax

import (
	"fmt"
	"github.com/liuyongshuai/goUtils"
	"github.com/liuyongshuai/runtofu/controller/admin"
	"strings"
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
type AdminAjaxBaseController struct {
	admin.AdminBaseController
}

//校验是否为Ajax请求
func (bc *AdminAjaxBaseController) Prepare() error {
	bc.UserInfo = bc.CheckLogin(true, func() {
		bc.Notice(nil, 100100, "登录校验失败，请重新登录")
	})
	if bc.UserInfo.Uid <= 0 {
		return fmt.Errorf("登录失败")
	}
	return nil
}

//返回json数据
func (bc *AdminAjaxBaseController) Notice(d interface{}, ret ...interface{}) {
	var errno int64 = 0
	var errmsg = "ok"
	f := ""
	if len(ret) > 0 {
		errno, _ = goUtils.MakeElemType(ret[0]).ToInt64()
	}
	if len(ret) > 1 {
		errmsg = goUtils.MakeElemType(ret[1]).ToString()
	}
	if len(ret) > 2 {
		f = goUtils.MakeElemType(ret[2]).ToString()
	}
	if strings.Count(errmsg, "%") > 0 {
		errmsg = fmt.Sprintf(errmsg, f)
	}
	bc.RenderJson(AjaxResponseData{
		Result: AjaxResponseResult{
			Errno:  errno,
			ErrMsg: errmsg,
		},
		Data: d,
	})
}
