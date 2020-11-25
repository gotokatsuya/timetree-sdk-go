package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// API endpoint base constants
const (
	APIEndpoint = "https://timetreeapis.com"
)

// Client type
type Client struct {
	accessToken string
	httpClient  *http.Client
	endpoint    *url.URL
}

// NewClient returns a newvclient instance.
func NewClient(accessToken string, httpClient *http.Client) (*Client, error) {
	c := &Client{
		accessToken: accessToken,
		httpClient:  httpClient,
	}
	u, err := url.Parse(APIEndpoint)
	if err != nil {
		return nil, err
	}
	c.endpoint = u
	return c, nil
}

// NewClientWithoutAccessToken returns a new pay client instance.
func NewClientWithoutAccessToken(httpClient *http.Client) (*Client, error) {
	return NewClient("", httpClient)
}

// WithHTTPClient function
func (c *Client) WithHTTPClient(httpClient *http.Client) *Client {
	c.httpClient = httpClient
	return c
}

// mergeQuery method
func (c *Client) mergeQuery(path string, q interface{}) (string, error) {
	v := reflect.ValueOf(q)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return path, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return path, err
	}

	qs, err := query.Values(q)
	if err != nil {
		return path, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewRequest method
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	switch method {
	case http.MethodGet, http.MethodDelete:
		if body != nil {
			merged, err := c.mergeQuery(path, body)
			if err != nil {
				return nil, err
			}
			path = merged
		}
	}
	u, err := c.endpoint.Parse(path)
	if err != nil {
		return nil, err
	}

	var reqBody io.ReadWriter
	switch method {
	case http.MethodPost, http.MethodPut:
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewBuffer(b)
		}
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.timetree.v1+json")
	req.Header.Set("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
	}
	return req, nil
}

// Do method
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer resp.Body.Close()

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
				return resp, err
			}
		}
	}
	return resp, err
}

type ErrorResponse struct {
	Type   string `json:"type,omitempty"`
	Status int    `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	// Errors string `json:"errors"`
}
