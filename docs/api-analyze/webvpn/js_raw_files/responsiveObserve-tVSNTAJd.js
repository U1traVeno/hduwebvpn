// HTTP/1.1 200 OK
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:47:32 GMT
// Etag: W/"6968b84e-41f"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 1055

import{x as d,s as a}from"./index-Cey6Kqla.js";import{F as o}from"./request-B-CfM8d9.js";const w=["xxxl","xxl","xl","lg","md","sm","xs"];function X(){const[,l]=o();return d(()=>{const i=(e=>({xs:`(max-width: ${e.screenXSMax}px)`,sm:`(min-width: ${e.screenSM}px)`,md:`(min-width: ${e.screenMD}px)`,lg:`(min-width: ${e.screenLG}px)`,xl:`(min-width: ${e.screenXL}px)`,xxl:`(min-width: ${e.screenXXL}px)`,xxxl:`{min-width: ${e.screenXXXL}px}`}))(l.value),t=new Map;let c=-1,n={};return{matchHandlers:{},dispatch:e=>(n=e,t.forEach(r=>r(n)),t.size>=1),subscribe(e){return t.size||this.register(),c+=1,t.set(c,e),e(n),c},unsubscribe(e){t.delete(e),t.size||this.unregister()},unregister(){Object.keys(i).forEach(e=>{const r=i[e],s=this.matchHandlers[r];s==null||s.mql.removeListener(s==null?void 0:s.listener)}),t.clear()},register(){Object.keys(i).forEach(e=>{const r=i[e],s=m=>{let{matches:h}=m;this.dispatch(a(a({},n),{[e]:h}))},x=window.matchMedia(r);x.addListener(s),this.matchHandlers[r]={mql:x,listener:s},s(x)})},responsiveMap:i}})}export{w as r,X as u};
