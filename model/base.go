/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     model
 * @date        2018-02-03 16:00
 */
package model

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/negoutils"
)

// 通过开放平台登录的cookie信息
type BlogCookieInfo struct {
	RuntofuUid int64  `json:"u"` //本地的用户ID
	CookieVal  string `json:"c"` //设置的cookie信息字段值
	Expire     int64  `json:"e"` //过期的时间戳，秒
}

type BaseModel struct {
	Table string //当前所用的表名
}

// 获取单行信息
func (m *BaseModel) FetchRow(cond map[string]interface{}) (map[string]negoutils.ElemType, error) {
	ret := make(map[string]negoutils.ElemType)
	rows, err := mDB.FetchCondRows(m.Table, cond)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	if len(rows) <= 0 {
		return ret, nil
	}
	ret = rows[0]
	return ret, nil
}

func (m *BaseModel) FetchList(cond map[string]interface{}, page int, pagesize int, orderby string) []map[string]negoutils.ElemType {
	if page <= 0 {
		page = 1
	}
	if pagesize < 1 {
		pagesize = 10
	}
	start := (page - 1) * pagesize
	where := ""
	cd, param := negoutils.FormatCond(cond, "AND")
	if len(param) > 0 {
		where = fmt.Sprintf("WHERE %s", cd)
	}
	fsql := fmt.Sprintf("SELECT * FROM %s %s %s LIMIT %d,%d", m.Table, where, orderby, start, pagesize)
	rows, err := mDB.FetchRows(fsql, param...)
	if err != nil {
		fmt.Println(err)
	}
	return rows
}

// 获取总数
func (m *BaseModel) FetchTotal(cond map[string]interface{}) int64 {
	where := ""
	cd, param := negoutils.FormatCond(cond, "AND")
	if len(param) > 0 {
		where = fmt.Sprintf("WHERE %s", cd)
	}
	fsql := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", m.Table, where)
	one, err := mDB.FetchOne(fsql, param...)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	ret, _ := one.ToInt64()
	return ret
}
