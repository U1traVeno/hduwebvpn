// HTTP/1.1 200 OK
// Accept-Ranges: bytes
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:48:45 GMT
// Etag: "6968b84e-35d"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 861

import{w as f}from"./index-DXhXwGTY.js";function s(o){return o!=null&&o===o.window}function D(o,a){var c,l;if(typeof window>"u")return 0;const e="scrollTop";let t=0;return s(o)?t=o.scrollY:o instanceof Document?t=o.documentElement[e]:(o instanceof HTMLElement||o)&&(t=o[e]),o&&!s(o)&&typeof t!="number"&&(t=(l=((c=o.ownerDocument)!==null&&c!==void 0?c:o).documentElement)===null||l===void 0?void 0:l[e]),t}function E(o){let a=arguments.length>1&&arguments[1]!==void 0?arguments[1]:{};const{getContainer:c=()=>window,callback:l,duration:e=450}=a,t=c(),w=D(t),p=Date.now(),m=()=>{const r=Date.now()-p,u=function(n,i,g,v){const d=g-i;return(n/=v/2)<1?d/2*n*n*n+i:d/2*((n-=2)*n*n+2)+i}(r>e?e:r,w,o,e);s(t)?t.scrollTo(window.scrollX,u):t instanceof Document?t.documentElement.scrollTop=u:t.scrollTop=u,r<e?f(m):typeof l=="function"&&l()};f(m)}export{D as g,E as s};
