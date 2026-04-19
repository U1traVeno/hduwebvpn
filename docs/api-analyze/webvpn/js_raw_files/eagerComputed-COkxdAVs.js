// HTTP/1.1 200 OK
// Accept-Ranges: bytes
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:48:45 GMT
// Etag: "6968b84e-158"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 344

import{u as e}from"./responsiveObserve-tVSNTAJd.js";import{E as n,B as o,C as t,aj as c}from"./index-Cey6Kqla.js";function f(){const u=n({});let s=null;const a=e();return o(()=>{s=a.value.subscribe(r=>{u.value=r})}),t(()=>{a.value.unsubscribe(s)}),u}function b(u){const s=n();return c(()=>{s.value=u()},{flush:"sync"}),s}export{b as e,f as u};
