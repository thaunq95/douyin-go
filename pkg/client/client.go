package client

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/thaunq95/douyin-go/pkg/abogus"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"
	ttwidRegisterURL = "https://ttwid.bytedance.com/ttwid/union/register/"
	ttwidPayload     = `{"region":"cn","aid":1768,"needFid":false,"service":"www.ixigua.com","migrate_info":{"ticket":"","source":"node"},"cbUrlProtocol":"https","union":true}`
)

// EncodeQueryParams encodes query parameters in a specific order required by Douyin's signature algorithm,
// which is required for correct signature verification by Douyin.
func EncodeQueryParams(params map[string]string) string {
	orderedKeys := []string{
		"device_platform", "aid", "channel", "pc_client_type",
		"version_code", "version_name", "cookie_enabled",
		"screen_width", "screen_height", "browser_language",
		"browser_platform", "browser_name", "browser_version",
		"browser_online", "engine_name", "engine_version",
		"os_name", "os_version", "cpu_core_num", "device_memory",
		"platform", "downlink", "effective_type", "from_user_page",
		"locate_query", "need_time_list", "pc_libra_divert",
		"publish_video_strategy_type", "round_trip_time",
		"show_live_replay_strategy", "time_list_query",
		"whale_cut_token", "update_version_code", "msToken",
		"max_cursor", "count", "sec_user_id", "aweme_id",
	}

	visited := make(map[string]bool)
	var pairs []string

	for _, k := range orderedKeys {
		if v, ok := params[k]; ok {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
			visited[k] = true
		}
	}

	for k, v := range params {
		if !visited[k] {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}

	return strings.Join(pairs, "&")
}

// Client wraps an http.Client with automatic header injection and cookie management
type Client struct {
	HTTPClient *http.Client
	UserAgent  string
	Cookie     string
}

// NewClient creates a new Douyin client instance with a standard cookie jar and optional proxy
func NewClient(cookie string, proxyStr string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{}
	if proxyStr != "" {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL %s: %w", proxyStr, err)
		}
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	httpClient := &http.Client{
		Jar:       jar,
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	return &Client{
		HTTPClient: httpClient,
		UserAgent:  defaultUserAgent,
		Cookie:     cookie,
	}, nil
}

// GenerateFalseMSToken generates a fake msToken that matches the format expected by Douyin APIs
func GenerateFalseMSToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 126)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b) + "=="
}

// FetchTTWid requests a valid ttwid cookie from Douyin's servers
func (c *Client) FetchTTWid() (string, error) {
	req, err := http.NewRequest("POST", ttwidRegisterURL, bytes.NewBuffer([]byte(ttwidPayload)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "ttwid" {
			return cookie.Value, nil
		}
	}

	// Try extracting from Set-Cookie header manually if Jar didn't capture it
	setCookie := resp.Header.Get("Set-Cookie")
	if setCookie != "" {
		parts := strings.Split(setCookie, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "ttwid=") {
				return strings.TrimPrefix(part, "ttwid="), nil
			}
		}
	}

	return "", fmt.Errorf("ttwid cookie not found in registry response")
}

// SendRequest performs an HTTP request and injects necessary headers/cookies
func (c *Client) SendRequest(method string, reqURL string, query map[string]string, body io.Reader) ([]byte, error) {
	// 1. Build Query Parameters
	u, err := url.Parse(reqURL)
	if err != nil {
		return nil, err
	}

	// Ensure msToken is present in query parameters, even if empty
	if _, ok := query["msToken"]; !ok {
		query["msToken"] = ""
	}

	// 2. Generate a_bogus signature using ordered query parameters
	encodedQuery := EncodeQueryParams(query)
	signature := abogus.GenerateABogus(encodedQuery, c.UserAgent)
	finalQuery := fmt.Sprintf("%s&a_bogus=%s", encodedQuery, url.QueryEscape(signature))

	u.RawQuery = finalQuery
	finalURL := u.String()

	// 3. Create Request
	req, err := http.NewRequest(method, finalURL, body)
	if err != nil {
		return nil, err
	}

	// 4. Set Headers
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Referer", "https://www.douyin.com/?recommend=1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")

	if c.Cookie != "" {
		req.Header.Set("Cookie", c.Cookie)
	}

	// 5. Send Request with Retries
	var resp *http.Response
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		resp, lastErr = c.HTTPClient.Do(req)
		if lastErr == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed after retries: %w", lastErr)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Get performs a GET request with signed query params
func (c *Client) Get(reqURL string, query map[string]string) ([]byte, error) {
	return c.SendRequest("GET", reqURL, query, nil)
}

// Post performs a POST request with signed query params
func (c *Client) Post(reqURL string, query map[string]string, body io.Reader) ([]byte, error) {
	return c.SendRequest("POST", reqURL, query, body)
}
