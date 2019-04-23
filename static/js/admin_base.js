/**
 * @author  Liu Yongshuai<liuyongshuai@hotmail.com>
 * @date    2018-02-09 16:37
 */
$( function(){
    /************菜单栏的隐藏和显示************/
    var sideBar = $( "div.sidebar-menu" );
    var mainContent = $( "div.main-content" );
    var mainMarginLeft = mainContent.css( "margin-left" );
    var expIcon = $( "span.tofu-expand-icon" );
    var TOGGLE_SIDEBAR_DISPLAY_FLAG = "runtofu_side_bar_show_flag";
    var flagVal = tofuUtils.store.get( TOGGLE_SIDEBAR_DISPLAY_FLAG, 1 );
    var showSideBar = function(){
        sideBar.show();
        mainContent.css( { "margin-left": mainMarginLeft } );
        expIcon.attr( "isShow", 1 );
        expIcon.removeClass( "glyphicon-resize-small" ).addClass( "glyphicon-move" );
    };
    var hideSideBar = function(){
        sideBar.hide();
        mainContent.css( { "margin-left": "0px" } );
        expIcon.attr( "isShow", 0 );
        expIcon.removeClass( "glyphicon-move" ).addClass( "glyphicon-resize-small" );
    };
    if( flagVal > 0 ){
        showSideBar();
    }
    else{
        hideSideBar();
    }
    expIcon.click( function(){
        var _this = $( this );
        var isShow = parseInt( _this.attr( "isShow" ).trim() );
        if( isShow > 0 ){
            hideSideBar();
            tofuUtils.store.set( TOGGLE_SIDEBAR_DISPLAY_FLAG, 0 );
        }
        else{
            showSideBar();
            tofuUtils.store.set( TOGGLE_SIDEBAR_DISPLAY_FLAG, 1 );
        }
    } );
    /************登录相关************/
    $( "#login_submit_action" ).click( function(){
        var name = $( "#name" ).val().trim();
        var passwd = $( "#passwd" ).val().trim();
        if( name.length <= 0 ){
            return tofuUtils.modalMsg( { content: "用户名不能为空" } );
        }
        if( passwd.length <= 0 ){
            return tofuUtils.modalMsg( { content: "密码不能为空" } );
        }
        tofuUtils.sendRequest( {
            url: "/ajax/login",
            args: "name=" + name + "&passwd=" + passwd,
            onSuccess: function(){
                location.href = "/"
            }
        } );
    } );
    /************退出登录相关************/
    $( "#adminNavigationLogout" ).click( function(){
        tofuUtils.sendRequest( {
            url: "/ajax/logout",
            onSuccess: function(){
                location.href = "/login"
            }
        } );
    } );
    /************修改密码相关************/
    $( "#adminNavigationChangePasswd" ).click( function(){
        $("#modify_user_passwd_action_dialog").modal( 'show' );
    });
    $( "#modify_user_passwd_action_submit_btn" ).click( function(){
        var oldPasswd = $( "#user_old_passwd" ).val().trim();
        var newPasswd = $( "#user_new_passwd" ).val().trim();
        var newPasswdAgain = $( "#user_new_passwd_again" ).val().trim();
        if( oldPasswd.length <= 0 ){
            return tofuUtils.modalMsg( { content: "原密码不能为空" } );
        }
        if( newPasswd.length <= 0 ){
            return tofuUtils.modalMsg( { content: "新密码不能为空" } );
        }
        if( newPasswd != newPasswdAgain ){
            return tofuUtils.modalMsg( { content: "确认新密码有误，请检查后重试" } );
        }
        var param = "old_passwd=" + oldPasswd + "&new_passwd=" + newPasswd;
        tofuUtils.sendRequest( {
            url: "/ajax/changepasswd",
            args: param,
            onSuccess: function(){
                $( "#adminNavigationLogout" ).trigger( "click" );
            }
        } );
    } );
} );