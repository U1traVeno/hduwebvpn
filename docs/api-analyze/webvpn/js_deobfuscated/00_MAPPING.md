# WebVPN JS 代码混淆-反混淆对照表

## 映射原则

1. **符号映射**: 混淆代码中的每个短符号 (`a`, `b`, `c`, `e`, `f`, `g`, `h`, `i`, `l`, `n`, `o`, `r`, `s`, `t`, `u`) 都对应一个具名符号
2. **函数映射**: 每个匿名函数都有具名函数名
3. **对象映射**: 每个混淆的对象属性都有可读的键名
4. **一一对应**: 反混淆后的代码与原始代码结构完全对应，只是符号被替换

---

## 文件 1: authentication-DJLiEriA.js

### 原始混淆代码
```javascript
import{r as t}from"./request-B-CfM8d9.js";function e(){return t({url:"/api/access/authentication/list",method:"get"})}function i(){return t({url:"/api/access/authentication/all",method:"get"})}export{e as a,i as g};
```

### 符号映射表
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `r` | import | `request` | 从 request 模块导入的请求函数 |
| `t` | local var | `request` | `r` 的本地别名 |
| `e` | function | `getAuthenticationList` | 获取认证方式列表 |
| `i` | function | `getAllAuthentication` | 获取所有认证方式 |
| `e as a` | export | `getAuthList` | 导出 e 函数 |
| `i as g` | export | `getAllAuthMethods` | 导出 i 函数 |

### 反混淆代码
```javascript
import { r as request } from "./request-B-CfM8d9.js";

function getAuthenticationList() {
    return request({ url: "/api/access/authentication/list", method: "get" });
}

function getAllAuthentication() {
    return request({ url: "/api/access/authentication/all", method: "get" });
}

export { getAuthenticationList as getAuthList, getAllAuthentication as getAllAuthMethods };
```

---

## 文件 2: authMethod-CtTsQ3_a.js

### 原始混淆代码
```javascript
import{r as t}from"./request-B-CfM8d9.js";function c(a){return t({url:"/api/access/auth-method-bind-local/start",method:"post",data:a})}function s(a){return t({url:"/api/access/auth-method/bind-local-account",method:"post",data:a})}function n(a){return t({url:"/api/access/half-account-bind/start",method:"post",data:a})}function r(a){return t({url:"/api/access/half-account-bind/finish",method:"post",data:a})}export{r as a,s as b,c,n as h};
```

### 符号映射表
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `r` | import | `request` | 从 request 模块导入 |
| `t` | local var | `request` | `r` 的本地别名 |
| `c` | function | `startBindLocalAccount` | 开始绑定本地账号 |
| `s` | function | `bindLocalAccount` | 绑定本地账号 |
| `n` | function | `startHalfAccountBind` | 开始半账号绑定 |
| `r` | function | `finishHalfAccountBind` | 完成半账号绑定 |
| `r as a` | export | `finishHalfBind` | 导出 r 函数 |
| `s as b` | export | `bindLocal` | 导出 s 函数 |
| `c` | export | `startBindLocal` | 导出 c 函数 |
| `n as h` | export | `startHalfBind` | 导出 n 函数 |

### 反混淆代码
```javascript
import { r as request } from "./request-B-CfM8d9.js";

function startBindLocalAccount(data) {
    return request({ url: "/api/access/auth-method-bind-local/start", method: "post", data: data });
}

function bindLocalAccount(data) {
    return request({ url: "/api/access/auth-method/bind-local-account", method: "post", data: data });
}

function startHalfAccountBind(data) {
    return request({ url: "/api/access/half-account-bind/start", method: "post", data: data });
}

function finishHalfAccountBind(data) {
    return request({ url: "/api/access/half-account-bind/finish", method: "post", data: data });
}

export {
    finishHalfAccountBind as finishHalfBind,
    bindLocalAccount as bindLocal,
    startBindLocalAccount,
    startHalfAccountBind as startHalfBind
};
```

---

## 文件 3: urlConf-D5f8gSXt.js

