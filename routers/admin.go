/**
 * @author      Liu Yongshuai
 * @package     routers
 * @date        2018-02-07 20:56
 */
package routers

import (
	"github.com/liuyongshuai/negoutils"
	"github.com/liuyongshuai/runtofu/controller/admin"
	"github.com/liuyongshuai/runtofu/controller/admin/ajax"
	"github.com/liuyongshuai/runtofu/controller/admin/system"
)

var AdminRouterList []*negoutils.RuntofuRouterItem

func init() {
	AdminRouterList = append(AdminRouterList,
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "test",
			Controller: &admin.TestController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "taglist",
			Controller: &admin.TagListController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "articlelist",
			Controller: &admin.ArticleListController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "articleinfo",
			Controller: &admin.ArticleInfoController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "login",
			Controller: &admin.LoginController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "",
			Controller: &admin.IndexController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "ajax/article/info",
			Controller: &ajax.AdminAjaxArticleInfoController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "ajax/article/modify",
			Controller: &ajax.AdminAjaxArticleModifyController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "ajax/article/new",
			Controller: &ajax.AdminAjaxArticleNewController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypeRegexp,
			Config:     "ajax/tag/(.+)",
			Param:      "action=$1",
			Controller: &ajax.AdminAjaxTagController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypeRegexp,
			Config:     "ajax/system/(.+)",
			Param:      "action=$1",
			Controller: &ajax.AdminAjaxSystemController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "ajax/login",
			Controller: &ajax.AdminAjaxLoginController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "ajax/logout",
			Controller: &ajax.AdminAjaxLogoutController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "ajax/changepasswd",
			Controller: &ajax.AdminAjaxChangePasswdController{},
		},
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "system/menu",
			Controller: &system.AdminSystemMenuController{},
		},
		//通过oauth登录的用户
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "oauth/userlist",
			Controller: &admin.OauthUserController{},
		},
	)
}
