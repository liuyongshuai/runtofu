// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        goweb_router_item.go
// @date        2018-02-02 17:51

package negoutils

import "reflect"

// 路由类型
const (
	InvalidRouterType  = iota
	RouterTypePathInfo //全路径匹配
	RouterTypeRegexp   //正则匹配
)

// 单个路由结构体
type RuntofuRouterItem struct {
	Type           int          //路由类型
	Config         string       //相关的配置
	Controller     interface{}  //所引用的控制层
	ControllerType reflect.Type //控制层的类型
	Param          string       //额外的参数
}

// 要缓存的路由
type RuntofuRouterCache struct {
	R *RuntofuRouterItem //匹配的路由项
	F RouterMatchFunc    //所要用的处理函数
}
