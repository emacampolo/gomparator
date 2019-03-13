package main

import (
	"context"
	"github.com/emacampolo/gomparator/internal/comparator"
	"github.com/emacampolo/gomparator/internal/platform/http"
	"github.com/urfave/cli"
	"go.uber.org/ratelimit"
	"io"
	"log"
	"os"
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
		cli.DurationFlag{
			Name:  "timeout",
			Value: http.DefaultTimeout,
			Usage: "requests timeout",
		},
		cli.IntFlag{
			Name:  "connections",
			Value: http.DefaultConnections,
			Usage: "max open idle connections per target host",
		},
	}

	app.Action = Action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Action(ctx *cli.Context) {
	file, err := os.Open(ctx.String("path"))
	if err != nil {
		log.Fatal(err)
	}
	defer cl(file)

	hosts := ctx.StringSlice("host")
	if len(hosts) != 2 {
		log.Fatal("invalid number of hosts provided")
	}

	headers := http.ParseHeaders(ctx.String("headers"))
	fetcher := http.New(http.Timeout(ctx.Duration("timeout")), http.Connections(ctx.Int("connections")))

	urls := comparator.NewReader(file, hosts)
	responses := comparator.NewProducer(context.Background(), urls, ctx.Int("workers"), headers,
		ratelimit.New(ctx.Int("ratelimit")), fetcher)
	comparator.Compare(responses, ctx.Bool("show-diff"), ctx.Bool("status-code-only"))
}

func cl(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
