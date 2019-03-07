package http

import (
	"github.com/emacampolo/gomparator/internal/platform/io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var httpClient = &http.Client{}

type Response struct {
	URL        *url.URL
	JSON       []byte
	StatusCode int
}

func New() Client {
	return Client{}
}

type Client struct{}

func (c Client) Fetch(host string, relPath string, headers map[string]string) (*Response, error) {
	url, err := url.Parse(relPath)
	if err != nil {
		return nil, err
	}

	queryString := url.Query()
	url.RawQuery = queryString.Encode()

	base, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	url = base.ResolveReference(url)
	resp, err := c.get(url.String(), headers)
	if err != nil {
		return nil, err
	}
	defer io.Close(resp.Body)

	json, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		URL:        url,
		StatusCode: resp.StatusCode,
		JSON:       json,
	}, nil
}

func (c Client) get(url string, headers map[string]string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
