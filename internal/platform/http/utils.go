package http

import (
	"log"
	"net/url"
	"strings"
)

func JoinPath(host string, relPath string) (*url.URL, error) {
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

func ParseHeaders(headers string) map[string]string {
	var result map[string]string

	h := strings.Split(headers, ",")
	result = make(map[string]string, len(h))

	for _, header := range h {
		if header == "" {
			continue
		}

		h := strings.Split(header, ":")
		if len(h) != 2 {
			log.Fatal("invalid header")
		}

		result[h[0]] = h[1]
	}

	return result
}
