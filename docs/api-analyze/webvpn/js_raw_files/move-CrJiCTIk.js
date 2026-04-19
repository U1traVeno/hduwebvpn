// HTTP/1.1 200 OK
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:48:45 GMT
// Etag: W/"6968b84e-7bd"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 1981

import{K as r}from"./request-B-CfM8d9.js";import{i as s}from"./styleChecker-TFxV7Lzl.js";const m=new r("antMoveDownIn",{"0%":{transform:"translate3d(0, 100%, 0)",transformOrigin:"0 0",opacity:0},"100%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1}}),f=new r("antMoveDownOut",{"0%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1},"100%":{transform:"translate3d(0, 100%, 0)",transformOrigin:"0 0",opacity:0}}),y=new r("antMoveLeftIn",{"0%":{transform:"translate3d(-100%, 0, 0)",transformOrigin:"0 0",opacity:0},"100%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1}}),p=new r("antMoveLeftOut",{"0%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1},"100%":{transform:"translate3d(-100%, 0, 0)",transformOrigin:"0 0",opacity:0}}),c=new r("antMoveRightIn",{"0%":{transform:"translate3d(100%, 0, 0)",transformOrigin:"0 0",opacity:0},"100%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1}}),O=new r("antMoveRightOut",{"0%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1},"100%":{transform:"translate3d(100%, 0, 0)",transformOrigin:"0 0",opacity:0}}),g={"move-up":{inKeyframes:new r("antMoveUpIn",{"0%":{transform:"translate3d(0, -100%, 0)",transformOrigin:"0 0",opacity:0},"100%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1}}),outKeyframes:new r("antMoveUpOut",{"0%":{transform:"translate3d(0, 0, 0)",transformOrigin:"0 0",opacity:1},"100%":{transform:"translate3d(0, -100%, 0)",transformOrigin:"0 0",opacity:0}})},"move-down":{inKeyframes:m,outKeyframes:f},"move-left":{inKeyframes:y,outKeyframes:p},"move-right":{inKeyframes:c,outKeyframes:O}},u=(t,n)=>{const{antCls:o}=t,a=`${o}-${n}`,{inKeyframes:i,outKeyframes:e}=g[n];return[s(a,i,e,t.motionDurationMid),{[`
        ${a}-enter,
        ${a}-appear
      `]:{opacity:0,animationTimingFunction:t.motionEaseOutCirc},[`${a}-leave`]:{animationTimingFunction:t.motionEaseInOutCirc}}]};export{u as i};
