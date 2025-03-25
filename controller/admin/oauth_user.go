/**
 * @author      Liu Yongshuai
 * @package     admin
 * @date        2018-03-21 12:06
 */
package admin

type OauthUserController struct {
	AdminBaseController
}

//校验是否登录
func (bc *OauthUserController) Run() {
	bc.TplName = "oauth_user.tpl"
	bc.RenderHtml()
}
