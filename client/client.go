package client

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/U1traVeno/hduwebvpn/handler"
	"github.com/U1traVeno/hduwebvpn/request"
	"github.com/U1traVeno/hduwebvpn/service"
	"github.com/U1traVeno/hduwebvpn/transport"
)

// ClientOption 是 Client 的配置选项
type ClientOption func(*Client)

// WithUsername 设置用户名
func WithUsername(username string) ClientOption {
	return func(c *Client) {
		c.username = username
	}
}

// WithPassword 设置密码
func WithPassword(password string) ClientOption {
	return func(c *Client) {
		c.password = password
	}
}

// WithMode 设置访问模式
func WithMode(mode transport.Mode) ClientOption {
	return func(c *Client) {
		c.mode = mode
		c.transport = c.newTransport(mode)
	}
}

// Client is the main VPN client, managing auth state and registered services
type Client struct {
	username   string
	password   string
	mode       transport.Mode
	transport  transport.Transport
	cookiejar  *cookiejar.Jar
	httpClient *http.Client
	services   map[string]interface{} // stores concrete *Service to avoid import cycle
}

// newTransport 根据 mode 创建对应的 Transport
func (c *Client) newTransport(mode transport.Mode) transport.Transport {
	switch mode {
	case transport.DirectMode:
		return &transport.DirectTransport{}
	case transport.WebVPNMode:
		return &transport.WebVPNTransport{}
	default:
		return &transport.WebVPNTransport{}
	}
}

// NewClient 创建新的客户端实例，同时完成 CAS SSO 认证
func NewClient(opts ...ClientOption) (*Client, error) {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		mode:      transport.WebVPNMode,
		services:  make(map[string]interface{}),
		cookiejar: jar,
		transport: &transport.WebVPNTransport{},
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 5 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// RegisterService 注册一个内网服务
func (c *Client) RegisterService(name string, baseURL string) interface{} {
	svc := service.NewService(name, baseURL)
	svc.SetClient(c)
	c.services[name] = svc
	return svc
}

// Service 获取已注册的服务
func (c *Client) Service(name string) interface{} {
	return c.services[name]
}

// Cookies 获取指定域名的 cookies
func (c *Client) Cookies(domain string) []*http.Cookie {
	// TODO: 从 cookiejar 中提取对应 domain 的 cookies
	return nil
}

// Do 发送请求，构建 Context 并启动中间件链
func (c *Client) Do(req *request.Request) (*request.Response, error) {
	return handler.Do(c, req)
}

// GetCookieJar returns the cookie jar (used by Transport interface)
func (c *Client) GetCookieJar() *cookiejar.Jar {
	return c.cookiejar
}

// GetHTTPClient returns the HTTP client (used by Transport interface)
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// GetTransport returns the transport (used for URL encoding/decoding)
func (c *Client) GetTransport() transport.Transport {
	return c.transport
}

// GetUsername returns the username
func (c *Client) GetUsername() string {
	return c.username
}

// GetPassword returns the password
func (c *Client) GetPassword() string {
	return c.password
}

// GetBaseURL parses and returns the base URL for a service
func GetBaseURL(baseURL string) (*url.URL, error) {
	return url.Parse(baseURL)
}
