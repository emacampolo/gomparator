package comparator

import (
	"bufio"
	"github.com/emacampolo/gomparator/internal/platform/http"
	"io"
	"net/url"
)

type URLPairResponse struct {
	Left, Right *URL
}

type URL struct {
	URL   *url.URL
	Error error
}

func NewReader(urls io.Reader, hosts []string) <-chan *URLPairResponse {
	out := make(chan *URLPairResponse)
	scanner := bufio.NewScanner(urls)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()

			leftUrl := &URL{}
			leftUrl.URL, leftUrl.Error = http.JoinPath(hosts[0], text)

			rightUrl := &URL{}
			rightUrl.URL, rightUrl.Error = http.JoinPath(hosts[1], text)

			out <- &URLPairResponse{Left: leftUrl, Right: rightUrl}
		}
		close(out)
	}()

	return out
}
