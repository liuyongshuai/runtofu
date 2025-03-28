package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/liuyongshuai/runtofu/model"
	"github.com/liuyongshuai/runtofu/negoutils"
	"html/template"
	"strings"
	"time"
)

var gPageSize int64 = 20

type RunToFuBaseController struct {
	negoutils.RuntofuController
}

// 获取已登录的用户信息
func (c *RunToFuBaseController) Prepare() error {
	runtofuUInfo := c.GetLoginUserInfo()
	c.AddTplData("runtofuUserInfo", runtofuUInfo)
	return nil
}

// 获取登录的用户信息，如果有的话
func (c *RunToFuBaseController) GetLoginUserInfo() (runtofuUserInfo model.RuntofuUserInfo) {
	//提取cookie里的相关值，它是有一定长度的
	cookieVal := c.Ctx.Input.GetCookie(model.BlogCookieKey)
	fmt.Println("cookieVal", cookieVal)
	cookieLen := len(cookieVal)
	if cookieLen <= 48 {
		fmt.Println("cookieLen <=48", cookieLen)
		return
	}

	//本地的用户ID信息
	var runtofuUid int64 = 0

	//前面16位随机数 + uid的base62 + 前面两项的md5值。别问为啥这么做，就是特么的感觉牛X些
	tmp := cookieVal[:cookieLen-32]
	runtofuUid = negoutils.Base62Decode(tmp[16:])
	md5 := cookieVal[cookieLen-32:]
	if md5 != strings.ToUpper(negoutils.MD5(tmp)) || runtofuUid <= 0 {
		fmt.Println("md5=", cookieLen, "\truntofuUid=", runtofuUid)
		return
	}

	//拿着取得到的uid从redis里提取session信息
	rkey := fmt.Sprintf(model.BLOG_SESSION_PREFIX, runtofuUid)
	rconn := model.RedisPool.Get()
	defer rconn.Close()
	tmpInfo, err := rconn.Do("get", rkey)
	//取失败则返回
	if err != nil || tmpInfo == nil {
		fmt.Println("redisInfo=", tmpInfo, "\terr=", err)
		return
	}

	//逐个字段的判断redis里的信息是否正确
	var cookieInfo model.BlogCookieInfo
	if tmpInfo != nil {
		err = json.Unmarshal(tmpInfo.([]byte), &cookieInfo)
		if err != nil {
			fmt.Println("json.Unmarshal failed\t", "\terr=", err)
			return
		}
	}

	//提取出来的session信息不对
	if cookieInfo.RuntofuUid <= 0 {
		fmt.Println("cookieInfo.RuntofuUid <= 0\t", "\tcookieInfo=", cookieInfo)
		return
	}

	//过期了
	if cookieInfo.Expire < time.Now().Unix() {
		fmt.Println("cookieInfo.Expire < time.Now().Unix()\t", "\tcookieInfo=", cookieInfo)
		return
	}

	//redis里的值跟cookie对不上
	if cookieInfo.CookieVal != cookieVal {
		fmt.Println("cookieInfo.CookieVal != cookieVal\t", "\tcookieInfo=", cookieInfo)
		return
	}

	//获取用户信息
	runtofuUserInfo, err = model.MRuntofuUser.GetRuntofuUserInfo(runtofuUid)
	if err != nil {
		fmt.Println("GetRuntofuUserInfo_ERROR", err)
	}
	fmt.Println("GetRuntofuUserInfo", runtofuUserInfo)
	return
}

// 专门渲染文章列表及右边的话题信息
func (c *RunToFuBaseController) articleListAndTags(articleList []model.ArticleInfo, totalNum int64) (template.HTML, error) {
	data := make(map[interface{}]interface{})
	tagList, _ := model.MTag.GetTagList(1, 50)
	data["articleList"] = articleList
	data["tagList"] = tagList
	data["pagination"] = c.getPagination(totalNum)

	buf := new(bytes.Buffer)
	err := c.Tpl.ExecuteTpl(buf, "article_list_and_tags.tpl", data)
	if err != nil {
		fmt.Println(err)
	}
	bufstr := buf.String()
	return template.HTML(bufstr), err
}

// 分页
func (c *RunToFuBaseController) getPagination(totalNum int64) template.HTML {
	reqUrl := c.Ctx.Request.URL
	rawQuery := reqUrl.RawQuery
	curPath := reqUrl.Path
	return negoutils.Pagination(curPath, rawQuery, totalNum, gPageSize, "page")
}
