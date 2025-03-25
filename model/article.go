/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     model
 * @date        2018-02-03 15:59
 */
package model

import (
	"fmt"
	"github.com/liuyongshuai/negoutils"
	"strings"
	"time"
)

// 实例化一个m层
func NewArticleModel() *ArticleModel {
	ret := &ArticleModel{}
	ret.Table = "article_info"
	return ret
}

// 文章的详细信息
type ArticleInfo struct {
	ArticleId      int64     `json:"article_id"`       //文章ID
	Title          string    `json:"title"`            //文章标题
	Content        string    `json:"content"`          //文章内容
	ShortDesc      string    `json:"short_desc"`       //简单的描述信息
	TagList        []TagInfo `json:"tag_list"`         //标签列表
	IsOrigin       bool      `json:"is_origin"`        //是否为原始
	CreateTime     int       `json:"create_time"`      //发表时间
	LastModifyTime int       `json:"last_modify_time"` //最后修改时间
	IsPublish      bool      `json:"is_publish"`       //是否已发布
	IsRec          bool      `json:"is_rec"`           //是否被推荐
}

type ArticleModel struct {
	BaseModel
}

// 添加一篇文章信息
func (m *ArticleModel) AddArticleInfo(title, content string, tagIds []int, isOrigin bool) (int64, bool, error) {
	var articleId int64
	var err error
	for i := 0; i < 3; i++ {
		articleId, err = MSnowFlake.NextId()
		if err == nil {
			break
		}
	}
	if articleId <= 0 {
		return 0, false, fmt.Errorf("generate article id failed")
	}
	for _, tid := range tagIds {
		if tid <= 0 {
			return 0, false, fmt.Errorf("invalid tag_id")
		}
	}
	if len(title) <= 0 {
		return 0, false, fmt.Errorf("invalid title")
	}

	//组装要写入的文章数据
	articleData := make(map[string]interface{})
	articleData["article_id"] = articleId
	articleData["title"] = title
	articleData["content"] = content
	articleData["is_origin"] = 1
	if !isOrigin {
		articleData["is_origin"] = 0
	}
	var tags []string
	for _, tid := range tagIds {
		tags = append(tags, negoutils.MakeElemType(tid).ToString())
	}
	articleData["tags"] = strings.Join(tags, ",")
	articleData["create_time"] = time.Now().Unix()
	articleData["last_modify_time"] = time.Now().Unix()

	//要写入的文章话题数据
	tagFields := []string{"article_id", "tag_id", "is_publish"}
	var tagData [][]interface{}
	for _, tid := range tagIds {
		var tmp []interface{}
		tmp = append(tmp, articleId, tid, 0)
		tagData = append(tagData, tmp)
	}

	//开启事务，写入之
	rtx, err := mDB.BeginTransaction()
	if err != nil {
		return 0, false, err
	}
	//写入文章数据
	_, b, e := mDB.InsertData(m.Table, articleData, false)
	if e != nil || !b {
		mDB.RollBackTransaction(rtx)
		return 0, false, fmt.Errorf("insert article data failed")
	}
	//写入话题数据
	_, b, e = mDB.InsertBatchData(MArticleTag.Table, tagFields, tagData, false)
	if e != nil || !b {
		mDB.RollBackTransaction(rtx)
		return 0, false, fmt.Errorf("insert article data failed")
	}
	mDB.CommitTransaction(rtx)
	go MArticleTag.UpdateTagContentNum(articleId, 0)
	return articleId, true, nil
}

// 清理文件信息
func (m *ArticleModel) DeleteArticleInfo(articleId int64) (bool, error) {
	if articleId <= 0 {
		return false, fmt.Errorf("invalid article id")
	}
	aInfo, err := m.GetArticleInfo(articleId)
	if err != nil {
		return false, err
	}
	if aInfo.ArticleId <= 0 {
		return false, fmt.Errorf("delete article failed:not exixts")
	}
	cond := make(map[string]interface{})
	cond["article_id"] = articleId
	//开启事务，写入之
	rtx, err := mDB.BeginTransaction()
	if err != nil {
		return false, err
	}
	_, b, e := mDB.DeleteData(m.Table, cond)
	if e != nil || !b {
		mDB.RollBackTransaction(rtx)
		return false, fmt.Errorf("delete article data failed")
	}
	_, b, e = mDB.DeleteData(MArticleTag.Table, cond)
	if e != nil || !b {
		mDB.RollBackTransaction(rtx)
		return false, fmt.Errorf("delete article data failed")
	}
	mDB.CommitTransaction(rtx)

	//更新话题下面的文章数量
	go func(tlist []TagInfo) {
		for _, tInfo := range tlist {
			go MArticleTag.UpdateTagContentNum(0, tInfo.TagId)
		}
	}(aInfo.TagList)

	return true, nil
}

