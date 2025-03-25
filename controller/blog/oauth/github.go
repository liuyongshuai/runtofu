/**
 * @author      Liu Yongshuai
 * @package     oauth
 * @date        2018-03-14 22:35
 */
package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/liuyongshuai/negoutils"
	"github.com/liuyongshuai/runtofu/controller/blog"
	"github.com/liuyongshuai/runtofu/model"
	"strconv"
	"strings"
	"time"
)

type GithubOauthController struct {
	blog.RunToFuBaseController
}

// 运行主逻辑
func (bc *GithubOauthController) Run() {
	bc.TplName = "oauth.tpl"
	code := bc.GetParam("code", "").ToString()
	tokenInfo, err := model.MGithubApi.GetAccessToken(code)
	if err != nil {
		fmt.Println("ERROR\t", err)
		bc.RenderHtml()
		return
	}

	//提取用户信息
	ghUserInfo, rawJson, err := model.MGithubApi.GetUserInfo(tokenInfo.AccessToken)
	if err != nil {
		fmt.Println("ERROR\t", err)
		bc.RenderHtml()
		return
	}

	localGithubUserInfo := model.GithubUserInfo{
		GithubUid:   ghUserInfo.Uid,
		LoginName:   ghUserInfo.LoginName,
		AvatarUrl:   ghUserInfo.AvatarUrl,
		HtmlUrl:     ghUserInfo.ProfileUrl,
		BIO:         ghUserInfo.Desc,
		RawJson:     rawJson,
		AccessToken: tokenInfo.AccessToken,
	}

	localUid, _ := model.MSnowFlake.NextId()

	//写用户信息入库
	cond := make(map[string]interface{})
	cond["github_uid"] = localGithubUserInfo.GithubUid
	total := model.MGithubUser.GetGithubUserTotal(cond)
	if total <= 0 {
		err = model.MGithubUser.AddGithubUserInfo(localGithubUserInfo)
		if err != nil {
			fmt.Println("ERROR\t", err, localGithubUserInfo)
			bc.RenderHtml()
			return
		}
	} else {
		data := make(map[string]interface{})
		data["avatar_url"] = localGithubUserInfo.AvatarUrl
		data["html_url"] = localGithubUserInfo.HtmlUrl
		data["bio"] = localGithubUserInfo.BIO
		data["access_token"] = localGithubUserInfo.AccessToken
		data["raw_json"] = localGithubUserInfo.RawJson
		data["access_token"] = localGithubUserInfo.AccessToken
		err = model.MGithubUser.UpdateGithubUserInfo(localGithubUserInfo.GithubUid, data)
		if err != nil {
			fmt.Println("ERROR\t", err, localGithubUserInfo)
			bc.RenderHtml()
			return
		}
	}

	//写入通用库
	thirdUidStr := strconv.FormatInt(localGithubUserInfo.GithubUid, 10)
	localUInfo := model.RuntofuUserInfo{
		Uid:        localUid,
		Name:       localGithubUserInfo.LoginName,
		Portrait:   localGithubUserInfo.AvatarUrl,
		ProfileUrl: localGithubUserInfo.HtmlUrl,
		Type:       2,
		ThirdUid:   thirdUidStr,
	}
	cond = make(map[string]interface{})
	cond["third_uid"] = localGithubUserInfo.GithubUid
	cond["type"] = 2
	total = model.MRuntofuUser.GetRuntofuUserTotal(cond)
	if total <= 0 {
		err = model.MRuntofuUser.AddRuntofuUserInfo(localUInfo)
		if err != nil {
			fmt.Println("ERROR\t", err)
			bc.RenderHtml()
			return
		}
	} else {
		tmpInfo, _ := model.MRuntofuUser.GetRuntofuUserList(cond, 1, 1)
		if len(tmpInfo) > 0 {
			data := make(map[string]interface{})
			data["portrait"] = localGithubUserInfo.AvatarUrl
			data["profile_url"] = localUInfo.ProfileUrl
			localUInfo = tmpInfo[0]
			model.MRuntofuUser.UpdateRuntofuUserInfo(localUInfo.Uid, data)
		}
	}

	//登录成功，设置一下cookie
	//前面16位随机数 + uid的base62 + 前面几项的md5值。别问为啥这么做，就是特么的感觉牛X些
	expireTime := time.Now().Unix() + 86400*300
	cookieVal := negoutils.RandomStr(16) + negoutils.Base62Encode(localUInfo.Uid)
	md5 := strings.ToUpper(negoutils.MD5(cookieVal))
	cookieVal += md5
	jsUinfo, err := json.Marshal(model.BlogCookieInfo{
		RuntofuUid: localUInfo.Uid,
		CookieVal:  cookieVal,
		Expire:     expireTime,
	})
	if err != nil {
		fmt.Println("ERROR\t", err)
		bc.RenderHtml()
		return
	}

	//将相关信息存入redis
	rkey := fmt.Sprintf(model.BLOG_SESSION_PREFIX, localUInfo.Uid)
	rconn := model.RedisPool.Get()
	defer rconn.Close()
	_, e := rconn.Do("setex", rkey, expireTime, string(jsUinfo))
	if e != nil {
		fmt.Println("redis.Do setex", e, string(jsUinfo))
		fmt.Println("ERROR\t", err)
		bc.RenderHtml()
		return
	}

	//写入cookie信息
	fmt.Println("start set cookie.......")
	bc.Ctx.Output.AddCookie(
		model.BlogCookieKey,
		cookieVal,
		expireTime,
		"/",                   //path
		bc.Ctx.Input.Domain(), //domain
		true,                  //secure
	)
	bc.RenderHtml()
}
