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
		loggerWg.Add(1)
		defer loggerWg.Done()
	loop:
		for {
			select {
			case e := <-logger.Queue():
				logger.Log(e)
			case <-exitSignal:
				break loop
			}
		}
		// make sure the queue is empty
		for {
			select {
			case e := <-logger.Queue():
				logger.Log(e)
			default:
				return
			}
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
