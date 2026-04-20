# hduwebvpn

2026.4.15 左右，杭电的 webvpn 的 API 进行了很大的更新，好消息是 webvpn 的各种接口都变得规整了许多，坏消息是我得重新扒一遍。

## Feature

本库希望提供一个通用的 hduclient，可用于封装各类受 webvpn 支持的 hdu 应用服务，可以配置使用 inner net/webvpn 访问模式，而无需重写 webvpn 和 sso 的认证登录和失效重新登录逻辑。

服务注册时需要提供所使用的内网 base url （例如 course.hdu.edu.cn），实现接口时，传入 api 路径而不是完整 url 。可以从每一次的请求的返回中得知请求实际上发送的 RealURL (例如，在 webvpn 模式下会经过 webvpn 转换)。

client 对于 Set-Cookies 响应头可以自动管理 cookies。同样支持手动从请求中提取 cookies，并存入 cookiesjar。

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

### 3. 注册内网服务并发起请求

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

### 4. 访问其他内网服务

```golang
// 注册教务系统
client.RegisterService("jw", "http://newjw.hdu.edu.cn")

// 访问教务系统 API
score, err := client.Service("jw").Get("/sso/driot4login")
```

### 5. 真实请求与响应

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

### 6. Cookie 管理

```golang
// 客户端自动管理 Set-Cookie
// 如需手动提取 cookie：
cookies := client.Cookies("course.hdu.edu.cn")
for _, c := range cookies {
    fmt.Printf("Name: %s, Value: %s\n", c.Name, c.Value)
}
```

### 6.1 自定义请求头（OAuth2 等场景）

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

### 7. Token 自动刷新

如果 token 过期，客户端会自动使用保存的凭据重新认证并重试请求：

```golang
// 当 token 失效时，Do 方法会：
// 1. 检测到 401 未授权响应
// 2. 使用保存的 username/password 重新获取 token
// 3. 重试失败的请求
resp, err := client.Service("course").Get("/api/some-resource")
```

### 8. 完整示例：访问智慧课堂

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
- **Mode**：`webvpn` 模式将 URL 转换为 `*.webvpn.hdu.edu.cn` 格式；`direct` 模式直接访问内网（需在校园网内）
- **自动重试**：token 失效时自动使用保存的凭据重新认证并重试请求

## 核心类型

### Response

`Do` 方法返回的 `Response` 包含真实请求与响应的完整信息：

```golang
type Response struct {
    RealReq    *RealRequest  // 实际发出的真实请求
    StatusCode int           // 真实响应状态码
    Header     http.Header   // 真实响应头
    Body       []byte        // 业务层响应体
}

type RealRequest struct {
    URL    string      // 真实请求 URL（webvpn 转换后）
    Method string      // GET/POST 等
    Header http.Header // 实际发出的请求头
    Body   []byte      // 实际发出的请求体
}
```