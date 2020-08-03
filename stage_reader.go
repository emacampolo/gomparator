package main

import (
	"bufio"
	"io"
	"net/url"
)

type URLPair struct {
	RelURL      string
	Left, Right URL
}

type URL struct {
	URL   *url.URL
	Error error
}

type reader struct {
	reader io.Reader
	hosts  []string
}

func (r *reader) Read() <-chan URLPair {
	stream := make(chan URLPair)
	go func() {
		defer close(stream)

		leftHost := r.hosts[0]
		rightHost := r.hosts[1]
		scanner := bufio.NewScanner(r.reader)
		for scanner.Scan() {
			text := scanner.Text()
			leftURL := URL{}
			leftURL.URL, leftURL.Error = joinPath(leftHost, text)

			rightURL := URL{}
			rightURL.URL, rightURL.Error = joinPath(rightHost, text)

			stream <- URLPair{RelURL: text, Left: leftURL, Right: rightURL}
		}
	}()

	return stream
}

func joinPath(host string, relPath string) (*url.URL, error) {
	u, err := url.Parse(relPath)
	if err != nil {
		return nil, err
	}

	queryString := u.Query()
	u.RawQuery = queryString.Encode()

	base, err := u.Parse(host)
	if err != nil {
		return nil, err
	}

	return base.ResolveReference(u), nil
}

func NewReader(r io.Reader, hosts []string) Reader {
	return &reader{
		reader: r,
		hosts:  hosts,
	}
}
