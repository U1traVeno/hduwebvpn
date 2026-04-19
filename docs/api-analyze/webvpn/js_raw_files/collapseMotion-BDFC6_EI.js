// HTTP/1.1 200 OK
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:47:32 GMT
// Etag: W/"6968b84e-6f1"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 1777

import{a3 as h}from"./index-Cey6Kqla.js";function g(e,n,t,s){for(var i=e.length,o=t+-1;++o<i;)if(n(e[o],o,e))return o;return-1}function p(e){return e!=e}function $(e,n){return!!(e!=null&&e.length)&&function(t,s,i){return s==s?function(o,u,m){for(var a=m-1,f=o.length;++a<f;)if(o[a]===u)return a;return-1}(t,s,i):g(t,p,i)}(e,n,0)>-1}const d=e=>({[e.componentCls]:{[`${e.antCls}-motion-collapse-legacy`]:{overflow:"hidden","&-active":{transition:`height ${e.motionDurationMid} ${e.motionEaseInOut},
        opacity ${e.motionDurationMid} ${e.motionEaseInOut} !important`}},[`${e.antCls}-motion-collapse`]:{overflow:"hidden",transition:`height ${e.motionDurationMid} ${e.motionEaseInOut},
        opacity ${e.motionDurationMid} ${e.motionEaseInOut} !important`}}});function c(e,n){return e.classList?e.classList.contains(n):` ${e.className} `.indexOf(` ${n} `)>-1}function l(e,n){e.classList?e.classList.add(n):c(e,n)||(e.className=`${e.className} ${n}`)}function r(e,n){if(e.classList)e.classList.remove(n);else if(c(e,n)){const t=e.className;e.className=` ${t} `.replace(` ${n} `," ")}}const v=function(){let e=arguments.length>0&&arguments[0]!==void 0?arguments[0]:"ant-motion-collapse";return{name:e,appear:!(arguments.length>1&&arguments[1]!==void 0)||arguments[1],css:!0,onBeforeEnter:n=>{n.style.height="0px",n.style.opacity="0",l(n,e)},onEnter:n=>{h(()=>{n.style.height=`${n.scrollHeight}px`,n.style.opacity="1"})},onAfterEnter:n=>{n&&(r(n,e),n.style.height=null,n.style.opacity=null)},onBeforeLeave:n=>{l(n,e),n.style.height=`${n.offsetHeight}px`,n.style.opacity=null},onLeave:n=>{setTimeout(()=>{n.style.height="0px",n.style.opacity="0"})},onAfterLeave:n=>{n&&(r(n,e),n.style&&(n.style.height=null,n.style.opacity=null))}}};export{$ as a,g as b,v as c,l as d,d as g,r};
