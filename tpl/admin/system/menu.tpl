{{ template "header" .}}
{{ static_js  "colorpickersliders.js" }}
{{ static_css  "colorpickersliders.css" }}
<script>
    var menuList = {{json_encode .leftMenuList}};
</script>

{{/**新建菜单的按钮**/}}
<div class="action-wrapper">
    <a href="javascript:void(0);" action="editmenu" type="button" class="btn btn-default pull-right">添加菜单</a>
</div>

{{/******菜单列表******/}}
{{$menuNum:=len .leftMenuList}}
{{if gt $menuNum 0}}
<table class="table table-bordered table-striped">
    <colgroup>
        <col class="col-xs-1">
        <col class="col-xs-1">
        <col class="col-xs-1">
        <col class="col-xs-1">
    </colgroup>
    <thead>
        <tr>
            <th>菜单名称</th>
            <th>子菜单数</th>
            <th>路径</th>
            <th>操作</th>
        </tr>
    </thead>
    <tbody>
    {{range .leftMenuList}}
        {{if gt .MenuInfo.MenuId 0}}
            <tr menu_id={{.MenuInfo.MenuId}}>
                <td>
                    {{leftMenuIcon .MenuInfo}}
                    {{.MenuInfo.MenuName}}
                </td>
                <td>
                    <button data-toggle="tooltip" data-placement="top" action="showsubmenu" type="button" class="btn btn-link">
                        {{.MenuInfo.ChildMenuNum}}
                    </button>
                </td>
                <td>{{.MenuInfo.MenuPath}}</td>
                <td>
                    <div class="btn-group btn-group-xs">
                        <a href="javascript:void(0);" action="editmenu" type="button" class="btn btn-link">编辑</a>
                        <a href="javascript:void(0);" action="delmenu" type="button" class="btn btn-link">删除</a>
                        <a href="javascript:void(0);" type="button" class="btn btn-link" action="upmenu">
                            <span class="glyphicon glyphicon-arrow-up"></span>
                        </a>
                        <a href="javascript:void(0);" type="button" class="btn btn-link" action="downmenu">
                            <span class="glyphicon glyphicon-arrow-down"></span>
                        </a>
                    </div>
                </td>
            </tr>
            {{if gt .MenuInfo.ChildMenuNum 0}}
            <tr parent_menu_id="{{.MenuInfo.MenuId}}" style="display: none">
                <td colspan="4">
                    <table class="table table-hover">
                        <colgroup>
                            <col class="col-xs-1">
                            <col class="col-xs-1">
                            <col class="col-xs-1">
                            <col class="col-xs-1">
                        </colgroup>
                        {{range .SubMenuList}}
                            <tr menu_id="{{.MenuId}}">
                                <td>
                                    {{leftMenuIcon .}}
                                    {{.MenuName}}
                                </td>
                                <th>---</th>
                                <td>
                                    {{.MenuPath}}
                                </td>
                                <td>
                                    <div class="btn-group btn-group-xs">
                                        <a href="javascript:void(0);" action="editmenu" type="button" class="btn btn-link">编辑</a>
                                        <a href="javascript:void(0);" action="delmenu" type="button" class="btn btn-link">删除</a>
                                        <a href="javascript:void(0);" type="button" class="btn btn-link" action="upmenu">
                                            <span class="glyphicon glyphicon-arrow-up"></span>
                                        </a>
                                        <a href="javascript:void(0);" type="button" class="btn btn-link" action="downmenu">
                                            <span class="glyphicon glyphicon-arrow-down"></span>
                                        </a>
                                    </div>
                                </td>
                            </tr>
                        {{end}}
                    </table>
                </td>
            </tr>
            {{end}}
        {{end}}
    {{end}}
    </tbody>
</table>
{{end}}

{{/*****添加菜单的弹出对话框*****/}}
<div id="menu_edit_dialog" class="modal" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-sm">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal">
                    <span aria-hidden="true">×</span>
                </button>
                <h4 class="modal-title">编辑菜单信息</h4>
            </div>
            <div class="modal-body">
                <form class="form-horizontal" role="form">
                    {{/*菜单名称*/}}
                    <div class="form-group">
                        <label for="menu_name" class="col-sm-3 control-label">名称</label>
                        <div class="col-sm-9">
                            <input value="" class="form-control" id="menu_name" placeholder="菜单名称" type="text">
                        </div>
                    </div>
                    {{/*父菜单信息*/}}
                    <div class="form-group">
                        <label for="menu_parent_menu" class="col-sm-3 control-label">父菜单</label>
                        <div class="col-sm-9">
                            <select class="form-control" id="menu_parent_menu"></select>
                        </div>
                    </div>
                    {{/*当前菜单的路径信息*/}}
                    <div class="form-group">
                        <label for="menu_path" class="col-sm-3 control-label">path</label>
                        <div class="col-sm-9">
                            <input value="" class="form-control" id="menu_path" placeholder="二级菜单的URL" type="text">
                        </div>
                    </div>
                    {{/*当前菜单在左侧菜单框显示的图标*/}}
                    <div class="form-group has-feedback">
                        <label for="menu_icon" class="col-sm-3 control-label">图标</label>
                        <div class="col-sm-9">
                            <input value="" readonly class="form-control" id="menu_icon_name" data-toggle="popover" placeholder="菜单左侧的图标" type="text">
                            <span class="glyphicon form-control-feedback"></span>
                        </div>
                    </div>
                    {{/*当前菜单的图标的颜色*/}}
                    <div class="form-group">
                        <label for="menu_icon_color" class="col-sm-3 control-label">图标颜色</label>
                        <div class="col-sm-9">
                            <input value="" readonly class="form-control" id="menu_icon_color" placeholder="图标的颜色" type="text">
                        </div>
                    </div>
                    {{/*提交操作*/}}
                    <div class="form-group">
                        <div class="col-sm-offset-4 col-sm-8">
                            <input value="" id="menu_menu_id" type="hidden">
                            <button id="edit_menu_submit_btn" class="btn btn-default">确定</button>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
{{ static_js  "admin_menu.js" }}
{{ template "footer" .}}