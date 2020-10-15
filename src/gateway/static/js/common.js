function get(url, data, successFunc) {
    if (url[0]=='/')
        url=url.substring(1)
    url = "/html/" + url

    if (data == null) {
        data = {}
    }
    // let lvToken = getCookie("token")
    // if (lvToken) {
    //     data["token"] = lvToken;
    // }

    $.ajax({
        type: "GET",
        url: url,
        data: data,
        dataType: "json",
        beforeSend: function (msg) {
            console.log("beforeSend", msg)
        },
        success: function (msg) {
            console.log("successs", msg)
            // alert( "Data Saved: " + msg );
            if (successFunc) {
                successFunc(msg)
                return;
            }
            alert(msg.msg)
            if (msg.url && msg.url.length > 0) {
                location.href = msg.url
            }
        },
        error: function (msg) {
            console.log("error", msg)
        }
    });
}

function post(url, data, successFunc,FuncResultDeal) {
    if (url[0]=='/')
        url=url.substring(1)
    url = "/html/" + url

    if (data == null) {
        data = {}
    }
    // let lvToken = getCookie("token")
    // if (lvToken) {
    //     data["token"] = lvToken;
    // }

    $.ajax({
        type: "POST",
        url: url,
        data: JSON.stringify(data),
        dataType: "json",
        beforeSend: function (msg) {
            console.log("beforeSend", msg)
        },
        success: function (msg) {
            console.log("successs", msg)
            // alert( "Data Saved: " + msg );
            if (successFunc) {
                successFunc(msg)
                return;
            }
            alert(msg.msg)
            if (FuncResultDeal) {
                FuncResultDeal(msg)
                return;
            }
            if (msg.url && msg.url.length > 0) {
                location.href = msg.url
            }
        },
        error: function (msg) {
            console.log("error", msg)
        }
    });
}

function postFile(url, data, successFunc) {
    url = "/gatewayapi/" + url

    // let lvToken = getCookie("token")
    // if (lvToken) {
    //     data["token"] = lvToken;
    // }

    var formData = new FormData();
    formData.append("file", $("input[type=file]")[0].files[0]);
    for (let item in data) {
        formData.append(item, data[item]);
    }

    // formData.append("token","");
    $.ajax({
        url: url, /*接口域名地址*/
        type: 'post',
        data: formData,
        dataType: "json",
        contentType: false,
        processData: false,
        success: function (msg) {
            console.log("successs", msg)
            // alert( "Data Saved: " + msg );
            if (successFunc) {
                successFunc(msg)
            }
            if (msg.errcode < 0) {
                alert(msg.msg)

                return
            }


        },
        error: function (res) {
            console.log(res);
            $("#" + fileid).val("");
        }
    });

}

