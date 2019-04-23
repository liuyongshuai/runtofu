/**
 * 插入文档模板
 * @author  Liu Yongshuai<liuyongshuai@sogou-inc.com>
 * @date    2017-03-10 0:35
 */
(function(){

    var factory = function( exports ){

        var pluginName = "doc-template-dialog";
        var $ = jQuery;

        exports.fn.DocTemplateDialog = function(){

            var _this = this;
            var cm = this.cm;
            var lang = this.lang;
            var editor = this.editor;
            var settings = this.settings;
            var path = settings.pluginPath + pluginName + "/";
            var cursor = cm.getCursor();
            var selection = cm.getSelection();
            var classPrefix = this.classPrefix;
            var dialogName = classPrefix + pluginName, dialog;
            var templateList = {};
            var optionList = "";

            cm.focus();

            if( editor.find( "." + dialogName ).length < 1 ){
                $.getJSON( path + pluginName.replace( "-dialog", "" ) + ".json", function( json ){
                    $( json ).each( function( index, data ){
                        optionList += "<option value='" + data.value + "'>" + data.desc + "</option>";
                        templateList[data.value] = data.template;
                    } );
                    var dialogHTML = "<div class=\"" + classPrefix + "form\">" +
                                     "<select class=\"form-control\" style=\"width:250px;\" data-ldapp-doc-template>" + optionList + "</select>";
                    dialog = _this.createDialog( {
                        title: "选择要套用的模板",
                        width: 320,
                        height: 150,
                        content: dialogHTML,
                        mask: settings.dialogShowMask,
                        drag: settings.dialogDraggable,
                        lockScreen: settings.dialogLockScreen,
                        maskStyle: {
                            opacity: settings.dialogMaskOpacity,
                            backgroundColor: settings.dialogMaskBgColor
                        },
                        buttons: {
                            enter: [
                                lang.buttons.enter, function(){
                                    var tem = this.find( "[data-ldapp-doc-template]" ).val();
                                    var str = templateList[tem];
                                    cm.replaceSelection( str );
                                    this.hide().lockScreen( false ).hideMask();
                                    return false;
                                }
                            ],
                            cancel: [
                                lang.buttons.cancel, function(){
                                    this.hide().lockScreen( false ).hideMask();
                                    return false;
                                }
                            ]
                        }
                    } );
                } );
            }
        };
    };

    // CommonJS/Node.js
    if( typeof require === "function" && typeof exports === "object" && typeof module === "object" ){
        module.exports = factory;
    }
    else if( typeof define === "function" )  // AMD/CMD/Sea.js
    {
        if( define.amd ){ // for Require.js
            define( ["editormd"], function( editormd ){
                factory( editormd );
            } );

        }
        else{ // for Sea.js
            define( function( require ){
                var editormd = require( "./../../editormd" );
                factory( editormd );
            } );
        }
    }
    else{
        factory( window.editormd );
    }

})();
