package fetcher

import (
	"github.com/emacampolo/gomparator/internal/comparator"
	"github.com/emacampolo/gomparator/internal/http"
	"log"
	"net/url"
)

func New() fetcher {
	return fetcher{}
}

type fetcher struct{}

func (fetcher) Fetch(host string, relPath string, headers map[string]string) (*comparator.Response, error) {
	u, err := url.Parse(relPath)
	if err != nil {
		log.Fatal(err)
	}

	queryString := u.Query()
	u.RawQuery = queryString.Encode()

	base, err := url.Parse(host)
	if err != nil {
		log.Fatal(err)
	}

	u = base.ResolveReference(u)
	response, err := http.Get(u.String(), headers)
	if err != nil {
		return nil, err
	}

	return &comparator.Response{
		URL:        u,
		StatusCode: response.StatusCode,
		JSON:       response.Body,
	}, nil
}
