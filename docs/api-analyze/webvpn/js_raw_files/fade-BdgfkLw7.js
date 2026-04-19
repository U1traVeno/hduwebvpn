// HTTP/1.1 200 OK
// Accept-Ranges: bytes
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:47:32 GMT
// Etag: "6968b84e-226"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 550

import{K as o}from"./request-B-CfM8d9.js";import{i as r}from"./styleChecker-TFxV7Lzl.js";const m=new o("antFadeIn",{"0%":{opacity:0},"100%":{opacity:1}}),c=new o("antFadeOut",{"0%":{opacity:1},"100%":{opacity:0}}),u=function(a){let i=arguments.length>1&&arguments[1]!==void 0&&arguments[1];const{antCls:e}=a,n=`${e}-fade`,t=i?"&":"";return[r(n,m,c,a.motionDurationMid,i),{[`
        ${t}${n}-enter,
        ${t}${n}-appear
      `]:{opacity:0,animationTimingFunction:"linear"},[`${t}${n}-leave`]:{animationTimingFunction:"linear"}}]};export{u as i};