function checkLogin() {
    let admintoken = this.getCookie("admintoken");
    let token = this.getCookie("token");
    if (location.href.indexOf("/admin") > 0) {
        if (admintoken) {
            let html = "<a href='/admin_user.html'>后台系统</a><a href='/admin_logout'>退出系统</a>";
            $("#logininfo").html(html)

            let usermenu = $("#usermenu")
            if (usermenu.length > 0) {
                html = '                <a href="/admin_pwd.html"><i class="fa fa-key " aria-hidden="true"></i> 修改密码</a>\n' +
                    '                <a href="/admin_addmanager.html"><i class="fa fa-user " aria-hidden="true"></i> 新增管理员</a>\n' +
                    '                <a href="/admin_user.html"><i class="fa fa-user " aria-hidden="true"></i> 用户管理</a>\n' +
                    '                <a href="/admin_projectshequ.html"><i class="fa fa-university " aria-hidden="true"></i> 社区版管理</a>\n' +
                    '                <a href="/admin_projectqiye.html"><i class="fa fa-university " aria-hidden="true"></i> 企业版管理</a>\n' +
                    '                <a href="/admin_projectpaylog.html"><i class="fa fa-university " aria-hidden="true"></i> 付费日志</a>\n'
                $("#usermenu").html(html)
            }
            return
        }
    } else {
        if (token && location.href.indexOf("/user") > 0) {
            this.post("token", null, function (data) {
                console.log(data)
                if (data.errcode >= 0) {
                    let html = "<a href='/user_info.html'>" + data.username + "</a><a href='/user_logout'>退出系统</a>";
                    $("#logininfo").html(html)

                    let usermenu = $("#usermenu")
                    if (usermenu.length > 0) {
                        html = '                <a href="/user_pwd.html"><i class="fa fa-key " aria-hidden="true"></i> 修改密码</a>\n' +
                            '                <a href="/user_info.html"><i class="fa fa-user " aria-hidden="true"></i> 用户认证</a>\n' +
                            '                <a href="/user_projectshequ.html"><i class="fa fa-university " aria-hidden="true"></i> 社区版管理</a>\n' +
                            '                <a href="/user_project.html"><i class="fa fa-university " aria-hidden="true"></i> 企业版管理</a>\n' +
                            '                <a href="/user_projectpaylog.html"><i class="fa fa-university " aria-hidden="true"></i> 支付日志</a>\n' +
                            '                <a href="/user_kefu.html"><i class="fa fa-link" aria-hidden="true"></i> 微信客服</a>'
                        $(usermenu).html(html)
                    }
                }
            })
        }
        if (admintoken && token) {
            let html = "<a href='/user_info.html'>会员中心</a><a href='/admin_user.html'>后台系统</a>";
            $("#logininfo").html(html)
        } else if (admintoken) {
            let html = "<a href='/admin_user.html'>后台系统</a><a href='/admin_logout'>退出系统</a>";
            $("#logininfo").html(html)
        } else if (token) {
            this.post("token", null, function (data) {
                console.log(data)
                if (data.errcode >= 0) {
                    let html = "<a href='/user_info.html'>" + data.username + "</a><a href='/user_logout'>退出系统</a>";
                    $("#logininfo").html(html)

                    let usermenu = $("#usermenu")
                    if (usermenu.length > 0) {
                        html = '                <a href="/user_pwd.html"><i class="fa fa-key " aria-hidden="true"></i> 修改密码</a>\n' +
                            '                <a href="/user_info.html"><i class="fa fa-user " aria-hidden="true"></i> 用户认证</a>\n' +
                            '                <a href="/user_projectshequ.html"><i class="fa fa-university " aria-hidden="true"></i> 社区版管理</a>\n' +
                            '                <a href="/user_project.html"><i class="fa fa-university " aria-hidden="true"></i> 企业版管理</a>\n' +
                            '                <a href="/user_projectpaylog.html"><i class="fa fa-university " aria-hidden="true"></i> 支付日志</a>\n' +
                            '                <a href="/user_kefu.html"><i class="fa fa-link" aria-hidden="true"></i> 微信客服</a>'
                        $(usermenu).html(html)
                    }
                }
            })
        }
    }


}

function setCookie(name, value,expire) {
    if (!expire){
        expire=30*24 * 60 * 60 * 1000
    }
    let exp = new Date();
    exp.setTime(exp.getTime() + expire);
    document.cookie = name + "=" + escape(value) + ";path=/;expires=" + exp.toGMTString();
}

function getCookie(name) {
    let arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg))
        return unescape(arr[2]);
    else
        return null;
}

function delCookie(name) {
    setCookie(name,"",-100)
}

//获取url中的参数
function getUrlParam(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); //构造一个含有目标参数的正则表达式对象
    var r = window.location.search.substr(1).match(reg);  //匹配目标参数
    if (r != null) return unescape(r[2]); return null; //返回参数值
}


function formsubmit(Vue, ctrid, posturl, FuncOverright, FuncResultDeal) {
    let lvCheck = pintuerCheckFormSubmit(ctrid);
    if (!lvCheck)
        return -1;

    let lvText = $("#" + ctrid + " input[type='text'] ,#" + ctrid + " input[type='number'] ,#" + ctrid + " input[type='email'] ,#" + ctrid + " input[type='password'],#" + ctrid + " input[type='hidden']");
    let lvRadio = $("#" + ctrid + " input[type='radio']");
    let lvCheckBox = $("#" + ctrid + " input[type='checkbox']");
    let lvSelect = $("#" + ctrid + " select");
    let lvTextArea = $("#" + ctrid + " textarea");

    let params = {};
    for (let i = 0; i < lvText.length; i++) {
        if ($(lvText[i]).attr("name") == undefined)
            continue;
        // params[$(lvText[i]).attr("name")] = encodeURIComponent($(lvText[i]).val());
        let ctlname = $(lvText[i]).attr("name")
        if (ctlname.indexOf('[]') > 0) {
            if (params[ctlname] == undefined)
                params[ctlname] = new Array()
            params[ctlname].push($(lvText[i]).val())
        } else
            params[$(lvText[i]).attr("name")] = $.trim($(lvText[i]).val());
    }
    for (let i = 0; i < lvTextArea.length; i++) {
        if ($(lvTextArea[i]).attr("name") == undefined)
            continue;
        // params[$(lvTextArea[i]).attr("name")] = encodeURIComponent($(lvTextArea[i]).val());
        params[$(lvTextArea[i]).attr("name")] = $.trim($(lvTextArea[i]).val());
    }
    for (let i = 0; i < lvRadio.length; i++) {
        if ($(lvRadio[i]).attr("name") == undefined)
            continue;
        if (lvRadio[i].checked)
            params[$(lvRadio[i]).attr("name")] = $(lvRadio[i]).val();
    }

    for (let i = 0; i < lvSelect.length; i++) {
        if ($(lvSelect[i]).attr("name") == undefined)
            continue;
        params[$(lvSelect[i]).attr("name")] = $(lvSelect[i]).val();
    }

    post(posturl, params, FuncOverright, FuncResultDeal);
}

