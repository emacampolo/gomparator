package comparator

import (
	"context"
	"github.com/emacampolo/gomparator/internal/platform/http"
	"go.uber.org/ratelimit"
	"net/url"
	"sync"
)

type Fetcher interface {
	Fetch(url string, headers map[string]string) (*http.Response, error)
}

type HostPairResponse struct {
	Left, Right *Host
}

type Host struct {
	StatusCode int
	Body       []byte
	URL        *url.URL
	Error      error
}

func NewProducer(ctx context.Context, urls <-chan *URLPairResponse, concurrency int, headers map[string]string,
	limiter ratelimit.Limiter, fetcher Fetcher) <-chan *HostPairResponse {

	ch := make(chan *HostPairResponse)
	go func() {
		var wg sync.WaitGroup

		for w := 0; w < concurrency; w++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					case u := <-urls:
						limiter.Take()
						ch <- &HostPairResponse{
							Left:  fetch(u.Left, fetcher, headers),
							Right: fetch(u.Right, fetcher, headers),
						}
					}
				}
			}()
		}

		wg.Wait()
		close(ch)
	}()

	return ch
}

func fetch(u *URL, fetcher Fetcher, headers map[string]string) *Host {
	host := &Host{}

	if u.Error != nil {
		host.Error = u.Error
	} else if response, err := fetcher.Fetch(u.URL.String(), headers); err == nil {
		host.URL = u.URL
		host.Body = response.Body
		host.StatusCode = response.StatusCode
	} else {
		host.Error = err
	}
	return host
}
