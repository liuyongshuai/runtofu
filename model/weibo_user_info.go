/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-03-18 21:43
 */
package model

import (
	"errors"
	"fmt"
	"github.com/liuyongshuai/runtofu/negoutils"
)

// 实例化一个m层
func NewWeiboUserModel() *WeiboUserModel {
	ret := &WeiboUserModel{}
	ret.Table = "weibo_user"
	return ret
}

// 详细信息
type WeiboUserInfo struct {
	WbUid                 int64  `json:"wb_uid"`
	ScreenName            string `json:"screen_name"`
	Name                  string `json:"name"`
	Desc                  string `json:"description"`
	ProfileImageUrl       string `json:"profile_image_url"`
	ProfileUrl            string `json:"profile_url"`
	RawJson               string `json:"raw_json"`
	AccessToken           string `json:"access_token"`
	AccessTokenExpireTime int64  `json:"access_token_expire_time"`
}

type WeiboUserModel struct {
	BaseModel
}

// 增
func (m *WeiboUserModel) AddWeiboUserInfo(uinfo WeiboUserInfo) (err error) {
	if uinfo.WbUid <= 0 {
		return errors.New("invalid weibo uid")
	}
	uData := make(map[string]interface{})
	uData["wb_uid"] = uinfo.WbUid
	uData["screen_name"] = uinfo.ScreenName
	uData["name"] = uinfo.Name
	uData["description"] = uinfo.Desc
	uData["profile_image_url"] = uinfo.ProfileImageUrl
	uData["profile_url"] = uinfo.ProfileUrl
	uData["raw_json"] = uinfo.RawJson
	uData["access_token"] = uinfo.AccessToken
	uData["access_token_expire_time"] = uinfo.AccessTokenExpireTime
	_, b, e := mDB.InsertData(m.Table, uData, false)
	if e != nil || !b {
		return fmt.Errorf("insert weibo uinfo data failed")
	}
	return nil
}

// 改
func (m *WeiboUserModel) UpdateWeiboUserInfo(uid int64, data map[string]interface{}) (err error) {
	if uid <= 0 {
		return fmt.Errorf("invalid weibo user id")
	}
	uinfo, err := m.GetWeiboUserInfo(uid)
	if err != nil || uinfo.WbUid <= 0 {
		return fmt.Errorf("weibo user info not exists")
	}
	cond := make(map[string]interface{})
	cond["wb_uid"] = uid
	_, b, e := mDB.UpdateData(m.Table, data, cond)
	if e != nil || !b {
		return fmt.Errorf("update weibo user info data failed")
	}
	return
}

// 查
func (m *WeiboUserModel) GetWeiboUserInfo(uid int64) (uinfo WeiboUserInfo, err error) {
	if uid <= 0 {
		return
	}
	cond := make(map[string]interface{})
	cond["wb_uid"] = uid
	row, err := m.FetchRow(cond)
	if err != nil || len(row) <= 0 {
		return
	}
	uinfo = formatWeiboUserInfo(row)
	return
}

// 提取文章列表，按时间倒序排序
func (m *WeiboUserModel) GetWeiboUserList(cond map[string]interface{}, page, pagesize int) ([]WeiboUserInfo, error) {
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `wb_uid` DESC")
	var ret []WeiboUserInfo
	for _, row := range rows {
		ret = append(ret, formatWeiboUserInfo(row))
	}
	return ret, nil
}

// 提取文章总数
func (m *WeiboUserModel) GetWeiboUserTotal(cond map[string]interface{}) int64 {
	return m.FetchTotal(cond)
}

// 格式化文章信息
func formatWeiboUserInfo(row map[string]negoutils.ElemType) (ret WeiboUserInfo) {
	ret.WbUid, _ = row["wb_uid"].ToInt64()
	ret.ScreenName = row["screen_name"].ToString()
	ret.Name = row["name"].ToString()
	ret.Desc = row["description"].ToString()
	ret.ProfileImageUrl = row["profile_image_url"].ToString()
	ret.ProfileUrl = row["profile_url"].ToString()
	ret.RawJson = row["raw_json"].ToString()
	ret.AccessToken = row["access_token"].ToString()
	ret.AccessTokenExpireTime, _ = row["access_token_expire_time"].ToInt64()
	return
}
