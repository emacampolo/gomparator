package comparator

import (
	"fmt"
	"github.com/emacampolo/gomparator/internal/json"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/ratelimit"
	"log"
	"net/url"
	"sync"
)

type Fetcher interface {
	Fetch(host string, relPath string, headers map[string]string) (*Response, error)
}

func New(fetcher Fetcher, rateLimiter ratelimit.Limiter) Comparator {
	return Comparator{
		fetcher,
		rateLimiter,
	}
}

type Comparator struct {
	Fetcher
	ratelimit.Limiter
}

type Response struct {
	URL        *url.URL
	JSON       []byte
	StatusCode int
}

func (r Response) IsOk() bool {
	return r.StatusCode == 200
}

func (comp Comparator) Compare(hosts []string, headers map[string]string, jobs <-chan string, wg *sync.WaitGroup,
	showDiff bool, statusCodeOnly bool) {
	defer wg.Done()

	comp.Limiter.Take()
	for relUrl := range jobs {
		first, err := comp.Fetch(hosts[0], relUrl, headers)
		if err != nil {
			log.Println(fmt.Sprintf("host: %s, path: %s", hosts[0], relUrl), err)
			continue
		}

		second, err := comp.Fetch(hosts[1], relUrl, headers)
		if err != nil {
			log.Println(fmt.Sprintf("host: %s, path: %s", hosts[1], relUrl), err)
			continue
		}

		if first.StatusCode == second.StatusCode {
			if statusCodeOnly {
				log.Println(fmt.Sprintf("ok status code %d url %s", first.StatusCode, relUrl))
			} else {
				compareResponses(first, second, relUrl, showDiff)
			}
		} else {
			log.Println(fmt.Sprintf("nok status code url %s, %s: %d - %s: %d",
				relUrl, first.URL.Host, first.StatusCode, second.URL.Host, second.StatusCode))
		}
	}
}

func compareResponses(first *Response, second *Response, relUrl string, showDiff bool) {
	equal := json.Equal(first.JSON, second.JSON)
	if equal {
		log.Println("ok")
	} else {
		if showDiff {
			j1, j2, err := json.Unmarshal(first.JSON, second.JSON)
			if err != nil {
				log.Fatalf("error unmarshaling from %s with error %v", relUrl, err)
			}
			log.Println(fmt.Sprintf("nok json diff url %s", relUrl), cmp.Diff(j1, j2))
		} else {
			log.Println(fmt.Sprintf("nok json diff url %s", relUrl))
		}
	}
}
