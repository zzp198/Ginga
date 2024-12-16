// ==UserScript==
// @name         B站广告屏蔽
// @namespace    http://tampermonkey.net/
// @version      2024-12-14
// @description  try to take over the world!
// @author       You
// @match        *://*.bilibili.com/*
// @require      https://code.jquery.com/jquery-3.6.0.min.js
// @icon         https://www.google.com/s2/favicons?sz=64&domain=bilibili.com
// @grant        none
// ==/UserScript==

(function () {
    'use strict';

    setInterval(function () {

        $(".slide-ad-exp").each(function () {
            $(this).remove();
            console.log("移除播放页广告");
        });

        $(".video-card-ad-small").each(function () {
            $(this).remove();
            console.log("移除播放页广告");
        });

        $(".bili-video-card__stats--text").each(function () {
            if ($(this).text() === "广告") {
                $(this).closest(".bili-video-card").remove();
                console.log("移除卡片广告");
            }
        });

        $(".vui_icon").each(function () {
            $(this).closest(".bili-video-card").remove();
            console.log("移除建模等广告");
        });

    }, 1000);
})();