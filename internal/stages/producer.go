package stages

import (
	"github.com/ecampolo/gomparator/internal/platform/http"
	"go.uber.org/ratelimit"
	"net/url"
	"sync"
)

type Fetcher interface {
	Fetch(url string, headers map[string]string) (*http.Response, error)
}

type HostsPair struct {
	RelURL      string
	Errors      []error
	Left, Right *Host
}

func (h *HostsPair) EqualStatusCode() bool {
	return h.Left.StatusCode == h.Right.StatusCode
}

func (h *HostsPair) HasErrors() bool {
	return len(h.Errors) >= 1
}

type Host struct {
	StatusCode int
	Body       []byte
	URL        *url.URL
	error      error
}

type Producer struct {
	concurrency int
	headers     map[string]string
	limiter     ratelimit.Limiter
	fetcher     Fetcher
}

func (p *Producer) Produce(in <-chan *URLPair) <-chan *HostsPair {
	stream := make(chan *HostsPair)
	go func() {
		defer close(stream)

		var wg sync.WaitGroup
		wg.Add(p.concurrency)

		for w := 0; w < p.concurrency; w++ {
			go func() {
				defer wg.Done()
				for val := range in {
					p.limiter.Take()
					stream <- p.produce(val)
				}
			}()
		}
		wg.Wait()
	}()

	return stream
}

func NewProducer(concurrency int, headers map[string]string, limiter ratelimit.Limiter, fetcher Fetcher) *Producer {
	return &Producer{
		concurrency: concurrency,
		headers:     headers,
		limiter:     limiter,
		fetcher:     fetcher,
	}
}

func (p *Producer) produce(u *URLPair) *HostsPair {
	work := func(u *URL) <-chan *Host {
		ch := make(chan *Host, 1)
		go func() {
			defer close(ch)
			ch <- p.fetch(u)
		}()
		return ch
	}

	leftCh := work(u.Left)
	rightCh := work(u.Right)

	lHost := <-leftCh
	rHost := <-rightCh

	response := &HostsPair{
		RelURL: u.RelURL,
		Left:   lHost,
		Right:  rHost,
	}

	if lHost.error != nil {
		response.Errors = append(response.Errors, lHost.error)
	}

	if rHost.error != nil {
		response.Errors = append(response.Errors, rHost.error)
	}

	return response
}

func (p *Producer) fetch(u *URL) *Host {
	host := &Host{}

	if u.Error != nil {
		host.error = u.Error
	} else if response, err := p.fetcher.Fetch(u.URL.String(), p.headers); err == nil {
		host.URL = u.URL
		host.Body = response.Body
		host.StatusCode = response.StatusCode
	} else {
		host.error = err
	}
	return host
}
