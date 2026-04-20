# hduwebvpn

2026.4.15 左右，杭电的 webvpn 的 API 进行了很大的更新，好消息是 webvpn 的各种接口都变得规整了许多，坏消息是我得重新扒一遍。

## Feature

本库希望提供一个通用的 hduclient，可用于封装各类受 webvpn 支持的 hdu 应用服务，可以配置使用 inner net/webvpn 访问模式，而无需重写 webvpn 和 sso 的认证登录和失效重新登录逻辑。

服务注册时需要提供所使用的内网 base url（例如 course.hdu.edu.cn），实现接口时，传入 api 路径而不是完整 url。可以从每一次的请求的返回中得知请求实际上发送的 RealURL（例如，在 webvpn 模式下会经过 webvpn 转换）。

client 对于 Set-Cookies 响应头可以自动管理 cookies。同样支持手动从请求中提取 cookies，并存入 cookiesjar。

## 核心架构设计：双重 URL 驱动

本库的核心设计思想是将**业务层地址 (Business URL)** 和**实际请求地址 (Real URL)** 彻底解耦。在这种架构下：

- `cookiejar` 和底层的 `http.Client` 完全是"瞎子"，它们只认 Real URL 并处理实际的网络通信
- Service 层的逻辑永远只看 Business URL
- `Transport` 充当了这两界之间的翻译官

### 禁用自动重定向

默认的 `http.Client` 会自动跟随 302。本库在 `CheckRedirect` 中返回 `http.ErrUseLastResponse`，将 3xx 响应拦截下来，交由 Middleware 处理。

### Transport 抽象

将 Direct 和 WebVPN 抽象为 Transport 接口，负责提供 `EncodeURL` (Business -> Real) 和 `DecodeURL` (Real -> Business) 的能力：

```go
type Transport interface {
    // Encode 将业务地址转换为实际请求地址
    Encode(businessURL *url.URL) *url.URL

    // Decode 将实际的重定向地址还原为业务地址
    // 如果该 RealURL 不属于本 transport 的常规包装结构（例如 WebVPN 登录页），返回 nil
    Decode(realURL *url.URL) *url.URL

    // IsAuthFailure 判断该实际重定向地址是否意味着当前 Transport 的通道认证（如 WebVPN 认证）失效
    IsAuthFailure(realURL *url.URL) bool

    // Reauth 重新执行通道层认证
    Reauth(client *Client) error
}
```

**Direct Transport:**
- `Encode`: 原样返回
- `Decode`: 原样返回
- `IsAuthFailure`: 永远返回 false（内网直连不存在通道层掉线）

**WebVPN Transport:**
- `Encode`: 加上 `https-xxx-443.webvpn.hdu.edu.cn` 壳子
- `Decode`: 剥离 WebVPN 壳子，还原为原始业务地址
- `IsAuthFailure`: 如果 Location 的 Host 不是 `*.webvpn.hdu.edu.cn` 且不是自身域名，说明掉出了 WebVPN 环境

## Gin 风格中间件链

本库采用 Gin 风格的中间件链设计，核心类型为 `Context` 和 `Handler`：

```go
type Handler func(*Context)

type Context struct {
    Request  *Request
    Response *Response
    index    int
    handlers []Handler
    client   *Client
    Err      error
}
```

通过 `Context.Next()` 递归调用链，实现"洋葱模型"。每个 Handler 可以：
- 在 `c.Next()` 之前做前置处理（如 URL 翻译）
- 在 `c.Next()` 之后做后置处理（如重定向检查）
- 使用 `c.Abort()` 阻止后续 handler 执行

## 请求流转图

