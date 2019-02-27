package utils

import (
	"github.com/urfave/cli"
	"io"
	"log"
	"strings"
)

func ParseHeaders(c *cli.Context) map[string]string {
	var result map[string]string

	headers := strings.Split(c.String("headers"), ",")
	result = make(map[string]string, len(headers))

	for _, header := range headers {
		if header == "" {
			continue
		}

		h := strings.Split(header, ":")
		if len(h) != 2 {
			log.Fatal("invalid header")
		}

		result[h[0]] = h[1]
	}

	return result
}

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
