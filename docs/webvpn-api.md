# WebVPN API 分析报告

## 概述

本报告基于 `webvpn.xlsx` 中保存的 HTTP 请求包流量分析，涵盖了杭州电子科技大学 WebVPN 系统（`webvpn.hdu.edu.cn`）的 API 结构、认证流程和请求封装机制。

**数据来源**：128 条 HTTP 请求/响应记录  
**主要域名**：

- `webvpn.hdu.edu.cn` - WebVPN 前端和 API（69 条记录）
- `sso.hdu.edu.cn` - 统一身份认证服务（58 条记录）
- `https-cas-hdu-edu-cn-443.webvpn.hdu.edu.cn` - CAS 回调（1 条记录）

---

## 一、关键 REST API 调用分析

### 1.1 认证相关 API

#### 获取认证方式列表

```
GET /api/access/authentication/list
```

返回可用的认证方式，当前仅配置了统一身份认证（CAS）：

```json
{
  "code": 0,
  "data": {
    "list": [{
      "externalId": "57AAnALn",
      "name": "统一身份认证",
      "authType": 4,
      "authOptions": {
        "embed": false,
        "forbidAutoJump": false,
        "forceAutoJump": true
      }
    }]
  }
}
```

**前端调用位置**：`index-CscyhzoS.js` (Login 页面组件)

关键代码逻辑：

```javascript
import { g as Y, a as Z } from "./authentication-DJLiEriA.js";
// ...
const { data: d } = await Z(); // 调用 getAuthList()
// d.list[0].externalId → 设置为 authId
t.authId = t.authCommonList[0].externalId;
```

**符号对应关系**：

| 混淆符号 | 原始函数/变量 | 文件 |
|----------|---------------|------|
| `Z` (as `a`) | `getAuthList` | `authentication-DJLiEriA.js` |
| `Y` (as `g`) | `getAllAuthMethods` | `authentication-DJLiEriA.js` |
| `v()` | 自动选择认证方式逻辑 | `index-CscyhzoS.js` |
| `f()` | 处理认证列表分类 | `index-CscyhzoS.js` |

#### externalId 的双重来源

`externalId` 在前端有两处来源：

1. **API 响应获取** (`/api/access/authentication/list`)：
   - 用于构建 SSO 登录 URL
   - `POST /api/access/auth/start` 请求体中传入

2. **URL 路由参数获取**：
   - 用于 CAS 回调验证
   - 格式：`/callback/cas/:externalId?ticket=ST-xxx`
   - 获取方式：`route.currentRoute.value.params.externalId`

关键代码（`08-cas-callback.js`）：

```javascript
const ticket = decodeURIComponent(route.query.ticket);
const externalId = route.currentRoute.value.params.externalId;
```

**符号对应关系**：

| 混淆符号 | 原始函数/变量 | 文件 |
|----------|---------------|------|
| `CasCallback` | CAS 回调组件 | `CasCallback-q97cqoLl.js` |
| `authenticate` (导入) | `authenticate` | `07-auth-bvec.js` |
| `generateDeviceId` (导入) | `generateDeviceId` | `util-DyruD4Ub.js` |

#### 发起认证

```
POST /api/access/auth/start
```

请求体：

```json
{
  "externalId": "57AAnALn",
  "data": "{\"callbackUrl\":\"https://webvpn.hdu.edu.cn/callback/cas/57AAnALn\"}"
}
```

响应：

```json
{
  "code": 0,
  "data": {
    "type": 2,
    "action": {
      "login_url": "https://sso.hdu.edu.cn/login?service=https://webvpn.hdu.edu.cn/callback/cas/57AAnALn"
    }
  }
}
```

#### 完成认证（获取 Token）

```
POST /api/access/auth/finish
```

请求体：

```json
{
  "externalId": "57AAnALn",
  "data": "{\"callbackUrl\":\"https://webvpn.hdu.edu.cn/callback/cas/57AAnALn\",\"ticket\":\"ST-1816171-Nr2J8hMaNdbwZiN-1Xba-c0uTaIrg-sso-64ccf764df-56tht\",\"deviceId\":\"8a5541b34915d358d5842ec56d040ef6\"}"
}
```

