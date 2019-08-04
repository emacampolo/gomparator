package http

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// DefaultTimeout is the default amount of time an Attacker waits for a request
// before it times out.
const DefaultTimeout = 30 * time.Second

type Response struct {
	Body       []byte
	StatusCode int
}

type Client struct {
	dialer *net.Dialer
	client *retryablehttp.Client
}

func New(opts ...func(*Client)) *Client {
	c := Client{}

	c.dialer = &net.Dialer{
		Timeout:   DefaultTimeout,
		KeepAlive: DefaultTimeout,
		DualStack: true,
	}

	c.client = retryablehttp.NewClient()
	c.client.Logger = nil
	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

// Timeout returns a functional option which sets the maximum amount of time
// an Attacker will wait for a request to be responded to.
func Timeout(d time.Duration) func(*Client) {
	return func(c *Client) {
		tr := c.client.HTTPClient.Transport.(*http.Transport)
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
	req, _ := retryablehttp.NewRequest("GET", url, nil)

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
