// HTTP/1.1 200 OK
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:48:45 GMT
// Etag: W/"6968b84e-52e"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 1326

import{u as s}from"./navCustom-CgDbtGc9.js";import{u as m}from"./user-V10Yxw5H.js";import{_ as c,d as u,c as l,b as d,o as p,f,r as g,h,W as y,j as w}from"./index-Cey6Kqla.js";import"./request-B-CfM8d9.js";import"./user-cwdmpEuP.js";import"./const-BHerojsC.js";import"./urlConf-D5f8gSXt.js";import"./util-DyruD4Ub.js";import"./styleChecker-TFxV7Lzl.js";import"./index-DXhXwGTY.js";import"./initDefaultProps-_hAYIFu4.js";import"./fade-BdgfkLw7.js";const I=u({name:"SiteNav",setup(){const o=m(),e=s(),n=h({});return y(async()=>{var r;if(await o.getUserInfo({onlyFetchUserInfo:!0}),(r=o.userInfo)==null?void 0:r.userId){if(await e.getCustomInfo(),e.getTitle&&(document.title=e.getTitle),e.getFavicon){let t=document.querySelector("link[rel*='icon']")||document.createElement("link");t.type="image/x-icon",t.rel="shortcut icon",t.href=e.getFavicon,document.getElementsByTagName("head")[0].appendChild(t)}if(e!=null&&e.getCss){let t=document.createElement("style");t.innerHTML=e==null?void 0:e.getCss,document.getElementsByTagName("head")[0].appendChild(t)}}else window.location.href=`${window.location.origin}/auth/`}),{...w(n),userStore:o}}}),v={key:0},H=c(I,[["render",function(o,e,n,r,t,C){var i;const a=g("router-view");return(i=o.userStore.userInfo)!=null&&i.userId?(p(),l("div",v,[f(a)])):d("",!0)}]]);export{H as default};
