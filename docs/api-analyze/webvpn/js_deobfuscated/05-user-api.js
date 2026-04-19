/**
 * 文件: user-cwdmpEuP.js → 05-user-api.js
 * 功能: 用户相关 API
 *
 * 注意: 导入 index-Cey6Kqla.js (Vue 框架) 和 request.js
 */

import { s as nothing } from "./index-Cey6Kqla.js";
import { r as request } from "./request.js";

const HTML_PROPS = `accept acceptcharset accesskey action allowfullscreen allowtransparency
alt async autocomplete autofocus autoplay capture cellpadding cellspacing challenge
charset checked classid classname colspan cols content contenteditable contextmenu
controls coords crossorigin data datetime default defer dir disabled download draggable
enctype form formaction formenctype formmethod formnovalidate formtarget frameborder
headers height hidden high href hreflang htmlfor for httpequiv icon id inputmode integrity
is keyparams keytype kind label lang list loop low manifest marginheight marginwidth max maxlength media
mediagroup method min minlength multiple muted name novalidate nonce open
optimum pattern placeholder poster preload radiogroup readonly rel required
reversed role rowspan rows sandbox scope scoped scrolling seamless selected
shape size sizes span spellcheck src srcdoc srclang srcset start step style
summary tabindex target title type usemap value width wmode wrap onCopy onCut onPaste onCompositionend onCompositionstart onCompositionupdate onKeydown
    onKeypress onKeyup onFocus onBlur onChange onInput onSubmit onClick onContextmenu onDoubleclick onDblclick
    onDrag onDragend onDragenter onDragexit onDragleave onDragover onDragstart onDrop onMousedown
    onMouseenter onMouseleave onMousemove onMouseout onMouseover onMouseup onSelect onTouchcancel
    onTouchend onTouchmove onTouchstart onTouchstartPassive onTouchmovePassive onScroll onWheel onAbort onCanplay onCanplaythrough
    onDurationchange onEmptied onEncrypted onEnded onError onLoadeddata onLoadedmetadata
    onLoadstart onPause onPlay onPlaying onProgress onRatechange onSeeked onSeeking onStalled onSuspend onTimeupdate onVolumechange onWaiting onLoad onError`.split(/[\s\n]+/);

function startsWith(str, prefix) {
    return str.indexOf(prefix) === 0;
}

function filterProps(props, mode) {
    let filter = mode === false
        ? { aria: true, data: true, attr: true }
        : mode === true
            ? { aria: true }
            : { aria: true, data: true, attr: true };

    const result = {};
    Object.keys(props).forEach(key => {
        if (
            (filter.aria && (key === "role" || startsWith(key, "aria-"))) ||
            (filter.data && startsWith(key, "data-")) ||
            (filter.attr && (HTML_PROPS.includes(key) || HTML_PROPS.includes(key.toLowerCase())))
        ) {
            result[key] = props[key];
        }
    });
    return result;
}

function getUserInfo() {
    return request({ url: "/api/access/user/info", method: "get" });
}

function logout() {
    return request({ url: "/api/access/user/logout", method: "post" });
}

function changePassword(data) {
    return request({ url: "/api/access/user/change-password", method: "post", data: data });
}

function changeInfo(data) {
    return request({ url: "/api/access/user/change-info", method: "post", data: data });
}

function changeEmail(data) {
    return request({ url: "/api/access/user/change-email", method: "post", data: data });
}

function changeMobile(data) {
    return request({ url: "/api/access/user/change-mobile", method: "post", data: data });
}

function changeAvatar(data) {
    return request({ url: "/api/access/user/change-avatar", method: "post", data: data });
}

function getAuthHistory(params) {
    return request({ url: "/api/access/user/auth/history", method: "get", params: params });
}

function getAccessLogList(params) {
    return request({ url: "/api/access/access-log/list", method: "get", params: params });
}

export {
    getAccessLogList as accessLogs,
    changeInfo as updateInfo,
    changePassword as updatePassword,
    changeEmail as updateEmail,
    changeMobile as updateMobile,
    getUserInfo as userInfo,
    getAuthHistory as authHistory,
    changeAvatar as updateAvatar,
    logout as userLogout,
    filterProps as filterHTMLProps
};
