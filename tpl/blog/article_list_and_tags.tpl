{{/*首页、推荐列表页等的格式是一样的，就抽出来了*/}}
{{$articleNum:=len .articleList}}
{{if gt $articleNum 0}}
<div class="tofu-item-list col-sm-8 col-md-8">
    {{ range .articleList}}
    <div class="tofu-item">
        <div class="title">
            <a href="/article/{{ .ArticleId }}" target="_blank">
                {{ .Title }}
            </a>
        </div>
        <div>
            {{ if .IsOrigin }}
                <span title='原创' class="origin ori1">原</span>
            {{ else }}
                <span title='转载引用' class="origin ori0">转</span>
            {{ end }}
            {{ range .TagList }}
                <a href="/tag/{{ .TagId }}" class="tag" target="_blank">{{ .TagName}}</a>
            {{ end }}
            <span class="ctime">{{ ftime .CreateTime "Y-m-d H:i"}}</span>
        </div>
    </div>
    {{end}}
    {{ .pagination }}
</div>
{{end}}
{{/*页面右边的热门话题列表*/}}
{{$tagNum:=len .tagList}}
{{if gt $tagNum 0}}
<div class="col-sm-4 col-md-4 tofu-left-bar">
    <div class="panel panel-default">
      <div class="panel-heading">标签列表</div>
      <div class="panel-body list-group">
          {{ range .tagList}}
          <button type="button" class="list-group-item">
            <a href="/tag/{{ .TagId }}" target="_blank">{{ .TagName }}</a>
            <span class="badge">{{ .ContentNum }}</span>
          </button>
          {{ end }}
      </div>
    </div>
</div>
{{end}}