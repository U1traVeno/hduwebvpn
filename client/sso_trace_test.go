package client_test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/U1traVeno/hduwebvpn/pkg/sso"
	"github.com/joho/godotenv"
)

func TestSSOLoginTrace(t *testing.T) {
	_ = godotenv.Load("../.env")

	username := os.Getenv("HDU_USER")
	password := os.Getenv("HDU_PASSWORD")

	t.Logf("username: %q", username)
	t.Logf("password: %q", password)
	t.Logf("password bytes: %v", []byte(password))

	if username == "" || password == "" {
		t.Skip("HDU_USER or HDU_PASSWORD not set")
	}

	jar, _ := cookiejar.New(nil)
	httpClient := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Step 1: Get auth list
	authListURL := "https://webvpn.hdu.edu.cn/api/access/authentication/list"
	req, _ := http.NewRequest("GET", authListURL, nil)
	req.Header.Set("Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("get auth list failed: %v", err)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("auth list status: %d, body: %s", resp.StatusCode, string(bodyBytes))

	var authListResp struct {
		Code int `json:"code"`
		Data struct {
			List []struct {
				ExternalID string `json:"externalId"`
			} `json:"list"`
		} `json:"data"`
	}
	json.Unmarshal(bodyBytes, &authListResp)
	externalID := authListResp.Data.List[0].ExternalID
	t.Logf("externalID from auth list: %s", externalID)

	// Step 2: Start auth with correct externalID
	callbackURL := fmt.Sprintf("https://webvpn.hdu.edu.cn/callback/cas/%s", externalID)
	startURL := "https://webvpn.hdu.edu.cn/api/access/auth/start"
	reqBody := fmt.Sprintf(`{"externalId":"%s","data":"{\"callbackUrl\":\"%s\"}"}`, externalID, callbackURL)
	req, _ = http.NewRequest("POST", startURL, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err = httpClient.Do(req)
	if err != nil {
		t.Fatalf("start auth failed: %v", err)
	}
	bodyBytes, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("start auth status: %d, body: %s", resp.StatusCode, string(bodyBytes))

	var startResp struct {
		Code int `json:"code"`
		Data struct {
			Type   int `json:"type"`
			Action struct {
				LoginURL string `json:"login_url"`
			} `json:"action"`
		} `json:"data"`
	}
	json.Unmarshal(bodyBytes, &startResp)
	ssoLoginURL := startResp.Data.Action.LoginURL
	if ssoLoginURL == "" {
		t.Fatal("no sso login url from start auth response")
	}
	t.Logf("SSO login URL: %s", ssoLoginURL)

	// Step 3: Get SSO login page to extract flowkey and cryptoKey
	req, _ = http.NewRequest("GET", ssoLoginURL, nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	resp, err = httpClient.Do(req)
	if err != nil {
		t.Fatalf("fetch SSO login page failed: %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("parse SSO login page: %v", err)
	}

	flowkey := strings.TrimSpace(doc.Find("#login-page-flowkey").Text())
	cryptoKey := strings.TrimSpace(doc.Find("#login-croypto").Text())
	t.Logf("flowkey: %q", flowkey)
	t.Logf("cryptoKey: %q", cryptoKey)

	if flowkey == "" || cryptoKey == "" {
		t.Fatal("missing flowkey or cryptoKey")
	}

	// Step 4: Encrypt password
	encryptedPwd, err := sso.EncryptPasswordAES(cryptoKey, password)
	if err != nil {
		t.Fatalf("encrypt password: %v", err)
	}
	t.Logf("encrypted password: %q", encryptedPwd)

	// Step 5: Submit login form
	form := url.Values{
		"username":     {username},
		"type":         {"UsernamePassword"},
		"_eventId":     {"submit"},
		"geolocation":  {""},
		"execution":    {flowkey},
		"captcha_code": {""},
		"croypto":      {cryptoKey},
		"password":     {encryptedPwd},
	}
	formEncoded := form.Encode()
	t.Logf("form encoded length: %d", len(formEncoded))
	t.Logf("form encoded: %s", formEncoded)

	req, _ = http.NewRequest("POST", ssoLoginURL, strings.NewReader(formEncoded))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", ssoLoginURL)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q,image/webp=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Origin", "https://sso.hdu.edu.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err = httpClient.Do(req)
	if err != nil {
		t.Fatalf("SSO login request failed: %v", err)
	}
	bodyBytes, _ = io.ReadAll(resp.Body)
	t.Logf("SSO login status: %d, body length: %d", resp.StatusCode, len(bodyBytes))
	t.Logf("SSO login headers: %v", resp.Header)

	doc2, _ := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	errorCode := strings.TrimSpace(doc2.Find("#login-error-code").Text())
	if errorCode != "" {
		t.Logf("error code from page: %s", errorCode)
	}

	location := resp.Header.Get("Location")
	if location != "" {
		t.Logf("Location header: %s", location)
		locURL, _ := url.Parse(location)
		if locURL != nil {
			t.Logf("ticket from location: %s", locURL.Query().Get("ticket"))
		}
	}
}
