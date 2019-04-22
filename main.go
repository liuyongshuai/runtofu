package main

import (
	"flag"
	"fmt"
	"github.com/kr/pretty"
	"github.com/liuyongshuai/runtofuAdmin/config"
	"github.com/liuyongshuai/runtofuAdmin/controller"
	"github.com/liuyongshuai/runtofuAdmin/model"
	"github.com/liuyongshuai/runtofuAdmin/routers"
	"github.com/liuyongshuai/runtofuAdmin/utils"
	"github.com/liuyongshuai/wego"
	"github.com/liuyongshuai/wego/context"
	"os"
)

func main() {
	// 命令行参数。
	var configPath string
	// para1: file handle para2: CLI para name para3: default value para4: desc info
	flag.StringVar(&configPath, "config", "conf/service.conf", "server config.")
	flag.Parse()

	// 解析配置。
	if err := config.GetConfiger().Init(configPath); err != nil {
		fmt.Printf("fail to read config.||err=%v||config=%v", err, configPath)
		os.Exit(1)
		return
	}
	conf := config.GetConfiger()
	fmt.Printf("%# v\n", pretty.Formatter(conf))

	//初始化model层
	err := model.Init(conf)
	if err != nil {
		panic(err)
	}

	//阻塞主进程
	ch := make(chan interface{})

	//admin管理系统的一些设置项
	adminApp := wego.NewWeGoAPP(). //新建一个app
					SetPort(conf.Http.Port).                                     //监听端口
					SetTplDir(conf.Http.TplDir).                                 //模板根目录
					SetTplExt(conf.Http.TplExt).                                 //模板扩展名称
					AddTplCommonData("SITE_NAME", conf.Http.SiteName).           //站点的名称
					AddTplCommonData("STATIC_PREFIX", conf.Common.StaticPrefix). //静态资源前缀
					AddTplFuncMap(utils.TplFuncs).                               //自定义的模板函数
					SetErrController(&controller.ErrorController{}).             //错误页面
					AddRouters(routers.AdminRouterList...).                      //路由信息
					SetRecoverFunc(func(ctx *context.WeGoContext) {              //panic时的处理函数
			if err := recover(); err != nil {
				errmsg := "url=" + ctx.Input.URI()
				fmt.Println(err, errmsg)
			}
		})
	go adminApp.Run()

	//阻塞之
	<-ch
}
