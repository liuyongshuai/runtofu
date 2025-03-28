{{define "header"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>{{.SITE_NAME}}</title>
    <link rel="Shortcut Icon" href="{{ static_image "favicon.ico"}}" type="image/x-icon"/>
    <link rel="Bookmark" href="{{ static_image "favicon.ico"}}"/>
    <meta name="referrer" content="origin-when-cross-origin">
    <meta content="IE=edge,chrome=1" http-equiv="X-UA-Compatible">
    <meta name="author" content="Liu Yongshuai">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script>
        var STATIC_PREFIX = "{{.STATIC_PREFIX}}";
    </script>
    {{ static_js  "base.js" }}
    {{ static_css  "base.css" }}
    {{ static_css  "admin.css" }}
    {{ static_js  "admin_base.js" }}
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

{{/*顶部导航栏*/}}
<div class="navbar" role="navigation">
    <div class="navbar-inner">
        <div class="container-fluid">
            {{if gt .userInfo.Uid 0}}
                <ul class="nav pull-right">
                    <li class="dropdown">
                        <a href="javascript:void(0)" role="button" class="dropdown-toggle" data-toggle="dropdown">
                            <span class="glyphicon glyphicon-user"></span>
                            {{.userInfo.RealName}}
                            <span class="caret"></span>
                        </a>
                       <ul class="dropdown-menu">
                            <li>
                                <a tabindex="-1" href="javascript:void(0);" id="adminNavigationLogout">
                                    <span class="glyphicon glyphicon-log-out"></span>
                                    退出登录
                                </a>
                            </li>
                            <li>
                                <a tabindex="-1" href="javascript:void(0);" id="adminNavigationChangePasswd">
                                    <span class="glyphicon glyphicon-cog"></span>
                                    修改密码
                                </a>
                            </li>
                        </ul>
                    </li>
                </ul>
            {{end}}
            <a class="brand" href="/">
                <span class="second">奔跑吧豆腐</span>
            </a>
        </div>
    </div>
</div>

{{/*左侧菜单*/}}
{{if gt .userInfo.Uid 0}}
<div class="sidebar-menu">
    {{$leftMenuLen:=len .leftMenuList}}
    {{if gt $leftMenuLen 0}}
        {{/*遍历所有的菜单信息*/}}
        {{range .leftMenuList}}
            {{/**一级菜单，它的href,有子菜单和无子菜单时不一样**/}}
            <a href="{{leftMenuUrl .MenuInfo}}" class="nav-header" {{if gt .MenuInfo.ChildMenuNum 0}}data-toggle="collapse"{{end}}>
                {{leftMenuIcon .MenuInfo}}
                {{.MenuInfo.MenuName}}
                {{if gt .MenuInfo.ChildMenuNum 0}}
                    <span class="glyphicon glyphicon-chevron-up"></span>
                {{end}}
            </a>
            {{/*如果有子菜单，遍历所有的子菜单*/}}
            {{if gt .MenuInfo.ChildMenuNum 0}}
                <ul id="left_menu_{{.MenuInfo.MenuId}}" class="nav nav-list collapse{{if inLeftMenu .SubMenuList $.SERVER_REQUEST_URL}} in{{end}}">
                    {{range .SubMenuList}}
                        <li {{if eq $.SERVER_REQUEST_URL .MenuPath}}class="active"{{end}}>
                            <a href="{{.MenuPath}}">
                                {{leftMenuIcon .}}
                                {{.MenuName}}
                            </a>
                        </li>
                    {{end}}
                </ul>
            {{end}}
        {{end}}
    {{end}}
</div>
{{end}}

{{/*主内容*/}}
{{if gt .userInfo.Uid 0}}
<div class="main-content" id="global-main-content">
    <div class="container-fluid">
        <div class="row-fluid">
{{end}}
{{end}}