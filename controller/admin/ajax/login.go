/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-11 22:57
 */
package ajax

import (
	"encoding/json"
	"fmt"
	"github.com/liuyongshuai/negoutils"
	"github.com/liuyongshuai/runtofu/model"
	"strings"
	"time"
)

type AdminAjaxLoginController struct {
	AdminAjaxBaseController
}

// 校验是否为Ajax请求
func (bc *AdminAjaxLoginController) Prepare() error {
	return nil
}

// 返回数据信息
func (bc *AdminAjaxLoginController) Run() {
	name := bc.GetParam("name", "").ToString()
	passwd := bc.GetParam("passwd", "").ToString()

	//cookie的默认时长为10天
	var cookieExpireVal int64 = 864000

	//开始登录操作
	if len(name) > 0 && len(passwd) > 0 {
		uinfo := model.MAdminUser.GetAdminUserInfoByLoginName(name)
		if uinfo.Uid <= 0 {
			bc.Notice(nil, 100100, "登录失败：获取用户信息失败")
			return
		}
		p := negoutils.MD5(passwd + uinfo.Passcode)
		if p != uinfo.Passwd {
			bc.Notice(nil, 100100, "登录失败：校验密码失败")
			return
		}
		//前面16位随机数 + uid的base62 + 前面两项的md5值。别问为啥这么做，就是特么的感觉牛X些
		cookieVal := negoutils.RandomStr(16) + negoutils.Base62Encode(uinfo.Uid)
		md5 := strings.ToUpper(negoutils.MD5(cookieVal))
		cookieVal += md5
		jsUinfo, err := json.Marshal(model.AdminCookieInfo{
			Uid:       uinfo.Uid,
			CookieVal: cookieVal,
			Expire:    time.Now().Unix() + cookieExpireVal,
		})
		if err != nil {
			fmt.Println("json.Marshal", err)
			bc.Notice(nil, 100100, "登录失败："+err.Error())
			return
		}
		//将相关信息存入redis
		rkey := fmt.Sprintf(model.ADMIN_SESSION_PREFIX, uinfo.Uid)
		rconn := model.RedisPool.Get()
		defer rconn.Close()
		_, e := rconn.Do("setex", rkey, cookieExpireVal, string(jsUinfo))
		if e != nil {
			fmt.Println("redis.Do setex", e, string(jsUinfo))
			bc.Notice(nil, 100100, "登录失败："+e.Error())
			return
		}
		bc.Ctx.Output.AddCookie(model.CookieKey, cookieVal, cookieExpireVal, "/")
	}
	bc.Notice(nil)
}
