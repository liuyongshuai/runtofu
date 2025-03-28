{{define "header"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <title>{{.SITE_NAME}}</title>
    <link rel="Shortcut Icon" type="image/x-icon" href="{{ static_image "favicon.ico"}}" />
    <link rel="Bookmark" href="{{ static_image "favicon.ico" }}" />
    <meta name="referrer" content="origin-when-cross-origin" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" />
    <meta name="apple-mobile-web-app-title" content="{{.SITE_NAME}}" />
    <meta name="renderer" content="webkit" />
    <meta name="author" content="liuyongshuai@hotmail.com" />
    <meta name="baidu-site-verification" content="youWOPbOW2" />
    <meta name="google-site-verification" content="Qj3ww6bloz8FHkDgOEx2p2aDWyzkeZmlOGGEW1zM03M" />
    <meta property="wb:webmaster" content="a97aa556624c7a14" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <script>var STATIC_PREFIX = "{{.STATIC_PREFIX}}";</script>
    <script src="/Users/liuyongshuai/mycode/liuyongshuai/runtofu/static/js/base.js"></script>
    <link type="text/css" rel="stylesheet" href="/Users/liuyongshuai/mycode/liuyongshuai/runtofu/static/css/base.css" />
    {{ static_js  "base.js" }}
    {{ static_css  "base.css" }}
    {{ static_css  "runtofu.css" }}
    <!--[if lt IE 9]>
     {{ static_js  "html5shiv.js" }}
    <![endif]-->
</head>
<!--[if lt IE 7 ]>
<body class="ie ie6"> <![endif]-->
<!--[if IE 7 ]>
<body class="ie ie7"> <![endif]-->
<!--[if IE 8 ]>
<body class="ie ie8"> <![endif]-->
<!--[if IE 9 ]>
<body class="ie ie9"> <![endif]-->
<!--[if (gt IE 9)|!(IE)]><!-->
<body>
{{/***************顶部的导航栏***************/}}
<header class="navbar navbar-fixed-top navbar-default" role="navigation" id="runtofu-header-navbar">
    <div class="container-fluid">
        <div class="row-fluid">
            <div class="col-sm-1 col-md-1"></div>
            <div class="col-sm-10 col-md-10">
                <div class="navbar-header pull-left">
                    <a class="navbar-brand" href="/" title="返回首页">
                        <span class="glyphicon glyphicon-home" style="color:green"></span>
                    </a>
                </div>
                <div class="navbar-header pull-left">
                    <a class="navbar-brand b" href="javascript:void(0);">奔跑的豆腐</a>
                </div>
                {{/*导航栏*/}}
                <div class="navbar-header tofu-nav-right">
                    <a class="navbar-brand" href="/taglist">
                        话题标签
                    </a>
                    <a class="navbar-brand" href="/reclist">
                        本站推荐
                    </a>
                    <a class="navbar-brand" href="/article/1518877692688732160">
                        关于本站
                    </a>
                    <a class="navbar-brand" href="/feedback">
                        留言板
                    </a>
                    <a class="navbar-brand" href="http://admin.runtofu.com" target="_blank" style="color:green;font-weight:bold;">
                        管理后台
                    </a>
                </div>
            </div>
            {{/*登录相关的图标*/}}
            <div class="col-sm-1 col-md-1 login-wrapper">
                {{ if le .runtofuUserInfo.Uid 0 }}
                <div class="login-icon">
                    <a href="{{ .WEIBO_OAUTH }}">
                        <img src="{{ static_image "logo/weibo.png"}}" />
                    </a>
                </div>
                <div class="login-icon">
                    <a href="{{ .GITHUB_OAUTH }}">
                       <svg aria-hidden="true" class="" height="25" version="1.1" viewBox="0 0 16 16" width="25"><path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"></path></svg>
                    </a>
                </div>
                {{else}}
                <div class="login-icon">
                    <a href="{{ .runtofuUserInfo.ProfileUrl }}" title="{{ .runtofuUserInfo.Name }}" target="_blank" data-toggle="tooltip" data-placement="bottom">
                        <img src="{{ .runtofuUserInfo.Portrait }}" />
                    </a>
                </div>
                {{end}}
            </div>
        </div>
    </div>
</header>

{{/************主内容**********/}}
<div class="main-content" id="global-main-content">
    <div class="container-fluid">
        <div class="row-fluid">
        <div class="col-sm-1 col-md-1"></div>
        <div id="tofu-content" class="col-sm-10 col-md-10 tofu-content">
{{end}}