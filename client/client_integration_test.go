package client_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/U1traVeno/hduwebvpn"
	"github.com/joho/godotenv"
)

func TestWebVPNClientAccessCourse(t *testing.T) {
	_ = godotenv.Load("../.env")

	username := os.Getenv("HDU_USER")
	password := os.Getenv("HDU_PASSWORD")

	t.Logf("username length: %d", len(username))
	t.Logf("password length: %d", len(password))
	t.Logf("password chars: %v", []byte(password))

	if username == "" || password == "" {
		t.Skip("HDU_USER or HDU_PASSWORD not set in .env")
	}

	client, err := hduwebvpn.NewClient(
		hduwebvpn.WithUsername(username),
		hduwebvpn.WithPassword(password),
	)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	svc := client.RegisterService("course", "https://course.hdu.edu.cn").(*hduwebvpn.Service)

	// 获取 transport 并触发 Reauth 以观察详细流程
	tp := client.GetTransport()
	t.Log("Starting Reauth...")
	err = tp.Reauth(client)
	if err != nil {
		t.Fatalf("Reauth failed: %v", err)
	}
	t.Log("Reauth successful")

	encoded := client.GetTransport().Encode(&url.URL{
		Scheme: "https",
		Host:   "course.hdu.edu.cn",
		Path:   "/",
	})
	expected := "https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn/"
	if encoded.String() != expected {
		t.Errorf("URL encoding wrong:\n  got:  %s\n  want: %s", encoded.String(), expected)
	}

	resp, err := svc.Get("/")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	t.Logf("Successfully accessed course.hdu.edu.cn via WebVPN, status: %d", resp.StatusCode)
}
