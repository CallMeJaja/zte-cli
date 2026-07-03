package router

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost     = "192.168.100.1"
	defaultUsername = "admin"
	defaultPassword = "admin"
	userAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

var (
	reLogintoken     = regexp.MustCompile(`createHiddenInput\("Frm_Logintoken",\s*"(\d+)"\)`)
	reLoginchecktoken = regexp.MustCompile(`createHiddenInput\("Frm_Loginchecktoken",\s*"(\d+)"\)`)
	reSessionToken   = regexp.MustCompile(`var session_token = "(\d+)";`)
	reTokenSetter    = regexp.MustCompile(`tokenInput\.setAttribute\("value",\s*(\d+)\)`)
)

// Client manages HTTP communication with the ZTE router.
type Client struct {
	Host          string
	Username      string
	Password      string
	BaseURL       string
	SessionToken  string
	httpClient    *http.Client
}

// NewClient creates a new router client with the given credentials.
func NewClient(host, username, password string) (*Client, error) {
	if host == "" {
		host = defaultHost
	}
	if username == "" {
		username = defaultUsername
	}
	if password == "" {
		password = defaultPassword
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &Client{
		Host:     host,
		Username: username,
		Password: password,
		BaseURL:  fmt.Sprintf("http://%s", host),
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 15 * time.Second,
		},
	}

	return client, nil
}

// Login authenticates with the router using the reverse-engineered flow:
// 1. GET login page to extract dynamic tokens
// 2. Generate random number, hash password with SHA256
// 3. POST login credentials
// 4. GET template.gch to establish session and extract session_token
func (c *Client) Login() (bool, error) {
	// Step 1: Fetch login page to get tokens
	req, err := http.NewRequest("GET", c.BaseURL+"/", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to reach router at %s: %w", c.Host, err)
	}
	defer resp.Body.Close()

	body := ReadResponseBody(resp)
	if body == "" {
		return false, fmt.Errorf("empty response from router")
	}

	// Extract tokens
	logintoken := extractMatch(reLogintoken, body)
	loginchecktoken := extractMatch(reLoginchecktoken, body)

	if logintoken == "" {
		logintoken = "4"
	}
	if loginchecktoken == "" {
		return false, fmt.Errorf("could not extract Frm_Loginchecktoken from login page")
	}

	// Step 2: Generate random number and hash password
	randNum := rand.Intn(89999999) + 10000000
	passwordHash := SHA256Hex(c.Password + strconv.Itoa(randNum))

	// Step 3: POST login
	payload := fmt.Sprintf(
		"action=login&Username=%s&Password=%s&Frm_Logintoken=%s&UserRandomNum=%d&Frm_Loginchecktoken=%s",
		c.Username, passwordHash, logintoken, randNum, loginchecktoken,
	)

	req, err = http.NewRequest("POST", c.BaseURL+"/", strings.NewReader(payload))
	if err != nil {
		return false, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body = ReadResponseBody(resp)

	if strings.Contains(body, "User information is error") {
		return false, fmt.Errorf("invalid username or password")
	}
	if strings.Contains(body, "failed for three times") {
		return false, fmt.Errorf("router locked due to too many failed attempts, wait 60 seconds")
	}

	// Step 4: Fetch template to establish session and get session_token
	req, err = http.NewRequest("GET", c.BaseURL+"/template.gch", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create template request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("template request failed: %w", err)
	}
	defer resp.Body.Close()

	body = ReadResponseBody(resp)

	if resp.StatusCode == 404 || strings.Contains(body, "404 Not Found") {
		return false, fmt.Errorf("login appeared to succeed but template.gch returned 404")
	}

	// Extract session_token from the main menu template
	c.SessionToken = extractMatch(reSessionToken, body)
	if c.SessionToken == "" {
		// Try alternative extraction from token setter
		c.SessionToken = extractMatch(reTokenSetter, body)
	}

	return true, nil
}

// GetPage fetches a .gch page from the router using the active session.
// It also refreshes the session_token from the response.
func (c *Client) GetPage(pageName string) (string, error) {
	url := fmt.Sprintf("%s/getpage.gch?pid=1002&nextpage=%s", c.BaseURL, pageName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body := ReadResponseBody(resp)

	if resp.StatusCode == 404 || strings.Contains(body, "404 Not Found") {
		return "", fmt.Errorf("page %s not found (404)", pageName)
	}

	// Refresh session token from page
	if token := extractMatch(reSessionToken, body); token != "" {
		c.SessionToken = token
	}

	return body, nil
}

// PostAction sends a form POST to a specific .gch page.
func (c *Client) PostAction(pageName string, formData string) (string, error) {
	url := fmt.Sprintf("%s/getpage.gch?pid=1002&nextpage=%s", c.BaseURL, pageName)

	req, err := http.NewRequest("POST", url, strings.NewReader(formData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body := ReadResponseBody(resp)

	// Refresh session token
	if token := extractMatch(reSessionToken, body); token != "" {
		c.SessionToken = token
	}

	return body, nil
}

// extractMatch returns the first capture group from a regex match, or empty string.
func extractMatch(re *regexp.Regexp, input string) string {
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
