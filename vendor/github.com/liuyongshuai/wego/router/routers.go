/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     router
 * @date        2018-02-02 17:39
 */
package router

import (
	"github.com/liuyongshuai/wego/context"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

//返回一个路由列表信息
func NewWeGoRouterList() *WeGoRouterList {
	ret := &WeGoRouterList{
		RCache: make(map[string]WeGoRouterCache),
		MFunc:  make(map[int]RouterMatchFunc),
		Mutex:  new(sync.RWMutex),
	}
	ret.AddMatchFunc(RouterTypePathInfo, matchPathInfoRouter).
		AddMatchFunc(RouterTypeRegexp, matchRegexpRouter)
	return ret
}

//一堆路由列表，带缓存
type WeGoRouterList struct {
	RList  []*WeGoRouterItem          //所有的路由列表信息
	RCache map[string]WeGoRouterCache //已经匹配过的缓存起来
	MFunc  map[int]RouterMatchFunc    //各种类型的处理函数
	Mutex  *sync.RWMutex
}

//添加处理函数
func (rs *WeGoRouterList) AddMatchFunc(t int, fn RouterMatchFunc) *WeGoRouterList {
	rs.MFunc[t] = fn
	return rs
}

//添加路由信息
func (rs *WeGoRouterList) AddRouter(r *WeGoRouterItem) *WeGoRouterList {
	if r == nil {
		return rs
	}
	reflectVal := reflect.ValueOf(r.Controller)
	r.ControllerType = reflect.Indirect(reflectVal).Type()
	rs.RList = append(rs.RList, r)
	return rs
}

//批量添加路由信息
func (rs *WeGoRouterList) AddRouters(r ...*WeGoRouterItem) *WeGoRouterList {
	if len(r) <= 0 {
		return rs
	}
	for _, i := range r {
		if i != nil {
			reflectVal := reflect.ValueOf(i.Controller)
			i.ControllerType = reflect.Indirect(reflectVal).Type()
			rs.RList = append(rs.RList, i)
		}
	}
	return rs
}

//开始匹配路由
func (rs *WeGoRouterList) Match(ctx *context.WeGoContext, req *http.Request) *WeGoRouterItem {
	if len(rs.RList) == 0 {
		return nil
	}
	//提取请求的URI
	path := req.URL.Path
	path = strings.Trim(path, "/")

	//解析请求的URI，提取参数，放到上下文环境里
	urlInfo, _ := url.Parse(req.RequestURI)
	vals := urlInfo.Query()
	for k, vs := range vals {
		if len(vs) > 0 {
			ctx.Input.SetParam(k, vs[0])
		} else {
			ctx.Input.SetParam(k, "")
		}
	}

	//先提取缓存里有没有
	if rc, ok := rs.RCache[path]; ok {
		router := rc.R
		fn := rc.F
		rter := fn(ctx, path, router)
		if rter != nil {
			return rter
		}
	}

	//开始匹配路由信息
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()
	var router *WeGoRouterItem = nil
	for _, rinfo := range rs.RList {
		if fn, ok := rs.MFunc[rinfo.Type]; ok {
			router = fn(ctx, path, rinfo)
			if router != nil {
				rs.RCache[path] = WeGoRouterCache{F: fn, R: router}
				return router
			}
		}
	}
	return nil
}
