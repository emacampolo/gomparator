package main

import (
	"bufio"
	"fmt"
	"github.com/emacampolo/gomparator/internal/fetcher"
	"go.uber.org/ratelimit"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Gomparator"
	app.Usage = "Compares API responses by status code and response body"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "the named file for reading URLs",
		},
		cli.StringSliceFlag{
			Name:  "host",
			Usage: "exactly 2 hosts must be specified",
		},
		cli.StringFlag{
			Name:  "headers",
			Usage: "headers separated by commas. eg: \"X-Auth-token: 0x123, X-Public: false\"",
		},
		cli.IntFlag{
			Name:  "ratelimit, r",
			Value: 25,
			Usage: "operation rate limit per second",
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
	defer cl(file)

	hosts := c.StringSlice("host")
	if len(hosts) != 2 {
		log.Fatal("invalid number of hosts provided")
	}

	scanner := bufio.NewScanner(file)
	rl := ratelimit.New(c.Int("ratelimit"))
	lineNumber := 1
	f := fetcher.New()
	headers := parseHeaders(c)

	for scanner.Scan() {
		rl.Take()
		relUrl := scanner.Text()

		var responses []*fetcher.Response
		for _, host := range hosts {
			response, err := f.Fetch(host, relUrl, headers)
			if err != nil {
				log.Println(fmt.Sprintf("line: %d, host: %s, path: %s", lineNumber, host, relUrl), err)
				break
			}
			responses = append(responses, response)
		}

		if len(responses) != 2 {
			continue
		}

		first := responses[0]
		second := responses[1]

		if first.IsOk() && second.IsOk() && cmp.Equal(first.JSON, second.JSON) {
			log.Println(fmt.Sprintf("line %d ok status code %d", lineNumber, 200))
		} else if first.IsOk() && second.IsOk() {
			log.Println(fmt.Sprintf("line %d nok json diff url %s", lineNumber, relUrl), cmp.Diff(first.JSON, second.JSON))
		} else if first.StatusCode == second.StatusCode {
			log.Println(fmt.Sprintf("line %d ok status code %d url %s", lineNumber, first.StatusCode, relUrl))
		} else {
			log.Println(fmt.Sprintf("line %d nok status code url %s, %s: %d - %s: %d",
				lineNumber, relUrl, first.URL.Host, first.StatusCode, second.URL.Host, second.StatusCode))
		}

		lineNumber++
	}
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

func cl(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
