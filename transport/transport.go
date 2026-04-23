package transport

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/U1traVeno/hduwebvpn/pkg/sso"
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
type WebVPNTransport struct {
	deviceID   string
	urlMap     map[string]string // 业务 URL host -> WebVPN URL host
	reverseMap map[string]string // WebVPN URL host -> 业务 URL host
}

const (
	webVPNHost    = "webvpn.hdu.edu.cn"
	authListURL   = "/api/access/authentication/list"
	authStartURL  = "/api/access/auth/start"
	authFinishURL = "/api/access/auth/finish"
	callbackPath  = "/callback/cas/"
	siteListURL   = "/api/access/nav/site-list"
)

// SiteInfo 存储站点信息，用于 URL 映射
type SiteInfo struct {
	ID        int
	Name      string
	RawURL    string
	WebVPNURL string
}

// getDeviceID 返回一个固定的 device ID（32-char MD5 格式）
func (t *WebVPNTransport) getDeviceID() string {
	// 此处需要为每一个 WebVPNTransport 维护一个独立的 deviceID, 避免本项目被识别为同一设备导致认证失败。
	if t.deviceID == "" {
		// 生成一个随机的 device ID（32-char MD5 格式）
		randomStr := fmt.Sprintf("%d", time.Now().UnixNano())
		t.deviceID = fmt.Sprintf("%x", md5.Sum([]byte(randomStr)))
	}
	return t.deviceID
}

