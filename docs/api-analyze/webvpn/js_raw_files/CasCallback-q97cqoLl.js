// HTTP/1.1 200 OK
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:48:44 GMT
// Etag: W/"6968b84e-599"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 1433

import{a as l}from"./auth-Bvec_zE4.js";import{_ as u,d as f,c as g,o as h,f as w,h as I,B as k,k as y,u as b,j as v}from"./index-Cey6Kqla.js";import{u as L}from"./user-V10Yxw5H.js";import{g as U}from"./util-DyruD4Ub.js";import{S as x}from"./index-o4p_ZbnI.js";import{b as B}from"./request-B-CfM8d9.js";import"./auth-Crs8kf_D.js";import"./authMethod-CtTsQ3_a.js";import"./user-cwdmpEuP.js";import"./const-BHerojsC.js";import"./urlConf-D5f8gSXt.js";import"./styleChecker-TFxV7Lzl.js";import"./index-DXhXwGTY.js";import"./initDefaultProps-_hAYIFu4.js";import"./fade-BdgfkLw7.js";const C=f({name:"CasCallback",setup(){const s=b(),r=y(),t=L(),a=I({authLoading:!1});return k(async()=>{var e;const o=await t.getUserInfo({unbind:!0});if(!((e=o==null?void 0:o.data)!=null&&e.userId))if(a.authLoading=!0,r.query.ticket){const i=decodeURIComponent(r.query.ticket),n=s.currentRoute.value.params.externalId,c=window.location.href.split("?")[0],d={externalId:n,data:JSON.stringify({callbackUrl:c,ticket:i,deviceId:await U()})};await(async m=>{try{const{code:p}=await l(m);p===0&&(a.authLoading=!1,t.getUserInfo(),console.log("userStore.showBindWechat",t.showBindWechat))}finally{a.authLoading=!1}})(d)}else B.error("认证失败")}),{...v(a)}}}),S={class:"callback-container"},K=u(C,[["render",function(s,r,t,a,o,e){const i=x;return h(),g("div",S,[w(i,{tip:"授权中...",size:"large"})])}],["__scopeId","data-v-571c06a7"]]);export{K as default};
