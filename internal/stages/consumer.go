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

func (c *Consumer) Consume(val *HostsPair) {
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
		leftJSON, lErr := unmarshal(val.Left)
		rightJSON, rErr := unmarshal(val.Right)
		if lErr != nil || rErr != nil {
			c.bar.IncrementError()
			c.log.Errorf("error occurred when unmarshal %s", val.RelURL)
			return
		}

		if !json.Equal(leftJSON, rightJSON) {
			c.bar.IncrementError()
			c.log.Warnf("json diff url %s", val.RelURL)
			return
		}

		c.bar.IncrementOk()
	} else {
		c.bar.IncrementError()
		c.log.Warnf("status code url %s, %s: %d - %s: %d",
			val.RelURL, val.Left.URL.Host, val.Left.StatusCode, val.Right.URL.Host, val.Right.StatusCode)
	}
}

func unmarshal(h *Host) (interface{}, error) {
	j, err := json.Unmarshal(h.Body)
	if err != nil {
		return nil, err
	}

	return j, nil
}
