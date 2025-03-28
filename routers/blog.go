/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     conf
 * @date        2018-02-03 15:34
 */
package routers

import (
	"github.com/liuyongshuai/runtofu/controller/blog"
	"github.com/liuyongshuai/runtofu/controller/blog/ajax"
	"github.com/liuyongshuai/runtofu/controller/blog/oauth"
	"github.com/liuyongshuai/runtofu/negoutils"
)

var BlogRouterList []*negoutils.RuntofuRouterItem

func init() {
	BlogRouterList = append(BlogRouterList,
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "test",
			Controller: &blog.TestController{},
		},
		//文章详细页面
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypeRegexp,
			Config:     `^article/([\d]+)$`,
			Param:      "aid=$1",
			Controller: &blog.ArticleController{},
		},
		//首页
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "",
			Controller: &blog.IndexController{},
		},
		//推荐文章列表页
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "reclist",
			Controller: &blog.RecListController{},
		},
		//留言板
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "feedback",
			Controller: &blog.FeedBackController{},
		},
		//全部话题标签列表
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypePathInfo,
			Config:     "taglist",
			Controller: &blog.TabListController{},
		},
		//标签话题下面的文章列表
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypeRegexp,
			Config:     `^tag/([\d]+)$`,
			Param:      "tag_id=$1",
			Controller: &blog.TagController{},
		},
		//Ajax请求：获取文章详细信息
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypeRegexp,
			Config:     `^ajax/article/info/([\d]+)$`,
			Param:      "aid=$1",
			Controller: &ajax.RunToFuAjaxArticleInfoController{},
		},
		//oauth登录授权：微博：oauth/weibo/auth、oauth/weibo/unauth
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypeRegexp,
			Config:     `^oauth/weibo/([a-zA-Z]+)$`,
			Param:      "action=$1",
			Controller: &oauth.WeiboOauthController{},
		},
		//oauth登录授权：github.com：oauth/github/auth、oauth/github/unauth
		&negoutils.RuntofuRouterItem{
			Type:       negoutils.RouterTypeRegexp,
			Config:     `^oauth/github/([a-zA-Z]+)$`,
			Param:      "action=$1",
			Controller: &oauth.GithubOauthController{},
		},
	)
}
