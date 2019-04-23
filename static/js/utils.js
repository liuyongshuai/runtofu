$( function(){
    /************显示回到顶部的浮动图标************/
    $( document ).scroll( function(){
        if( parseInt( $( document ).height() ) > parseInt( $( window ).height() ) && parseInt( $( document ).scrollTop() ) > 100 ){
            $( "div.goto-top" ).show();
        }
        else{
            $( "div.goto-top" ).hide();
        }
    } );
    /************显示工具条提示************/
    $( "#global-main-content" ).tooltip( {
        trigger: "hover",
        html: true,
        selector: "[title]"
    } );
    $( "*[data-toggle=\"tooltip\"]" ).tooltip( {
        trigger: "hover",
        html: true
    } );
} );


/**********通用工具**********/
var tofuUtils = (function(){
    return {
        /**
         * *******************************************************************
         * 产生随机字符串
         *
         * @param    length    随机串长度
         * @param    type    类型：1仅数字、2仅字母、3字母数字混合
         * *******************************************************************
         */
        random: function( length, type ){
            type = parseInt( type );
            var char1 = [
                '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'
            ];
            var char2 = [
                'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
                'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
                'u', 'v', 'w', 'x', 'y', 'z'
            ];
            var char3 = [
                '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
                'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
                'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
                'u', 'v', 'w', 'x', 'y', 'z'
            ];
            var res = "";
            var data;
            if( type === 1 ){
                data = char1;
            }
            else if( type === 2 ){
                data = char2;
            }
            else{
                data = char3;
            }
            for( var i = 0; i < length; i++ ){
                var id = Math.floor( Math.random() * data.length );
                res += data[id];
            }
            return res;
        },

        /**
         * *******************************************************************
         * 生成伪UUID，格式：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx（8-4-4-4-12）
         * *******************************************************************
         */
        juuid: function(){
            return this.random( 8, 3 ) + "-" + this.random( 4, 3 ) + "-" + this.random( 4, 3 ) + "-" + this.random( 4, 3 ) + "-" + this.random( 12, 3 );
        },

        /**
         * *******************************************************************
         * 用模态弹框提示错误信息，只有一个确认关闭框
         * *******************************************************************
         */
        modalMsg: function( options ){
            var defaultShowModalErrorsOption = {
                content: "", /*默认内容*/
                type: "danger", /*类型：danger/warning/info/success/primary*/
                title: "出错啦！", /*默认标题*/
                onConfirm: ""/*点击确定时的回调函数*/
            };
            var _opt = $.extend( {}, defaultShowModalErrorsOption, options );
            var confirmFunc = _opt.onConfirm;
            var t = BootstrapDialog.TYPE_DANGER;
            switch( _opt.type ){
                case "danger":
                    t = BootstrapDialog.TYPE_DANGER;
                    break;
                case "warning":
                    t = BootstrapDialog.TYPE_WARNING;
                    break;
                case "info":
                    t = BootstrapDialog.TYPE_INFO;
                    break;
                case "success":
                    t = BootstrapDialog.TYPE_SUCCESS;
                    break;
                case "primary":
                    t = BootstrapDialog.TYPE_PRIMARY;
                    break;
                case "default":
                    t = BootstrapDialog.TYPE_DEFAULT;
                    break;
            }
            BootstrapDialog.show( {
                type: t,
                title: _opt.title,
                message: _opt.content,
                closable: true,
                size: BootstrapDialog.SIZE_SMALL,
                buttons: [
                    {
                        label: '确定',
                        cssClass: 'btn-' + _opt.type,
                        action: function( dialog ){
                            if( $.isFunction( confirmFunc ) ){
                                confirmFunc( dialog );
                            }
                            dialog.close();
                        }
                    }
                ]
            } );
            return false;
        },

        /**
         * *******************************************************************
         * 发起POST/GET请求
         * *******************************************************************
         */
        sendRequest: function( options ){
            var _this = this;
            var defaultSendPostRequestOptions = {
                url: "", /*要请求的URL*/
                args: "", /*参数：如“a=1&b=2”*/
                async: true, /*是否为异步请求，为false表示同步请求，此时浏览器将被锁住*/
                onSuccess: null, /*请求成功时回调的函数*/
                retRawData: false,//是否保留原始数据
                onError: null, /*出错时的回调函数*/
                onComplete: null/*请求完成时的回调函数*/
            };
            var _opt = $.extend( {}, defaultSendPostRequestOptions, options );
            var succFunc = _opt.onSuccess;
            var errorFunc = _opt.onError;
            var completeFunc = _opt.onComplete;
            var retRawData = _opt.retRawData;
            try{
                $.ajax( {
                    type: "POST",
                    url: _opt.url,
                    data: _opt.args,
                    dataType: "json",
                    async: _opt.async,
                    xhrFields: {
                        withCredentials: true
                    },
                    success: function( ret ){
                        _this.loading.end();
                        if( $.isFunction( completeFunc ) ){
                            completeFunc( ret );
                        }
                        if( retRawData ){
                            if( $.isFunction( succFunc ) ){
                                succFunc( ret );
                            }
                            return;
                        }
                        if( $.isEmptyObject( ret ) ){
                            return _this.modalMsg( { content: ret, type: "danger" } );
                        }
                        if( $.type( ret ) === "string" ){
                            ret = JSON.parse( ret );
                        }
                        var code = parseInt( ret.result.errno );
                        if( code !== 0 ){
                            if( $.isFunction( errorFunc ) ){
                                errorFunc( ret );
                            }
                            else{
                                _this.modalMsg( { content: ret.result.errmsg, type: "danger" } );
                            }
                            return;
                        }
                        if( $.isFunction( succFunc ) ){
                            succFunc( ret.data );
                        }
                    },
                    error: function( XMLHttpRequest, textStatus, errorThrown ){
                        _this.loading.end();
                        var errmsg = "Request.status = " + XMLHttpRequest.status;
                        if( XMLHttpRequest.responseText.length > 0 ){
                            errmsg = "\nRequest.responseText = " + XMLHttpRequest.responseText;
                        }
                        if( errorThrown.length > 0 ){
                            errmsg += "\nerrorThrown = " + errorThrown;
                        }
                        if( textStatus.length > 0 ){
                            errmsg += "\ntextStatus = " + textStatus;
                        }
                        _this.modalMsg( { content: errmsg, type: "danger", title: "请求失败" } );
                    }
                } );
            }
            catch( e ){
                _this.modalMsg( { content: e.message, type: "danger" } );
            }
        },
        /**
         * *******************************************************************
         * 获取指定DOM上的指定属性值列表
         * *******************************************************************
         */
        getAttrs: function( selector, attrName ){
            var ret = [];
            $( selector ).each( function(){
                ret.push( $( this ).attr( attrName ) );
            } );
            return ret;
        },

        /**
         * *******************************************************************
         * 记录日志信息
         * *******************************************************************
         */
        log: function(){
            try{
                if( window.console && window.console.log ){
                    console.log.apply( this, arguments );
                }
            }
            catch( e ){
            }
        },

        /**
         * *******************************************************************
         * 显示加载中的gif动图
         * *******************************************************************
         */
        loading: (function(){
            var LOADING_PIC_START_NUM = Math.floor( 1000 * Math.random() );
            var LOADING_PICLIST = [];
            for( var i = 0; i <= 34; i++ ){
                LOADING_PICLIST.push( STATIC_PREFIX + "/images/loading/" + i + ".gif" );
            }
            return {
                start: function(){
                    this.end();
                    var len = LOADING_PICLIST.length;
                    var num = (LOADING_PIC_START_NUM++) % len;
                    var img = LOADING_PICLIST[num];
                    $( "body" ).append( '<div class="loading"><div class="modal fade in"><div class="modal-dialog"><img src="' + img + '"></div></div><div class="modal-backdrop fade in"></div></div>' );
                },
                end: function(){
                    $( "div.loading" ).remove();
                }
            }
        })(),
        /**
         * *******************************************************************
         * 返回当前时间戳（类PHP中的time()）
         * *******************************************************************
         */
        time: function( micro ){
            var tm = new Date().getTime();
            if( micro ){
                return tm;
            }
            return parseInt( tm / 1000 );
        },

        /**
         * *******************************************************************
         * 将时间戳格式化为指定格式：格式同PHP中的date()函数
         * 参见：http://php.net/manual/zh/function.date.php
         * *******************************************************************
         */
        date: function( format, timestamp ){
            var a, jsdate = ((timestamp) ? new Date( timestamp * 1000 ) : new Date());
            var pad = function( n, c ){
                if( (n = n + "").length < c ){
                    return new Array( ++c - n.length ).join( "0" ) + n;
                }
                else{
                    return n;
                }
            };
            var txt_weekdays = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
            var txt_ordin = { 1: "st", 2: "nd", 3: "rd", 21: "st", 22: "nd", 23: "rd", 31: "st" };
            var txt_months = [
                "",
                "January",
                "February",
                "March",
                "April",
                "May",
                "June",
                "July",
                "August",
                "September",
                "October",
                "November",
                "December"
            ];
            var f = {
                d: function(){
                    return pad( f.j(), 2 );
                },
                D: function(){
                    t = f.l();
                    return t.substr( 0, 3 );
                },
                j: function(){
                    return jsdate.getDate();
                },
                l: function(){
                    return txt_weekdays[f.w()];
                },
                N: function(){
                    return f.w() + 1;
                },
                S: function(){
                    return txt_ordin[f.j()] ? txt_ordin[f.j()] : 'th';
                },
                w: function(){
                    return jsdate.getDay();
                },
                z: function(){
                    return (jsdate - new Date( jsdate.getFullYear() + "/1/1" )) / 864e5 >> 0;
                },
                W: function(){
                    var a = f.z(), b = 364 + f.L() - a;
                    var nd2, nd = (new Date( jsdate.getFullYear() + "/1/1" ).getDay() || 7) - 1;
                    if( b <= 2 && ((jsdate.getDay() || 7) - 1) <= 2 - b ){
                        return 1;
                    }
                    else{
                        if( a <= 2 && nd >= 4 && a >= (6 - nd) ){
                            nd2 = new Date( jsdate.getFullYear() - 1 + "/12/31" );
                            return date( "W", Math.round( nd2.getTime() / 1000 ) );
                        }
                        else{
                            return (1 + (nd <= 3 ? ((a + nd) / 7) : (a - (7 - nd)) / 7) >> 0);
                        }
                    }
                },
                F: function(){
                    return txt_months[f.n()];
                },
                m: function(){
                    return pad( f.n(), 2 );
                },
                M: function(){
                    t = f.F();
                    return t.substr( 0, 3 );
                },
                n: function(){
                    return jsdate.getMonth() + 1;
                },
                t: function(){
                    var n;
                    if( (n = jsdate.getMonth() + 1) == 2 ){
                        return 28 + f.L();
                    }
                    else{
                        if( n & 1 && n < 8 || !(n & 1) && n > 7 ){
                            return 31;
                        }
                        else{
                            return 30;
                        }
                    }
                },
                L: function(){
                    var y = f.Y();
                    return (!(y & 3) && (y % 1e2 || !(y % 4e2))) ? 1 : 0;
                },
                Y: function(){
                    return jsdate.getFullYear();
                },
                y: function(){
                    return (jsdate.getFullYear() + "").slice( 2 );
                },
                a: function(){
                    return jsdate.getHours() > 11 ? "pm" : "am";
                },
                A: function(){
                    return f.a().toUpperCase();
                },
                B: function(){
                    var off = (jsdate.getTimezoneOffset() + 60) * 60;
                    var theSeconds = (jsdate.getHours() * 3600) +
                                     (jsdate.getMinutes() * 60) +
                                     jsdate.getSeconds() + off;
                    var beat = Math.floor( theSeconds / 86.4 );
                    if( beat > 1000 ) beat -= 1000;
                    if( beat < 0 ) beat += 1000;
                    if( (String( beat )).length == 1 ) beat = "00" + beat;
                    if( (String( beat )).length == 2 ) beat = "0" + beat;
                    return beat;
                },
                g: function(){
                    return jsdate.getHours() % 12 || 12;
                },
                G: function(){
                    return jsdate.getHours();
                },
                h: function(){
                    return pad( f.g(), 2 );
                },
                H: function(){
                    return pad( jsdate.getHours(), 2 );
                },
                i: function(){
                    return pad( jsdate.getMinutes(), 2 );
                },
                s: function(){
                    return pad( jsdate.getSeconds(), 2 );
                },
                O: function(){
                    var t = pad( Math.abs( jsdate.getTimezoneOffset() / 60 * 100 ), 4 );
                    if( jsdate.getTimezoneOffset() > 0 ) t = "-" + t;
                    else t = "+" + t;
                    return t;
                },
                P: function(){
                    var O = f.O();
                    return (O.substr( 0, 3 ) + ":" + O.substr( 3, 2 ));
                },
                c: function(){
                    return f.Y() + "-" + f.m() + "-" + f.d() + "T" + f.h() + ":" + f.i() + ":" + f.s() + f.P();
                },
                U: function(){
                    return Math.round( jsdate.getTime() / 1000 );
                }
            };
            return format.replace( /[\\]?([a-zA-Z])/g, function( t, s ){
                if( t != s ){
                    ret = s;
                }
                else if( f[s] ){
                    ret = f[s]();
                }
                else{
                    ret = s;
                }
                return ret;
            } );
        },

        /**
         * *******************************************************************
         * 高仿PHP中的mktime，顺序：小时、分钟、秒、月、日、年
         * *******************************************************************
         */
        mktime: function(){
            var d = new Date();
            var r = arguments;
            var i = 0;
            var e = ['Hours', 'Minutes', 'Seconds', 'Month', 'Date', 'FullYear'];
            for( i = 0; i < e.length; i++ ){
                if( typeof r[i] === 'undefined' ){
                    r[i] = d['get' + e[i]]();
                    r[i] += (i === 3);
                }
                else{
                    r[i] = parseInt( r[i], 10 );
                    if( isNaN( r[i] ) ){
                        return false;
                    }
                }
            }
            r[5] += (r[5] >= 0 ? (r[5] <= 69 ? 2e3 : (r[5] <= 100 ? 1900 : 0)) : 0);
            d.setFullYear( r[5], r[3] - 1, r[4] );
            d.setHours( r[0], r[1], r[2] );
            var time = d.getTime();
            return (time / 1e3 >> 0) - (time < 0);
        },

        /**
         * *******************************************************************
         * 高仿PHP中的strtotime
         * @param text
         * @param now
         * *******************************************************************
         */
        strtotime: function( text, now ){
            text = this.strval( text );
            var parsed;
            var match;
            var today;
            var year;
            var date;
            var days;
            var ranges;
            var len;
            var times;
            var regex;
            var i;
            var fail = false;
            if( !text ){
                return fail
            }
            text = text.replace( /^\s+|\s+$/g, '' ).replace( /\s{2,}/g, ' ' ).replace( /[\t\r\n]/g, '' ).toLowerCase();
            var pattern = new RegExp( [
                '^(\\d{1,4})',
                '([\\-\\.\\/:])',
                '(\\d{1,2})',
                '([\\-\\.\\/:])',
                '(\\d{1,4})',
                '(?:\\s(\\d{1,2}):(\\d{2})?:?(\\d{2})?)?',
                '(?:\\s([A-Z]+)?)?$'
            ].join( '' ) );
            match = text.match( pattern );

            if( match && match[2] === match[4] ){
                if( match[1] > 1901 ){
                    switch( match[2] ){
                        case '-':
                            // YYYY-M-D
                            if( match[3] > 12 || match[5] > 31 ){
                                return fail;
                            }
                            return new Date( match[1], parseInt( match[3], 10 ) - 1, match[5], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                        case '.':
                            // YYYY.M.D is not parsed by strtotime()
                            return fail;
                        case '/':
                            // YYYY/M/D
                            if( match[3] > 12 || match[5] > 31 ){
                                return fail;
                            }
                            return new Date( match[1], parseInt( match[3], 10 ) - 1, match[5], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                    }
                }
                else if( match[5] > 1901 ){
                    switch( match[2] ){
                        case '-':
                            // D-M-YYYY
                            if( match[3] > 12 || match[1] > 31 ){
                                return fail;
                            }
                            return new Date( match[5], parseInt( match[3], 10 ) - 1, match[1], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                        case '.':
                            // D.M.YYYY
                            if( match[3] > 12 || match[1] > 31 ){
                                return fail;
                            }
                            return new Date( match[5], parseInt( match[3], 10 ) - 1, match[1], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                        case '/':
                            // M/D/YYYY
                            if( match[1] > 12 || match[3] > 31 ){
                                return fail;
                            }
                            return new Date( match[5], parseInt( match[1], 10 ) - 1, match[3], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                    }
                }
                else{
                    switch( match[2] ){
                        case '-':
                            // YY-M-D
                            if( match[3] > 12 || match[5] > 31 || (match[1] < 70 && match[1] > 38) ){
                                return fail;
                            }
                            year = match[1] >= 0 && match[1] <= 38 ? +match[1] + 2000 : match[1];
                            return new Date( year, parseInt( match[3], 10 ) - 1, match[5], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                        case '.':
                            // D.M.YY or H.MM.SS
                            if( match[5] >= 70 ){
                                // D.M.YY
                                if( match[3] > 12 || match[1] > 31 ){
                                    return fail;
                                }
                                return new Date( match[5], parseInt( match[3], 10 ) - 1, match[1], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                            }
                            if( match[5] < 60 && !match[6] ){
                                // H.MM.SS
                                if( match[1] > 23 || match[3] > 59 ){
                                    return fail;
                                }
                                today = new Date();
                                return new Date( today.getFullYear(), today.getMonth(), today.getDate(), match[1] || 0, match[3] || 0, match[5] || 0, match[9] || 0 ) / 1000;
                            }
                            // invalid format, cannot be parsed
                            return fail;
                        case '/':
                            // M/D/YY
                            if( match[1] > 12 || match[3] > 31 || (match[5] < 70 && match[5] > 38) ){
                                return fail;
                            }
                            year = match[5] >= 0 && match[5] <= 38 ? +match[5] + 2000 : match[5];
                            return new Date( year, parseInt( match[1], 10 ) - 1, match[3], match[6] || 0, match[7] || 0, match[8] || 0, match[9] || 0 ) / 1000;
                        case ':':
                            // HH:MM:SS
                            if( match[1] > 23 || match[3] > 59 || match[5] > 59 ){
                                return fail;
                            }
                            today = new Date();
                            return new Date( today.getFullYear(), today.getMonth(), today.getDate(), match[1] || 0, match[3] || 0, match[5] || 0 ) / 1000;
                    }
                }
            }

            if( text === 'now' ){
                return now === null || isNaN( now ) ? new Date().getTime() / 1000 | 0 : now | 0;
            }
            if( !isNaN( parsed = Date.parse( text ) ) ){
                return parsed / 1000 | 0;
            }
            pattern = new RegExp( [
                '^([0-9]{4}-[0-9]{2}-[0-9]{2})',
                '[ t]',
                '([0-9]{2}:[0-9]{2}:[0-9]{2}(\\.[0-9]+)?)',
                '([\\+-][0-9]{2}(:[0-9]{2})?|z)'
            ].join( '' ) );
            match = text.match( pattern );
            if( match ){
                if( match[4] === 'z' ){
                    match[4] = 'Z';
                }
                else if( match[4].match( /^([+-][0-9]{2})$/ ) ){
                    match[4] = match[4] + ':00';
                }
                if( !isNaN( parsed = Date.parse( match[1] + 'T' + match[2] + match[4] ) ) ){
                    return parsed / 1000 | 0;
                }
            }

            date = now ? new Date( now * 1000 ) : new Date();
            days = {
                'sun': 0,
                'mon': 1,
                'tue': 2,
                'wed': 3,
                'thu': 4,
                'fri': 5,
                'sat': 6
            };
            ranges = {
                'yea': 'FullYear',
                'mon': 'Month',
                'day': 'Date',
                'hou': 'Hours',
                'min': 'Minutes',
                'sec': 'Seconds'
            };

            function lastNext( type, range, modifier ){
                var diff;
                var day = days[range];
                if( typeof day !== 'undefined' ){
                    diff = day - date.getDay();
                    if( diff === 0 ){
                        diff = 7 * modifier;
                    }
                    else if( diff > 0 && type === 'last' ){
                        diff -= 7;
                    }
                    else if( diff < 0 && type === 'next' ){
                        diff += 7;
                    }
                    date.setDate( date.getDate() + diff )
                }
            }

            function process( val ){
                var splt = val.split( ' ' );
                var type = splt[0];
                var range = splt[1].substring( 0, 3 );
                var typeIsNumber = /\d+/.test( type );
                var ago = splt[2] === 'ago';
                var num = (type === 'last' ? -1 : 1) * (ago ? -1 : 1);
                if( typeIsNumber ){
                    num *= parseInt( type, 10 );
                }
                if( ranges.hasOwnProperty( range ) && !splt[1].match( /^mon(day|\.)?$/i ) ){
                    return date['set' + ranges[range]]( date['get' + ranges[range]]() + num );
                }
                if( range === 'wee' ){
                    return date.setDate( date.getDate() + (num * 7) );
                }
                if( type === 'next' || type === 'last' ){
                    lastNext( type, range, num );
                }
                else if( !typeIsNumber ){
                    return false;
                }
                return true;
            }

            times = '(years?|months?|weeks?|days?|hours?|minutes?|min|seconds?|sec|sunday|sun\\.?|monday|mon\\.?|tuesday|tue\\.?|wednesday|wed\\.?|thursday|thu\\.?|friday|fri\\.?|saturday|sat\\.?)';
            regex = '([+-]?\\d+\\s' + times + '|' + '(last|next)\\s' + times + ')(\\sago)?';
            match = text.match( new RegExp( regex, 'gi' ) );
            if( !match ){
                return fail;
            }
            for( i = 0, len = match.length; i < len; i++ ){
                if( !process( match[i] ) ){
                    return fail;
                }
            }

            return (date.getTime() / 1000);
        },

        /**
         * *******************************************************************
         * 高仿PHP中的empty
         * @param mixedVar
         * *******************************************************************
         */
        empty: function( mixedVar ){
            var undef;
            var key;
            var i;
            var len;
            var emptyValues = [undef, null, false, 0, '', '0'];
            for( i = 0, len = emptyValues.length; i < len; i++ ){
                if( mixedVar === emptyValues[i] ){
                    return true;
                }
            }
            if( typeof mixedVar === 'object' ){
                for( key in mixedVar ){
                    if( mixedVar.hasOwnProperty( key ) ){
                        return false;
                    }
                }
                return true;
            }
            return false;
        },

        /**
         * *******************************************************************
         * 高仿PHP中的intval
         * @param mixedVar
         * @param base
         * *******************************************************************
         */
        intval: function( mixedVar, base ){
            var tmp;
            var type = typeof mixedVar;
            if( type === 'boolean' ){
                return +mixedVar;
            }
            else if( type === 'string' ){
                mixedVar = this.trim( mixedVar );
                tmp = parseInt( mixedVar, base || 10 );
                return (isNaN( tmp ) || !isFinite( tmp )) ? 0 : tmp
            }
            else if( type === 'number' && isFinite( mixedVar ) ){
                return mixedVar | 0;
            }
            else{
                return 0;
            }
        },

        /**
         * *******************************************************************
         * 高仿PHP中setcookie
         * @param name
         * @param value
         * @param expires
         * @param path
         * @param domain
         * @param secure
         * *******************************************************************
         */
        setcookie: function( name, value, expires, path, domain, secure ){
            return this.setrawcookie( name, encodeURIComponent( value ), expires, path, domain, secure );
        },

        /**
         * *******************************************************************
         * 高仿PHP中的setrawcookie
         * @param name
         * @param value
         * @param expires
         * @param path
         * @param domain
         * @param secure
         * *******************************************************************
         */
        setrawcookie: function( name, value, expires, path, domain, secure ){
            if( typeof window === 'undefined' ){
                return true;
            }
            if( typeof expires === 'string' && (/^\d+$/).test( expires ) ){
                expires = parseInt( expires, 10 );
            }
            if( expires instanceof Date ){
                expires = expires.toUTCString();
            }
            else if( typeof expires === 'number' ){
                expires = (new Date( expires * 1e3 )).toUTCString();
            }
            var r = [name + '=' + value];
            var i = '';
            var s = {
                expires: expires,
                path: path,
                domain: domain
            };
            for( i in s ){
                if( s.hasOwnProperty( i ) ){
                    s[i] && r.push( i + '=' + s[i] );
                }
            }
            if( secure ){
                r.push( 'secure' );
            }
            window.document.cookie = r.join( ';' );
            return true;
        },

        /**
         * *******************************************************************
         * 高仿PHP中的trim函数
         * @param str
         * @param charlist
         * *******************************************************************
         */
        trim: function( str, charlist ){
            str = this.strval( str );
            var whitespace = [
                ' ', '\n', '\r', '\t', '\f', '\x0b', '\xa0', '\u2000', '\u2001', '\u2002', '\u2003', '\u2004',
                '\u2005', '\u2006', '\u2007', '\u2008', '\u2009', '\u200a', '\u200b', '\u2028', '\u2029',
                '\u3000'
            ].join( '' );
            var l = 0;
            var i = 0;
            str += '';
            if( charlist ){
                whitespace = (charlist + '').replace( /([[\]().?/*{}+$^:])/g, '$1' );
            }
            l = str.length;
            for( i = 0; i < l; i++ ){
                if( whitespace.indexOf( str.charAt( i ) ) === -1 ){
                    str = str.substring( i );
                    break;
                }
            }
            l = str.length;
            for( i = l - 1; i >= 0; i-- ){
                if( whitespace.indexOf( str.charAt( i ) ) === -1 ){
                    str = str.substring( 0, i + 1 );
                    break;
                }
            }
            return whitespace.indexOf( str.charAt( 0 ) ) === -1 ? str : '';
        },

        /**
         * *******************************************************************
         * 高仿PHP中的rtrim
         * @param str
         * @param charlist
         * *******************************************************************
         */
        rtrim: function( str, charlist ){
            str = this.strval( str );
            charlist = !charlist ? ' \\s\u00A0' : (charlist + '').replace( /([[\]().?/*{}+$^:])/g, '\\$1' );
            var re = new RegExp( '[' + charlist + ']+$', 'g' );
            return (str + '').replace( re, '' );
        },

        /**
         * *******************************************************************
         * 高仿PHP中的ltrim
         * @param str
         * @param charlist
         * *******************************************************************
         */
        ltrim: function( str, charlist ){
            str = this.strval( str );
            charlist = !charlist ? ' \\s\u00A0' : (charlist + '').replace( /([[\]().?/*{}+$^:])/g, '$1' );
            var re = new RegExp( '^[' + charlist + ']+', 'g' );
            return (str + '').replace( re, '' );
        },

        /**
         * *******************************************************************
         * 是否全是数字
         * @param str
         * *******************************************************************
         */
        is_all_numeric: function( str ){
            str = this.strval( str );
            if( $.type( str ) !== "string" ){
                return false;
            }
            var trim = this.trim( str, '0123456789' );
            if( str.length <= 0 || trim.length > 0 ){
                return false;
            }
            return true;
        },

        /**
         * *******************************************************************
         * 是否是合法的URL
         * @param url
         * *******************************************************************
         */
        is_valid_url: function( url ){
            if( $.type( url ) !== "string" ){
                return false;
            }
            return /^http[s]?:\/\/([\w-]+\.)+[\w-]+([\w-./?%&=]*)?$/i.test( url );
        },
        /**
         * *******************************************************************
         * 转为全角字符
         * @param str
         * @returns {string}
         * *******************************************************************
         */
        to_DBC: function( str ){
            str = this.strval( str );
            var tmp = "";
            for( var i = 0; i < str.length; i++ ){
                if( str.charCodeAt( i ) === 32 ){
                    tmp = tmp + String.fromCharCode( 12288 );
                }
                else if( str.charCodeAt( i ) < 127 ){
                    tmp = tmp + String.fromCharCode( str.charCodeAt( i ) + 65248 );
                }
                else{
                    tmp += String.fromCharCode( str.charCodeAt( i ) );
                }
            }
            return tmp;
        },
        /**
         * *******************************************************************
         * 转为半角字符
         * @param str
         * @returns {string}
         * *******************************************************************
         */
        to_SBC: function( str ){
            str = this.strval( str );
            var result = "";
            for( var i = 0; i < str.length; i++ ){
                var code = str.charCodeAt( i );
                if( code >= 65281 && code <= 65373 ){
                    var d = str.charCodeAt( i ) - 65248;
                    result += String.fromCharCode( d );
                }
                else if( code === 12288 ){
                    var d = str.charCodeAt( i ) - 12288 + 32;
                    result += String.fromCharCode( d );
                }
                else{
                    result += v.charAt( i );
                }
            }
            return result;
        },
        /**
         * *******************************************************************
         * 是否存在中文
         * @param str
         * @returns {string}
         * *******************************************************************
         */
        is_exist_chinese: function( str ){
            str = this.strval( str );
            return !/^[\x00-\xff]*$/.test( str );
        },
        /**
         * *******************************************************************
         * 将指定的值转为字符串
         * @param value
         * *******************************************************************
         */
        strval: function( value ){
            var type = typeof value;
            switch( type ){
                case 'boolean':
                    return value ? '1' : '';
                case 'string':
                    return value;
                case 'number':
                    if( isNaN( value ) ){
                        return 'NAN';
                    }
                    if( !isFinite( value ) ){
                        return (value < 0 ? '-' : '') + 'INF';
                    }
                    return value + '';
                case 'undefined':
                    return '';
                case 'object':
                    if( Array.isArray( value ) || value !== null ){
                        return JSON.stringify( value );
                    }
                    return '';
                case 'function':
                default:
                    throw new Error( 'Unsupported value type' );
            }
        },

        /**
         * *******************************************************************
         * 临时存储在本地
         * *******************************************************************
         */
        store: (function(){
            var win = (typeof window != 'undefined' ? window : global),
                doc = win.document,
                localStorageName = 'localStorage',
                scriptTag = 'script',
                storage;

            function isLocalStorageNameSupported(){
                try{
                    return (localStorageName in win && win[localStorageName])
                }
                catch( err ){
                    return false
                }
            }

            var ret = {
                set: function( key, value ){
                },
                get: function( key, defaultVal ){
                },
                has: function( key ){
                    return this.get( key ) !== undefined;
                },
                remove: function( key ){
                },
                clear: function(){
                },
                forEach: function(){
                },
                getAll: function(){
                    var ret = {};
                    this.forEach( function( key, val ){
                        ret[key] = val;
                    } );
                    return ret;
                },
                serialize: function( value ){
                    return JSON.stringify( value );
                },
                deserialize: function( value ){
                    if( typeof value != 'string' ){
                        return undefined
                    }
                    try{
                        return JSON.parse( value );
                    }
                    catch( e ){
                        return value || undefined;
                    }
                }
            };
            if( isLocalStorageNameSupported() ){
                storage = win[localStorageName];
                ret.set = function( key, val ){
                    if( val === undefined ){
                        return ret.remove( key )
                    }
                    storage.setItem( key, ret.serialize( val ) );
                    return val
                };
                ret.get = function( key, defaultVal ){
                    var val = ret.deserialize( storage.getItem( key ) );
                    return (val === undefined ? defaultVal : val)
                };
                ret.remove = function( key ){
                    storage.removeItem( key )
                };
                ret.clear = function(){
                    storage.clear()
                };
                ret.forEach = function( callback ){
                    if( $.isFunction( callback ) ){
                        for( var i = 0; i < storage.length; i++ ){
                            var key = storage.key( i );
                            callback( key, ret.get( key ) );
                        }
                    }
                };
            }
            else if( doc && doc.documentElement.addBehavior ){
                var storageOwner,
                    storageContainer;
                try{
                    storageContainer = new ActiveXObject( 'htmlfile' );
                    storageContainer.open();
                    storageContainer.write( '<' + scriptTag + '>document.w=window</' + scriptTag + '><iframe src="/favicon.ico"></iframe>' );
                    storageContainer.close();
                    storageOwner = storageContainer.w.frames[0].document;
                    storage = storageOwner.createElement( 'div' );
                }
                catch( e ){
                    storage = doc.createElement( 'div' );
                    storageOwner = doc.body;
                }
                var withIEStorage = function( storeFunction ){
                    return function(){
                        var args = Array.prototype.slice.call( arguments, 0 );
                        args.unshift( storage );
                        storageOwner.appendChild( storage );
                        storage.addBehavior( '#default#userData' );
                        storage.load( localStorageName );
                        var result = storeFunction.apply( store, args );
                        storageOwner.removeChild( storage );
                        return result;
                    }
                };

                var forbiddenCharsRegex = new RegExp( "[!\"#$%&'()*+,/\\\\:;<=>?@[\\]^`{|}~]", "g" );
                var ieKeyFix = function( key ){
                    return key.replace( /^d/, '___$&' ).replace( forbiddenCharsRegex, '___' );
                };
                ret.set = withIEStorage( function( storage, key, val ){
                    key = ieKeyFix( key );
                    if( val === undefined ){
                        return ret.remove( key )
                    }
                    storage.setAttribute( key, ret.serialize( val ) );
                    storage.save( localStorageName );
                    return val;
                } );
                ret.get = withIEStorage( function( storage, key, defaultVal ){
                    key = ieKeyFix( key );
                    var val = ret.deserialize( storage.getAttribute( key ) );
                    return (val === undefined ? defaultVal : val);
                } );
                ret.remove = withIEStorage( function( storage, key ){
                    key = ieKeyFix( key );
                    storage.removeAttribute( key );
                    storage.save( localStorageName );
                } );
                ret.clear = withIEStorage( function( storage ){
                    var attributes = storage.XMLDocument.documentElement.attributes;
                    storage.load( localStorageName );
                    for( var i = attributes.length - 1; i >= 0; i-- ){
                        storage.removeAttribute( attributes[i].name );
                    }
                    storage.save( localStorageName );
                } );
                ret.forEach = withIEStorage( function( storage, callback ){
                    if( $.isFunction( callback ) ){
                        var attributes = storage.XMLDocument.documentElement.attributes;
                        for( var i = 0, attr; attr = attributes[i]; ++i ){
                            callback( attr.name, ret.deserialize( storage.getAttribute( attr.name ) ) );
                        }
                    }
                } );
            }
            return ret;
        })(),
        /**
         * *******************************************************************
         * JS版的格式打印函数
         * *******************************************************************
         */
        sprintf: function(){
            var regex = /%%|%(\d+\$)?([-+'#0 ]*)(\*\d+\$|\*|\d+)?(?:\.(\*\d+\$|\*|\d+))?([scboxXuideEfFgG])/g;
            var a = arguments;
            var i = 0;
            var format = a[i++];
            var _pad = function( str, len, chr, leftJustify ){
                if( !chr ){
                    chr = ' ';
                }
                var padding = (str.length >= len) ? '' : new Array( 1 + len - str.length >>> 0 ).join( chr );
                return leftJustify ? str + padding : padding + str;
            };
            var justify = function( value, prefix, leftJustify, minWidth, zeroPad, customPadChar ){
                var diff = minWidth - value.length;
                if( diff > 0 ){
                    if( leftJustify || !zeroPad ){
                        value = _pad( value, minWidth, customPadChar, leftJustify );
                    }
                    else{
                        value = [
                            value.slice( 0, prefix.length ),
                            _pad( '', diff, '0', true ),
                            value.slice( prefix.length )
                        ].join( '' );
                    }
                }
                return value;
            };
            var _formatBaseX = function( value, base, prefix, leftJustify, minWidth, precision, zeroPad ){
                // Note: casts negative numbers to positive ones
                var number = value >>> 0;
                prefix = (prefix && number && {
                    '2': '0b',
                    '8': '0',
                    '16': '0x'
                }[base]) || '';
                value = prefix + _pad( number.toString( base ), precision || 0, '0', false );
                return justify( value, prefix, leftJustify, minWidth, zeroPad );
            };
            // _formatString()
            var _formatString = function( value, leftJustify, minWidth, precision, zeroPad, customPadChar ){
                if( precision !== null && precision !== undefined ){
                    value = value.slice( 0, precision );
                }
                return justify( value, '', leftJustify, minWidth, zeroPad, customPadChar );
            };
            // doFormat()
            var doFormat = function( substring, valueIndex, flags, minWidth, precision, type ){
                var number, prefix, method, textTransform, value;
                if( substring === '%%' ){
                    return '%';
                }
                // parse flags
                var leftJustify = false;
                var positivePrefix = '';
                var zeroPad = false;
                var prefixBaseX = false;
                var customPadChar = ' ';
                var flagsl = flags.length;
                var j;
                for( j = 0; j < flagsl; j++ ){
                    switch( flags.charAt( j ) ){
                        case ' ':
                            positivePrefix = ' ';
                            break;
                        case '+':
                            positivePrefix = '+';
                            break;
                        case '-':
                            leftJustify = true;
                            break;
                        case "'":
                            customPadChar = flags.charAt( j + 1 );
                            break;
                        case '0':
                            zeroPad = true;
                            customPadChar = '0';
                            break;
                        case '#':
                            prefixBaseX = true;
                            break;
                    }
                }
                // parameters may be null, undefined, empty-string or real valued
                // we want to ignore null, undefined and empty-string values
                if( !minWidth ){
                    minWidth = 0;
                }
                else if( minWidth === '*' ){
                    minWidth = +a[i++];
                }
                else if( minWidth.charAt( 0 ) === '*' ){
                    minWidth = +a[minWidth.slice( 1, -1 )];
                }
                else{
                    minWidth = +minWidth;
                }
                // Note: undocumented perl feature:
                if( minWidth < 0 ){
                    minWidth = -minWidth;
                    leftJustify = true;
                }
                if( !isFinite( minWidth ) ){
                    throw new Error( 'sprintf: (minimum-)width must be finite' );
                }
                if( !precision ){
                    precision = 'fFeE'.indexOf( type ) > -1 ? 6 : (type === 'd') ? 0 : undefined;
                }
                else if( precision === '*' ){
                    precision = +a[i++];
                }
                else if( precision.charAt( 0 ) === '*' ){
                    precision = +a[precision.slice( 1, -1 )];
                }
                else{
                    precision = +precision;
                }
                // grab value using valueIndex if required?
                value = valueIndex ? a[valueIndex.slice( 0, -1 )] : a[i++];
                switch( type ){
                    case 's':
                        return _formatString( value + '', leftJustify, minWidth, precision, zeroPad, customPadChar );
                    case 'c':
                        return _formatString( String.fromCharCode( +value ), leftJustify, minWidth, precision, zeroPad );
                    case 'b':
                        return _formatBaseX( value, 2, prefixBaseX, leftJustify, minWidth, precision, zeroPad );
                    case 'o':
                        return _formatBaseX( value, 8, prefixBaseX, leftJustify, minWidth, precision, zeroPad );
                    case 'x':
                        return _formatBaseX( value, 16, prefixBaseX, leftJustify, minWidth, precision, zeroPad );
                    case 'X':
                        return _formatBaseX( value, 16, prefixBaseX, leftJustify, minWidth, precision, zeroPad ).toUpperCase();
                    case 'u':
                        return _formatBaseX( value, 10, prefixBaseX, leftJustify, minWidth, precision, zeroPad );
                    case 'i':
                    case 'd':
                        number = +value || 0;
                        // Plain Math.round doesn't just truncate
                        number = Math.round( number - number % 1 );
                        prefix = number < 0 ? '-' : positivePrefix;
                        value = prefix + _pad( String( Math.abs( number ) ), precision, '0', false );
                        return justify( value, prefix, leftJustify, minWidth, zeroPad );
                    case 'e':
                    case 'E':
                    case 'f': // @todo: Should handle locales (as per setlocale)
                    case 'F':
                    case 'g':
                    case 'G':
                        number = +value;
                        prefix = number < 0 ? '-' : positivePrefix;
                        method = ['toExponential', 'toFixed', 'toPrecision']['efg'.indexOf( type.toLowerCase() )];
                        textTransform = ['toString', 'toUpperCase']['eEfFgG'.indexOf( type ) % 2];
                        value = prefix + Math.abs( number )[method]( precision );
                        return justify( value, prefix, leftJustify, minWidth, zeroPad )[textTransform]();
                    default:
                        return substring;
                }
            };
            return format.replace( regex, doFormat );
        }
    };
})();


/**
 * *************************************
 * 用小tips作用轻提示
 * *************************************
 */
;(function( $ ){
    var defaultShowTipMsgOption = {
        delegate: null,//在动态页面中，用于代理的DOM元素，可选
        title: "",//提示的标题，可选
        content: "",//内容，必选
        type: "",//类型：danger/warning/info/success/primary
        placement: "top",//tips的显示位置， top、bottom、left、right、auto
        timeout: 2,//超时自动消失，秒数
        onHide: ""//弹框关闭时的回调函数
    };
    /**
     * 入口
     * @param    options        选项
     */
    $.fn.showTipMsg = function( options ){
        var _opt = $.extend( {}, defaultShowTipMsgOption, options );
        var _this = $( this );
        _this.popover( {
            title: _opt.title,
            content: _opt.content,
            placement: _opt.placement,
            selector: _opt.delegate,
            template: '<div class="popover ' + "popover-" + _opt.type + '" role="tooltip"><div class="arrow"></div><h3 class="popover-title"></h3><div class="popover-content"></div></div>'
        } ).popover( 'show' ).on( 'shown.bs.popover', function(){
            var timeout = parseInt( _opt.timeout );
            if( timeout > 0 ){
                var tipObj = this;
                window.setTimeout( function(){
                    $( tipObj ).popover( 'destroy' );
                }, 1000 * parseInt( timeout ) );
            }
        } ).on( 'hidden.bs.popover', function(){
            if( typeof _opt.onHide === 'function' ){
                _opt.onHide();
            }
        } );
        return _this;
    };
})( jQuery );
/**
 * *************************************
 * 显示确认框
 * *************************************
 */
;(function( $ ){
    var defaultShowConfirmOptions = {
        title: "是否确认此操作？",//默认的显示内容
        placement: "top",//弹框显示的位置：top/bottom/left/right
        onConfirm: "",//点确定时的回调函数
        onCancel: ""//点取消时的回调函数
    };
    /**
     * 入口
     * @param    options        选项
     */
    $.fn.showConfirm = function( options ){
        var _opt = $.extend( {}, defaultShowConfirmOptions, options );
        var confirmFunc = _opt.onConfirm;
        var cancelFunc = _opt.onCancel;
        var _this = $( this );
        _this.confirmation( "destroy" );
        /******************设置默认的回调函数******************/
        if( typeof confirmFunc !== 'function' ){
            confirmFunc = function(){
                _this.confirmation( "destroy" );
            };
        }
        if( typeof cancelFunc !== 'function' ){
            cancelFunc = function(){
                _this.confirmation( "destroy" );
            };
        }
        _this.confirmation( {
            singleton: true,
            title: _opt.title,
            placement: _opt.placement,
            popout: true,
            onConfirm: confirmFunc,
            onCancel: cancelFunc,
            btnOkClass: "btn btn-danger btn-sm confirm-ok",
            btnOkLabel: "确定",
            btnCancelLabel: "取消",
            container: "body",
            href: "javascript:void(0);"
        } ).confirmation( 'show' );
        return _this;
    };
})( jQuery );
/**
 * *************************************
 * 在小弹框里修改单个字段的值
 * *************************************
 */
;(function( $ ){
    var defaultPopoverEditOptions = {
        title: "修改操作",
        placement: "top",//弹框显示的位置：top/bottom/left/right
        type: "input",//修改字段时用的类型：input/select/date/datetime/textarea/file
        /*通用可选的选项*/
        css: "",//赋给输入框的CSS样式，如“width:200px;height:30px”
        onConfirm: "",//点击确定之后的回调函数，参数为输入的DOM对象
        onBeforeShow: "",//弹框弹出之前的回调函数，无参数
        onAfterShow: "",//显示之后的回调函数，参数为输入的DOM元素
        attrs: "",//给输入对象添加的额外属性，如{wenao:1,abc:2}
        beforeCheckbox: false,//前面加checkbox元素
        checkboxText: "",//checkbox的名称
        /*当type==input/date/datetime/textarea时的默认值*/
        defaultVal: null,//为null时不用它，非null时才用
        /*当type==file时用的属性*/
        fileUrl: "",//文件要上传的地址
        fileText: "选择要上传的文件",//大按钮上显示的文本
        fileName: "runtofu_file",//input[file]的name
        onUploadSuccess: null,
        /*当type==select时的选项，其中selectData和selectUrl二必选一，其余可选*/
        selectWidth: 150,//框的宽度
        selectData: "",//数据源，json对象，格式要求[{value:"",text:""},........]
        selectUrl: "",//远程数据源，要求返回的数据里data字段的格式[{value:"",text:""},........]
        selectFilter: ""//选项过滤函数，返回true的才写到options里去
    };
    /**
     * 入口
     * @param    options        选项
     */
    $.fn.popoverEdit = function( options ){
        var _opt = $.extend( {}, defaultPopoverEditOptions, options );
        var tid = tofuUtils.juuid();
        var content = "";
        var _this = $( this );
        _this.popover( "destroy" );
        var originVal = _opt.defaultVal == null ? _this.text() : _opt.defaultVal;
        var css = _opt.css;
        var t = _opt.type;
        var attlist = "";
        $.each( _opt.attrs, function( k, v ){
            if( k.length > 0 ){
                attlist += " " + k + "=\"" + v + "\" ";
            }
        } );
        /*如果前面有checkbox元素*/
        var cbx = "";
        if( _opt.beforeCheckbox ){
            cbx = '<label class="checkbox-inline"><input type="checkbox">' + _opt.checkboxText + '</label>';
        }
        /************这三种类型的输入对象是一致的************/
        if( t == "date" || t == "datetime" || t == "input" ){
            if( css.length <= 0 && (t == "date" || t == "datetime") ){
                css = (t == "date") ? "width: 100px" : "width: 120px";
            }
            content = '<div class="form-inline" role="form" elid="' + tid + '"><div class="form-group form-edit">' + cbx + '<input type="text" class="form-control input-sm" id="' + tid + '" value="' + originVal + '" style="' + css + '"' + attlist + '></div><button class="btn btn-primary btn-sm"><span class="glyphicon glyphicon-ok"></span></button></div>';
        }
        /************下拉选择框特殊处理，数据来源可用现成的或远程加载************/
        else if( t == "select" ){
            /*默认过滤函数*/
            var filter = _opt.selectFilter;
            if( typeof filter !== 'function' ){
                filter = function( sd ){
                    return true;
                };
            }
            /*如果提供的selectData为空则远程请求加载数据*/
            var selectOp = "";
            if( typeof _opt.selectData != 'object' || _opt.selectData == null || !_opt.selectData ){
                if( _opt.selectUrl.length <= 0 ){
                    return tofuUtils.modalMsg( {
                        type: "danger",
                        content: "popoverEdit的type为select时，选项selectData和selectUrl不能同时为空！"
                    } );
                }
                tofuUtils.sendRequest( {
                    url: _opt.selectUrl,
                    onSuccess: function( ret ){
                        op = getSelectOption( ret, originVal, filter, attlist );
                        $( "#" + tid ).append( op ).selectpicker( 'refresh' );
                    }
                } );
            }
            else{
                selectOp = getSelectOption( _opt.selectData, originVal, filter, attlist );
            }
            if( css.length <= 0 ){
                css = "min-width: 120px";
            }
            content = '<div class="form-inline" role="form"><div class="form-group">' + cbx + '<select id="' + tid + '" style="' + css + '">' + selectOp + '</select></div><button class="btn btn-primary btn-sm btn-st"><span class="glyphicon glyphicon-ok"></span></button></div>';
        }
        /************textarea也特殊处理************/
        else if( t == "textarea" ){
            content = '<div class="form-group form-textarea"><textarea id="' + tid + '" style="' + css + '"' + attlist + '>' + originVal + '</textarea></div>' + cbx + '<button class="btn btn-primary btn-sm btn-te"><span class="glyphicon glyphicon-ok"></span></button>';
        }
        /************处理上传文件部分************/
        else if( t == "file" ){
            attlist = "";
            $.each( _opt.attrs, function( k, v ){
                if( k.length > 0 ){
                    attlist += '<input type="hidden" name="' + k + '" value="' + v + '">';
                }
            } );
            content = '<form class="form-inline" role="form" method="POST" action="' + _opt.fileUrl + '"><button style="width:100%" type="button" class="btn btn-default">' + _opt.fileText + '</button>' + attlist + '<input id="' + tid + '" name="' + _opt.fileName + '" type="file" class="up-file"></form>';
        }
        else{
            tofuUtils.modalMsg( {
                type: "danger",
                content: "popoverEdit暂不支持类型：" + t + "，目前只支持input/select/date/datetime/textarea/file"
            } );
            return _this;
        }
        $( '*' ).popover( 'hide' );
        _this.popover( {
            html: true,
            content: content,
            placement: _opt.placement,
            title: _opt.title,
            container: "body",
            template: '<div class="popover" role="tooltip"><div class="arrow"></div><div class="popover-header"><div class="popover-close"><button class="close">×</button></div><div class="popover-title"></div></div><div class="popover-content"></div></div>'
        } ).on( 'show.bs.popover', function(){
            if( typeof _opt.onBeforeShow === 'function' ){
                _opt.onBeforeShow();
            }
        } ).on( 'shown.bs.popover', function(){
            if( typeof _opt.onAfterShow === 'function' ){
                _opt.onAfterShow( $( "#" + tid ) );
            }
        } );
        _this.popover( 'show' );
        var _el = $( "#" + tid );
        /******************激活下拉选择框******************/
        if( t == "select" ){
            _el.selectpicker( {
                maxOptions: 1,
                width: _opt.selectWidth + "px"
            } );
        }
        /******************激活日期选择框******************/
        if( t == "date" ){
            _el.datetimepicker( {
                value: originVal,
                timepicker: false,
                formatDate: "Y-m-d",
                format: "Y-m-d"
            } );
        }
        /******************激活日期、时间选择框******************/
        if( t == "datetime" ){
            _el.datetimepicker( {
                value: originVal,
                format: "Y-m-d H:i"
            } );
        }
        /******************上传文件的AjaxSubmit******************/
        if( t == 'file' ){
            var form = _el.parents( "form" );
            form.find( "button" ).click( function(){
                _el.trigger( "click" );
            } );
            _el.change( function(){
                form.ajaxSubmit( {
                    success: function( retobj ){
                        var errno = parseInt( retobj.result.errno );
                        if( errno !== 0 ){
                            tofuUtils.modalMsg( {
                                type: "danger",
                                content: "上传文件出错：" + retobj.result.errmsg
                            } );
                        }
                        else{
                            if( typeof _opt.onUploadSuccess === 'function' ){
                                _opt.onUploadSuccess( retobj.data );
                            }
                        }
                    },
                    error: function( XMLHttpRequest, textStatus, errorThrown ){
                        tofuUtils.loading.end();
                        var errmsg = "Request.status = " + XMLHttpRequest.status;
                        if( XMLHttpRequest.responseText.length > 0 ){
                            errmsg = "\nRequest.responseText = " + XMLHttpRequest.responseText;
                        }
                        if( errorThrown.length > 0 ){
                            errmsg += "\nerrorThrown = " + errorThrown;
                        }
                        if( textStatus.length > 0 ){
                            errmsg += "\ntextStatus = " + textStatus;
                        }
                        tofuUtils.modalMsg( { content: errmsg, type: "danger", title: "请求失败" } );
                    }
                } );
            } );
        }
        /******************添加销毁事件******************/
        _el.parents( "div.popover" ).find( "button.close" ).click( function(){
            _this.popover( "destroy" );
            return false;
        } );
        /******************添加确认事件******************/
        _el.parents( "div.popover-content" ).find( "span.glyphicon-ok" ).parent().click( function(){
            if( typeof  _opt.onConfirm === "function" ){
                _opt.onConfirm( _el );
            }
            _this.popover( "destroy" );
            return false;
        } );
        return _this;
    };

    /**
     * 组装select类型的option列表
     */
    function getSelectOption( op, originVal, filter, attlist ){
        var ret = "";
        var attr = "";
        $( op ).each( function( index, val ){
            if( !filter( val ) ){
                return;
            }
            attr = attlist;
            if( originVal == val.text ){
                attr += " selected=\"true\" ";
            }
            $.each( val, function( k, v ){
                if( k != 'value' && k != 'text' ){
                    attr += " " + k + "=\"" + v + "\" ";
                }
            } );
            ret += '<option value="' + val.value + '" ' + attr + '>' + val.text + '</option>';
        } );
        return ret;
    }
})( jQuery );
/**
 * *************************************
 * Bootstrap模态框可拖动
 * *************************************
 */
;(function( $ ){
    var oldModal = $.fn.modal;
    $.fn.modal = function( o, _r ){
        var _this = $( this );
        if( !_this.attr( 'ifbindmv' ) ){
            _this.attr( 'isbindmv', '1' );
            var _head = _this.find( '.modal-header' );
            var _dialog = _this.find( '.modal-dialog' );
            _head.css( "cursor", "move" );
            var move = {
                isMove: false,
                left: 0,
                top: 0
            };
            _this.on( 'mousemove', function( e ){
                if( !move.isMove ) return;
                _dialog.offset( {
                    top: e.pageY - move.top,
                    left: e.pageX - move.left
                } );
            } ).on( 'mouseup', function( e ){
                move.isMove = false;
            } );
            _head.on( 'mousedown', function( e ){
                move.isMove = true;
                var offset = _dialog.offset();
                move.left = e.pageX - offset.left;
                move.top = e.pageY - offset.top;
            } );
        }
        return oldModal.call( this, o, _r );
    }
})( jQuery );