/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     ajax
 * @date        2018-02-03 13:31
 */
package ajax

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/controller"
	"github.com/liuyongshuai/runtofu/controller/admin"
	"github.com/liuyongshuai/runtofu/negoutils"
	"strings"
)

// Ajax层的基类
type AdminAjaxBaseController struct {
	admin.AdminBaseController
}

// 校验是否为Ajax请求
func (bc *AdminAjaxBaseController) Prepare() error {
	bc.UserInfo = bc.CheckLogin(true, func() {
		bc.Notice(nil, 100100, "登录校验失败，请重新登录")
	})
	if bc.UserInfo.Uid <= 0 {
		return fmt.Errorf("登录失败")
	}
	return nil
}

// 返回json数据
func (bc *AdminAjaxBaseController) Notice(d interface{}, ret ...interface{}) {
	var errno int64 = 0
	var errmsg = "ok"
	f := ""
	if len(ret) > 0 {
		errno, _ = negoutils.MakeElemType(ret[0]).ToInt64()
	}
	if len(ret) > 1 {
		errmsg = negoutils.MakeElemType(ret[1]).ToString()
	}
	if len(ret) > 2 {
		f = negoutils.MakeElemType(ret[2]).ToString()
	}
	if strings.Count(errmsg, "%") > 0 {
		errmsg = fmt.Sprintf(errmsg, f)
	}
	bc.RenderJson(controller.AjaxResponseResult{
		Errcode: errno,
		ErrMsg:  errmsg,
		Data:    d,
	})
}
