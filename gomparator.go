package main

import (
	"bufio"
	"fmt"
	"github.com/emacampolo/gomparator/internal/fetcher"
	"go.uber.org/ratelimit"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/urfave/cli"
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
			Usage: "exactly 2 hosts must be specified. eg: --host 'http://hostA --host 'http://hostB'",
		},
		cli.StringFlag{
			Name:  "headers",
			Usage: "headers separated by commas. eg: \"X-Auth-token: 0x123, X-Public: false\"",
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

	limiter := ratelimit.New(c.Int("ratelimit"))
	f := fetcher.New()
	headers := parseHeaders(c)

	jobs := make(chan string)

	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			jobs <- scanner.Text()
		}
		close(jobs)
	}()

	wg := new(sync.WaitGroup)

	for w := 0; w < c.Int("workers"); w++ {
		wg.Add(1)
		go doWork(f, hosts, headers, jobs, wg, limiter)
	}

	wg.Wait()
}

func doWork(fetcher fetcher.Fetcher, hosts []string, headers map[string]string, jobs <-chan string, wg *sync.WaitGroup, limiter ratelimit.Limiter) {
	defer wg.Done()

	limiter.Take()
	for relUrl := range jobs {
		first, err := fetcher.Fetch(hosts[0], relUrl, headers)
		if err != nil {
			log.Println(fmt.Sprintf("host: %s, path: %s", hosts[0], relUrl), err)
			continue
		}

		second, err := fetcher.Fetch(hosts[1], relUrl, headers)
		if err != nil {
			log.Println(fmt.Sprintf("host: %s, path: %s", hosts[1], relUrl), err)
			continue
		}

		if first.IsOk() && second.IsOk() && reflect.DeepEqual(first.JSON, second.JSON) {
			log.Println(fmt.Sprintf("ok status code %d", 200))
		} else if first.IsOk() && second.IsOk() {
			log.Println(fmt.Sprintf("nok json diff url %s", relUrl))
		} else if first.StatusCode == second.StatusCode {
			log.Println(fmt.Sprintf("ok status code %d url %s", first.StatusCode, relUrl))
		} else {
			log.Println(fmt.Sprintf("nok status code url %s, %s: %d - %s: %d",
				relUrl, first.URL.Host, first.StatusCode, second.URL.Host, second.StatusCode))
		}
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
