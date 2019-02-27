package http

import (
	"github.com/emacampolo/gomparator/internal/utils"
	"io/ioutil"
	"net/http"
)

// Response
type Response struct {
	Body       []byte
	StatusCode int
}

// A Client is an HTTP client.
// It wraps net/http's client and add some methods for making HTTP request easier.
type httpClient struct {
	*http.Client
}

func (c *httpClient) get(url string, headers map[string]string) (*Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer utils.Close(resp.Body)

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Body: p, StatusCode: resp.StatusCode,}, nil
}

var client = &httpClient{&http.Client{}}

func Get(url string, headers map[string]string) (*Response, error) {
	return client.get(url, headers)
}
