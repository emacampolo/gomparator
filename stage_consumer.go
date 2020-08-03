package main

import (
	"github.com/sirupsen/logrus"
)

type consumer struct {
	statusCodeOnly bool
	bar            *ProgressBar
	log            *logrus.Logger
	exclude        string
}

func NewConsumer(statusCodeOnly bool, bar *ProgressBar, log *logrus.Logger, exclude string) Consumer {
	return &consumer{
		statusCodeOnly: statusCodeOnly,
		bar:            bar,
		log:            log,
		exclude:        exclude,
	}
}

func (c *consumer) Consume(val HostsPair) {
	if val.HasErrors() {
		c.bar.IncrementError()
		for _, v := range val.Errors {
			c.log.Errorln(v)
		}
		return
	}

	if val.EqualStatusCode() && c.statusCodeOnly {
		c.bar.IncrementOk()
		return
	}

	if !val.EqualStatusCode() {
		c.bar.IncrementError()
		c.log.Warnf("found status code diff: url %s, %s: %d - %s: %d",
			val.RelURL, val.Left.URL.Host, val.Left.StatusCode, val.Right.URL.Host, val.Right.StatusCode)
		return
	}

	leftJSON, err := unmarshal(val.Left.Body)
	if err != nil {
		c.bar.IncrementError()
		c.log.Errorf("could not unmarshal json: url %s: %v", val.RelURL, err)
		return
	}

	rightJSON, err := unmarshal(val.Right.Body)
	if err != nil {
		c.bar.IncrementError()
		c.log.Errorf("could not unmarshal json: url %s: %v", val.RelURL, err)
		return
	}

	if c.exclude != "" {
		Remove(leftJSON, c.exclude)
		Remove(rightJSON, c.exclude)
	}

	if !Equal(leftJSON, rightJSON) {
		c.bar.IncrementError()
		c.log.Warnf("found json diff: url %s", val.RelURL)
		return
	}

	c.bar.IncrementOk()
}

func unmarshal(b []byte) (interface{}, error) {
	j, err := Unmarshal(b)
	if err != nil {
		return nil, err
	}

	return j, nil
}
