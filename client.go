package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	host             = "api.anthropic.com"
	endpoint         = "v1/complete"
	apiKeyHeader     = "X-Api-Key"
	apiVersionHeader = "Anthropic-Version"
	defaultVersion   = "2023-06-01"
)

// Client is a client for the Anthropic API.
type Client struct {
	key, version string
}

// NewClient returns a client with the given API key.
func NewClient(key string) *Client {
	return &Client{key: key, version: defaultVersion}
}

// SetVersion set's the value passed in the |Anthropic-Version| header for requests.
// The default value is "2023-06-01".
func (c *Client) SetVersion(version string) {
	c.version = version
}

// NewCompletion returns a completion response from the API.
func (c *Client) NewCompletion(ctx context.Context, req *Request) (*Response, error) {
	var b, err = c.post(ctx, endpoint, req)
	if err != nil {
		return nil, err
	}

	var resp = &Response{}
	if err = json.Unmarshal(b, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// TODO: Implement Streaming Response.

func (c *Client) post(ctx context.Context, path string, payload any) ([]byte, error) {
	var b, err = json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var u = url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}

	var req *http.Request
	req, err = c.newRequest(ctx, "POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = interpretResponse(resp); err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) newRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	var req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(apiVersionHeader, c.version)
	req.Header.Set(apiKeyHeader, c.key)

	return req, nil
}

func interpretResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var b, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("code: %d, unable to read response body", resp.StatusCode)
		}
		return fmt.Errorf("code: %d, error: %s", resp.StatusCode, string(b))
	}

	return nil
}
