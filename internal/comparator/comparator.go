package comparator

import (
	"github.com/emacampolo/gomparator/internal/platform/json"
	"github.com/google/go-cmp/cmp"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[gomparator] ", 2)

func Compare(hosts <-chan *HostPairResponse, showDiff bool, statusCodeOnly bool) {
	for h := range hosts {

		if h.Left.Error != nil {
			logger.Printf("error %v", h.Left.Error)
			continue
		}

		if h.Right.Error != nil {
			logger.Printf("error %v", h.Right.Error)
			continue
		}

		if h.Left.StatusCode == h.Right.StatusCode && statusCodeOnly {
			logger.Printf("ok status code %d url %s?%s", h.Left.StatusCode, h.Left.URL.Path, h.Left.URL.RawQuery)
		} else if h.Left.StatusCode == h.Right.StatusCode {
			if j1, j2 := unmarshal(h.Left), unmarshal(h.Right); j1 == nil || j2 == nil {
				continue
			} else if json.Equal(j1, j2) {
				logger.Println("ok")
			} else {
				if showDiff {
					logger.Printf("nok json diff url %s?%s \n%s", h.Left.URL.Path, h.Left.URL.RawQuery, cmp.Diff(j1, j2))
				} else {
					logger.Printf("nok json diff url %s?%s", h.Left.URL.Path, h.Left.URL.RawQuery)
				}
			}
		} else {
			logger.Printf("nok status code url %s?%s, %s: %d - %s: %d",
				h.Left.URL.Path, h.Left.URL.RawQuery, h.Left.URL.Host, h.Left.StatusCode, h.Right.URL.Host, h.Right.StatusCode)
		}
	}
}

func unmarshal(h *Host) interface{} {
	j, err := json.Unmarshal(h.Body)
	if err != nil {
		logger.Printf("nok error unmarshaling from %s with error %v", h.URL.String(), err)
		return nil
	}

	return j
}