```
用户调用 (User Request)
     │  携带 BusinessURL (例如: https://course.hdu.edu.cn/api)
     ▼
┌─────────────────────────────────────────────────────────────┐
│ serviceAuthHandler (业务层视角)                               │
│                                                              │
│  - 调用 c.Next() 透传给下一层                                │
│                                                              │
│  ▲ (c.Next() 返回后)                                         │
│  │ 检查 resp.BusinessRedirectURL:                            │
│  │ 如果重定向到了 sso.hdu.edu.cn -> 触发 Service SSO 认证      │
│  │ 认证完成后 -> 重新调用本 handler 再次发起请求               │
│  │ 否则 -> 将 resp 向上返回给用户                            │
└─────────────────────────────────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────────────────────────────────┐
│ transportHandler (翻译层 & 通道层视角)                        │
│                                                              │
│  - 根据当前 Mode (WebVPN/Direct) 将 BusinessURL 翻译为 RealURL│
│  - 将 RealURL 赋给 Request                                   │
│  - 调用 c.Next() 透传给下一层                                │
│                                                              │
│  ▲ (c.Next() 返回后)                                         │
│  │ 检查 resp.StatusCode (是否 3xx):                          │
│  │ 1. 获取 Header 中的 Location (即 RealRedirectURL)         │
│  │ 2. 判断是否是 WebVPN 自身的掉线重定向 (非 .webvpn 域名):   │
│  │    -> 触发 WebVPN SSO 认证                                 │
│  │    -> 认证完成后 -> 重新调用本 handler 再次发起请求        │
│  │ 3. 如果是正常的通道内重定向 (例如 webvpn 包装的 SSO 链接): │
│  │    -> Decode 还原为 BusinessRedirectURL 赋给 Response      │
│  │    -> 将 resp 向上返回给 serviceAuthHandler               │
└─────────────────────────────────────────────────────────────┘
     │
     ▼
┌─────────────────────────────────────────────────────────────┐
│ execBaseDo (底层)                                           │
│  - 从 CookieJar 中获取 RealURL 的 Cookie                     │
│  - 发起真实的 HTTP 请求 (禁止自动重定向)                     │
│  - 自动将 Set-Cookie 写入 CookieJar (关联 RealURL)           │
│  - 返回原始 Response                                         │
└─────────────────────────────────────────────────────────────┘
```

## 核心类型

### Request

`Request` 是业务层构建的请求：

```go
type Request struct {
    service     *Service
    method      string
    BusinessURL *url.URL     // 业务地址，例如 https://course.hdu.edu.cn/api
    RealURL     *url.URL     // 实际地址，由 Transport 转换填入
    Header      http.Header
    body        []byte
}
```

### Response

`Response` 包含解耦后的响应信息：

```go
type Response struct {
    RealReq             *RealRequest
    RawResponse         *http.Response // 原始的 http.Response
    StatusCode          int
    Header              http.Header
    body                []byte

    // 核心新增：用于分离重定向逻辑
    RealRedirectURL     *url.URL // 从 Location Header 直接解析出的实际重定向地址
    BusinessRedirectURL *url.URL // 经过 Transport 解码还原后的业务层重定向地址
}
```

### RealRequest

`RealRequest` 是实际发出的真实请求：

```go
type RealRequest struct {
    URL     string      // 真实请求 URL（webvpn 转换后）
    Method  string      // GET/POST 等
    Header  http.Header // 实际发出的请求头
    Body    []byte      // 实际发出的请求体
}
```

## 使用示例

### 1. 创建客户端（包含认证信息）

```golang
package main

import (
    "fmt"
    hduwebvpn "hduwebvpn"
)

func main() {
    // 创建客户端，同时传入认证信息，每个 client 独立管理自己的登录状态
    client := hduwebvpn.NewClient(
        hduwebvpn.WithUsername("24270123"),
        hduwebvpn.WithPassword("password"),
    )

    // 或指定模式（webvpn / direct）
    client := hduwebvpn.NewClient(
        hduwebvpn.WithUsername("24270123"),
        hduwebvpn.WithPassword("password"),
        hduwebvpn.WithMode(hduwebvpn.DirectMode),  // 校内直连
    )
}
```

客户端创建时会自动完成 CAS SSO 认证流程（包括获取 ticket、提交验证、获取 token）。

### 2. 注册内网服务并发起请求

```golang
// 注册服务：智慧课堂 (course.hdu.edu.cn)
client.RegisterService("course", "https://course.hdu.edu.cn")

// 方式一：使用便捷方法（自动设置必要的 Header）
userInfo, err := client.Service("course").Get("/api/access/user/info")
// userInfo: {"code": 0, "data": {"userId": 615, "nickname": "24270123", ...}}

// 方式二：构建 Request 对象，完全控制请求头（类似 http.Client）
req, _ := client.Service("course").NewRequest("GET", "/api/access/user/info", nil)
req.Header.Set("Accept", "application/json")
req.Header.Set("X-Custom-Header", "custom-value")
resp, err := client.Service("course").Do(req)
```

API 设计尽量接近 `net/http` 的风格，通过 `Request.Header` 管理请求头。

### 3. 访问其他内网服务

```golang
// 注册教务系统
client.RegisterService("jw", "http://newjw.hdu.edu.cn")

// 访问教务系统 API
score, err := client.Service("jw").Get("/sso/driot4login")
```

