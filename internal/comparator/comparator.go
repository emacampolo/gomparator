package comparator

import (
	"fmt"
	"github.com/emacampolo/gomparator/internal/platform/http"
	"github.com/emacampolo/gomparator/internal/platform/json"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/ratelimit"
	"log"
	"os"
	"sync"
)

var logger = log.New(os.Stdout, "[comparator] ", 0)

type Fetcher interface {
	Fetch(host string, relPath string, headers map[string]string) (*http.Response, error)
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

func (comp Comparator) Compare(hosts []string, headers map[string]string, jobs <-chan string, wg *sync.WaitGroup,
	showDiff bool, statusCodeOnly bool) {
	defer wg.Done()

	comp.Limiter.Take()
	for relUrl := range jobs {
		first, err := comp.Fetch(hosts[0], relUrl, headers)
		if err != nil {
			logger.Println(fmt.Sprintf("host: %s, path: %s", hosts[0], relUrl), err)
			continue
		}

		second, err := comp.Fetch(hosts[1], relUrl, headers)
		if err != nil {
			logger.Println(fmt.Sprintf("host: %s, path: %s", hosts[1], relUrl), err)
			continue
		}

		if first.StatusCode == second.StatusCode {
			if statusCodeOnly {
				logger.Println(fmt.Sprintf("ok status code %d url %s", first.StatusCode, relUrl))
			} else {
				compareResponses(first, second, relUrl, showDiff)
			}
		} else {
			logger.Println(fmt.Sprintf("nok status code url %s, %s: %d - %s: %d",
				relUrl, first.URL.Host, first.StatusCode, second.URL.Host, second.StatusCode))
		}
	}
}

func compareResponses(first *http.Response, second *http.Response, relUrl string, showDiff bool) {
	equal, err := json.Equal(first.JSON, second.JSON)
	if err != nil {
		log.Println(fmt.Sprintf("error unmarshaling from %s with error %v", relUrl, err))
	}

	if equal {
		logger.Println("ok")
	} else {
		if showDiff {
			j1, j2, _ := json.Unmarshal(first.JSON, second.JSON)
			logger.Println(fmt.Sprintf("nok json diff url %s", relUrl), cmp.Diff(j1, j2))
		} else {
			logger.Println(fmt.Sprintf("nok json diff url %s", relUrl))
		}
	}
}