### 原始混淆代码
```javascript
const a=`${window.location.origin}/callback`,c={oauth2CallbackUrl:`${a}/oauth2`,casCallbackUrl:`${a}/cas`,larkCallbackUrl:`${a}/lark`,wechatEnCallbackUrl:`${a}/wechat-enterprise`,wechatEnPcCallbackUrl:`${a}/wechat-en-pc`,dingTalkCallbackUrl:`${a}/ding-talk`,wechatOfficialRzkcPcCallbackUrl:`${a}/wechat-official-rzkc-pc`,wechatOfficialRzkcCallbackUrl:`${a}/wechat-official-rzkc`,alipayPcCallbackUrl:`${a}/alipay-pc`,alipayCallbackUrl:`${a}/alipay`,ecardCallbackUrl:`${a}/ecard`},r=async()=>{const l=`${window.location.origin}/site-nav/`;return`${localStorage.getItem("returnUrl")?localStorage.getItem("returnUrl"):l}${localStorage.getItem("returnUrlHash")?localStorage.getItem("returnUrlHash"):""}`};export{c,r as h};
```

### 符号映射表
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `a` | const | `baseCallbackUrl` | 基础回调 URL |
| `c` | const | `callbackUrls` | 回调 URL 配置对象 |
| `r` | async function | `getReturnUrl` | 获取返回 URL |
| `l` | local const | `defaultReturnUrl` | 默认返回 URL |
| `c` | export | `callbackConfig` | 导出 c 对象 |
| `r as h` | export | `getHomeUrl` | 导出 r 函数 |

### 反混淆代码
```javascript
const baseCallbackUrl = `${window.location.origin}/callback`;

const callbackUrls = {
    oauth2CallbackUrl: `${baseCallbackUrl}/oauth2`,
    casCallbackUrl: `${baseCallbackUrl}/cas`,
    larkCallbackUrl: `${baseCallbackUrl}/lark`,
    wechatEnCallbackUrl: `${baseCallbackUrl}/wechat-enterprise`,
    wechatEnPcCallbackUrl: `${baseCallbackUrl}/wechat-en-pc`,
    dingTalkCallbackUrl: `${baseCallbackUrl}/ding-talk`,
    wechatOfficialRzkcPcCallbackUrl: `${baseCallbackUrl}/wechat-official-rzkc-pc`,
    wechatOfficialRzkcCallbackUrl: `${baseCallbackUrl}/wechat-official-rzkc`,
    alipayPcCallbackUrl: `${baseCallbackUrl}/alipay-pc`,
    alipayCallbackUrl: `${baseCallbackUrl}/alipay`,
    ecardCallbackUrl: `${baseCallbackUrl}/ecard`
};

const getReturnUrl = async () => {
    const defaultReturnUrl = `${window.location.origin}/site-nav/`;
    return `${localStorage.getItem("returnUrl") ? localStorage.getItem("returnUrl") : defaultReturnUrl}${localStorage.getItem("returnUrlHash") ? localStorage.getItem("returnUrlHash") : ""}`;
};

export { callbackUrls as callbackConfig, getReturnUrl as getHomeUrl };
```

---

## 文件 4: navCustom-CgDbtGc9.js

### 原始混淆代码
```javascript
import{V as s}from"./index-Cey6Kqla.js";import{r as t}from"./request-B-CfM8d9.js";function n(){return t({url:"/api/access/nav/site-list",method:"get"})}function c(){return t({url:"/api/access/nav/favorite-sites",method:"post"})}function i(o){return t({url:"/api/access/nav/add-to-favorites",method:"post",data:o})}function u(o){return t({url:"/api/access/nav/remove-from-favorites",method:"post",data:o})}function m(){return t({url:"/api/access/nav/config",method:"get"})}function f(){return t({url:"/api/access/imufe/code",method:"post"})}const g=s("navCustom",{state:()=>({customInfo:{banner:"",copyright:"",title:"网站导航",favicon:"",css:"",logo:"",siteTypeTagColor:"",userNameColor:"",userNamePriority:"username",personalCenterBgColor:"",data:{problemDescription:""}}}),actions:{async getCustomInfo(){const{code:o,data:e}=await t({url:"/api/access/nav/custom",method:"get"});o===0&&e&&(this.customInfo={...e})}},getters:{getBanner:o=>o.customInfo.banner,getCopyright:o=>o.customInfo.copyright,getTitle:o=>o.customInfo.title,getFavicon:o=>o.customInfo.favicon,getCss:o=>o.customInfo.css,getLogo:o=>o.customInfo.logo,getSiteTypeTagColor:o=>o.customInfo.siteTypeTagColor,getUserNameColor:o=>o.customInfo.userNameColor,getUserNamePriority:o=>o.customInfo.userNamePriority?o.customInfo.userNamePriority:"username",getPersonalCenterBgColor:o=>o.customInfo.personalCenterBgColor,getProblemDescription:o=>{var e;return((e=o.customInfo.data)==null?void 0:e.problemDescription)||""}}});export{i as a,n as b,c,m as d,f as g,u as r,g as u};
```