// 更新一篇文章信息
func (m *ArticleModel) UpdateArticleInfo(articleId int64, data map[string]interface{}) (bool, error) {
	if articleId <= 0 {
		return false, fmt.Errorf("invalid article id")
	}
	articleInfo, err := m.GetArticleInfo(articleId)
	if err != nil || articleInfo.ArticleId <= 0 {
		return false, fmt.Errorf("article info not exists")
	}
	cond := make(map[string]interface{})
	cond["article_id"] = articleId
	data["last_modify_time"] = time.Now().Unix()
	_, b, e := mDB.UpdateData(m.Table, data, cond)
	if e != nil || !b {
		return false, fmt.Errorf("update article data failed")
	}
	articleInfo, _ = m.GetArticleInfo(articleId)
	//是否修改了标签
	if _, ok := data["tags"]; ok {
		mDB.DeleteData(MArticleTag.Table, cond)

		//要写入的文章话题数据
		tagFields := []string{"article_id", "tag_id", "is_publish"}
		var tagData [][]interface{}
		for _, tinfo := range articleInfo.TagList {
			var tmp []interface{}
			tmp = append(tmp, articleId, tinfo.TagId, articleInfo.IsPublish)
			tagData = append(tagData, tmp)
		}
		_, b, e = mDB.InsertBatchData(MArticleTag.Table, tagFields, tagData, false)
		if e != nil || !b {
			return false, fmt.Errorf("insert article data failed")
		}
	}
	//是否修改了发布状态
	if _, ok := data["is_publish"]; ok {
		data := make(map[string]interface{})
		data["is_publish"] = articleInfo.IsPublish
		mDB.UpdateData(MArticleTag.Table, data, cond)
	}
	go MArticleTag.UpdateTagContentNum(articleId, 0)
	return true, nil
}

// 提取一篇文章信息
func (m *ArticleModel) GetArticleInfo(aid int64) (ArticleInfo, error) {
	ainfo := ArticleInfo{}
	if aid <= 0 {
		return ainfo, fmt.Errorf("invalid article id")
	}
	cond := make(map[string]interface{})
	cond["article_id"] = aid
	row, err := m.FetchRow(cond)
	if err != nil || len(row) <= 0 {
		return ainfo, err
	}
	ainfo = formatArticleInfo(row)
	return ainfo, nil
}

// 提取一堆文章信息
func (m *ArticleModel) GetArticleInfos(aids []int64) (map[int64]ArticleInfo, error) {
	ret := make(map[int64]ArticleInfo)
	if len(aids) <= 0 {
		return ret, fmt.Errorf("invalid article ids")
	}
	cond := make(map[string]interface{})
	cond["article_id:in"] = aids
	rows := m.FetchList(cond, 1, len(aids), "")
	for _, row := range rows {
		aid, _ := row["article_id"].ToInt64()
		ret[aid] = formatArticleInfo(row)
	}
	return ret, nil
}

// 提取文章列表，按时间倒序排序
func (m *ArticleModel) GetArticleList(cond map[string]interface{}, page, pagesize int) ([]ArticleInfo, error) {
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `article_id` DESC")
	var ret []ArticleInfo
	for _, row := range rows {
		ret = append(ret, formatArticleInfo(row))
	}
	return ret, nil
}

// 提取文章总数
func (m *ArticleModel) GetArticleTotal(cond map[string]interface{}) int64 {
	return m.FetchTotal(cond)
}

// 格式化文章信息
func formatArticleInfo(row map[string]negoutils.ElemType) ArticleInfo {
	ainfo := ArticleInfo{}
	ainfo.ArticleId, _ = row["article_id"].ToInt64()
	ainfo.Title = row["title"].ToString()
	ainfo.Content = row["content"].ToString()
	ainfo.ShortDesc = row["short_desc"].ToString()
	ts := strings.Split(row["tags"].ToString(), ",")
	for _, t := range ts {
		tid, _ := negoutils.MakeElemType(t).ToInt()
		if tid <= 0 {
			continue
		}
		tInfo, _ := MTag.GetTagInfo(tid)
		ainfo.TagList = append(ainfo.TagList, tInfo)
	}
	origin, _ := row["is_origin"].ToInt()
	ainfo.IsOrigin = false
	if origin > 0 {
		ainfo.IsOrigin = true
	}
	publish, _ := row["is_publish"].ToInt()
	ainfo.IsPublish = false
	if publish > 0 {
		ainfo.IsPublish = true
	}
	rec, _ := row["is_rec"].ToInt()
	ainfo.IsRec = false
	if rec > 0 {
		ainfo.IsRec = true
	}
	ainfo.CreateTime, _ = row["create_time"].ToInt()
	ainfo.LastModifyTime, _ = row["last_modify_time"].ToInt()
	return ainfo
}
