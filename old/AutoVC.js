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

    for (var i = 0; i < datacenter.options.length; i++) {
        if (datacenter.options[i].text.includes("US") && !datacenter.options[i].text.includes("US4")) {
            datacenter.options[i].selected = true;
            break;
        }
        if (datacenter.options[i].text.includes("CA1")) {
            datacenter.options[i].selected = true;
            break;
        }
        if (datacenter.options[i].text.includes("EU1")) {
            datacenter.options[i].selected = true;
            break;
        }
        if (datacenter.options[i].text.includes("1")) {
            datacenter.options[i].selected = true;
            break;
        }
    }

    document.getElementById("os").value = 2;
    document.getElementById("password").value = 123456;
    document.getElementById("purpose").value = 2;

    var elements = document.getElementsByName('agreement[]');
    for (i = 0; i < elements.length; i++) {
        elements[i].checked = true;
    }

    // 图片验证码获取焦点,这样进来输入答案直接点提交即可最大节省时间
    var result = document.getElementById("result");
    if (result != null) {
        result.focus();
    }

    // 展示大约用时
    var time = 0;
    var create_btn = document.getElementById("create_btn");
    setInterval(function () {
        time++;
        create_btn.textContent = "大约用时: " + time / 10.0 + "秒";
        if (time >= 130) {
            create_btn.click(); // 等待13秒自动提交，实际用时还要考虑加载和网络延迟，抢不抢的过看同行力度
        }
    }, 100);

    // 自动点击hCaptcha,使用开发者扩展**NopeCHA**就行

})();