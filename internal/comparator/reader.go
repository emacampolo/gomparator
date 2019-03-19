package comparator

import (
	"bufio"
	"github.com/emacampolo/gomparator/internal/platform/http"
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
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan *URLPairResponse)
	scanner := bufio.NewScanner(file)

	go func() {
		for scanner.Scan() {
			text := scanner.Text()

			leftUrl := &URL{}
			leftUrl.URL, leftUrl.Error = http.JoinPath(hosts[0], text)

			rightUrl := &URL{}
			rightUrl.URL, rightUrl.Error = http.JoinPath(hosts[1], text)

			out <- &URLPairResponse{Left: leftUrl, Right: rightUrl}
		}
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
		close(out)
	}()

	return out
}
