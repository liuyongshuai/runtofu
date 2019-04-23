/**
 * @author  Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date    2018-02-09 17:19
 */
$( function(){
    /*********markdown内容实时预览*********/
    var localStoreKey = "runtofu_articleinfo_local_" + articleId;
    var wikiEditorMarkdown = null;
    var titleDomObj = $( "#article_title_for_edit" );
    /*********如果本地的缓存不是最新的则请求服务端*********/
    tofuUtils.sendRequest( {
        url: "/ajax/article/info",
        args: "article_id=" + articleId,
        onSuccess: function( articleInfo ){
            var content = "";
            var title = "";
            if( !$.isEmptyObject( articleInfo ) ){
                content = articleInfo.content;
                title = articleInfo.title;
                document.title = title;
            }
            titleDomObj.text( title );
            /*********初始化显示markdown格式的编辑器*********/
            wikiEditorMarkdown = editormd( "edit_article_wrapper_markdown", {
                width: "100%",
                toolbar: true,
                autoHeight: true,
                watch: true,
                path: STATIC_PREFIX + "/editor.md/lib/",
                markdown: content,
                saveHTMLToTextarea: true,
                toolbarIcons: "runtofu",
                htmlDecode: true,
                imageUpload: true,
                imageFormats: ["jpg", "jpeg", "gif", "png", "bmp", "PNG", "JPG", "JPEG", "GIF"],
                imageUploadURL: "/ajax/system/uploadImage",
                placeholder: "点击此处输入。。",
                tocm: true,
                toc: true,
                tocDropdown: true,
                tex: true,
                flowChart: true,
                sequenceDiagram: true,
                taskList: true,
                readOnly: isHaveModifyPriv <= 0,
                onfullscreen: function(){
                    $( "div.edit-foot-bar" ).hide();
                },
                onfullscreenExit: function(){
                    $( "div.edit-foot-bar" ).show();
                },
                onchange: function(){
                    var localContent = this.getValue();
                    var localTitle = titleDomObj.text();
                    tofuUtils.store.set( localStoreKey, {
                        "localContent": localContent,
                        "localTitle": localTitle,
                        "localLastModifyTime": tofuUtils.time()
                    } );
                }
            } );
            /*********检测是否有本地缓存*********/
            var localInfo = tofuUtils.store.get( localStoreKey );
            if( !$.isEmptyObject( localInfo ) && localInfo.localContent.length > 0 && localInfo.localLastModifyTime > 0 ){
                var obj = $( "#load_local_articleinfo_cache" );
                obj.show();
                obj.find( "span" ).text( "最后修改于 " + tofuUtils.date( "Y-m-d H:i:s", localInfo.localLastModifyTime ) );
                obj.find( "button[op=\"load\"]" ).click( function(){
                    wikiEditorMarkdown.setValue( localInfo.localContent );
                    titleDomObj.text( localInfo.localTitle );
                    obj.remove();
                } );
                obj.find( "button[op=\"clear\"]" ).click( function(){
                    tofuUtils.store.remove( localStoreKey );
                    obj.remove();
                } );
            }
        }
    } );
    /*********显示和隐藏下面的固定的bar*********/
    $( "span.show-foot-bar" ).click( function(){
        var _this = $( this );
        if( _this.hasClass( "glyphicon-eye-open" ) ){
            $( "div.edit-foot-bar" ).hide();
            _this.removeClass( "glyphicon-eye-open" ).addClass( "glyphicon-eye-close" );
        }
        else{
            $( "div.edit-foot-bar" ).show();
            _this.removeClass( "glyphicon-eye-close" ).addClass( "glyphicon-eye-open" );
        }
    } );
    /*********编辑文章标题*********/
    $( "#article_title_for_edit" ).click( function(){
        var _this = $( this );
        _this.popoverEdit( {
            title: "修改名称",
            css: "width:500px",
            onConfirm: function( _el ){
                tofuUtils.sendRequest( {
                    url: "/ajax/article/modify",
                    args: "article_id=" + articleId + "&title=" + encodeURIComponent( _el.val() ),
                    onSuccess: function(){
                        _this.text( _el.val() );
                    }
                } );
            }
        } );
    } );
    /*********提交文章信息*********/
    $( "#save_article_submit_btn" ).click( function(){
        var btn = $( this );
        var lbtn = btn.button( 'loading' );
        btn.attr( "disabled", "disabled" );
        var mk = wikiEditorMarkdown.getMarkdown();
        tofuUtils.loading.start();
        tofuUtils.sendRequest( {
            url: "/ajax/article/modify",
            args: "article_id=" + articleId + "&content=" + encodeURIComponent( mk ),
            onSuccess: function(){
                tofuUtils.store.remove( localStoreKey );
                lbtn.button( 'reset' );
                btn.removeAttr( "disabled" );
                location.href = "/articleinfo?article_id=" + articleId;
            }
        } );
    } );
} );
