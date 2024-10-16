// ==UserScript==
// @name         Auto VC
// @namespace    http://tampermonkey.net/
// @version      2024-10-10
// @description  try to take over the world!
// @author       Free Server
// @match        https://free.vps.vc/create-vps
// @icon         https://www.google.com/s2/favicons?sz=64&domain=vps.vc
// @grant        none
// ==/UserScript==

(function () {
    'use strict';

    var datacenter = document.getElementById("datacenter");

    //datacenter.options.add(new Option("US4-CHI", "value3"));
    //datacenter.options.add(new Option("EU1-CHI", "value1"));
    //datacenter.options.add(new Option("US2-CHI", "value3"));
    //datacenter.options.add(new Option("US3-CHI", "value3"));
    //datacenter.options.add(new Option("US1-CHI", "value1"));
    //datacenter.options.add(new Option("CA1-CHI", "value2"));

    // 自动优选地区和填表
    datacenter.size = datacenter.options.length;
    if (datacenter.options.length === 1 && datacenter.options[0].text === "-select-") {
        location.reload();
        return
    }

    if (datacenter.options.length > 1) {
        datacenter.options[1].selected = true;
    }

    var flag = false;
    var i = 0;

    for (i = 0; i < datacenter.options.length; i++) {
        if (datacenter.options[i].text.includes("US1")) {
            datacenter.options[i].selected = true;
            flag = true;
            break;
        }
    }

    if (!flag) {
        for (i = 0; i < datacenter.options.length; i++) {
            if (datacenter.options[i].text.includes("CA1")) {
                datacenter.options[i].selected = true;
                flag = true;
                break;
            }
        }
    }

    if (!flag) {
        for (i = 0; i < datacenter.options.length; i++) {
            if (datacenter.options[i].text.includes("1")) {
                datacenter.options[i].selected = true;
                break;
            }
        }
    }

    document.getElementById("os").value = 2;
    document.getElementById("password").value = 123456;
    document.getElementById("purpose").value = 1;

    var elements = document.getElementsByName('agreement[]');
    for (i = 0; i < elements.length; i++) {
        elements[i].checked = true;
    }

    // 图片验证码获取焦点,这样进来输入答案直接点提交即可最大节省时间
    var result = document.getElementById("result");
    if (result != null) {
        result.focus();
    }

    // 展示大约用时,别太嚣张
    var time = 2;
    var create_btn = document.getElementById("create_btn");
    setInterval(function () {
        time++;
        create_btn.textContent = "大约用时: " + time;
    }, 1000);

    // 自动点击hCaptcha, 注意:"需要配合无障碍cookie和--disable-web-security"不然会出错
    var hCaptchaID = 0;
    hCaptchaID = setInterval(function () {
        var iframes = document.getElementsByTagName("iframe");

        for (var i = 0; i < iframes.length; i++) {
            var iframeDoc = iframes[i].contentDocument || iframes[i].contentWindow.document;
            var anchor = iframeDoc.querySelector("#anchor");
            if (anchor != null) {
                anchor.click();
                clearInterval(hCaptchaID);
            }
        }
    }, 100);

})();