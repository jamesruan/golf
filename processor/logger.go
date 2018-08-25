package processor

import (
	"github.com/jamesruan/golf/event"
)

const DefaultQueueCapacity = 1000

type LoggerP struct {
	name   string
	logger event.Logger
}

func NewLoggerP(name string, logger event.Logger) *LoggerP {
	return &LoggerP{
		name:   name,
		logger: logger,
	}
}

func (p LoggerP) Name() string {
	return p.name
}

func (p LoggerP) Process(e *event.Event) {
	p.logger.Log(e)
	return
}
