package hduwebvpn

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
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
		c.transport = c.newTransport(mode)
	}
}

// =============================================================================
// Transport 接口 - 负责 Business URL 与 Real URL 之间的相互转换
// =============================================================================

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
	Reauth(client *Client) error
}

// DirectTransport 直连模式 Transport
type DirectTransport struct{}

func (t *DirectTransport) Encode(businessURL *url.URL) *url.URL {
	return businessURL
}

func (t *DirectTransport) Decode(realURL *url.URL) *url.URL {
	return realURL
}

func (t *DirectTransport) IsAuthFailure(realURL *url.URL) bool {
	return false // 直连模式不存在通道层掉线
}

func (t *DirectTransport) Reauth(client *Client) error {
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

func (t *WebVPNTransport) Reauth(client *Client) error {
	// TODO: 实现 WebVPN 通道层重认证
	return nil
}

// =============================================================================
// Client
// =============================================================================

// Client 是 WebVPN 客户端，每个实例独立管理自己的认证状态
type Client struct {
	username  string
	password  string
	mode      Mode
	transport Transport
	cookiejar *cookiejar.Jar
	httpClient *http.Client
	services  map[string]*Service
}

// newTransport 根据 mode 创建对应的 Transport
func (c *Client) newTransport(mode Mode) Transport {
	switch mode {
	case DirectMode:
		return &DirectTransport{}
	case WebVPNMode:
		return &WebVPNTransport{}
	default:
		return &WebVPNTransport{}
	}
}

// NewClient 创建新的客户端实例，同时完成 CAS SSO 认证
func NewClient(opts ...ClientOption) (*Client, error) {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		mode:      WebVPNMode,
		services:  make(map[string]*Service),
		cookiejar: jar,
		transport: &WebVPNTransport{},
	}

	for _, opt := range opts {
		opt(c)
	}

	// 配置 http.Client，禁用自动重定向
	c.httpClient = &http.Client{
		Jar: c.cookiejar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 拦截 3xx，由 Middleware 处理
		},
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

// =============================================================================
// Context & Handler
// =============================================================================

type Handler func(*Context)

type Context struct {
	Request  *Request
	Response *Response
	index    int
	handlers []Handler
	client   *Client
	Err      error
}

func NewContext(client *Client, req *Request, handlers []Handler) *Context {
	return &Context{
		Request:  req,
		Response: &Response{},
		index:    -1,
		handlers: handlers,
		client:   client,
	}
}

func (c *Context) Next() {
	if c.Err != nil {
		return
	}
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	} else {
		c.execBaseDo()
	}
}

func (c *Context) Abort() {
	c.index = len(c.handlers)
}

func (c *Context) execBaseDo() {
	httpReq, err := http.NewRequest(c.Request.method, c.Request.RealURL.String(), nil)
	if err != nil {
		c.Err = err
		return
	}

	for k, v := range c.Request.Header {
		httpReq.Header[k] = v
	}

	httpResp, err := c.client.httpClient.Do(httpReq)
	if err != nil {
		c.Err = err
		return
	}

	c.Response = &Response{
		RealReq: &RealRequest{
			URL:     c.Request.RealURL.String(),
			Method:  c.Request.method,
			Header:  httpReq.Header,
			Body:    c.Request.body,
		},
		RawResponse: httpResp,
		StatusCode:  httpResp.StatusCode,
		Header:      httpResp.Header,
	}
}

// =============================================================================
// Middleware 链
// =============================================================================

// isRedirect 判断状态码是否为重定向
func isRedirect(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

// transportHandler 处理通道层认证（WebVPN 掉线重连）
func transportHandler(c *Context) {
	c.Request.RealURL = c.client.transport.Encode(c.Request.BusinessURL)

	c.Next()

	if c.Err != nil {
		return
	}

	resp := c.Response
	if isRedirect(resp.StatusCode) {
		realLoc, err := url.Parse(resp.Header.Get("Location"))
		if err == nil {
			resp.RealRedirectURL = realLoc

			if c.client.transport.IsAuthFailure(realLoc) {
				if authErr := c.client.transport.Reauth(c.client); authErr != nil {
					c.Err = fmt.Errorf("transport reauth failed: %w", authErr)
					return
				}
				transportHandler(c)
				return
			}

			resp.BusinessRedirectURL = c.client.transport.Decode(realLoc)
		}
	}
}

// serviceAuthHandler 处理业务层认证（CAS SSO）
func serviceAuthHandler(c *Context) {
	c.Next()

	if c.Err != nil {
		return
	}

	resp := c.Response
	if isRedirect(resp.StatusCode) && resp.BusinessRedirectURL != nil {
		if strings.Contains(resp.BusinessRedirectURL.Host, "sso.hdu.edu.cn") {
			if authErr := c.Request.service.doSSO(); authErr != nil {
				c.Err = fmt.Errorf("service sso auth failed: %w", authErr)
				return
			}
			serviceAuthHandler(c)
			return
		}
	}
}

// =============================================================================
// Service
// =============================================================================

// Service 代表一个内网服务
type Service struct {
	name    string
	baseURL string
	client  *Client
}

// NewRequest 创建一个新的请求，完全控制请求头（类似 http.Client）
func (s *Service) NewRequest(method, path string, body []byte) (*Request, error) {
	base, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	// 拼接 business URL
	businessURL := base.ResolveReference(&url.URL{Path: path})

	req := &Request{
		service:     s,
		method:      method,
		BusinessURL: businessURL,
		Header:      make(http.Header),
		body:        body,
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

// doSSO 执行业务层 SSO 认证
func (s *Service) doSSO() error {
	// TODO: 实现 CAS SSO 认证流程
	return nil
}

// =============================================================================
// Do - 请求发送（连接 Middleware 链）
// =============================================================================

// Do 发送请求，构建 Context 并启动中间件链
func (c *Client) Do(req *Request) (*Response, error) {
	handlers := []Handler{
		serviceAuthHandler,
		transportHandler,
	}
	ctx := NewContext(c, req, handlers)
	ctx.Next()
	return ctx.Response, ctx.Err
}

// =============================================================================
// Request / Response / RealRequest
// =============================================================================

// Request 是用户构建的请求（业务层概念）
type Request struct {
	service     *Service
	method      string
	BusinessURL *url.URL // 业务地址，例如 https://course.hdu.edu.cn/api
	RealURL     *url.URL // 实际地址，由 Transport 转换填入
	Header      http.Header
	body        []byte
}

// RealRequest 是实际发出的真实请求
type RealRequest struct {
	URL    string      // 真实请求 URL（webvpn 转换后）
	Method string      // GET/POST 等
	Header http.Header // 实际发出的请求头
	Body   []byte      // 实际发出的请求体
}

// Response 是请求响应
type Response struct {
	RealReq             *RealRequest
	RawResponse         *http.Response // 原始的 http.Response
	StatusCode          int
	Header              http.Header
	body                []byte

	// 用于分离重定向逻辑
	RealRedirectURL     *url.URL // 从 Location Header 直接解析出的实际重定向地址
	BusinessRedirectURL *url.URL // 经过 Transport 解码还原后的业务层重定向地址
}

// Body 返回响应体
func (r *Response) Body() []byte {
	return r.body
}
