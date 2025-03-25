/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-14 17:34
 */
package ajax

import (
	"github.com/liuyongshuai/runtofu/model"
)

type AdminAjaxChangePasswdController struct {
	AdminAjaxBaseController
}

//修改密码
func (bc *AdminAjaxChangePasswdController) Run() {
	if bc.UserInfo.Uid <= 0 {
		bc.Notice(nil, 100100, "修改密码失败：获取用户信息失败")
		return
	}
	oldPasswd := bc.GetParam("old_passwd", "").ToString()
	newPasswd := bc.GetParam("new_passwd", "").ToString()
	if len(oldPasswd) <= 0 {
		bc.Notice(nil, 100100, "修改密码失败：原密码不能为空")
		return
	}
	if len(newPasswd) <= 0 {
		bc.Notice(nil, 100100, "修改密码失败：新密码不能为空")
		return
	}
	b, e := model.MAdminUser.ChangePasswd(bc.UserInfo.Uid, oldPasswd, newPasswd)
	if !b {
		errmsg := "修改密码失败"
		if e != nil {
			errmsg += e.Error()
		}
		bc.Notice(nil, 100100, errmsg)
		return
	}
	bc.Notice(nil)
}
