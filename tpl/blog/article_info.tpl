{{ template "header" .}}
{{ editor_js "lib/marked.min.js" }}
{{ editor_js "lib/prettify.min.js"  }}
{{ editor_js "lib/raphael.min.js"  }}
{{ editor_js "lib/underscore.min.js"  }}
{{ editor_js "lib/sequence-diagram.min.js"  }}
{{ editor_js "lib/flowchart.min.js"  }}
{{ editor_js "lib/jquery.flowchart.min.js"  }}
{{ editor_css "css/editormd.preview.css"  }}
{{ editor_js "editormd.js"  }}
<script type="text/javascript">var aid = "{{.aid}}";</script>
<div class="content-wrapper">
    {{/******文章顶部的标题、标签、创建时间、是否原创等******/}}
    <div class="header">
        {{/*标题，加载时会替换此处*/}}
        <div class="title" role="title">大标题</div>
        <div class="s">
            {{/*是否原创占位符*/}}
            <span class="isOrigin"></span>
            {{/*文章所属的话题标签等信息*/}}
            <span class="t">
                标签：
                <span role="tags">
                    <a>gin</a>、
                    <a>golang</a>
                </span>
            </span>
            {{/*创建时间及最后的修改时间信息*/}}
            <span class="pull-right">
                创建时间：<span role="ctime">2018-01-01 22:22</span>
                &nbsp;&nbsp;&nbsp;&nbsp;
                最后修改时间：<span role="lmtime">2018-01-01 22:22</span>
            </span>
        </div>
    </div>
    {{/******通过Ajax异步加载的文章详细内容******/}}
    <div class="markdown-body" id="tofu_info_markdown_wrapper"></div>
</div>
{{ static_js "article.js"  }}

{{ template "footer" .}}