### 符号映射表
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `s` | import | `createStore` | Pinia 的 createStore |
| `t` | import | `request` | 请求函数 |
| `n` | function | `getSiteList` | 获取站点列表 |
| `c` | function | `getFavoriteSites` | 获取收藏站点 |
| `i` | function | `addToFavorites` | 添加收藏 |
| `u` | function | `removeFromFavorites` | 移除收藏 |
| `m` | function | `getNavConfig` | 获取导航配置 |
| `f` | function | `getImufeCode` | 获取 imufe 代码 |
| `g` | const | `navCustomStore` | 导航自定义 Store |
| `o` | param | `data` | 数据参数 |
| `e` | local var | `data` | 响应数据 |
| `i as a` | export | `addToFavorites as addFavorite` | |
| `n as b` | export | `getSiteList as siteList` | |
| `c` | export | `getFavoriteSites as favoriteSites` | |
| `m as d` | export | `getNavConfig as navConfig` | |
| `f as g` | export | `getImufeCode as imufeCode` | |
| `u as r` | export | `removeFromFavorites as removeFavorite` | |
| `g as u` | export | `navCustomStore as navStore` | |

---

## 文件 5: user-cwdmpEuP.js

### 原始混淆代码 (关键部分)
```javascript
import{s as i}from"./index-Cey6Kqla.js";import{r as o}from"./request-B-CfM8d9.js";const s=`accept acceptcharset...`.split(/[\s\n]+/);
function c(e,a){return e.indexOf(a)===0}
function d(e){...}
function p(){return o({url:"/api/access/user/info",method:"get"})}
function m(){return o({url:"/api/access/user/logout",method:"post"})}
function h(e){return o({url:"/api/access/user/change-password",method:"post",data:e})}
...
export{w as a,g as b,h as c,f as d,y as e,p as f,b as g,v as h,m as l,d as p};
```

### 符号映射表
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `i` | import | `nothing` | 未使用 |
| `o` | import | `request` | 请求函数 |
| `s` | const | `HTML_PROPS` | HTML 属性列表 |
| `c` | function | `startsWith` | 字符串前缀检查 |
| `d` | function | `filterProps` | 过滤属性 |
| `p` | function | `getUserInfo` | 获取用户信息 |
| `m` | function | `logout` | 用户登出 |
| `h` | function | `changePassword` | 修改密码 |
| `g` | function | `changeInfo` | 修改用户信息 |
| `f` | function | `changeEmail` | 修改邮箱 |
| `y` | function | `changeMobile` | 修改手机 |
| `v` | function | `changeAvatar` | 修改头像 |
| `b` | function | `getAuthHistory` | 获取认证历史 |
| `w` | function | `getAccessLogList` | 获取访问日志 |
| `w as a` | export | `getAccessLogList as accessLogs` | |
| `g as b` | export | `changeInfo as updateInfo` | |
| `h as c` | export | `changePassword as updatePassword` | |
| `f as d` | export | `changeEmail as updateEmail` | |
| `y as e` | export | `changeMobile as updateMobile` | |
| `p as f` | export | `getUserInfo as userInfo` | |
| `b as g` | export | `getAuthHistory as authHistory` | |
| `v as h` | export | `changeAvatar as updateAvatar` | |
| `m as l` | export | `logout as userLogout` | |
| `d as p` | export | `filterProps as filterHTMLProps` | |

---

## 文件 6: auth-Crs8kf_D.js

### 原始混淆代码
```javascript
import{d as s,c}from"./auth-Crs8kf_D.js";import{b as n,c as t}from"./authMethod-CtTsQ3_a.js";import{u as r}from"./user-V10Yxw5H.js";const o=r(),d=async a=>o.userInfo.needToBindLocalAccount?await t(a):await c(a),f=async a=>o.userInfo.needToBindLocalAccount?await n(a):await s(a);export{f as a,d as h};
```

