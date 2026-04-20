package handler

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/U1traVeno/hduwebvpn/request"
	"github.com/U1traVeno/hduwebvpn/transport"
)

// ClientInterface defines what the middleware needs from the client
type ClientInterface interface {
	GetHTTPClient() *http.Client
	GetCookieJar() *cookiejar.Jar
	GetTransport() transport.Transport
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
}

// NewContext constructs a Context with handlers
func NewContext(client ClientInterface, req *request.Request, handlers []Handler) *Context {
	return &Context{
		Request:  req,
		Response: &request.Response{},
		index:    -1,
		handlers: handlers,
		client:   client,
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
		RealReq: &request.RealRequest{
			URL:    c.Request.RealURL.String(),
			Method: c.Request.Method,
			Header: httpReq.Header,
			Body:   c.Request.Body,
		},
		RawResponse: httpResp,
		StatusCode:  httpResp.StatusCode,
		Header:      httpResp.Header,
	}
}

// isRedirect 判断状态码是否为重定向
func isRedirect(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
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
		if strings.Contains(resp.BusinessRedirectURL.Host, "sso.hdu.edu.cn") {
			// TODO: 实现 CAS SSO 认证流程
			serviceAuthHandler(c)
			return
		}
	}
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