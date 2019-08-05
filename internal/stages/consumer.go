package stages

import (
	"github.com/Sirupsen/logrus"
	"github.com/ecampolo/gomparator/internal/platform/json"
)

type Consumer struct {
	statusCodeOnly bool
	bar            *ProgressBar
	log            *logrus.Logger
}

func NewConsumer(statusCodeOnly bool, bar *ProgressBar, log *logrus.Logger) *Consumer {
	return &Consumer{
		statusCodeOnly: statusCodeOnly,
		bar:            bar,
		log:            log,
	}
}

func (c *Consumer) Consume(val HostsPair) {
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

	if val.EqualStatusCode() {
		leftJSON, err := unmarshal(val.Left)
		if err != nil {
			c.bar.IncrementError()
			c.log.Errorf("could not unmarshal json: url %s: %v", val.RelURL, err)
			return
		}

		rightJSON, err := unmarshal(val.Right)
		if err != nil {
			c.bar.IncrementError()
			c.log.Errorf("could not unmarshal json: url %s: %v", val.RelURL, err)
			return
		}

		if !json.Equal(leftJSON, rightJSON) {
			c.bar.IncrementError()
			c.log.Warnf("found json diff: url %s", val.RelURL)
			return
		}

		c.bar.IncrementOk()
	} else {
		c.bar.IncrementError()
		c.log.Warnf("found status code diff: url %s, %s: %d - %s: %d",
			val.RelURL, val.Left.URL.Host, val.Left.StatusCode, val.Right.URL.Host, val.Right.StatusCode)
	}
}

func unmarshal(h Host) (interface{}, error) {
	j, err := json.Unmarshal(h.Body)
	if err != nil {
		return nil, err
	}

	return j, nil
}
