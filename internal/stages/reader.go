package stages

import (
	"bufio"
	"io"
	"net/url"

	"github.com/ecampolo/gomparator/internal/platform/http"
)

type URLPair struct {
	RelURL      string
	Left, Right URL
}

type URL struct {
	URL   *url.URL
	Error error
}

type Reader struct {
	reader io.Reader
	hosts  []string
}

func (r *Reader) Read() <-chan URLPair {
	stream := make(chan URLPair)
	go func() {
		defer close(stream)

		leftHost := r.hosts[0]
		rightHost := r.hosts[1]
		scanner := bufio.NewScanner(r.reader)
		for scanner.Scan() {
			text := scanner.Text()
			leftUrl := URL{}
			leftUrl.URL, leftUrl.Error = http.JoinPath(leftHost, text)

			rightUrl := URL{}
			rightUrl.URL, rightUrl.Error = http.JoinPath(rightHost, text)

			stream <- URLPair{RelURL: text, Left: leftUrl, Right: rightUrl}
		}
	}()

	return stream
}

func NewReader(reader io.Reader, hosts []string) *Reader {
	return &Reader{
		reader: reader,
		hosts:  hosts,
	}
}
