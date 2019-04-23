/**
 * @author  Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date    2018-02-10 22:09
 */
$( function(){
    /************添加弹框************/
    $( "#add_article_tag_action" ).click( function(){
        $( "#article_tag" ).val( "" );
        $( "#add_article_tag_action_dialog" ).modal( 'show' );
    } );
    /************提交操作************/
    $( "#add_article_tag_action_submit_btn" ).click( function(){
        var tag_name = $( "#article_tag" ).val().trim();
        if( tag_name.length <= 0 ){
            return tofuUtils.modalMsg( { type: "danger", content: "话题名称不能为空" } );
        }
        tofuUtils.sendRequest( {
            url: "/ajax/tag/add",
            args: "tag_name=" + encodeURIComponent( tag_name ),
            onSuccess: function( tagInfo ){
                $( "span.enum-item[tag_id=\"" + tagInfo.tag_id + "\"]" ).remove();
                $( "#article_tag_item_list" ).append( '<span class="enum-item" tag_id="' + tagInfo.tag_id + '">' + tagInfo.tag_name + '<span class="glyphicon glyphicon-remove"></span></span>' );
                $( "#add_article_tag_action_dialog" ).modal( 'hide' );
            }
        } );
    } );
    /************删除操作************/
    $( "#article_tag_item_list" ).delegate( "span.enum-item[tag_id] > span.glyphicon-remove", "click", function(){
        var _this = $( this );
        var tag_id = _this.parent().attr( "tag_id" );
        tofuUtils.sendRequest( {
            url: "/ajax/tag/delete",
            args: "tag_id=" + tag_id,
            onSuccess: function(){
                _this.parent().remove();
            }
        } );
    } );
} );