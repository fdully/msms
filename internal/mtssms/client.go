package mtssms

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

type Provider interface {
	Send(ctx context.Context, payload *MTSRequest) (*MTSResponse, error)
	Info(ctx context.Context, payload *MTSInfoRequest) (*MTSInfoResponse, error)
}

// client provides an API client.
type client struct {
	baseURL    *url.URL
	httpClient *http.Client
	login      string
	password   string
	signuture  string
}

// NewClient creates a new HTTP client with the given base URL.
func New(base, login, password, sign string) (Provider, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	client := &client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL:   u,
		login:     login,
		password:  password,
		signuture: sign,
	}

	return client, nil
}

// Sending messages.
func (c *client) Send(ctx context.Context, payload *MTSRequest) (*MTSResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/http-api/v1/messages", payload)
	if err != nil {
		return nil, err
	}

	var out MTSResponse
	if err := c.doOK(req, &out); err != nil {
		return &out, err
	}
	return &out, nil
}

// get info about messages.
func (c *client) Info(ctx context.Context, payload *MTSInfoRequest) (*MTSInfoResponse, error) {
	req, err := c.newRequest(ctx, "POST", "/http-api/v1/messages/info", payload)
	if err != nil {
		return nil, err
	}

	var out MTSInfoResponse
	if err := c.doOK(req, &out); err != nil {
		return &out, err
	}
	return &out, nil
}

// newRequest creates a new request with the given method, path (relative to the
// baseURL), and optional body. If the body is given, it's encoded as json.
func (c *client) newRequest(ctx context.Context, method, pth string, body interface{}) (*http.Request, error) {
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

	if c.login != "" && c.password != "" {
		req.SetBasicAuth(c.login, c.password)
	}

	return req, nil
}

// doOK is like do, but expects a 200 response.
func (c *client) doOK(req *http.Request, out interface{}) error {
	resp, err := c.do(req, out)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("expected 200 response, got %d", resp.StatusCode)
	}
	return nil
}

// errorResponse is used to extract an error from the response, if it exists.
// This is a fallback for when all else fails.
type errorResponse struct {
	Error1 string `json:"error"`
	Error2 string `json:"Error"`
}

// Error returns the error string, if any.
func (e *errorResponse) Error() string {
	if e.Error1 != "" {
		return e.Error1
	}
	return e.Error2
}

// do executes the request and decodes the result into out. It returns the http
// response. It does NOT do error checking on the response code.
func (c *client) do(req *http.Request, out interface{}) (*http.Response, error) {
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

	var errResp errorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error() != "" {
		return nil, fmt.Errorf("%s: error response from API: %s, err: %w, body: %s",
			errPrefix, errResp.Error(), err, body)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return nil, fmt.Errorf("%s: failed to decode JSON response: %w: body: %s",
			errPrefix, err, body)
	}
	return resp, nil
}
