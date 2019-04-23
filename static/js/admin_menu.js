/**
 * @author  Liu Yongshuai
 * @date    2016-11-24 23:27
 */
var iconList="glyphicon-asterisk glyphicon-plus glyphicon-eur glyphicon-minus glyphicon-cloud glyphicon-envelope glyphicon-pencil glyphicon-glass glyphicon-music glyphicon-search glyphicon-heart glyphicon-star glyphicon-star-empty glyphicon-user glyphicon-film glyphicon-th-large glyphicon-th glyphicon-th-list glyphicon-ok glyphicon-remove glyphicon-zoom-in glyphicon-zoom-out glyphicon-off glyphicon-signal glyphicon-cog glyphicon-trash glyphicon-home glyphicon-file glyphicon-time glyphicon-road glyphicon-download-alt glyphicon-download glyphicon-upload glyphicon-inbox glyphicon-play-circle glyphicon-repeat glyphicon-refresh glyphicon-list-alt glyphicon-lock glyphicon-flag glyphicon-headphones glyphicon-volume-off glyphicon-volume-down glyphicon-volume-up glyphicon-qrcode glyphicon-barcode glyphicon-tag glyphicon-tags glyphicon-book glyphicon-bookmark glyphicon-print glyphicon-camera glyphicon-font glyphicon-bold glyphicon-italic glyphicon-text-height glyphicon-text-width glyphicon-align-left glyphicon-align-center glyphicon-align-right glyphicon-align-justify glyphicon-list glyphicon-indent-left glyphicon-indent-right glyphicon-facetime-video glyphicon-picture glyphicon-map-marker glyphicon-adjust glyphicon-tint glyphicon-edit glyphicon-share glyphicon-check glyphicon-move glyphicon-step-backward glyphicon-fast-backward glyphicon-backward glyphicon-play glyphicon-pause glyphicon-stop glyphicon-forward glyphicon-fast-forward glyphicon-step-forward glyphicon-eject glyphicon-plus-sign glyphicon-minus-sign glyphicon-remove-sign glyphicon-ok-sign glyphicon-question-sign glyphicon-info-sign glyphicon-screenshot glyphicon-remove-circle glyphicon-ok-circle glyphicon-ban-circle glyphicon-arrow-left glyphicon-arrow-right glyphicon-arrow-up glyphicon-arrow-down glyphicon-share-alt glyphicon-resize-full glyphicon-resize-small glyphicon-exclamation-sign glyphicon-gift glyphicon-leaf glyphicon-fire glyphicon-eye-open glyphicon-eye-close glyphicon-warning-sign glyphicon-plane glyphicon-calendar glyphicon-random glyphicon-comment glyphicon-magnet glyphicon-chevron-up glyphicon-chevron-down glyphicon-chevron-left glyphicon-chevron-right glyphicon-retweet glyphicon-shopping-cart glyphicon-folder-close glyphicon-folder-open glyphicon-resize-vertical glyphicon-resize-horizontal glyphicon-hdd glyphicon-bullhorn glyphicon-bell glyphicon-certificate glyphicon-thumbs-up glyphicon-thumbs-down glyphicon-hand-right glyphicon-hand-left glyphicon-hand-up glyphicon-hand-down glyphicon-circle-arrow-right glyphicon-circle-arrow-left glyphicon-circle-arrow-up glyphicon-circle-arrow-down glyphicon-globe glyphicon-wrench glyphicon-tasks glyphicon-filter glyphicon-briefcase glyphicon-fullscreen glyphicon-dashboard glyphicon-paperclip glyphicon-heart-empty glyphicon-link glyphicon-phone glyphicon-pushpin glyphicon-usd glyphicon-gbp glyphicon-sort glyphicon-sort-by-alphabet glyphicon-sort-by-alphabet-alt glyphicon-sort-by-order glyphicon-sort-by-order-alt glyphicon-sort-by-attributes glyphicon-sort-by-attributes-alt glyphicon-unchecked glyphicon-expand glyphicon-collapse-down glyphicon-collapse-up glyphicon-log-in glyphicon-flash glyphicon-log-out glyphicon-new-window glyphicon-record glyphicon-save glyphicon-open glyphicon-saved glyphicon-import glyphicon-export glyphicon-send glyphicon-floppy-disk glyphicon-floppy-saved glyphicon-floppy-remove glyphicon-floppy-save glyphicon-floppy-open glyphicon-credit-card glyphicon-transfer glyphicon-cutlery glyphicon-header glyphicon-compressed glyphicon-earphone glyphicon-phone-alt glyphicon-tower glyphicon-stats glyphicon-sd-video glyphicon-hd-video glyphicon-subtitles glyphicon-sound-stereo glyphicon-sound-dolby glyphicon-sound-5-1 glyphicon-sound-6-1 glyphicon-sound-7-1 glyphicon-copyright-mark glyphicon-registration-mark glyphicon-cloud-download glyphicon-cloud-upload glyphicon-tree-conifer glyphicon-tree-deciduous glyphicon-cd glyphicon-save-file glyphicon-open-file glyphicon-level-up glyphicon-copy glyphicon-paste glyphicon-alert glyphicon-equalizer glyphicon-king glyphicon-queen glyphicon-pawn glyphicon-bishop glyphicon-knight glyphicon-baby-formula glyphicon-tent glyphicon-blackboard glyphicon-bed glyphicon-apple glyphicon-erase glyphicon-hourglass glyphicon-lamp glyphicon-duplicate glyphicon-piggy-bank glyphicon-scissors glyphicon-bitcoin glyphicon-yen glyphicon-scale glyphicon-ice-lolly glyphicon-ice-lolly-tasted glyphicon-education glyphicon-option-horizontal glyphicon-option-vertical glyphicon-menu-hamburger glyphicon-modal-window glyphicon-oil glyphicon-grain glyphicon-sunglasses glyphicon-text-size glyphicon-text-color glyphicon-text-background glyphicon-object-align-top glyphicon-object-align-bottom glyphicon-object-align-horizontal glyphicon-object-align-left glyphicon-object-align-vertical glyphicon-object-align-right glyphicon-triangle-right glyphicon-triangle-left glyphicon-triangle-bottom glyphicon-triangle-top glyphicon-console glyphicon-superscript glyphicon-subscript glyphicon-menu-left glyphicon-menu-right glyphicon-menu-down glyphicon-menu-up".split(" ");
$( function(){
    /************删除菜单项************/
    $( "a[action='delmenu']" ).click( function(){
        var trObj = $( this ).parents( "tr[menu_id]" );
        var menu_id = parseInt( trObj.attr( "menu_id" ) );
        tofuUtils.sendRequest( {
            url: "/ajax/system/delMenu",
            args: "menu_id=" + menu_id,
            onSuccess: function(){
                location.reload(true);
            }
        } );
    } );
    var menuParentSelect = $( "#menu_parent_menu" );
    /***********初始化父菜单选择下拉框的数据***********/
    var initParentSelect = function(){
        var options = '<option value="0">&nbsp;</option>';
        menuParentSelect.empty();
        $.each( $.parseJSON( menuList ), function( index, info ){
            if( parseInt( info.menuInfo.menu_id ) > 0 ){
                options += '<option value="' + info.menuInfo.menu_id + '">' + info.menuInfo.menu_name + '</option>';
            }
        } );
        menuParentSelect.append( options );
    };
    /***********初始化颜色选择器***********/
    var initIconColor = function( color ){
        $( "#menu_icon_color" ).ColorPickerSliders( {
            placement: 'right',
            flat: false,
            color: color,
            title: "选择菜单图标的颜色",
            sliders: false,
            swatches: false,
            hsvpanel: true,
            previewformat: 'hex'
        } );
    };
    /************编辑菜单信息************/
    $( "a[action=\"editmenu\"]" ).click( function(){
        $( "#menu_name" ).val( "" );
        $( "#menu_path" ).val( "" );
        initParentSelect();
        $( "#menu_icon_name" ).val( "" );
        $( "#menu_icon_color" ).val( "" ).css( "background-color", "white" );
        var trObj = $( this ).parents( "tr[menu_id]" );
        var menu_id = 0;
        if( trObj.length > 0 ){
            menu_id = parseInt( trObj.attr( "menu_id" ) );
        }
        $( "#menu_menu_id" ).val( menu_id );
        tofuUtils.sendRequest( {
            url: "/ajax/system/getMenuInfo",
            args: "menu_id=" + menu_id,
            onSuccess: function( menuInfo ){
                var originColor = "";
                if( !$.isEmptyObject( menuInfo ) ){
                    $( "#menu_name" ).val( menuInfo.menu_name );
                    $( "#menu_path" ).val( menuInfo.menu_path );
                    $( "#menu_icon_name" ).val( menuInfo.icon_name ).parent().find( "span.glyphicon" ).addClass( menuInfo.icon_name );
                    originColor = menuInfo.icon_color;
                    if( parseInt( menuInfo.parent_menu_id ) > 0 ){
                        $( "#menu_parent_menu option[value=0]" ).remove();
                        $( "#menu_parent_menu option[value=\"" + menuInfo.parent_menu_id + "\"]" ).attr( "selected", "true" );
                    }
                    else{
                        menuParentSelect.empty();
                        menuParentSelect.attr( "disabled", true );
                    }
                }
                menuParentSelect.selectpicker( "destroy" ).selectpicker( { maxOptions: 1 } );
                initIconColor( originColor );
                $( "#menu_edit_dialog" ).modal( 'show' );
            }
        } );
    } );
    /************提交菜单信息************/
    $( "#edit_menu_submit_btn" ).click( function(){
        var menu_id = $( "#menu_menu_id" ).val();
        var name = $( "#menu_name" ).val();
        var path = $( "#menu_path" ).val();
        var icon = $( "#menu_icon_name" ).val();
        var icon_color = $( "#menu_icon_color" ).val();
        var parent_menu_id = $( "#menu_parent_menu" ).val();
        var args = "menu_id=" + menu_id + "&name=" + encodeURIComponent( name ) + "&path=" + encodeURIComponent( path );
        args += "&icon=" + icon + "&icon_color=" + icon_color;
        args += "&parent_menu_id=" + parent_menu_id;
        tofuUtils.sendRequest( {
            url: "/ajax/system/modifyMenuInfo",
            args: args,
            onSuccess: function(){
                location.reload( true );
            }
        } );
    } );
    /************展开子菜单************/
    $( "button[action='showsubmenu']" ).click( function(){
        var menu_id = parseInt( $( this ).parents( "tr[menu_id]" ).attr( "menu_id" ) );
        $( "tr[parent_menu_id][parent_menu_id!=\"" + menu_id + "\"]" ).hide();
        var sub = $( "tr[parent_menu_id=\"" + menu_id + "\"]" );
        if( sub.is( ":visible" ) ){
            sub.hide();
            location.hash = "";
        }
        else{
            sub.show();
            location.hash = "#" + menu_id;
        }
    } );
    /************处理URL后跟的锚点信息************/
    var menu_id = parseInt( location.hash.substring( 1 ) );
    if( menu_id > 0 ){
        $( "tr[parent_menu_id=\"" + menu_id + "\"]" ).show();
    }
    /************图标选择器************/
    $( "#menu_icon_name" ).popover( {
        html: true,
        placement: "right",
        title: "选择菜单图标",
        container: "body",
        content: function(){
            var icons = "<div class='icon-picker'>";
            $.each( iconList, function( index, val ){
                icons += '<span class="glyphicon ' + val + '"></span>';
            } );
            icons += "</div>";
            return icons;
        }
    } ).on( 'shown.bs.popover', function(){
        $( "div.icon-picker > span" ).tooltip( {
            html: true,
            placement: "auto",
            title: function(){
                var cs = "enlarge";
                var cls = $( this ).attr( "class" );
                var t = "<span class='" + cls + " " + cs + "'></span>";
                return t;
            }
        } );
    } );
    $( "body" ).delegate( "div.icon-picker > span", "click", function(){
        var cls = $( this ).attr( "class" );
        cls = cls.split( " " );
        cls = cls[1];
        var oldCls = $( "#menu_icon_name" ).val();
        $( "#menu_icon_name" ).val( cls ).parent().find( "span.glyphicon" ).removeClass( oldCls ).addClass( cls );
        $( '#menu_icon_name' ).popover( 'hide' )
    } );
    /************上调话题顺序************/
    $( "a[action='upmenu']" ).click( function(){
        var trObj = $( this ).parents( "tr[menu_id]" );
        var prev = trObj.prev();
        var menu_id = parseInt( trObj.attr( "menu_id" ) );
        tofuUtils.sendRequest( {
            url: "/ajax/system/upMenu",
            args: "menu_id=" + menu_id,
            onSuccess: function(){
                location.reload( true );
            }
        } );
    } );
    /************下调话题顺序************/
    $( "a[action='downmenu']" ).click( function(){
        var trObj = $( this ).parents( "tr[menu_id]" );
        var next = trObj.next();
        var menu_id = parseInt( trObj.attr( "menu_id" ) );
        tofuUtils.sendRequest( {
            url: "/ajax/system/downMenu",
            args: "menu_id=" + menu_id,
            onSuccess: function(){
                location.reload( true );
            }
        } );
    } );
} );