### 符号映射表
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `d` | import | `finishAuth` | 完成认证 (来自 auth-Crs8kf_D.js) |
| `c` | import | `startAuth` | 开始认证 (来自 auth-Crs8kf_D.js) |
| `n` | import | `bindLocalAccount` | 绑定本地账号 |
| `t` | import | `startBindLocalAccount` | 开始绑定本地账号 |
| `u` | import | `loginStore` | 登录 Store |
| `r` | local var | `loginStore` | `u` 的别名 |
| `o` | const | `storeInstance` | store 实例 |
| `d` | async function | `handleAuthWithLocal` | 处理需要绑定本地账号的认证 |
| `f` | async function | `handleAuthWithoutLocal` | 处理不需要绑定本地账号的认证 |
| `a` | param | `credentials` | 认证凭据 |
| `f as a` | export | `handleAuthWithoutLocal as auth` | |
| `d as h` | export | `handleAuthWithLocal as authWithBind` | |

### 反混淆代码
```javascript
import { d as finishAuth, c as startAuth } from "./auth-Crs8kf_D.js";
import { b as bindLocalAccount, n as startBindLocalAccount } from "./authMethod-CtTsQ3_a.js";
import { u as loginStore } from "./user-V10Yxw5H.js";

const storeInstance = loginStore();

const handleAuthWithLocal = async (credentials) => {
    return storeInstance.userInfo.needToBindLocalAccount
        ? await startBindLocalAccount(credentials)
        : await startAuth(credentials);
};

const handleAuthWithoutLocal = async (credentials) => {
    return storeInstance.userInfo.needToBindLocalAccount
        ? await bindLocalAccount(credentials)
        : await finishAuth(credentials);
};

export {
    handleAuthWithoutLocal as auth,
    handleAuthWithLocal as authWithBind
};
```

---

## 文件 7: auth-Bvec_zE4.js

### 原始混淆代码
```javascript
import{d as s,c}from"./auth-Crs8kf_D.js";import{b as n,c as t}from"./authMethod-CtTsQ3_a.js";import{u as r}from"./user-V10Yxw5H.js";const o=r(),d=async a=>o.userInfo.needToBindLocalAccount?await t(a):await c(a),f=async a=>o.userInfo.needToBindLocalAccount?await n(a):await s(a);export{f as a,d as h};
```

### 符号映射表
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `d` | import | `finishAuth` | 完成认证 |
| `c` | import | `startAuth` | 开始认证 |
| `n` | import | `bindLocalAccount` | 绑定本地账号 |
| `t` | import | `startBindLocalAccount` | 开始绑定本地账号 |
| `u` | import | `loginStore` | 登录 Store |
| `r` | local var | `loginStore` | `u` 的别名 |
| `o` | const | `storeInstance` | store 实例 |
| `d` | async function | `authWithLocal` | 需要本地账号的认证 |
| `f` | async function | `authWithoutLocal` | 不需要本地账号的认证 |
| `f as a` | export | `authWithoutLocal as authenticate` | |
| `d as h` | export | `authWithLocal as authenticateWithBinding` | |

---

## 文件 8: CasCallback-q97cqoLl.js

