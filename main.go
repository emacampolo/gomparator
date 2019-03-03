package main

import (
	"github.com/emacampolo/gomparator/internal/comparator"
	"github.com/emacampolo/gomparator/internal/fetcher"
	"github.com/emacampolo/gomparator/internal/utils"
	"github.com/urfave/cli"
	"go.uber.org/ratelimit"

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

	headers := utils.ParseHeaders(c)
	lines := utils.ReadFile(file)
	comp := comparator.New(fetcher.New(), ratelimit.New(c.Int("ratelimit")))

	wg := new(sync.WaitGroup)
	for w := 0; w < c.Int("workers"); w++ {
		wg.Add(1)
		go comp.Compare(hosts, headers, lines, wg, c.Bool("show-diff"))
	}

	wg.Wait()
}
