{{ template "header" .}}
{{ static_js  "admin_articlelist.js" }}

{{/******顶部的查询栏******/}}
<div class="panel panel-warning panel-search">
    <div class="panel-heading">查询文章</div>
    <div class="panel-body">
        <div class="form-inline">
            <div class="form-group">
                <div class="input-group">
                    <div class="input-group-addon">发布状态：</div>
                    <select class="form-control" id="search_is_publish">
                        <option value="-1">全部状态</option>
                        <option value="0" {{ if eq .is_publish  0}}selected{{ end }} >未发布</option>
                        <option value="1" {{ if eq .is_publish  1}}selected{{ end }} >已发布</option>
                    </select>
                </div>
            </div>
            <div class="form-group">
                <div class="input-group">
                    <div class="input-group-addon">推荐状态：</div>
                    <select class="form-control" id="search_is_rec">
                        <option value="-1">全部状态</option>
                        <option value="0" {{ if eq .is_rec  0}}selected{{ end }} >未推荐</option>
                        <option value="1" {{ if eq .is_rec  1}}selected{{ end }} >已推荐</option>
                    </select>
                </div>
            </div>
            <div class="form-group">
                <div class="input-group text200">
                    <div class="input-group-addon">开始时间：</div>
                    <input value="{{ .sctime }}" id="search_sctime" class="form-control" type="text">
                </div>
            </div>
            <div class="form-group">
                <div class="input-group text200">
                    <div class="input-group-addon">结束时间：</div>
                    <input value="{{ .ectime }}" id="search_ectime" class="form-control" type="text">
                </div>
            </div>
            <div class="form-group">
                <button id="search_btn" class="btn btn-default">查询</button>
            </div>
            <div class="form-group">
                <a type="button" href="javascript:void(0);" act="modify_article_meta_action" class="btn btn-primary">创建新文章</a>
            </div>
        </div>
    </div>
</div>

{{/******文章列表******/}}
{{$articleNum:=len .articleList}}
{{if gt $articleNum 0}}
    {{.pagination}}
    <table class="table table-bordered table-striped">
        <colgroup>
            <col class="col-xs-2">
            <col class="col-xs-1">
            <col class="col-xs-1">
            <col class="col-xs-1">
            <col class="col-xs-1">
            <col class="col-xs-1">
            <col class="col-xs-1">
        </colgroup>
        <thead>
            <tr>
                <th>标题</th>
                <th>发布</th>
                <th>推荐</th>
                <th>话题标签</th>
                <th>创建时间</th>
                <th>修改时间</th>
                <th>操作</th>
            </tr>
        </thead>
        <tbody>
            {{range .articleList}}
                <tr article_id="{{.ArticleId}}">
                    <td>
                        {{ if .IsOrigin }}
                            <span title='原创' class="origin ori1">原</span>
                        {{ else }}
                            <span title='转载引用' class="origin ori0">转</span>
                        {{ end }}
                        <a href="https://runtofu.com/article/{{.ArticleId}}" target="_blank">
                            {{.Title}}
                        </a>
                    </td>
                    <td>
                        {{if .IsPublish}}
                            <span style="color:green;font-weight:bold;">是</span>
                        {{else}}
                            <span style="color:red;font-weight:bold;">否</span>
                        {{end}}
                    </td>
                    <td>
                        {{if .IsRec}}
                            <span style="color:green;font-weight:bold;">是</span>
                        {{else}}
                            <span style="color:red;font-weight:bold;">否</span>
                        {{end}}
                    </td>
                    <td>{{range .TagList}}{{.TagName}}<br/>{{end}}</td>
                    <td>{{ftime .CreateTime "Y-m-d H:i"}}</td>
                    <td>{{ftime .LastModifyTime "Y-m-d H:i"}}</td>
                    <td>
                        <div class="btn-group btn-group-xs">
                            <a href="/articleinfo?article_id={{.ArticleId}}" type="button" class="btn btn-link" target="_blank">编辑正文</a>
                            <a href="javascript:void(0);" type="button" class="btn btn-link" act="modify_article_meta_action">元数据</a>
                            {{if .IsPublish}}
                                <a href="javascript:void(0);" type="button" class="btn btn-link" action="unpublish_article_action">取消发布</a>
                            {{else}}
                                <a href="javascript:void(0);" type="button" class="btn btn-link" action="publish_article_action">发布</a>
                            {{end}}
                            {{if .IsRec}}
                                <a href="javascript:void(0);" type="button" class="btn btn-link" action="unrec_article_action">取消推荐</a>
                            {{else}}
                                <a href="javascript:void(0);" type="button" class="btn btn-link" action="rec_article_action">推荐</a>
                            {{end}}
                            <a href="javascript:void(0);" type="button" class="btn btn-link" action="del_article_action">删除</a>
                        </div>
                    </td>
                </tr>
            {{end}}
        </tbody>
    </table>
    {{.pagination}}
{{else}}
    <div class="warning_box">暂无文章信息</div>
{{end}}

{{/******编辑文章元信息的弹框******/}}
{{if gt .userInfo.Uid 0}}
<div id="modiofy_article_metainfo_action_dialog" class="modal" tabindex="-1" role="dialog">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal">
                    <span aria-hidden="true">×</span>
                </button>
                <h4 class="modal-title">添加文章</h4>
            </div>
            <div class="modal-body">
                <div class="form-horizontal">
                    <div class="form-group">
                        <label class="col-sm-2 control-label">文章标题</label>
                        <div class="col-sm-9">
                            <input value="" class="form-control" id="article_title" type="text">
                        </div>
                    </div>
                    <div class="form-group">
                        <label class="col-sm-2 control-label">话题标签</label>
                        <div class="col-sm-9">
                            {{range .allTagList}}
                            <span class="enum-item eitem-unsel" tag_id="{{.TagId}}">{{.TagName}}</span>
                            {{end}}
                        </div>
                    </div>
                    <div class="form-group">
                        <div class="col-sm-offset-2 col-sm-10">
                            <div class="checkbox">
                                <label>
                                    <input name="article_is_origin" id="article_is_origin" type="checkbox" value="1"/>是否原创
                                </label>
                            </div>
                        </div>
                    </div>
                    <div class="form-group">
                        <div class="col-sm-offset-5 col-sm-8">
                            <input type="hidden" id="modify_article_id" value="0">
                            <button id="modiofy_article_metainfo_action_submit_btn" class="btn btn-default">确定</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
{{end}}
{{ template "footer" .}}