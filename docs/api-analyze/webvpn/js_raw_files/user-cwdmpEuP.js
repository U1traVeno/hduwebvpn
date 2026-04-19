// HTTP/1.1 200 OK
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:47:32 GMT
// Etag: W/"6968b84e-bbf"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 3007

import{s as i}from"./index-Cey6Kqla.js";import{r as o}from"./request-B-CfM8d9.js";const s=`accept acceptcharset accesskey action allowfullscreen allowtransparency
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
    onLoadstart onPause onPlay onPlaying onProgress onRatechange onSeeked onSeeking onStalled onSuspend onTimeupdate onVolumechange onWaiting onLoad onError`.split(/[\s\n]+/);function c(e,a){return e.indexOf(a)===0}function d(e){let a,t=arguments.length>1&&arguments[1]!==void 0&&arguments[1];a=t===!1?{aria:!0,data:!0,attr:!0}:t===!0?{aria:!0}:i({},t);const r={};return Object.keys(e).forEach(n=>{(a.aria&&(n==="role"||c(n,"aria-"))||a.data&&c(n,"data-")||a.attr&&(s.includes(n)||s.includes(n.toLowerCase())))&&(r[n]=e[n])}),r}function p(){return o({url:"/api/access/user/info",method:"get"})}function m(){return o({url:"/api/access/user/logout",method:"post"})}function h(e){return o({url:"/api/access/user/change-password",method:"post",data:e})}function g(e){return o({url:"/api/access/user/change-info",method:"post",data:e})}function f(e){return o({url:"/api/access/user/change-email",method:"post",data:e})}function y(e){return o({url:"/api/access/user/change-mobile",method:"post",data:e})}function v(e){return o({url:"/api/access/user/change-avatar",method:"post",data:e})}function b(e){return o({url:"/api/access/user/auth/history",method:"get",params:e})}function w(e){return o({url:"/api/access/access-log/list",method:"get",params:e})}export{w as a,g as b,h as c,f as d,y as e,p as f,b as g,v as h,m as l,d as p};
