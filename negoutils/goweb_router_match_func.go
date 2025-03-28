// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        goweb_router_match_func.go
// @date        2018-02-02 17:51

package negoutils

import (
	"fmt"
	"regexp"
	"strings"
)

// 路由解析函数
type RouterMatchFunc func(*RuntofuContext, string, *RuntofuRouterItem) *RuntofuRouterItem

/*
*
匹配全路径路由

	:arg	只配置数字类型
	:arg:	可以匹配任意类型

如：
/ggtest/:id	匹配下面的URL，并把 id作为参数回传给controller层

	/ggtest/6666
	/ggtest/89999

/ggtest/:name:	匹配下面的URL。并把name作为参数回传给controller层

	/ggtest/wendao
	/ggtest/abc
	/ggtest/44444
*/
func matchPathInfoRouter(ctx *RuntofuContext, uri string, rt *RuntofuRouterItem) *RuntofuRouterItem {
	//请求的URL的切片
	pathinfo := strings.Split(uri, "/")
	//事先配置好的路由的切片，要和URL逐个比对，若有:arg、:arg:这样的还要替换
	config := strings.Split(rt.Config, "/")
	if len(pathinfo) > len(config) {
		return nil
	}
	for key, val := range config {
		//当前段以冒号“:”开头
		if strings.HasPrefix(val, ":") {
			if len(pathinfo) > key {
				//非数字必须  :arg:，数字只能 :arg
				if !IsAllNumber(pathinfo[key]) && !strings.HasSuffix(val, ":") {
					return nil
				}
				//提取匹配上的参数名和参数值，并回填到Input对象里去
				pk, pv := strings.Trim(val, ":"), pathinfo[key]
				ctx.Input.SetParam(pk, pv)
			}
			//如果不是以冒号“:”开头，则请求URL的各段必须和配置的路由的各段完全相等才算匹配上
		} else if len(pathinfo) <= key || pathinfo[key] != val {
			return nil
		}
	}
	return rt
}

/*
*
匹配正则路由，直接用正则表达式去匹配请求的URL，并把捕获的参数回传到Param配置里,如
config 为 `^ggtest/aid(\w+?)/cid(\d+)$`，Param 为 aid=$1&cid=$2
则将请求中aid后面的字符串挑出来赋给aid，cid后面的字符串挑出来赋给cid
*/
func matchRegexpRouter(ctx *RuntofuContext, uri string, rt *RuntofuRouterItem) *RuntofuRouterItem {
	//先判断请求的URLj否匹配给定的正则表达式
	match, err := regexp.MatchString(rt.Config, uri)
	if err != nil || !match {
		return nil
	}
	//要替换的参数配置，如a=$1&b=$2，1、2对应正则表达式里的捕获的参数
	arg := rt.Param
	reg := regexp.MustCompile(rt.Config)
	//从请求的URL里提取出被捕获的参数部分，返回值为包含切片的切片，取第一个即为匹配的参数
	tmpArg := reg.FindAllStringSubmatch(uri, -1)
	routeArgs := []string{}
	if len(tmpArg) > 0 {
		routeArgs = tmpArg[0]
	}
	//匹配出来的数据，0项为url，其他项即为正则捕获的参数
	if len(routeArgs) > 0 {
		routeArgs = routeArgs[1:]
	}
	//开始遍历所有捕获的参数，并替换配置好的$1、$2等
	for i, val := range routeArgs {
		j := i + 1
		repstr := fmt.Sprintf("$%d", j)
		arg = strings.Replace(arg, repstr, val, -1)
	}
	//将替换好的参数切开、回填到Input对象里去
	tmp := strings.Split(arg, "&")
	for _, a := range tmp {
		if !strings.Contains(a, "=") {
			continue
		}
		t1 := strings.Split(a, "=")
		ctx.Input.SetParam(t1[0], t1[1])
	}
	return rt
}
