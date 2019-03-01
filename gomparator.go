package main

import (
	"github.com/emacampolo/gomparator/internal/fetcher"
	"github.com/emacampolo/gomparator/internal/json"
	"github.com/emacampolo/gomparator/internal/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/urfave/cli"
	"go.uber.org/ratelimit"

	"fmt"
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
			Usage: "whether or not it shows differences when comparision fails",
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
	defer utils.Close(file)

	hosts := c.StringSlice("host")
	if len(hosts) != 2 {
		log.Fatal("invalid number of hosts provided")
	}

	rateLimiter := ratelimit.New(c.Int("ratelimit"))
	fetcher := fetcher.New()
	headers := utils.ParseHeaders(c)
	lines := utils.ReadFile(file)
	wg := new(sync.WaitGroup)

	for w := 0; w < c.Int("workers"); w++ {
		wg.Add(1)
		go doWork(c, fetcher, hosts, headers, lines, wg, rateLimiter)
	}

	wg.Wait()
}

func doWork(c *cli.Context, fetcher fetcher.Fetcher, hosts []string, headers map[string]string, jobs <-chan string, wg *sync.WaitGroup, limiter ratelimit.Limiter) {
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

		if first.IsOk() && second.IsOk() {
			equal := json.Equal(first.JSON, second.JSON)
			if equal {
				log.Println("ok")
			} else {
				j1, j2, err := json.Unmarshal(first.JSON, second.JSON)
				if err != nil {
					log.Fatalf("error unmarshaling from %s with error %v", relUrl, err)
				}

				if c.Bool("show-diff") {
					log.Println(fmt.Sprintf("nok json diff url %s", relUrl), cmp.Diff(j1, j2))
				} else {
					log.Println(fmt.Sprintf("nok json diff url %s", relUrl))
				}
			}
		} else if first.StatusCode == second.StatusCode {
			log.Println(fmt.Sprintf("ok status code %d url %s", first.StatusCode, relUrl))
		} else {
			log.Println(fmt.Sprintf("nok status code url %s, %s: %d - %s: %d",
				relUrl, first.URL.Host, first.StatusCode, second.URL.Host, second.StatusCode))
		}
	}
}
