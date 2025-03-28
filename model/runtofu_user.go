/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-02-08 13:31
 */
package model

import (
	"errors"
	"fmt"
	"github.com/liuyongshuai/runtofu/negoutils"
)

type RuntofuUserInfo struct {
	Uid        int64  `json:"uid"`         //用户ID
	Name       string `json:"name"`        //名称
	Portrait   string `json:"portrait"`    //头像
	ProfileUrl string `json:"profile_url"` //个人主页地址
	Type       int64  `json:"type"`        //用户类型：1-微博、2-github
	ThirdUid   string `json:"third_uid"`   //第三方的用户ID
}

// 实例化一个m层
func NewRuntofuUserModel() *RuntofuUserModel {
	ret := &RuntofuUserModel{}
	ret.Table = "runtofu_user"
	return ret
}

type RuntofuUserModel struct {
	BaseModel
}

// 增
func (m *RuntofuUserModel) AddRuntofuUserInfo(uinfo RuntofuUserInfo) (err error) {
	if uinfo.Uid <= 0 {
		return errors.New("invalid Runtofu uid")
	}
	uData := make(map[string]interface{})
	uData["uid"] = uinfo.Uid
	uData["name"] = uinfo.Name
	uData["profile_url"] = uinfo.ProfileUrl
	uData["portrait"] = uinfo.Portrait
	uData["type"] = uinfo.Type
	uData["third_uid"] = uinfo.ThirdUid
	_, b, e := mDB.InsertData(m.Table, uData, false)
	if e != nil || !b {
		return fmt.Errorf("insert Runtofu uinfo data failed")
	}
	return nil
}

// 改
func (m *RuntofuUserModel) UpdateRuntofuUserInfo(uid int64, data map[string]interface{}) (err error) {
	if uid <= 0 {
		return fmt.Errorf("invalid runtofu user id")
	}
	uinfo, err := m.GetRuntofuUserInfo(uid)
	if err != nil || uinfo.Uid <= 0 {
		return fmt.Errorf("runtofu user info not exists")
	}
	cond := make(map[string]interface{})
	cond["uid"] = uid
	_, b, e := mDB.UpdateData(m.Table, data, cond)
	if e != nil || !b {
		return fmt.Errorf("update Runtofu user info data failed")
	}
	return
}

// 查
func (m *RuntofuUserModel) GetRuntofuUserInfo(uid int64) (uinfo RuntofuUserInfo, err error) {
	if uid <= 0 {
		return
	}
	cond := make(map[string]interface{})
	cond["uid"] = uid
	row, err := m.FetchRow(cond)
	if err != nil || len(row) <= 0 {
		return
	}
	uinfo = formatRuntofuUserInfo(row)
	return
}

// 提取文章列表，按时间倒序排序
func (m *RuntofuUserModel) GetRuntofuUserList(cond map[string]interface{}, page, pagesize int) ([]RuntofuUserInfo, error) {
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `uid` DESC")
	var ret []RuntofuUserInfo
	for _, row := range rows {
		ret = append(ret, formatRuntofuUserInfo(row))
	}
	return ret, nil
}

// 提取文章总数
func (m *RuntofuUserModel) GetRuntofuUserTotal(cond map[string]interface{}) int64 {
	return m.FetchTotal(cond)
}

// 格式化文章信息
func formatRuntofuUserInfo(row map[string]negoutils.ElemType) (ret RuntofuUserInfo) {
	ret.Uid, _ = row["uid"].ToInt64()
	ret.Name = row["name"].ToString()
	ret.Portrait = row["portrait"].ToString()
	ret.Type, _ = row["type"].ToInt64()
	ret.ThirdUid = row["third_uid"].ToString()
	ret.ProfileUrl = row["profile_url"].ToString()
	return
}
