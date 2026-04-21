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
	RawResponse *http.Response // 原始的 http.Response
	StatusCode  int
	Header      http.Header
	body        []byte

	// 用于分离重定向逻辑
	RealRedirectURL     *url.URL // 从 Location Header 直接解析出的实际重定向地址
	BusinessRedirectURL *url.URL // 经过 Transport 解码还原后的业务层重定向地址
}