响应（包含 `webvpn-token`）：

```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

响应头中同时设置 Cookie：

```
Set-Cookie: webvpn-token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...; Path=/; Domain=hdu.edu.cn; Max-Age=31536000; HttpOnly
```

#### Token 解码内容

webvpn-token 为 JWT 格式，解码后包含：

- `userId`: 用户 ID（如 615）
- `userName`: 学号（如 24270123）
- `authType`: 认证类型（4 = CAS）
- `salt`: 盐值用于后续签名验证
- `exp`: 过期时间戳

#### 获取认证配置

```
GET /api/access/authentication/conf
```

返回会话管理、安全策略等配置：

```json
{
  "code": 0,
  "data": {
    "sessionConf": {
      "maxLifeTime": 86400,
      "sessionExpireAtBrowserClose": false
    },
    "ipAuthLockConf": {
      "open": true,
      "lockDuration": 600,
      "maxLoginAttempts": 10
    },
    "accountAuthLockConf": {
      "open": true,
      "lockDuration": 600,
      "maxLoginAttempts": 6
    }
  }
}
```

### 1.2 用户信息 API

```
GET /api/access/user/info
```

未授权时返回：

```json
{"code": 401, "message": "未授权", "data": null}
```

授权后返回：

```json
{
  "code": 0,
  "data": {
    "userId": 615,
    "username": "",
    "nickname": "24270123",
    "fullName": "",
    "groups": ["sso认证组", "默认"],
    "authType": 4,
    "needTriggerTFA": false,
    "needChangePwd": false
  }
}
```

### 1.3 导航与站点 API

#### 获取站点列表

```
GET /api/access/nav/site-list
```

返回内网资源站点目录，按分类组织：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "list": [
      {
        "name": "帮助中心",
        "sites": [
          {
            "id": 157,
            "name": "常见问题",
            "url": "https://https-vpnhelp-hdu-edu-cn-443.webvpn.hdu.edu.cn/",
            "icon": "AppstoreOutlined",
            "color": "#5398f7",
            "iconUrl": "",
            "sort": 3,
            "rawURL": "https://vpnhelp.hdu.edu.cn/"
          },
          {
            "id": 253,
            "name": "新版内网门户",
            "url": "https://https-web-hdu-edu-cn-443.webvpn.hdu.edu.cn",
            "icon": "AppstoreOutlined",
            "color": "#5398f7",
            "iconUrl": "",
            "sort": 2,
            "rawURL": "https://web.hdu.edu.cn"
          }
        ]
      },

      ......
      
      {
        "name": "教学科研",
        "sites": [
          {
            "id": 153,
            "name": "教务系统",
            "url": "https://http-newjw-hdu-edu-cn-80.webvpn.hdu.edu.cn/sso/driot4login",
            "icon": "AppstoreOutlined",
            "color": "#5398f7",
            "iconUrl": "",
            "sort": 30,
            "rawURL": "http://newjw.hdu.edu.cn/sso/driot4login"
          },
          {
            "id": 155,
            "name": "智慧课堂（校园网内请直接访问https://course.hdu.edu.cn/）",
            "url": "https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn/#/home",
            "icon": "AppstoreOutlined",
            "color": "#5398f7",
            "iconUrl": "",
            "sort": 27,
            "rawURL": "https://course.hdu.edu.cn/#/home"
          }
        ]
      }
    ]
  }
}
```

#### 获取自定义导航配置

```
GET /api/access/nav/custom
```

#### 获取导航配置

```
GET /api/access/nav/config
```

```json
{"code": 0, "data": {"config": {"showUrl": true}}}
```

#### 收藏站点管理

```
POST /api/access/nav/favorite-sites
```

请求需要携带 `webvpn-token` Cookie，Content-Length 为 0 时表示获取收藏列表。

#### Web Terminal 列表

```
GET /api/access/web-terminal/list
```

