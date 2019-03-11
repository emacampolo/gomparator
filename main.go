package main

import (
	"github.com/emacampolo/gomparator/internal/comparator"
	"github.com/emacampolo/gomparator/internal/platform/http"
	"github.com/emacampolo/gomparator/internal/platform/io"
	"github.com/urfave/cli"
	"go.uber.org/ratelimit"
	"strings"

	"log"
	"os"
	"sync"
)

func main() {
	app := cli.NewApp()
	app.Name = "Gomparator"
	app.Usage = "Compares API responses by status code and response body"
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "the named file for reading URL. It should contain one column only with a rel path. eg: /v1/cards?query=123",
		},
		cli.StringSliceFlag{
			Name:  "host",
			Usage: "exactly 2 hosts must be specified. eg: --host 'http://host1.com --host 'http://host2.com'",
		},
		cli.StringFlag{
			Name:  "headers",
			Usage: `headers separated by commas. eg: "X-Auth-Token: token, X-Public: false"`,
		},
		cli.IntFlag{
			Name:  "ratelimit, r",
			Value: 5,
			Usage: "operation rate limit per second",
		},
		cli.IntFlag{
			Name:  "workers, w",
			Value: 1,
			Usage: "number of workers running concurrently",
		},
		cli.BoolFlag{
			Name:  "show-diff",
			Usage: "whether or not it shows differences when comparison fails",
		},
		cli.BoolFlag{
			Name:  "status-code-only",
			Usage: "whether or not it only compares status code ignoring response body",
		},
	}

	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(c *cli.Context) {
	file, err := os.Open(c.String("path"))
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close(file)

	hosts := c.StringSlice("host")
	if len(hosts) != 2 {
		log.Fatal("invalid number of hosts provided")
	}

	headers := parseHeaders(c)
	lines := io.ReadFile(file)
	comp := comparator.New(http.New(), ratelimit.New(c.Int("ratelimit")))

	var wg sync.WaitGroup
	for w := 0; w < c.Int("workers"); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			comp.Compare(hosts, headers, lines, c.Bool("show-diff"), c.Bool("status-code-only"))
		}()
	}

	wg.Wait()
}

func parseHeaders(c *cli.Context) map[string]string {
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
