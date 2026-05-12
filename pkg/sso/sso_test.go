package sso

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func TestTryGetTicketFromSession_Success(t *testing.T) {
	// 模拟 SSO 服务器：GET 登录页直接 302 重定向到 callback（已有活跃会话）
	ssoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/callback/cas/test-external-id?ticket=ST-test-ticket-123", http.StatusFound)
	}))
	defer ssoSrv.Close()

	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	ticket, redirectURL, err := tryGetTicketFromSession(
		context.Background(),
		httpClient,
		ssoSrv.URL+"/login?service=https://webvpn.hdu.edu.cn/callback/cas/test-external-id",
	)
	if err != nil {
		t.Fatalf("tryGetTicketFromSession failed: %v", err)
	}
	if ticket != "ST-test-ticket-123" {
		t.Errorf("expected ticket 'ST-test-ticket-123', got %q", ticket)
	}
	if redirectURL == "" {
		t.Error("expected non-empty redirectURL")
	}
	t.Logf("ticket: %s, redirectURL: %s", ticket, redirectURL)
}

func TestTryGetTicketFromSession_NoSession(t *testing.T) {
	// 模拟 SSO 服务器：GET 登录页返回 200 登录表单（没有活跃会话）
	ssoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><body>
			<span id="login-page-flowkey">test-flowkey</span>
			<span id="login-croypto">dGVzdC1jcnlwdG8tMDEyMzQ1Njc=</span>
		</body></html>`))
	}))
	defer ssoSrv.Close()

	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	_, _, err := tryGetTicketFromSession(
		context.Background(),
		httpClient,
		ssoSrv.URL+"/login",
	)
	if err == nil {
		t.Fatal("expected error for non-redirect response")
	}
	t.Logf("got expected error: %v", err)
}

func TestTryGetTicketFromSession_NoTicketInLocation(t *testing.T) {
	// 模拟 SSO 服务器：302 但 Location 不带 ticket
	ssoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/some-other-page", http.StatusFound)
	}))
	defer ssoSrv.Close()

	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	_, _, err := tryGetTicketFromSession(
		context.Background(),
		httpClient,
		ssoSrv.URL+"/login",
	)
	if err == nil {
		t.Fatal("expected error when Location has no ticket")
	}
	t.Logf("got expected error: %v", err)
}

func TestAuth_SessionReuse(t *testing.T) {
	// 模拟完整的已认证场景：
	// SSO 登录页 → 302 → callback?ticket=xxx（已有会话，跳过表单登录）
	callbackCalled := false

	callbackSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callbackCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer callbackSrv.Close()

	// SSO 服务器返回绝对 URL 重定向到 callback server
	redirectTarget := callbackSrv.URL + "/callback/cas/test-ext?ticket=ST-reuse-ticket-456"
	ssoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", redirectTarget)
		w.WriteHeader(http.StatusFound)
	}))
	defer ssoSrv.Close()

	loginURL := ssoSrv.URL + "/login"

	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	ticket, err := Auth(
		context.Background(),
		httpClient,
		loginURL,
		"testuser",
		"testpassword",
	)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	if ticket != "ST-reuse-ticket-456" {
		t.Errorf("expected ticket 'ST-reuse-ticket-456', got %q", ticket)
	}
	if !callbackCalled {
		t.Error("callback URL was never followed")
	}
	t.Logf("session reuse successful, ticket: %s", ticket)
}

func TestAuth_NormalFlow(t *testing.T) {
	// 模拟正常的登录流程：无活跃会话，需要获取 flowkey/cryptoKey 并提交表单
	// 使用有效的 16 字节 AES 密钥 (base64 编码)
	cryptoKey := base64.StdEncoding.EncodeToString([]byte("test-key-0123456")) // 16 bytes
	var formSubmitCalled bool

	callbackSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer callbackSrv.Close()

	redirectTarget := callbackSrv.URL + "/callback/cas/normal-ext?ticket=ST-normal-ticket-789"

	ssoSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			w.Header().Set("Content-Type", "text/html")
			// 使用 span 而非 input，因为 goquery.Text() 对 input 返回空
			w.Write([]byte(`<html><body>
				<span id="login-page-flowkey">normal-flowkey</span>
				<span id="login-croypto">` + cryptoKey + `</span>
			</body></html>`))
		case r.Method == "POST":
			formSubmitCalled = true
			w.Header().Set("Location", redirectTarget)
			w.WriteHeader(http.StatusFound)
		}
	}))
	defer ssoSrv.Close()

	loginURL := ssoSrv.URL + "/login"

	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: jar}

	ticket, err := Auth(
		context.Background(),
		httpClient,
		loginURL,
		"testuser",
		"testpassword",
	)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	if ticket != "ST-normal-ticket-789" {
		t.Errorf("expected ticket 'ST-normal-ticket-789', got %q", ticket)
	}
	if !formSubmitCalled {
		t.Error("login form was never submitted")
	}
	t.Logf("normal flow successful, ticket: %s", ticket)
}
