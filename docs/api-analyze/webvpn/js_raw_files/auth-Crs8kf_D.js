// HTTP/1.1 200 OK
// Accept-Ranges: bytes
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:47:32 GMT
// Etag: "6968b84e-33d"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 829

import{r as a}from"./request-B-CfM8d9.js";function e(t){return a({url:"/api/access/auth/start",method:"post",data:t})}function u(t){return a({url:"/api/access/auth/finish",method:"post",data:t})}function o(t){return a({url:"/api/access/auth/tfa",method:"post",data:t})}function r(){return a({url:"/api/access/auth/tfa-config",method:"post"})}function n(){return a({url:"/api/access/auth/user-notice-info",method:"get"})}function c(t){return a({url:"/api/access/auth/consume-session",method:"get",params:t})}function i(t){return a({url:"/api/access/auth/reset-password",method:"post",data:t})}function h(t){return a({url:"/api/access/auth/wechat-log",method:"post",data:t})}function p(t){return a({url:"/api/access/auth/session-token",method:"post",data:t})}export{o as a,r as b,e as c,u as d,c as e,h as f,n as g,i as r,p as s};
