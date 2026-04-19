# WebVPN JS 反混淆文件索引

## 文件清单

| 文件 | 对应原始文件 | 功能 |
|------|------------|------|
| `request.js` | request-B-CfM8d9.js | HTTP 请求封装 |
| `00_MAPPING.md` | - | 混淆-反混淆符号对照表 |
| `01-authentication.js` | authentication-DJLiEriA.js | 认证方式列表 API |
| `02-auth-method.js` | authMethod-CtTsQ3_a.js | 本地账号绑定 API |
| `03-url-conf.js` | urlConf-D5f8gSXt.js | URL 配置 |
| `04-nav-custom.js` | navCustom-CgDbtGc9.js | 导航站点 API + Store |
| `05-user-api.js` | user-cwdmpEuP.js | 用户信息 API |
| `06-auth-crs.js` | auth-Crs8kf_D.js | 认证核心 API |
| `07-auth-bvec.js` | auth-Bvec_zE4.js | 认证流程封装 |
| `08-cas-callback.js` | CasCallback-q97cqoLl.js | CAS 回调组件 |

---

## 导入关系图

```
request.js (HTTP 封装)
    ↑
    │ 导入
    │
01-authentication.js ──────────→ request.js
02-auth-method.js ─────────────→ request.js
03-url-conf.js (无外部导入)
04-nav-custom.js ──────────────→ index-Cey6Kqla.js (Vue框架)
         │                              request.js
         │
         └──────────────────────────────→ 05-user-api.js ──────────→ index-Cey6Kqla.js (Vue框架)
         │                                              request.js
         │
         └──────────────────────────────────────────────→ 06-auth-crs.js ────────────→ request.js
         │                                              (auth-core)
         │
         └──────────────────────────────────────────────→ 07-auth-bvec.js ─────────→ 06-auth-crs.js
         │                                              (auth-wrapper)      02-auth-method.js
         │                                                               user-V10Yxw5H.js (未反混淆)
         │
         └──────────────────────────────────────────────────────────────────→ 08-cas-callback.js ──→ 07-auth-bvec.js
                                                                             (CasCallback)      05-user-api.js
                                                                                              user-V10Yxw5H.js (未反混淆)
                                                                                              util-DyruD4Ub.js (未反混淆)
                                                                                              index-o4p_ZbnI.js (未反混淆)
                                                                                              request.js
```

---

## 未反混淆的文件 (保持原名)

以下文件尚未反混淆，保持原始文件名作为导入路径：

| 原始文件 | 说明 |
|---------|------|
| `index-Cey6Kqla.js` | Vue 框架核心 (refs, reactive, computed 等) |
| `user-V10Yxw5H.js` | 用户 Pinia Store |
| `util-DyruD4Ub.js` | 工具函数 (含 deviceId 生成) |
| `index-o4p_ZbnI.js` | Vue 组件 (Spin 等) |
| `const-BHerojsC.js` | 常量配置 |
| `styleChecker-TFxV7Lzl.js` | 样式检查器 |
| `index-DXhXwGTY.js` | 索引模块 |
| `initDefaultProps-_hAYIFu4.js` | 默认属性初始化 |
| `fade-BdgfkLw7.js` | 动画效果 |

---

## API 端点汇总

### 认证相关
| 端点 | 方法 | 文件 |
|------|------|------|
| `/api/access/auth/start` | POST | 06-auth-crs.js |
| `/api/access/auth/finish` | POST | 06-auth-crs.js |
| `/api/access/auth/tfa` | POST | 06-auth-crs.js |
| `/api/access/auth/tfa-config` | POST | 06-auth-crs.js |
| `/api/access/auth/user-notice-info` | GET | 06-auth-crs.js |
| `/api/access/auth/consume-session` | GET | 06-auth-crs.js |
| `/api/access/auth/reset-password` | POST | 06-auth-crs.js |
| `/api/access/auth/wechat-log` | POST | 06-auth-crs.js |
| `/api/access/auth/session-token` | POST | 06-auth-crs.js |

### 用户相关
| 端点 | 方法 | 文件 |
|------|------|------|
| `/api/access/user/info` | GET | 05-user-api.js |
| `/api/access/user/logout` | POST | 05-user-api.js |
| `/api/access/user/change-password` | POST | 05-user-api.js |
| `/api/access/user/change-info` | POST | 05-user-api.js |
| `/api/access/user/change-email` | POST | 05-user-api.js |
| `/api/access/user/change-mobile` | POST | 05-user-api.js |
| `/api/access/user/change-avatar` | POST | 05-user-api.js |
| `/api/access/user/auth/history` | GET | 05-user-api.js |
| `/api/access/access-log/list` | GET | 05-user-api.js |

### 导航相关
| 端点 | 方法 | 文件 |
|------|------|------|
| `/api/access/nav/site-list` | GET | 04-nav-custom.js |
| `/api/access/nav/favorite-sites` | POST | 04-nav-custom.js |
| `/api/access/nav/add-to-favorites` | POST | 04-nav-custom.js |
| `/api/access/nav/remove-from-favorites` | POST | 04-nav-custom.js |
| `/api/access/nav/config` | GET | 04-nav-custom.js |
| `/api/access/nav/custom` | GET | 04-nav-custom.js |

### 认证方式
| 端点 | 方法 | 文件 |
|------|------|------|
| `/api/access/authentication/list` | GET | 01-authentication.js |
| `/api/access/authentication/all` | GET | 01-authentication.js |

### 本地账号绑定
| 端点 | 方法 | 文件 |
|------|------|------|
| `/api/access/auth-method-bind-local/start` | POST | 02-auth-method.js |
| `/api/access/auth-method/bind-local-account` | POST | 02-auth-method.js |
| `/api/access/half-account-bind/start` | POST | 02-auth-method.js |
| `/api/access/half-account-bind/finish` | POST | 02-auth-method.js |

---

## 关键发现

### 1. externalId 的获取方式

**源码证据** (`08-cas-callback.js`):
```javascript
const externalId = route.currentRoute.value.params.externalId;
```

`externalId` 从 **URL 路由参数** 获取，格式为 `/callback/cas/:externalId`

### 2. logout 不是自动流程的一部分

`logout()` 是用户**主动触发**的操作，不是认证流程的自动步骤。

### 3. 完整的 CAS 回调流程

```javascript
// 08-cas-callback.js 中的流程:

onMounted(async () => {
    // 1. 检查是否已登录
    const userResult = await store.getUserInfo({ unbind: true });

    // 2. 如果未登录且有 ticket
    if (!(userResult?.data?.userId) && route.query.ticket) {

        // 3. 获取 ticket 和 externalId
        const ticket = decodeURIComponent(route.query.ticket);
        const externalId = route.currentRoute.value.params.externalId;

        // 4. 构建认证请求
        const authPayload = {
            externalId: externalId,
            data: JSON.stringify({
                callbackUrl: window.location.href.split("?")[0],
                ticket: ticket,
                deviceId: await generateDeviceId()
            })
        };

        // 5. 调用认证 API
        await authenticate(authPayload);
    }
});
```