---

## 二、webvpn-token 获取流程详解

完整的认证流程如下：

### 流程图

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           用户访问 WebVPN                                 │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  1. GET /api/access/authentication/list                                  │
│     获取可用认证方式（当前为统一身份认证）                                   │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  2. POST /api/access/auth/start                                          │
│     发起认证请求，获取 SSO 登录 URL                                       │
│     返回: login_url = "https://sso.hdu.edu.cn/login?service=..."         │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  3. 浏览器重定向到 sso.hdu.edu.cn 进行身份验证                             │
│     GET https://sso.hdu.edu.cn/login?service=callback_url               │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  4. 用户提交用户名/密码到 SSO                                             │
│     POST https://sso.hdu.edu.cn/login (form submission)                 │
│     携带: username, type=UsernamePassword, _eventId=submit, execution  │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  5. SSO 返回重定向到 WebVPN Callback                                      │
│     Location: https://webvpn.hdu.edu.cn/callback/cas/{externalId}       │
│                             ?ticket=ST-xxxxxx                          │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  6. GET /callback/cas/{externalId}?ticket=ST-xxx                        │
│     WebVPN 验证 CAS Ticket                                               │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  7. POST /api/access/auth/finish                                         │
│     提交 ticket 和 deviceId 完成认证                                     │
│     body: {"externalId": "...", "ticket": "...", "deviceId": "..."}     │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│  8. 获取 webvpn-token                                                    │
│     响应 Set-Cookie: webvpn-token=JWT...                                 │
│     响应 body: {"code": 0, "data": {"token": "JWT..."}}                 │
└─────────────────────────────────────────────────────────────────────────┘
```

### deviceId 生成机制

从请求包中可以看到 `deviceId: "8a5541b34915d358d5842ec56d040ef6"`，这是一个 32 位的 MD5 格式字符串，用于设备标识。

### Ticket 有效期

CAS Ticket 格式：`ST-1816171-Nr2J8hMaNdbwZiN-1Xba-c0uTaIrg-sso-64ccf764df-56tht`

---

## 三、浏览器端请求封装机制

### 3.1 URL 封装规则

WebVPN 对内网资源的 URL 进行特殊封装，格式为：

```
https://{protocol}-{hostname}-{port}.webvpn.hdu.edu.cn/{path}
```

**示例**：

| 原始内网 URL | WebVPN 封装后 URL |
|-------------|------------------|
| `https://vpnhelp.hdu.edu.cn/` | `https://https-vpnhelp-hdu-edu-cn-443.webvpn.hdu.edu.cn/` |
| `https://web.hdu.edu.cn` | `https://https-web-hdu-edu-cn-443.webvpn.hdu.edu.cn` |
| `http://chat.hdu.edu.cn/chat/...` | `https://http-chat-hdu-edu-cn-80.webvpn.hdu.edu.cn/chat/...` |
| `http://course.hdu.edu.cn/...` | 通过 CAS 回调处理 |

**URL 封装规则解析**：

- `https://` → `https-`
- `http://` → `http-`
- 域名中的 `.` 替换为 `-`
- 端口号直接拼接在域名最后
- 原始路径保持不变

**重要发现：此 URL 转换不在前端 JavaScript 中实现**

经过对所有反混淆和原始 JavaScript 文件的全面搜索，**未在前端代码中找到**将内网 URL 转换为 webvpn 格式的逻辑。

**结论**：此 URL 转换是在 **服务器端（nginx + 后端）** 完成的：

- `GET /api/access/nav/site-list` API 返回的 `url` 字段已经是转换后的 webvpn 格式
- `rawURL` 字段是原始内网 URL
- 前端仅负责渲染这些 URL，不做格式转换

相关文件：

| 文件 | 功能 |
|------|------|
| `navCustom-CgDbtGc9.js` | 调用 `/api/access/nav/site-list` 获取站点列表 |
| `convertBaseUrl.js` | 仅处理前端内部路由和 base URL，不处理 webvpn URL 格式 |

