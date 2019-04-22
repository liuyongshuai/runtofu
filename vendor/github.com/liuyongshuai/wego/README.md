# tofugo
golang的超轻型的Mini型Web框架，特点如下：

* 支持简单的请求/响应时的上下文
* 支持路径式的路由和正则式的路由两种方式
* 提供了简单的Controller层的封装，以应付常用的http请求
* 提供了对模板的简单支持，每个应用对应一个模板根目录

示例见：https://github.com/liuyongshuai/runtofu

目录结构如下：
```
├── README.md
├── app.go //为APP的创建提供了一些基本的方法
├── context    //对所有的请求提供了一个上下文环境
│   ├── context.go
│   ├── input.go
│   └── output.go
├── controller //定义了controller的基类
│   └── controller.go
├── router  //请求的路由
│   ├── match_func.go
│   ├── router_item.go
│   └── routers.go
├── tofugo.go //这是一个 ServeHTTP Handler的封装，提供了路由转发、初始化上下文、传递参数等操作
└── tplparser //解析tpl模板信息
    ├── tpl.go
    └── tplfunc.go
```
* controller：此基类提供了参数获取、请求类型判断、赋数据给模板、响应信息等。
* router：提供了较为简单的路由转发功能，一种是按预定好的path，另一种是正则匹配的方式。但提供了添加自己的路由函数的接口方法，可以自定义路由匹配。
* tofugo.go：这里提供了操作controller、router、context的所有封装，它就是处理请求的handler。
* app.go：就比较简单了，就是把handler的接口对外暴露一下。


