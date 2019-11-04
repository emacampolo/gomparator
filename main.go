package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
	app.Version = "1.5"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path",
			Usage: "specifies the file from which to read targets. It should contain one column only with a rel path. eg: /v1/cards?query=123",
		},
		cli.StringSliceFlag{
			Name:  "host",
			Usage: "targeted hosts. Exactly 2 must be specified. eg: -host 'http://host1.com -host 'http://host2.com'",
		},
		cli.StringSliceFlag{
			Name:  "header, H",
			Usage: "headers to be used in the http call",
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
			Value: DefaultTimeout,
			Usage: "request timeout",
		},
		cli.DurationFlag{
			Name:  "duration",
			Value: 0,
			Usage: "duration of the comparision [0 = forever]",
		},
		cli.StringFlag{
			Name:  "exclude",
			Usage: "excludes a value from both json for the specified path. A path is a series of keys separated by a dot or #",
		},
	}

	app.Action = Action
	_ = app.Run(os.Args)
}

type options struct {
	filePath       string
	hosts          []string
	headers        []string
	timeout        time.Duration
	duration       time.Duration
	workers        int
	rateLimit      int
	statusCodeOnly bool
	maxBody        int64
	exclude        string
}

func Action(cli *cli.Context) {
	opts := parseFlags(cli)
	headers := parseHeaders(opts.headers)

	fetcher := NewHTTPClient(
		Timeout(opts.timeout),
		MaxBody(opts.maxBody))

	ctx, cancel := createContext(opts)
	defer cancel()

	file := openFile(opts)
	defer file.Close()

	logFile := createTmpFile()
	defer logFile.Close()

	log.Printf("created log temp file in %s", logFile.Name())
	log.SetOutput(logFile)

	lines := getTotalLines(file)
	// Once we count the number of lines that will be used as total for the progress bar we reset
	// the pointer to the beginning of the file since it is much faster than closing and reopening
	_, err := file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	bar := NewProgressBar(lines)
	bar.Start()

	reader := NewReader(file, opts.hosts)
	producer := NewProducer(opts.workers, headers,
		ratelimit.New(opts.rateLimit), fetcher)
	comparator := NewConsumer(opts.statusCodeOnly, bar, log.StandardLogger(), opts.exclude)
	p := New(reader, producer, ctx, comparator)

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
	opts.headers = cli.StringSlice("header")
	opts.timeout = cli.Duration("timeout")
	opts.duration = cli.Duration("duration")
	opts.workers = cli.Int("workers")
	opts.rateLimit = cli.Int("ratelimit")
	opts.statusCodeOnly = cli.Bool("status-code-only")
	if opts.statusCodeOnly {
		opts.maxBody = 0
	} else {
		opts.maxBody = DefaultMaxBody
	}
	opts.exclude = cli.String("exclude")
	return opts
}

func parseHeaders(h []string) map[string]string {
	result := make(map[string]string, len(h))

	for _, header := range h {
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