// getAuthList 获取认证方式列表，提取 externalId
func (t *WebVPNTransport) getAuthList(ctx context.Context, httpClient *http.Client) (string, error) {
	authListURL := fmt.Sprintf("https://%s%s", webVPNHost, authListURL)

	req, err := http.NewRequestWithContext(ctx, "GET", authListURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request auth list: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth list failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read auth list response: %w", err)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				ExternalID string `json:"externalId"`
				Name       string `json:"name"`
				AuthType   int    `json:"authType"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("parse auth list response: %w", err)
	}

	if result.Code != 0 || len(result.Data.List) == 0 {
		return "", fmt.Errorf("auth list returned code %d", result.Code)
	}

	return result.Data.List[0].ExternalID, nil
}

// startAuth 发起认证，获取 SSO 登录 URL
func (t *WebVPNTransport) startAuth(ctx context.Context, httpClient *http.Client, externalID, callbackURL string) (string, error) {
	startURL := fmt.Sprintf("https://%s%s", webVPNHost, authStartURL)

	reqBody := map[string]interface{}{
		"externalId": externalID,
		"data":       fmt.Sprintf(`{"callbackUrl":"%s"}`, callbackURL),
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal start auth request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", startURL, strings.NewReader(string(reqBodyBytes)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request start auth: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("start auth failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read start auth response: %w", err)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			Type   int `json:"type"`
			Action struct {
				LoginURL string `json:"login_url"`
			} `json:"action"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("parse start auth response: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("start auth returned code %d", result.Code)
	}

	return result.Data.Action.LoginURL, nil
}

// finishAuth 完成认证，获取 webvpn-token
func (t *WebVPNTransport) finishAuth(ctx context.Context, httpClient *http.Client, externalID, callbackURL, ticket string) error {
	finishURL := fmt.Sprintf("https://%s%s", webVPNHost, authFinishURL)

	reqBody := map[string]interface{}{
		"externalId": externalID,
		"data": fmt.Sprintf(`{"callbackUrl":"%s","ticket":"%s","deviceId":"%s"}`,
			callbackURL, ticket, t.getDeviceID()),
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal finish auth request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", finishURL, strings.NewReader(string(reqBodyBytes)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request finish auth: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("finish auth failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read finish auth response: %w", err)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return fmt.Errorf("parse finish auth response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("finish auth returned code %d", result.Code)
	}

	// 获取站点列表，更新 URL 映射
	err = t.FetchSiteList(ctx, httpClient)
	if err != nil {
		return fmt.Errorf("fetch site list: %w", err)
	}

	return nil
}

// Reauth 重新执行 WebVPN 通道认证，获取新的 webvpn-token
func (t *WebVPNTransport) Reauth(client interface{}) error {
	// 获取 HTTP client 和凭证
	var httpClient *http.Client
	var username, password string

	switch c := client.(type) {
	case interface {
		GetHTTPClient() *http.Client
		GetUsername() string
		GetPassword() string
	}:
		httpClient = c.GetHTTPClient()
		username = c.GetUsername()
		password = c.GetPassword()
	default:
		return fmt.Errorf("unsupported client type: %T", client)
	}

	if httpClient == nil {
		return errors.New("nil HTTP client")
	}

	ctx := context.Background()

	// Step 1: 获取 externalId
	externalID, err := t.getAuthList(ctx, httpClient)
	if err != nil {
		return fmt.Errorf("get auth list: %w", err)
	}

	callbackURL := fmt.Sprintf("https://%s%s%s", webVPNHost, callbackPath, externalID)

	// Step 2: 获取 SSO 登录 URL
	ssoLoginURL, err := t.startAuth(ctx, httpClient, externalID, callbackURL)
	if err != nil {
		return fmt.Errorf("start auth: %w", err)
	}

	// Step 3: 登录 SSO，获取 ticket
	ticket, err := sso.Auth(ctx, httpClient, ssoLoginURL, username, password)
	if err != nil {
		return fmt.Errorf("SSO login: %w", err)
	}

	// Step 4: 完成认证，获取 token（Set-Cookie 会自动存入 cookie jar）
	err = t.finishAuth(ctx, httpClient, externalID, callbackURL, ticket)
	if err != nil {
		return fmt.Errorf("finish auth: %w", err)
	}

	return nil
}

// FetchSiteList 获取站点列表并构建 URL 映射表
func (t *WebVPNTransport) FetchSiteList(ctx context.Context, httpClient *http.Client) error {
	listURL := fmt.Sprintf("https://%s%s", webVPNHost, siteListURL)

	req, err := http.NewRequestWithContext(ctx, "GET", listURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request site list: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("site list failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read site list response: %w", err)
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			List []struct {
				Name  string `json:"name"`
				Sites []struct {
					ID      int    `json:"id"`
					Name    string `json:"name"`
					URL     string `json:"url"`
					Icon    string `json:"icon"`
					Color   string `json:"color"`
					IconURL string `json:"iconUrl"`
					Sort    int    `json:"sort"`
					RawURL  string `json:"rawURL"`
				} `json:"sites"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return fmt.Errorf("parse site list response: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("site list returned code %d: %s", result.Code, result.Message)
	}

	// 构建 URL 映射表
	t.urlMap = make(map[string]string)
	t.reverseMap = make(map[string]string)

	for _, category := range result.Data.List {
		for _, site := range category.Sites {
			if site.RawURL == "" || site.URL == "" {
				continue
			}

			rawURL, err := url.Parse(site.RawURL)
			if err != nil {
				continue
			}

			webVPNURL, err := url.Parse(site.URL)
			if err != nil {
				continue
			}

			// 存储 host 映射
			t.urlMap[rawURL.Host] = webVPNURL.Host
			t.reverseMap[webVPNURL.Host] = rawURL.Host
		}
	}

	return nil
}

func (t *WebVPNTransport) Encode(businessURL *url.URL) *url.URL {
	if businessURL == nil {
		return nil
	}

	// 先查 urlMap
	if webVPNHost, ok := t.urlMap[businessURL.Host]; ok {
		encoded := *businessURL
		encoded.Host = webVPNHost
		return &encoded
	}

	// 如果 urlMap 中没有，使用动态转换逻辑
	// 格式: https://https-{host}-{port}.webvpn.hdu.edu.cn
	// 例如: https://course.hdu.edu.cn -> https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn
	host := businessURL.Host
	scheme := businessURL.Scheme
	port := businessURL.Port()

	// 构建 WebVPN 格式的 host
	// host 如 "course.hdu.edu.cn" -> "https-course-hdu-edu-cn-443"
	webVPNHostPart := strings.ReplaceAll(host, ".", "-")
	if port != "" && port != "443" && port != "80" {
		webVPNHostPart = fmt.Sprintf("%s-%s", webVPNHostPart, port)
	} else {
		webVPNHostPart = fmt.Sprintf("%s-443", webVPNHostPart)
	}

	encoded := *businessURL
	encoded.Host = fmt.Sprintf("https-%s.webvpn.hdu.edu.cn", webVPNHostPart)
	encoded.Scheme = scheme
	return &encoded
}

func (t *WebVPNTransport) Decode(realURL *url.URL) *url.URL {
	if realURL == nil {
		return nil
	}

	// 先查 reverseMap
	if businessHost, ok := t.reverseMap[realURL.Host]; ok {
		decoded := *realURL
		decoded.Host = businessHost
		return &decoded
	}

	// 动态逆向解析:
	//   https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn -> https://course.hdu.edu.cn
	//   https://http-course-hdu-edu-cn-80.webvpn.hdu.edu.cn  -> http://course.hdu.edu.cn
	if !strings.HasSuffix(realURL.Host, ".webvpn.hdu.edu.cn") {
		return nil
	}

	hostWithoutSuffix := strings.TrimSuffix(realURL.Host, ".webvpn.hdu.edu.cn")

	// 提取原始 scheme（https- 或 http-）
	var scheme string
	var hostPart string
	switch {
	case strings.HasPrefix(hostWithoutSuffix, "https-"):
		scheme = "https"
		hostPart = strings.TrimPrefix(hostWithoutSuffix, "https-")
	case strings.HasPrefix(hostWithoutSuffix, "http-"):
		scheme = "http"
		hostPart = strings.TrimPrefix(hostWithoutSuffix, "http-")
	default:
		return nil
	}

	// 处理端口后缀：-443、-80 或其他 -{port}
	// 从末尾开始找最后一个 "-"，若其后是纯数字则视为端口
	if lastDash := strings.LastIndex(hostPart, "-"); lastDash != -1 {
		portStr := hostPart[lastDash+1:]
		if _, err := strconv.Atoi(portStr); err == nil {
			hostPart = hostPart[:lastDash]
		}
	}

	originalHost := strings.ReplaceAll(hostPart, "-", ".")

	decoded := *realURL
	decoded.Scheme = scheme
	decoded.Host = originalHost
	return &decoded
}

func (t *WebVPNTransport) IsAuthFailure(realURL *url.URL) bool {
	// 如果重定向到了非 webvpn 域名，说明掉出了 WebVPN 环境
	if realURL.Host == "" {
		return false
	}
	// webvpn 认证失效会被重定向到类似:
	// https://webvpn.hdu.edu.cn?returnUrl=https%3A%2F%2Fhttps-course-hdu-edu-cn-443.webvpn.hdu.edu.cn
	// 此时 Host 为 webvpn.hdu.edu.cn, 这种情况已经被下面的 HasSuffix 判断捕获（因为它不以 '.webvpn.hdu.edu.cn' 结尾）。
	return !strings.HasSuffix(realURL.Host, ".webvpn.hdu.edu.cn")
}
