package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// DefaultTimeout is the default amount of time for a request before it times out.
const DefaultTimeout = 30 * time.Second

// DefaultMaxBody is the default max number of bytes to be read from response bodies.
// Defaults to no limit.
const DefaultMaxBody = int64(-1)

type Response struct {
	Body       []byte
	StatusCode int
}

type Client struct {
	httpClient      *http.Client
	retryableClient *retryablehttp.Client
	maxBody         int64
}

func NewHTTPClient(opts ...func(*Client)) *Client {
	c := Client{}

	c.retryableClient = retryablehttp.NewClient()
	c.retryableClient.Logger = nil
	c.httpClient = c.retryableClient.HTTPClient
	c.maxBody = DefaultMaxBody

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

// Timeout returns a functional option which sets the maximum amount of time for a request before it times out.
func Timeout(d time.Duration) func(*Client) {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// MaxBody returns a functional option which limits the max number of bytes
// read from response bodies. Set to -1 to disable any limits.
func MaxBody(n int64) func(*Client) {
	return func(a *Client) { a.maxBody = n }
}

func (c *Client) Fetch(url string, headers map[string]string) (*Response, error) {
	res := Response{}

	resp, err := c.get(url, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res.StatusCode = resp.StatusCode

	body := io.Reader(resp.Body)
	if c.maxBody >= 0 {
		body = io.LimitReader(resp.Body, c.maxBody)
	}

	res.Body, err = ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) get(url string, headers map[string]string) (*http.Response, error) {
	req, _ := retryablehttp.NewRequest("GET", url, nil)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.retryableClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
