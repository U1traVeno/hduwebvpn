package transport

import (
	"net/url"
	"strings"
)

// Mode 表示访问模式
type Mode int

const (
	WebVPNMode Mode = iota // webvpn 模式，URL 会被转换为 *.webvpn.hdu.edu.cn 格式
	DirectMode             // direct 模式，直接访问内网（需在校园网内）
)

// Transport 接口定义了 URL 编码/解码和认证重试能力
type Transport interface {
	// Encode 将业务地址转换为实际请求地址
	Encode(businessURL *url.URL) *url.URL

	// Decode 将实际的重定向地址还原为业务地址
	// 如果该 RealURL 不属于本 transport 的常规包装结构（例如 WebVPN 登录页），返回 nil
	Decode(realURL *url.URL) *url.URL

	// IsAuthFailure 判断该实际重定向地址是否意味着通道认证失效
	IsAuthFailure(realURL *url.URL) bool

	// Reauth 重新执行通道层认证
	Reauth(client interface{}) error
}

// DirectTransport 直连模式 Transport
type DirectTransport struct{}

// Client interface for Reauth
type Client interface {
	GetHTTPClient() interface{}
	GetCookieJar() interface{}
}

func (t *DirectTransport) Encode(businessURL *url.URL) *url.URL {
	return businessURL
}

func (t *DirectTransport) Decode(realURL *url.URL) *url.URL {
	return realURL
}

func (t *DirectTransport) IsAuthFailure(realURL *url.URL) bool {
	return false // 直连模式不存在通道层掉线
}

func (t *DirectTransport) Reauth(client interface{}) error {
	return nil // 直连模式无需重认证
}

// WebVPNTransport WebVPN 模式 Transport
type WebVPNTransport struct{}

func (t *WebVPNTransport) Encode(businessURL *url.URL) *url.URL {
	// 格式: https://https-{host}-{port}.webvpn.hdu.edu.cn
	// 例如: https://course.hdu.edu.cn -> https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn
	return nil
}

func (t *WebVPNTransport) Decode(realURL *url.URL) *url.URL {
	// 逆向解析: https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn -> https://course.hdu.edu.cn
	return nil
}

func (t *WebVPNTransport) IsAuthFailure(realURL *url.URL) bool {
	// 如果重定向到了非 webvpn 域名，说明掉出了 WebVPN 环境
	if realURL.Host == "" {
		return false
	}
	return !strings.HasSuffix(realURL.Host, ".webvpn.hdu.edu.cn")
}

func (t *WebVPNTransport) Reauth(client interface{}) error {
	// TODO: 实现 WebVPN 通道层重认证
	return nil
}