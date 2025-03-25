/**
 * @author      Liu Yongshuai
 * @package     model
 * @date        2018-02-14 17:05
 */
package model

import (
	"fmt"
	"github.com/liuyongshuai/goUtils"
)

//管理后台的基本菜单信息
type AdminMenuInfo struct {
	MenuId       int    `json:"menu_id"`        //菜单ID
	MenuName     string `json:"menu_name"`      //菜单名称
	MenuPath     string `json:"menu_path"`      //菜单路径
	IconName     string `json:"icon_name"`      //图标名称
	IconColor    string `json:"icon_color"`     //图标的颜色
	ParentMenuId int    `json:"parent_menu_id"` //父菜单ID
	ChildMenuNum int    `json:"child_menu_num"` //子菜单数量
	MenuSort     int    `json:"menu_sort"`      //菜单在同级中的排序
}

//管理后台的菜单层级关系
type AdminMenuList struct {
	MenuInfo    AdminMenuInfo   `json:"menuInfo"`    //当前的菜单信息
	SubMenuList []AdminMenuInfo `json:"subMenuList"` //所有的子菜单列表
}

//实例化一个m层
func NewAdminMenuModel() *AdminMenuModel {
	ret := &AdminMenuModel{}
	ret.Table = "admin_menu"
	return ret
}

type AdminMenuModel struct {
	BaseModel
}

//添加一条菜单信息
func (m *AdminMenuModel) AddMenuInfo(mInfo AdminMenuInfo) (int, bool, error) {
	if len(mInfo.MenuName) <= 0 {
		return 0, false, fmt.Errorf("invalid menu name")
	}
	if mInfo.ParentMenuId > 0 {
		pMenuInfo, err := m.GetAdminMenuInfo(mInfo.ParentMenuId)
		if err != nil {
			return 0, false, fmt.Errorf("父菜单不存在")
		}
		if pMenuInfo.ParentMenuId != 0 {
			return 0, false, fmt.Errorf("父菜单必须为一级菜单（目前只支持二级菜单深度）")
		}
	}
	data := make(map[string]interface{})
	tmp, _ := mDB.FetchOne("SELECT MAX(`menu_id`) FROM `" + m.Table + "`")
	maxMenuId, _ := tmp.ToInt()
	if maxMenuId > 0 {
		data["menu_id"] = maxMenuId + 1
	}
	data["menu_name"] = mInfo.MenuName
	data["menu_path"] = mInfo.MenuPath
	data["icon_name"] = mInfo.IconName
	data["icon_color"] = mInfo.IconColor
	data["parent_menu_id"] = mInfo.ParentMenuId
	data["menu_sort"] = mInfo.MenuSort
	insertId, b, e := mDB.InsertData(m.Table, data, false)
	go m.UpdateMenuChildNum(mInfo.ParentMenuId)
	return int(insertId), b, e
}

//更新菜单信息
func (m *AdminMenuModel) UpdateMenuInfo(menuId int, data map[string]interface{}) (bool, error) {
	cond := make(map[string]interface{})
	cond["menu_id"] = menuId
	_, b, e := mDB.UpdateData(m.Table, data, cond)
	return b, e
}

//删除单个菜单
func (m *AdminMenuModel) DeleteMenuInfo(menuId int) (bool, error) {
	mInfo, err := m.GetAdminMenuInfo(menuId)
	if err != nil || mInfo.MenuId <= 0 {
		return false, fmt.Errorf("删除菜单失败：获取当前菜单信息失败")
	}
	if mInfo.ChildMenuNum > 0 {
		return false, fmt.Errorf("删除菜单失败：当前菜单下面还有子菜单")
	}
	cond := make(map[string]interface{})
	cond["menu_id"] = menuId
	_, b, e := mDB.DeleteData(m.Table, cond)
	go m.UpdateMenuChildNum(mInfo.ParentMenuId)
	return b, e
}

//获取菜单信息
func (m *AdminMenuModel) GetAdminMenuInfo(menuId int) (ret AdminMenuInfo, err error) {
	cond := make(map[string]interface{})
	cond["menu_id"] = menuId
	row, err := m.FetchRow(cond)
	if err != nil {
		return ret, err
	}
	ret = formatAdminMenuInfo(row)
	return ret, nil
}

//获取菜单列表
func (m *AdminMenuModel) GetAdminMenuList(cond map[string]interface{}, page int, pagesize int) (ret []AdminMenuInfo) {
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `menu_sort` DESC,`menu_id` ASC")
	for _, row := range rows {
		ret = append(ret, formatAdminMenuInfo(row))
	}
	return ret
}

//提取总数
func (m *AdminMenuModel) GetAdminMenuTotal(cond map[string]interface{}) int64 {
	return m.FetchTotal(cond)
}

//更新子菜单数量
func (m *AdminMenuModel) UpdateMenuChildNum(menuId int) {
	if menuId <= 0 {
		return
	}
	mInfo, _ := m.GetAdminMenuInfo(menuId)
	if mInfo.ParentMenuId > 0 {
		return
	}
	cond := make(map[string]interface{})
	cond["parent_menu_id"] = menuId
	total := m.GetAdminMenuTotal(cond)
	data := make(map[string]interface{})
	data["child_menu_num"] = total
	m.UpdateMenuInfo(menuId, data)
}

//获取所有的菜单列表，供左侧菜单直接用的
func (m *AdminMenuModel) GetAllAdminMenuList() (ret []AdminMenuList) {
	//先查一级菜单（可直接缓存起来）
	cond := make(map[string]interface{})
	cond["parent_menu_id"] = 0
	topTotal := m.GetAdminMenuTotal(cond)
	topList := m.GetAdminMenuList(cond, 1, int(topTotal))
	for _, mInfo := range topList {
		amList := AdminMenuList{MenuInfo: mInfo}
		if mInfo.ChildMenuNum > 0 {
			cond := make(map[string]interface{})
			cond["parent_menu_id"] = mInfo.MenuId
			t := m.GetAdminMenuTotal(cond)
			amList.SubMenuList = m.GetAdminMenuList(cond, 1, int(t))
		}
		ret = append(ret, amList)
	}
	ret = append(ret, mSystemMenuList)
	return ret
}

//格式化菜单信息
func formatAdminMenuInfo(row map[string]goUtils.ElemType) (ret AdminMenuInfo) {
	if len(row) <= 0 {
		return
	}
	ret.MenuId, _ = row["menu_id"].ToInt()
	ret.MenuName = row["menu_name"].ToString()
	ret.MenuPath = row["menu_path"].ToString()
	ret.IconName = row["icon_name"].ToString()
	ret.IconColor = row["icon_color"].ToString()
	ret.ParentMenuId, _ = row["parent_menu_id"].ToInt()
	ret.ChildMenuNum, _ = row["child_menu_num"].ToInt()
	ret.MenuSort, _ = row["menu_sort"].ToInt()
	return ret
}
