package worker

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Worker defines a job runner
type Worker interface {
	// Start the worker
	Start() error
	// Stop the worker
	Stop() error
}

// RepeatedTask invokes function repeatedly
// with a specified execution interval
type RepeatedTask struct {
	task     func() error
	interval time.Duration
	ticker   *time.Ticker
	quit     chan struct{}
	logger   logrus.FieldLogger
}

// NewRepeatedTask creates new RepeatedTask of
// function with given execution interval
func NewRepeatedTask(f func() error, interval time.Duration, logger logrus.FieldLogger) *RepeatedTask {
	return &RepeatedTask{
		task:     f,
		interval: interval,
		logger:   logger.WithField("component", "repeatedtask"),
	}
}

// Start invocation of the function with an interval
func (rt *RepeatedTask) Start() {
	rt.ticker = time.NewTicker(rt.interval)
	rt.quit = make(chan struct{})

	go func() {
		for {
			select {
			case <-rt.ticker.C:
				if err := rt.task(); err != nil {
					rt.logger.Error("error on start the invocation of the function with an interval: ", err)
				}
			case <-rt.quit:
				rt.ticker.Stop()
				return
			}
		}
	}()
}

// Stop repeated task execution
func (rt *RepeatedTask) Stop() {
	close(rt.quit)
}
