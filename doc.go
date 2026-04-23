// Package hduwebvpn provides a client library for accessing Hangzhou Dianzi University
// intranet services through both WebVPN (off-campus) and Direct (on-campus) modes.
//
// It handles the complexity of WebVPN authentication, CAS SSO login, URL encoding/decoding,
// cookie management, and automatic re-authentication on token expiration. Users can register
// intranet services (e.g., course.hdu.edu.cn) and send requests using business-level URLs
// while the library transparently translates them to real network addresses.
//
// # Core Design
//
// The library is built around three key abstractions:
//
//   - Transport: Encodes business URLs to real URLs (WebVPN wraps hostnames) and decodes
//     them back. Also detects authentication failures and triggers re-auth.
//   - Client: Manages credentials, cookie jar, HTTP client, and registered services.
//     Each Client instance is fully isolated with its own login state.
//   - Service: Represents a registered intranet service. Provides convenient Get/Post
//     methods as well as full-control NewRequest for custom headers.
//
// Requests flow through a Gin-style middleware chain (transportHandler → serviceAuthHandler
// → execBaseDo), allowing pre/post processing like URL translation and redirect handling.
//
// # Usage
//
// Create a client with your credentials. By default it uses WebVPN mode:
//
//	client, err := hduwebvpn.NewClient(
//	    hduwebvpn.WithUsername("24270123"),
//	    hduwebvpn.WithPassword("password"),
//	)
//
// For on-campus direct access without WebVPN tunneling:
//
//	client, err := hduwebvpn.NewClient(
//	    hduwebvpn.WithUsername("24270123"),
//	    hduwebvpn.WithPassword("password"),
//	    hduwebvpn.WithMode(hduwebvpn.DirectMode),
//	)
//
// Register a service and send requests. RegisterService returns an interface{} that
// must be type-asserted to *hduwebvpn.Service:
//
//	client.RegisterService("course", "https://course.hdu.edu.cn")
//	svc := client.Service("course").(*hduwebvpn.Service)
//
//	// Simple GET
//	resp, err := svc.Get("/api/access/user/info")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(resp.StatusCode)
//	fmt.Println(string(resp.Body))
//
// For full control over headers, build a Request manually:
//
//	req, err := svc.NewRequest("GET", "/api/endpoint", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	req.Header.Set("Accept", "application/json")
//	req.Header.Set("X-Custom-Header", "custom-value")
//	resp, err := svc.Do(req)
//
// Access cookies managed by the client:
//
//	cookies := client.Cookies("course.hdu.edu.cn")
//	for _, c := range cookies {
//	    fmt.Printf("%s=%s\n", c.Name, c.Value)
//	}
//
// Manually trigger re-authentication (useful for probing login flow):
//
//	tp := client.GetTransport()
//	err = tp.Reauth(client)
//
// Inspect URL encoding/decoding (WebVPN mode converts hosts to the webvpn.hdu.edu.cn domain):
//
//	businessURL, _ := url.Parse("https://course.hdu.edu.cn/path")
//	realURL := client.GetTransport().Encode(businessURL)
//	fmt.Println(realURL.String()) // https://https-course-hdu-edu-cn-443.webvpn.hdu.edu.cn/path
//
//	decoded := client.GetTransport().Decode(realURL)
//	fmt.Println(decoded.String()) // https://course.hdu.edu.cn/path
//
// # Automatic Behaviors
//
//   - Cookie Management: The underlying http.Client uses a cookie jar. Set-Cookie
//     headers are automatically stored and sent for matching domains.
//   - Auto Re-auth: If a request fails due to expired WebVPN or service authentication,
//     the library automatically re-authenticates using the stored username/password and
//     retries the request.
//   - Redirect Handling: HTTP 3xx redirects are followed automatically by the underlying
//     http.Client. After the request completes, middleware inspects the final request URL
//     to detect authentication failures (e.g., redirected to an SSO login page).
//
// # Response Body
//
// Response body data is available directly via resp.Body ([]byte). The body has already
// been read and closed by the middleware chain.
package hduwebvpn
