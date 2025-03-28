// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        goweb_router.go
// @date        2018-02-02 17:51

package negoutils

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

// 返回一个路由列表信息
func NewRuntofuRouterList() *RuntofuRouterList {
	ret := &RuntofuRouterList{
		RCache: make(map[string]RuntofuRouterCache),
		MFunc:  make(map[int]RouterMatchFunc),
		Mutex:  new(sync.RWMutex),
	}
	ret.AddMatchFunc(RouterTypePathInfo, matchPathInfoRouter).
		AddMatchFunc(RouterTypeRegexp, matchRegexpRouter)
	return ret
}

// 一堆路由列表，带缓存
type RuntofuRouterList struct {
	RList  []*RuntofuRouterItem          //所有的路由列表信息
	RCache map[string]RuntofuRouterCache //已经匹配过的缓存起来
	MFunc  map[int]RouterMatchFunc       //各种类型的处理函数
	Mutex  *sync.RWMutex
}

// 添加处理函数
func (rs *RuntofuRouterList) AddMatchFunc(t int, fn RouterMatchFunc) *RuntofuRouterList {
	rs.MFunc[t] = fn
	return rs
}

// 添加路由信息
func (rs *RuntofuRouterList) AddRouter(r *RuntofuRouterItem) *RuntofuRouterList {
	if r == nil {
		return rs
	}
	reflectVal := reflect.ValueOf(r.Controller)
	r.ControllerType = reflect.Indirect(reflectVal).Type()
	rs.RList = append(rs.RList, r)
	return rs
}

// 批量添加路由信息
func (rs *RuntofuRouterList) AddRouters(r ...*RuntofuRouterItem) *RuntofuRouterList {
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

// 开始匹配路由
func (rs *RuntofuRouterList) Match(ctx *RuntofuContext, req *http.Request) *RuntofuRouterItem {
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
	var router *RuntofuRouterItem = nil
	for _, rinfo := range rs.RList {
		if fn, ok := rs.MFunc[rinfo.Type]; ok {
			router = fn(ctx, path, rinfo)
			if router != nil {
				rs.RCache[path] = RuntofuRouterCache{F: fn, R: router}
				return router
			}
		}
	}
	return nil
}
