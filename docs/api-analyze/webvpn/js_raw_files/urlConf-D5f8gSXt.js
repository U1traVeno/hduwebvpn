// HTTP/1.1 200 OK
// Accept-Ranges: bytes
// Content-Type: application/javascript
// Date: Sun, 19 Apr 2026 08:47:32 GMT
// Etag: "6968b84e-2ce"
// Last-Modified: Thu, 15 Jan 2026 09:50:06 GMT
// Server: nginx
// Content-Length: 718

const a=`${window.location.origin}/callback`,c={oauth2CallbackUrl:`${a}/oauth2`,casCallbackUrl:`${a}/cas`,larkCallbackUrl:`${a}/lark`,wechatEnCallbackUrl:`${a}/wechat-enterprise`,wechatEnPcCallbackUrl:`${a}/wechat-en-pc`,dingTalkCallbackUrl:`${a}/ding-talk`,wechatOfficialRzkcPcCallbackUrl:`${a}/wechat-official-rzkc-pc`,wechatOfficialRzkcCallbackUrl:`${a}/wechat-official-rzkc`,alipayPcCallbackUrl:`${a}/alipay-pc`,alipayCallbackUrl:`${a}/alipay`,ecardCallbackUrl:`${a}/ecard`},r=async()=>{const l=`${window.location.origin}/site-nav/`;return`${localStorage.getItem("returnUrl")?localStorage.getItem("returnUrl"):l}${localStorage.getItem("returnUrlHash")?localStorage.getItem("returnUrlHash"):""}`};export{c,r as h};
