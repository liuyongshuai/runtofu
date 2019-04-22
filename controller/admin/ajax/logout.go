/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-13 13:58
 */
package ajax

type AdminAjaxLogoutController struct {
	AdminAjaxBaseController
}

//返回数据信息
func (bc *AdminAjaxLogoutController) Run() {
	adminUserInfo := bc.CheckLogin(true, func() {
		bc.Notice(nil, 100100, "登录校验失败，请重新登录")
	})
	if adminUserInfo.Uid <= 0 {
		return
	}
	bc.Logout(adminUserInfo.Uid)
	bc.Notice(nil)
}