/**
 * Comment
 */
function LoadValid() {
    $('textarea, input, select').blur(function () {
        let e = $(this);
        console.log(e);
        if (e.attr("validate")) {
            //e.closest('.field').find(".input-help").remove();
            e.parent().find(".input-help").remove();
            let $checkdata = e.attr("validate").split(',');
            let $checkvalue = e.val();
            let $checkstate = true;
            let $checktext = "";
            if (e.attr("placeholder") == $checkvalue) {
                $checkvalue = "";
            }
            if ($checkvalue != "" || e.attr("validate").indexOf("required") >= 0) {
                for (let i = 0; i < $checkdata.length; i++) {
                    let $checktype = $checkdata[i].split(':');
                    if (!$webValid(e, $checktype[0], $checkvalue)) {
                        $checkstate = false;

                        if ($.trim($checktype[1]) == "") {
                            $checktext = $txtValid($checktype[0]);//"不能为空！";

                            break;
                        }
                        if ($checktext != "") {
                            //$checktext = $checktext + ",";
                            break;
                        }
                        $checktext = $checktext + $checktype[1] + "";

                    }
                }
            }

            if ($checkstate) {
                $(e).css("border-color", "");
                e.parent().removeClass("check-error");
                e.parent().find(".input-help").remove();
                e.parent().addClass("check-success");
            } else {
                $(e).css("border-color", "red");
                if ($(e).attr("placeholder") != "" && $(e).attr("placeholder") != undefined)
                    $checktext = "请输入" + $(e).attr("placeholder");
                else $checktext = "不能为空"
                e.parent().removeClass("check-success");
                e.parent().addClass("check-error");
                $checktext = "";
                e.parent().append('<span class="input-help">' + $checktext + '</span>');
            }

        }
    });
    let $webValid = function (element, type, value) {
        value = $.trim(value);
        let $choice = value.replace(/(^\s*)|(\s*$)/g, "");
        switch (type) {
            case "required":
                return /[^(^\s*)|(\s*$)]/.test($choice);
                break;
            case "chinese":
                return /^[\u0391-\uFFE5]+$/.test($choice);
                break;
            case "number":
                return /^\d+$/.test($choice);
                break;
            case "integer":
                return /^[-\+]?\d+$/.test($choice);
                break;
            case "plusinteger":
                return /^[+]?\d+$/.test($choice);
                break;
            case "double":
                return /^[-\+]?\d+(\.\d+)?$/.test($choice);
                break;
            case "plusdouble":
                return /^[+]?\d+(\.\d+)?$/.test($choice);
                break;
            case "english":
                return /^[A-Za-z]+$/.test($choice);
                break;
            case "lettterint":
                return /^[0-9A-Za-z]+$/.test($choice);
                break;
            case "username":
                return /^[a-z]\w{3,}$/i.test($choice);
                break;
            //case "mobile": return /^((\(\d{3}\))|(\d{3}\-))?13[0-9]\d{8}?$|15[89]\d{8}?$|170\d{8}?$|147\d{8}?$/.test($choice); break;
            case "mobile":
                return /^((\(\d{3}\))|(\d{3}\-))?1[3-9][0-9]\d{8}?$/.test($choice);
                break;
            case "phone":
                return /^((\(\d{2,3}\))|(\d{3}\-))?(\(0\d{2,3}\)|0\d{2,3}-)?[1-9]\d{6,7}(\-\d{1,4})?$/.test($choice);
                break;
            case "tel":
                return /^((\(\d{3}\))|(\d{3}\-))?13[0-9]\d{8}?$|15[89]\d{8}?$|170\d{8}?$|147\d{8}?$/.test($choice) || /^((\(\d{2,3}\))|(\d{3}\-))?(\(0\d{2,3}\)|0\d{2,3}-)?[1-9]\d{6,7}(\-\d{1,4})?$/.test($choice);
                break;
            case "email":
                return /^[^@]+@[^@]+\.[^@]+$/.test($choice);
                break;
            case "url":
                return /^http:\/\/[A-Za-z0-9]+\.[A-Za-z0-9]+[\/=\?%\-&_~`@[\]\':+!]*([^<>\"\"])*$/.test($choice);
                break;
            case "urlpath":
                return /^[A-Za-z0-9\/]+$/.test($choice);
                break;
            case "ip":
                return /^[\d\.]{7,15}$/.test($choice);
                break;
            case "qq":
                return /^[1-9]\d{4,10}$/.test($choice);
                break;
            case "currency":
                return /^\d+(\.\d+)?$/.test($choice);
                break;
            case "zip":
                return /^[1-9]\d{5}$/.test($choice);
                break;
            case "radio":
                let radio = element.closest('form').find('input[name="' + element.attr("name") + '"]:checked').length;
                return eval(radio == 1);
                break;
            case "json":
                if ($choice == "") return true
                try {
                    var obj = JSON.parse($choice);
                    if (typeof obj == 'object' && obj) {
                        return true
                    } else {
                        return false
                    }
                } catch (e) {
                    return false
                }
                break;
            case "jsonsingle":
                if ($choice == "") return true
                try {
                    var obj = JSON.parse($choice);
                    if (typeof obj == 'object' && obj) {
                        let objitem
                        for (objitem in obj) {
                            if (typeof obj[objitem] == 'object' && obj[objitem]) {
                                return false
                            }
                        }

                        return true
                    } else {
                        return false
                    }
                } catch (e) {
                    return false
                }
                break;
            default:
                let $test = type.split('#');
                if ($test.length > 1) {
                    switch ($test[0]) {
                        case "compare":
                            return eval(Number($choice) + $test[1]);
                            break;
                        case "regexp":
                            return new RegExp($test[1], "gi").test($choice);
                            break;
                        case "length":
                            let $length;
                            if (element.attr("type") == "checkbox") {
                                $length = element.closest('form').find('input[name="' + element.attr("name") + '"]:checked').length;
                            } else {
                                $length = $choice.replace(/[\u4e00-\u9fa5]/g, "***").length;
                            }
                            return eval($length + $test[1]);
                            break;
                        case "ajax":
                            let $getdata;
                            let $url = $test[1] + $choice;
                            $.ajaxSetup({async: false});
                            $.getJSON($url, function (data) {
                                $getdata = data.getdata;
                            });
                            if ($getdata == "true") {
                                return true;
                            }
                            break;
                        case "repeat":
                            return $choice == jQuery('input[name="' + $test[1] + '"]').eq(0).val();
                            break;
                        default:
                            return true;
                            break;
                    }
                    break;
                } else {
                    return true;
                }
        }
    };

    let $txtValid = function (type) {
        switch (type) {
            case "required":
                return "不能为空";
                break;
            case "chinese":
                return "必须为中文";
                break;
            case "number":
                return "必须为数字";
                break;
            case "integer":
                return "必须为整数";
                break;
            case "plusinteger":
                return "必须为正整数";
                break;
            case "double":
                return "必须为小叔数字";
                break;
            case "plusdouble":
                return "必须为正小数";
                break;
            case "english":
                return "必须为字母";
                break;
            case "username":
                return "3位以上的字母";
            //case "mobile": return /^((\(\d{3}\))|(\d{3}\-))?13[0-9]\d{8}?$|15[89]\d{8}?$|170\d{8}?$|147\d{8}?$/.test($choice); break;
            case "mobile":
                return "必须为正确的手机号码";
            case "phone":
                return "必须为电话号码";
            case "tel":
                return "必须为手机号码或电话号码";
            case "email":
                return "必须为邮件地址";
            case "url":
                return "必须为网络地址";
            case "urlpath":
                return "必须为网络地址路径";
            case "ip":
                return "必须为IP地址";
            case "qq":
                return "必须为QQ号码";
            case "currency":
                return "必须为货币数字";
            case "zip":
                return "必须为邮编";
            case "radio":
                return "单选框未选";
                break;
            case "json":
                return "JSON格式有误";
                break;
            case "jsonsingle":
                return "JSON数据必须是一层数据";
                break;
            default:
                return "验证失败";
        }
    }
}


/**
 * 判断需要验证的控件
 */
function pintuerCheckFormSubmit(formid) {
    let lvList = $('#' + formid + ' input[validate],textarea[validate],select[validate]');
    if (lvList.length == 0) return true;

    var $events = lvList.data("events");
    if (!$events || !$events["blur"]) {
        this.LoadValid();
    }


    for (let i = 0; i < lvList.length; i++) {
        if ($(lvList[i]).is(":hidden"))
            continue;
        $(lvList[i]).trigger("blur");
    }
    let numError = $('#' + formid + ' .input-help').length;
    if (numError) {
        $('#' + formid + ' .input-help').first().parent().find('input[validate],textarea[validate],select[validate]').first().focus().select();
        return false;
    }
    return true;
}
