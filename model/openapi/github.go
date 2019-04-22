/**
 * @author      Liu Yongshuai
 * @package     openapi
 * @date        2018-03-17 12:15
 */
package openapi

//https://developer.github.com/apps/building-oauth-apps/authorization-options-for-oauth-apps/

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuyongshuai/goUtils"
	"github.com/liuyongshuai/runtofu/configer"
	"strings"
)

type GithubOpenApi struct {
	Conf configer.OauthGithubConf
}

//初始化各种配置信息
func (gh *GithubOpenApi) InitConf(conf configer.OauthGithubConf) {
	gh.Conf = conf
}

//获取要跳到oauth认证的页面地址
func (gh *GithubOpenApi) GetAuthorizeUrl() string {
	retUrl := "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email,read:user"
	retUrl = fmt.Sprintf(retUrl, gh.Conf.ClientId, gh.Conf.CallBackUrl)
	return retUrl
}

type GithubAccessToken struct {
	AccessToken string `json:"access_token"`
}

//获取授权过的Access Token
func (gh *GithubOpenApi) GetAccessToken(authCode string) (ret GithubAccessToken, err error) {
	apiUrl := "https://github.com/login/oauth/access_token"
	httpReq := goUtils.NewHttpClient(apiUrl, context.Background())
	httpReq.SetUrl(apiUrl)
	httpReq.AddField("client_id", gh.Conf.ClientId)
	httpReq.AddField("client_secret", gh.Conf.ClientSecret)
	httpReq.AddField("code", authCode)
	httpReq.AddField("redirect_uri", gh.Conf.CallBackUrl)
	resp, err := httpReq.Post()
	if err != nil {
		fmt.Println("ERROR_GetAccessToken", err)
		return
	}
	jsonStr := resp.GetBodyString()
	fmt.Println(apiUrl, jsonStr)
	tmp := strings.Split(jsonStr, "&")
	for _, t := range tmp {
		if !strings.Contains(t, "access_token") {
			continue
		}
		tt := strings.Split(t, "=")
		if len(tt) == 2 {
			ret.AccessToken = tt[1]
		}
	}
	return
}

type ApiGithubUserInfo struct {
	LoginName  string `json:"login"`
	Uid        int64  `json:"id"`
	AvatarUrl  string `json:"avatar_url"`
	ProfileUrl string `json:"html_url"`
	Desc       string `json:"bio"`
}

//获取用户信息
func (gh *GithubOpenApi) GetUserInfo(token string) (ret ApiGithubUserInfo, rawJson string, err error) {
	apiUrl := "https://api.github.com/user?access_token=%s"
	apiUrl = fmt.Sprintf(apiUrl, token)
	httpReq := goUtils.NewHttpClient(apiUrl, context.Background())
	httpReq.SetUrl(apiUrl)
	resp, err := httpReq.Get()
	if err != nil {
		fmt.Println("ERROR_GetUserInfo", err)
		return
	}
	rawJson = resp.GetBodyString()
	fmt.Println(apiUrl, rawJson)
	err = json.Unmarshal([]byte(rawJson), &ret)
	return
}