### 4. 真实请求与响应

`Response` 包含完整真实请求（`RealReq`）和真实响应信息，透明化 webvpn 的 URL 转换：

```golang
resp, err := client.Service("course").Do(req)
if err != nil {
    panic(err)
}

// 真实请求 URL（webvpn 转换后）
fmt.Println(resp.RealReq.URL)
// https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn/api/access/user/info

// 实际发出的请求头（含 webvpn-token Cookie）
fmt.Println(resp.RealReq.Header.Get("Cookie"))

// 业务响应
fmt.Println(resp.StatusCode) // 200
fmt.Println(string(resp.Body)) // {"code": 0, ...}
```

### 5. Cookie 管理

```golang
// 客户端自动管理 Set-Cookie
// 如需手动提取 cookie：
cookies := client.Cookies("course.hdu.edu.cn")
for _, c := range cookies {
    fmt.Printf("Name: %s, Value: %s\n", c.Name, c.Value)
}
```

### 5.1 自定义请求头（OAuth2 等场景）

某些内网服务有独立的认证体系，需要从响应中提取 token 并在后续请求中携带：

```golang
client.RegisterService("course", "https://course.hdu.edu.cn")

// 1. 构建请求，自由设置任意 Header
req, _ := client.Service("course").NewRequest("GET", "/jy-application-patrol-class/oauth2/token", nil)
req.Header.Set("tenantid", "RBAC")
req.Header.Set("Accept", "application/json, text/plain, */*")

// 2. 发送请求
resp, err := client.Service("course").Do(req)

// 3. 用户自行解析响应 JSON，提取 jwt_token
var result struct {
    Result struct{ JwtToken string } `json:"result"`
}
json.Unmarshal(resp.Body(), &result)
jwtToken := result.Result.JwtToken

// 4. 后续请求携带该 token
req2, _ := client.Service("course").NewRequest("GET", "/jy-application-patrol-class/authority/me", nil)
req2.Header.Set("jwt-token", jwtToken)
resp2, _ := client.Service("course").Do(req2)
```

### 6. Token 自动刷新

如果 token 过期，客户端会自动使用保存的凭据重新认证并重试请求：

```golang
// 当 token 失效时，Do 方法会：
// 1. 检测到 401 未授权响应
// 2. 使用保存的 username/password 重新获取 token
// 3. 重试失败的请求
resp, err := client.Service("course").Get("/api/some-resource")
```

### 7. 完整示例：访问智慧课堂

```golang
package main

import (
    "fmt"
    "log"
    hduwebvpn "hduwebvpn"
)

func main() {
    // 创建客户端时即完成认证，每个 client 独立
    client, err := hduwebvpn.NewClient(
        hduwebvpn.WithUsername("24270123"),
        hduwebvpn.WithPassword("password"),
    )
    if err != nil {
        log.Fatalf("创建客户端失败: %v", err)
    }

    // 注册智慧课堂服务
    client.RegisterService("course", "https://course.hdu.edu.cn")

    // 获取用户信息
    userInfo, err := client.Service("course").Get("/api/access/user/info")
    if err != nil {
        log.Fatalf("请求失败: %v", err)
    }
    fmt.Printf("用户信息: %s\n", userInfo.Body())

    // 获取课程列表
    courses, err := client.Service("course").Get("/api/course/list")
    fmt.Printf("课程列表: %s\n", courses.Body())
}
```

## 设计说明

- **Client**：每个 Client 实例独立管理自己的认证信息（username/password）和会话状态，绑定时即完成登录
- **Service**：代表一个内网服务（如 course.hdu.edu.cn），注册时指定内网 base URL
- **Transport**：负责 Business URL 与 Real URL 之间的相互转换，是 WebVPN/Direct 模式的抽象接口
- **Mode**：`webvpn` 模式将 URL 转换为 `*.webvpn.hdu.edu.cn` 格式；`direct` 模式直接访问内网（需在校园网内）
- **双重 Middleware**：Transport Middleware 处理通道层认证（如 WebVPN 掉线重连），Service Auth Middleware 处理业务层认证（如 CAS SSO）
- **自动重试**：无论是通道层还是业务层认证失效，都会自动重新认证并重试请求
- **CookieJar 行为自洽**：由于底层始终使用 Real URL 操作 cookiejar，业务系统的 Cookie 和 WebVPN Token 能正确分离管理
