<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>{{.SITE_NAME}}</title>
    <link rel="Shortcut Icon" href="{{ static_image "favicon.ico" }}" type="image/x-icon"/>
    <link rel="Bookmark" href="{{ static_image "favicon.ico" }}"/>
    <meta content="IE=edge,chrome=1" http-equiv="X-UA-Compatible">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script>
        var STATIC_PREFIX = "{{.STATIC_PREFIX}}";
    </script>
    {{ static_js  "base.js" }}
    {{ static_css  "base.css" }}
    {{ static_css  "admin.css" }}
    {{ static_js  "admin_base.js" }}
    {{ editor_js "editormd.js"  }}
    {{ editor_css "css/editormd.css"  }}
    <!--[if lt IE 9]>
     {{ static_js  "html5shiv.js" }}
    <![endif]-->
    <script>
        var articleId = "{{.articleId}}";
        var isHaveModifyPriv = {{if gt .userInfo.Uid 0}}1{{else}}0{{end}};
    </script>
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
    {{ static_js  "admin_articleinfo.js" }}
    <div class="container-fluid">
        <div class="row">
            <div class="col-md-12" id="load_local_articleinfo_cache" style="display: none">
                <p class="bg-danger local-storage-tips">发现此内容存在本地缓存，是否载入？
                    <span></span>
                    <button class="btn btn-primary btn-xs" op="load">载入本地缓存数据</button>
                    <button class="btn btn-primary btn-xs" op="clear">清除本地缓存数据</button>
                </p>
            </div>
        </div>
        <div class="row">
            <div class="col-md-12">
                <div class="edit-wrapper" id="edit_article_wrapper_markdown"></div>
            </div>
        </div>
    </div>
    <span class="glyphicon glyphicon-eye-open show-foot-bar"></span>
    <div class="form-group edit-foot-bar">
        <div class="col-sm-10">
            <span id="article_title_for_edit" class="editable" title="文章标题" data-toggle="tooltip"></span>
        </div>
        <div class="col-sm-2">
            <button id="save_article_submit_btn" class="btn btn-default">提交保存</button>
            <span>
                <a target="_blank" href="https://runtofu.com/article/{{.articleId}}">
                    detail
                </a>
            </span>
        </div>
    </div>
</body>
</html>