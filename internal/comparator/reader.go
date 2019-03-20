package comparator

import (
	"bufio"
	"github.com/emacampolo/gomparator/internal/platform/http"
	"io"
	"log"
	"net/url"
	"os"
)

type URLPairResponse struct {
	Left, Right *URL
}

type URL struct {
	URL   *url.URL
	Error error
}

func NewReader(filePath string, hosts []string) <-chan *URLPairResponse {
	out := make(chan *URLPairResponse)
	go func() {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer cl(file)

		scanner := bufio.NewScanner(file)
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

func cl(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
