/**
 * @author      Liu Yongshuai
 * @package     openapi
 * @date        2018-03-17 12:15
 */
package openapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liuyongshuai/negoutils"
	"github.com/liuyongshuai/runtofu/confutils"
	"strings"
)

type WeiboOpenApi struct {
	Conf confutils.OauthWeiboConf
}

// 初始化各种配置信息
func (wb *WeiboOpenApi) InitConf(conf confutils.OauthWeiboConf) {
	wb.Conf = conf
}

// 获取要跳到oauth认证的页面地址
// http://open.weibo.com/wiki/Oauth2/authorize
func (wb *WeiboOpenApi) GetAuthorizeUrl() string {
	retUrl := wb.Conf.ApiUrl + "/oauth2/authorize"
	retUrl = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=email,follow_app_official_microblog",
		retUrl, wb.Conf.AppKey, wb.Conf.CallBackUrl)
	return retUrl
}

type WeiboAccessToken struct {
	AccessToken string `json:"access_token"` //用户授权的唯一票据，用于调用微博的开放接口，同时也是第三方应用验证微博用户登录的唯一票据，第三方应用应该用该票据和自己应用内的用户建立唯一影射关系，来识别登录状态，不能使用本返回值里的UID字段来做登录识别。
	ExpiresIn   int64  `json:"expires_in"`   //access_token的生命周期，单位是秒数。
	Uid         string `json:"uid"`          //授权用户的UID，本字段只是为了方便开发者，减少一次user/show接口调用而返回的，第三方应用不能用此字段作为用户登录状态的识别，只有access_token才是用户授权的唯一票据。
	IsRealName  string `json:"isRealName"`
}

// 获取授权过的Access Token
// http://open.weibo.com/wiki/Oauth2/access_token
func (wb *WeiboOpenApi) GetAccessToken(authCode string) (ret WeiboAccessToken, err error) {
	ret = WeiboAccessToken{}
	url := wb.Conf.ApiUrl + "/oauth2/access_token"
	httpReq := negoutils.NewHttpClient(url, context.Background())
	httpReq.SetUrl(url)
	httpReq.AddField("client_id", wb.Conf.AppKey)
	httpReq.AddField("client_secret", wb.Conf.AppSecret)
	httpReq.AddField("grant_type", "authorization_code")
	httpReq.AddField("code", authCode)
	httpReq.AddField("redirect_uri", wb.Conf.CallBackUrl)
	resp, err := httpReq.Post()
	if err != nil {
		fmt.Println("ERROR_GetAccessToken", err)
		return
	}
	jsonStr := resp.GetBodyString()
	fmt.Println(jsonStr)
	if strings.Contains(jsonStr, "error") {
		err = errors.New(jsonStr)
		return
	}
	err = json.Unmarshal([]byte(jsonStr), &ret)
	return
}

type WeiboTokenInfo struct {
	Uid      string `json:"uid"`
	AppKey   string `json:"appkey"`
	CreateAt string `json:"create_at"`
	ExpireIn string `json:"expire_in"`
}

// 查询用户access_token的授权相关信息，包括授权时间，过期时间和scope权限。
// http://open.weibo.com/wiki/Oauth2/get_token_info
func (wb *WeiboOpenApi) GetTokenInfo(token string) (ret WeiboTokenInfo, err error) {
	ret = WeiboTokenInfo{}
	url := wb.Conf.ApiUrl + "/oauth2/get_token_info"
	httpReq := negoutils.NewHttpClient(url, context.Background())
	httpReq.SetUrl(url)
	httpReq.AddField("access_token", token)
	resp, err := httpReq.Post()
	if err != nil {
		fmt.Println("ERROR_GetTokenInfo", err)
		return
	}
	jsonStr := resp.GetBodyString()
	if strings.Contains(jsonStr, "error") {
		err = errors.New(jsonStr)
		return
	}
	err = json.Unmarshal([]byte(jsonStr), &ret)
	return
}

type ApiWeiboUserInfo struct {
	Uid             int64  `json:"id"`                //用户ID
	UidStr          string `json:"idstr"`             //用户ID的字符串形式
	ScreenName      string `json:"screen_name"`       //用户昵称
	Name            string `json:"name"`              //友好显示名称
	Desc            string `json:"description"`       //用户个人描述
	ProfileImageURl string `json:"profile_image_url"` //用户头像地址，50×50像素
	ProfileUrl      string `json:"profile_url"`       //用户的微博统一URL地址
}

// 根据用户ID获取用户信息
// http://open.weibo.com/wiki/2/users/show
func (wb *WeiboOpenApi) GetUserInfo(token, uid string) (ret ApiWeiboUserInfo, rawJson string, err error) {
	ret = ApiWeiboUserInfo{}
	url := wb.Conf.ApiUrl + "/2/users/show.json"
	url = fmt.Sprintf("%s?uid=%s&access_token=%s", url, uid, token)
	httpReq := negoutils.NewHttpClient(url, context.Background())
	httpReq.SetUrl(url)
	httpReq.SetTimeout(10)
	resp, err := httpReq.Get()
	if err != nil {
		fmt.Println("ERROR_GetUserInfo", err)
		return
	}
	rawJson = resp.GetBodyString()
	if strings.Contains(rawJson, "error") {
		err = errors.New(rawJson)
		return
	}
	err = json.Unmarshal([]byte(rawJson), &ret)
	return
}
