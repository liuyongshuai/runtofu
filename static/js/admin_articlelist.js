/**
 * @author  Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date    2018-02-09 12:43
 */
$( function(){
    /*************查询的开始和结束时间*************/
    var t = new Date();
    var m = parseInt( t.getMonth() ) + 1;
    if( m < 10 ){
        m = "0" + m;
    }
    var today = t.getFullYear() + "-" + m + "-" + t.getDate();
    $( '#search_sctime' ).datetimepicker( {
        timepicker: false,
        formatDate: "Y-m-d",
        format: "Y-m-d",
        maxDate: today,
        onSelectDate: function(){
            var stime = $( "#search_sctime" ).val();
            if( stime.length > 0 ){
                $( '#search_ectime' ).datetimepicker( { minDate: stime } );
            }
        }
    } );
    $( '#search_ectime' ).datetimepicker( {
        timepicker: false,
        formatDate: "Y-m-d",
        format: "Y-m-d",
        maxDate: today,
        onSelectDate: function(){
            var etime = $( "#search_ectime" ).val();
            if( etime.length > 0 ){
                $( '#search_sctime' ).datetimepicker( { maxDate: etime } );
            }
        }
    } );
    //搜索操作
    $( "#search_btn" ).click( function(){
        var is_publish = $( "#search_is_publish" ).val();
        var is_rec = $( "#search_is_rec" ).val();
        var sctime = $( "#search_sctime" ).val();
        var ectime = $( "#search_ectime" ).val();
        location.href = "?is_publish=" + is_publish + "&sctime=" + sctime + "&ectime=" + ectime + "&is_rec=" + is_rec;
    } );
    //选择话题标签
    $( "#modiofy_article_metainfo_action_dialog span.enum-item" ).click( function(){
        if( $( this ).hasClass( "eitem-unsel" ) ){
            $( this ).removeClass( "eitem-unsel" );
            $( this ).addClass( "eitem-sel" );
        }
        else{
            $( this ).addClass( "eitem-unsel" );
            $( this ).removeClass( "eitem-sel" );
        }
    } );
    var dialogObj = $( "#modiofy_article_metainfo_action_dialog" );
    //修改文章的元信息
    $( "a[act=\"modify_article_meta_action\"]" ).click( function(){
        var trObj = $( this ).parents( "tr[article_id]" );
        var aid = 0;
        if( trObj ){
            aid = trObj.attr( "article_id" );
        }
        $( "#modify_article_id" ).val( aid );
        $( "#article_title" ).val( "" );
        $( "#article_is_origin" ).prop( "checked", false );
        dialogObj.find( "span.enum-item[tag_id]" ).removeClass( "eitem-sel" ).addClass( "eitem-unsel" );
        tofuUtils.sendRequest( {
            url: "/ajax/article/info",
            args: "article_id=" + aid,
            onSuccess: function( ainfo ){
                if( !$.isEmptyObject( ainfo ) ){
                    $( "#article_title" ).val( ainfo.title );
                    if( ainfo.is_origin ){
                        $( "#article_is_origin" ).prop( "checked", true );
                    }
                    $( ainfo.tag_list ).each( function( k, v ){
                        dialogObj.find( "span.enum-item[tag_id=" + v.tag_id + "]" ).removeClass( "eitem-unsel" ).addClass( "eitem-sel" );
                    } );
                }
                dialogObj.modal( 'show' );
            }
        } );

    } );
    //添加文章/修改文章信息
    $( "#modiofy_article_metainfo_action_submit_btn" ).click( function(){
        var article_id = $( "#modify_article_id" ).val();
        var title = $( "#article_title" ).val().trim();
        var tags = [];
        dialogObj.find( "span.enum-item[tag_id]" ).each( function(){
            if( $( this ).hasClass( "eitem-sel" ) ){
                var tid = $( this ).attr( "tag_id" );
                tags.push( tid )
            }
        } );
        if( title.length <= 0 ){
            return tofuUtils.modalMsg( { content: "标题不能为空" } );
        }
        if( tags.length <= 0 ){
            return tofuUtils.modalMsg( { content: "话题不能为空" } );
        }
        var is_origin = $( "#article_is_origin" ).is( ":checked" ) ? 1 : 0;
        var param = "article_id=" + article_id + "&is_origin=" + is_origin;
        param += "&title=" + encodeURIComponent( title );
        param += "&tags=" + tags.join( "," );
        param += "&action=modify";
        var url = "/ajax/article/modify";
        if( article_id <= 0 ){
            url = "/ajax/article/new";
        }
        tofuUtils.sendRequest( {
            url: url,
            args: param,
            onSuccess: function(){
                location.reload( true );
            }
        } );
    } );
    //修改文章属性信息
    $( "tr[article_id] a[action]" ).click( function(){
        var trObj = $( this ).parents( "tr[article_id]" );
        var aid = trObj.attr( "article_id" );
        var action = $( this ).attr( "action" );
        var param = "article_id=" + aid;
        if( action == "unpublish_article_action" ){//取消发布
            param += "&is_publish=0";
        }
        if( action == "publish_article_action" ){//发布
            param += "&is_publish=1";
        }
        if( action == "unrec_article_action" ){//取消推荐
            param += "&is_rec=0";
        }
        if( action == "rec_article_action" ){//推荐
            param += "&is_rec=1";
        }
        if( action == "del_article_action" ){//delete
            param += "&action=delete";
        }
        tofuUtils.sendRequest( {
            url: "/ajax/article/modify",
            args: param,
            onSuccess: function(){
                location.reload( true );
            }
        } );
    } );
} );