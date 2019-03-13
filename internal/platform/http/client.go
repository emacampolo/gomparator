package http

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	// DefaultTimeout is the default amount of time an Attacker waits for a request
	// before it times out.
	DefaultTimeout = 30 * time.Second
	// DefaultConnections is the default amount of max open idle connections per
	// target host.
	DefaultConnections = 10000
)

type Response struct {
	Body       []byte
	StatusCode int
}

type Client struct {
	dialer *net.Dialer
	client http.Client
}

func New(opts ...func(*Client)) *Client {
	c := &Client{}

	c.dialer = &net.Dialer{
		Timeout:   DefaultTimeout,
		KeepAlive: 30 * time.Second,
	}

	c.client = http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           c.dialer.DialContext,
			ResponseHeaderTimeout: DefaultTimeout,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConnsPerHost:   DefaultConnections,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Connections returns a functional option which sets the number of maximum idle
// open connections per target host.
func Connections(n int) func(*Client) {
	return func(c *Client) {
		tr := c.client.Transport.(*http.Transport)
		tr.MaxIdleConnsPerHost = n
	}
}

// Timeout returns a functional option which sets the maximum amount of time
// an Attacker will wait for a request to be responded to.
func Timeout(d time.Duration) func(*Client) {
	return func(c *Client) {
		tr := c.client.Transport.(*http.Transport)
		tr.ResponseHeaderTimeout = d
		c.dialer.Timeout = d
		tr.DialContext = c.dialer.DialContext
	}
}

func (c *Client) Fetch(url string, headers map[string]string) (*Response, error) {
	resp, err := c.get(url, headers)
	if err != nil {
		return nil, err
	}
	defer cl(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func (c *Client) get(url string, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func cl(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
