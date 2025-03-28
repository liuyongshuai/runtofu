/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     model
 * @date        2018-02-03 16:01
 */
package model

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/negoutils"
)

// 实例化一个m层
func NewTagModel() *TagModel {
	ret := &TagModel{}
	ret.Table = "tag"
	return ret
}

// 标签（分类）信息
type TagInfo struct {
	TagId      int    `json:"tag_id"`      //标签ID
	TagName    string `json:"tag_name"`    //标签名称
	ContentNum int    `json:"content_num"` //定时更新的内容数量
}

type TagModel struct {
	BaseModel
}

// 添加一个话题信息
func (m *TagModel) AddTagInfo(name string) (int, bool, error) {
	if len(name) <= 0 {
		return 0, false, fmt.Errorf("invalid tag name")
	}
	data := make(map[string]interface{})
	data["tag_name"] = name
	i, b, e := mDB.InsertData(m.Table, data, false)
	return int(i), b, e
}

// 删除一个话题信息
func (m *TagModel) DeleteTagInfo(tid int) (bool, error) {
	if tid <= 0 {
		return false, fmt.Errorf("invalid tag id")
	}
	total := MArticleTag.GetArticleTotal(tid)
	if total > 0 {
		return false, fmt.Errorf("there are %d article in this tag", total)
	}
	cond := make(map[string]interface{})
	cond["tag_id"] = tid
	_, r, e := mDB.DeleteData(m.Table, cond)
	return r, e
}

// 更新话题信息
func (m *TagModel) UpdateTagInfo(tid int, data map[string]interface{}) (bool, error) {
	if tid <= 0 {
		return false, fmt.Errorf("invalid tag id")
	}
	cond := make(map[string]interface{})
	cond["tag_id"] = tid
	_, _, e := mDB.UpdateData(m.Table, data, cond)
	if e != nil {
		return false, e
	}
	return true, nil
}

// 提取单个话题信息
func (m *TagModel) GetTagInfo(tid int) (TagInfo, error) {
	tinfo := TagInfo{}
	if tid <= 0 {
		return tinfo, fmt.Errorf("invalid tag id")
	}
	cond := make(map[string]interface{})
	cond["tag_id"] = tid
	row, err := m.FetchRow(cond)
	if err != nil || len(row) <= 0 {
		return tinfo, err
	}
	tinfo = formatTagInfo(row)
	return tinfo, nil
}

// 提取批量话题信息
func (m *TagModel) GetTagInfos(tids []int) (map[int]TagInfo, error) {
	ret := make(map[int]TagInfo)
	if len(tids) <= 0 {
		return ret, fmt.Errorf("invalid tag ids")
	}
	cond := make(map[string]interface{})
	cond["tag_id:in"] = tids
	rows := m.FetchList(cond, 1, len(tids), "")
	for _, row := range rows {
		tid, _ := row["tag_id"].ToInt()
		ret[tid] = formatTagInfo(row)
	}
	return ret, nil
}

// 提取批量话题信息，按内容数量排序
func (m *TagModel) GetTagList(page int, pagesize int) ([]TagInfo, error) {
	cond := make(map[string]interface{})
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `content_num` DESC")
	var ret []TagInfo
	for _, row := range rows {
		ret = append(ret, formatTagInfo(row))
	}
	return ret, nil
}

// 提取所有的话题ID
func (m *TagModel) GetAllTagIds() (ret []int) {
	cond := make(map[string]interface{})
	resp, err := mDB.FetchCondRows(m.Table, cond, "tag_id")
	if err != nil {
		return
	}
	for _, info := range resp {
		tid, _ := info["tag_id"].ToInt()
		if tid > 0 {
			ret = append(ret, tid)
		}
	}
	return
}

// 提取总数
func (m *TagModel) GetTagTotal() int64 {
	cond := make(map[string]interface{})
	return m.FetchTotal(cond)
}

// 更新话题下的内容数量
func (m *TagModel) UpdateTagContentNum(tid, cnum int) {
	data := make(map[string]interface{})
	data["content_num"] = cnum
	m.UpdateTagInfo(tid, data)
}

// 格式化话题信息
func formatTagInfo(row map[string]negoutils.ElemType) TagInfo {
	tinfo := TagInfo{}
	tinfo.TagId, _ = row["tag_id"].ToInt()
	tinfo.TagName = row["tag_name"].ToString()
	tinfo.ContentNum, _ = row["content_num"].ToInt()
	return tinfo
}
