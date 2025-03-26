// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @file        main.go
// @date        2025-03-25 16:12

package main

import (
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"github.com/liuyongshuai/negoutils"
	"github.com/liuyongshuai/runtofu/confutils"
	"github.com/liuyongshuai/runtofu/controller/admin"
	"github.com/liuyongshuai/runtofu/controller/blog"
	"github.com/liuyongshuai/runtofu/model"
	"github.com/liuyongshuai/runtofu/routers"
	"os"
	"runtime"
)

func main() {
	// 命令行参数。
	var configPath string
	// para1: file handle para2: CLI para name para3: default value para4: desc info
	flag.StringVar(&configPath, "config", "conf/service.conf", "server config.")
	flag.Parse()

	// 解析配置。
	if err := confutils.GetConfiger().Init(configPath); err != nil {
		fmt.Printf("fail to read config.||err=%v||config=%v", err, configPath)
		os.Exit(1)
		return
	}
	conf := confutils.GetConfiger()
	fmt.Printf("%# v\n", pretty.Formatter(conf))

	//初始化model层
	err := model.Init(conf)
	if err != nil {
		panic(err)
	}

	//阻塞主进程
	ch := make(chan interface{})

	//APP的一些设置项
	runtofuApp := negoutils.NewRuntofuAPP(). //新建一个app
		SetPort(conf.Http.Port). //监听端口
		SetTplDir(conf.Http.TplDir). //模板根目录
		SetTplExt(conf.Http.TplExt). //模板扩展名称
		AddTplCommonData("SITE_NAME", conf.Http.SiteName). //站点的名称
		AddTplCommonData("STATIC_PREFIX", conf.Common.StaticPrefix). //静态资源前缀
		AddTplCommonData("WEIBO_OAUTH", model.MWeiboApi.GetAuthorizeUrl()). //weibo登录的跳转地址
		AddTplCommonData("GITHUB_OAUTH", model.MGithubApi.GetAuthorizeUrl()). //github登录的跳转地址
		AddTplCommonData("runtofuUserInfo", model.RuntofuUserInfo{}). //登录后的用户信息，先赋个空值
		AddTplFuncMap(confutils.TplFuncs). //自定义的模板函数
		SetErrController(&blog.ErrorController{}). //错误页面
		AddRouters(routers.BlogRouterList...). //路由信息
		SetRecoverFunc(func(ctx *negoutils.RuntofuContext) { //panic时的处理函数
			if err := recover(); err != nil {
				errmsg := "url=" + ctx.Input.URI()
				fmt.Println(err, errmsg)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				stackStr := negoutils.ByteToStr(buf[0:stackSize])
				fmt.Println(stackStr)
			}
		})
	go runtofuApp.Run()

	//admin管理系统的一些设置项
	adminApp := negoutils.NewRuntofuAPP(). //新建一个app
		SetPort(conf.Admin.Port). //监听端口
		SetTplDir(conf.Admin.TplDir). //模板根目录
		SetTplExt(conf.Admin.TplExt). //模板扩展名称
		AddTplCommonData("SITE_NAME", conf.Http.SiteName). //站点的名称
		AddTplCommonData("STATIC_PREFIX", conf.Common.StaticPrefix). //静态资源前缀
		AddTplFuncMap(confutils.TplFuncs). //自定义的模板函数
		SetErrController(&admin.ErrorController{}). //错误页面
		AddRouters(routers.AdminRouterList...). //路由信息
		SetRecoverFunc(func(ctx *negoutils.RuntofuContext) { //panic时的处理函数
			if err := recover(); err != nil {
				errmsg := "url=" + ctx.Input.URI()
				fmt.Println(err, errmsg)
				buf := make([]byte, 1<<16)
				stackSize := runtime.Stack(buf, false)
				stackStr := negoutils.ByteToStr(buf[0:stackSize])
				fmt.Println(stackStr)
			}
		})
	go adminApp.Run()

	//阻塞之
	<-ch
}
