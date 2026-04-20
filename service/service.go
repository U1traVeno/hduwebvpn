package service

import (
	"net/url"

	"github.com/U1traVeno/hduwebvpn/request"
)

// Doer interface for making requests
type Doer interface {
	Do(req *request.Request) (*request.Response, error)
}

// Service represents an registered intranet service
type Service struct {
	name    string
	baseURL string
	doer    Doer
}

// NewService creates a new service
func NewService(name string, baseURL string) *Service {
	return &Service{
		name:    name,
		baseURL: baseURL,
	}
}

// SetClient sets the doer for this service
func (s *Service) SetClient(c Doer) {
	s.doer = c
}

// NewRequest 创建一个新的请求，完全控制请求头（类似 http.Client）
func (s *Service) NewRequest(method, path string, body []byte) (*request.Request, error) {
	base, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	// 拼接 business URL
	businessURL := base.ResolveReference(&url.URL{Path: path})

	req := &request.Request{
		Service:     s,
		Method:      method,
		BusinessURL: businessURL,
		Header:      make(map[string][]string),
		Body:        body,
	}
	return req, nil
}

// Get 便捷方法，发起 GET 请求
func (s *Service) Get(path string) (*request.Response, error) {
	req, err := s.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return s.doer.Do(req)
}

// Post 便捷方法，发起 POST 请求
func (s *Service) Post(path string, body []byte) (*request.Response, error) {
	req, err := s.NewRequest("POST", path, body)
	if err != nil {
		return nil, err
	}
	return s.doer.Do(req)
}

// doSSO 执行业务层 SSO 认证
func (s *Service) DoSSO() error {
	// TODO: 实现 CAS SSO 认证流程
	return nil
}