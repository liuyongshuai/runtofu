/**
 * @author      Liu Yongshuai
 * @package     routers
 * @date        2018-02-07 20:56
 */
package routers

import (
	"github.com/liuyongshuai/runtofu/controller"
	"github.com/liuyongshuai/runtofu/controller/ajax"
	"github.com/liuyongshuai/runtofu/controller/system"
	"github.com/liuyongshuai/runtofu/goweb/router"
)

var AdminRouterList []*router.WeGoRouterItem

func init() {
	AdminRouterList = append(AdminRouterList,
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "test",
			Controller: &controller.TestController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "taglist",
			Controller: &controller.TagListController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "articlelist",
			Controller: &controller.ArticleListController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "articleinfo",
			Controller: &controller.ArticleInfoController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "login",
			Controller: &controller.LoginController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "",
			Controller: &controller.IndexController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/article/info",
			Controller: &ajax.AdminAjaxArticleInfoController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/article/modify",
			Controller: &ajax.AdminAjaxArticleModifyController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/article/new",
			Controller: &ajax.AdminAjaxArticleNewController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypeRegexp,
			Config:     "ajax/tag/(.+)",
			Param:      "action=$1",
			Controller: &ajax.AdminAjaxTagController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypeRegexp,
			Config:     "ajax/system/(.+)",
			Param:      "action=$1",
			Controller: &ajax.AdminAjaxSystemController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/login",
			Controller: &ajax.AdminAjaxLoginController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/logout",
			Controller: &ajax.AdminAjaxLogoutController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "ajax/changepasswd",
			Controller: &ajax.AdminAjaxChangePasswdController{},
		},
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "system/menu",
			Controller: &system.AdminSystemMenuController{},
		},
		//通过oauth登录的用户
		&router.WeGoRouterItem{
			Type:       router.RouterTypePathInfo,
			Config:     "oauth/userlist",
			Controller: &controller.OauthUserController{},
		},
	)
}
