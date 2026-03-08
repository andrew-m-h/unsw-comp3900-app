//go:build e2e || integration

package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// API path constants (relative to baseURL).
const (
	APIHealth    = "/api/health"
	APIGuestbook = "/api/guestbook/"
)

type acceptJSONTransport struct {
	base http.RoundTripper
}

func (t *acceptJSONTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json")
	return t.base.RoundTrip(req)
}

// Client sends HTTP requests with standard JSON headers (Accept, Content-Type where needed).
// Used by both e2e and integration tests.
type Client struct {
	BaseURL string
	*http.Client
}

// NewClient returns a client that uses baseURL and injects Accept: application/json via a RoundTripper.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		Client: &http.Client{
			Transport: &acceptJSONTransport{base: http.DefaultTransport},
		},
	}
}

// GetJSON performs a GET request to baseURL+path, expects 200 OK, and decodes the response body into v.
func (c *Client) GetJSON(path string, v any) error {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("GET %s: status %d, want 200", path, resp.StatusCode)
	}
	return DecodeJSON(resp, v)
}

// PostJSON performs a POST request to baseURL+path with body as JSON.
func (c *Client) PostJSON(path string, body any) (*http.Response, error) {
	enc, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(enc))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

// PostJSONExpectCreated performs a POST request, expects 201 Created, and decodes the response body into v.
func (c *Client) PostJSONExpectCreated(path string, body any, v any) error {
	resp, err := c.PostJSON(path, body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		return fmt.Errorf("POST %s: status %d, want 201", path, resp.StatusCode)
	}
	return DecodeJSON(resp, v)
}

// DecodeJSON decodes resp.Body into v and closes the body.
func DecodeJSON(resp *http.Response, v any) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
