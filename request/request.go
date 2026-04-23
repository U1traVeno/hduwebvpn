package request

import (
	"net/http"
	"net/url"
)

// Request 是用户构建的请求（业务层概念）
type Request struct {
	Service     interface{} // *service.Service - 接口类型避免循环依赖
	Method      string
	BusinessURL *url.URL // 业务地址，例如 https://course.hdu.edu.cn/api
	RealURL     *url.URL // 实际地址，由 Transport 转换填入
	Header      http.Header
	Body        []byte
}

// Response 是请求响应
type Response struct {
	BusinessReqURL *url.URL       // 实际请求的业务地址。注意，经过重定向后可能与 Request.BusinessURL 不同。
	RawResponse    *http.Response // 原始的 http.Response
	StatusCode     int
	Header         http.Header
	Body           []byte
}
