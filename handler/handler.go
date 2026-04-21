package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/U1traVeno/hduwebvpn/pkg/crypto"
	"github.com/U1traVeno/hduwebvpn/request"
	"github.com/U1traVeno/hduwebvpn/transport"
)

const (
	// SSOHost 常量（用于域名比较）
	SSOHost = "sso.hdu.edu.cn"
	CASHost = "cas.hdu.edu.cn"
)

var (
	// ErrLoginFailed 表示 SSO 登录失败。
	ErrLoginFailed = errors.New("sso login failed")
	// ErrGetFlowkey 表示获取 flowkey/crypto 失败。
	ErrGetFlowkey = errors.New("failed to get flowkey/crypto")
)

// ClientInterface defines what the middleware needs from the client
type ClientInterface interface {
	GetHTTPClient() *http.Client
	GetCookieJar() *cookiejar.Jar
	GetTransport() transport.Transport
	GetUsername() string
	GetPassword() string
}

// Handler is the middleware function type
type Handler func(*Context)

// Context carries state through the middleware chain
type Context struct {
	Request  *request.Request
	Response *request.Response
	index    int
	handlers []Handler
	client   ClientInterface
	Err      error
	logger   *slog.Logger
}

// NewContext constructs a Context with handlers
func NewContext(client ClientInterface, req *request.Request, handlers []Handler) *Context {
	return &Context{
		Request:  req,
		Response: &request.Response{},
		index:    -1,
		handlers: handlers,
		client:   client,
		logger:   slog.Default(),
	}
}

// Next advances to the next handler in the chain
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

// Abort halts the handler chain early
func (c *Context) Abort() {
	c.index = len(c.handlers)
}

// execBaseDo executes the actual HTTP request
func (c *Context) execBaseDo() {
	httpReq, err := http.NewRequest(c.Request.Method, c.Request.RealURL.String(), nil)
	if err != nil {
		c.Err = err
		return
	}

	for k, v := range c.Request.Header {
		httpReq.Header[k] = v
	}

	httpResp, err := c.client.GetHTTPClient().Do(httpReq)
	if err != nil {
		c.Err = err
		return
	}

	c.Response = &request.Response{
		RawResponse: httpResp,
		StatusCode:  httpResp.StatusCode,
		Header:      httpResp.Header,
	}
}

// isRedirect 判断状态码是否为重定向
func isRedirect(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

// isSSOHost 判断 host 是否为 SSO 相关域名
func isSSOHost(host string) bool {
	return host == SSOHost || host == CASHost ||
		strings.HasPrefix(host, "sso-") || strings.HasPrefix(host, "cas-") ||
		strings.HasPrefix(host, "sso.") || strings.HasPrefix(host, "cas.")
}

// transportHandler 处理通道层认证（WebVPN 掉线重连）
func transportHandler(c *Context) {
	t := c.client.GetTransport()
	c.Request.RealURL = t.Encode(c.Request.BusinessURL)

	c.Next()

	if c.Err != nil {
		return
	}

	resp := c.Response
	if isRedirect(resp.StatusCode) {
		realLoc, err := url.Parse(resp.Header.Get("Location"))
		if err == nil {
			resp.RealRedirectURL = realLoc

			if t.IsAuthFailure(realLoc) {
				if authErr := t.Reauth(c.client); authErr != nil {
					c.Err = fmt.Errorf("transport reauth failed: %w", authErr)
					return
				}
				transportHandler(c)
				return
			}

			resp.BusinessRedirectURL = t.Decode(realLoc)
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
		if isSSOHost(resp.BusinessRedirectURL.Host) {
			username := c.client.GetUsername()
			password := c.client.GetPassword()

			if err := c.doServiceSSO(username, password); err != nil {
				c.Err = err
				return
			}
			// 重新执行请求
			serviceAuthHandler(c)
			return
		}
	}
}

// doServiceSSO 执行业务层 SSO 认证流程
func (c *Context) doServiceSSO(username, password string) error {
	ctx := context.Background()
	t := c.client.GetTransport()

	// WebVPN 模式下，将 SSO URL 转换为 WebVPN 格式
	LoginURL := t.Encode(&url.URL{Scheme: "https", Host: SSOHost, Path: "/login"})

	// Step 1: 获取 flowkey 和 cryptoKey（使用 WebVPN 转换后的 URL）
	flowkey, cryptoKey, err := c.getFlowkeyCryptoFrom(ctx, LoginURL.String())
	if err != nil {
		return fmt.Errorf("get flowkey/crypto: %w", err)
	}

	// Step 2: 加密密码
	encryptedPwd, err := crypto.EncryptPasswordAES(cryptoKey, password)
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	// Step 3: POST 登录表单
	form := url.Values{
		"username":     {username},
		"type":         {"UsernamePassword"},
		"_eventId":     {"submit"},
		"geolocation":  {""},
		"execution":    {flowkey},
		"captcha_code": {""},
		"croypto":      {cryptoKey}, // typo in original API
		"password":     {encryptedPwd},
	}

	loginReq, err := http.NewRequestWithContext(ctx, "POST", LoginURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginReq.Header.Set("Referer", LoginURL.String())
	loginReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	loginReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q,image/webp=0.9,*/*;q=0.8")
	loginReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	loginReq.Header.Set("Origin", "https://sso.hdu.edu.cn")

	loginResp, err := c.client.GetHTTPClient().Do(loginReq)
	if err != nil {
		return err
	}
	defer func() { _ = loginResp.Body.Close() }()

	// Step 4: 检查响应，如果是 4xx 则认为登录失败，返回登录失败错误，如果是 3xx 则认为登录成功（SSO 成功后会重定向回业务系统）
	if loginResp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(loginResp.Body)
		return fmt.Errorf("%w: status %d, body: %s", ErrLoginFailed, loginResp.StatusCode, string(bodyBytes))
	}
	if !isRedirect(loginResp.StatusCode) {
		bodyBytes, _ := io.ReadAll(loginResp.Body)
		return fmt.Errorf("%w: expected redirect but got status %d, body: %s", ErrLoginFailed, loginResp.StatusCode, string(bodyBytes))
	}
	finalHost := loginResp.Header.Get("Location")
	if finalHost == "" {
		return fmt.Errorf("%w: missing Location header", ErrLoginFailed)
	}
	c.logger.InfoContext(ctx, "sso login successful", "username", username, "final_host", finalHost)
	return nil
}

// getFlowkeyCryptoFrom 从 SSO 登录页提取 flowkey 和 cryptoKey
func (c *Context) getFlowkeyCryptoFrom(ctx context.Context, loginURL string) (flowkey, cryptoKey string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", loginURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Referer", loginURL)

	resp, err := c.client.GetHTTPClient().Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	flowkey = doc.Find("#login-page-flowkey").Text()
	cryptoKey = doc.Find("#login-croypto").Text()

	if flowkey == "" || cryptoKey == "" {
		return "", "", ErrGetFlowkey
	}

	return strings.TrimSpace(flowkey), strings.TrimSpace(cryptoKey), nil
}

// Do sends a request by building a Context and starting the middleware chain
func Do(client ClientInterface, req *request.Request) (*request.Response, error) {
	handlers := []Handler{
		serviceAuthHandler,
		transportHandler,
	}
	ctx := NewContext(client, req, handlers)
	ctx.Next()
	return ctx.Response, ctx.Err
}
