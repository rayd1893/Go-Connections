package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func New(base string) (*Client, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	client := &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: u,
	}

	return client, nil
}

func (c *Client) NewRequest(ctx context.Context, method, pth string, body interface{}) (*http.Request, error) {
	pth = strings.TrimPrefix(pth, "/")
	u := c.baseURL.ResolveReference(&url.URL{Path: pth})

	var b bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&b).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) Do(req *http.Request, out interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if out == nil {
		return resp, nil
	}

	errPrefix := fmt.Sprintf("%s %s - %d", strings.ToUpper(req.Method), req.URL.String(), resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read body: %w", errPrefix, err)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return nil, fmt.Errorf("%s: response content-type is not application/json (got %s): body: %s",
			errPrefix, ct, body)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return nil, fmt.Errorf("%s: failed to decode JSON response: %w: body: %s",
			errPrefix, err, body)
	}
	return resp, nil
}

func (c *Client) DoWithStatus(req *http.Request, status int, out interface{}) error {
	resp, err := c.Do(req, out)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == status {
		return fmt.Errorf("expected %d response, got %d", status, resp.StatusCode)
	}
	return nil
}
