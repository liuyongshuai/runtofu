/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-03-20 23:10
 */
package model

import (
	"errors"
	"fmt"
	"github.com/liuyongshuai/negoutils"
)

// 实例化一个m层
func NewGithubUserModel() *GithubUserModel {
	ret := &GithubUserModel{}
	ret.Table = "github_user"
	return ret
}

// 详细信息
type GithubUserInfo struct {
	GithubUid   int64  `json:"github_uid"`
	LoginName   string `json:"login_name"`
	AvatarUrl   string `json:"avatar_url"`
	HtmlUrl     string `json:"html_url"`
	BIO         string `json:"bio"`
	RawJson     string `json:"raw_json"`
	AccessToken string `json:"access_token"`
}

type GithubUserModel struct {
	BaseModel
}

// 增
func (m *GithubUserModel) AddGithubUserInfo(uinfo GithubUserInfo) (err error) {
	if uinfo.GithubUid <= 0 {
		return errors.New("invalid Github uid")
	}
	uData := make(map[string]interface{})
	uData["github_uid"] = uinfo.GithubUid
	uData["login_name"] = uinfo.LoginName
	uData["avatar_url"] = uinfo.AvatarUrl
	uData["html_url"] = uinfo.HtmlUrl
	uData["bio"] = uinfo.BIO
	uData["raw_json"] = uinfo.RawJson
	uData["access_token"] = uinfo.AccessToken
	_, b, e := mDB.InsertData(m.Table, uData, false)
	if e != nil || !b {
		return fmt.Errorf("insert github uinfo data failed")
	}
	return nil
}

// 改
func (m *GithubUserModel) UpdateGithubUserInfo(uid int64, data map[string]interface{}) (err error) {
	if uid <= 0 {
		return fmt.Errorf("invalid Github user id")
	}
	uinfo, err := m.GetGithubUserInfo(uid)
	if err != nil || uinfo.GithubUid <= 0 {
		return fmt.Errorf("github user info not exists")
	}
	cond := make(map[string]interface{})
	cond["github_uid"] = uid
	_, b, e := mDB.UpdateData(m.Table, data, cond)
	if e != nil || !b {
		return fmt.Errorf("update github user info data failed")
	}
	return
}

// 查
func (m *GithubUserModel) GetGithubUserInfo(uid int64) (uinfo GithubUserInfo, err error) {
	if uid <= 0 {
		return
	}
	cond := make(map[string]interface{})
	cond["github_uid"] = uid
	row, err := m.FetchRow(cond)
	if err != nil || len(row) <= 0 {
		return
	}
	uinfo = formatGithubUserInfo(row)
	return
}

// 提取文章列表，按时间倒序排序
func (m *GithubUserModel) GetGithubUserList(cond map[string]interface{}, page, pagesize int) ([]GithubUserInfo, error) {
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `github_uid` DESC")
	var ret []GithubUserInfo
	for _, row := range rows {
		ret = append(ret, formatGithubUserInfo(row))
	}
	return ret, nil
}

// 提取文章总数
func (m *GithubUserModel) GetGithubUserTotal(cond map[string]interface{}) int64 {
	return m.FetchTotal(cond)
}

// 格式化文章信息
func formatGithubUserInfo(row map[string]negoutils.ElemType) (ret GithubUserInfo) {
	ret.GithubUid, _ = row["github_uid"].ToInt64()
	ret.LoginName = row["login_name"].ToString()
	ret.AvatarUrl = row["avatar_url"].ToString()
	ret.HtmlUrl = row["html_url"].ToString()
	ret.BIO = row["bio"].ToString()
	ret.RawJson = row["raw_json"].ToString()
	ret.AccessToken = row["access_token"].ToString()
	return
}
