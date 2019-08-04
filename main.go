package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ecampolo/gomparator/internal/pipeline"
	"github.com/ecampolo/gomparator/internal/platform/http"
	"github.com/ecampolo/gomparator/internal/stages"
	"github.com/urfave/cli"
	"go.uber.org/ratelimit"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
		FullTimestamp:   true,
	})

	log.SetOutput(os.Stdout)
}

func main() {
	app := cli.NewApp()
	app.Name = "Gomparator"
	app.Usage = "Compares API responses by status code and response body"
	app.Version = "1.1"

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
			Name:  "status-code-only",
			Usage: "whether or not it only compares status code ignoring response body",
		},
		cli.DurationFlag{
			Name:  "timeout",
			Value: http.DefaultTimeout,
			Usage: "requests timeout",
		},
		cli.DurationFlag{
			Name:  "duration",
			Value: 0,
			Usage: "duration of the comparision [0 = forever]",
		},
	}

	app.Action = Action
	_ = app.Run(os.Args)
}

type options struct {
	filePath       string
	hosts          []string
	headers        string
	timeout        time.Duration
	duration       time.Duration
	connections    int
	workers        int
	rateLimit      int
	statusCodeOnly bool
}

func Action(cli *cli.Context) {
	opts := parseFlags(cli)
	headers := http.ParseHeaders(opts.headers)
	fetcher := http.New(http.Timeout(opts.timeout))

	ctx, cancel := createContext(opts)
	defer cancel()

	file := openFile(opts)
	defer cl(file)

	logFile := createTmpFile()
	defer cl(logFile)

	log.Printf("created log temp file in %s", logFile.Name())
	log.SetOutput(logFile)

	lines := getTotalLines(file)
	// Once we count the number of lines that will be used as total for the progress bar we reset
	// the pointer to the beginning of the file since it is much faster than closing and reopening
	_, err := file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	bar := stages.NewProgressBar(lines)
	bar.Start()

	reader := stages.NewReader(file, opts.hosts)
	producer := stages.NewProducer(opts.workers, headers,
		ratelimit.New(opts.rateLimit), fetcher)
	comparator := stages.NewConsumer(opts.statusCodeOnly, bar, log.StandardLogger())
	p := pipeline.New(reader, producer, ctx, comparator)

	p.Run()
	bar.Stop()
}

func createContext(opts *options) (context.Context, context.CancelFunc) {
	var ctx context.Context
	var cancel context.CancelFunc
	t := opts.duration
	if t == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		// The request has a timeout, so create a context that is
		// canceled automatically when the timeout expires.
		ctx, cancel = context.WithTimeout(context.Background(), t)
	}
	return ctx, cancel
}

func openFile(opts *options) *os.File {
	file, err := os.Open(opts.filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func createTmpFile() *os.File {
	now := time.Now()
	logFile, err := ioutil.TempFile("", fmt.Sprintf("gomparator.%s.*.txt", now.Format("20060102")))
	if err != nil {
		log.Fatal(err)
	}
	return logFile
}

func getTotalLines(reader io.Reader) int {
	scanner := bufio.NewScanner(reader)

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanLines)

	// Count the lines.
	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}

func parseFlags(cli *cli.Context) *options {
	opts := &options{}

	if opts.hosts = cli.StringSlice("host"); len(opts.hosts) != 2 {
		log.Fatal("invalid number of hosts provided")
	}

	opts.filePath = cli.String("path")
	opts.headers = cli.String("headers")
	opts.timeout = cli.Duration("timeout")
	opts.connections = cli.Int("connections")
	opts.duration = cli.Duration("duration")
	opts.workers = cli.Int("workers")
	opts.rateLimit = cli.Int("ratelimit")
	opts.statusCodeOnly = cli.Bool("status-code-only")
	return opts
}

func cl(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
