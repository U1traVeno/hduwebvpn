package sso

import (
	"context"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	// ErrLoginFailed 表示 SSO 登录失败。
	ErrLoginFailed = errors.New("sso login failed")
	// ErrGetFlowkey 表示获取 flowkey/crypto 失败。
	ErrGetFlowkey = errors.New("failed to get flowkey/crypto")
)

const (
	// SSOHost 常量（用于域名比较）
	SSOHost = "sso.hdu.edu.cn"
	CASHost = "cas.hdu.edu.cn"
)

func IsAuthFailure(host string) bool {
	// 判断是否需要重新认证
	// 这里简单判断如果被重定向到了 SSO 登录页，则认为是认证失败
	return host == SSOHost || host == CASHost ||
		strings.HasPrefix(host, "sso-") || strings.HasPrefix(host, "cas-") ||
		strings.HasPrefix(host, "sso.") || strings.HasPrefix(host, "cas.")
}

// getFlowkeyCryptoFrom 从 SSO 登录页提取 flowkey 和 cryptoKey
func getFlowkeyCryptoFrom(ctx context.Context, httpClient *http.Client, loginURL string) (flowkey, cryptoKey string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", loginURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Referer", loginURL)

	resp, err := httpClient.Do(req)
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

// EncryptPasswordAES 使用 AES-ECB + PKCS7 padding 对密码进行加密。
// 这是 HDU SSO 系统的标准密码加密方式。
//
// 参数：
//   - cryptoKey: Base64 编码的 AES 密钥（从 SSO 登录页获取）
//   - password: 明文密码
//
// 返回：
//   - Base64 编码的加密密文
//   - 错误（如果加密失败）
func EncryptPasswordAES(cryptoKey, password string) (string, error) {
	// 1. 解码 Base64 密钥
	keyBytes, err := base64.StdEncoding.DecodeString(cryptoKey)
	if err != nil {
		return "", fmt.Errorf("decode crypto key: %w", err)
	}

	// 2. 验证 AES 密钥长度（16/24/32 字节）
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", fmt.Errorf("invalid crypto key length: got %d bytes, want 16/24/32", len(keyBytes))
	}

	// 3. 创建 AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("create aes cipher: %w", err)
	}

	// 4. PKCS7 padding
	plaintext := []byte(password)
	blockSize := block.BlockSize()
	padding := blockSize - len(plaintext)%blockSize
	padtext := make([]byte, len(plaintext)+padding)
	copy(padtext, plaintext)
	for i := len(plaintext); i < len(padtext); i++ {
		padtext[i] = byte(padding)
	}

	// 5. AES-ECB 加密（逐块独立加密）
	ciphertext := make([]byte, len(padtext))
	for i := 0; i < len(padtext); i += blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], padtext[i:i+blockSize])
	}

	// 6. Base64 编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func MD5Hash(text string) string {
	// 计算 MD5 哈希
	hash := md5.Sum([]byte(text))
	// 返回十六进制字符串
	return fmt.Sprintf("%x", hash)
}

// Auth 执行完整的 SSO 登录流程，返回登录成功后的 ticket。
func Auth(
	ctx context.Context,
	httpClient *http.Client,
	ssoLoginURL,
	username,
	password string,
) (string, error) {
	flowkey, cryptoKey, err := getFlowkeyCryptoFrom(ctx, httpClient, ssoLoginURL)
	if err != nil {
		return "", fmt.Errorf("get flowkey/crypto: %w", err)
	}

	encryptedPwd, err := EncryptPasswordAES(cryptoKey, password)
	if err != nil {
		return "", fmt.Errorf("encrypt password: %w", err)
	}

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

	req, err := http.NewRequestWithContext(ctx, "POST", ssoLoginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ssoLoginURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Origin", "https://sso.hdu.edu.cn")

	// 使用临时 client 禁止自动重定向，以便从 302 响应的 Location header 中直接提取 ticket
	tempClient := &http.Client{
		Transport:     httpClient.Transport,
		Jar:           httpClient.Jar,
		Timeout:       httpClient.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}

	resp, err := tempClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("SSO login request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("SSO login failed: expected redirect but got status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", fmt.Errorf("SSO login failed: missing Location header")
	}

	u, err := url.Parse(loc)
	if err != nil {
		return "", fmt.Errorf("SSO login failed: invalid Location %q", loc)
	}
	ticket := u.Query().Get("ticket")
	if ticket == "" {
		return "", fmt.Errorf("SSO login failed: no ticket in redirect location")
	}

	// 使用与原 client 完成后续重定向链，
	// 确保服务端在跳转过程中设置的 cookies（如业务系统 session）被正确保存。
	followReq, err := http.NewRequestWithContext(ctx, http.MethodGet, loc, nil)
	if err == nil {
		followResp, err := httpClient.Do(followReq)
		if err == nil && followResp != nil {
			_, _ = io.Copy(io.Discard, followResp.Body)
			_ = followResp.Body.Close()
		}
	}

	return ticket, nil
}
