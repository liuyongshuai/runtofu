{{ template "header" .}}
<script>
    var menuList = {{json_encode .leftMenuList}};
</script>

{{/**新建菜单的按钮**/}}
<div class="action-wrapper">
    <a href="javascript:void(0);" action="editMenuInfo" type="button" class="btn btn-default pull-right">添加菜单</a>
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
<script>
$(function () {
    /***********展示/隐藏给菜单添加用户的十字按钮***********/
    $("span.glyphicon-plus").parent().mouseover(function () {
        $(this).find("span.glyphicon-plus").show()
    }).mouseout(function () {
        $(this).find("span.glyphicon-plus").hide();
    });
    /***********给菜单添加用户***********/
    $("span.glyphicon-plus").click(function () {
        var menu_id = parseInt($(this).parents("tr[menu_id]").attr("menu_id"));
        if (isNaN(menu_id)) {
            menu_id = 0;
        }
        var _this = $(this);
        _this.popoverEdit({
            type: "select",
            selectUrl: "/ajax/system/get_nopriv_userlist?menu_id=" + menu_id,
            onConfirm: function (input) {
                var user_name = input.val();
                poiUtils.sendRequest({
                    url: "/ajax/system/add_usermenu",
                    args: "user_name=" + user_name + "&menu_id=" + menu_id,
                    onSuccess: function () {
                        var userName = input.val();
                        var txt = input.find("option:selected").text();
                        _this.before('<span user_name="' + userName + '" class="enum-item">' + txt + '<span class="glyphicon glyphicon-remove"></span></span>');
                    }
                });
            }
        });
    });
    /***********删除菜单下的用户***********/
    $("table.table").delegate("span.glyphicon-remove", "click", function () {
        var _this = $(this);
        var user_name = _this.parent().attr("user_name");
        var menu_id = parseInt(_this.parents("tr[menu_id]").attr("menu_id"));
        poiUtils.sendRequest({
            url: "/ajax/system/del_usermenu",
            args: "user_name=" + user_name + "&menu_id=" + menu_id,
            onSuccess: function () {
                _this.parent().remove();
            }
        });
    });
    /************删除菜单项************/
    $("a[action='delMenuInfo']").click(function () {
        var _this = $(this);
        var menu_id = parseInt(_this.parents("tr[menu_id]").attr("menu_id"));
        poiUtils.sendRequest({
            url: "/ajax/system/del_menu",
            args: "menu_id=" + menu_id,
            onSuccess: function () {
                _this.parent().parent().remove();
            }
        });
    });
    /***********初始化父菜单选择下拉框的数据***********/
    var initParentSelect = function () {
        var options = '<option value="0">&nbsp;</option>';
        $("#menu_parent_menu").empty();
        $.each(menuList, function (index, info) {
            options += '<option value="' + info.menu_info.menu_id + '">' + info.menu_info.menu_name + '</option>';
        });
        $("#menu_parent_menu").append(options);
    };
    /***********初始化颜色选择器***********/
    var initIconColor = function (color) {
        $("#menu_icon_color").ColorPickerSliders({
            placement: 'right',
            flat: false,
            color: color,
            title: "选择菜单图标的颜色",
            sliders: false,
            swatches: false,
            hsvpanel: true,
            previewformat: 'hex'
        });
    };
    /************编辑菜单信息************/
    $("[action=\"editMenuInfo\"]").click(function () {
        $("#menu_name").val("");
        $("#menu_path").val("");
        initParentSelect();
        $("#menu_icon").val("");
        $("#menu_icon_color").val("").css("background-color", "white");
        $("#menu_desc").text("");
        $("#menu_type").prop("checked", false);
        var menu_id = parseInt($(this).parents("tr[menu_id]").attr("menu_id"));
        menu_id = parseInt(menu_id);
        if (isNaN(menu_id)) {
            menu_id = 0;
        }
        $("#menu_menu_id").val(menu_id);
        $("#menu_parent_menu").parents("div.form-group").show();
        $("#menu_path").parents("div.form-group").show();
        $("#menu_type").parents("div.form-group").show();
        poiUtils.sendRequest({
            url: "/ajax/system/getMenuInfo",
            args: "menu_id=" + menu_id,
            onSuccess: function (menu_info) {
                var originColor = "";
                if (!$.isEmptyObject(menu_info)) {
                    $("#menu_name").val(menu_info.menu_name);
                    $("#menu_path").val(menu_info.menu_path);
                    $("#menu_icon").val(menu_info.menu_icon).parent().find("span.glyphicon").addClass(menu_info.menu_icon);
                    originColor = menu_info.menu_icon_color;
                    $("#menu_desc").text(menu_info.menu_desc);
                    if (parseInt(menu_info.menu_type) <= 1) {
                        $("#menu_type").prop("checked", true);
                    }
                    //如果有父菜单，说明是子菜单
                    if (parseInt(menu_info.menu_parent_id) > 0) {
                        $("#menu_type").parents("div.form-group").hide();
                        $("#menu_parent_menu option[value=0]").remove();
                        $("#menu_parent_menu option[value=\"" + menu_info.menu_parent_id + "\"]").attr("selected", "true");

                    } else {
                        $("#menu_parent_menu").empty();
                        $("#menu_parent_menu").attr("disabled", true);
                        $("#menu_parent_menu").parents("div.form-group").hide();
                        $("#menu_path").parents("div.form-group").hide();
                    }
                }
                $("#menu_parent_menu").selectpicker("destroy").selectpicker({maxOptions: 1});
                initIconColor(originColor);
                $("#menu_edit_dialog").modal('show');
            }
        });
    });
    /************上调顺序************/
    $("a[action='upMenuSort']").click(function () {
        var trObj = $(this).parents("tr[menu_id]");
        var prev = trObj.prev();
        var menu_id = parseInt(trObj.attr("menu_id"));
        if (isNaN(menu_id) || menu_id <= 0) {
            return poiUtils.modalMsg({content: "get menu_id failed"});
        }
        poiUtils.sendRequest({
            url: "/ajax/system/up_menusort",
            args: "menu_id=" + menu_id,
            onSuccess: function () {
                window.location.reload();
            }
        });
    });
    /************下调顺序************/
    $("a[action='downMenuSort']").click(function () {
        var trObj = $(this).parents("tr[menu_id]");
        var next = trObj.next();
        var menu_id = parseInt(trObj.attr("menu_id"));
        if (isNaN(menu_id) || menu_id <= 0) {
            return poiUtils.modalMsg({content: "get menu_id failed"});
        }
        poiUtils.sendRequest({
            url: "/ajax/system/down_menusort",
            args: "menu_id=" + menu_id,
            onSuccess: function () {
                window.location.reload();
            }
        });
    });
    /************提交菜单信息************/
    $("#edit_menu_submit_btn").click(function () {
        var menu_id = $("#menu_menu_id").val();
        menu_id = parseInt(menu_id);
        if (isNaN(menu_id)) {
            menu_id = 0;
        }
        var name = $("#menu_name").val();
        var path = $("#menu_path").val();
        var icon = $("#menu_icon").val();
        var icon_color = $("#menu_icon_color").val();
        var desc = $("#menu_desc").val();
        var parent_menu_id = $("#menu_parent_menu").val();
        parent_menu_id = parseInt(parent_menu_id);
        if (isNaN(parent_menu_id)) {
            parent_menu_id = 0;
        }
        var menu_type = $("#menu_type").is(":checked") ? 1 : 10;
        var args = "menu_id=" + menu_id + "&menu_name=" + encodeURIComponent(name) + "&menu_path=" + encodeURIComponent(path);
        args += "&menu_icon=" + icon + "&menu_icon_color=" + icon_color + "&menu_desc=" + encodeURIComponent(desc);
        args += "&menu_parent_id=" + parent_menu_id + "&menu_type=" + menu_type;
        poiUtils.sendRequest({
            url: "/ajax/system/modify_menuinfo",
            args: args,
            onSuccess: function () {
                location.reload(true);
            }
        });
        return false;
    });
    /************展开子菜单************/
    $("button[action='showsubmenu']").click(function () {
        var menu_id = parseInt($(this).parents("tr[menu_id]").attr("menu_id"));
        $("tr[parent_menu_id][parent_menu_id!=\"" + menu_id + "\"]").hide();
        var sub = $("tr[parent_menu_id=\"" + menu_id + "\"]");
        if (sub.is(":visible")) {
            sub.hide();
            location.hash = "";
        } else {
            sub.show();
            location.hash = "#" + menu_id;
        }
    });
    /************处理URL后跟的锚点信息************/
    var menu_id = parseInt(location.hash.substring(1));
    if (menu_id > 0) {
        $("tr[parent_menu_id=\"" + menu_id + "\"]").show();
    }
    /************图标选择器************/
    $("#menu_icon").popover({
        html: true,
        placement: "right",
        title: "选择菜单图标",
        container: "body",
        content: function () {
            var icons = "<div class='icon-picker'>";
            $.each(bootstrapIconList, function (index, val) {
                icons += '<span class="glyphicon ' + val + '"></span>';
            });
            icons += "</div>";
            return icons;
        }
    }).on('shown.bs.popover', function () {
        $("div.icon-picker > span").tooltip({
            html: true,
            placement: "auto",
            title: function () {
                var cs = "enlarge";
                var cls = $(this).attr("class");
                var t = "<span class='" + cls + " " + cs + "'></span>";
                return t;
            }
        });
    });
    $("body").delegate("div.icon-picker > span", "click", function () {
        var cls = $(this).attr("class");
        cls = cls.split(" ");
        cls = cls[1];
        var oldCls = $("#menu_icon").val();
        $("#menu_icon").val(cls).parent().find("span.glyphicon").removeClass(oldCls).addClass(cls);
        $('#menu_icon').popover('hide')
    });
});
</script>
{{ template "footer" .}}