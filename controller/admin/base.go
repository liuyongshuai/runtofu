package admin

import (
	"encoding/json"
	"fmt"
	"github.com/liuyongshuai/goUtils"
	"github.com/liuyongshuai/runtofu/goweb/controller"
	"github.com/liuyongshuai/runtofu/model"
	"github.com/liuyongshuai/runtofu/utils"
	"html/template"
	"strings"
	"time"
)

var gPageSize int64 = 30

type AdminBaseController struct {
	controller.RuntofuController
	UserInfo model.AdminUserInfo //登录后的用户的信息
}

func (c *AdminBaseController) Prepare() error {
	c.UserInfo = c.CheckLogin(true, func() {
		c.Ctx.Redirect("/login")
	})
	c.SetLeftMenu()
	return nil
}

//校验是否登录
func (c *AdminBaseController) CheckLogin(mustLogin bool, errFn func()) (ret model.AdminUserInfo) {
	c.TplData["userInfo"] = ret

	//提取cookie里的相关值，它是有一定长度的
	cookieVal := c.Ctx.Input.GetCookie(model.CookieKey)
	if len(cookieVal) <= 48 && mustLogin {
		fmt.Println("[checkLogin]提取cookie失败")
		errFn()
		return
	}

	var uid int64 = 0
	cookieLen := len(cookieVal)
	//从cookie里截取uid
	if cookieLen > 48 {
		//前面16位随机数 + uid的base62 + 前面两项的md5值。别问为啥这么做，就是特么的感觉牛X些
		tmp := cookieVal[:cookieLen-32]
		uid = goUtils.Base62Decode(tmp[16:])
		md5 := cookieVal[cookieLen-32:]
		if md5 != strings.ToUpper(goUtils.MD5(tmp)) && mustLogin {
			uid = 0
			fmt.Println("[checkLogin]校验cookie失败，md5对不上")
			errFn()
			return
		}
	}
	if uid <= 0 && mustLogin {
		fmt.Println("[checkLogin]提取用户uid失败")
		errFn()
		return
	}

	//拿着取得到的uid从redis里提取session信息
	rkey := fmt.Sprintf(model.ADMIN_SESSION_PREFIX, uid)
	rconn := model.RedisPool.Get()
	defer rconn.Close()
	tmpInfo, err := rconn.Do("get", rkey)
	//取失败则返回
	if (err != nil || tmpInfo == nil) && mustLogin {
		fmt.Println("[checkLogin]读取redis中的session信息失败", err)
		errFn()
		return
	}

	//逐个字段的判断redis里的信息是否正确
	var craw model.AdminCookieInfo
	if tmpInfo != nil {
		err = json.Unmarshal(tmpInfo.([]byte), &craw)
		if err != nil && mustLogin {
			fmt.Println("[checkLogin]解析redis中的信息失败", err)
			errFn()
			return
		}
	}

	//提取出来的session信息不对
	if craw.Uid <= 0 && mustLogin {
		fmt.Println("[checkLogin]从redis中提取的session信息uid非法", err)
		errFn()
		return
	}

	//过期了
	if craw.Expire < time.Now().Unix() && mustLogin {
		fmt.Println("[checkLogin]登录已过期", err)
		c.Logout(craw.Uid)
		errFn()
		return
	}

	//redis里的值跟cookie对不上
	if craw.CookieVal != cookieVal && mustLogin {
		fmt.Println("[checkLogin]cookie值跟session的匹配不上", err)
		c.Logout(craw.Uid)
		errFn()
		return
	}

	//提取用户信息
	uinfo := model.MAdminUser.GetAdminUserInfoByUid(craw.Uid)
	if uinfo.Uid > 0 {
		c.TplData["userInfo"] = uinfo
		ret = uinfo
	} else if mustLogin {
		c.Logout(craw.Uid)
		fmt.Println("[checkLogin]提取用户信息失败", err)
		errFn()
	}
	return
}

//退出登录信息
func (c *AdminBaseController) Logout(uid int64) {
	//先清理cookie信息
	c.Ctx.Output.AddCookie(model.CookieKey, "", -1, "/")
	if uid <= 0 {
		return
	}
	rkey := fmt.Sprintf(model.ADMIN_SESSION_PREFIX, uid)
	rconn := model.RedisPool.Get()
	defer rconn.Close()
	rconn.Do("del", rkey)
}

//分页
func (c *AdminBaseController) getPagination(totalNum int64) template.HTML {
	reqUrl := c.Ctx.Request.URL
	rawQuery := reqUrl.RawQuery
	curPath := reqUrl.Path
	return utils.Pagination(curPath, rawQuery, totalNum, gPageSize, "page")
}

//设置左侧的菜单栏
func (c *AdminBaseController) SetLeftMenu() {
	if _, ok := c.TplData["leftMenuList"]; ok {
		return
	}
	if c.UserInfo.Uid <= 0 {
		c.TplData["leftMenuList"] = []model.AdminMenuList{}
	} else {
		c.TplData["leftMenuList"] = model.MAdminMenu.GetAllAdminMenuList()
	}
	//左边一级菜单的路径
	c.Tpl.AddTplFunc("leftMenuUrl", func(mInfo model.AdminMenuInfo) template.HTML {
		ret := fmt.Sprintf("#left_menu_%d", mInfo.MenuId)
		if mInfo.ChildMenuNum <= 0 {
			ret = mInfo.MenuPath
		}
		return template.HTML(ret)
	})
	//左则菜单的图标
	c.Tpl.AddTplFunc("leftMenuIcon", func(mInfo model.AdminMenuInfo) template.HTML {
		if len(mInfo.IconName) <= 0 {
			return template.HTML("")
		}
		style := ""
		if len(mInfo.IconColor) > 0 {
			style = fmt.Sprintf("style=\"color:%s\"", mInfo.IconColor)
		}
		ret := fmt.Sprintf("<span class=\"glyphicon %s\" %s></span>", mInfo.IconName, style)
		return template.HTML(ret)
	})
	//当前菜单是否在菜则子菜单里，要展开用
	c.Tpl.AddTplFunc("inLeftMenu", func(subList []model.AdminMenuInfo, curPath string) bool {
		for _, mInfo := range subList {
			if curPath == mInfo.MenuPath {
				return true
			}
		}
		return false
	})
}