### 3.2 Cookie 传递

认证成功后，浏览器需要携带 `webvpn-token` Cookie 访问受保护资源：

```
Cookie: webvpn-token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3.3 SSO 内部 API 调用

在 SSO（sso.hdu.edu.cn）内部，还涉及以下 API：

| API Path | 用途 |
|----------|------|
| `/linkid/protected/api/dictconfig/get` | 获取字典配置 |
| `/linkid/api/aggregate/identitycategory/protected/all/get` | 获取身份类别 |
| `/api/protected/wechat/checkEqualUser` | 微信绑定检查 |
| `/api/protected/auth/method/getActive` | 获取可用认证方法 |
| `/api/protected/user/findCaptchaCount/{username}` | 验证码次数检查 |

---

## 四、API 响应格式规范

所有 WebVPN API 遵循统一的响应格式：

```json
{
  "code": 0,           // 0=成功，非0=失败
  "message": "ok",     // 状态描述
  "data": {}           // 响应数据，可为 null
}
```

**常见响应码**：

- `0`: 成功
- `401`: 未授权（无有效 token）
- `200`: SSO 内部接口成功

**CORS 响应头**：

```
Access-Control-Allow-Origin: https://webvpn.hdu.edu.cn
Access-Control-Allow-Credentials: true
Access-Control-Allow-Methods: GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE
```

---

## 五、前端资源文件

### 5.1 反混淆文件映射

反混淆后的 JS 文件位于 `webvpn_js_deobfuscated/` 目录：

| 反混淆文件名 | 原始文件名 | 功能 |
|-------------|-----------|------|
| `request.js` | `request-B-CfM8d9.js` | HTTP 请求封装 |
| `01-authentication.js` | `authentication-DJLiEriA.js` | 认证方式列表 API |
| `02-auth-method.js` | `authMethod-CtTsQ3_a.js` | 本地账号绑定 API |
| `03-url-conf.js` | `urlConf-D5f8gSXt.js` | URL 配置 |
| `04-nav-custom.js` | `navCustom-CgDbtGc9.js` | 导航站点 API + Store |
| `05-user-api.js` | `user-cwdmpEuP.js` | 用户信息 API |
| `06-auth-crs.js` | `auth-Crs8kf_D.js` | 认证核心 API |
| `07-auth-bvec.js` | `auth-Bvec_zE4.js` | 认证流程封装 |
| `08-cas-callback.js` | `CasCallback-q97cqoLl.js` | CAS 回调组件 |

### 5.2 未反混淆的关键文件

| 文件名 | 大小 | 说明 |
|--------|------|------|
| `const-BHerojsC.js` | 219KB | 常量配置（包含认证类型枚举等） |
| `user-V10Yxw5H.js` | - | 用户 Pinia Store |
| `util-DyruD4Ub.js` | - | 工具函数（deviceId 生成等） |
| `index-Cey6Kqla.js` | - | Vue 框架核心 |
| `index-o4p_ZbnI.js` | - | Vue 组件（Spin 等） |

---

## 六、关键发现

1. **JWT Token 认证**：webvpn-token 使用 HS256 签名的 JWT，包含用户 ID、学号、认证类型等信息，存储在 HttpOnly Cookie 中。

2. **CAS 单点登录**：WebVPN 依赖 SSO 的 CAS 协议实现身份认证，externalId `57AAnALn` 是唯一的认证配置标识。externalId 通过两处获取：API 响应（构建 SSO 登录 URL）和 URL 路由参数（CAS 回调验证）。

3. **设备指纹**：使用 MD5 格式的 deviceId 进行设备标识。

4. **统一会话管理**：支持同一用户多设备登录检查（enableUniqueSessionConf）。

5. **URL 双向转换**：webvpn URL 格式转换在**服务器端（nginx + 后端）**完成，不在前端 JavaScript 中。前端通过 `/api/access/nav/site-list` API 接收的 `url` 字段已经是转换后的格式。

6. **IP/账号锁定**：具备登录失败锁定机制（IP 和账号维度）。
