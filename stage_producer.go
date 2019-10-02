package main

import (
	"net/url"
	"sync"

	"go.uber.org/ratelimit"
)

type Fetcher interface {
	Fetch(url string, headers map[string]string) (*Response, error)
}

type HostsPair struct {
	RelURL      string
	Errors      []error
	Left, Right Host
}

func (h HostsPair) EqualStatusCode() bool {
	return h.Left.StatusCode == h.Right.StatusCode
}

func (h HostsPair) HasErrors() bool {
	return len(h.Errors) > 0
}

type Host struct {
	StatusCode int
	Body       []byte
	URL        *url.URL
	Error      error
}

type producer struct {
	concurrency int
	headers     map[string]string
	limiter     ratelimit.Limiter
	fetcher     Fetcher
}

func (p *producer) Produce(in <-chan URLPair) <-chan HostsPair {
	stream := make(chan HostsPair)
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

func NewProducer(concurrency int, headers map[string]string, limiter ratelimit.Limiter, fetcher Fetcher) Producer {
	return &producer{
		concurrency: concurrency,
		headers:     headers,
		limiter:     limiter,
		fetcher:     fetcher,
	}
}

func (p *producer) produce(u URLPair) HostsPair {
	work := func(u URL) <-chan Host {
		ch := make(chan Host, 1)
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

	response := HostsPair{
		RelURL: u.RelURL,
		Left:   lHost,
		Right:  rHost,
	}

	if lHost.Error != nil {
		response.Errors = append(response.Errors, lHost.Error)
	}

	if rHost.Error != nil {
		response.Errors = append(response.Errors, rHost.Error)
	}

	return response
}

func (p *producer) fetch(u URL) Host {
	host := Host{}

	if u.Error != nil {
		host.Error = u.Error
		return host
	}

	response, err := p.fetcher.Fetch(u.URL.String(), p.headers)
	if err != nil {
		host.Error = err
		return host
	}

	host.URL = u.URL
	host.Body = response.Body
	host.StatusCode = response.StatusCode
	return host
}
