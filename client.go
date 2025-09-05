// client.go
package rolestrategy

import (
	"net/http"
)

// Client is the Jenkins Role Strategy SDK client
type Client struct {
	BaseURL    string
	Username   string
	APIToken   string
	HTTPClient *http.Client
}

// NewClient creates a new RoleStrategy client
func NewClient(baseURL, username, apiToken string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8080/jenkins"
	}
	// Ensure baseURL ends with '/'
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}
	return &Client{
		BaseURL:    baseURL,
		Username:   username,
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}

// newRequest creates a new HTTP request with Basic Auth
func (c *Client) newRequest(method, path string, body string) (*http.Request, error) {
	u := c.BaseURL + path
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.APIToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if body != "" {
		req.ContentLength = int64(len(body))
		req.Body = nopReadCloser{str: body}
	}

	return req, nil
}

// nopReadCloser implements io.ReadCloser without Close doing anything
type nopReadCloser struct{ str string }

func (nopReadCloser) Close() error { return nil }
func (r nopReadCloser) Read(b []byte) (int, error) {
	n := copy(b, r.str)
	if n >= len(r.str) {
		r.str = ""
		return n, nil
	}
	r.str = r.str[n:]
	return n, nil
}
