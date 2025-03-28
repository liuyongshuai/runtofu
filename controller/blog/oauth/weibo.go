/**
 * @author      Liu Yongshuai
 * @package     oauth
 * @date        2018-03-14 22:35
 */
package oauth

import (
	"encoding/json"
	"fmt"
	"github.com/liuyongshuai/runtofu/controller/blog"
	"github.com/liuyongshuai/runtofu/model"
	"github.com/liuyongshuai/runtofu/negoutils"
	"strconv"
	"strings"
	"time"
)

type WeiboOauthController struct {
	blog.RunToFuBaseController
}

// 运行主逻辑
func (bc *WeiboOauthController) Run() {
	bc.TplName = "oauth.tpl"
	code := bc.GetParam("code", "").ToString()
	atInfo, err := model.MWeiboApi.GetAccessToken(code)
	if err != nil {
		fmt.Println("ERROR\t", err)
		bc.RenderHtml()
		return
	}

	//提取微博的用户信息
	wbUserInfo, rawJson, err := model.MWeiboApi.GetUserInfo(atInfo.AccessToken, atInfo.Uid)
	if err != nil {
		fmt.Println("ERROR\t", err)
		bc.RenderHtml()
		return
	}
	localWeiboUserInfo := model.WeiboUserInfo{
		WbUid:                 wbUserInfo.Uid,
		ScreenName:            wbUserInfo.ScreenName,
		Name:                  wbUserInfo.Name,
		Desc:                  wbUserInfo.Desc,
		ProfileImageUrl:       wbUserInfo.ProfileImageURl,
		ProfileUrl:            wbUserInfo.ProfileUrl,
		RawJson:               rawJson,
		AccessToken:           atInfo.AccessToken,
		AccessTokenExpireTime: time.Now().Unix() + atInfo.ExpiresIn,
	}

	localUid, _ := model.MSnowFlake.NextId()

	//写用户信息入库
	cond := make(map[string]interface{})
	cond["wb_uid"] = localWeiboUserInfo.WbUid
	total := model.MWeiboUser.GetWeiboUserTotal(cond)
	if total <= 0 {
		err = model.MWeiboUser.AddWeiboUserInfo(localWeiboUserInfo)
		if err != nil {
			fmt.Println("ERROR\t", err, localWeiboUserInfo)
			bc.RenderHtml()
			return
		}
	} else {
		data := make(map[string]interface{})
		data["screen_name"] = localWeiboUserInfo.ScreenName
		data["name"] = localWeiboUserInfo.Name
		data["profile_image_url"] = localWeiboUserInfo.ProfileImageUrl
		data["raw_json"] = localWeiboUserInfo.RawJson
		data["access_token"] = localWeiboUserInfo.AccessToken
		data["access_token_expire_time"] = localWeiboUserInfo.AccessTokenExpireTime
		err = model.MWeiboUser.UpdateWeiboUserInfo(localWeiboUserInfo.WbUid, data)
		if err != nil {
			fmt.Println("ERROR\t", err, localWeiboUserInfo)
			bc.RenderHtml()
			return
		}
	}

	//写入通用库
	thirdUidStr := strconv.FormatInt(localWeiboUserInfo.WbUid, 10)
	localUInfo := model.RuntofuUserInfo{
		Uid:        localUid,
		Name:       localWeiboUserInfo.ScreenName,
		Portrait:   localWeiboUserInfo.ProfileImageUrl,
		ProfileUrl: "https://weibo.com/" + localWeiboUserInfo.ProfileUrl,
		Type:       1,
		ThirdUid:   thirdUidStr,
	}
	cond = make(map[string]interface{})
	cond["third_uid"] = localWeiboUserInfo.WbUid
	cond["type"] = 1
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
			data["name"] = localWeiboUserInfo.ScreenName
			data["portrait"] = localWeiboUserInfo.ProfileImageUrl
			data["profile_url"] = localUInfo.ProfileUrl
			localUInfo = tmpInfo[0]
			model.MRuntofuUser.UpdateRuntofuUserInfo(localUInfo.Uid, data)
		}
	}

	//登录成功，设置一下cookie
	//前面16位随机数 + uid的base62 + 前面几项的md5值。别问为啥这么做，就是特么的感觉牛X些
	cookieVal := negoutils.RandomStr(16) + negoutils.Base62Encode(localUInfo.Uid)
	md5 := strings.ToUpper(negoutils.MD5(cookieVal))
	cookieVal += md5
	jsUinfo, err := json.Marshal(model.BlogCookieInfo{
		RuntofuUid: localUInfo.Uid,
		CookieVal:  cookieVal,
		Expire:     localWeiboUserInfo.AccessTokenExpireTime,
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
	_, e := rconn.Do("setex", rkey, localWeiboUserInfo.AccessTokenExpireTime, string(jsUinfo))
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
		localWeiboUserInfo.AccessTokenExpireTime,
		"/",                   //path
		bc.Ctx.Input.Domain(), //domain
		true,                  //secure
	)
	bc.RenderHtml()
}
