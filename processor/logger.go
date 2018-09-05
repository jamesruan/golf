package processor

import (
	"github.com/jamesruan/golf/event"
)

type loggerP struct {
	name   string
	logger event.Logger
}

// LoggerP returns a processor that handle event to a logger
func NewLoggerP(name string, logger event.Logger) P {
	go func() {
		for e := range logger.Queue() {
			logger.Log(e)
		}
	}()
	return &loggerP{
		name:   name,
		logger: logger,
	}
}

func (p loggerP) Name() string {
	return p.name
}

func (p loggerP) Process(e *event.Event) {
	p.logger.Enqueue(e)
	return
}
