$( function(){
    tofuUtils.loading.start();
    tofuUtils.sendRequest( {
        url: "/ajax/article/info/" + aid,
        onSuccess: function( articleInfo ){
            $( "div[role=title]" ).text( articleInfo.title );
            var tmp = [];
            //所有的话题标签列表
            $( articleInfo.tag_list ).each( function( k, v ){
                tmp.push( "<a href=\"/tag/" + v.tag_id + "\" target='_blank'>" + v.tag_name + "</a>" )
            } );
            $( "span[role=tags]" ).html( tmp.join( "、" ) );
            //创建时间
            $( "span[role=ctime]" ).text( tofuUtils.date( "Y-m-d H:i", articleInfo.create_time ) );
            //最后修改时间
            $( "span[role=lmtime]" ).text( tofuUtils.date( "Y-m-d H:i", articleInfo.last_modify_time ) );
            //原创和转载的标志
            if( articleInfo.is_origin ){
                $( "span.isOrigin" ).replaceWith( $( "<span title='原创' class=\"origin ori1\">原</span>" ) );
            }
            else{
                $( "span.isOrigin" ).replaceWith( $( "<span title='转载引用' class=\"origin ori0\">转</span>" ) );
            }
            wikiEditor = editormd.markdownToHTML( "tofu_info_markdown_wrapper", {
                markdown: articleInfo.content,
                htmlDecode: true,
                taskList: true,
                tocm: true,
                tocDropdown: true,
                toc: true,
                tex: true,
                flowChart: true,
                sequenceDiagram: true
            } );
            document.title = articleInfo.title;
        },
        onError: function( ret ){
            tofuUtils.modalMsg( {
                content: ret.result.errmsg + "，点确定返回首页",
                type: "danger",
                onConfirm: function(){
                    location.href = "/"
                }
            } );
        }
    } );

} );
