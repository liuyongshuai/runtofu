/**
 * @author      Liu Yongshuai
 * @package     routers
 * @date        2018-02-07 20:56
 */
package routers

import (
	"github.com/liuyongshuai/runtofu/controller/admin"
	"github.com/liuyongshuai/runtofu/controller/admin/ajax"
	"github.com/liuyongshuai/runtofu/controller/admin/system"
	"github.com/liuyongshuai/runtofu/goweb/router"
)

var AdminRouterList []*router.RuntofuRouterItem

func init() {
	AdminRouterList = append(AdminRouterList,
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "test",
			Controller: &admin.TestController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "taglist",
			Controller: &admin.TagListController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "articlelist",
			Controller: &admin.ArticleListController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "articleinfo",
			Controller: &admin.ArticleInfoController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "login",
			Controller: &admin.LoginController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "",
			Controller: &admin.IndexController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/article/info",
			Controller: &ajax.AdminAjaxArticleInfoController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/article/modify",
			Controller: &ajax.AdminAjaxArticleModifyController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/article/new",
			Controller: &ajax.AdminAjaxArticleNewController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypeRegexp,
			Config:     "ajax/tag/(.+)",
			Param:      "action=$1",
			Controller: &ajax.AdminAjaxTagController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypeRegexp,
			Config:     "ajax/system/(.+)",
			Param:      "action=$1",
			Controller: &ajax.AdminAjaxSystemController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/login",
			Controller: &ajax.AdminAjaxLoginController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/logout",
			Controller: &ajax.AdminAjaxLogoutController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/changepasswd",
			Controller: &ajax.AdminAjaxChangePasswdController{},
		},
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "system/menu",
			Controller: &system.AdminSystemMenuController{},
		},
		//通过oauth登录的用户
		&router.RuntofuRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "oauth/userlist",
			Controller: &admin.OauthUserController{},
		},
	)
}
