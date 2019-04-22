/**
 * @author      Liu Yongshuai
 * @package     admin
 * @date        2018-02-11 20:35
 */
package admin

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/model"
)

type LoginController struct {
	AdminBaseController
}

func (bc *LoginController) Prepare() error {
	return nil
}

//校验是否登录
func (bc *LoginController) Run() {
	bc.UserInfo = bc.CheckLogin(false, func() {})
	if bc.UserInfo.Uid > 0 {
		bc.Ctx.Redirect("/")
		return
	}
	bc.Ctx.Output.AddCookie(model.CookieKey, "", -1, "/")
	bc.SetLeftMenu()
	bc.TplName = "login.tpl"
	err := bc.RenderHtml()
	fmt.Println(err)
}
