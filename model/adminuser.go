/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-02-11 15:41
 */
package model

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/negoutils"
)

const (
	ADMIN_USER_TYPE_SUPER  = 100 //超级用户
	ADMIN_USER_TYPE_NORMAL = 1   //普通用户
)

type AdminUserInfo struct {
	Uid       int64  //用户ID
	Type      int    //用户类型：100超级管理员、1普通用户
	LoginName string //登录名称
	RealName  string //真实姓名
	Passwd    string //密码
	Passcode  string //生成密码的code
}

// 实例化一个m层
func NewAdminUserModel() *AdminUserModel {
	ret := &AdminUserModel{}
	ret.Table = "admin_user"
	return ret
}

// 管理后台的cookie信息
type AdminCookieInfo struct {
	Uid       int64  `json:"uid"`        //用户ID
	CookieVal string `json:"cookie_val"` //设置的cookie信息字段值
	Expire    int64  `json:"expire"`     //过期的时间戳，nano
}

type AdminUserModel struct {
	BaseModel
}

// 根据登录名称获取用户信息
func (m *AdminUserModel) GetAdminUserInfoByLoginName(loginName string) (ret AdminUserInfo) {
	cond := make(map[string]interface{})
	cond["login_name"] = loginName
	rows := m.GetAdminUserList(cond, 1, 1)
	if len(rows) > 0 {
		ret = rows[0]
		return
	}
	return ret
}

// 提取用户信息
func (m *AdminUserModel) GetAdminUserInfoByUid(uid int64) (ret AdminUserInfo) {
	cond := make(map[string]interface{})
	cond["uid"] = uid
	rows := m.GetAdminUserList(cond, 1, 1)
	if len(rows) > 0 {
		ret = rows[0]
		return
	}
	return ret
}

// 提取用户信息
func (m *AdminUserModel) GetAdminUserList(cond map[string]interface{}, page int, pagesize int) (ret []AdminUserInfo) {
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `uid` ASC")
	for _, row := range rows {
		ret = append(ret, formatAdminUserInfo(row))
	}
	return ret
}

// 获取用户总数
func (m *AdminUserModel) GetAdminUserTotal(cond map[string]interface{}) int64 {
	return m.FetchTotal(cond)
}

// 更新用户信息
func (m *AdminUserModel) UpdateAdminUserInfo(uid int64, data map[string]interface{}) (bool, error) {
	if uid <= 0 {
		return false, fmt.Errorf("invalid uid")
	}
	cond := make(map[string]interface{})
	cond["uid"] = uid
	_, b, e := mDB.UpdateData(m.Table, data, cond)
	return b, e
}

// 添加管理后台用户
func (m *AdminUserModel) AddAdminUserInfo(lname, rname, passwd string) (int64, error) {
	passcode := negoutils.RandomStr(16)
	isql := fmt.Sprintf("SELECT MAX(`uid`) FROM `%s`", m.Table)
	tmp, _ := mDB.FetchOne(isql)
	maxUid, _ := tmp.ToInt()
	passwd = negoutils.MD5(passwd + passcode)
	data := make(map[string]interface{})
	if maxUid <= 0 {
		maxUid = 10000
	}
	var uid int64 = 0
	var b = false
	var e error
	for i := 0; i < 3; i++ {
		maxUid++
		data["uid"] = maxUid
		data["login_name"] = lname
		data["real_name"] = rname
		data["passwd"] = passwd
		data["passcode"] = passcode
		uid, b, e = mDB.InsertData(m.Table, data, false)
		if e != nil || !b {
			continue
		}
		return uid, nil
	}
	return 0, fmt.Errorf("insert admin user info failed")
}

// 修改密码
func (m *AdminUserModel) ChangePasswd(uid int64, oldPasswd, newPasswd string) (bool, error) {
	uinfo := m.GetAdminUserInfoByUid(uid)
	if uinfo.Uid <= 0 {
		return false, fmt.Errorf("用户信息不存在")
	}
	if oldPasswd == newPasswd {
		return false, fmt.Errorf("新密码不能跟原密码相同")
	}
	dbPasswd := negoutils.MD5(oldPasswd + uinfo.Passcode)
	if dbPasswd != uinfo.Passwd {
		return false, fmt.Errorf("原密码校验失败")
	}
	data := make(map[string]interface{})
	newPasscode := negoutils.RandomStr(16)
	data["passwd"] = negoutils.MD5(newPasswd + newPasscode)
	data["passcode"] = newPasscode
	return m.UpdateAdminUserInfo(uid, data)
}

// 格式化用户信息
func formatAdminUserInfo(row map[string]negoutils.ElemType) (ret AdminUserInfo) {
	if len(row) <= 0 {
		return
	}
	ret.Uid, _ = row["uid"].ToInt64()
	ret.LoginName = row["login_name"].ToString()
	ret.RealName = row["real_name"].ToString()
	ret.Passwd = row["passwd"].ToString()
	ret.Passcode = row["passcode"].ToString()
	ret.Type, _ = row["type"].ToInt()
	return ret
}
