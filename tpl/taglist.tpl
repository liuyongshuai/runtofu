{{ template "header" .}}
{{/*顶部的相关操作栏*/}}
<div class="panel panel-warning panel-search">
    <div class="panel-heading">相关操作</div>
    <div class="panel-body">
        <div class="form-inline">
            <div class="form-group">
                <button id="add_article_tag_action" type="button" class="btn btn-default">添加话题标签</button>
            </div>
        </div>
    </div>
</div>

{{/*顶部的相关操作栏*/}}
{{/******文章列表******/}}
{{$tagNum:=len .tagList}}
{{if gt $tagNum 0}}
<div id="article_tag_item_list">
    {{range .tagList}}
        <span class="enum-item" tag_id="{{.TagId}}">
            {{.TagName}}（{{.ContentNum}}）
            <span class="glyphicon glyphicon-remove"></span>
        </span>
    {{end}}
</div>
{{else}}
    <div class="warning_box">暂无话题信息</div>
{{end}}

{{/*添加话题的弹框*/}}
<div id="add_article_tag_action_dialog" class="modal" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-sm">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal">
                    <span aria-hidden="true">×</span>
                </button>
                <h4 class="modal-title">添加标签</h4>
            </div>
            <div class="modal-body">
                <div class="form-horizontal">
                    <div class="form-group">
                        <label class="col-sm-3 control-label">话题标签</label>
                        <div class="col-sm-8">
                            <input value="" class="form-control" id="article_tag" placeholder="" type="text">
                        </div>
                    </div>
                    <div class="form-group">
                        <div class="col-sm-offset-5 col-sm-8">
                            <button id="add_article_tag_action_submit_btn" class="btn btn-default">确定</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
{{ static_js  "admin_taglist.js" }}
{{ template "footer" .}}