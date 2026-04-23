package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"

	"github.com/U1traVeno/hduwebvpn/pkg/sso"
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
	defer func() { _ = httpResp.Body.Close() }()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		c.Err = fmt.Errorf("read response body: %w", err)
		return
	}

	c.Response = &request.Response{
		RawResponse: httpResp,
		StatusCode:  httpResp.StatusCode,
		Header:      httpResp.Header,
		Body:        bodyBytes,
	}
}

// transportHandler 处理通道层认证（WebVPN 掉线重连）
func transportHandler(c *Context) {
	t := c.client.GetTransport()
	c.Request.RealURL = t.Encode(c.Request.BusinessURL)

	c.Next()

	if c.Err != nil {
		return
	}

	if decoded := c.client.GetTransport().Decode(c.Response.RawResponse.Request.URL); decoded != nil {
		c.Response.BusinessReqURL = decoded
	} else {
		c.Response.BusinessReqURL = c.Response.RawResponse.Request.URL
	}

	resp := c.Response
	// http.Client 已经默认会跟随重定向。因此，最终的响应应该要么是被重定向到 webvpn 的登录页，要么就是最终的业务响应。
	// 其中如果是前者，那么 RawResponse.Request.URL 应当是 sso.hdu.edu.cn, 此时进行 webvpn Reauth
	// 如果是后者，那么 RawResponse.Request.URL 应当是业务系统的 URL，此时通过了 webvpn 的认证，无需 ReAuth。
	if t.IsAuthFailure(resp.RawResponse.Request.URL) {
		c.logger.InfoContext(
			context.Background(),
			"received WebVPN SSO redirect",
		)
		if err := t.Reauth(c.client); err != nil {
			c.Err = fmt.Errorf("transport reauth failed: %w", err)
			return
		}
		// 重新执行请求
		transportHandler(c)
		return
	}
}

// serviceAuthHandler 处理业务层认证（CAS SSO）
func serviceAuthHandler(c *Context) {
	c.Next()

	if c.Err != nil {
		return
	}

	resp := c.Response
	// http.Client 已经默认会跟随重定向。因此，最终的响应应该要么是被重定向到 SSO 登录页，要么就是最终的业务响应。
	// 如果是前者，则直接执行 SSO 认证；如果是后者，则说明请求已经成功完成，无需再认证。
	if resp.BusinessReqURL == nil || !sso.IsAuthFailure(resp.BusinessReqURL.Host) {
		return
	}

	c.logger.InfoContext(
		context.Background(),
		"received SSO redirect",
		"location", resp.Header.Get("Location"),
	)
	username := c.client.GetUsername()
	password := c.client.GetPassword()

	if _, err := sso.Auth(context.Background(), c.client.GetHTTPClient(), resp.RawResponse.Request.URL.String(), username, password); err != nil {
		c.Err = fmt.Errorf("sso auth failed: %w", err)
		return
	}
	// 重新执行请求
	serviceAuthHandler(c)
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
