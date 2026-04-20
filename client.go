package hduwebvpn

import (
	"net/http"
	"net/http/cookiejar"
)

// Mode 表示访问模式
type Mode int

const (
	WebVPNMode Mode = iota // webvpn 模式，URL 会被转换为 *.webvpn.hdu.edu.cn 格式
	DirectMode             // direct 模式，直接访问内网（需在校园网内）
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
func WithMode(mode Mode) ClientOption {
	return func(c *Client) {
		c.mode = mode
	}
}

// Client 是 WebVPN 客户端，每个实例独立管理自己的认证状态
type Client struct {
	username  string
	password  string
	mode      Mode
	cookiejar *cookiejar.Jar
	services  map[string]*Service
}

// NewClient 创建新的客户端实例，同时完成 CAS SSO 认证
func NewClient(opts ...ClientOption) (*Client, error) {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		mode:      WebVPNMode,
		services:  make(map[string]*Service),
		cookiejar: jar,
	}

	for _, opt := range opts {
		opt(c)
	}

	// TODO: 实现 CAS SSO 认证流程
	// 1. GET /api/access/authentication/list 获取认证方式
	// 2. POST /api/access/auth/start 获取 SSO 登录 URL
	// 3. 访问 SSO 获取 ticket
	// 4. POST /api/access/auth/finish 完成认证，获取 webvpn-token
	// 5. 将 token 存入 cookiejar

	return c, nil
}

// RegisterService 注册一个内网服务
// name: 服务名称，如 "course"
// baseURL: 内网 base URL，如 "https://course.hdu.edu.cn"
func (c *Client) RegisterService(name string, baseURL string) *Service {
	svc := &Service{
		name:    name,
		baseURL: baseURL,
		client:  c,
	}
	c.services[name] = svc
	return svc
}

// Service 获取已注册的服务
func (c *Client) Service(name string) *Service {
	return c.services[name]
}

// Cookies 获取指定域名的 cookies
func (c *Client) Cookies(domain string) []*http.Cookie {
	// TODO: 从 cookiejar 中提取对应 domain 的 cookies
	return nil
}

// Service 代表一个内网服务
type Service struct {
	name    string
	baseURL string
	client  *Client
}

// NewRequest 创建一个新的请求，完全控制请求头（类似 http.Client）
func (s *Service) NewRequest(method, path string, body []byte) (*Request, error) {
	// TODO: 根据 mode 将 baseURL + path 转换为真实 URL
	// webvpn 模式: https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn
	// direct 模式: https://course.hdu.edu.cn
	realURL := ""

	req := &Request{
		service: s,
		method:  method,
		path:    path,
		realURL: realURL,
		Header:  make(http.Header),
		body:    body,
	}
	return req, nil
}

// Get 便捷方法，发起 GET 请求
func (s *Service) Get(path string) (*Response, error) {
	req, err := s.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(req)
}

// Post 便捷方法，发起 POST 请求
func (s *Service) Post(path string, body []byte) (*Response, error) {
	req, err := s.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	return s.client.Do(req)
}

// Do 发送请求
func (c *Client) Do(req *Request) (*Response, error) {
	// TODO:
	// 1. 根据 mode 将 req.path 转换为真实 URL
	// 2. 构建真实请求（http.Request），设置 Header、Cookie 等
	// 3. 记录 RealReq 信息（URL、Method、Header、Body）
	// 4. 发送请求
	// 5. 自动处理 Set-Cookie 到 cookiejar
	// 6. 检测 401，自动重新认证并重试
	// 7. 返回包含 RealReq 的 Response

	realReq := &RealRequest{
		URL:    req.realURL,
		Method: req.method,
		Header: req.Header,
		Body:   req.body,
	}

	resp := &Response{
		RealReq:    realReq,
		StatusCode: 200,
		Header:     make(http.Header),
		body:       []byte{},
	}

	return resp, nil
}

// Request 是用户构建的请求（业务层概念）
type Request struct {
	service *Service
	method  string
	path    string
	realURL string // 转换后的真实 URL
	Header  http.Header
	body    []byte
}

// RealRequest 是实际发出的真实请求
type RealRequest struct {
	URL    string
	Method string
	Header http.Header
	Body   []byte
}

// Response 是请求响应
type Response struct {
	RealReq    *RealRequest // 实际发出的真实请求
	StatusCode int          // 真实响应状态码
	Header     http.Header // 真实响应头
	body       []byte       // 业务层响应体
}

// Body 返回响应体
func (r *Response) Body() []byte {
	return r.body
}
