package apiclient

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"time"

	"github.com/shashank-sharma/metadata/internal/logger"
)

type APIClient struct {
	HTTPClient *http.Client
	BaseURL    string
	token      string
	devToken   string
}

func New(baseUrl string, timeout time.Duration) *APIClient {
	// Enforcing the trailing slash for BaseURL
	if !endsWithSlash(baseUrl) {
		baseUrl += "/"
	}

	return &APIClient{
		HTTPClient: &http.Client{Timeout: timeout},
		BaseURL:    baseUrl,
	}
}

func (c *APIClient) SetToken(token string) {
	c.token = token
}

func (c *APIClient) SetDevToken(token string) {
	c.devToken = token
}

func (c *APIClient) NewRequestWithParams(method, route string, params map[string]string, headers map[string]string) (*http.Request, error) {
	reqUrl, err := c.url(route)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if len(params) > 0 {
		query := req.URL.Query()
		for key, value := range params {
			query.Add(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}

	c.logRequest(req)
	return req, nil
}

func (c *APIClient) NewRequest(method, route string, body io.Reader, headers map[string]string) (*http.Request, error) {
	reqUrl, err := c.url(route)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, reqUrl, body)
	if err != nil {
		return nil, err
	}

	// Add the provided headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	c.logRequest(req)

	return req, nil
}

func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}

	if c.devToken != "" {
		req.Header.Set("AuthSyncToken", c.devToken)
	}

	if req.Method == "POST" {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	c.logResponse(resp)

	return resp, nil
}

func (c *APIClient) url(route string) (string, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, route)
	return u.String(), nil
}

func (c *APIClient) logRequest(req *http.Request) {
	dump, err := httputil.DumpRequestOut(req, true)
	dumpString := string(dump)
	if err != nil {
		logger.LogDebug("Error dumping request:", err)
	}

	if len(dumpString) > 200 {
		dumpString = dumpString[:200] + "..."
	}

	if false {
		logger.LogDebug("API request:", dumpString)
	}
}

func (c *APIClient) logResponse(resp *http.Response) {
	originalBody := resp.Body
	defer originalBody.Close()

	bodyBytes, err := io.ReadAll(originalBody)
	if err != nil {
		logger.LogDebug("Error reading response body:", err)
		return
	}

	dump, err := httputil.DumpResponse(resp, false) // false because we handle the response body separately
	if err != nil {
		logger.LogDebug("Error dumping response:", err)
		return
	}

	stringBody := string(bodyBytes)
	if len(stringBody) > 200 {
		stringBody = stringBody[:200] + "..."
	}

	if false {
		logger.Debug.Printf("API response: %s\n", string(dump))
		logger.Debug.Printf("Response body: %s\n", stringBody)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}

func endsWithSlash(url string) bool {
	return url[len(url)-1] == '/'
}