### 原始混淆代码 (关键部分)
```javascript
import{a as l}from"./auth-Bvec_zE4.js";import{_ as u,d as f,c as g,o as h,f as w,h as I,B as k,k as y,u as b,j as v}from"./index-Cey6Kqla.js";import{u as L}from"./user-V10Yxw5H.js";import{g as U}from"./util-DyruD4Ub.js";import{S as x}from"./index-o4p_ZbnI.js";import{b as B}from"./request-B-CfM8d9.js";import"./auth-Crs8kf_D.js";import"./authMethod-CtTsQ3_a.js";import"./user-cwdmpEuP.js";import"./const-BHerojsC.js";import"./urlConf-D5f8gSXt.js";import"./styleChecker-TFxV7Lzl.js";import"./index-DXhXwGTY.js";import"./initDefaultProps-_hAYIFu4.js";import"./fade-BdgfkLw7.js";const C=f({name:"CasCallback",setup(){const s=b(),r=y(),t=L(),a=I({authLoading:!1});return k(async()=>{var e;const o=await t.getUserInfo({unbind:!0});if(!((e=o==null?void 0:o.data)!=null&&e.userId))if(a.authLoading=!0,r.query.ticket){const i=decodeURIComponent(r.query.ticket),n=s.currentRoute.value.params.externalId,c=window.location.href.split("?")[0],d={externalId:n,data:JSON.stringify({callbackUrl:c,ticket:i,deviceId:await U()})};await(async m=>{try{const{code:p}=await l(m);p===0&&(a.authLoading=!1,t.getUserInfo(),console.log("userStore.showBindWechat",t.showBindWechat))}finally{a.authLoading=!1}})(d)}else B.error("认证失败")}),{...v(a)}}}),S={class:"callback-container"},K=u(C,[["render",function(s,r,t,a,o,e){const i=x;return h(),g("div",S,[w(i,{tip:"授权中...",size:"large"})])}],["__scopeId","data-v-571c06a7"]]);export{K as default};
```

### 重要符号映射
| 混淆符号 | 类型 | 反混淆名称 | 说明 |
|---------|------|-----------|------|
| `l` | import | `authenticate` | 来自 auth-Bvec_zE4 的认证函数 |
| `u` | import | `defineComponent` | Vue defineComponent |
| `f` | import | `withDefaults` | Vue withDefaults |
| `c` | import | `createElementVNode` | h() 函数 |
| `h` | import | `createBaseVNode` | 创建普通元素 |
| `w` | import | `Loading` | 加载组件 |
| `I` | import | `ref` | Vue ref |
| `k` | import | `onMounted` | Vue onMounted |
| `y` | import | `useRoute` | Vue useRoute |
| `b` | import | `useRouter` | Vue useRouter |
| `v` | import | `toRefs` | Vue toRefs |
| `L` | import | `loginStore` | 来自 user-V10Yxw5H |
| `U` | import | `generateDeviceId` | 来自 util-DyruD4Ub |
| `x` | import | `Spin` | 加载 Spin 组件 |
| `B` | import | `message` | 全局消息 API |
| `s` | local const | `router` | useRouter() |
| `r` | local const | `route` | useRoute() |
| `t` | local const | `store` | loginStore() |
| `a` | local const | `state` | ref({ authLoading: false }) |
| `i` | local var | `ticket` | URL 中的 ticket |
| `n` | local var | `externalId` | 路由参数中的 externalId |
| `c` | local var | `callbackUrl` | 当前页面 URL（不含 query） |
| `d` | local const | `authPayload` | 认证请求数据 |
| `m` | param | `payload` | 认证载荷 |

### 关键代码逻辑
```javascript
// 1. 从路由参数获取 externalId (这是关键发现！)
const externalId = route.currentRoute.value.params.externalId;

// 2. 从 URL query 获取 ticket
const ticket = decodeURIComponent(route.query.ticket);

// 3. 生成回调 URL (不含 query)
const callbackUrl = window.location.href.split("?")[0];

// 4. 生成设备 ID
const deviceId = await generateDeviceId();

// 5. 组装认证数据
const authPayload = {
    externalId: externalId,
    data: JSON.stringify({
        callbackUrl: callbackUrl,
        ticket: ticket,
        deviceId: deviceId
    })
};

// 6. 调用认证 API
const { code } = await authenticate(authPayload);
if (code === 0) {
    authLoading = false;
    store.getUserInfo();
}
```

---

## 总结：认证流程的真实顺序

根据 `CasCallback-q97cqoLl.js` 的源码，完整的 CAS 回调流程是：

```
1. 用户从 SSO 重定向回来: /callback/cas/:externalId?ticket=ST-xxx
         ↓
2. Vue 组件 mounted
         ↓
3. 调用 store.getUserInfo({ unbind: true }) 检查是否已登录
         ↓
4. 如果 userId 不存在 且 URL 中有 ticket:
   4a. 从路由参数获取 externalId
   4b. 从 URL query 获取 ticket
   4c. 生成 callbackUrl 和 deviceId
   4d. 调用 authenticate(payload) → POST /api/access/auth/finish
         ↓
5. 认证成功后:
   5a. 关闭 loading
   5b. 调用 store.getUserInfo() 获取完整用户信息
```
