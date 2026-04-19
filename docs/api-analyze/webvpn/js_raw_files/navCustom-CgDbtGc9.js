// HTTP/1.1 200 OK
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:48:45 GMT
// Etag: W/"6968b84e-602"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 1538

import{V as s}from"./index-Cey6Kqla.js";import{r as t}from"./request-B-CfM8d9.js";function n(){return t({url:"/api/access/nav/site-list",method:"get"})}function c(){return t({url:"/api/access/nav/favorite-sites",method:"post"})}function i(o){return t({url:"/api/access/nav/add-to-favorites",method:"post",data:o})}function u(o){return t({url:"/api/access/nav/remove-from-favorites",method:"post",data:o})}function m(){return t({url:"/api/access/nav/config",method:"get"})}function f(){return t({url:"/api/access/imufe/code",method:"post"})}const g=s("navCustom",{state:()=>({customInfo:{banner:"",copyright:"",title:"网站导航",favicon:"",css:"",logo:"",siteTypeTagColor:"",userNameColor:"",userNamePriority:"username",personalCenterBgColor:"",data:{problemDescription:""}}}),actions:{async getCustomInfo(){const{code:o,data:e}=await t({url:"/api/access/nav/custom",method:"get"});o===0&&e&&(this.customInfo={...e})}},getters:{getBanner:o=>o.customInfo.banner,getCopyright:o=>o.customInfo.copyright,getTitle:o=>o.customInfo.title,getFavicon:o=>o.customInfo.favicon,getCss:o=>o.customInfo.css,getLogo:o=>o.customInfo.logo,getSiteTypeTagColor:o=>o.customInfo.siteTypeTagColor,getUserNameColor:o=>o.customInfo.userNameColor,getUserNamePriority:o=>o.customInfo.userNamePriority?o.customInfo.userNamePriority:"username",getPersonalCenterBgColor:o=>o.customInfo.personalCenterBgColor,getProblemDescription:o=>{var e;return((e=o.customInfo.data)==null?void 0:e.problemDescription)||""}}});export{i as a,n as b,c,m as d,f as g,u as r,g as u};
