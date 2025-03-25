/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     model
 * @date        2018-02-05 12:20
 */
package model

import "github.com/liuyongshuai/negoutils"

// 实例化一个m层
func NewArticleTagModel() *ArticleTagModel {
	ret := &ArticleTagModel{}
	ret.Table = "article_tag"
	return ret
}

type ArticleTagModel struct {
	BaseModel
}

// 更新话题下面的文章数量
func (m *ArticleTagModel) UpdateTagContentNum(articleId int64, tagId int) {
	var tagIds []int
	if tagId > 0 {
		tagIds = append(tagIds, tagId)
	}
	if articleId > 0 {
		ainfo, err := MArticle.GetArticleInfo(articleId)
		if err == nil && ainfo.ArticleId > 0 {
			for _, tinfo := range ainfo.TagList {
				tagIds = append(tagIds, tinfo.TagId)
			}
		}
	}
	for _, tid := range tagIds {
		total := m.GetArticleTotal(tid)
		MTag.UpdateTagContentNum(tid, int(total))
	}
}

// 提取文章列表
func (m *ArticleTagModel) GetArticleList(tid, page, pagesize int) ([]int64, error) {
	cond := make(map[string]interface{})
	cond["tag_id"] = tid
	cond["is_publish"] = 1
	rows := m.FetchList(cond, page, pagesize, "ORDER BY `sort` DESC")
	var ret []int64
	for _, row := range rows {
		aid, _ := row["article_id"].ToInt64()
		ret = append(ret, aid)
	}
	return ret, nil
}

// 提取文章总数
func (m *ArticleTagModel) GetArticleTotal(tid ...interface{}) int64 {
	cond := make(map[string]interface{})
	cond["is_publish"] = 1
	if len(tid) > 0 {
		td, _ := negoutils.MakeElemType(tid[0]).ToInt()
		if td > 0 {
			cond["tag_id"] = td
		}
	}
	return m.FetchTotal(cond)
}
